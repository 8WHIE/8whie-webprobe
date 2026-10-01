package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/8whie/8whie-webprobe/internal/config"
	"github.com/8whie/8whie-webprobe/internal/engine"
	"github.com/8whie/8whie-webprobe/internal/filters"
	"github.com/8whie/8whie-webprobe/internal/httpclient"
	"github.com/8whie/8whie-webprobe/internal/output"
	"github.com/8whie/8whie-webprobe/internal/validation"
	"github.com/8whie/8whie-webprobe/pkg/model"
)

// TestWordAndLineCounting verifies accurate byte analysis.
func TestWordAndLineCounting(t *testing.T) {
	testData := []byte("Hello world from 8WHIE WebProbe\nSecond line has four words\nThird line")
	words, lines := filters.CountWordsAndLines(testData)

	if words != 12 {
		t.Errorf("Expected 12 words, got %d", words)
	}
	if lines != 3 {
		t.Errorf("Expected 3 lines, got %d", lines)
	}

	emptyWords, emptyLines := filters.CountWordsAndLines([]byte(""))
	if emptyWords != 0 || emptyLines != 0 {
		t.Errorf("Expected 0 words and 0 lines for empty buffer, got %d words, %d lines", emptyWords, emptyLines)
	}
}

// TestFilterCriteriaEvaluation tests all filtering and matching mechanisms.
func TestFilterCriteriaEvaluation(t *testing.T) {
	criteria := model.FilterCriteria{
		MatchStatusCodes:  []int{200, 301},
		FilterStatusCodes: []int{404},
		MatchSizes:        []int64{100, 200},
		FilterWords:       []int{99},
	}

	filter := filters.NewResponseFilter(criteria)

	// Case 1: Matching 200 with size 100
	resp1 := &model.ProbeResponse{
		StatusCode:    200,
		ContentLength: 100,
		WordCount:     15,
		LineCount:     2,
	}
	matched, _ := filter.Evaluate(resp1)
	if !matched {
		t.Errorf("Expected resp1 to match criteria")
	}

	// Case 2: Status 404 (should be filtered)
	resp2 := &model.ProbeResponse{
		StatusCode:    404,
		ContentLength: 100,
		WordCount:     15,
		LineCount:     2,
	}
	matched2, _ := filter.Evaluate(resp2)
	if matched2 {
		t.Errorf("Expected resp2 (404) to be rejected by filter")
	}

	// Case 3: Status 200 but unmatched size 50
	resp3 := &model.ProbeResponse{
		StatusCode:    200,
		ContentLength: 50,
		WordCount:     15,
		LineCount:     2,
	}
	matched3, _ := filter.Evaluate(resp3)
	if matched3 {
		t.Errorf("Expected resp3 (unmatched size) to be rejected")
	}

	// Case 4: Filtered word count 99
	resp4 := &model.ProbeResponse{
		StatusCode:    200,
		ContentLength: 100,
		WordCount:     99,
		LineCount:     2,
	}
	matched4, _ := filter.Evaluate(resp4)
	if matched4 {
		t.Errorf("Expected resp4 (filtered word count) to be rejected")
	}
}

// TestConfigParsing verifies CLI flag compilation.
func TestConfigParsing(t *testing.T) {
	args := []string{
		"-u", "https://example.com/FUZZ",
		"-w", "/tmp/paths.txt",
		"-c", "40",
		"-rate", "100",
		"-mc", "200,302",
		"-fc", "404,500",
		"-ms", "1024,2048",
		"-of", "json",
	}

	cfg, err := config.ParseFlags(args)
	if err != nil {
		t.Fatalf("Failed to parse valid flags: %v", err)
	}

	if cfg.TargetURL != "https://example.com/FUZZ" {
		t.Errorf("Unexpected URL: %s", cfg.TargetURL)
	}
	if cfg.Concurrency != 40 {
		t.Errorf("Unexpected concurrency: %d", cfg.Concurrency)
	}
	if cfg.RateLimit != 100 {
		t.Errorf("Unexpected rate limit: %d", cfg.RateLimit)
	}
	if len(cfg.Criteria.MatchStatusCodes) != 2 || cfg.Criteria.MatchStatusCodes[0] != 200 {
		t.Errorf("Unexpected match status codes: %v", cfg.Criteria.MatchStatusCodes)
	}
	if len(cfg.Criteria.FilterStatusCodes) != 2 || cfg.Criteria.FilterStatusCodes[0] != 404 {
		t.Errorf("Unexpected filter status codes: %v", cfg.Criteria.FilterStatusCodes)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("Unexpected format: %s", cfg.OutputFormat)
	}
}

// TestValidationErrors verifies sanity guardrails.
func TestValidationErrors(t *testing.T) {
	// Missing URL
	cfg1 := config.NewDefaultConfig()
	cfg1.WordlistPath = "/tmp/paths.txt"
	if err := validation.ValidateExecution(cfg1); err == nil {
		t.Errorf("Expected error for missing URL")
	}

	// Missing placeholder
	cfg2 := config.NewDefaultConfig()
	cfg2.TargetURL = "https://example.com/admin"
	cfg2.WordlistPath = "/tmp/paths.txt"
	if err := validation.ValidateExecution(cfg2); err == nil {
		t.Errorf("Expected error for missing FUZZ placeholder")
	}

	// Concurrency limits
	cfg3 := config.NewDefaultConfig()
	cfg3.TargetURL = "https://example.com/FUZZ"
	cfg3.Concurrency = 9999
	if err := cfg3.ValidateSanity(); err == nil {
		t.Errorf("Expected error for out-of-bounds concurrency")
	}

	// Invalid HTTP method
	cfg4 := config.NewDefaultConfig()
	cfg4.TargetURL = "https://example.com/FUZZ"
	cfg4.Method = "HACK"
	if err := cfg4.ValidateSanity(); err == nil {
		t.Errorf("Expected error for unsupported HTTP method")
	}
}

// TestIntegrationMockServer performs a real end-to-end discovery run against a local mock HTTP server.
func TestIntegrationMockServer(t *testing.T) {
	// 1. Setup local mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/admin":
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("Forbidden area"))
		case "/login":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Login Form Content"))
		case "/api":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("API v1 endpoints"))
		case "/redirect":
			http.Redirect(w, r, "/login", http.StatusFound)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("404 Page Not Found"))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	// 2. Prepare mock wordlist file
	tmpDir := t.TempDir()
	wordlistPath := filepath.Join(tmpDir, "words.txt")
	wordsContent := "admin\nlogin\napi\nmissing1\nmissing2\nredirect\n"
	if err := os.WriteFile(wordlistPath, []byte(wordsContent), 0644); err != nil {
		t.Fatalf("Failed to create temporary wordlist: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "results.json")

	// 3. Configure Discovery Scan
	cfg := config.NewDefaultConfig()
	cfg.TargetURL = server.URL + "/FUZZ"
	cfg.WordlistPath = wordlistPath
	cfg.Concurrency = 2
	cfg.FilterCodesRaw = "404" // Exclude 404
	cfg.OutputFile = outputPath
	cfg.OutputFormat = "json"
	cfg.Quiet = true
	_ = cfg.CompileCriteria()

	if err := validation.ValidateExecution(cfg); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	w, err := output.NewWriter(cfg)
	if err != nil {
		t.Fatalf("Writer init failed: %v", err)
	}

	eng, err := engine.NewEngine(cfg, w)
	if err != nil {
		t.Fatalf("Engine init failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := eng.Run(ctx)
	if err != nil {
		t.Fatalf("Engine run failed: %v", err)
	}

	if err := w.Close(stats); err != nil {
		t.Fatalf("Writer close failed: %v", err)
	}

	// 4. Verify outcomes: should find admin (403), login (200), api (200), redirect (302) = 4 matches
	if stats.MatchedResults != 4 {
		t.Errorf("Expected 4 matched results (admin, login, api, redirect), got %d", stats.MatchedResults)
	}
	if stats.FilteredResults != 2 {
		t.Errorf("Expected 2 filtered 404 results, got %d", stats.FilteredResults)
	}

	// Check JSON output file exists and has content
	jsonBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if !strings.Contains(string(jsonBytes), "8WHIE WebProbe") {
		t.Errorf("Expected report to contain 8WHIE WebProbe branding")
	}
	if !strings.Contains(string(jsonBytes), "/login") {
		t.Errorf("Expected report to include /login endpoint")
	}
}

// TestHttpClientConnectionPool verifies client creation and single probe execution.
func TestHttpClientConnectionPool(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Audit") != "Authorized" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "Echo success")
	}))
	defer ts.Close()

	cfg := config.NewDefaultConfig()
	client, err := httpclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("Client creation failed: %v", err)
	}

	req := &model.ProbeRequest{
		ID:     1,
		URL:    ts.URL,
		Method: "GET",
		Headers: map[string]string{
			"X-Audit": "Authorized",
		},
		Payload:   "test",
		Timestamp: time.Now(),
	}

	resp, err := client.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
	if resp.WordCount != 2 {
		t.Errorf("Expected 2 words ('Echo success'), got %d", resp.WordCount)
	}
}
