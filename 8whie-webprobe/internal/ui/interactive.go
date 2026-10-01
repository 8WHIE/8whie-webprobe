package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/8whie/8whie-webprobe/internal/config"
	"github.com/8whie/8whie-webprobe/internal/engine"
	"github.com/8whie/8whie-webprobe/internal/output"
	"github.com/8whie/8whie-webprobe/internal/validation"
)

// InteractiveMenu presents the guided terminal management console.
type InteractiveMenu struct {
	cfg    *config.Config
	reader *bufio.Reader
}

// NewInteractiveMenu initializes the interactive workflow.
func NewInteractiveMenu(cfg *config.Config) *InteractiveMenu {
	return &InteractiveMenu{
		cfg:    cfg,
		reader: bufio.NewReader(os.Stdin),
	}
}

// DisplayOriginalBanner prints the custom 8WHIE WebProbe ASCII header.
func DisplayOriginalBanner() {
	banner := `
 ===============================================================================
   ___ _ _ _ _  _ ___   _ _ _     _   ___          _         
  ( _ ) | | | || |_ _| | | | |___| |_| _ \_ _ ___ | |__  ___ 
  / _ \_  _ | __ || |  | | | / -_) '_ \  _/ '_/ _ \| '_ \/ -_)
  \___/ |_| |_||_|___| |_____/\___|_,__/_| |_| \___/|_.__/\___|

   8WHIE WebProbe — HTTP Discovery & Security Research
   Developed by 8WHIE (Android • Termux • Linux • Security Research)
   YouTube: 8WHIE | Instagram: @aaynkt | Telegram: @Arnxkt | Channel: t.me/whiee
 ===============================================================================`
	fmt.Println(banner)
}

// Run executes the main interactive prompt loop.
func (im *InteractiveMenu) Run(ctx context.Context) error {
	DisplayOriginalBanner()
	validation.PrintSecurityNotice(false)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n[!] Interactive session terminated.")
			return nil
		default:
		}

		im.printMainMenu()
		choice := im.promptString("Select an option (1-10): ")

		switch strings.TrimSpace(choice) {
		case "1":
			im.runDiscovery(ctx)
		case "2":
			im.configureTarget()
		case "3":
			im.configureWordlist()
		case "4":
			im.configureRequestSettings()
		case "5":
			im.configureFilters()
		case "6":
			im.configureOutput()
		case "7":
			im.viewCurrentConfig()
		case "8":
			im.displayHelp()
		case "9":
			im.displayAbout()
		case "10", "q", "exit":
			fmt.Println("\n[+] Exiting 8WHIE WebProbe. Good luck with your authorized research!")
			return nil
		default:
			fmt.Println("\n[!] Invalid selection. Please choose an option from 1 to 10.")
		}
	}
}

func (im *InteractiveMenu) printMainMenu() {
	fmt.Println("\n--- [ 8WHIE WebProbe :: Interactive Console ] ---")
	fmt.Println("  1. Start Web Discovery")
	fmt.Println("  2. Configure Target")
	fmt.Println("  3. Configure Wordlist")
	fmt.Println("  4. Request Settings (Method, Headers, Body, TLS, Proxy)")
	fmt.Println("  5. Response Filters (Codes, Sizes, Words, Lines)")
	fmt.Println("  6. Output Settings (File, Format, Verbose, Quiet)")
	fmt.Println("  7. View Configuration")
	fmt.Println("  8. Help & Documentation")
	fmt.Println("  9. About & Author")
	fmt.Println(" 10. Exit")
	fmt.Println("--------------------------------------------------")
}

func (im *InteractiveMenu) promptString(label string) string {
	fmt.Print(label)
	text, _ := im.reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func (im *InteractiveMenu) configureTarget() {
	fmt.Println("\n[Target Configuration]")
	current := im.cfg.TargetURL
	if current == "" {
		current = "(none)"
	}
	fmt.Printf("Current Target: %s\n", current)
	fmt.Printf("Current Placeholder: %s\n", im.cfg.Placeholder)

	val := im.promptString("Enter Target URL (e.g. https://target.local/FUZZ, leave empty to keep): ")
	if val != "" {
		im.cfg.TargetURL = val
	}

	ph := im.promptString(fmt.Sprintf("Enter Placeholder token [%s]: ", im.cfg.Placeholder))
	if ph != "" {
		im.cfg.Placeholder = ph
	}
	fmt.Println("[+] Target updated.")
}

func (im *InteractiveMenu) configureWordlist() {
	fmt.Println("\n[Wordlist Configuration]")
	current := im.cfg.WordlistPath
	if current == "" {
		current = "(none)"
	}
	fmt.Printf("Current Wordlist: %s\n", current)

	val := im.promptString("Enter Wordlist file path (leave empty to keep): ")
	if val != "" {
		if _, err := os.Stat(val); err != nil {
			fmt.Printf("[!] Warning: Wordlist file %q not found or inaccessible: %v\n", val, err)
		}
		im.cfg.WordlistPath = val
		fmt.Println("[+] Wordlist path updated.")
	}
}

func (im *InteractiveMenu) configureRequestSettings() {
	fmt.Println("\n[Request Settings]")
	fmt.Printf("1. HTTP Method (Current: %s)\n", im.cfg.Method)
	fmt.Printf("2. Concurrency Workers (Current: %d)\n", im.cfg.Concurrency)
	fmt.Printf("3. Rate Limit RPS (Current: %d)\n", im.cfg.RateLimit)
	fmt.Printf("4. Timeout Seconds (Current: %v)\n", im.cfg.Timeout)
	fmt.Printf("5. Follow Redirects (Current: %t)\n", im.cfg.FollowRedirects)
	fmt.Printf("6. Insecure TLS (Current: %t)\n", im.cfg.InsecureTLS)
	fmt.Printf("7. Custom Proxy (Current: %s)\n", im.cfg.ProxyURL)
	fmt.Printf("8. Custom Headers (Total: %d)\n", len(im.cfg.Headers))
	fmt.Printf("9. Request Body Data\n")

	choice := im.promptString("Choose setting to modify (or Enter to back): ")
	switch choice {
	case "1":
		m := im.promptString("HTTP Method [GET, HEAD, POST, PUT, DELETE, OPTIONS]: ")
		if m != "" {
			im.cfg.Method = strings.ToUpper(m)
		}
	case "2":
		cStr := im.promptString("Concurrency [1-500]: ")
		if c, err := strconv.Atoi(cStr); err == nil && c >= 1 && c <= 500 {
			im.cfg.Concurrency = c
		}
	case "3":
		rStr := im.promptString("Rate limit RPS (0 = unlimited): ")
		if r, err := strconv.Atoi(rStr); err == nil && r >= 0 {
			im.cfg.RateLimit = r
		}
	case "4":
		tStr := im.promptString("Timeout in seconds: ")
		if t, err := strconv.Atoi(tStr); err == nil && t > 0 {
			im.cfg.Timeout = time.Duration(t) * time.Second
		}
	case "5":
		red := im.promptString("Follow redirects? (y/n): ")
		im.cfg.FollowRedirects = strings.HasPrefix(strings.ToLower(red), "y")
	case "6":
		ins := im.promptString("Allow insecure TLS? (y/n, CAUTION): ")
		im.cfg.InsecureTLS = strings.HasPrefix(strings.ToLower(ins), "y")
		if im.cfg.InsecureTLS {
			validation.PrintInsecureTLSWarning()
		}
	case "7":
		px := im.promptString("Proxy URL (e.g. http://127.0.0.1:8080): ")
		im.cfg.ProxyURL = px
	case "8":
		h := im.promptString("Add Header 'Name: Value': ")
		if parts := strings.SplitN(h, ":", 2); len(parts) == 2 {
			im.cfg.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			fmt.Println("[+] Header added.")
		}
	case "9":
		b := im.promptString("Request Body data: ")
		im.cfg.Body = b
	}
}

func (im *InteractiveMenu) configureFilters() {
	fmt.Println("\n[Response Filters & Matchers]")
	fmt.Printf("1. Match Status Codes (Current: %s)\n", im.cfg.MatchCodesRaw)
	fmt.Printf("2. Filter Status Codes (Current: %s)\n", im.cfg.FilterCodesRaw)
	fmt.Printf("3. Match Response Sizes (Current: %s)\n", im.cfg.MatchSizesRaw)
	fmt.Printf("4. Filter Response Sizes (Current: %s)\n", im.cfg.FilterSizesRaw)
	fmt.Printf("5. Match Word Counts (Current: %s)\n", im.cfg.MatchWordsRaw)
	fmt.Printf("6. Filter Word Counts (Current: %s)\n", im.cfg.FilterWordsRaw)
	fmt.Printf("7. Match Line Counts (Current: %s)\n", im.cfg.MatchLinesRaw)
	fmt.Printf("8. Filter Line Counts (Current: %s)\n", im.cfg.FilterLinesRaw)

	choice := im.promptString("Choose filter to modify (or Enter to back): ")
	switch choice {
	case "1":
		im.cfg.MatchCodesRaw = im.promptString("Match Codes (e.g. 200,204,301,302): ")
	case "2":
		im.cfg.FilterCodesRaw = im.promptString("Filter Codes (e.g. 404,500): ")
	case "3":
		im.cfg.MatchSizesRaw = im.promptString("Match Sizes (e.g. 0,1024): ")
	case "4":
		im.cfg.FilterSizesRaw = im.promptString("Filter Sizes: ")
	case "5":
		im.cfg.MatchWordsRaw = im.promptString("Match Words: ")
	case "6":
		im.cfg.FilterWordsRaw = im.promptString("Filter Words: ")
	case "7":
		im.cfg.MatchLinesRaw = im.promptString("Match Lines: ")
	case "8":
		im.cfg.FilterLinesRaw = im.promptString("Filter Lines: ")
	}
	_ = im.cfg.CompileCriteria()
	fmt.Println("[+] Filters compiled.")
}

func (im *InteractiveMenu) configureOutput() {
	fmt.Println("\n[Output Settings]")
	fmt.Printf("Current Output File: %s\n", im.cfg.OutputFile)
	fmt.Printf("Current Format: %s (text, json, csv, md)\n", im.cfg.OutputFormat)
	fmt.Printf("Quiet Mode: %t | Verbose Mode: %t\n", im.cfg.Quiet, im.cfg.Verbose)

	f := im.promptString("Output File (leave empty for terminal only): ")
	if f != "" {
		im.cfg.OutputFile = f
	}
	fmtStr := im.promptString("Output Format [text, json, csv, md]: ")
	if fmtStr != "" {
		im.cfg.OutputFormat = strings.ToLower(fmtStr)
	}
}

func (im *InteractiveMenu) viewCurrentConfig() {
	fmt.Println("\n================= CURRENT CONFIGURATION =================")
	fmt.Printf("  Target URL        : %s\n", im.cfg.TargetURL)
	fmt.Printf("  Placeholder       : %s\n", im.cfg.Placeholder)
	fmt.Printf("  Wordlist Path     : %s\n", im.cfg.WordlistPath)
	fmt.Printf("  Method            : %s\n", im.cfg.Method)
	fmt.Printf("  Concurrency       : %d\n", im.cfg.Concurrency)
	fmt.Printf("  Rate Limit        : %d req/sec\n", im.cfg.RateLimit)
	fmt.Printf("  Timeout           : %v\n", im.cfg.Timeout)
	fmt.Printf("  Follow Redirects  : %t\n", im.cfg.FollowRedirects)
	fmt.Printf("  Insecure TLS      : %t\n", im.cfg.InsecureTLS)
	fmt.Printf("  Proxy             : %s\n", im.cfg.ProxyURL)
	fmt.Printf("  Match Codes       : %s\n", im.cfg.MatchCodesRaw)
	fmt.Printf("  Filter Codes      : %s\n", im.cfg.FilterCodesRaw)
	fmt.Printf("  Output File       : %s (%s)\n", im.cfg.OutputFile, im.cfg.OutputFormat)
	fmt.Println("=========================================================")
}

func (im *InteractiveMenu) displayHelp() {
	fmt.Println("\n--- [ 8WHIE WebProbe Help & Quick Guide ] ---")
	fmt.Println("Usage example in CLI:")
	fmt.Println("  webprobe -u https://target.local/FUZZ -w wordlist.txt -mc 200,301,302")
	fmt.Println("Options:")
	fmt.Println("  -u, --url          Target URL containing FUZZ")
	fmt.Println("  -w, --wordlist     Candidate path dictionary")
	fmt.Println("  -c, --concurrency  Worker concurrency (default 20)")
	fmt.Println("  -rate              Requests per second limit")
	fmt.Println("  -mc / -fc          Match / Filter HTTP status codes")
	fmt.Println("  -o, --output       Save findings to file")
	fmt.Println("  -of, --format      Output format (text, json, csv, md)")
	fmt.Println("----------------------------------------------")
}

func (im *InteractiveMenu) displayAbout() {
	fmt.Println("\n--- [ About 8WHIE WebProbe ] ---")
	fmt.Println("  Project   : 8WHIE WebProbe")
	fmt.Println("  Author    : 8WHIE")
	fmt.Println("  Version   : 1.0.0")
	fmt.Println("  License   : MIT License (Copyright (c) 2026 8WHIE)")
	fmt.Println("  Platforms : Android • Termux • Linux • macOS • Windows")
	fmt.Println("  YouTube   : 8WHIE")
	fmt.Println("  Instagram : @aaynkt")
	fmt.Println("  Telegram  : @Arnxkt")
	fmt.Println("  Channel   : https://t.me/whiee")
	fmt.Println("---------------------------------")
}

func (im *InteractiveMenu) runDiscovery(ctx context.Context) {
	if err := validation.ValidateExecution(im.cfg); err != nil {
		fmt.Printf("\n[!] Cannot start discovery: %v\n", err)
		return
	}

	if im.cfg.InsecureTLS {
		validation.PrintInsecureTLSWarning()
	}

	w, err := output.NewWriter(im.cfg)
	if err != nil {
		fmt.Printf("\n[!] Output error: %v\n", err)
		return
	}

	eng, err := engine.NewEngine(im.cfg, w)
	if err != nil {
		fmt.Printf("\n[!] Engine error: %v\n", err)
		return
	}

	w.PrintBanner()
	stats, err := eng.Run(ctx)
	if err != nil && ctx.Err() == nil {
		fmt.Printf("\n[!] Discovery terminated with error: %v\n", err)
	}

	w.PrintSummary(stats)
	_ = w.Close(stats)
}
