# MiMo Gateway API Reference

## Base URL
```
http://localhost:4200
```

## Authentication
All requests require an API key in the Authorization header:
```
Authorization: Bearer <your-api-key>
```

## Endpoints

### POST /v1/chat/completions

Create a chat completion.

**Request:**
```json
{
  "model": "mimo-auto",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ],
  "stream": false
}
```

**Response:**
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "mimo-auto",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hi there! How can I help you?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 12,
    "total_tokens": 22
  }
}
```

### GET /v1/models

List available models.

**Response:**
```json
{
  "object": "list",
  "data": [
    {
      "id": "mimo-auto",
      "object": "model",
      "created": 1677652288,
      "owned_by": "xiaomi"
    }
  ]
}
```

### GET /health

Health check.

**Response:**
```json
{
  "status": "ok",
  "service": "mimo-gateway",
  "version": "1.0.0",
  "binary": "auto-detect",
  "uptime": "2h34m"
}
```

## Rate Limiting

- 100 requests per minute per IP
- Returns 429 Too Many Requests when exceeded

## Error Responses

All errors follow the OpenAI format:
```json
{
  "error": {
    "message": "Error description",
    "type": "invalid_request_error",
    "param": null,
    "code": null
  }
}
```
