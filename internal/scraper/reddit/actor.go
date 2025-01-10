package reddit

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/anthdm/hollywood/actor"
	regexp "github.com/h2so5/goback/regexp"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"github.com/kerosiinikone/serverless-scraper/util"
	"golang.org/x/net/html"
)

type RedditScraperWorker struct {
	Post models.RedditPostDetails
	Mpid *actor.PID
}

func NewActor(post models.RedditPostDetails, mpid *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &RedditScraperWorker{
			Post: post,
			Mpid: mpid,
		}
	}
}

func (s *RedditScraperWorker) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		p, err := s.Scrape(ctx)
		if err != nil {
			s.handleError(ctx, err)
			return
		}
		ctx.Send(s.Mpid, p)
		ctx.Engine().Poison(ctx.PID())
		return
	}
}

func (s *RedditScraperWorker) Scrape(ctx *actor.Context) (models.ForumTree, error) {
	var (
		res *http.Response
	)
	res, err := fetchHttpResponse(getHeaders(s.Post.Subreddit), s.createPostLink())
	if err != nil {
		return models.ForumTree{}, err
	}
	defer res.Body.Close()

	doc, err := processHTTPBody(res)
	if err != nil {
		return models.ForumTree{}, err
	}
	p, err := s.extractPostContent(doc)
	if err != nil {
		return models.ForumTree{}, err
	}
	c, err := s.extractCommentTree(doc)
	if err != nil {
		return models.ForumTree{}, err
	}
	return models.ForumTree{
		Id: p.Id,
		Selftext: p.Selftext,
		Subreddit: p.Subreddit,
		Comments: c,
	}, nil
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

func (s *RedditScraperWorker) handleError(ctx *actor.Context, err error) {
	fmt.Println(err)
	ctx.Send(s.Mpid, models.ForumTree{})
	ctx.Engine().Poison(ctx.PID())
}

func (s *RedditScraperWorker) createPostLink() string {
	t := strings.Join(strings.Split(strings.ToLower(s.Post.Title), " "), "_")
	return fmt.Sprintf("https://old.reddit.com/r/%s/comments/%s/%s", s.Post.Subreddit, s.Post.Id, t)
}	

func (s *RedditScraperWorker) extractPostContent(doc *html.Node) (models.ForumTree, error) {
	var (
		selftext string
	)
	if doc == nil {
		return models.ForumTree{}, fmt.Errorf("No document found")
	}
	n := htmlquery.FindOne(doc, "//div[@data-type='link']")
	if n == nil {
		return models.ForumTree{}, fmt.Errorf("No post found")
	}
	pDiv := htmlquery.FindOne(n, "//div[@class='md']")
	if pDiv == nil {
		return models.ForumTree{}, fmt.Errorf("No selftext found")
	}
	ps := htmlquery.Find(pDiv, "//p")
	for _, p := range ps {
		finalText := sanitizeText(htmlquery.InnerText(p))
		if containsFilteredKeywords(finalText) {
			continue
		}
		selftext += finalText
	}
	if selftext == "" {
		return models.ForumTree{}, fmt.Errorf("No selftext")
	}
	return models.ForumTree{
		Id: s.Post.Id,
		Selftext: selftext,
		Subreddit: s.Post.Subreddit,
		Comments: make([]models.ReplyTree, 0),
	}, nil
}

func (s *RedditScraperWorker) extractCommentTree(doc *html.Node) ([]models.ReplyTree, error) {
	var (
		traverseComment func(n *html.Node) []models.ReplyTree
	)
	traverseComment = func(n *html.Node) []models.ReplyTree {
		var (
			comments []models.ReplyTree
		)
		c := htmlquery.FindOne(n, "//div[@data-type='comment']")
		if c == nil {
			return []models.ReplyTree{}
		}
		comments = append(comments, models.ReplyTree{
			Body: extractCommentText(c),
			Replies: traverseComment(c),
		})
		for cc := c.NextSibling; cc != nil; cc = cc.NextSibling {
			if isComment(cc) {
				comments = append(comments, models.ReplyTree{
					Body: extractCommentText(cc),
					Replies: traverseComment(cc),
				})
			}
		}
		return comments
	}
	return traverseComment(doc), nil
}

func processHTTPBody(res *http.Response) (*html.Node, error) {
	body, err := parseBodyIfGzip(res.Body, res)
	if err != nil {
		return nil, err
	}
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func parseBodyIfGzip(b io.ReadCloser, res *http.Response) ([]byte, error) {
	var (
		body []byte
	)
	if res.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(res.Body)
		if err != nil {
			return []byte{}, err
		}
		defer gz.Close()
		body, err = io.ReadAll(gz)
		if err != nil {
			return []byte{}, err
		}
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return []byte{}, err
		}
		return body, nil
	}
	return body, nil
}

func extractCommentText(n *html.Node) string {
	commentDiv := htmlquery.FindOne(n, "//div[@class='md']")
    if commentDiv == nil {
        return ""
    }
    commentBody := htmlquery.FindOne(commentDiv, "//p")
    if commentBody == nil {
        return ""
    }
    return sanitizeText(htmlquery.InnerText(commentBody))
}

func sanitizeText(text string) string {
    re := regexp.MustCompile(`[^\x00-\x7F]`)
    asciiText := re.ReplaceAllString(text, "")
    re = regexp.MustCompile(`"`)
    return re.ReplaceAllString(asciiText, "'")
}

func containsFilteredKeywords(text string) bool {
    for _, keyword := range util.Filter {
        if strings.Contains(text, keyword) {
            return true
        }
    }
    return false
}

func isComment(n *html.Node) bool {
    for _, attr := range n.Attr {
        if attr.Key == "data-type" && attr.Val == "comment" {
            return true
        }
    }
    return false
}
