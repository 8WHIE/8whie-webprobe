package engine

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/8whie/8whie-webprobe/internal/config"
	"github.com/8whie/8whie-webprobe/internal/filters"
	"github.com/8whie/8whie-webprobe/internal/httpclient"
	"github.com/8whie/8whie-webprobe/internal/output"
	"github.com/8whie/8whie-webprobe/pkg/model"
)

// Engine drives the endpoint discovery pipeline.
type Engine struct {
	cfg        *config.Config
	client     *httpclient.Client
	filter     *filters.ResponseFilter
	writer     *output.Writer
	stats      model.ScanStats
	activeJobs sync.WaitGroup
}

// NewEngine initializes a discovery engine with coordinated components.
func NewEngine(cfg *config.Config, writer *output.Writer) (*Engine, error) {
	client, err := httpclient.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("engine initialization failed: %w", err)
	}

	filter := filters.NewResponseFilter(cfg.Criteria)

	return &Engine{
		cfg:    cfg,
		client: client,
		filter: filter,
		writer: writer,
	}, nil
}

// Run executes the discovery workflow across candidate payloads in the wordlist.
func (e *Engine) Run(ctx context.Context) (model.ScanStats, error) {
	e.stats.StartTime = time.Now()

	// Open wordlist file
	file, err := os.Open(e.cfg.WordlistPath)
	if err != nil {
		return e.stats, fmt.Errorf("cannot open wordlist: %w", err)
	}
	defer file.Close()

	// Channels for coordinated processing
	jobsChan := make(chan model.ProbeRequest, e.cfg.Concurrency*2)
	resultsChan := make(chan model.ProbeResult, e.cfg.Concurrency*2)

	// Rate limiter ticker if configured
	var rateTicker *time.Ticker
	if e.cfg.RateLimit > 0 {
		interval := time.Second / time.Duration(e.cfg.RateLimit)
		rateTicker = time.NewTicker(interval)
		defer rateTicker.Stop()
	}

	// 1. Start Result Collector Goroutine
	var collectorWg sync.WaitGroup
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for res := range resultsChan {
			if res.Matched {
				atomic.AddInt64(&e.stats.MatchedResults, 1)
				e.writer.WriteResult(res)
			} else {
				atomic.AddInt64(&e.stats.FilteredResults, 1)
			}
		}
	}()

	// 2. Start Worker Pool
	var workerWg sync.WaitGroup
	for i := 0; i < e.cfg.Concurrency; i++ {
		workerWg.Add(1)
		go func(workerID int) {
			defer workerWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobsChan:
					if !ok {
						return
					}

					resp, err := e.client.Execute(ctx, &job)
					atomic.AddInt64(&e.stats.SentRequests, 1)

					if err != nil {
						atomic.AddInt64(&e.stats.NetworkErrors, 1)
						if e.cfg.Verbose {
							fmt.Fprintf(os.Stderr, "[!] Network error (%s): %v\n", job.URL, err)
						}
						continue
					}

					matched, note := e.filter.Evaluate(resp)
					result := model.ProbeResult{
						Request:   job,
						Response:  *resp,
						Matched:   matched,
						MatchNote: note,
					}

					select {
					case resultsChan <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}(i)
	}

	// 3. Read Wordlist & Dispatch Jobs
	scanner := bufio.NewScanner(file)
	jobID := 0
	placeholder := e.cfg.Placeholder
	if placeholder == "" {
		placeholder = "FUZZ"
	}

dispatchLoop:
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			break dispatchLoop
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		jobID++
		e.stats.TotalCandidates++

		// Prepare probe request with substituted values
		targetURL := strings.ReplaceAll(e.cfg.TargetURL, placeholder, line)
		body := strings.ReplaceAll(e.cfg.Body, placeholder, line)

		headers := make(map[string]string, len(e.cfg.Headers))
		for k, v := range e.cfg.Headers {
			kSub := strings.ReplaceAll(k, placeholder, line)
			vSub := strings.ReplaceAll(v, placeholder, line)
			headers[kSub] = vSub
		}

		req := model.ProbeRequest{
			ID:        jobID,
			URL:       targetURL,
			Method:    e.cfg.Method,
			Headers:   headers,
			Body:      body,
			Payload:   line,
			Timestamp: time.Now(),
		}

		// Enforce rate limiter pacing
		if rateTicker != nil {
			select {
			case <-rateTicker.C:
			case <-ctx.Done():
				break dispatchLoop
			}
		}

		select {
		case jobsChan <- req:
		case <-ctx.Done():
			break dispatchLoop
		}
	}

	close(jobsChan)
	workerWg.Wait()
	close(resultsChan)
	collectorWg.Wait()

	if err := scanner.Err(); err != nil && ctx.Err() == nil {
		return e.stats, fmt.Errorf("error reading wordlist stream: %w", err)
	}

	e.stats.EndTime = time.Now()
	e.stats.Duration = e.stats.EndTime.Sub(e.stats.StartTime)
	if e.stats.Duration > 0 {
		e.stats.RequestsPerSecond = float64(e.stats.SentRequests) / e.stats.Duration.Seconds()
	}

	return e.stats, nil
}
