# Cursor-Tap

[中文](./README.md) | English

A tool for intercepting and analyzing Cursor IDE's gRPC traffic. Decrypts TLS, deserializes protobuf, and displays AI conversations in real-time.

## Why

Cursor talks to its backend entirely via gRPC (Connect Protocol). The body is binary protobuf. Burp and Fiddler just show unreadable bytes. Cursor doesn't publish proto definitions either.

This tool decrypts traffic into readable JSON and shows each streaming frame in real-time.

## How It Works

1. **MITM Proxy** - Sits between Cursor and api2.cursor.sh, decrypts TLS with self-signed CA
2. **Proto Extraction** - Extracts proto definitions from Cursor's JS bundle (`protobuf-es` compiled output)
3. **Real-time Parsing** - Parses Connect Protocol envelope framing, deserializes each protobuf frame
4. **WebUI** - Pushes frames via WebSocket, four-panel layout for service tree / calls / frames / details

## Quick Start

### Using Nx (Recommended)

```bash
pnpm install

# Start proxy + API + WebUI together
pnpm dev

# Or start separately
pnpm dev:api
pnpm dev:tap
pnpm dev:web

# Other commands
pnpm build:tap
pnpm build:api
pnpm build:web
pnpm exec nx run proto:generate
pnpm exec nx test tap
pnpm exec nx test api
```

### Manual Setup

#### 1. Start API and Proxy

```bash
# Terminal 1: API server (REST + WebSocket)
go run ./apps/api start

# Terminal 2: MITM proxy (writes SQLite, notifies API)
go run ./apps/tap start --http-parse
```

Proxy listens on `localhost:8080` (HTTP) and `localhost:1080` (SOCKS5). API listens on `localhost:9090`.

#### 2. Configure Cursor

```bash
# Windows
set HTTP_PROXY=http://localhost:8080
set HTTPS_PROXY=http://localhost:8080
set NODE_EXTRA_CA_CERTS=C:\path\to\ca.crt

# macOS/Linux
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
export NODE_EXTRA_CA_CERTS=/path/to/ca.crt
```

CA certificate is auto-generated at `~/.cursor-tap/ca/ca.crt` on first run.

#### 3. Start WebUI

```bash
pnpm dev:web
# or: cd apps/web && pnpm dev
```

Open `http://localhost:3000`.

## Project Structure

```
├── apps/
│   ├── tap/            # Go MITM proxy (Nx: tap)
│   ├── api/            # Management API server (Nx: api)
│   └── web/            # Next.js Web Inspector (Nx: web)
├── internal/
│   ├── apiserver/      # Standalone API HTTP server
│   ├── recordstore/    # SQLite record persistence
│   ├── notify/         # Proxy → API update notifications
│   ├── ca/             # Self-signed CA, dynamic certs
│   ├── config/         # ProxyConfig / APIConfig
│   ├── proxy/          # HTTP/SOCKS5 proxy
│   ├── mitm/           # TLS MITM, HTTP/2 bridge
│   ├── httpstream/     # HTTP/gRPC/SSE parsing
│   └── api/            # REST + WebSocket routes
├── packages/proto/     # Protobuf definitions (Nx: proto)
├── tests/              # Go unit tests
├── tools/              # Dev utilities (ext / restore / inline)
└── docs/               # Architecture docs & reverse notes
```

See [Architecture](./docs/architecture.md) for details.

## What You Can See

- `AiService/RunSSE` - AI conversation channel (thinking, text, tool calls)
- `BidiService/BidiAppend` - User messages and tool results
- `AiService/StreamCpp` - Code completion
- `CppService/RecordCppFate` - Completion accept/reject feedback
- `AiService/Batch` - User behavior telemetry
- And dozens more...

## Disclaimer

For educational and research purposes only.

## Documentation

- [Architecture](./docs/architecture.md) — monorepo structure, modules, data flow, and API (Chinese)
- [Reverse Engineering Notes 1](./docs/cursor-reverse-notes-1.md) — detailed MITM and proto extraction process (Chinese)
