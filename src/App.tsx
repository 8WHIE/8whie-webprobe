import React, { useState, useEffect, useMemo, useRef } from 'react';
import { 
  Terminal, Shield, Play, RotateCcw, Copy, Check, Download, 
  Settings, BookOpen, HelpCircle, FileText, ExternalLink, 
  Filter, Cpu, Globe, Sliders, ChevronDown, ChevronRight,
  AlertTriangle, Layers, Code, Zap
} from 'lucide-react';

interface ProbeSimResult {
  id: number;
  payload: string;
  url: string;
  statusCode: number;
  contentLength: number;
  wordCount: number;
  lineCount: number;
  durationMs: number;
  redirectUrl?: string;
}

const SAMPLE_WORDLISTS: Record<string, string[]> = {
  api: [
    'api', 'api/v1', 'api/v2', 'health', 'healthz', 'metrics',
    'graphql', 'swagger', 'swagger.json', 'openapi.json', 'docs',
    'auth', 'oauth/token', 'users', 'status', 'ping', 'version',
    'debug', 'trace', 'info', 'schema'
  ],
  admin: [
    'admin', 'administrator', 'panel', 'dashboard', 'console',
    'manage', 'login', 'portal', 'control', 'superadmin',
    'backend', 'auth/login', 'account', 'signin', 'access'
  ],
  config: [
    'config.json', 'settings.json', 'package.json', 'docker-compose.yml',
    'web.config', '.env.example', 'robots.txt', 'sitemap.xml',
    'server.key', 'backup.zip', 'dump.sql', 'database.yml'
  ],
  routes: [
    'about', 'contact', 'terms', 'privacy', 'blog', 'support',
    'static', 'assets', 'images', 'upload', 'files', 'media',
    'downloads', 'releases', 'changelog', 'forum', 'faq'
  ]
};

const DOC_SECTIONS = [
  {
    id: 'usage',
    title: 'Usage Guide',
    content: `### Getting Started with 8WHIE WebProbe

8WHIE WebProbe substitutes the placeholder token (\`FUZZ\` by default) in your target URL, request body, or headers with words from your chosen dictionary file.

#### Basic Syntax
\`\`\`bash
webprobe -u <TARGET_URL> -w <WORDLIST> [OPTIONS]
\`\`\`

#### Quick Examples

1. **Standard Endpoint Discovery:**
\`\`\`bash
webprobe -u https://target.local/FUZZ -w paths.txt
\`\`\`

2. **Match Specific HTTP Status Codes:**
\`\`\`bash
webprobe -u https://target.local/FUZZ -w paths.txt -mc 200,301,302,403
\`\`\`

3. **Gentle Rate-Limited Scan:**
\`\`\`bash
webprobe -u https://target.local/FUZZ -w paths.txt -c 5 -rate 10 -t 15
\`\`\`

4. **Interactive Console Mode:**
\`\`\`bash
webprobe -i
\`\`\`
`
  },
  {
    id: 'config',
    title: 'CLI Configuration Reference',
    content: `### Complete Command Flags

| Flag | Long Flag | Default | Description |
| :--- | :--- | :--- | :--- |
| \`-u\` | \`--url\` | *Required* | Target URL with placeholder token (e.g. https://target.local/FUZZ) |
| \`-w\` | \`--wordlist\` | *Required* | Path to candidate dictionary wordlist |
| \`-p\` | \`--placeholder\` | \`FUZZ\` | Custom replacement placeholder string |
| \`-X\` | \`--method\` | \`GET\` | HTTP request method (GET, HEAD, POST, PUT, DELETE, OPTIONS) |
| \`-H\` | \`--header\` | *None* | Custom header 'Key: Value' (repeatable) |
| \`-d\` | \`--data\` | *None* | Explicit HTTP request body data |
| \`-c\` | \`--concurrency\` | \`20\` | Concurrent worker goroutines (1 to 500) |
| \`-rate\` | \`--rate-limit\` | \`0\` | Rate limit in requests per second (0 = unrestricted) |
| \`-t\` | \`--timeout\` | \`10\` | HTTP request timeout in seconds |
| \`-r\` | \`--redirects\` | \`false\` | Follow HTTP redirects (up to 10 hops) |
| \`-k\` | \`--insecure\` | \`false\` | Skip TLS certificate verification (shows warning) |
| \`-x\` | \`--proxy\` | *None* | Diagnostic HTTP/SOCKS proxy URL |
| \`-o\` | \`--output\` | *None* | Destination file path to save findings |
| \`-of\` | \`--format\` | \`text\` | Output format: text, json, csv, md |
| \`-q\` | \`--quiet\` | \`false\` | Quiet mode: prints only findings |
| \`-v\` | \`--verbose\` | \`false\` | Verbose diagnostics: prints network errors |
| \`-i\` | \`--interactive\` | \`false\` | Launch interactive terminal interface |
| \`-V\` | \`--version\` | | Display version, brand, and social links |
| \`-h\` | \`--help\` | | Display comprehensive command manual |
`
  },
  {
    id: 'filtering',
    title: 'Filtering & Matching',
    content: `### Response Filtering Engine

Isolate valid endpoints and eliminate repetitive false positives using multi-dimensional criteria:

#### Status Code Matching & Filtering
- \`-mc 200,204,301,302\`: Only display findings returning these exact HTTP status codes.
- \`-fc 404,500\`: Discard all findings returning 404 Not Found or 500 Server Error (default: \`-fc 404\`).

#### Response Size Filtering
- \`-fs 4528\`: Filter out generic customized error pages that share a static size of 4,528 bytes.
- \`-ms 1024,2048\`: Match responses of exact known size thresholds.

#### Word & Line Count Filtering
- \`-fw 12\`: Filter responses containing exactly 12 words.
- \`-fl 1\`: Filter single-line responses.
`
  },
  {
    id: 'output',
    title: 'Output & Reporting Formats',
    content: `### Machine-Readable Data Pipelines

8WHIE WebProbe supports 4 export formats via \`-of\`:

1. **Terminal High-Contrast Text (\`text\`):**
   Formatted tables with ANSI color highlighting (2xx bold green, 3xx cyan, 4xx yellow, 5xx red).

2. **Structured JSON (\`json\`):**
   Full audit trail including candidate count, duration, rate, and array of matches with millisecond timestamps.

3. **CSV (\`csv\`):**
   Directly importable into spreadsheet analysis applications or SIEM databases.

4. **GitHub-Flavored Markdown (\`md\`):**
   Ready-to-publish summary table for penetration testing documentation and client audit reports.
`
  },
  {
    id: 'security',
    title: 'Security & Authorized-Use Policy',
    content: `### Authorized Use Only

8WHIE WebProbe is engineered exclusively for authorized diagnostic, defensive, and security assessment workflows against infrastructure you own or have explicit written permission to test.

#### Security Safeguards
- **Zero Exploit Payloads:** Does not generate exploits, reverse shells, or malicious payloads.
- **Strict TLS Verification:** Certificate validation enabled by default; explicitly warns if \`-k\` is passed.
- **Bounded Concurrency & Rate Limiting:** Prevents denial-of-service against fragile backends.
- **Memory Capping:** Response stream reading capped at 5MB per probe to prevent memory exhaustion.
- **Zero Telemetry / Credential Leaks:** No phone-home telemetry; private headers are never broadcast externally.
`
  }
];

const FAQS = [
  {
    q: 'What is 8WHIE WebProbe?',
    a: '8WHIE WebProbe is an original open-source command-line tool written in Go for controlled HTTP discovery, endpoint enumeration, request pacing, and response filtering. Developed by 8WHIE for authorized security research.'
  },
  {
    q: 'Is 8WHIE WebProbe affiliated with ffuf or any existing tools?',
    a: 'No. WebProbe is a completely new, independently designed and authored project under the 8WHIE brand. It features original package architecture, worker engine, banner artwork, terminal UI, and documentation.'
  },
  {
    q: 'Does WebProbe work on Android / Termux?',
    a: 'Yes! WebProbe compiles and runs natively inside Android Termux using standard Go tools (pkg install golang git make). It is lightweight, fast, and needs zero external runtime libraries.'
  },
  {
    q: 'Does WebProbe require root permissions?',
    a: 'No. WebProbe executes as a standard user-space binary. It does not open raw sockets or require administrative privileges.'
  },
  {
    q: 'Can I use custom headers and authentication tokens?',
    a: 'Yes. You can supply multiple -H "Name: Value" flags to include Authorization tokens, API keys, or custom cookies. Always exercise caution and never store credentials in shared bash histories.'
  },
  {
    q: 'How do I prevent overwhelming the target server?',
    a: 'Use the -rate flag to enforce requests per second (e.g. -rate 10) and set a conservative worker count with -c (e.g. -c 5).'
  }
];

export default function App() {
  const [activeTab, setActiveTab] = useState<'simulator' | 'builder' | 'wordlists' | 'docs' | 'faq'>('simulator');
  const [copiedText, setCopiedText] = useState<string | null>(null);

  // Command Builder State
  const [targetUrl, setTargetUrl] = useState('https://staging.internal.corp/FUZZ');
  const [wordlistPath, setWordlistPath] = useState('wordlists/endpoints.txt');
  const [placeholder, setPlaceholder] = useState('FUZZ');
  const [method, setMethod] = useState('GET');
  const [concurrency, setConcurrency] = useState(20);
  const [rateLimit, setRateLimit] = useState(0);
  const [timeout, setTimeoutSec] = useState(10);
  const [matchCodes, setMatchCodes] = useState('200,301,302,403');
  const [filterCodes, setFilterCodes] = useState('404');
  const [filterSize, setFilterSize] = useState('');
  const [customHeader, setCustomHeader] = useState('');
  const [outputFile, setOutputFile] = useState('');
  const [outputFormat, setOutputFormat] = useState('text');
  const [followRedirects, setFollowRedirects] = useState(false);
  const [insecureTls, setInsecureTls] = useState(false);
  const [quietMode, setQuietMode] = useState(false);

  // Terminal Simulator State
  const [simRunning, setSimRunning] = useState(false);
  const [simProgress, setSimProgress] = useState(0);
  const [simResults, setSimResults] = useState<ProbeSimResult[]>([]);
  const [simViewFormat, setSimViewFormat] = useState<'text' | 'json' | 'csv'>('text');
  const [simStats, setSimStats] = useState({
    sent: 0,
    matched: 0,
    filtered: 0,
    elapsedMs: 0,
    rps: 0
  });

  // Wordlist Explorer State
  const [selectedWordlistKey, setSelectedWordlistKey] = useState<string>('api');
  const [wordlistSearch, setWordlistSearch] = useState('');

  // Documentation Browser State
  const [selectedDocId, setSelectedDocId] = useState('usage');

  // FAQ Accordion State
  const [expandedFaq, setExpandedFaq] = useState<number | null>(0);

  // Copy helper
  const handleCopy = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    setCopiedText(label);
    setTimeout(() => setCopiedText(null), 2000);
  };

  // Compile active CLI command string
  const generatedCommand = useMemo(() => {
    const parts = ['webprobe'];
    if (targetUrl) parts.push(`-u "${targetUrl}"`);
    if (wordlistPath) parts.push(`-w ${wordlistPath}`);
    if (placeholder && placeholder !== 'FUZZ') parts.push(`-p "${placeholder}"`);
    if (method !== 'GET') parts.push(`-X ${method}`);
    if (concurrency !== 20) parts.push(`-c ${concurrency}`);
    if (rateLimit > 0) parts.push(`-rate ${rateLimit}`);
    if (timeout !== 10) parts.push(`-t ${timeout}`);
    if (matchCodes) parts.push(`-mc ${matchCodes}`);
    if (filterCodes && filterCodes !== '404') parts.push(`-fc ${filterCodes}`);
    if (filterSize) parts.push(`-fs ${filterSize}`);
    if (customHeader) parts.push(`-H "${customHeader}"`);
    if (followRedirects) parts.push('-r');
    if (insecureTls) parts.push('-k');
    if (outputFile) parts.push(`-o ${outputFile}`);
    if (outputFormat !== 'text') parts.push(`-of ${outputFormat}`);
    if (quietMode) parts.push('-q');
    return parts.join(' ');
  }, [
    targetUrl, wordlistPath, placeholder, method, concurrency, 
    rateLimit, timeout, matchCodes, filterCodes, filterSize, 
    customHeader, followRedirects, insecureTls, outputFile, 
    outputFormat, quietMode
  ]);

  // Terminal Simulator Simulation Loop
  const simTimerRef = useRef<NodeJS.Timeout | null>(null);

  const startSimulation = () => {
    setSimRunning(true);
    setSimResults([]);
    setSimProgress(0);
    setSimStats({ sent: 0, matched: 0, filtered: 0, elapsedMs: 0, rps: 0 });

    const candidates = [
      { path: 'login', code: 200, size: 2840, words: 412, lines: 76, dur: 38 },
      { path: 'admin', code: 403, size: 310, words: 28, lines: 9, dur: 45 },
      { path: 'api/v1', code: 200, size: 1420, words: 160, lines: 34, dur: 32 },
      { path: 'notfound1', code: 404, size: 180, words: 14, lines: 4, dur: 22 },
      { path: 'metrics', code: 200, size: 4890, words: 820, lines: 142, dur: 54 },
      { path: 'notfound2', code: 404, size: 180, words: 14, lines: 4, dur: 20 },
      { path: 'redirect', code: 301, size: 165, words: 12, lines: 3, dur: 28, redir: 'https://staging.internal.corp/login' },
      { path: 'healthz', code: 200, size: 42, words: 2, lines: 1, dur: 18 },
      { path: 'debug', code: 403, size: 310, words: 28, lines: 9, dur: 36 },
      { path: 'notfound3', code: 404, size: 180, words: 14, lines: 4, dur: 21 },
      { path: 'swagger.json', code: 200, size: 12450, words: 1980, lines: 340, dur: 62 },
      { path: 'backup.sql', code: 404, size: 180, words: 14, lines: 4, dur: 24 },
      { path: 'config.json', code: 403, size: 310, words: 28, lines: 9, dur: 40 },
    ];

    let currentIndex = 0;
    const startTime = Date.now();

    simTimerRef.current = setInterval(() => {
      if (currentIndex >= candidates.length) {
        if (simTimerRef.current) clearInterval(simTimerRef.current);
        setSimRunning(false);
        return;
      }

      const item = candidates[currentIndex];
      currentIndex++;

      const is404 = item.code === 404;
      const isMatched = !is404;

      setSimStats(prev => {
        const sent = prev.sent + 1;
        const matched = prev.matched + (isMatched ? 1 : 0);
        const filtered = prev.filtered + (is404 ? 1 : 0);
        const elapsed = Math.max(10, Date.now() - startTime);
        const rps = (sent / (elapsed / 1000));
        return { sent, matched, filtered, elapsedMs: elapsed, rps: Math.round(rps * 10) / 10 };
      });

      setSimProgress(Math.round((currentIndex / candidates.length) * 100));

      if (isMatched) {
        setSimResults(prev => [
          ...prev,
          {
            id: currentIndex,
            payload: item.path,
            url: `https://staging.internal.corp/${item.path}`,
            statusCode: item.code,
            contentLength: item.size,
            wordCount: item.words,
            lineCount: item.lines,
            durationMs: item.dur,
            redirectUrl: item.redir
          }
        ]);
      }
    }, 280);
  };

  const stopSimulation = () => {
    if (simTimerRef.current) clearInterval(simTimerRef.current);
    setSimRunning(false);
  };

  const resetSimulation = () => {
    stopSimulation();
    setSimResults([]);
    setSimProgress(0);
    setSimStats({ sent: 0, matched: 0, filtered: 0, elapsedMs: 0, rps: 0 });
  };

  useEffect(() => {
    return () => {
      if (simTimerRef.current) clearInterval(simTimerRef.current);
    };
  }, []);

  // Filtered wordlist items
  const currentWordlistItems = useMemo(() => {
    const list = SAMPLE_WORDLISTS[selectedWordlistKey] || [];
    if (!wordlistSearch) return list;
    return list.filter(item => item.toLowerCase().includes(wordlistSearch.toLowerCase()));
  }, [selectedWordlistKey, wordlistSearch]);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
      
      {/* Top Bar Contract: 3 zones */}
      <header className="border-b border-slate-800 bg-slate-950/80 backdrop-blur sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          
          {/* Zone 1: Single text element wordmark */}
          <div className="flex items-center gap-3">
            <span className="font-mono font-bold text-lg tracking-tight text-white flex items-center gap-2">
              <span className="text-cyan-400">8WHIE</span> WebProbe
            </span>
          </div>

          {/* Zone 2: 4-6 clean text navigation links */}
          <nav className="hidden md:flex items-center gap-8 text-sm font-medium text-slate-400">
            <button 
              onClick={() => setActiveTab('simulator')}
              className={`hover:text-cyan-400 transition-colors ${activeTab === 'simulator' ? 'text-cyan-400 font-semibold' : ''}`}
            >
              Simulator
            </button>
            <button 
              onClick={() => setActiveTab('builder')}
              className={`hover:text-cyan-400 transition-colors ${activeTab === 'builder' ? 'text-cyan-400 font-semibold' : ''}`}
            >
              CLI Generator
            </button>
            <button 
              onClick={() => setActiveTab('wordlists')}
              className={`hover:text-cyan-400 transition-colors ${activeTab === 'wordlists' ? 'text-cyan-400 font-semibold' : ''}`}
            >
              Wordlists
            </button>
            <button 
              onClick={() => setActiveTab('docs')}
              className={`hover:text-cyan-400 transition-colors ${activeTab === 'docs' ? 'text-cyan-400 font-semibold' : ''}`}
            >
              Documentation
            </button>
            <button 
              onClick={() => setActiveTab('faq')}
              className={`hover:text-cyan-400 transition-colors ${activeTab === 'faq' ? 'text-cyan-400 font-semibold' : ''}`}
            >
              FAQ
            </button>
          </nav>

          {/* Zone 3: 1-2 primary actions */}
          <div className="flex items-center gap-3">
            <button 
              onClick={() => handleCopy('git clone https://github.com/8whie/8whie-webprobe.git', 'clone')}
              className="px-3.5 py-1.5 text-xs font-mono bg-slate-900 hover:bg-slate-800 text-slate-300 border border-slate-700 rounded-md transition-colors flex items-center gap-1.5"
            >
              {copiedText === 'clone' ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
              <span>git clone</span>
            </button>
            <a 
              href="https://t.me/whiee" 
              target="_blank" 
              rel="noreferrer"
              className="px-3.5 py-1.5 text-xs font-medium text-slate-950 bg-cyan-400 hover:bg-cyan-300 rounded-md transition-colors flex items-center gap-1.5 shadow-sm shadow-cyan-950"
            >
              <ExternalLink className="w-3.5 h-3.5" />
              <span>Channel</span>
            </a>
          </div>

        </div>
      </header>

      {/* Main Content Viewport */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-6 py-8">
        
        {/* Hero Section */}
        <section className="mb-10 pb-8 border-b border-slate-800">
          <div className="max-w-4xl">
            
            {/* Zero-Pill Unboxed Metadata Header */}
            <div className="flex items-center gap-2 text-xs font-mono text-cyan-400 mb-3 tracking-wide">
              <span>Android</span>
              <span aria-hidden="true" className="text-slate-600">·</span>
              <span>Termux</span>
              <span aria-hidden="true" className="text-slate-600">·</span>
              <span>Linux</span>
              <span aria-hidden="true" className="text-slate-600">·</span>
              <span>macOS</span>
              <span aria-hidden="true" className="text-slate-600">·</span>
              <span>Windows</span>
              <span aria-hidden="true" className="text-slate-600">·</span>
              <span>MIT License</span>
            </div>

            <h1 className="text-3xl sm:text-4xl font-bold tracking-tight text-white mb-3" style={{ textWrap: 'balance' }}>
              Original Command-Line Utility for HTTP Endpoint Discovery & Request Analysis
            </h1>
            
            <p className="text-slate-400 text-base leading-relaxed mb-6 max-w-3xl">
              8WHIE WebProbe is an independently engineered Go CLI utility providing controlled URL enumeration, 
              concurrency throttling, response size and status code filtering, and machine-readable data export for authorized 
              security testing and systems verification.
            </p>

            {/* Quick Command Banner */}
            <div className="bg-slate-900 border border-slate-800 rounded-lg p-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3 font-mono text-xs text-slate-300">
              <div className="flex items-center gap-2.5 overflow-x-auto">
                <span className="text-cyan-400 font-bold">$</span>
                <span className="text-slate-200">webprobe -u https://target.local/FUZZ -w paths.txt -mc 200,301,302</span>
              </div>
              <button 
                onClick={() => handleCopy('webprobe -u https://target.local/FUZZ -w paths.txt -mc 200,301,302', 'quickcmd')}
                className="self-end sm:self-auto px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded transition-colors flex items-center gap-1.5 shrink-0"
              >
                {copiedText === 'quickcmd' ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
                <span>Copy</span>
              </button>
            </div>

            {/* Mandatory Authorized Use Advisory Notice */}
            <div className="mt-4 p-3 bg-amber-950/30 border border-amber-800/60 rounded-md text-amber-300/90 text-xs flex items-start gap-2.5">
              <Shield className="w-4 h-4 text-amber-400 shrink-0 mt-0.5" />
              <div>
                <span className="font-semibold text-amber-200">Authorized Use Only: </span> 
                This tool is engineered strictly for diagnostic assessments and authorized security research against targets you own or have explicit written permission to test.
              </div>
            </div>

          </div>
        </section>

        {/* Interactive Workspace Navigation Tabs */}
        <div className="flex items-center gap-2 border-b border-slate-800 mb-8 overflow-x-auto pb-px">
          <button 
            onClick={() => setActiveTab('simulator')}
            className={`px-4 py-2 text-sm font-medium transition-colors flex items-center gap-2 border-b-2 whitespace-nowrap ${
              activeTab === 'simulator' 
                ? 'border-cyan-400 text-cyan-400' 
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <Terminal className="w-4 h-4" />
            <span>Interactive Simulator</span>
          </button>
          
          <button 
            onClick={() => setActiveTab('builder')}
            className={`px-4 py-2 text-sm font-medium transition-colors flex items-center gap-2 border-b-2 whitespace-nowrap ${
              activeTab === 'builder' 
                ? 'border-cyan-400 text-cyan-400' 
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <Sliders className="w-4 h-4" />
            <span>CLI Generator</span>
          </button>

          <button 
            onClick={() => setActiveTab('wordlists')}
            className={`px-4 py-2 text-sm font-medium transition-colors flex items-center gap-2 border-b-2 whitespace-nowrap ${
              activeTab === 'wordlists' 
                ? 'border-cyan-400 text-cyan-400' 
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <Layers className="w-4 h-4" />
            <span>Wordlist Arsenal</span>
          </button>

          <button 
            onClick={() => setActiveTab('docs')}
            className={`px-4 py-2 text-sm font-medium transition-colors flex items-center gap-2 border-b-2 whitespace-nowrap ${
              activeTab === 'docs' 
                ? 'border-cyan-400 text-cyan-400' 
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <BookOpen className="w-4 h-4" />
            <span>Documentation</span>
          </button>

          <button 
            onClick={() => setActiveTab('faq')}
            className={`px-4 py-2 text-sm font-medium transition-colors flex items-center gap-2 border-b-2 whitespace-nowrap ${
              activeTab === 'faq' 
                ? 'border-cyan-400 text-cyan-400' 
                : 'border-transparent text-slate-400 hover:text-slate-200'
            }`}
          >
            <HelpCircle className="w-4 h-4" />
            <span>FAQ</span>
          </button>
        </div>

        {/* TAB 1: INTERACTIVE SIMULATOR */}
        {activeTab === 'simulator' && (
          <div className="space-y-6">
            
            {/* Simulator Controls & Telemetry Header */}
            <div className="bg-slate-900 border border-slate-800 rounded-lg p-4">
              <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
                
                <div className="flex items-center gap-3">
                  {!simRunning ? (
                    <button 
                      onClick={startSimulation}
                      className="px-4 py-2 bg-cyan-400 hover:bg-cyan-300 text-slate-950 font-medium text-xs rounded-md transition-colors flex items-center gap-2 shadow-sm"
                    >
                      <Play className="w-3.5 h-3.5 fill-current" />
                      <span>Start Probe Scan</span>
                    </button>
                  ) : (
                    <button 
                      onClick={stopSimulation}
                      className="px-4 py-2 bg-amber-500 hover:bg-amber-400 text-slate-950 font-medium text-xs rounded-md transition-colors flex items-center gap-2"
                    >
                      <span>Pause Scan</span>
                    </button>
                  )}

                  <button 
                    onClick={resetSimulation}
                    className="px-3 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs rounded-md transition-colors flex items-center gap-1.5"
                  >
                    <RotateCcw className="w-3.5 h-3.5" />
                    <span>Reset</span>
                  </button>

                  <div className="text-xs font-mono text-slate-400 pl-2">
                    Target: <span className="text-slate-200">https://staging.internal.corp/FUZZ</span>
                  </div>
                </div>

                {/* Live Telemetry Counter */}
                <div className="flex items-center gap-6 font-mono text-xs">
                  <div>
                    <span className="text-slate-500">Sent: </span>
                    <span className="text-white font-semibold tabular-nums">{simStats.sent}</span>
                  </div>
                  <div>
                    <span className="text-slate-500">Matched: </span>
                    <span className="text-emerald-400 font-semibold tabular-nums">{simStats.matched}</span>
                  </div>
                  <div>
                    <span className="text-slate-500">Filtered: </span>
                    <span className="text-slate-400 font-semibold tabular-nums">{simStats.filtered}</span>
                  </div>
                  <div>
                    <span className="text-slate-500">Speed: </span>
                    <span className="text-cyan-400 font-semibold tabular-nums">{simStats.rps} req/s</span>
                  </div>
                </div>

                {/* View Format Selector */}
                <div className="flex items-center gap-1 bg-slate-950 p-1 border border-slate-800 rounded-md">
                  <button 
                    onClick={() => setSimViewFormat('text')}
                    className={`px-2.5 py-1 text-xs rounded transition-colors ${simViewFormat === 'text' ? 'bg-slate-800 text-cyan-300' : 'text-slate-400'}`}
                  >
                    Terminal
                  </button>
                  <button 
                    onClick={() => setSimViewFormat('json')}
                    className={`px-2.5 py-1 text-xs rounded transition-colors ${simViewFormat === 'json' ? 'bg-slate-800 text-cyan-300' : 'text-slate-400'}`}
                  >
                    JSON
                  </button>
                  <button 
                    onClick={() => setSimViewFormat('csv')}
                    className={`px-2.5 py-1 text-xs rounded transition-colors ${simViewFormat === 'csv' ? 'bg-slate-800 text-cyan-300' : 'text-slate-400'}`}
                  >
                    CSV
                  </button>
                </div>

              </div>

              {/* Progress Bar */}
              <div className="mt-3 w-full bg-slate-950 h-1 rounded-full overflow-hidden">
                <div 
                  className="bg-cyan-400 h-full transition-all duration-200" 
                  style={{ width: `${simProgress}%` }}
                />
              </div>
            </div>

            {/* Virtual Terminal Screen */}
            <div className="bg-slate-950 border border-slate-800 rounded-lg overflow-hidden shadow-2xl font-mono text-xs">
              
              {/* Terminal Window Header Bar */}
              <div className="bg-slate-900 border-b border-slate-800 px-4 py-2.5 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-2.5 h-2.5 rounded-full bg-red-500/80" />
                  <div className="w-2.5 h-2.5 rounded-full bg-amber-500/80" />
                  <div className="w-2.5 h-2.5 rounded-full bg-emerald-500/80" />
                  <span className="text-slate-400 text-xs pl-2 font-mono">webprobe :: terminal output</span>
                </div>
                <div className="text-slate-500 text-xs">
                  {simResults.length} findings matched
                </div>
              </div>

              {/* Terminal Body */}
              <div className="p-4 overflow-x-auto min-h-[360px] max-h-[500px] overflow-y-auto space-y-1">
                
                {/* Banner */}
                <pre className="text-cyan-400 leading-tight select-none">
{`  ___ _ _ _ _  _ ___   _ _ _     _   ___          _         
 ( _ ) | | | || |_ _| | | | |___| |_| _ \\_ _ ___ | |__  ___ 
 / _ \\_  _ | __ || |  | | | / -_) '_ \\  _/ '_/ _ \\| '_ \\/ -_)
 \\___/ |_| |_||_|___| |_____/\\___|_,__/_| |_| \\___/|_.__/\\___|`}
                </pre>
                
                <div className="text-slate-400 text-xs py-2 border-b border-slate-800/80">
                  <div>:: 8WHIE WebProbe           : v1.0.0</div>
                  <div>:: Target                   : https://staging.internal.corp/FUZZ</div>
                  <div>:: Wordlist                 : wordlists/endpoints.txt</div>
                  <div>:: Filter                   : -fc 404</div>
                  <div>:: Concurrency / Rate       : 20 workers / unrestricted</div>
                </div>

                {/* View Format Rendering */}
                {simViewFormat === 'text' && (
                  <div className="pt-2">
                    <div className="text-slate-500 font-bold pb-1 border-b border-slate-800 flex gap-6">
                      <span className="w-16">STATUS</span>
                      <span className="w-20">SIZE</span>
                      <span className="w-16">WORDS</span>
                      <span className="w-16">LINES</span>
                      <span className="w-20">DURATION</span>
                      <span className="flex-1">PAYLOAD / REDIRECT</span>
                    </div>

                    {simResults.length === 0 ? (
                      <div className="text-slate-600 py-10 text-center">
                        Click "Start Probe Scan" to initiate controlled endpoint discovery simulation.
                      </div>
                    ) : (
                      simResults.map(r => {
                        let statusColor = 'text-emerald-400';
                        if (r.statusCode >= 300 && r.statusCode < 400) statusColor = 'text-cyan-400';
                        if (r.statusCode === 401 || r.statusCode === 403) statusColor = 'text-amber-400';
                        if (r.statusCode >= 500) statusColor = 'text-red-400';

                        return (
                          <div key={r.id} className="flex gap-6 py-1 hover:bg-slate-900/60 rounded px-1 transition-colors">
                            <span className={`w-16 font-bold ${statusColor}`}>{r.statusCode}</span>
                            <span className="w-20 text-slate-300 tabular-nums">{r.contentLength} B</span>
                            <span className="w-16 text-slate-400 tabular-nums">{r.wordCount}</span>
                            <span className="w-16 text-slate-400 tabular-nums">{r.lineCount}</span>
                            <span className="w-20 text-slate-400 tabular-nums">{r.durationMs}ms</span>
                            <span className="flex-1 text-slate-200">
                              <span className="text-white">{r.payload}</span>
                              {r.redirectUrl && <span className="text-slate-500 ml-2">→ {r.redirectUrl}</span>}
                            </span>
                          </div>
                        );
                      })
                    )}
                  </div>
                )}

                {simViewFormat === 'json' && (
                  <pre className="text-emerald-300 p-2 overflow-x-auto text-xs">
                    {JSON.stringify({
                      tool: '8WHIE WebProbe',
                      version: '1.0.0',
                      target: 'https://staging.internal.corp/FUZZ',
                      statistics: simStats,
                      results: simResults
                    }, null, 2)}
                  </pre>
                )}

                {simViewFormat === 'csv' && (
                  <pre className="text-slate-300 p-2 overflow-x-auto text-xs">
                    {`ID,Payload,URL,Status,Size,Words,Lines,DurationMs,Redirect\n` +
                      simResults.map(r => 
                        `${r.id},${r.payload},${r.url},${r.statusCode},${r.contentLength},${r.wordCount},${r.lineCount},${r.durationMs},${r.redirectUrl || ''}`
                      ).join('\n')
                    }
                  </pre>
                )}

              </div>
            </div>

          </div>
        )}

        {/* TAB 2: CLI GENERATOR */}
        {activeTab === 'builder' && (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
            
            {/* Configuration Inputs */}
            <div className="lg:col-span-7 space-y-5">
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-5 space-y-4">
                <h3 className="text-base font-semibold text-white">Target & Wordlist Configuration</h3>
                
                <div>
                  <label className="block text-xs font-mono text-slate-400 mb-1.5">
                    Target URL with Placeholder (-u)
                  </label>
                  <input 
                    type="text" 
                    value={targetUrl}
                    onChange={e => setTargetUrl(e.target.value)}
                    className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    placeholder="https://example.com/FUZZ"
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Wordlist Path (-w)</label>
                    <input 
                      type="text" 
                      value={wordlistPath}
                      onChange={e => setWordlistPath(e.target.value)}
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Placeholder Token (-p)</label>
                    <input 
                      type="text" 
                      value={placeholder}
                      onChange={e => setPlaceholder(e.target.value)}
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">HTTP Method (-X)</label>
                    <select 
                      value={method}
                      onChange={e => setMethod(e.target.value)}
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    >
                      <option value="GET">GET</option>
                      <option value="POST">POST</option>
                      <option value="HEAD">HEAD</option>
                      <option value="PUT">PUT</option>
                      <option value="DELETE">DELETE</option>
                      <option value="OPTIONS">OPTIONS</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Concurrency (-c): {concurrency}</label>
                    <input 
                      type="range" 
                      min="1" 
                      max="100" 
                      value={concurrency}
                      onChange={e => setConcurrency(Number(e.target.value))}
                      className="w-full accent-cyan-400 cursor-pointer"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Rate Limit (-rate): {rateLimit > 0 ? `${rateLimit} rps` : 'Unlimited'}</label>
                    <input 
                      type="range" 
                      min="0" 
                      max="200" 
                      step="5"
                      value={rateLimit}
                      onChange={e => setRateLimit(Number(e.target.value))}
                      className="w-full accent-cyan-400 cursor-pointer"
                    />
                  </div>
                </div>

              </div>

              {/* Filtering Controls */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-5 space-y-4">
                <h3 className="text-base font-semibold text-white">Filtering & Response Rules</h3>
                
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Match Status Codes (-mc)</label>
                    <input 
                      type="text" 
                      value={matchCodes}
                      onChange={e => setMatchCodes(e.target.value)}
                      placeholder="200,301,302"
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Filter Status Codes (-fc)</label>
                    <input 
                      type="text" 
                      value={filterCodes}
                      onChange={e => setFilterCodes(e.target.value)}
                      placeholder="404,500"
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Filter Sizes In Bytes (-fs)</label>
                    <input 
                      type="text" 
                      value={filterSize}
                      onChange={e => setFilterSize(e.target.value)}
                      placeholder="e.g. 4528"
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-mono text-slate-400 mb-1.5">Custom Header (-H)</label>
                    <input 
                      type="text" 
                      value={customHeader}
                      onChange={e => setCustomHeader(e.target.value)}
                      placeholder="Authorization: Bearer <TOKEN>"
                      className="w-full bg-slate-950 border border-slate-800 rounded px-3 py-2 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                    />
                  </div>
                </div>

                {/* Checkbox Options */}
                <div className="pt-2 flex flex-wrap gap-4 text-xs font-mono">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input 
                      type="checkbox" 
                      checked={followRedirects}
                      onChange={e => setFollowRedirects(e.target.checked)}
                      className="accent-cyan-400"
                    />
                    <span>Follow Redirects (-r)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input 
                      type="checkbox" 
                      checked={insecureTls}
                      onChange={e => setInsecureTls(e.target.checked)}
                      className="accent-cyan-400"
                    />
                    <span>Allow Insecure TLS (-k)</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input 
                      type="checkbox" 
                      checked={quietMode}
                      onChange={e => setQuietMode(e.target.checked)}
                      className="accent-cyan-400"
                    />
                    <span>Quiet Mode (-q)</span>
                  </label>
                </div>

              </div>
            </div>

            {/* Live Command Preview Box */}
            <div className="lg:col-span-5 space-y-4">
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-5 sticky top-24">
                <div className="flex items-center justify-between mb-3">
                  <h3 className="text-sm font-semibold text-white">Generated CLI Command</h3>
                  <button 
                    onClick={() => handleCopy(generatedCommand, 'generated')}
                    className="px-3 py-1 bg-cyan-400 hover:bg-cyan-300 text-slate-950 font-medium text-xs rounded transition-colors flex items-center gap-1.5"
                  >
                    {copiedText === 'generated' ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
                    <span>Copy Command</span>
                  </button>
                </div>

                <div className="p-4 bg-slate-950 border border-slate-800 rounded-md font-mono text-xs text-cyan-300 break-all leading-relaxed">
                  {generatedCommand}
                </div>

                {/* Explanations of Active Flags */}
                <div className="mt-4 space-y-2 border-t border-slate-800 pt-4 text-xs font-mono">
                  <div className="text-slate-400 font-semibold text-[11px] uppercase tracking-wider">Active Flags Breakdown:</div>
                  <div className="text-slate-400">
                    <span className="text-white">-u:</span> Target containing placeholder token
                  </div>
                  <div className="text-slate-400">
                    <span className="text-white">-c {concurrency}:</span> {concurrency} concurrent discovery worker goroutines
                  </div>
                  {rateLimit > 0 && (
                    <div className="text-slate-400">
                      <span className="text-white">-rate {rateLimit}:</span> Maximum {rateLimit} requests dispatched per second
                    </div>
                  )}
                  {matchCodes && (
                    <div className="text-slate-400">
                      <span className="text-white">-mc {matchCodes}:</span> Matches only HTTP status {matchCodes}
                    </div>
                  )}
                  {insecureTls && (
                    <div className="text-amber-400">
                      <span>-k:</span> Bypasses TLS cert validation (staging labs only)
                    </div>
                  )}
                </div>

              </div>
            </div>

          </div>
        )}

        {/* TAB 3: WORDLIST ARSENAL */}
        {activeTab === 'wordlists' && (
          <div className="space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-slate-900 border border-slate-800 rounded-lg p-4">
              
              <div className="flex items-center gap-2">
                {Object.keys(SAMPLE_WORDLISTS).map(key => (
                  <button 
                    key={key}
                    onClick={() => setSelectedWordlistKey(key)}
                    className={`px-3 py-1.5 text-xs font-mono uppercase rounded transition-colors ${
                      selectedWordlistKey === key 
                        ? 'bg-cyan-400 text-slate-950 font-bold' 
                        : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                    }`}
                  >
                    {key}
                  </button>
                ))}
              </div>

              <div className="flex items-center gap-3">
                <input 
                  type="text" 
                  value={wordlistSearch}
                  onChange={e => setWordlistSearch(e.target.value)}
                  placeholder="Filter paths..."
                  className="bg-slate-950 border border-slate-800 rounded px-3 py-1.5 text-xs font-mono text-white focus:outline-none focus:border-cyan-400"
                />
                <button 
                  onClick={() => {
                    const text = currentWordlistItems.join('\n');
                    const blob = new Blob([text], { type: 'text/plain' });
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = `${selectedWordlistKey}-wordlist.txt`;
                    a.click();
                    URL.revokeObjectURL(url);
                  }}
                  className="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-mono rounded transition-colors flex items-center gap-1.5"
                >
                  <Download className="w-3.5 h-3.5" />
                  <span>Download .txt</span>
                </button>
              </div>

            </div>

            {/* Wordlist Table Grid */}
            <div className="bg-slate-900 border border-slate-800 rounded-lg overflow-hidden">
              <div className="p-3 bg-slate-950 border-b border-slate-800 text-xs font-mono text-slate-400 flex justify-between">
                <span>Wordlist Entries: {currentWordlistItems.length} items</span>
                <span>Ready for WebProbe -w flag</span>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2 p-4 font-mono text-xs">
                {currentWordlistItems.map((item, idx) => (
                  <div key={idx} className="bg-slate-950 border border-slate-800/80 rounded px-3 py-2 text-slate-300 hover:border-cyan-500/50 transition-colors flex items-center justify-between">
                    <span className="text-cyan-300">/{item}</span>
                    <button 
                      onClick={() => handleCopy(item, `word-${idx}`)}
                      className="text-slate-500 hover:text-slate-300 p-0.5"
                    >
                      {copiedText === `word-${idx}` ? <Check className="w-3 h-3 text-emerald-400" /> : <Copy className="w-3 h-3" />}
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* TAB 4: DOCUMENTATION */}
        {activeTab === 'docs' && (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
            
            {/* Sidebar nav */}
            <div className="lg:col-span-3 space-y-1">
              <div className="text-xs font-mono text-slate-500 uppercase tracking-wider px-3 mb-2">Guides</div>
              {DOC_SECTIONS.map(sec => (
                <button 
                  key={sec.id}
                  onClick={() => setSelectedDocId(sec.id)}
                  className={`w-full text-left px-3 py-2 rounded text-xs font-medium transition-colors flex items-center justify-between ${
                    selectedDocId === sec.id 
                      ? 'bg-slate-800 text-cyan-400 font-semibold' 
                      : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900'
                  }`}
                >
                  <span>{sec.title}</span>
                  {selectedDocId === sec.id && <ChevronRight className="w-3.5 h-3.5" />}
                </button>
              ))}
            </div>

            {/* Document Content Viewport */}
            <div className="lg:col-span-9 bg-slate-900 border border-slate-800 rounded-lg p-6">
              {(() => {
                const current = DOC_SECTIONS.find(s => s.id === selectedDocId);
                if (!current) return null;
                return (
                  <div className="prose prose-invert max-w-none text-xs sm:text-sm font-sans space-y-4">
                    <h2 className="text-xl font-bold text-white mb-4 pb-2 border-b border-slate-800">
                      {current.title}
                    </h2>
                    <div className="whitespace-pre-wrap font-mono text-xs text-slate-300 bg-slate-950 p-4 rounded-md border border-slate-800 leading-relaxed overflow-x-auto">
                      {current.content}
                    </div>
                  </div>
                );
              })()}
            </div>

          </div>
        )}

        {/* TAB 5: FAQ */}
        {activeTab === 'faq' && (
          <div className="max-w-3xl mx-auto space-y-3">
            {FAQS.map((faq, idx) => {
              const isOpen = expandedFaq === idx;
              return (
                <div key={idx} className="bg-slate-900 border border-slate-800 rounded-lg overflow-hidden">
                  <button 
                    onClick={() => setExpandedFaq(isOpen ? null : idx)}
                    className="w-full text-left p-4 text-sm font-semibold text-white flex items-center justify-between hover:bg-slate-800/40 transition-colors"
                  >
                    <span>{faq.q}</span>
                    <ChevronDown className={`w-4 h-4 text-slate-400 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
                  </button>
                  {isOpen && (
                    <div className="px-4 pb-4 pt-1 text-xs text-slate-300 leading-relaxed border-t border-slate-800/50 bg-slate-950/40">
                      {faq.a}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}

      </main>

      {/* Footer */}
      <footer className="border-t border-slate-800 bg-slate-950 py-8 px-6 mt-16">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-slate-500 font-mono">
          <div>
            <span className="text-slate-400 font-bold">8WHIE WebProbe</span> · Developed by <span className="text-cyan-400">8WHIE</span> (Copyright © 2026)
          </div>
          <div className="flex items-center gap-6">
            <a href="https://t.me/whiee" target="_blank" rel="noreferrer" className="hover:text-cyan-400 transition-colors">Telegram Channel</a>
            <a href="https://instagram.com/aaynkt" target="_blank" rel="noreferrer" className="hover:text-cyan-400 transition-colors">Instagram</a>
            <a href="https://github.com/8whie" target="_blank" rel="noreferrer" className="hover:text-cyan-400 transition-colors">GitHub</a>
          </div>
        </div>
      </footer>

    </div>
  );
}
