package civicauth

import (
	"testing"
	"time"
)

func TestInMemoryTokenStorage(t *testing.T) {
	storage := NewInMemoryTokenStorage()

	// Test storing tokens
	tokens := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	err := storage.Store("user123", tokens)
	if err != nil {
		t.Fatalf("Failed to store tokens: %v", err)
	}

	// Test retrieving tokens
	retrievedTokens, err := storage.Retrieve("user123")
	if err != nil {
		t.Fatalf("Failed to retrieve tokens: %v", err)
	}

	if retrievedTokens.AccessToken != tokens.AccessToken {
		t.Errorf("Expected access token %s, got %s", tokens.AccessToken, retrievedTokens.AccessToken)
	}

	if retrievedTokens.RefreshToken != tokens.RefreshToken {
		t.Errorf("Expected refresh token %s, got %s", tokens.RefreshToken, retrievedTokens.RefreshToken)
	}

	// Test deleting tokens
	err = storage.Delete("user123")
	if err != nil {
		t.Fatalf("Failed to delete tokens: %v", err)
	}

	// Verify tokens are deleted
	_, err = storage.Retrieve("user123")
	if err == nil {
		t.Error("Expected error when retrieving deleted tokens, got nil")
	}
}

func TestInMemoryTokenStorageErrors(t *testing.T) {
	storage := NewInMemoryTokenStorage()

	// Test storing with empty user ID
	tokens := &TokenResponse{AccessToken: "test"}
	err := storage.Store("", tokens)
	if err == nil {
		t.Error("Expected error when storing with empty user ID, got nil")
	}

	// Test retrieving with empty user ID
	_, err = storage.Retrieve("")
	if err == nil {
		t.Error("Expected error when retrieving with empty user ID, got nil")
	}

	// Test deleting with empty user ID
	err = storage.Delete("")
	if err == nil {
		t.Error("Expected error when deleting with empty user ID, got nil")
	}

	// Test retrieving non-existent user
	_, err = storage.Retrieve("nonexistent")
	if err == nil {
		t.Error("Expected error when retrieving non-existent user, got nil")
	}
}

func TestIsTokenExpired(t *testing.T) {
	// Test with valid expiry
	tokens := &TokenResponse{ExpiresIn: 3600} // 1 hour
	isExpired := IsTokenExpired(tokens, time.Now().Add(-30*time.Minute))
	if isExpired {
		t.Error("Token should not be expired")
	}

	// Test with expired token
	isExpired = IsTokenExpired(tokens, time.Now().Add(-2*time.Hour))
	if !isExpired {
		t.Error("Token should be expired")
	}

	// Test with no expiry info
	tokens = &TokenResponse{ExpiresIn: 0}
	isExpired = IsTokenExpired(tokens, time.Now().Add(-24*time.Hour))
	if isExpired {
		t.Error("Token without expiry info should not be considered expired")
	}
}
