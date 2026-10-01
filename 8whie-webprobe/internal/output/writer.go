package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/8whie/8whie-webprobe/internal/config"
	"github.com/8whie/8whie-webprobe/pkg/model"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
)

// Writer coordinates terminal and file reporting of discovery results.
type Writer struct {
	cfg        *config.Config
	stdout     io.Writer
	fileHandle *os.File
	mu         sync.Mutex
	results    []model.ProbeResult
}

// NewWriter initializes output streams based on user configuration.
func NewWriter(cfg *config.Config) (*Writer, error) {
	w := &Writer{
		cfg:     cfg,
		stdout:  os.Stdout,
		results: make([]model.ProbeResult, 0, 128),
	}

	if cfg.OutputFile != "" {
		f, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open output file %q: %w", cfg.OutputFile, err)
		}
		w.fileHandle = f
	}

	return w, nil
}

// Close flushes final structured data and closes opened output files.
func (w *Writer) Close(stats model.ScanStats) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fileHandle != nil {
		defer w.fileHandle.Close()

		switch strings.ToLower(w.cfg.OutputFormat) {
		case "json":
			type JSONReport struct {
				Tool       string              `json:"tool"`
				Version    string              `json:"version"`
				Target     string              `json:"target"`
				Wordlist   string              `json:"wordlist"`
				Statistics model.ScanStats     `json:"statistics"`
				Results    []model.ProbeResult `json:"results"`
			}
			report := JSONReport{
				Tool:       config.BrandName,
				Version:    config.Version,
				Target:     w.cfg.TargetURL,
				Wordlist:   w.cfg.WordlistPath,
				Statistics: stats,
				Results:    w.results,
			}
			enc := json.NewEncoder(w.fileHandle)
			enc.SetIndent("", "  ")
			if err := enc.Encode(report); err != nil {
				return fmt.Errorf("failed to serialize JSON output: %w", err)
			}

		case "csv":
			csvWriter := csv.NewWriter(w.fileHandle)
			_ = csvWriter.Write([]string{"ID", "Payload", "URL", "Status", "Size", "Words", "Lines", "DurationMs", "Redirect"})
			for _, r := range w.results {
				_ = csvWriter.Write([]string{
					strconv.Itoa(r.Request.ID),
					r.Request.Payload,
					r.Request.URL,
					strconv.Itoa(r.Response.StatusCode),
					strconv.FormatInt(r.Response.ContentLength, 10),
					strconv.Itoa(r.Response.WordCount),
					strconv.Itoa(r.Response.LineCount),
					strconv.FormatInt(r.Response.Duration.Milliseconds(), 10),
					r.Response.RedirectURL,
				})
			}
			csvWriter.Flush()

		case "md":
			w.writeMarkdown(w.fileHandle, stats)

		default: // plain text table
			for _, r := range w.results {
				fmt.Fprintf(w.fileHandle, "[%d] %s -> Size: %d, Words: %d, Lines: %d, Time: %v\n",
					r.Response.StatusCode, r.Request.URL, r.Response.ContentLength,
					r.Response.WordCount, r.Response.LineCount, r.Response.Duration)
			}
		}
	}

	return nil
}

// PrintBanner outputs the discovery configuration header.
func (w *Writer) PrintBanner() {
	if w.cfg.Quiet {
		return
	}

	cBold := colorBold
	cCyan := colorCyan
	cDim := colorDim
	cReset := colorReset

	if w.cfg.NoColor {
		cBold, cCyan, cDim, cReset = "", "", "", ""
	}

	banner := fmt.Sprintf(`%s%s
  ___ _ _ _ _  _ ___   _ _ _     _   ___          _         
 ( _ ) | | | || |_ _| | | | |___| |_| _ \_ _ ___ | |__  ___ 
 / _ \_  _ | __ || |  | | | / -_) '_ \  _/ '_/ _ \| '_ \/ -_)
 \___/ |_| |_||_|___| |_____/\___|_,__/_| |_| \___/|_.__/\___|
%s%s
 :: 8WHIE WebProbe           : v%s
 :: Security Research & Tool : 8WHIE (@Arnxkt | https://t.me/whiee)
 :: Target                   : %s
 :: Wordlist                 : %s
 :: Method                   : %s
 :: Concurrency / Rate       : %d workers / %d rps
 :: Timeout / Redirects      : %v / %t
%s`,
		cCyan, cBold, cReset, cDim,
		config.Version,
		w.cfg.TargetURL,
		w.cfg.WordlistPath,
		w.cfg.Method,
		w.cfg.Concurrency,
		w.cfg.RateLimit,
		w.cfg.Timeout,
		w.cfg.FollowRedirects,
		cReset,
	)

	fmt.Fprintln(w.stdout, banner)
	w.PrintTableHeader()
}

// PrintTableHeader outputs tabular column titles.
func (w *Writer) PrintTableHeader() {
	if w.cfg.Quiet {
		return
	}

	cBold := colorBold
	cDim := colorDim
	cReset := colorReset

	if w.cfg.NoColor {
		cBold, cDim, cReset = "", "", ""
	}

	header := fmt.Sprintf("%s%s%-8s %-10s %-8s %-8s %-10s %-25s %s%s",
		cBold, cDim, "STATUS", "SIZE", "WORDS", "LINES", "DURATION", "PAYLOAD", "REDIRECT", cReset)
	sep := strings.Repeat("-", 80)

	fmt.Fprintln(w.stdout, header)
	fmt.Fprintln(w.stdout, sep)
}

// WriteResult outputs a single matched discovery finding.
func (w *Writer) WriteResult(res model.ProbeResult) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.results = append(w.results, res)

	if w.cfg.Quiet {
		// In quiet mode, just print payload or URL
		fmt.Fprintln(w.stdout, res.Request.Payload)
		return
	}

	statusColor := colorGreen
	if res.Response.StatusCode >= 300 && res.Response.StatusCode < 400 {
		statusColor = colorCyan
	} else if res.Response.StatusCode == 401 || res.Response.StatusCode == 403 {
		statusColor = colorYellow
	} else if res.Response.StatusCode >= 500 {
		statusColor = colorRed
	}

	cReset := colorReset
	if w.cfg.NoColor {
		statusColor = ""
		cReset = ""
	}

	durStr := fmt.Sprintf("%dms", res.Response.Duration.Milliseconds())
	redir := ""
	if res.Response.RedirectURL != "" {
		redir = fmt.Sprintf("-> %s", res.Response.RedirectURL)
	}

	row := fmt.Sprintf("%s%-8d%s %-10d %-8d %-8d %-10s %-25s %s",
		statusColor,
		res.Response.StatusCode,
		cReset,
		res.Response.ContentLength,
		res.Response.WordCount,
		res.Response.LineCount,
		durStr,
		res.Request.Payload,
		redir,
	)

	fmt.Fprintln(w.stdout, row)
}

// PrintSummary outputs ending scan statistics.
func (w *Writer) PrintSummary(stats model.ScanStats) {
	if w.cfg.Quiet {
		return
	}

	cBold := colorBold
	cCyan := colorCyan
	cReset := colorReset

	if w.cfg.NoColor {
		cBold, cCyan, cReset = "", "", ""
	}

	fmt.Fprintln(w.stdout, strings.Repeat("=", 80))
	fmt.Fprintf(w.stdout, "%s%sDiscovery Completed in %v%s\n", cBold, cCyan, stats.Duration.Round(time.Millisecond), cReset)
	fmt.Fprintf(w.stdout, "  Total Candidates: %d | Sent Requests: %d | Matched: %d | Errors: %d\n",
		stats.TotalCandidates, stats.SentRequests, stats.MatchedResults, stats.NetworkErrors)
	fmt.Fprintf(w.stdout, "  Average Speed: %.1f req/sec\n", stats.RequestsPerSecond)
	if w.cfg.OutputFile != "" {
		fmt.Fprintf(w.stdout, "  Output Saved: %s (%s format)\n", w.cfg.OutputFile, w.cfg.OutputFormat)
	}
	fmt.Fprintln(w.stdout, strings.Repeat("=", 80))
}

// writeMarkdown generates a clean GitHub-compatible Markdown table.
func (w *Writer) writeMarkdown(out io.Writer, stats model.ScanStats) {
	fmt.Fprintf(out, "# 8WHIE WebProbe Discovery Report\n\n")
	fmt.Fprintf(out, "- **Target**: `%s`\n", w.cfg.TargetURL)
	fmt.Fprintf(out, "- **Wordlist**: `%s`\n", w.cfg.WordlistPath)
	fmt.Fprintf(out, "- **Duration**: %v\n", stats.Duration.Round(time.Millisecond))
	fmt.Fprintf(out, "- **Requests**: %d (Speed: %.1f req/s)\n", stats.SentRequests, stats.RequestsPerSecond)
	fmt.Fprintf(out, "- **Matched Findings**: %d\n\n", len(w.results))

	fmt.Fprintf(out, "| Status | Size (Bytes) | Words | Lines | Duration | Payload | Redirect |\n")
	fmt.Fprintf(out, "| :---: | :---: | :---: | :---: | :---: | :--- | :--- |\n")

	for _, r := range w.results {
		fmt.Fprintf(out, "| %d | %d | %d | %d | %dms | `%s` | %s |\n",
			r.Response.StatusCode,
			r.Response.ContentLength,
			r.Response.WordCount,
			r.Response.LineCount,
			r.Response.Duration.Milliseconds(),
			r.Request.Payload,
			r.Response.RedirectURL,
		)
	}
}
