package reddit

import (
	"context"
	"testing"

	"github.com/anthdm/hollywood/actor"
	"github.com/kerosiinikone/serverless-scraper/internal/scraper"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestDecrementPageCount(t *testing.T) {
    pageCount = 6
    count := decrementPageCount()
    assert.Equal(t, 5, count)
    count = decrementPageCount()
    assert.Equal(t, 4, count)
}

func TestScrape(t *testing.T) {
	req := &scraper.APIRequest{
		ID:        "123",
		ClientID:  "123",
		Subreddit: "golang",
	}
	rs := New(&scraper.Config{
		MaxDepth: 10,
		DelayMax: 10,
		DelayMin: 5,
	}, nil, nil)

	ctx := context.Background()
	out := make(chan struct{})

	err := rs.Scrape(ctx, req, out)
	
	assert.Nil(t, err)
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
		Id: "123",
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
	apiURL := "https://www.reddit.com/r/golang.json?limit=10"

	err := rs.requestAndPipePost(apiURL, "", pipe)

	assert.Nil(t, err)
}

