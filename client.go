// File: sockshttp/client.go
package sockshttp

import (
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
)

var DefaultEnvVar = "SOCKS_PROXY_URL"

// NewClientWithEnv returns an *http.Client using a SOCKS5 proxy URL
// from the specified environment variable (or "SOCKS_PROXY_URL" if empty).
// If the variable is unset or empty, a default http.Client is returned.
func NewClientWithEnv(envVar string) (*http.Client, error) {
	if envVar == "" {
		envVar = DefaultEnvVar
	}
	proxyURLStr := os.Getenv(envVar)
	if proxyURLStr == "" {
		return &http.Client{}, nil
	}

	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		return nil, err
	}

	dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return nil, err
	}

	transport := &http.Transport{
		DialContext: contextDialer.DialContext,
	}
	return &http.Client{Transport: transport}, nil
}
