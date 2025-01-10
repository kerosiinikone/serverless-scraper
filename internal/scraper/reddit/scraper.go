package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	regexp "github.com/h2so5/goback/regexp"
	"github.com/kerosiinikone/serverless-scraper/infra/queue"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"github.com/kerosiinikone/serverless-scraper/util"
	"golang.org/x/exp/rand"
)

type RedditScraper struct {
	cfg *scraper.Config
	s *s3manager.Uploader
    q *sqs.SQS
    a *util.Analyzer
    cb func(*scraper.APIRequest) error
}

var (
    pageCount = 1
    initialBackoff = time.Second
	mu sync.Mutex
)

func New(cfg *scraper.Config, storage *s3manager.Uploader, queue *sqs.SQS) RedditScraper {
	rs := RedditScraper{
		cfg: cfg,
		s: storage,
		q: queue,
        a: util.NewAnalyzer(),
	}
    rs.cb = rs.sendQueueMessage
    return rs
}

func (rs *RedditScraper) Scrape(ctx context.Context, request *scraper.APIRequest, out chan struct{}) error {
    var (
        finished = make(chan struct{})
        pipe = make(chan models.RedditPostDetails)
        apiURL = fmt.Sprintf("https://www.reddit.com/r/%s.json?limit=%d", request.Subreddit, rs.cfg.MaxDepth)
    )

    engine, err := actor.NewEngine(actor.NewEngineConfig())
    if err != nil {
        return fmt.Errorf("failed to create actor engine: %w", err)
    }
    managerPID := engine.Spawn(NewManager(finished, request, rs.s), "manager")

    go func() {
        if err := rs.requestAndPipePost(apiURL, "", pipe); err != nil {
            log.Println("Error requesting and piping: ", err)
            close(out)
        }
    }()

    return rs.acceptLoop(ctx, pipe, managerPID, engine, finished, request, out)
}

func (rs *RedditScraper) acceptLoop(ctx context.Context, pipe chan models.RedditPostDetails, managerPID *actor.PID, engine *actor.Engine, finished <-chan struct{}, request *scraper.APIRequest, out chan struct{}) error {
    for {
        select {
        case post := <-pipe:
            log.Println("Processed post: ", post)
            if err := rs.processAndDispatchPost(post, managerPID, engine) ; err != nil {
                log.Println("Failed to process post: ", err)
            }
        case <-finished:
            return rs.Close(out, request)
        case <-ctx.Done():
            return rs.Close(out, request)
        }
    }
}

func (rs *RedditScraper) Close(out chan struct{}, r *scraper.APIRequest) error {
    err := rs.cb(r)
    if err != nil {
        return err
    }
    close(out)
    return nil
}

func (rs *RedditScraper) processAndDispatchPost(post models.RedditPostDetails, managerPID *actor.PID, engine *actor.Engine) error {
    reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+|\\n")
    if err != nil {
        return fmt.Errorf("failed to compile regex: %w", err)
    }

    engine.Send(managerPID, models.RedditPostDetails{
        Id:        post.Id,
        Title:     reg.ReplaceAllString(post.Title, ""),
        Subreddit: post.Subreddit,
    })

    delay := time.Duration(rand.Intn(rs.cfg.DelayMax-rs.cfg.DelayMin) + rs.cfg.DelayMin)
    time.Sleep(time.Second * delay)

    return nil
}

func (rs *RedditScraper) requestAndPipePost(url string, after string, out chan<- models.RedditPostDetails) error {
    var response models.RedditPostResponse

    reqURL := url
    if after != "" {
        reqURL = fmt.Sprintf("%s&after=%s", url, after)
    }
    resp, err := fetchHttpResponse(map[string]string{}, reqURL)
    if err != nil {
		return err
	}
    defer resp.Body.Close()
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }

    if pageCount > 1 {
        go func() {
            if err := rs.requestAndPipePost(url, response.Data.After, out); err != nil {
                log.Println("Error requesting and piping: ", err)
            }
        }()
    }
    if err := rs.a.FilterPosts(response.Data.Children, out); err != nil {
        return fmt.Errorf("failed to filter posts: %w", err)
    }

    decrementPageCount()
    
    return nil
}

func fetchHttpResponse(headers Headers, url string) (*http.Response, error) {
	var res *http.Response
	
	proxy, err := util.Proxy()
	if err != nil {
		return nil, err
	}
	bc := util.NewBackoffCaller(headers, initialBackoff, proxy)
	res, err = bc.Call(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (rs *RedditScraper) sendQueueMessage(request *scraper.APIRequest) error {
    return queue.SendMessage(rs.q, &models.QueueMessage{
        ClientID:  request.ClientID,
        RequestID: request.ID,
    })
}

func decrementPageCount() {
    mu.Lock()
    defer mu.Unlock()
    pageCount = int(math.Max(0, float64(pageCount-1)))
}
