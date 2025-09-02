package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"captured.ventures/civic-auth-go/pkg/civicauth"
)

// Session storage (in production, use a proper session store)
var sessions = make(map[string]*SessionData)

type SessionData struct {
	State        string
	CodeVerifier string
	UserID       string
}

func main() {
	// Get configuration from environment variables
	config := civicauth.DefaultConfig()
	config.ClientID = getEnv("CIVIC_CLIENT_ID", "your-client-id")
	config.ClientSecret = getEnv("CIVIC_CLIENT_SECRET", "your-client-secret")
	config.RedirectURL = getEnv("CIVIC_REDIRECT_URL", "http://localhost:8080/callback")
	config.Issuer = getEnv("CIVIC_ISSUER", "https://auth.civicauth.com")

	// Create the Civic Auth client
	client, err := civicauth.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Civic Auth client: %v", err)
	}

	// Create token storage and managers
	storage := civicauth.NewInMemoryTokenStorage()
	tokenManager := civicauth.NewTokenManager(client)
	refreshManager := civicauth.NewTokenRefreshManager(client, storage)

	// Set up HTTP handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler(client))
	http.HandleFunc("/callback", callbackHandler(client, tokenManager, storage))
	http.HandleFunc("/profile", profileHandler(refreshManager))
	http.HandleFunc("/logout", logoutHandler(client))

	fmt.Println("Starting server on :8080")
	fmt.Println("Visit http://localhost:8080 to test the integration")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Civic Auth Integration Demo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .button { display: inline-block; padding: 10px 20px; margin: 10px 0; 
                 background: #007cba; color: white; text-decoration: none; border-radius: 5px; }
        .button:hover { background: #005a87; }
    </style>
</head>
<body>
    <h1>Civic Auth Integration Demo</h1>
    <p>This is a demo of the Civic Auth OIDC integration.</p>
    <a href="/login" class="button">Login with Civic Auth</a>
    <br><br>
    <a href="/profile" class="button">View Profile (requires login)</a>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func loginHandler(client *civicauth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate authorization flow parameters
		authURL, state, codeVerifier, err := client.CreateAuthorizationFlow()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create authorization flow: %v", err), http.StatusInternalServerError)
			return
		}

		// Store session data (in production, use proper session management)
		sessionID := generateSessionID()
		sessions[sessionID] = &SessionData{
			State:        state,
			CodeVerifier: codeVerifier,
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
		})

		// Redirect to authorization URL
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

func callbackHandler(client *civicauth.Client, tokenManager *civicauth.TokenManager, storage civicauth.TokenStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session data
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Session not found", http.StatusBadRequest)
			return
		}

		session, exists := sessions[cookie.Value]
		if !exists {
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}

		// Validate state parameter
		state := r.URL.Query().Get("state")
		if state != session.State {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		// Get authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			errorParam := r.URL.Query().Get("error")
			errorDesc := r.URL.Query().Get("error_description")
			http.Error(w, fmt.Sprintf("Authorization failed: %s - %s", errorParam, errorDesc), http.StatusBadRequest)
			return
		}

		// Exchange code for tokens
		tokens, err := client.ExchangeCodeForTokens(r.Context(), code, session.CodeVerifier)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to exchange code for tokens: %v", err), http.StatusInternalServerError)
			return
		}

		// Validate ID token if present
		var userID string
		if tokens.IDToken != "" {
			claims, err := tokenManager.ValidateIDToken(r.Context(), tokens.IDToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to validate ID token: %v", err), http.StatusInternalServerError)
				return
			}
			userID = claims.Subject
		} else {
			// If no ID token, get user info from userinfo endpoint
			userInfo, err := client.GetUserInfo(r.Context(), tokens.AccessToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
				return
			}
			userID = userInfo.Sub
		}

		// Store tokens
		if err := storage.Store(userID, tokens); err != nil {
			http.Error(w, fmt.Sprintf("Failed to store tokens: %v", err), http.StatusInternalServerError)
			return
		}

		// Update session with user ID
		session.UserID = userID

		// Redirect to profile page
		http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
	}
}

func profileHandler(refreshManager *civicauth.TokenRefreshManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session data
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		session, exists := sessions[cookie.Value]
		if !exists || session.UserID == "" {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// Get valid tokens (will refresh if needed)
		tokens, err := refreshManager.GetValidToken(r.Context(), session.UserID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get valid tokens: %v", err), http.StatusInternalServerError)
			return
		}

		// Get user information
		userInfo, err := refreshManager.Client.GetUserInfo(r.Context(), tokens.AccessToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
			return
		}

		// Display user profile
		html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>User Profile - Civic Auth Demo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .profile { background: #f5f5f5; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .button { display: inline-block; padding: 10px 20px; margin: 10px 5px; 
                 background: #007cba; color: white; text-decoration: none; border-radius: 5px; }
        .logout { background: #dc3545; }
        .button:hover { opacity: 0.8; }
    </style>
</head>
<body>
    <h1>User Profile</h1>
    <div class="profile">
        <h3>User Information</h3>
        <p><strong>Subject ID:</strong> %s</p>
        <p><strong>Name:</strong> %s</p>
        <p><strong>Email:</strong> %s</p>
        <p><strong>Email Verified:</strong> %t</p>
        <p><strong>Username:</strong> %s</p>
        <p><strong>Profile:</strong> %s</p>
    </div>
    <a href="/" class="button">Home</a>
    <a href="/logout" class="button logout">Logout</a>
</body>
</html>`,
			userInfo.Sub,
			userInfo.Name,
			userInfo.Email,
			userInfo.EmailVerified,
			userInfo.PreferredUsername,
			userInfo.Profile,
		)

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}
}

func logoutHandler(client *civicauth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session data
		cookie, err := r.Cookie("session_id")
		if err == nil {
			if session, exists := sessions[cookie.Value]; exists && session.UserID != "" {
				// Clear stored tokens
				// In a real implementation, you'd get the ID token for logout
				delete(sessions, cookie.Value)
			}
		}

		// Clear session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		})

		// Generate logout URL (you would typically include the ID token hint)
		logoutURL, err := client.GetLogoutURL("http://localhost:8080", "")
		if err != nil {
			// If logout URL generation fails, just redirect to home
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Redirect to logout URL
		http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateSessionID() string {
	// In production, use a proper session ID generator
	return fmt.Sprintf("session_%d", len(sessions))
}
