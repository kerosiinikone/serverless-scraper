package models

type DataEntry struct {
	Post      ForumTree `json:"data"`
	ClientID  string    `json:"client_id"`
	RequestID string    `json:"request_id"`
}

type QueueMessage struct {
	ClientID  string `json:"client_id"`
	RequestID string `json:"request_id"`
}