# 8WHIE WebProbe — Configuration Reference

This guide covers all command-line flags and runtime parameters available in `8WHIE WebProbe`.

| Flag | Long Flag | Default | Description |
| :--- | :--- | :--- | :--- |
| `-u` | `--url` | *Required* | Target URL with placeholder (e.g., `https://target.local/FUZZ`) |
| `-w` | `--wordlist` | *Required* | Path to candidate wordlist dictionary |
| `-p` | `--placeholder` | `FUZZ` | Custom placeholder token |
| `-X` | `--method` | `GET` | HTTP Request Method (`GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `OPTIONS`) |
| `-H` | `--header` | `none` | Custom header `'Key: Value'` (can be repeated) |
| `-d` | `--data` | `none` | Explicit HTTP request body data |
| `-c` | `--concurrency` | `20` | Worker pool concurrency count (1 to 500) |
| `-rate` | `--rate-limit` | `0` | Rate limit in requests per second (0 = unrestricted) |
| `-t` | `--timeout` | `10` | Request timeout in seconds |
| `-r` | `--redirects` | `false` | Follow HTTP redirects (up to 10 hops) |
| `-k` | `--insecure` | `false` | Allow untrusted TLS certificates (warning: insecure) |
| `-x` | `--proxy` | `none` | Forward traffic through HTTP/SOCKS proxy |
| `-mc` | `--match-code` | `none` | Match specific HTTP status codes (e.g. `200,301,302`) |
| `-fc` | `--filter-code` | `404` | Filter out specific HTTP status codes |
| `-ms` | `--match-size` | `none` | Match response byte sizes |
| `-fs` | `--filter-size` | `none` | Filter out response byte sizes |
| `-mw` | `--match-words` | `none` | Match response word counts |
| `-fw` | `--filter-words` | `none` | Filter out response word counts |
| `-ml` | `--match-lines` | `none` | Match response line counts |
| `-fl` | `--filter-lines` | `none` | Filter out response line counts |
| `-o` | `--output` | `none` | Destination file for results |
| `-of` | `--format` | `text` | Output format (`text`, `json`, `csv`, `md`) |
| `-q` | `--quiet` | `false` | Quiet mode: prints only payload or findings |
| `-v` | `--verbose` | `false` | Print diagnostic messages and network warnings |
| `--no-color` | | `false` | Disable ANSI terminal color codes |
| `-i` | `--interactive`| `false` | Launch interactive menu |
| `-V` | `--version` | | Display version and author information |
| `-h` | `--help` | | Display command manual |
