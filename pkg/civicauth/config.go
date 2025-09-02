package civicauth

import (
	"fmt"
	"net/http"
	"time"
)

// Config holds the configuration for the Civic Auth OIDC client
type Config struct {
	// ClientID is the OAuth2 client ID for your application
	ClientID string

	// ClientSecret is the OAuth2 client secret for your application
	ClientSecret string

	// RedirectURL is the callback URL where users will be redirected after authentication
	RedirectURL string

	// Issuer is the OIDC issuer URL (e.g., https://auth.civic.com)
	Issuer string

	// Scopes are the OAuth2 scopes to request (default: ["openid", "profile", "email"])
	Scopes []string

	// HTTPClient is the HTTP client to use for requests (optional)
	HTTPClient *http.Client

	// Timeout for HTTP requests (default: 30 seconds)
	Timeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Scopes:     []string{"openid", "profile", "email"},
		HTTPClient: &http.Client{},
		Timeout:    30 * time.Second,
	}
}

// Validate checks that the configuration is valid
func (c *Config) Validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("client ID is required")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("client secret is required")
	}
	if c.RedirectURL == "" {
		return fmt.Errorf("redirect URL is required")
	}
	if c.Issuer == "" {
		return fmt.Errorf("issuer URL is required")
	}
	if len(c.Scopes) == 0 {
		c.Scopes = []string{"openid", "profile", "email"}
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}

	c.HTTPClient.Timeout = c.Timeout

	return nil
}

// OIDCProvider represents the OIDC provider metadata
type OIDCProvider struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	JwksURI               string `json:"jwks_uri"`
	EndSessionEndpoint    string `json:"end_session_endpoint,omitempty"`
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
}

// UserInfo represents the OIDC user information
type UserInfo struct {
	Sub               string `json:"sub"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Website           string `json:"website,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Gender            string `json:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Zoneinfo          string `json:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	PhoneVerified     bool   `json:"phone_number_verified,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}

// Claims represents ID token claims
type Claims struct {
	Issuer       string `json:"iss"`
	Subject      string `json:"sub"`
	Audience     string `json:"aud"`
	Expiry       int64  `json:"exp"`
	IssuedAt     int64  `json:"iat"`
	Nonce        string `json:"nonce,omitempty"`
	AuthTime     int64  `json:"auth_time,omitempty"`
	SessionState string `json:"session_state,omitempty"`

	// Standard profile claims
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Website           string `json:"website,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Gender            string `json:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Zoneinfo          string `json:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	PhoneVerified     bool   `json:"phone_number_verified,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}
