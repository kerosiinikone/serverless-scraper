package util

import (
	"testing"
	"time"
)

func TestCaller(t *testing.T) {
	t.Run("Test NewBackoffCaller", func(t *testing.T) {
		p, err := Proxy()
		if err != nil {
			t.Errorf("Proxy returned error: %v", err)
		}
		bc := NewBackoffCaller(map[string]string{}, 1, p)
		if bc == nil {
			t.Error("NewBackoffCaller returned nil")
		}
	})
	t.Run("Test Call", func(t *testing.T) {
		p, err := Proxy()
		if err != nil {
			t.Errorf("Proxy returned error: %v", err)
		}
		bc := NewBackoffCaller(map[string]string{}, 1, p)
		resp, err := bc.Call("https://www.google.com")
		if err != nil {
			t.Errorf("Call returned error: %v", err)
		}
		if resp == nil {
			t.Error("Call returned nil response")
		}
	})
	t.Run("Test Call with Backoffs", func(t *testing.T) {
		p, err := Proxy()
		if err != nil {
			t.Errorf("Proxy returned error: %v", err)
		}
		bc := NewBackoffCaller(map[string]string{}, 1, p)
		start := time.Now()
		resp, err := bc.Call("https://httpstat.us/429")
		if err == nil {
			t.Error("Call did not return error")
		}
		if resp != nil {
			t.Error("Call returned response")
		}
		// backoffTime := (int(bc.initialBackoff) << bc.maxRetries - 1)
		if time.Since(start) < 2*time.Second {
			t.Error("Call did not backoff")
		}
	})
}