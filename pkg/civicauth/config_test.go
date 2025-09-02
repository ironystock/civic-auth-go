package civicauth

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	expectedScopes := []string{"openid", "profile", "email"}
	if len(config.Scopes) != len(expectedScopes) {
		t.Errorf("Expected %d scopes, got %d", len(expectedScopes), len(config.Scopes))
	}

	for i, scope := range expectedScopes {
		if config.Scopes[i] != scope {
			t.Errorf("Expected scope %s at index %d, got %s", scope, i, config.Scopes[i])
		}
	}

	if config.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout of 30s, got %v", config.Timeout)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/callback",
				Issuer:       "https://auth.civic.com",
			},
			expectError: false,
		},
		{
			name: "missing client ID",
			config: &Config{
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/callback",
				Issuer:       "https://auth.civic.com",
			},
			expectError: true,
		},
		{
			name: "missing client secret",
			config: &Config{
				ClientID:    "test-client-id",
				RedirectURL: "http://localhost:8080/callback",
				Issuer:      "https://auth.civic.com",
			},
			expectError: true,
		},
		{
			name: "missing redirect URL",
			config: &Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Issuer:       "https://auth.civic.com",
			},
			expectError: true,
		},
		{
			name: "missing issuer",
			config: &Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				RedirectURL:  "http://localhost:8080/callback",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError && err == nil {
				t.Error("Expected validation error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	config := &Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Issuer:       "https://auth.civic.com",
	}

	err := config.Validate()
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Check that defaults were applied
	expectedScopes := []string{"openid", "profile", "email"}
	if len(config.Scopes) != len(expectedScopes) {
		t.Errorf("Expected %d default scopes, got %d", len(expectedScopes), len(config.Scopes))
	}

	if config.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}

	if config.Timeout != 30*time.Second {
		t.Error("Expected default timeout to be set")
	}
}
