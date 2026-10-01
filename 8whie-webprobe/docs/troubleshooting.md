# 8WHIE WebProbe — Troubleshooting & Diagnostics

Common questions and resolutions when running `webprobe`.

### 1. `[!] Validation error: no placeholder token "FUZZ" discovered`

**Cause:** The URL provided in `-u` does not include the replacement keyword `FUZZ`, and neither the body nor headers include it.  
**Fix:** Specify where candidate words should be inserted:
```bash
webprobe -u https://example.com/FUZZ -w wordlist.txt
```
If you prefer a custom token (e.g., `{TARGET}`), supply `-p "{TARGET}"`:
```bash
webprobe -u https://example.com/{TARGET} -w wordlist.txt -p "{TARGET}"
```

---

### 2. High Network Errors or Timeouts

**Cause:** The target server or network gateway is being overwhelmed by too many simultaneous connections, or a firewall is dropping packets.  
**Fix:**
- Reduce concurrency from 20 to 5 or 2: `-c 5`
- Enforce a rate limit (requests per second): `-rate 10`
- Increase the timeout duration: `-t 20`
- Enable verbose diagnostics: `-v` to inspect exact network failures.

---

### 3. `[!] CAUTION: TLS certificate validation is disabled`

**Cause:** You passed `-k` / `--insecure`.  
**Note:** This is intended solely for internal lab environments or self-signed staging servers. Do not use this when testing over untrusted public networks as communications could be intercepted.

---

### 4. Excessive Findings / Everything Returns 200 OK

**Cause:** The target website is a Single Page Application (SPA) or has a wildcard catch-all route configured that returns HTTP 200 for all paths.  
**Fix:**
1. Check the response size of an arbitrary invalid path (e.g. `/thispathdoesnotexist12345`).
2. If it returns 200 with 1,540 bytes, filter out that specific byte size:
   ```bash
   webprobe -u https://target.local/FUZZ -w wordlist.txt -fs 1540
   ```
3. Or filter by line/word counts (`-fl` or `-fw`).

---

### 5. Running on Android / Termux

**Setup in Termux:**
```bash
pkg update && pkg install golang git
git clone https://github.com/8whie/8whie-webprobe.git
cd 8whie-webprobe
make build
./bin/webprobe -h
```
Ensure you have network permissions enabled for the Termux app.
