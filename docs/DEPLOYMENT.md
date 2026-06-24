# Deployment Guide

## Docker Deployment

### Quick Start
```bash
# Clone the repository
git clone https://github.com/10xdev4u-alt/mimo-gateway.git
cd mimo-gateway

# Start services
docker compose -f docker-compose.dev.yml up -d

# Check status
curl http://localhost:4200/health
```

### Production Deployment
```bash
# Build production image
docker build -t mimo-gateway:latest .

# Run with environment variables
docker run -d \
  -p 4200:4200 \
  -e JWT_SECRET=your-secret-here \
  -e DATABASE_URL=postgres://user:pass@host:5432/db \
  mimo-gateway:latest
```

## Manual Deployment

### Prerequisites
- Go 1.24+
- Node.js 20+
- PostgreSQL 16+
- Redis 7+

### Steps
```bash
# Install dependencies
cd apps/api && go mod download

# Build
go build -o mimo-gateway ./cmd/server

# Run
./mimo-gateway
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_PORT | Server port | 4200 |
| JWT_SECRET | JWT signing secret | (required) |
| DATABASE_URL | PostgreSQL connection | (required) |
| REDIS_URL | Redis connection | redis://localhost:6380 |
| MIMO_BIN_PATH | Path to .mimocode binary | auto-detect |

## Health Checks

```bash
# Basic health
curl http://localhost:4200/health

# Detailed status
curl http://localhost:4200/v1/models
```
