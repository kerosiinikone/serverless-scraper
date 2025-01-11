package util

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/derekparker/trie"
	"github.com/drankou/go-vader/vader"
	"github.com/kerosiinikone/serverless-scraper/pkg/models"
)

//go:embed data/vader_lexicon.txt
var vaderLexicon []byte 

var maxDatapoints = 80

type Analyzer struct {
	Analyzer *vader.SentimentIntensityAnalyzer
}

func NewAnalyzer() *Analyzer {
	sia := vader.SentimentIntensityAnalyzer{
		LexiconMap: vader.MakeLexiconMap(string(vaderLexicon)),
	}
	return &Analyzer{
		Analyzer: &sia,
	}
} 

// Optimize
func (a *Analyzer) FilterPosts(posts []models.PostContainer, out chan<- models.RedditPostDetails) error {
	t := trie.New()
	for _, keyword := range Keywords {
		t.Add(keyword, nil)
	}
	outerLoop:
	for _, post := range posts {
		p := post.Post
		if maxDatapoints == 0 {
			break
		}
		foundKeywords := make(map[string]bool)
		for fk := range Filter {
			if strings.Contains(p.Title, Filter[fk]) || strings.Contains(p.Selftext, Filter[fk]) {
				continue outerLoop
			}
		}
		words := strings.Fields(p.Selftext + " " + p.Title)
		for _, word := range words {
			if t.HasKeysWithPrefix(word) {
				for _, k := range t.PrefixSearch(strings.ToLower(word)) {
					foundKeywords[k] = true
				}
			}
		}
		score := a.Analyzer.PolarityScores(fmt.Sprintf("%s %s", p.Title, p.Selftext))
		if len(foundKeywords) > 0 && score["compound"] < -0.7 {
			maxDatapoints--
			out <- models.RedditPostDetails{
				Id:        p.Id,
				Title:     p.Title,
				Subreddit: p.Subreddit,
			}
		}
	}
	return nil
}