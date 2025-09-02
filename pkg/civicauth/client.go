package civicauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is the main OIDC client for Civic Auth
type Client struct {
	config   *Config
	provider *OIDCProvider
}

// NewClient creates a new Civic Auth OIDC client
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	client := &Client{
		config: config,
	}

	// Discover OIDC provider metadata
	if err := client.discoverProvider(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to discover provider: %w", err)
	}

	return client, nil
}

// discoverProvider fetches the OIDC provider metadata
func (c *Client) discoverProvider(ctx context.Context) error {
	wellKnownURL := strings.TrimSuffix(c.config.Issuer, "/") + "/.well-known/openid_configuration"
	
	req, err := http.NewRequestWithContext(ctx, "GET", wellKnownURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch provider metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("provider metadata request failed with status: %d", resp.StatusCode)
	}

	var provider OIDCProvider
	if err := json.NewDecoder(resp.Body).Decode(&provider); err != nil {
		return fmt.Errorf("failed to decode provider metadata: %w", err)
	}

	c.provider = &provider
	return nil
}

// generateCodeChallenge generates a PKCE code challenge
func generateCodeChallenge() (codeVerifier, codeChallenge string, err error) {
	// Generate code verifier (43-128 characters)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	codeVerifier = base64.RawURLEncoding.EncodeToString(b)

	// Generate code challenge (SHA256 hash of verifier)
	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge = base64.RawURLEncoding.EncodeToString(h[:])

	return codeVerifier, codeChallenge, nil
}

// generateState generates a random state parameter
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// AuthCodeURLOptions holds options for generating the authorization URL
type AuthCodeURLOptions struct {
	State         string
	Nonce         string
	CodeChallenge string
	Prompt        string // none, login, consent, select_account
	MaxAge        int    // Maximum age of authentication in seconds
	LoginHint     string // Hint about the user's identity
}

// GetAuthCodeURL generates the authorization URL for the OAuth2 flow
func (c *Client) GetAuthCodeURL(opts *AuthCodeURLOptions) (string, error) {
	if c.provider == nil {
		return "", fmt.Errorf("provider not initialized")
	}

	params := url.Values{
		"response_type":   []string{"code"},
		"client_id":       []string{c.config.ClientID},
		"redirect_uri":    []string{c.config.RedirectURL},
		"scope":           []string{strings.Join(c.config.Scopes, " ")},
		"response_mode":   []string{"query"},
	}

	// Add optional parameters
	if opts != nil {
		if opts.State != "" {
			params.Set("state", opts.State)
		}
		if opts.Nonce != "" {
			params.Set("nonce", opts.Nonce)
		}
		if opts.CodeChallenge != "" {
			params.Set("code_challenge", opts.CodeChallenge)
			params.Set("code_challenge_method", "S256")
		}
		if opts.Prompt != "" {
			params.Set("prompt", opts.Prompt)
		}
		if opts.MaxAge > 0 {
			params.Set("max_age", fmt.Sprintf("%d", opts.MaxAge))
		}
		if opts.LoginHint != "" {
			params.Set("login_hint", opts.LoginHint)
		}
	}

	authURL := c.provider.AuthorizationEndpoint + "?" + params.Encode()
	return authURL, nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens
func (c *Client) ExchangeCodeForTokens(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	if c.provider == nil {
		return nil, fmt.Errorf("provider not initialized")
	}

	data := url.Values{
		"grant_type":    []string{"authorization_code"},
		"client_id":     []string{c.config.ClientID},
		"client_secret": []string{c.config.ClientSecret},
		"code":          []string{code},
		"redirect_uri":  []string{c.config.RedirectURL},
	}

	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.provider.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshToken refreshes an access token using a refresh token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if c.provider == nil {
		return nil, fmt.Errorf("provider not initialized")
	}

	data := url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{c.config.ClientID},
		"client_secret": []string{c.config.ClientSecret},
		"refresh_token": []string{refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.provider.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo retrieves user information using an access token
func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	if c.provider == nil {
		return nil, fmt.Errorf("provider not initialized")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.provider.UserinfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %w", err)
	}

	return &userInfo, nil
}

// GetLogoutURL generates the logout URL
func (c *Client) GetLogoutURL(postLogoutRedirectURI, idTokenHint string) (string, error) {
	if c.provider == nil || c.provider.EndSessionEndpoint == "" {
		return "", fmt.Errorf("logout endpoint not available")
	}

	params := url.Values{}
	
	if postLogoutRedirectURI != "" {
		params.Set("post_logout_redirect_uri", postLogoutRedirectURI)
	}
	
	if idTokenHint != "" {
		params.Set("id_token_hint", idTokenHint)
	}

	logoutURL := c.provider.EndSessionEndpoint
	if len(params) > 0 {
		logoutURL += "?" + params.Encode()
	}

	return logoutURL, nil
}

// Helper function to create a full authorization flow with PKCE
func (c *Client) CreateAuthorizationFlow() (authURL, state, codeVerifier string, err error) {
	// Generate state parameter
	state, err = generateState()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Generate PKCE parameters
	codeVerifier, codeChallenge, err := generateCodeChallenge()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate code challenge: %w", err)
	}

	// Generate authorization URL
	opts := &AuthCodeURLOptions{
		State:         state,
		CodeChallenge: codeChallenge,
	}

	authURL, err = c.GetAuthCodeURL(opts)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate auth URL: %w", err)
	}

	return authURL, state, codeVerifier, nil
}
