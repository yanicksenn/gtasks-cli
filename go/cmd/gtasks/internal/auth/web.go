package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	oauth2_v2 "google.golang.org/api/oauth2/v2"
)

const (
	oauthServiceURL = "https://oauth-hub-dev-849914933450.us-central1.run.app/oauth/gtasks_cli.json"
	exchangeServiceURL = "https://oauth-hub-dev-849914933450.us-central1.run.app/exchange/gtasks_cli.json"
	authRedirectURI    = "http://localhost:8080/callback"
)

// LoginViaWebFlow orchestrates the web-based authentication process.
func LoginViaWebFlow(ctx context.Context) (string, error) {
	tokenChan := make(chan *oauth2.Token)
	errChan := make(chan error)

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if errStr := query.Get("error"); errStr != "" {
			errDesc := query.Get("error_description")
			http.Error(w, fmt.Sprintf("Authentication failed: %s - %s", errStr, errDesc), http.StatusBadRequest)
			errChan <- fmt.Errorf("authentication service returned an error: %s (%s)", errStr, errDesc)
			return
		}

		authCode := query.Get("code")
		if authCode == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			errChan <- fmt.Errorf("callback did not return authorization code")
			return
		}

		exchangeURL := fmt.Sprintf("%s?code=%s&redirect_uri=%s", exchangeServiceURL, authCode, authRedirectURI)
		resp, err := http.Post(exchangeURL, "application/json", nil)
		if err != nil {
			http.Error(w, "Failed to exchange authorization code for token", http.StatusInternalServerError)
			errChan <- fmt.Errorf("failed to exchange authorization code for token: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to exchange authorization code for token", http.StatusInternalServerError)
			errChan <- fmt.Errorf("token exchange failed with status: %s", resp.Status)
			return
		}

		var token oauth2.Token
		if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
			http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
			errChan <- fmt.Errorf("failed to decode token response: %w", err)
			return
		}

		fmt.Fprintln(w, "Authentication successful! You can close this window.")
		tokenChan <- &token
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start local server: %w", err)
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	loginURL := fmt.Sprintf("%s?redirect_uri=%s", oauthServiceURL, authRedirectURI)
	fmt.Printf("Your browser should open automatically. If not, please visit:\n%s\n", loginURL)
	if err := browser.OpenURL(loginURL); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not open browser: %v\n", err)
	}

	var token *oauth2.Token
	select {
	case token = <-tokenChan:
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	
	svc, err := oauth2_v2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to create oauth2 service: %w", err)
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user info: %w", err)
	}

	cache, err := loadTokenCache()
	if err != nil {
		return "", err
	}
	if err := cache.save(userInfo.Email, token); err != nil {
		return "", err
	}

	return userInfo.Email, nil
}

