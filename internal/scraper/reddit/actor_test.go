package reddit

import (
	"testing"

	"github.com/kerosiinikone/serverless-scraper/pkg/models"
)

// Integration test
func TestScrape_Actor(t *testing.T) {
	s := &RedditScraperWorker{
		Post: models.RedditPostDetails{
			Id: "1hxw1yf",
			Subreddit: "golang",
			Title: "How Many Bottles of Water Have You Drunk Today?",
		},
	}

	tree, err := s.Scrape(nil)

	if err != nil {
		t.Error(err)
	}
	if tree.Id != s.Post.Id {
		t.Error("Id mismatch")
	}
	if tree.Subreddit != "golang" {
		t.Error("Subreddit mismatch")
	}
}

// Unit tests
func TestProcessHTTPBody(t *testing.T) {
	var (
		headers = map[string]string{}
		link    = "https://old.reddit.com/r/golang/comments/1hxw1yf/how_many_bottles_of_water_have_you_drunk_today/"
	)
	resp, err := fetchHttpResponse(headers, link)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	doc, err := processHTTPBody(resp)
	if err != nil {
		t.Error(err)
	}
	if doc == nil {
		t.Error("Document is nil")
	}
}

func TestExtractPostContent(t *testing.T) {
	s := &RedditScraperWorker{
		Post: models.RedditPostDetails{
			Id: "1hxw1yf",
			Subreddit: "golang",
			Title: "How Many Bottles of Water Have You Drunk Today?",
		},
	}
	res, err := fetchHttpResponse(map[string]string{}, s.createPostLink())
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	doc, err := processHTTPBody(res)
	if err != nil {
		t.Error(err)
	}
	p, err := s.extractPostContent(doc)
	if err != nil {
		t.Error(err)
	}
	if p.Id != s.Post.Id {
		t.Error("Id mismatch")
	}
	if p.Subreddit != s.Post.Subreddit {
		t.Error("Subreddit mismatch")
	}
}

func TestExtractPostCommentTree(t *testing.T) {
	s := &RedditScraperWorker{
		Post: models.RedditPostDetails{
			Id: "1hxw1yf",
			Subreddit: "golang",
			Title: "How Many Bottles of Water Have You Drunk Today?",
		},
	}
	res, err := fetchHttpResponse(map[string]string{}, s.createPostLink())
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	doc, err := processHTTPBody(res)
	if err != nil {
		t.Error(err)
	}
	_, err = s.extractCommentTree(doc)
	if err != nil {
		t.Error(err)
	}
}

// ...