# 8WHIE WebProbe — Security & Authorized-Use Policy

## 1. Authorized Use Policy

`8WHIE WebProbe` is engineered strictly as an HTTP request and response diagnostic tool for:
- Systems and domains you own.
- Systems where you have explicit written authorization (e.g. penetration testing agreement, bug bounty scope, internal organization vulnerability management).

**Prohibited Activities:**
- Probing third-party servers without prior authorization.
- Attempting denial-of-service by setting unreasonable concurrency against targets.
- Circumventing security measures on unauthorized infrastructure.

The authors and contributors of `8WHIE WebProbe` assume no liability and are not responsible for any misuse or damage caused by this utility.

---

## 2. Safety Design & Secure Defaults

- **No Malicious Payloads**: WebProbe does not generate exploits, shellcodes, or malicious intrusion sequences. It simply sends HTTP requests with words substituted from the operator's chosen wordlist.
- **Strict TLS by Default**: Insecure certificate skipping is disabled by default. If enabled via `-k`, a prominent warning is emitted.
- **Controlled Concurrency**: Concurrency is strictly bounded and rate limits are available to prevent unintended server saturation.
- **Safe Memory Consumption**: Response bodies are read with bounded buffers (5MB max) to eliminate memory exhaustion attacks from malicious target responses.
- **No Secret Leaking**: Headers and authorization tokens provided in CLI flags or configuration are never logged to public telemetry or external servers.

---

## 3. Reporting Vulnerabilities

If you discover a security issue or vulnerability within `8WHIE WebProbe` itself, please contact:
- **Developer**: 8WHIE
- **Telegram**: [@Arnxkt](https://t.me/Arnxkt)
- **Instagram**: [@aaynkt](https://instagram.com/aaynkt)
- **Telegram Channel**: [https://t.me/whiee](https://t.me/whiee)
