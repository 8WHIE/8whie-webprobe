package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/8whie/8whie-webprobe/pkg/model"
)

// Version information for 8WHIE WebProbe.
const (
	Version   = "1.0.0"
	BrandName = "8WHIE WebProbe"
	Author    = "8WHIE"
	Telegram  = "@Arnxkt / https://t.me/whiee"
	Instagram = "@aaynkt"
	YouTube   = "8WHIE"
)

// HeaderList supports repeated -H flags on the command line.
type HeaderList []string

func (h *HeaderList) String() string {
	return strings.Join(*h, ", ")
}

func (h *HeaderList) Set(val string) error {
	*h = append(*h, val)
	return nil
}

// Config represents all execution options for 8WHIE WebProbe.
type Config struct {
	TargetURL       string
	WordlistPath    string
	Method          string
	RawHeaders      HeaderList
	Headers         map[string]string
	Body            string
	Concurrency     int
	RateLimit       int // Requests per second (0 = unlimited)
	Timeout         time.Duration
	FollowRedirects bool
	InsecureTLS     bool
	ProxyURL        string
	Placeholder     string

	// Filters and Matchers
	MatchCodesRaw  string
	FilterCodesRaw string
	MatchSizesRaw  string
	FilterSizesRaw string
	MatchWordsRaw  string
	FilterWordsRaw string
	MatchLinesRaw  string
	FilterLinesRaw string

	Criteria model.FilterCriteria

	// Output controls
	OutputFile   string
	OutputFormat string // text, json, csv, md
	Quiet        bool
	Verbose      bool
	NoColor      bool

	// CLI operation modes
	Interactive bool
	ShowVersion bool
	ShowHelp    bool
}

// NewDefaultConfig initializes configuration with safe research defaults.
func NewDefaultConfig() *Config {
	return &Config{
		Method:          "GET",
		Headers:         make(map[string]string),
		Concurrency:     20,
		RateLimit:       0,
		Timeout:         10 * time.Second,
		FollowRedirects: false,
		InsecureTLS:     false,
		Placeholder:     "FUZZ",
		FilterCodesRaw:  "404",
		OutputFormat:    "text",
	}
}

// ParseFlags parses command-line arguments into a Config instance.
func ParseFlags(args []string) (*Config, error) {
	cfg := NewDefaultConfig()

	fs := flag.NewFlagSet("webprobe", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	// Target & Wordlist
	fs.StringVar(&cfg.TargetURL, "u", "", "Target URL containing placeholder (e.g., https://target.local/FUZZ)")
	fs.StringVar(&cfg.TargetURL, "url", "", "Alias for -u")
	fs.StringVar(&cfg.WordlistPath, "w", "", "Path to wordlist file containing candidate payloads")
	fs.StringVar(&cfg.WordlistPath, "wordlist", "", "Alias for -w")
	fs.StringVar(&cfg.Placeholder, "p", "FUZZ", "Placeholder token to replace in URL, headers, or body")
	fs.StringVar(&cfg.Placeholder, "placeholder", "FUZZ", "Alias for -p")

	// HTTP Options
	fs.StringVar(&cfg.Method, "X", "GET", "HTTP method (GET, HEAD, POST, PUT, DELETE, OPTIONS)")
	fs.StringVar(&cfg.Method, "method", "GET", "Alias for -X")
	fs.Var(&cfg.RawHeaders, "H", "Custom HTTP header 'Name: Value' (can be repeated)")
	fs.Var(&cfg.RawHeaders, "header", "Alias for -H")
	fs.StringVar(&cfg.Body, "d", "", "HTTP request body data")
	fs.StringVar(&cfg.Body, "data", "", "Alias for -d")
	fs.BoolVar(&cfg.FollowRedirects, "r", false, "Follow HTTP redirects (up to 10 hops)")
	fs.BoolVar(&cfg.FollowRedirects, "redirects", false, "Alias for -r")
	fs.BoolVar(&cfg.InsecureTLS, "k", false, "Allow insecure TLS connections (disables certificate validation)")
	fs.BoolVar(&cfg.InsecureTLS, "insecure", false, "Alias for -k")
	fs.StringVar(&cfg.ProxyURL, "x", "", "HTTP proxy address (e.g., http://127.0.0.1:8080)")
	fs.StringVar(&cfg.ProxyURL, "proxy", "", "Alias for -x")

	// Concurrency & Timing
	fs.IntVar(&cfg.Concurrency, "c", 20, "Number of concurrent probe workers (1-500)")
	fs.IntVar(&cfg.Concurrency, "concurrency", 20, "Alias for -c")
	fs.IntVar(&cfg.RateLimit, "rate", 0, "Rate limit in requests per second (0 = unrestricted)")
	var timeoutSec int
	fs.IntVar(&timeoutSec, "t", 10, "Request timeout in seconds")
	fs.IntVar(&timeoutSec, "timeout", 10, "Alias for -t")

	// Filtering & Matching
	fs.StringVar(&cfg.MatchCodesRaw, "mc", "", "Match HTTP status codes (comma-separated, e.g. 200,204,301,302)")
	fs.StringVar(&cfg.FilterCodesRaw, "fc", "404", "Filter out HTTP status codes (comma-separated, e.g. 404,500)")
	fs.StringVar(&cfg.MatchSizesRaw, "ms", "", "Match response byte sizes (comma-separated, e.g. 512,1024)")
	fs.StringVar(&cfg.FilterSizesRaw, "fs", "", "Filter out response byte sizes (comma-separated)")
	fs.StringVar(&cfg.MatchWordsRaw, "mw", "", "Match response word counts (comma-separated)")
	fs.StringVar(&cfg.FilterWordsRaw, "fw", "", "Filter out response word counts (comma-separated)")
	fs.StringVar(&cfg.MatchLinesRaw, "ml", "", "Match response line counts (comma-separated)")
	fs.StringVar(&cfg.FilterLinesRaw, "fl", "", "Filter out response line counts (comma-separated)")

	// Output
	fs.StringVar(&cfg.OutputFile, "o", "", "File path to save results")
	fs.StringVar(&cfg.OutputFile, "output", "", "Alias for -o")
	fs.StringVar(&cfg.OutputFormat, "of", "text", "Output format: text, json, csv, md")
	fs.StringVar(&cfg.OutputFormat, "format", "text", "Alias for -of")
	fs.BoolVar(&cfg.Quiet, "q", false, "Quiet mode: display only matched findings")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "Alias for -q")
	fs.BoolVar(&cfg.Verbose, "v", false, "Verbose mode: print detailed diagnostics")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Alias for -v")
	fs.BoolVar(&cfg.NoColor, "no-color", false, "Disable ANSI color formatting")

	// Modes
	fs.BoolVar(&cfg.Interactive, "i", false, "Launch interactive terminal menu")
	fs.BoolVar(&cfg.Interactive, "interactive", false, "Alias for -i")
	fs.BoolVar(&cfg.ShowVersion, "V", false, "Display version and build information")
	fs.BoolVar(&cfg.ShowVersion, "version", false, "Alias for -V")
	fs.BoolVar(&cfg.ShowHelp, "h", false, "Display comprehensive usage instructions")
	fs.BoolVar(&cfg.ShowHelp, "help", false, "Alias for -h")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg.Timeout = time.Duration(timeoutSec) * time.Second

	// Process headers
	for _, h := range cfg.RawHeaders {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			cfg.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Parse filter criteria
	if err := cfg.CompileCriteria(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// CompileCriteria parses raw comma-separated values into numeric slices.
func (c *Config) CompileCriteria() error {
	var err error
	if c.MatchCodesRaw != "" {
		if c.Criteria.MatchStatusCodes, err = ParseIntList(c.MatchCodesRaw); err != nil {
			return fmt.Errorf("invalid -mc (match codes): %w", err)
		}
	}
	if c.FilterCodesRaw != "" {
		if c.Criteria.FilterStatusCodes, err = ParseIntList(c.FilterCodesRaw); err != nil {
			return fmt.Errorf("invalid -fc (filter codes): %w", err)
		}
	}
	if c.MatchSizesRaw != "" {
		if c.Criteria.MatchSizes, err = ParseInt64List(c.MatchSizesRaw); err != nil {
			return fmt.Errorf("invalid -ms (match sizes): %w", err)
		}
	}
	if c.FilterSizesRaw != "" {
		if c.Criteria.FilterSizes, err = ParseInt64List(c.FilterSizesRaw); err != nil {
			return fmt.Errorf("invalid -fs (filter sizes): %w", err)
		}
	}
	if c.MatchWordsRaw != "" {
		if c.Criteria.MatchWords, err = ParseIntList(c.MatchWordsRaw); err != nil {
			return fmt.Errorf("invalid -mw (match words): %w", err)
		}
	}
	if c.FilterWordsRaw != "" {
		if c.Criteria.FilterWords, err = ParseIntList(c.FilterWordsRaw); err != nil {
			return fmt.Errorf("invalid -fw (filter words): %w", err)
		}
	}
	if c.MatchLinesRaw != "" {
		if c.Criteria.MatchLines, err = ParseIntList(c.MatchLinesRaw); err != nil {
			return fmt.Errorf("invalid -ml (match lines): %w", err)
		}
	}
	if c.FilterLinesRaw != "" {
		if c.Criteria.FilterLines, err = ParseIntList(c.FilterLinesRaw); err != nil {
			return fmt.Errorf("invalid -fl (filter lines): %w", err)
		}
	}
	return nil
}

// ParseIntList splits a comma-delimited string into integers.
func ParseIntList(raw string) ([]int, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	tokens := strings.Split(raw, ",")
	results := make([]int, 0, len(tokens))
	for _, t := range tokens {
		item := strings.TrimSpace(t)
		if item == "" {
			continue
		}
		val, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("expected numeric integer, got %q", item)
		}
		results = append(results, val)
	}
	return results, nil
}

// ParseInt64List splits a comma-delimited string into 64-bit integers.
func ParseInt64List(raw string) ([]int64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	tokens := strings.Split(raw, ",")
	results := make([]int64, 0, len(tokens))
	for _, t := range tokens {
		item := strings.TrimSpace(t)
		if item == "" {
			continue
		}
		val, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("expected numeric 64-bit integer, got %q", item)
		}
		results = append(results, val)
	}
	return results, nil
}

// ValidateSanity verifies configuration constraints before engine startup.
func (c *Config) ValidateSanity() error {
	if c.Concurrency < 1 || c.Concurrency > 500 {
		return errors.New("concurrency (-c) must be between 1 and 500 workers")
	}
	if c.RateLimit < 0 {
		return errors.New("rate limit (-rate) cannot be negative")
	}
	if c.Timeout <= 0 {
		return errors.New("timeout (-t) must be greater than 0")
	}
	c.Method = strings.ToUpper(strings.TrimSpace(c.Method))
	validMethods := map[string]bool{
		"GET": true, "HEAD": true, "POST": true, "PUT": true,
		"DELETE": true, "OPTIONS": true, "PATCH": true,
	}
	if !validMethods[c.Method] {
		return fmt.Errorf("unsupported HTTP method %q; supported methods are GET, HEAD, POST, PUT, DELETE, OPTIONS, PATCH", c.Method)
	}
	c.OutputFormat = strings.ToLower(strings.TrimSpace(c.OutputFormat))
	validFormats := map[string]bool{
		"text": true, "json": true, "csv": true, "md": true,
	}
	if !validFormats[c.OutputFormat] {
		return fmt.Errorf("unsupported output format %q; allowed: text, json, csv, md", c.OutputFormat)
	}
	return nil
}
