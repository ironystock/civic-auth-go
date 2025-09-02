package main

import (
	"fmt"
	"log"
	"os"

	"captured.ventures/civic-auth-go/pkg/civicauth"
)

func main() {
	// Create configuration
	config := civicauth.DefaultConfig()
	config.ClientID = getEnv("CIVIC_CLIENT_ID", "your-client-id")
	config.ClientSecret = getEnv("CIVIC_CLIENT_SECRET", "your-client-secret")
	config.RedirectURL = getEnv("CIVIC_REDIRECT_URL", "http://localhost:8080/callback")
	config.Issuer = getEnv("CIVIC_ISSUER", "https://auth.civicauth.com")

	// Create client
	client, err := civicauth.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 1: Generate authorization URL with PKCE
	fmt.Println("=== Authorization Code Flow with PKCE ===")
	authURL, state, codeVerifier, err := client.CreateAuthorizationFlow()
	if err != nil {
		log.Fatalf("Failed to create authorization flow: %v", err)
	}

	fmt.Printf("1. Visit this URL to authorize the application:\n%s\n\n", authURL)
	fmt.Printf("2. State parameter: %s\n", state)
	fmt.Printf("3. Code verifier (keep secret): %s\n\n", codeVerifier)

	// In a real CLI app, you would:
	// 1. Open the browser with the auth URL
	// 2. Start a temporary HTTP server to receive the callback
	// 3. Extract the authorization code from the callback
	// 4. Exchange the code for tokens

	// Example 2: Manual authorization URL generation
	fmt.Println("=== Manual Authorization URL Generation ===")
	opts := &civicauth.AuthCodeURLOptions{
		State:  "my-custom-state",
		Prompt: "consent", // Force consent screen
	}

	manualAuthURL, err := client.GetAuthCodeURL(opts)
	if err != nil {
		log.Fatalf("Failed to generate auth URL: %v", err)
	}
	fmt.Printf("Manual auth URL: %s\n\n", manualAuthURL)

	// Example 3: Token validation (if you have an ID token)
	fmt.Println("=== Token Management Example ===")
	_ = civicauth.NewTokenManager(client) // tokenManager created for demonstration
	fmt.Println("Token manager created - ready to validate ID tokens")

	// Example 4: Storage demonstration
	fmt.Println("\n=== Token Storage Example ===")
	storage := civicauth.NewInMemoryTokenStorage()

	// Simulate storing tokens
	sampleTokens := &civicauth.TokenResponse{
		AccessToken:  "sample-access-token",
		RefreshToken: "sample-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	err = storage.Store("user123", sampleTokens)
	if err != nil {
		log.Printf("Failed to store tokens: %v", err)
	} else {
		fmt.Println("Tokens stored successfully")
	}

	// Retrieve tokens
	retrievedTokens, err := storage.Retrieve("user123")
	if err != nil {
		log.Printf("Failed to retrieve tokens: %v", err)
	} else {
		fmt.Printf("Retrieved tokens for user: access_token=%s, token_type=%s\n",
			retrievedTokens.AccessToken, retrievedTokens.TokenType)
	}

	fmt.Println("\n=== Configuration Summary ===")
	fmt.Printf("Client ID: %s\n", config.ClientID)
	fmt.Printf("Issuer: %s\n", config.Issuer)
	fmt.Printf("Redirect URL: %s\n", config.RedirectURL)
	fmt.Printf("Scopes: %v\n", config.Scopes)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
