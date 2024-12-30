package models

type RedditPostDetails struct {
	Id        string `json:"id"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
}

type RedditPost struct {
	Title     string `json:"title"`
	Selftext  string `json:"selftext"`
	Id        string `json:"id"`
	Subreddit string `json:"subreddit"`
}

type PostContainer struct {
	Post RedditPost `json:"data"`
}

type RedditPostResponse struct {
	Data struct {
		After    string          `json:"after"`
		Children []PostContainer `json:"children"`
	}
}

type ForumTree struct {
	Id        string      `json:"id"`
	Subreddit string      `json:"subreddit"`
	Selftext  string      `json:"selftext"`
	Comments  []ReplyTree `json:"comments"`
}

type ReplyTree struct {
	Body    string      `json:"body"`
	Replies []ReplyTree `json:"replies"`
}