# Security

## API Key Authentication

All requests require an API key in the Authorization header:

```
Authorization: Bearer mg_your_api_key_here
```

## Rate Limiting

- 100 requests per minute per IP
- Returns 429 when exceeded

## Input Validation

- Request size limit: 10MB
- Input sanitization for XSS prevention
- JSON validation for all endpoints

## Security Headers

- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Referrer-Policy: strict-origin-when-cross-origin
- Content-Security-Policy: default-src 'self'

## Environment Variables

Sensitive values should be set via environment variables:

- JWT_SECRET: JWT signing secret (required)
- DATABASE_URL: Database connection (required)
- API keys and credentials
