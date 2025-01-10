package util

import (
	"bufio"
	_ "embed"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/rand"
)

//go:embed data/proxy.txt
var proxyData []byte 

var proxies []string

func Proxy() (func(*http.Request) (*url.URL, error), error) {
	proxies := loadProxies()
	if len(proxies) == 0 {
		return nil, nil
	}
	p := proxies[rand.Intn(len(proxies))]
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
	scanner := bufio.NewScanner(strings.NewReader(string(proxyData)))
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	return proxies
}
