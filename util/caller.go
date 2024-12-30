package util

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type BackoffCaller struct {
	c *http.Client
	headers map[string]string
	maxRetries int 
	initialBackoff time.Duration
	proxy func(*http.Request) (*url.URL, error)
	ctx context.Context
}

func NewBackoffCaller(headers map[string]string, initialBackoff time.Duration, proxy func(*http.Request) (*url.URL, error)) *BackoffCaller { 
	bc := &BackoffCaller{
		headers: headers,
		initialBackoff: initialBackoff,
		proxy: proxy,
		maxRetries: 5,
	}
	bc.c = &http.Client{
        Transport: &http.Transport{Proxy: proxy},
    }
	return bc
}

func (bc *BackoffCaller) Call(reqUrl string) (*http.Response, error) {
	var (
        resp *http.Response
    )
    req, err := http.NewRequest("GET", reqUrl, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    for k, v := range bc.headers {
        req.Header.Add(k, v)
    }
    for retries := 0; retries < bc.maxRetries; retries++ {
        resp, err = bc.c.Do(req)
        if err == nil && resp.StatusCode == http.StatusOK {
            return resp, nil
        }
        if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
            log.Printf("Received 429, retrying in %v...", bc.initialBackoff)
            time.Sleep(bc.initialBackoff)
            bc.initialBackoff *= 2
			continue
        } else {
            return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
        }
    }
    return nil, fmt.Errorf("max retries reached for %s", reqUrl)
}