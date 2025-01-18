package util

import (
	"bufio"
	_ "embed"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/exp/rand"
)

//go:embed data/proxy.txt
var proxyData []byte

var proxies []string

func Proxy() (func(*http.Request) (*url.URL, error), error) {
	prs := loadProxies()
	if len(prs) == 0 {
		return nil, nil
	}
	p := prs[rand.Intn(len(prs))]
	proxyUrl, err := url.Parse(p)
	if err != nil {
		return nil, err
	}
	return http.ProxyURL(proxyUrl), nil
}

func loadProxies() []string {
	if len(proxies) != 0 {
		return proxies
	}
	if os.Getenv("TEST_PROXY") != "" {
		return []string{os.Getenv("TEST_PROXY")}
	}
	scanner := bufio.NewScanner(strings.NewReader(string(proxyData)))
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	return proxies
}
