package scrape

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Code based on https://github.com/jychp/cloudflare-bypass

var TOKEN_VALUE_ENV = "TOKEN_VALUE"

var PROXYS_URL = []string{}

const (
	WORKERS_FILE = ".workers"
	TOKEN_HEADER = "Px-Token"
	HOST_HEADER  = "Px-Host"
	IP_HEADER    = "Px-IP"
	FAKE_IP      = "1.2.3.4"
)

type CFProxy struct {
	Token     string
	ProxyHost string
	FakeIP    string
	UserAgent string
	Client    *http.Client
}

var proxyCount = -1

func loadProxyUrls() ([]string, error) {
	file, err := os.Open(WORKERS_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, strings.TrimSpace(scanner.Text()))
	}
	if len(text) == 0 {
		return nil, fmt.Errorf("The workers file is empty")
	}
	return text, nil
}

func GetCFProxy() (*CFProxy, error) {
	if len(PROXYS_URL) == 0 {
		urls, err := loadProxyUrls()
		if err != nil {
			return nil, fmt.Errorf("Error loading workers: %s", err)
		}
		PROXYS_URL = urls
	}
	proxyCount += 1
	if proxyCount >= len(PROXYS_URL) {
		proxyCount = 0
	}
	proxyHost := PROXYS_URL[proxyCount]
	return NewProxy(proxyHost, UserAgent, FAKE_IP)
}

func NewProxy(proxyHost string, ua string, fakeIP string) (*CFProxy, error) {
	token := os.Getenv(TOKEN_VALUE_ENV)
	if token == "" {
		return nil, fmt.Errorf("%s environment variable is not set", TOKEN_VALUE_ENV)
	}
	return &CFProxy{
		Token:     token,
		ProxyHost: proxyHost,
		UserAgent: ua,
		FakeIP:    fakeIP,
		Client:    &http.Client{Timeout: 2 * time.Second},
	}, nil
}

func (proxy *CFProxy) Get(URL string) (*http.Response, error) {
	return proxy.DoRequest("GET", URL, make(map[string]string))
}

func (proxy *CFProxy) DoRequest(method string, URL string, headers map[string]string) (*http.Response, error) {
	parsedURI, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	proxyfiedUrl := fmt.Sprintf("%s://%s%s", parsedURI.Scheme, proxy.ProxyHost, parsedURI.Path+"?"+parsedURI.RawQuery)
	req, err := http.NewRequest(method, proxyfiedUrl, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", proxy.UserAgent)
	req.Header.Set(HOST_HEADER, parsedURI.Hostname())
	req.Header.Set(IP_HEADER, proxy.FakeIP)
	req.Header.Set(TOKEN_HEADER, proxy.Token)
	return proxy.Client.Do(req)
}
