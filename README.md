# Civic Auth Go SDK

A comprehensive Go SDK for integrating with Civic Auth's OIDC/OAuth2 authentication service.

## Features

- ✅ Full OIDC/OAuth2 implementation
- ✅ PKCE (Proof Key for Code Exchange) support
- ✅ JWT ID token validation with JWK
- ✅ Automatic token refresh
- ✅ Configurable token storage
- ✅ Comprehensive error handling
- ✅ Production-ready examples

## Installation

```bash
go get captured.ventures/civic-auth-go
```

## Quick Start

### Basic Setup

```go
package main

import (
    "captured.ventures/civic-auth-go/pkg/civicauth"
    "log"
)

func main() {
    // Configure the client
    config := civicauth.DefaultConfig()
    config.ClientID = "your-civic-client-id"
    config.ClientSecret = "your-civic-client-secret"
    config.RedirectURL = "http://localhost:8080/callback"
    config.Issuer = "https://auth.civic.com"

    // Create the client
    client, err := civicauth.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Generate authorization URL with PKCE
    authURL, state, codeVerifier, err := client.CreateAuthorizationFlow()
    if err != nil {
        log.Fatalf("Failed to create authorization flow: %v", err)
    }

    // Redirect user to authURL...
    // After callback, exchange code for tokens:
    tokens, err := client.ExchangeCodeForTokens(ctx, code, codeVerifier)
    if err != nil {
        log.Fatalf("Failed to exchange tokens: %v", err)
    }

    // Get user information
    userInfo, err := client.GetUserInfo(ctx, tokens.AccessToken)
    if err != nil {
        log.Fatalf("Failed to get user info: %v", err)
    }

    fmt.Printf("User: %s (%s)\n", userInfo.Name, userInfo.Email)
}
```

## Configuration

The `Config` struct supports the following options:

```go
type Config struct {
    ClientID     string        // Your Civic Auth client ID
    ClientSecret string        // Your Civic Auth client secret
    RedirectURL  string        // Callback URL for your application
    Issuer       string        // OIDC issuer URL (e.g., https://auth.civic.com)
    Scopes       []string      // OAuth2 scopes (default: ["openid", "profile", "email"])
    HTTPClient   *http.Client  // Custom HTTP client (optional)
    Timeout      time.Duration // Request timeout (default: 30 seconds)
}
```

### Environment Variables

For security, you can configure the client using environment variables:

- `CIVIC_CLIENT_ID`: Your client ID
- `CIVIC_CLIENT_SECRET`: Your client secret
- `CIVIC_REDIRECT_URL`: Your callback URL
- `CIVIC_ISSUER`: The OIDC issuer URL

## Authentication Flow

### 1. Generate Authorization URL

```go
// Simple flow with PKCE
authURL, state, codeVerifier, err := client.CreateAuthorizationFlow()

// Or customize the authorization request
opts := &civicauth.AuthCodeURLOptions{
    State:     "your-state-value",
    Nonce:     "your-nonce-value",
    Prompt:    "consent",           // Force consent screen
    MaxAge:    3600,               // Max age of authentication
    LoginHint: "user@example.com",  // Hint about user identity
}
authURL, err := client.GetAuthCodeURL(opts)
```

### 2. Handle Callback

```go
// Extract code and state from callback URL
code := r.URL.Query().Get("code")
state := r.URL.Query().Get("state")

// Validate state parameter (important for security)
if state != expectedState {
    // Handle invalid state
    return
}

// Exchange code for tokens
tokens, err := client.ExchangeCodeForTokens(ctx, code, codeVerifier)
if err != nil {
    // Handle error
    return
}
```

### 3. Validate and Use Tokens

```go
// Create token manager for JWT validation
tokenManager := civicauth.NewTokenManager(client)

// Validate ID token
if tokens.IDToken != "" {
    claims, err := tokenManager.ValidateIDToken(ctx, tokens.IDToken)
    if err != nil {
        // Handle invalid token
        return
    }
    userID := claims.Subject
}

// Get user information using access token
userInfo, err := client.GetUserInfo(ctx, tokens.AccessToken)
if err != nil {
    // Handle error
    return
}
```

## Token Management

### Token Storage

The SDK provides a `TokenStorage` interface for persisting tokens:

```go
type TokenStorage interface {
    Store(userID string, tokens *TokenResponse) error
    Retrieve(userID string) (*TokenResponse, error)
    Delete(userID string) error
}
```

Built-in implementations:

```go
// In-memory storage (for development/testing)
storage := civicauth.NewInMemoryTokenStorage()

// Store tokens
err := storage.Store("user123", tokens)

// Retrieve tokens
tokens, err := storage.Retrieve("user123")
```

### Automatic Token Refresh

Use `TokenRefreshManager` for automatic token refresh:

```go
storage := civicauth.NewInMemoryTokenStorage()
refreshManager := civicauth.NewTokenRefreshManager(client, storage)

// This will automatically refresh the token if needed
validTokens, err := refreshManager.GetValidToken(ctx, "user123")
```

### Manual Token Refresh

```go
// Refresh tokens manually
newTokens, err := client.RefreshToken(ctx, refreshToken)
if err != nil {
    // Handle refresh error (user may need to re-authenticate)
    return
}
```

## ID Token Validation

The SDK automatically validates ID tokens against Civic Auth's public keys:

```go
tokenManager := civicauth.NewTokenManager(client)
claims, err := tokenManager.ValidateIDToken(ctx, idToken)
if err != nil {
    // Token is invalid
    return
}

// Access user claims
userID := claims.Subject
email := claims.Email
name := claims.Name
```

The validation process:
1. Verifies the JWT signature using Civic Auth's public keys
2. Validates the issuer and audience claims
3. Checks token expiration
4. Returns parsed claims

## Logout

Generate a logout URL to properly sign out users:

```go
logoutURL, err := client.GetLogoutURL("http://localhost:8080", idToken)
if err != nil {
    // Handle error
    return
}

// Redirect user to logoutURL
http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
```

## Error Handling

The SDK provides detailed error messages for debugging:

```go
tokens, err := client.ExchangeCodeForTokens(ctx, code, codeVerifier)
if err != nil {
    // Errors include context about what failed
    log.Printf("Token exchange failed: %v", err)
    
    // Check for specific error types
    if strings.Contains(err.Error(), "invalid_grant") {
        // Handle invalid authorization code
    }
    return
}
```

## Examples

### Web Application

See [`examples/web_server.go`](examples/web_server.go) for a complete web server implementation with:
- Login/logout flows
- Session management
- User profile display
- Token refresh handling

Run the example:

```bash
export CIVIC_CLIENT_ID="your-client-id"
export CIVIC_CLIENT_SECRET="your-client-secret"
export CIVIC_ISSUER="https://auth.civic.com"
go run examples/web_server.go
```

### Command Line

See [`examples/cli_example.go`](examples/cli_example.go) for CLI usage patterns.

```bash
go run examples/cli_example.go
```

## Production Considerations

### Security

1. **State Parameter**: Always validate the state parameter to prevent CSRF attacks
2. **PKCE**: Use the PKCE flow for public clients (the SDK does this by default)
3. **HTTPS**: Always use HTTPS in production for redirect URLs
4. **Token Storage**: Use secure, encrypted storage for tokens in production
5. **Token Validation**: Always validate ID tokens before trusting claims

### Performance

1. **HTTP Client**: Reuse HTTP clients and connections
2. **JWK Caching**: The SDK automatically caches JWKs for performance
3. **Token Caching**: Implement proper token storage to avoid unnecessary refreshes

### Error Handling

1. **Retry Logic**: Implement retry logic for network failures
2. **Graceful Degradation**: Handle cases where authentication services are unavailable
3. **User Experience**: Provide clear error messages to users

## API Reference

### Client Methods

- `NewClient(config *Config) (*Client, error)` - Create a new client
- `CreateAuthorizationFlow() (authURL, state, codeVerifier string, err error)` - Generate full auth flow
- `GetAuthCodeURL(opts *AuthCodeURLOptions) (string, error)` - Generate authorization URL
- `ExchangeCodeForTokens(ctx context.Context, code, codeVerifier string) (*TokenResponse, error)` - Exchange code for tokens
- `RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)` - Refresh tokens
- `GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error)` - Get user information
- `GetLogoutURL(postLogoutRedirectURI, idTokenHint string) (string, error)` - Generate logout URL

### Token Manager Methods

- `NewTokenManager(client *Client) *TokenManager` - Create token manager
- `ValidateIDToken(ctx context.Context, idToken string) (*Claims, error)` - Validate ID token

### Storage Methods

- `NewInMemoryTokenStorage() *InMemoryTokenStorage` - Create in-memory storage
- `Store(userID string, tokens *TokenResponse) error` - Store tokens
- `Retrieve(userID string) (*TokenResponse, error)` - Retrieve tokens
- `Delete(userID string) error` - Delete tokens

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Testing

```bash
go test ./pkg/civicauth/...
```

## Project Structure

```
civic-auth-go/
├── pkg/civicauth/        # Core SDK library
│   ├── config.go         # Configuration and types
│   ├── client.go         # Main OIDC client
│   ├── tokens.go         # Token management utilities
│   ├── config_test.go    # Configuration tests
│   └── tokens_test.go    # Token utility tests
├── examples/             # Usage examples
│   ├── web_server.go     # Web application example
│   └── cli_example.go    # CLI example
├── bin/                  # Built executables
├── README.md             # Comprehensive documentation
├── Makefile             # Build and test automation
├── LICENSE              # MIT License
└── go.mod               # Go module definition
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues related to this SDK, please open an issue on GitHub.

For Civic Auth service questions, contact Civic support.
