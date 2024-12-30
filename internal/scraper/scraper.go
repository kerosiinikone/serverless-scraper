package scraper

import "context"

type Scraper interface {
	Scrape(ctx context.Context, req *APIRequest, out chan struct{}) error
}

type Config struct {
	MaxDepth int
	DelayMax int
	DelayMin int
}

type APIRequest struct {
	ID        string `json:"request_id"`
	ClientID  string `json:"client_id"`
	Subreddit string `json:"subreddit"`
}