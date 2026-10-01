# 8WHIE WebProbe

```text
  ___ _ _ _ _  _ ___   _ _ _     _   ___          _         
 ( _ ) | | | || |_ _| | | | |___| |_| _ \_ _ ___ | |__  ___ 
 / _ \_  _ | __ || |  | | | / -_) '_ \  _/ '_/ _ \| '_ \/ -_)
 \___/ |_| |_||_|___| |_____/\___|_,__/_| |_| \___/|_.__/\___|
```

> **8WHIE WebProbe** — An original open-source CLI utility for HTTP endpoint discovery, content enumeration, request analysis, and web security research workflows.
>
> *Android • Termux • Linux • macOS • Windows • Security Research*

---

## ⚡ Overview

**8WHIE WebProbe** is a fast, concurrent command-line utility designed for legitimate security researchers, system administrators, and web application developers. It provides controlled HTTP endpoint discovery, URL path enumeration, request pacing, and multi-parameter response analysis against systems you own or are explicitly authorized to test.

Built entirely in idiomatic Go with zero external runtime dependencies, 8WHIE WebProbe offers high throughput, minimal memory overhead, rate limiting, and structured output formats.

---

## 🔒 Authorized-Use Policy & Legal Warning

> **CRITICAL LEGAL NOTICE:**  
> 8WHIE WebProbe is provided strictly for lawful security research, defensive auditing, and authorized diagnostic testing.  
> Performing unauthorized scans, penetration tests, or discovery against network services without explicit written permission from the system owner is illegal and unethical. The developer assumes no responsibility for unauthorized misuse.

---

## ✨ Features

- **High-Performance Worker Pool**: Fully bounded concurrency (1–500 workers) with graceful shutdown handling (SIGINT/SIGTERM).
- **Flexible URL & Request Placeholders**: Dynamic keyword replacement (`FUZZ` by default, customizable via `-p`) in URLs, headers, or request bodies.
- **Granular Response Filtering**: Match or filter findings by HTTP status code (`-mc` / `-fc`), response byte length (`-ms` / `-fs`), word count (`-mw` / `-fw`), and line count (`-ml` / `-fl`).
- **Precision Rate Limiting**: Built-in request-per-second (`-rate`) throttling to safeguard fragile backend services and staging environments.
- **Multiple Output Formats**: Human-readable colored terminal output, JSON (`-of json`), CSV (`-of csv`), and GitHub-flavored Markdown (`-of md`).
- **Interactive Terminal Console**: Full interactive menu (`webprobe -i`) for guided target configuration and probe execution.
- **HTTP/SOCKS Proxy Support**: Configurable proxy routing (`-x`) for auditing traffic through diagnostic proxies.
- **Controlled Redirect Handling**: Optional redirect following (`-r`) with automatic loop avoidance.
- **Safe TLS Defaults**: Strict certificate validation by default, with explicit warnings when testing internal staging servers with insecure certificates (`-k`).
- **Pure Go Implementation**: Zero external C-library dependencies; compiles to a single standalone static binary.

---

## 💻 Supported Platforms

- **Linux**: AMD64, ARM64, 32-bit ARM (Ubuntu, Debian, Kali Linux, Arch, Fedora, Alpine)
- **Android / Termux**: ARMv7, ARM64
- **macOS**: Apple Silicon (M1/M2/M3/M4) and Intel x86_64
- **Windows**: Windows 10/11, Windows Server (x86_64)

---

## 📦 Requirements

- **Go**: 1.21, 1.22, or 1.23 (if compiling from source)
- **Make** (optional, for convenient build targets)

---

## 🚀 Installation & Build

### 1. Build from Source

```bash
# Clone the repository
git clone https://github.com/8whie/8whie-webprobe.git
cd 8whie-webprobe

# Build binary using Make
make build

# Verify build
./bin/webprobe -V
```

### 2. Install to Go Path

```bash
go install ./cmd/webprobe
```

### 3. Termux on Android Setup

```bash
pkg update -y && pkg install -y golang git make
git clone https://github.com/8whie/8whie-webprobe.git
cd 8whie-webprobe
make build
./bin/webprobe -h
```

---

## 📖 Quick Start & Basic Usage

### Basic Path Discovery
```bash
webprobe -u https://example.com/FUZZ -w /path/to/wordlist.txt
```

### Match Specific Status Codes
```bash
webprobe -u https://example.com/FUZZ -w wordlist.txt -mc 200,301,302,403
```

### Filter Out Cluttered 404s and 500s
```bash
webprobe -u https://example.com/FUZZ -w wordlist.txt -fc 404,500
```

### Rate-Limited Diagnostic Probe
```bash
webprobe -u https://target.local/FUZZ -w wordlist.txt -c 5 -rate 10 -t 15
```

### Custom Headers & Authorization
```bash
webprobe -u https://target.local/FUZZ -w paths.txt \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Audit-Team: SecurityOps"
```

### Export Results to JSON
```bash
webprobe -u https://target.local/FUZZ -w wordlist.txt -o findings.json -of json
```

### Launch Interactive Console
```bash
webprobe -i
```

---

## 🛠️ Complete CLI Options Reference

```text
TARGET & WORDLIST:
  -u, --url string         Target URL with placeholder (e.g. https://target.local/FUZZ)
  -w, --wordlist string    Path to candidate wordlist dictionary
  -p, --placeholder string Placeholder token (default "FUZZ")

HTTP REQUEST CONTROLS:
  -X, --method string      HTTP method (GET, HEAD, POST, PUT, DELETE, OPTIONS) (default "GET")
  -H, --header string      Custom HTTP header (repeatable)
  -d, --data string        Explicit HTTP request body data
  -r, --redirects          Follow HTTP redirects (up to 10 hops)
  -t, --timeout int        HTTP request timeout in seconds (default 10)
  -x, --proxy string       Proxy server address (e.g. http://127.0.0.1:8080)
  -k, --insecure           Allow untrusted TLS certificates (warning: insecure)

CONCURRENCY & PACING:
  -c, --concurrency int    Concurrent worker count (1 to 500) (default 20)
  -rate int                Enforce rate limit in requests per second (0 = unrestricted)

MATCHING & FILTERING:
  -mc string               Match HTTP status codes (e.g. 200,204,301,302)
  -fc string               Filter out HTTP status codes (default "404")
  -ms string               Match response byte sizes (e.g. 512,1024)
  -fs string               Filter out response byte sizes
  -mw string               Match response word counts
  -fw string               Filter out response word counts
  -ml string               Match response line counts
  -fl string               Filter out response line counts

OUTPUT & REPORTING:
  -o, --output string      Destination file for findings
  -of, --format string     Output format: text, json, csv, md (default "text")
  -q, --quiet              Quiet mode: output findings only
  -v, --verbose            Verbose diagnostics: network warnings
  --no-color               Disable ANSI terminal color output

MODES:
  -i, --interactive        Start interactive guided terminal interface
  -V, --version            Display version and author information
  -h, --help               Display full help manual
```

---

## 🧪 Testing & Validation

The codebase includes comprehensive unit tests and an integration test utilizing an isolated Go mock HTTP test server (`httptest.NewServer`):

```bash
# Run unit & integration tests
go test -v -race ./tests/...

# Run Go static analysis
go vet ./...

# Verify code formatting
gofmt -l .
```

---

## 📚 Documentation Index

Detailed documentation guides are available in the [`docs/`](docs/) directory:
- [Usage Guide](docs/usage.md)
- [Configuration Reference](docs/configuration.md)
- [Filtering & Matching Strategies](docs/filtering.md)
- [Output Formats & Reporting](docs/output.md)
- [Troubleshooting & FAQ](docs/troubleshooting.md)
- [Authorized Use & Security Policy](docs/security.md)
- [Frequently Asked Questions (FAQ)](docs/faq.md)

---

## 🤝 Contribution Guidelines

Contributions are welcome! Please follow these standards:
1. Ensure all code adheres to standard `gofmt` style.
2. Verify all tests pass cleanly (`go test -v ./tests/...`).
3. Maintain zero external runtime dependencies where practical.
4. Open a pull request describing your improvements clearly.

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

Copyright (c) 2026 **8WHIE**

---

## 🌐 Brand & Developer Information

Developed by **8WHIE**

- **YouTube**: [8WHIE](https://youtube.com)
- **Instagram**: [@aaynkt](https://instagram.com/aaynkt)
- **Telegram**: [@Arnxkt](https://t.me/Arnxkt)
- **Telegram Channel**: [https://t.me/whiee](https://t.me/whiee)
- **GitHub**: [8whie](https://github.com/8whie)
