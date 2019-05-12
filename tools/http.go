package tools

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPCurl helpers for http requests
type HTTPCurl struct {
	HTTPClient *http.Client
	Timeout    time.Duration
}

// NewHTTPCurl new http helper object
func NewHTTPCurl() *HTTPCurl {
	h := &HTTPCurl{
		Timeout: time.Duration(30),
	}
	h.Init()
	return h

}

// Init start prep
func (g *HTTPCurl) Init() {
	//web
	g.HTTPClient = &http.Client{
		Timeout: time.Duration(g.Timeout * time.Second),
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   g.Timeout * time.Second,
				KeepAlive: 0,
			}).Dial,
			TLSHandshakeTimeout: g.Timeout * time.Second,
		},
	}
}

// Get request via get
func (g *HTTPCurl) Get(url string) (string, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", -1, err
	}
	//settings
	req.Close = true
	req.Header.Set("Connection", "close")
	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", -1, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	// read the body
	return strings.TrimSpace(string(contents)), resp.StatusCode, nil
}
