package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/8whie/8whie-webprobe/internal/config"
	"github.com/8whie/8whie-webprobe/internal/engine"
	"github.com/8whie/8whie-webprobe/internal/output"
	"github.com/8whie/8whie-webprobe/internal/ui"
	"github.com/8whie/8whie-webprobe/internal/validation"
)

func main() {
	// Setup graceful interrupt handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cfg.ShowVersion {
		printVersion()
		os.Exit(0)
	}

	if cfg.ShowHelp {
		printComprehensiveHelp()
		os.Exit(0)
	}

	// Interactive Mode Dispatch
	if cfg.Interactive {
		menu := ui.NewInteractiveMenu(cfg)
		if err := menu.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Interactive error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Standard CLI Execution
	if err := validation.ValidateExecution(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "[!] Validation error: %v\n", err)
		fmt.Fprintln(os.Stderr, "    Run 'webprobe -h' for usage syntax and examples.")
		os.Exit(1)
	}

	validation.PrintSecurityNotice(cfg.Quiet)
	if cfg.InsecureTLS {
		validation.PrintInsecureTLSWarning()
	}

	writer, err := output.NewWriter(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to configure output: %v\n", err)
		os.Exit(1)
	}

	writer.PrintBanner()

	eng, err := engine.NewEngine(cfg, writer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to initialize engine: %v\n", err)
		os.Exit(1)
	}

	stats, err := eng.Run(ctx)
	if err != nil && ctx.Err() == nil {
		fmt.Fprintf(os.Stderr, "[!] Engine error during execution: %v\n", err)
		_ = writer.Close(stats)
		os.Exit(1)
	}

	if ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "\n[!] Operation cancelled by operator signal.")
	}

	writer.PrintSummary(stats)
	if err := writer.Close(stats); err != nil {
		fmt.Fprintf(os.Stderr, "[!] Failed to flush output: %v\n", err)
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("%s v%s\n", config.BrandName, config.Version)
	fmt.Printf("Developed by %s (Android • Termux • Linux • Security Research)\n", config.Author)
	fmt.Printf("YouTube: %s | Instagram: %s | Telegram: %s\n", config.YouTube, config.Instagram, config.Telegram)
	fmt.Println("License: MIT License (Copyright (c) 2026 8WHIE)")
}

func printComprehensiveHelp() {
	ui.DisplayOriginalBanner()
	helpText := `
NAME:
   8WHIE WebProbe - HTTP Endpoint Discovery & Security Research Tool

SYNOPSIS:
   webprobe -u <TARGET_URL> -w <WORDLIST> [OPTIONS]
   webprobe -i  (Launch interactive terminal mode)

DESCRIPTION:
   8WHIE WebProbe is a fast, original, open-source command-line utility for
   controlled HTTP discovery, path enumeration, and response analysis.
   Engineered for legitimate security researchers, developers, and system
   administrators testing authorized infrastructure.

AUTHORIZATION NOTICE:
   Use of this tool against targets without prior explicit written authorization
   is strictly prohibited. You are personally responsible for complying with all
   applicable computer security and network access laws.

TARGET & WORDLIST:
   -u, --url string
         Target URL containing a placeholder token (e.g. https://target.local/FUZZ)
   -w, --wordlist string
         Path to dictionary file with candidate payloads (one per line)
   -p, --placeholder string
         Placeholder token in URL, body, or headers (default "FUZZ")

HTTP REQUEST CONTROLS:
   -X, --method string
         HTTP request method (GET, HEAD, POST, PUT, DELETE, OPTIONS) (default "GET")
   -H, --header string
         Custom header 'Header-Name: Value' (can be repeated for multiple headers)
   -d, --data string
         Explicit HTTP request body data
   -r, --redirects
         Follow HTTP redirects (up to 10 hops) (default false)
   -t, --timeout int
         HTTP request timeout in seconds (default 10)
   -x, --proxy string
         Proxy server address (e.g. http://127.0.0.1:8080)
   -k, --insecure
         Allow untrusted TLS certificates (warning: insecure, disabled by default)

PERFORMANCE & PACING:
   -c, --concurrency int
         Number of concurrent worker routines (1 to 500) (default 20)
   -rate int
         Enforce rate limit in requests per second (0 = unrestricted)

RESPONSE MATCHING & FILTERING:
   -mc string
         Match only specific HTTP status codes (e.g. 200,204,301,302)
   -fc string
         Filter out specific HTTP status codes (default "404")
   -ms string
         Match specific response byte sizes (e.g. 512,1024)
   -fs string
         Filter out specific response byte sizes
   -mw string
         Match response word counts
   -fw string
         Filter out response word counts
   -ml string
         Match response line counts
   -fl string
         Filter out response line counts

OUTPUT & FORMATTING:
   -o, --output string
         Destination file path to save findings
   -of, --format string
         Output format: text, json, csv, md (default "text")
   -q, --quiet
         Quiet mode: print only discovered findings without banners
   -v, --verbose
         Verbose diagnostics: output detailed network warnings
   --no-color
         Disable ANSI terminal color output

OPERATIONAL MODES:
   -i, --interactive
         Start guided interactive terminal interface
   -V, --version
         Display version, brand, and social links
   -h, --help
         Display this comprehensive manual

EXAMPLES:
   1. Basic endpoint discovery:
      webprobe -u https://api.internal.local/FUZZ -w /usr/share/wordlists/common.txt

   2. Filter default 404s and match active endpoints:
      webprobe -u https://target.local/FUZZ -w paths.txt -mc 200,301,302,403

   3. Low-speed stealth/rate-limited diagnostic probe:
      webprobe -u https://target.local/FUZZ -w paths.txt -c 5 -rate 10

   4. Custom authorization and headers:
      webprobe -u https://target.local/FUZZ -w paths.txt -H "Authorization: Bearer <TOKEN>" -H "X-Source: Audit"

   5. Export structured results to JSON for CI/CD audit:
      webprobe -u https://target.local/FUZZ -w paths.txt -o findings.json -of json

   6. Test custom POST payload with JSON body:
      webprobe -u https://target.local/api -w ids.txt -X POST -d '{"userId": "FUZZ"}' -H "Content-Type: application/json"

TROUBLESHOOTING:
   - "target URL must contain placeholder": Ensure your URL contains 'FUZZ' (or custom -p token).
   - "connection refused": Verify the web server is running and accessible on the specified port.
   - High network error rates: Lower worker count (-c) or add a rate limit (-rate).
`
	fmt.Println(helpText)
}
