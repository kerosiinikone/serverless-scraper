package reddit

import (
	"context"
	"testing"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"github.com/kerosiinikone/serverless-scraper/util"
	"github.com/stretchr/testify/assert"
)

// Integration test
func TestScrape(t *testing.T) {
	var list []*scraper.APIRequest

	req := &scraper.APIRequest{
		ID:        "123",
		ClientID:  "123",
		Subreddit: "golang",
	}
	rs := &RedditScraper{
		cfg: &scraper.Config{
			MaxDepth: 2,
			DelayMax: 2,
			DelayMin: 1,
		},
		cb: func(*scraper.APIRequest) error {
			list = append(list, req)
			return nil
		},
		a: util.NewAnalyzer(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	out := make(chan struct{})

	err := rs.Scrape(ctx, req, out)

	select {
	case <-ctx.Done():
		assert.Nil(t, err)
	}
	if len(list) == 0 {
		t.Error("Request list is empty, scraper failed")
	}
}

// Unit tests
func TestDecrementPageCount(t *testing.T) {
    pageCount = 6
    decrementPageCount()
    assert.Equal(t, 5, pageCount)
    decrementPageCount()
    assert.Equal(t, 4, pageCount)
}

func TestFetchHttpResponse(t *testing.T) {
	var (
		headers = map[string]string{}
		link    = "https://old.reddit.com/r/golang/comments/1hxw1yf/how_many_bottles_of_water_have_you_drunk_today/"	
	)
	resp, err := fetchHttpResponse(headers, link)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("Response is nil")
	}
}

func TestProcessAndDispatchPost(t *testing.T) {
	rs := New(&scraper.Config{
		MaxDepth: 2,
		DelayMax: 2,
		DelayMin: 1,
	}, nil, nil)

	engine, err := actor.NewEngine(actor.NewEngineConfig())
		if err != nil {
			t.Error(err)
		}
    managerPID := engine.Spawn(NewManager(nil, nil, rs.s), "manager")

	post := models.RedditPostDetails{
		Id: "1hxw1yf",
		Subreddit: "golang",
		Title: "How Many Bottles of Water Have You Drunk Today?",
	}

	err = rs.processAndDispatchPost(post, managerPID, engine)

	assert.Nil(t, err)
}

func TestRequestAndPipePost(t *testing.T) {
	rs := New(&scraper.Config{
		MaxDepth: 2,
		DelayMax: 2,
		DelayMin: 1,
	}, nil, nil)

	pipe := make(chan models.RedditPostDetails)
	apiURL := "https://old.reddit.com/r/golang.json?limit=10"

	err := rs.requestAndPipePost(apiURL, "", pipe)

	assert.Nil(t, err)
}

