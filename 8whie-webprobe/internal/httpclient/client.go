package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/8whie/8whie-webprobe/internal/config"
	"github.com/8whie/8whie-webprobe/internal/filters"
	"github.com/8whie/8whie-webprobe/pkg/model"
)

const maxResponseBodyCap = 5 * 1024 * 1024 // 5 MB read cap per response

// Client manages outbound HTTP discovery probes with robust connection reuse.
type Client struct {
	httpClient *http.Client
	cfg        *config.Config
}

// NewClient constructs an optimized HTTP probe client respecting configuration.
func NewClient(cfg *config.Config) (*Client, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.InsecureTLS, // Explicitly controlled; defaults to false
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   cfg.Timeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          cfg.Concurrency * 2,
		MaxIdleConnsPerHost:   cfg.Concurrency,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   cfg.Timeout,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
		DisableCompression:    false,
	}

	if cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy configuration: %w", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	if !cfg.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		}
	}

	return &Client{
		httpClient: client,
		cfg:        cfg,
	}, nil
}

// Execute performs an isolated HTTP probe request and computes response metrics.
func (c *Client) Execute(ctx context.Context, req *model.ProbeRequest) (*model.ProbeResponse, error) {
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to construct HTTP request: %w", err)
	}

	// Apply default User-Agent unless overridden
	httpReq.Header.Set("User-Agent", model.DefaultUserAgent)
	httpReq.Header.Set("Accept", "*/*")

	// Apply user-defined headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	startTime := time.Now()
	httpResp, err := c.httpClient.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		return &model.ProbeResponse{
			StatusCode: 0,
			Duration:   duration,
			Error:      err.Error(),
		}, err
	}
	defer httpResp.Body.Close()

	// Safely read body with boundary
	limitedReader := io.LimitReader(httpResp.Body, maxResponseBodyCap)
	var bodyBuf bytes.Buffer
	readBytes, readErr := bodyBuf.ReadFrom(limitedReader)
	if readErr != nil && readErr != io.EOF {
		return &model.ProbeResponse{
			StatusCode: httpResp.StatusCode,
			Duration:   duration,
			Error:      fmt.Sprintf("body read error: %v", readErr),
		}, readErr
	}

	contentLength := httpResp.ContentLength
	if contentLength < 0 {
		contentLength = readBytes
	}

	bodyBytes := bodyBuf.Bytes()
	words, lines := filters.CountWordsAndLines(bodyBytes)

	redirectURL := ""
	if loc, locErr := httpResp.Location(); locErr == nil && loc != nil {
		redirectURL = loc.String()
	}

	return &model.ProbeResponse{
		StatusCode:    httpResp.StatusCode,
		ContentLength: contentLength,
		WordCount:     words,
		LineCount:     lines,
		Duration:      duration,
		RedirectURL:   redirectURL,
		ContentType:   httpResp.Header.Get("Content-Type"),
		ServerHeader:  httpResp.Header.Get("Server"),
	}, nil
}
