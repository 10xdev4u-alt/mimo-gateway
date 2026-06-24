# Troubleshooting

## Common Issues

### Gateway won't start
- Check that port 4200 is not in use
- Verify environment variables are set
- Check database connection

### 403 Forbidden from MiMo API
- This is expected - direct API calls are blocked
- Use the binary backend instead
- Ensure .mimocode binary is installed

### High latency
- Binary initialization takes ~5-7 seconds
- Subsequent requests should be faster
- Check network connectivity

### Build errors
- Ensure Go 1.24+ is installed
- Run `go mod tidy` to fix dependencies
- Check for import conflicts

## Debug Mode

Set environment variable for debug logging:
```bash
APP_ENV=development DEBUG=true ./mimo-gateway
```

## Health Check

```bash
curl http://localhost:4200/health
```
