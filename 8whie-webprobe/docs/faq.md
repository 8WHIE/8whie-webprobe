# 8WHIE WebProbe — Frequently Asked Questions (FAQ)

### What is 8WHIE WebProbe?
`8WHIE WebProbe` is an original open-source command-line tool written in Go for controlled HTTP discovery, endpoint enumeration, request pacing, and response filtering. It was independently created by **8WHIE** for legitimate web security testing, penetration testing audits, and developer diagnostics.

### What is HTTP endpoint discovery?
HTTP endpoint discovery (also referred to as content discovery or path enumeration) is the practice of methodically checking an HTTP service against a dictionary list of standard path names, file extensions, or API routes to discover exposed endpoints, configuration files, staging directories, or legacy scripts.

### How do I test a website I own?
1. Ensure your web application is accessible (e.g., `https://my-authorized-site.com`).
2. Run WebProbe with a wordlist:
   ```bash
   webprobe -u https://my-authorized-site.com/FUZZ -w wordlist.txt -mc 200,301,302
   ```
3. Analyze the results to ensure sensitive internal admin panels or debugging files are not unintentionally exposed to the public.

### Does WebProbe require root privileges?
No. WebProbe executes as a standard unprivileged user space binary. It does not open raw sockets or require administrative permissions.

### Does WebProbe work on Linux?
Yes. WebProbe is natively compiled for Linux on AMD64, ARM64, and 32-bit ARM architectures.

### Does WebProbe work on Android / Termux?
Yes! WebProbe compiles and runs smoothly inside Android Termux. You can install Go inside Termux (`pkg install golang git`), clone the repository, run `make build`, and execute `webprobe`.

### Can I use custom headers?
Yes. Use the `-H` flag (can be repeated multiple times). For example:
```bash
webprobe -u https://target.local/FUZZ -w paths.txt -H "Authorization: Bearer SECRET_TOKEN" -H "X-Audit: Authorized"
```
**Safety Warning:** Never paste private production credentials into shell histories or public scripts.

### Can I limit request speed?
Yes. Use `-rate <N>` to enforce a maximum of `N` requests per second. You can also adjust worker concurrency via `-c <N>`.

### Can I filter HTTP responses?
Yes. You can match or filter by status codes (`-mc`, `-fc`), byte sizes (`-ms`, `-fs`), word counts (`-mw`, `-fw`), and line counts (`-ml`, `-fl`).

### Can results be exported to JSON?
Yes. Specify `-o results.json -of json`. WebProbe also supports CSV (`-of csv`), Markdown (`-of md`), and standard text (`-of text`).

### Is WebProbe affiliated with Kali Linux?
No. WebProbe is an independent open-source project by 8WHIE. It can be installed and used on Kali Linux, Debian, Ubuntu, Fedora, Arch, Alpine, Termux, macOS, or Windows.

### Is WebProbe affiliated with ffuf?
No. WebProbe is a completely new, independently designed and authored project. It does not share source code, functions, design algorithms, or project structure with ffuf.

### Can I use WebProbe against public websites?
Only if you have explicit, documented permission or ownership of the target. Scanning systems without permission is unlawful in many jurisdictions and violates terms of service. Always test responsibly.
