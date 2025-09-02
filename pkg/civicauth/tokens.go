package civicauth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWK represents a JSON Web Key
type JWK struct {
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
	X5t string   `json:"x5t"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// JWKSet represents a set of JSON Web Keys
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// TokenManager handles token operations
type TokenManager struct {
	Client  *Client
	jwkSet  *JWKSet
	jwkCache map[string]*rsa.PublicKey
}

// NewTokenManager creates a new token manager
func NewTokenManager(client *Client) *TokenManager {
	return &TokenManager{
		Client:  client,
		jwkCache: make(map[string]*rsa.PublicKey),
	}
}

// fetchJWKSet fetches the JWK set from the provider
func (tm *TokenManager) fetchJWKSet(ctx context.Context) error {
	if tm.Client.provider == nil {
		return fmt.Errorf("provider not initialized")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", tm.Client.provider.JwksURI, nil)
	if err != nil {
		return fmt.Errorf("failed to create JWK request: %w", err)
	}

	resp, err := tm.Client.config.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWK set: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWK request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read JWK response: %w", err)
	}

	var jwkSet JWKSet
	if err := json.Unmarshal(body, &jwkSet); err != nil {
		return fmt.Errorf("failed to decode JWK set: %w", err)
	}

	tm.jwkSet = &jwkSet
	return nil
}

// getPublicKey gets the public key for the given key ID
func (tm *TokenManager) getPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	// Check cache first
	if key, exists := tm.jwkCache[kid]; exists {
		return key, nil
	}

	// Fetch JWK set if not already done
	if tm.jwkSet == nil {
		if err := tm.fetchJWKSet(ctx); err != nil {
			return nil, err
		}
	}

	// Find the key with matching kid
	var jwk *JWK
	for _, key := range tm.jwkSet.Keys {
		if key.Kid == kid {
			jwk = &key
			break
		}
	}

	if jwk == nil {
		// Refetch JWK set in case it was updated
		if err := tm.fetchJWKSet(ctx); err != nil {
			return nil, fmt.Errorf("failed to refetch JWK set: %w", err)
		}

		for _, key := range tm.jwkSet.Keys {
			if key.Kid == kid {
				jwk = &key
				break
			}
		}

		if jwk == nil {
			return nil, fmt.Errorf("key with kid %s not found", kid)
		}
	}

	// Convert JWK to RSA public key
	publicKey, err := tm.jwkToRSAPublicKey(jwk)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JWK to RSA public key: %w", err)
	}

	// Cache the key
	tm.jwkCache[kid] = publicKey

	return publicKey, nil
}

// jwkToRSAPublicKey converts a JWK to an RSA public key
func (tm *TokenManager) jwkToRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	// Try X.509 certificate first
	if len(jwk.X5c) > 0 {
		certData, err := base64.StdEncoding.DecodeString(jwk.X5c[0])
		if err != nil {
			return nil, fmt.Errorf("failed to decode X.509 certificate: %w", err)
		}

		cert, err := x509.ParseCertificate(certData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse X.509 certificate: %w", err)
		}

		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate does not contain RSA public key")
		}

		return rsaKey, nil
	}

	// Fall back to N and E parameters
	if jwk.N == "" || jwk.E == "" {
		return nil, fmt.Errorf("JWK missing required parameters")
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N parameter: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E parameter: %w", err)
	}

	// Convert bytes to big integers
	n := new(rsa.PublicKey)
	n.N = new(rsa.PublicKey).N.SetBytes(nBytes)

	// E is usually 65537, but decode from bytes to be safe
	e := 0
	for _, b := range eBytes {
		e = e*256 + int(b)
	}
	n.E = e

	return n, nil
}

// ValidateIDToken validates an ID token
func (tm *TokenManager) ValidateIDToken(ctx context.Context, idToken string) (*Claims, error) {
	// Parse the token without verification first to get the header
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		// Get the key ID from the token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("token header missing kid")
		}

		// Get the public key
		publicKey, err := tm.getPublicKey(ctx, kid)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		// Ensure the signing method is RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse and verify ID token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	// Convert claims to our Claims struct
	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to get token claims")
	}

	claims := &Claims{}

	// Convert map claims to struct
	claimsJSON, err := json.Marshal(claimsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	if err := json.Unmarshal(claimsJSON, claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Validate issuer
	if claims.Issuer != tm.Client.config.Issuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", tm.Client.config.Issuer, claims.Issuer)
	}

	// Validate audience
	if claims.Audience != tm.Client.config.ClientID {
		return nil, fmt.Errorf("invalid audience: expected %s, got %s", tm.Client.config.ClientID, claims.Audience)
	}

	// Validate expiry
	if time.Unix(claims.Expiry, 0).Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// IsTokenExpired checks if a token is expired based on the expires_in value
func IsTokenExpired(tokenResp *TokenResponse, issuedAt time.Time) bool {
	if tokenResp.ExpiresIn <= 0 {
		return false // No expiry information
	}

	expiryTime := issuedAt.Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return time.Now().After(expiryTime)
}

// TokenStorage interface for storing and retrieving tokens
type TokenStorage interface {
	Store(userID string, tokens *TokenResponse) error
	Retrieve(userID string) (*TokenResponse, error)
	Delete(userID string) error
}

// InMemoryTokenStorage is a simple in-memory token storage implementation
type InMemoryTokenStorage struct {
	tokens map[string]*TokenResponse
}

// NewInMemoryTokenStorage creates a new in-memory token storage
func NewInMemoryTokenStorage() *InMemoryTokenStorage {
	return &InMemoryTokenStorage{
		tokens: make(map[string]*TokenResponse),
	}
}

// Store stores tokens for a user
func (s *InMemoryTokenStorage) Store(userID string, tokens *TokenResponse) error {
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}
	s.tokens[userID] = tokens
	return nil
}

// Retrieve retrieves tokens for a user
func (s *InMemoryTokenStorage) Retrieve(userID string) (*TokenResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	
	tokens, exists := s.tokens[userID]
	if !exists {
		return nil, errors.New("tokens not found for user")
	}
	
	return tokens, nil
}

// Delete deletes tokens for a user
func (s *InMemoryTokenStorage) Delete(userID string) error {
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}
	
	delete(s.tokens, userID)
	return nil
}

// TokenRefreshManager automatically refreshes tokens when needed
type TokenRefreshManager struct {
	Client  *Client
	storage TokenStorage
}

// NewTokenRefreshManager creates a new token refresh manager
func NewTokenRefreshManager(client *Client, storage TokenStorage) *TokenRefreshManager {
	return &TokenRefreshManager{
		Client:  client,
		storage: storage,
	}
}

// GetValidToken gets a valid access token, refreshing if necessary
func (trm *TokenRefreshManager) GetValidToken(ctx context.Context, userID string) (*TokenResponse, error) {
	// Retrieve stored tokens
	tokens, err := trm.storage.Retrieve(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tokens: %w", err)
	}

	// For simplicity, we'll assume token needs refresh if we have a refresh token
	// In a real implementation, you'd check the token's expiry time
	if tokens.RefreshToken != "" {
		// Try to refresh the token
		newTokens, err := trm.Client.RefreshToken(ctx, tokens.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// If no new refresh token was provided, keep the old one
		if newTokens.RefreshToken == "" {
			newTokens.RefreshToken = tokens.RefreshToken
		}

		// Store the new tokens
		if err := trm.storage.Store(userID, newTokens); err != nil {
			return nil, fmt.Errorf("failed to store refreshed tokens: %w", err)
		}

		return newTokens, nil
	}

	return tokens, nil
}
