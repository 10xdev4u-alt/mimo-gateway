# MiMo Gateway Architecture

## Overview

MiMo Gateway is an OpenAI-compatible API proxy for the MiMo Auto Free API. It routes requests through the official `.mimocode` binary to bypass TLS fingerprinting (JA3) that blocks direct API calls.

## Components

### Backend (Go + Grit)
- **API Server**: Gin-based HTTP server on port 4200
- **MiMo Proxy Handler**: Wraps the `.mimocode` binary
- **Auth System**: JWT-based authentication
- **Observability**: Pulse integration for metrics

### Frontend (React + Vite)
- **Web App**: Public-facing interface (port 4201)
- **Admin Dashboard**: Bifrost-themed admin panel (port 4202)

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /v1/chat/completions | OpenAI-compatible chat |
| GET | /v1/models | List available models |
| GET | /health | Gateway health check |

## Binary Backend

The gateway uses the `.mimocode` binary as a backend because:
1. MiFE gateway blocks direct HTTP calls (JA3 fingerprinting)
2. Only the compiled binary passes TLS verification
3. Binary handles auth (fingerprint → JWT) automatically

## Deployment

```bash
# Docker
docker compose up -d

# Manual
cd apps/api && go run cmd/server/main.go
```
