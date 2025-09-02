# AGENTS.md

This file provides context and instructions for AI coding agents working on the Civic Auth Go SDK.

## Project Overview

The Civic Auth Go SDK is a comprehensive Go library for integrating with Civic Auth's OIDC/OAuth2 authentication service. It provides:
- Complete OIDC/OAuth2 client implementation with PKCE support
- JWT ID token validation with automatic JWK key fetching
- Token management with configurable storage
- Production-ready examples and comprehensive documentation

## Setup Commands

- Install dependencies: `make install-deps` or `go mod tidy`
- Run tests: `make test`
- Build examples: `make build`
- Run linting: `make lint`
- Format code: `make fmt`
- Clean build artifacts: `make clean`
- Run all checks: `make ci`

## Project Structure

```
civic-auth-go/
├── pkg/civicauth/         # Core SDK library
│   ├── config.go          # Configuration and types
│   ├── client.go          # Main OIDC client implementation
│   ├── tokens.go          # Token management utilities
│   └── *_test.go          # Unit tests
├── examples/              # Usage examples
│   ├── web_server.go      # Complete web application example
│   └── cli_example.go     # Command-line usage example
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

## Code Style and Conventions

- Follow standard Go conventions and `gofmt` formatting
- Use `gofmt` for consistent code formatting (run `make fmt`)
- Prefer explicit error handling over panics
- Use context.Context for cancellation and timeouts
- Follow Go naming conventions (exported vs unexported)
- Add godoc comments for all exported functions and types
- Use structured logging where applicable
- Keep functions focused and single-purpose

### Naming Patterns
- Use descriptive variable names
- Prefer `ctx` for context.Context parameters
- Use `err` for error return values
- Struct methods should have receiver names that are short but clear
- Constants should be in ALL_CAPS with underscores

## Testing Instructions

- All tests are in `*_test.go` files alongside the code they test
- Run `make test` to execute the full test suite
- Individual package tests: `go test ./pkg/civicauth/...`
- Test coverage: `make test-coverage` (generates coverage.html)
- Tests should be deterministic and not rely on external services
- Use table-driven tests for multiple test cases
- Mock external dependencies when necessary

### Test Requirements
- New features must include unit tests
- Test coverage should not decrease
- Tests should be fast and reliable
- Use descriptive test names that explain what is being tested

## Build and Development

- Go version: 1.19+ recommended
- No external build dependencies beyond standard Go toolchain
- Use `make build` to compile examples
- Binary outputs go to `bin/` directory (gitignored)
- Development workflow: `make dev` runs install-deps, lint, and test

## Security Considerations

- Never commit secrets or API keys
- Use environment variables for sensitive configuration
- Validate all inputs, especially from external sources
- Follow OWASP guidelines for authentication flows
- JWT tokens should always be validated before use
- Use PKCE for OAuth2 authorization code flow
- State parameters must be validated to prevent CSRF

### Configuration Security
- ClientSecret should never be logged or exposed
- Use secure storage for refresh tokens in production
- Implement proper session management in web applications
- Always use HTTPS in production environments

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT token parsing and validation
- Standard library only for core functionality
- Minimal external dependencies by design
- All dependencies are tracked in go.mod/go.sum

## Common Tasks

### Adding New Features
1. Create feature branch: `git checkout -b feature/feature-name`
2. Implement feature with tests
3. Run full test suite: `make test`
4. Check code style: `make lint`
5. Update documentation as needed
6. Create pull request with clear description

### Debugging
- Use Go's built-in debugging tools
- Add temporary logging with structured output
- Check error messages for context about failures
- Use `go vet` to catch common issues

### Working with Examples
- Examples are in `examples/` directory
- Each example is self-contained
- Run examples with environment variables set:
  ```bash
  export CIVIC_CLIENT_ID="your-client-id"
  export CIVIC_CLIENT_SECRET="your-client-secret"
  export CIVIC_ISSUER="https://auth.civic.com"
  go run examples/web_server.go
  ```

## API Design Principles

- Follow Go idioms and standard library patterns
- Return errors as the last return value
- Use interfaces for extensibility (e.g., TokenStorage)
- Provide sensible defaults in configuration
- Make zero values useful where possible
- Use context for cancellation and request scoping

## Documentation Standards

- All exported functions and types must have godoc comments
- Comments should explain what the function does, not how
- Include usage examples in godoc when helpful
- Update README.md for user-facing changes
- Maintain CHANGELOG.md for version history

## Error Handling

- Use Go's standard error handling patterns
- Wrap errors with context using fmt.Errorf
- Define custom error types for specific error conditions
- Never ignore errors without good reason
- Provide helpful error messages with context

## Performance Considerations

- HTTP clients should be reused, not created per request
- Cache JWK keys to avoid repeated fetches
- Use connection pooling for HTTP requests
- Avoid blocking operations in hot paths
- Profile code for performance bottlenecks if needed

## Integration Notes

- SDK is designed to work with any OIDC-compliant provider
- Primary target is Civic Auth but should work with others
- Examples demonstrate real-world usage patterns
- Token storage interface allows for custom implementations
- Built for both server-side and CLI applications
