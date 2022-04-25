package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("proxy-response", "this is a message added by proxy")
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
	return proxy, nil
}

func PingPong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ping Pong from Reverse Proxy Server"))
}

func TestProxy(t *testing.T) {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("http://127.0.0.1:1234")
	if err != nil {
		t.Fatal(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/ping", PingPong)

	if err := http.ListenAndServe(":8888", nil); err != nil {
		t.Fatal(err)
	}
}
