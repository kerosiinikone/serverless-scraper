package util

import (
	"context"
	"testing"
	"time"

	"github.com/kerosiinikone/serverless-scraper/pkg/models"
	"golang.org/x/exp/rand"
)

func TestFilter(t *testing.T) {
	t.Run("Test NewAnalyzer", func(t *testing.T) {
		a := NewAnalyzer()
		if a == nil {
			t.Error("NewAnalyzer returned nil")
		}
	})
	t.Run("Test FilterPosts", func(t *testing.T) {
		a := NewAnalyzer()
		randomKeyword := Keywords[rand.Intn(len(Keywords))]
		posts := []models.PostContainer{
			{
				Post: models.RedditPost{
					Title:    randomKeyword,
					Selftext: "This is a very bad post about a very bad thing that is very bad",
				},
			},
		}

		out := make(chan models.RedditPostDetails)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		go func() {
			err := a.FilterPosts(posts, out)
			if err != nil {
				t.Errorf("FilterPosts returned error: %v", err)
			}
		}()

		select {
			case p := <-out: 
				if p.Title != randomKeyword {
					t.Error("FilterPosts returned the wrong post")
				}
			case <-ctx.Done():
				t.Error("FilterPosts timed out")
		}
	})
}