# 8WHIE WebProbe — Filtering & Matching Guide

Effective web security research requires isolating meaningful signals from repetitive noise. `8WHIE WebProbe` provides granular, multi-dimensional response filtering.

## Overview of Matchers vs. Filters

- **Matchers (`-m*`)**: Explicitly whitelist responses that satisfy specified criteria. Any response not meeting the match condition is discarded.
- **Filters (`-f*`)**: Explicitly blacklist responses that satisfy specified criteria. Any response meeting the filter condition is discarded.

When both are used simultaneously, **both** conditions must be satisfied: the response must match the whitelisted criteria and not match any filtered criteria.

---

## 1. Status Code Filtering

Default behavior filters out status code `404` (`-fc 404`).

### Matching Specific Status Codes

To view only successful responses (200 OK) or redirects:
```bash
webprobe -u https://target.local/FUZZ -w paths.txt -mc 200,301,302
```

### Filtering Out Clutter

If the target web server returns generic `500` server errors or custom `403` forbidden pages that clutter the output:
```bash
webprobe -u https://target.local/FUZZ -w paths.txt -fc 404,500,403
```

---

## 2. Response Size Filtering

When web servers respond with custom default pages or generic "Under Construction" banners that all share an identical byte length:

```bash
# Filter out responses that are exactly 4528 bytes:
webprobe -u https://target.local/FUZZ -w paths.txt -fs 4528

# Filter out multiple static error page lengths:
webprobe -u https://target.local/FUZZ -w paths.txt -fs 0,142,4528
```

---

## 3. Word Count and Line Count Filtering

Dynamic pages may change byte size due to timestamps or randomized session identifiers, but keep an identical word or line count.

```bash
# Filter out pages with 42 words:
webprobe -u https://target.local/FUZZ -w paths.txt -fw 42

# Filter out pages with 1 single line:
webprobe -u https://target.local/FUZZ -w paths.txt -fl 1
```

---

## Recommended Research Workflow

1. Send an intentional non-existent path probe:
   ```bash
   curl -s -i https://target.local/this-endpoint-does-not-exist-8whie
   ```
2. Inspect the HTTP status code, Content-Length, words, and lines.
3. Configure `webprobe` with `-fs`, `-fw`, or `-fc` matching that baseline non-existent response.
4. Run your probe with high precision.
