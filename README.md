# MiMo Gateway

> OpenAI-compatible API gateway for MiMo Auto Free API

[![CI](https://github.com/10xdev4u-alt/mimo-gateway/actions/workflows/ci.yml/badge.svg)](https://github.com/10xdev4u-alt/mimo-gateway/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## What is this?

MiMo Gateway is a proxy that routes requests to the MiMo Auto Free API through the official `.mimocode` binary. It provides a standard OpenAI-compatible interface, so you can use any OpenAI SDK or tool with MiMo's free models.

**Why a proxy?** The MiMo API gateway blocks direct HTTP calls via TLS fingerprinting (JA3). Only the official binary passes verification. This proxy wraps that binary and exposes a clean API.

## Features

- **OpenAI-compatible** — Drop-in replacement for `/v1/chat/completions`
- **Streaming support** — SSE streaming for real-time responses
- **Bifrost dashboard** — Dark-themed admin panel with real-time stats
- **Rate limiting** — Configurable per-IP rate limits
- **API key auth** — Secure access with API keys
- **Docker ready** — One-command deployment
- **CI/CD** — GitHub Actions for tests and security scanning

## Quick Start

```bash
# Clone
git clone https://github.com/10xdev4u-alt/mimo-gateway.git
cd mimo-gateway

# Docker (recommended)
docker compose -f docker-compose.dev.yml up -d

# Or manual
cd apps/api
go mod download
go run ./cmd/server
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/v1/chat/completions` | Chat completion (OpenAI format) |
| `GET` | `/v1/models` | List available models |
| `GET` | `/health` | Health check |

## Usage

```bash
# Test with curl
curl http://localhost:4200/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "mimo-auto",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

# Or with any OpenAI SDK
export OPENAI_BASE_URL=http://localhost:4200/v1
export OPENAI_API_KEY=your-api-key
```

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Client    │────▶│  MiMo Proxy  │────▶│  .mimocode  │
│  (OpenAI)   │     │   (Go API)   │     │   binary    │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │   React     │
                    │  Dashboard  │
                    └─────────────┘
```

## Dashboard

Access the admin dashboard at `http://localhost:4202`

- Real-time request stats
- API key management
- Request logs
- Model information
- System health

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | API server port | `4200` |
| `JWT_SECRET` | JWT signing secret | (required) |
| `DATABASE_URL` | PostgreSQL connection | (required) |
| `MIMO_BIN_PATH` | Path to .mimocode binary | auto-detect |
| `RATE_LIMIT` | Requests per minute | `100` |

## Development

```bash
# Backend
cd apps/api
go mod download
go run ./cmd/server

# Frontend
cd apps/admin
pnpm install
pnpm dev
```

## Testing

```bash
# Go tests
cd apps/api
go test -v ./...

# Lint
go vet ./...
```

## Deployment

```bash
# Docker
docker build -t mimo-gateway:latest .
docker run -p 4200:4200 mimo-gateway:latest

# Or with docker-compose
docker compose -f docker-compose.dev.yml up -d
```

## Credits

- Built with [Grit Framework](https://gritframework.dev)
- Uses [MiMo Auto Free API](https://mimo.xiaomi.com/coder) by Xiaomi
- Dashboard theme: Bifrost (custom dark theme)

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Powered by MiMo Auto Free API** 🔥
