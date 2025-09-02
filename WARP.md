# WARP.md

This file provides comprehensive context about the Civic Auth Go SDK for Warp's AI terminal and other development tools.

## Project Overview

The Civic Auth Go SDK is a production-ready Go library that provides seamless integration with Civic Auth's OIDC/OAuth2 authentication service. Built with security, performance, and developer experience in mind.

### Key Features
- 🔐 Complete OIDC/OAuth2 client implementation with PKCE support
- 🔑 JWT ID token validation with automatic JWK key fetching and caching
- 🔄 Automatic token refresh capabilities with configurable storage
- 🛡️ Security-first approach with state validation and CSRF protection
- 📦 Minimal dependencies (only standard library + JWT library)
- 🎯 Production-ready with comprehensive error handling
- 📚 Extensive documentation and working examples

## Architecture

### Core Components

1. **Configuration Layer** (`config.go` - 144 lines)
   - `Config` struct with validation
   - Default configuration factory
   - OIDC provider metadata handling
   - Token response and user info structures

2. **Client Layer** (`client.go` - 324 lines)
   - Main `Client` struct for OIDC operations
   - Authorization code flow with PKCE
   - Token exchange and refresh
   - User information retrieval
   - Provider discovery and logout URLs

3. **Token Management** (`tokens.go` - 363 lines)
   - JWT validation with JWK key management
   - Token storage interface and in-memory implementation
   - Automatic token refresh manager
   - RSA public key conversion utilities

### Dependencies

```go
module captured.ventures/civic-auth-go

go 1.24.6

require github.com/golang-jwt/jwt/v5 v5.3.0
```

**Why minimal dependencies?**
- Reduces attack surface and maintenance overhead
- Faster builds and smaller binary size
- Better compatibility across Go versions
- Relies on Go's excellent standard library

## Code Structure

```
pkg/civicauth/           # Core SDK (1,061 total lines)
├── config.go            # Configuration and types (144 lines)
├── client.go            # Main OIDC client (324 lines) 
├── tokens.go            # Token management (363 lines)
├── config_test.go       # Configuration tests (130 lines)
└── tokens_test.go       # Token utility tests (100 lines)

examples/                # Usage demonstrations (394 total lines)
├── web_server.go        # Complete web app example (294 lines)
└── cli_example.go       # CLI usage patterns (100 lines)
```

## API Surface

### Core Types

```go
// Main client for OIDC operations
type Client struct {
    config   *Config
    provider *OIDCProvider
}

// Configuration with sensible defaults
type Config struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
    Issuer       string
    Scopes       []string
    HTTPClient   *http.Client
    Timeout      time.Duration
}

// Token management and JWT validation
type TokenManager struct {
    Client   *Client
    jwkSet   *JWKSet
    jwkCache map[string]*rsa.PublicKey
}
```

### Key Functions

**Client Operations:**
- `NewClient(config *Config) (*Client, error)` - Initialize client
- `CreateAuthorizationFlow()` - Generate auth URL with PKCE
- `ExchangeCodeForTokens()` - Exchange authorization code for tokens
- `RefreshToken()` - Refresh access tokens
- `GetUserInfo()` - Retrieve user information
- `GetLogoutURL()` - Generate logout URLs

**Token Management:**
- `ValidateIDToken()` - Verify JWT signatures and claims
- `NewTokenManager()` - Create token validation manager
- `NewInMemoryTokenStorage()` - Create token storage

**Utilities:**
- `IsTokenExpired()` - Check token expiration
- `DefaultConfig()` - Create configuration with defaults

## Security Model

### Authentication Flow
1. **PKCE Generation**: Creates cryptographically secure code challenge/verifier pairs
2. **State Validation**: Prevents CSRF attacks with random state parameters
3. **JWT Verification**: Validates ID tokens using Civic Auth's public keys
4. **Token Refresh**: Securely handles token lifecycle management

### Security Features
- ✅ PKCE (Proof Key for Code Exchange) by default
- ✅ State parameter validation for CSRF protection
- ✅ JWT signature verification with JWK key rotation support
- ✅ Automatic JWK key caching with refresh capability
- ✅ Secure token storage interface
- ✅ HTTPS enforcement recommendations

### Threat Mitigation
- **CSRF Attacks**: State parameter validation
- **Code Interception**: PKCE implementation
- **Token Theft**: JWT validation and proper storage patterns
- **Replay Attacks**: Nonce support and expiration validation

## Development Workflow

### Setup Commands
```bash
# Install dependencies
make install-deps

# Run tests  
make test

# Build examples
make build

# Run linting
make lint

# Format code
make fmt

# Run all CI checks
make ci
```

### Testing Strategy
- **Unit Tests**: Core functionality with table-driven tests
- **Integration Tests**: Token validation and HTTP client behavior
- **Example Tests**: Verify example code compiles and runs
- **Coverage**: `make test-coverage` generates coverage reports

### Code Quality
- **Linting**: `go vet` + `gofmt` checks
- **Error Handling**: Comprehensive error wrapping with context
- **Documentation**: Godoc comments for all exported functions
- **Examples**: Working code samples for real-world usage

## Usage Patterns

### Basic Authentication Flow
```go
// 1. Create client
client, err := civicauth.NewClient(config)

// 2. Generate authorization URL
authURL, state, codeVerifier, err := client.CreateAuthorizationFlow()

// 3. Handle callback and exchange code
tokens, err := client.ExchangeCodeForTokens(ctx, code, codeVerifier)

// 4. Validate and use tokens
userInfo, err := client.GetUserInfo(ctx, tokens.AccessToken)
```

### Advanced Token Management
```go
// Token storage and refresh
storage := civicauth.NewInMemoryTokenStorage()
refreshManager := civicauth.NewTokenRefreshManager(client, storage)
validTokens, err := refreshManager.GetValidToken(ctx, userID)

// JWT validation
tokenManager := civicauth.NewTokenManager(client)
claims, err := tokenManager.ValidateIDToken(ctx, tokens.IDToken)
```

## Examples

### Web Server (`examples/web_server.go`)
- Complete web application with login/logout flows
- Session management with cookies
- User profile display
- Token refresh handling
- Production-ready error handling

**Key handlers:**
- `/` - Home page with login link
- `/login` - Initiates OAuth flow with PKCE
- `/callback` - Handles OAuth callback and token exchange
- `/profile` - Displays user information (requires auth)
- `/logout` - Clears session and redirects to logout URL

### CLI Example (`examples/cli_example.go`)
- Command-line usage demonstration
- Configuration examples
- Token storage simulation
- Authorization URL generation

## Production Considerations

### Performance
- **HTTP Client Reuse**: Single client instance per application
- **JWK Caching**: Automatic caching of public keys for JWT validation
- **Connection Pooling**: Leverages Go's HTTP client connection pooling
- **Minimal Allocations**: Efficient memory usage patterns

### Scalability
- **Stateless Design**: No server-side state requirements
- **Token Storage Interface**: Pluggable storage backends (Redis, database, etc.)
- **Context Support**: Proper cancellation and timeout handling
- **Concurrent Safe**: Thread-safe operations throughout

### Monitoring
- **Structured Errors**: Detailed error context for debugging
- **HTTP Status Codes**: Proper status code handling and reporting
- **Logging Integration**: Compatible with standard Go logging patterns
- **Metrics Ready**: Easy integration with monitoring systems

## Integration Points

### Storage Backends
```go
type TokenStorage interface {
    Store(userID string, tokens *TokenResponse) error
    Retrieve(userID string) (*TokenResponse, error)
    Delete(userID string) error
}
```

Implement custom storage for:
- Redis for distributed applications
- Database for persistent storage
- Encrypted file storage
- Cloud storage services

### HTTP Middleware
The SDK works well with popular Go web frameworks:
- Gin/Echo middleware for route protection
- net/http middleware for standard library usage
- gRPC interceptors for API services
- Custom authentication middleware

### Logging Integration
```go
// Compatible with popular loggers
import "log/slog"
import "github.com/sirupsen/logrus"
import "go.uber.org/zap"
```

## Environment Configuration

### Required Environment Variables
```bash
CIVIC_CLIENT_ID="your-client-id"           # OAuth2 client ID
CIVIC_CLIENT_SECRET="your-client-secret"   # OAuth2 client secret  
CIVIC_ISSUER="https://auth.civic.com"      # OIDC issuer URL
CIVIC_REDIRECT_URL="http://localhost:8080/callback"  # Callback URL
```

### Optional Configuration
```bash
CIVIC_SCOPES="openid profile email"        # OAuth2 scopes
CIVIC_TIMEOUT="30s"                        # HTTP request timeout
```

## Error Handling

The SDK provides detailed error context:

```go
// Network errors
if err != nil {
    if strings.Contains(err.Error(), "context deadline exceeded") {
        // Handle timeout
    }
}

// Authentication errors  
if strings.Contains(err.Error(), "invalid_grant") {
    // Handle invalid authorization code
}

// Token validation errors
if strings.Contains(err.Error(), "token has expired") {
    // Handle expired tokens
}
```

## Future Roadmap

### Planned Enhancements
- [ ] Device flow support for CLI applications
- [ ] Client credentials flow for service-to-service auth
- [ ] Token introspection endpoint support
- [ ] Metrics and observability integration
- [ ] Additional storage backends (Redis, SQL)
- [ ] Rate limiting and circuit breaker patterns

### Extension Points
- Custom HTTP transports
- Pluggable JWT validators
- Custom claim validation
- Storage encryption layers
- Audit logging hooks

## Troubleshooting

### Common Issues

**"Provider not initialized"**
- Ensure CIVIC_ISSUER is set correctly
- Check network connectivity to issuer
- Verify issuer URL returns valid OIDC metadata

**"Invalid state parameter"**
- State mismatch indicates potential CSRF attack
- Ensure state is properly stored between auth request and callback
- Check session management implementation

**"Token validation failed"**
- Clock skew between client and server
- Invalid issuer or audience claims
- Expired tokens or network issues fetching JWK keys

### Debug Mode
Set detailed logging to debug issues:
```go
config.HTTPClient.Transport = &loggingTransport{}
```

## Community & Support

- **GitHub**: https://github.com/ironystock/civic-auth-go
- **Issues**: Report bugs and feature requests via GitHub issues
- **Contributions**: See CONTRIBUTING.md for development guidelines
- **License**: MIT License - see LICENSE file

## AI Agent Notes

This codebase is well-suited for AI-assisted development:
- Clear separation of concerns
- Comprehensive test coverage
- Detailed error messages
- Extensive documentation
- Production-ready examples
- Security-first design principles

The AGENTS.md file provides additional context for AI coding assistants working on this project.
