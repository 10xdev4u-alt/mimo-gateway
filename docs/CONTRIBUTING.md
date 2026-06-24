# Contributing

## Development Setup

1. Fork and clone the repository
2. Install dependencies
3. Start development server

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

## Code Style

- Go: Follow `gofmt` conventions
- TypeScript: Use Prettier
- Commit messages: Conventional commits format

## Pull Requests

1. Create a feature branch
2. Make your changes
3. Add tests if applicable
4. Submit a pull request

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
