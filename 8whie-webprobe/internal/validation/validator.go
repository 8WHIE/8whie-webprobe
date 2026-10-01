package validation

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/8whie/8whie-webprobe/internal/config"
)

// ValidateExecution inspects all runtime parameters to guarantee safe, valid operation.
func ValidateExecution(cfg *config.Config) error {
	if cfg.TargetURL == "" {
		return errors.New("target URL (-u) is required; run 'webprobe -h' for usage instructions")
	}

	parsedURL, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return fmt.Errorf("malformed target URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme %q; target must begin with http:// or https://", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return errors.New("target URL must contain a valid hostname or IP address")
	}

	placeholder := cfg.Placeholder
	if placeholder == "" {
		placeholder = "FUZZ"
	}

	hasPlaceholderInURL := strings.Contains(cfg.TargetURL, placeholder)
	hasPlaceholderInBody := strings.Contains(cfg.Body, placeholder)
	hasPlaceholderInHeader := false
	for k, v := range cfg.Headers {
		if strings.Contains(k, placeholder) || strings.Contains(v, placeholder) {
			hasPlaceholderInHeader = true
			break
		}
	}

	if !hasPlaceholderInURL && !hasPlaceholderInBody && !hasPlaceholderInHeader {
		return fmt.Errorf("no placeholder token %q discovered in target URL, request body, or headers", placeholder)
	}

	if cfg.WordlistPath == "" {
		return errors.New("wordlist path (-w) is required; provide a valid dictionary file path")
	}

	info, err := os.Stat(cfg.WordlistPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("wordlist file %q does not exist", cfg.WordlistPath)
		}
		return fmt.Errorf("unable to access wordlist file %q: %w", cfg.WordlistPath, err)
	}

	if info.IsDir() {
		return fmt.Errorf("wordlist path %q is a directory, expected a regular file", cfg.WordlistPath)
	}

	if cfg.ProxyURL != "" {
		proxyParsed, err := url.Parse(cfg.ProxyURL)
		if err != nil || (proxyParsed.Scheme != "http" && proxyParsed.Scheme != "https" && proxyParsed.Scheme != "socks5") {
			return fmt.Errorf("invalid proxy URL %q; must be http://, https://, or socks5://", cfg.ProxyURL)
		}
	}

	if err := cfg.ValidateSanity(); err != nil {
		return err
	}

	return nil
}

// PrintSecurityNotice writes the authorization advisory to the designated output.
func PrintSecurityNotice(quiet bool) {
	if quiet {
		return
	}
	fmt.Fprintln(os.Stderr, "================================================================================")
	fmt.Fprintln(os.Stderr, "[!] AUTHORIZED USE NOTICE: 8WHIE WebProbe is designed exclusively for authorized")
	fmt.Fprintln(os.Stderr, "    security research, diagnostic testing, and system verification on targets")
	fmt.Fprintln(os.Stderr, "    you own or have explicit written authorization to test.")
	fmt.Fprintln(os.Stderr, "================================================================================")
}

// PrintInsecureTLSWarning alerts the operator when TLS verification is bypassed.
func PrintInsecureTLSWarning() {
	fmt.Fprintln(os.Stderr, "[!] CAUTION: TLS certificate validation is disabled (-k/--insecure).")
	fmt.Fprintln(os.Stderr, "    Traffic can be intercepted or manipulated by network intermediaries.")
}
