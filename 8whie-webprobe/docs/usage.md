# 8WHIE WebProbe — Usage Guide

`8WHIE WebProbe` is an original open-source command-line tool built for HTTP endpoint discovery, URL enumeration, and request response analysis.

## Operational Prerequisites

Ensure you have:
1. Written permission or ownership of the target domain/IP.
2. A dictionary wordlist (e.g., paths, routes, IDs, filenames).
3. Network connectivity to the target host.

## Basic Syntax

```bash
webprobe -u <TARGET_URL> -w <WORDLIST> [OPTIONS]
```

Where `<TARGET_URL>` contains the placeholder keyword `FUZZ` by default.

### 1. Basic Path Discovery

```bash
webprobe -u https://example.com/FUZZ -w wordlist.txt
```

WebProbe substitutes each line from `wordlist.txt` into the position of `FUZZ`, issues controlled HTTP requests, and prints non-404 status codes.

### 2. Targeting File Extensions

```bash
webprobe -u https://example.com/api/v1/FUZZ.json -w endpoints.txt
```

### 3. Interactive Console Mode

Launch the interactive text-based console menu to configure and run probes step by step:

```bash
webprobe -i
```

## Advanced Workflows

### Subdomain and Host Header Discovery

WebProbe allows placing the `FUZZ` placeholder inside headers:

```bash
webprobe -u https://192.168.1.50/ -H "Host: FUZZ.internal.corp" -w subdomains.txt -mc 200,302
```

### API Endpoint Discovery with POST Payloads

```bash
webprobe -u https://api.internal.local/v2/items \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"action": "lookup", "target": "FUZZ"}' \
  -w candidates.txt \
  -mc 200
```

### Rate-Limited Diagnostic Scanning

When auditing production or staging environments that require gentle network traffic:

```bash
webprobe -u https://target.local/FUZZ -w paths.txt -c 5 -rate 10 -t 15
```
This restricts execution to 5 concurrent workers and at most 10 requests per second.
