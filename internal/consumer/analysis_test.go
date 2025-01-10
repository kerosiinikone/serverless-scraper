package consumer

import (
	"encoding/json"
	"testing"

	"github.com/kerosiinikone/serverless-scraper/pkg/models"
)

// Unit tests
func TestFormatMessages(t *testing.T) {
	var unmarshaled Messages

	messages := []models.DataEntry{
		{
			Post: models.ForumTree{
				Selftext: "Hello",
			},
		},
		{
			Post: models.ForumTree{
				Selftext: "World",

			},
		},
	}
	formatted, err := formatMessages(messages)
	if err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal(formatted, &unmarshaled); err != nil {
		t.Error(err)
	}
	if unmarshaled.Messages[0].Body != "Hello" {
		t.Error("Selftext mismatch")
	}
	if unmarshaled.Messages[1].Body != "World" {
		t.Error("Selftext mismatch")
	}
}
