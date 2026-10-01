# 8WHIE WebProbe — Output Formats & Reporting

`8WHIE WebProbe` is designed for both direct human analysis in the terminal and machine-readable data pipelines (CI/CD pipelines, automated vulnerability scanners, audit archives).

## 1. Terminal Output (`text`)

The default mode produces a high-contrast tabular stream:

```text
STATUS   SIZE       WORDS    LINES    DURATION   PAYLOAD                   REDIRECT
--------------------------------------------------------------------------------
200      3512       450      85       42ms       index.html                
301      185        12       4        28ms       admin                     -> https://target.local/admin/
403      280        22       8        35ms       .git                      
```

### Color Coding Standards

- **2xx (Success)**: Bold Green
- **3xx (Redirect)**: Cyan
- **401 / 403 (Client Restriction)**: Yellow
- **5xx (Server Error)**: Red

Disable colors using `--no-color` for plain terminal output or piping.

---

## 2. JSON Output (`json`)

Structured JSON output provides comprehensive metadata, runtime statistics, and an array of discovered findings.

```bash
webprobe -u https://target.local/FUZZ -w paths.txt -o results.json -of json
```

Example JSON structure:

```json
{
  "tool": "8WHIE WebProbe",
  "version": "1.0.0",
  "target": "https://target.local/FUZZ",
  "wordlist": "paths.txt",
  "statistics": {
    "total_candidates": 500,
    "sent_requests": 500,
    "matched_results": 4,
    "filtered_results": 496,
    "network_errors": 0,
    "duration": 4250000000,
    "requests_per_second": 117.6
  },
  "results": [
    {
      "request": {
        "id": 14,
        "url": "https://target.local/admin",
        "method": "GET",
        "payload": "admin"
      },
      "response": {
        "status_code": 403,
        "content_length": 280,
        "word_count": 22,
        "line_count": 8,
        "duration_ms": 35000000
      },
      "matched": true
    }
  ]
}
```

---

## 3. CSV Output (`csv`)

Standard comma-separated format for importing into spreadsheets or database tables:

```bash
webprobe -u https://target.local/FUZZ -w paths.txt -o results.csv -of csv
```

Header layout:
```csv
ID,Payload,URL,Status,Size,Words,Lines,DurationMs,Redirect
14,admin,https://target.local/admin,403,280,22,8,35,
```

---

## 4. Markdown Output (`md`)

Creates a GitHub-flavored Markdown table, ideal for automated security audit reports and documentation.

```bash
webprobe -u https://target.local/FUZZ -w paths.txt -o report.md -of md
```
