# Architecture

## Overview

MiMo Gateway is an OpenAI-compatible API proxy for the MiMo Auto Free API.

## Components

### Backend (Go + Gint)
- **API Server**: Gin-based HTTP server on port 4200
- **MiMo Proxy Handler**: Wraps the .mimocode binary
- **Auth System**: JWT-based authentication
- **Observability**: Pulse integration for metrics

### Frontend (React + Vite)
- **Web App**: Public-facing interface (port 4201)
- **Admin Dashboard**: Bifrost-themed admin panel (port 4202)

## Data Flow

1. Client sends request to /v1/chat/completions
2. Proxy handler extracts prompt from messages
3. Binary backend processes the request
4. Response is returned in OpenAI format

## Binary Backend

The gateway uses the .mimocode binary because:
1. MiFE gateway blocks direct HTTP calls
2. Only the compiled binary passes TLS verification
3. Binary handles auth automatically

## Security

- API key authentication
- Rate limiting per IP
- Input validation
- Security headers
