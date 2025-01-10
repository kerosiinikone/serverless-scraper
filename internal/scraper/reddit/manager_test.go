package reddit

import (
	"testing"
	"time"
)

// Unit tests
func TestResetTimer(t *testing.T) {
	finished := make(chan struct{})
	m := &Manager{
		timerStopped: false,
		finished:     finished,
	}
	start := time.Now()
	m.resetTimer(func() error {
		close(finished)
		return nil
	})
	if m.timer == nil {
		t.Error("Timer is nil")
	}
	select {
	case <-finished:
		if time.Since(start) < 8*time.Second {
			t.Error("Timer did not trigger")
		}
	}
}
