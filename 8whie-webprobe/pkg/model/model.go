package model

import (
	"net/http"
	"time"
)

// ProbeRequest encapsulates outgoing HTTP request parameters for an individual probe.
type ProbeRequest struct {
	ID        int               `json:"id"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      string            `json:"body,omitempty"`
	Payload   string            `json:"payload"`
	Timestamp time.Time         `json:"timestamp"`
}

// ProbeResponse holds the observed metrics and headers from an HTTP response.
type ProbeResponse struct {
	StatusCode    int           `json:"status_code"`
	ContentLength int64         `json:"content_length"`
	WordCount     int           `json:"word_count"`
	LineCount     int           `json:"line_count"`
	Duration      time.Duration `json:"duration_ms"`
	RedirectURL   string        `json:"redirect_url,omitempty"`
	ContentType   string        `json:"content_type,omitempty"`
	ServerHeader  string        `json:"server,omitempty"`
	Error         string        `json:"error,omitempty"`
}

// ProbeResult pairs a probe request with its corresponding response and match status.
type ProbeResult struct {
	Request   ProbeRequest  `json:"request"`
	Response  ProbeResponse `json:"response"`
	Matched   bool          `json:"matched"`
	MatchNote string        `json:"match_note,omitempty"`
}

// FilterCriteria defines filtering and matching thresholds for responses.
type FilterCriteria struct {
	MatchStatusCodes  []int   `json:"match_status_codes,omitempty"`
	FilterStatusCodes []int   `json:"filter_status_codes,omitempty"`
	MatchSizes        []int64 `json:"match_sizes,omitempty"`
	FilterSizes       []int64 `json:"filter_sizes,omitempty"`
	MatchWords        []int   `json:"match_words,omitempty"`
	FilterWords       []int   `json:"filter_words,omitempty"`
	MatchLines        []int   `json:"match_lines,omitempty"`
	FilterLines       []int   `json:"filter_lines,omitempty"`
}

// ScanStats tracks aggregate metrics during and after discovery execution.
type ScanStats struct {
	TotalCandidates   int           `json:"total_candidates"`
	SentRequests      int64         `json:"sent_requests"`
	MatchedResults    int64         `json:"matched_results"`
	FilteredResults   int64         `json:"filtered_results"`
	NetworkErrors     int64         `json:"network_errors"`
	StartTime         time.Time     `json:"start_time"`
	EndTime           time.Time     `json:"end_time"`
	Duration          time.Duration `json:"duration"`
	RequestsPerSecond float64       `json:"requests_per_second"`
}

// ProbeProgress provides real-time progress information to consumers.
type ProbeProgress struct {
	CurrentIndex int
	TotalCount   int
	CurrentRPS   float64
	LastResult   *ProbeResult
}

// StandardHeaders provides default HTTP client headers.
var DefaultUserAgent = "8WHIE-WebProbe/1.0.0 (+https://github.com/8whie/8whie-webprobe; Security-Research)"

// HeaderMap converts http.Header to flat map.
func HeaderMap(h http.Header) map[string]string {
	m := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}
