package reddit

import "fmt"

type Headers map[string]string

func getHeaders(s string) Headers {
	return map[string]string{
		"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding":    "gzip, deflate, br",
		"Accept-Language":    "en-US,en;q=0.5",
		"Cache-Control":      "no-cache",
		"Pragma":             "no-cache",
		"Referer":            fmt.Sprintf("https://old.reddit.com/r/%s/", s),
		"Sec-Ch-Ua":          "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
		"Sec-Ch-Ua-Mobile":   "?0",
		"Sec-Ch-Ua-Platform": "\"macOS\"",
		"Sec-Fetch-Dest":     "document",
		"Sec-Fetch-Mode":     "navigate",
		"Sec-Fetch-Site":     "same-origin",
		"Sec-Fetch-User":     "?1",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
}
