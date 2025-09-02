# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2024-09-02

### Added
- Initial release of Civic Auth Go SDK
- Complete OIDC/OAuth2 client implementation with PKCE support
- JWT ID token validation with automatic JWK key fetching and caching
- Configurable token storage interface with in-memory implementation
- Automatic token refresh capabilities
- Production-ready web server example with complete authentication flow
- Command-line interface example demonstrating SDK usage
- Comprehensive test suite with unit tests
- MIT License for open source compatibility
- Detailed documentation with usage examples and best practices
- Build automation with Makefile
- Support for all standard OIDC scopes and custom parameters
- Secure session management patterns
- Error handling with detailed context
- Environment variable configuration support

### Security
- PKCE (Proof Key for Code Exchange) implementation for enhanced security
- State parameter validation to prevent CSRF attacks
- JWT signature verification using Civic Auth's public keys
- Secure token storage patterns and interfaces
- HTTPS enforcement recommendations in documentation

[Unreleased]: https://github.com/captured-ventures/civic-auth-go/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/captured-ventures/civic-auth-go/releases/tag/v1.0.0
