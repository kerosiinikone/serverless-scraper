package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
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
}

var (
    pageCount = 6
    initialBackoff = time.Second
	mu sync.Mutex
)

func New(cfg *scraper.Config, storage *s3manager.Uploader, queue *sqs.SQS) RedditScraper {
	return RedditScraper{
		cfg: cfg,
		s: storage,
		q: queue,
        a: util.NewAnalyzer(),
	}
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
        }
    }()

    return rs.acceptLoop(ctx, pipe, managerPID, engine, finished, request, out)
}

func (rs *RedditScraper) acceptLoop(ctx context.Context, pipe chan models.RedditPostDetails, managerPID *actor.PID, engine *actor.Engine, finished <-chan struct{}, request *scraper.APIRequest, out chan struct{}) error {
    for {
        select {
        case post := <-pipe:
            if err := rs.processAndDispatchPost(post, managerPID, engine) ; err != nil {
                // Continue to next post if processing fails -> do not return
                log.Println("Failed to process post: ", err)
            }
        case <-finished:
            err := rs.sendQueueMessage(request)
            if err != nil {
                return fmt.Errorf("Failed to send queue message: %w", err)
            }
            close(out)
            return nil
        case <-ctx.Done():
            err := rs.sendQueueMessage(request)
            if err != nil {
                return fmt.Errorf("Failed to send queue message: %w", err)
            }
            close(out)
            return nil
        }
    }
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

    if decrementPageCount() == 0 {
        return nil
    }
    resp, err := fetchHttpResponse(map[string]string{}, url)
    if err != nil {
		return err
	}
    defer resp.Body.Close()
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }
    go func() {
        if err := rs.requestAndPipePost(url, response.Data.After, out); err != nil {
            log.Println("Error requesting and piping: ", err)
        }
    }()
    if err := rs.a.FilterPosts(response.Data.Children, out); err != nil {
        return fmt.Errorf("failed to filter posts: %w", err)
    }
    
    return nil
}

func (rs *RedditScraper) sendQueueMessage(request *scraper.APIRequest) error {
    return queue.SendMessage(rs.q, &models.QueueMessage{
        ClientID:  request.ClientID,
        RequestID: request.ID,
    })
}

func decrementPageCount() int {
    mu.Lock()
    defer mu.Unlock()
    pageCount = int(math.Max(0, float64(pageCount-1)))
    return pageCount
}
