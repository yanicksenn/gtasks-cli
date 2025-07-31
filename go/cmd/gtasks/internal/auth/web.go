package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	oauth2_v2 "google.golang.org/api/oauth2/v2"
)

const (
	authServiceURL = "https://oauth-hub-dev-849914933450.us-central1.run.app/project1"
	redirectURI    = "http://localhost:8080/callback"
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

		accessToken := query.Get("access_token")                                     
        refreshToken := query.Get("refresh_token")                                   
        tokenType := query.Get("token_type")                                         
        expiresInStr := query.Get("expires_in")                                     
        if accessToken == "" {                                                       
            http.Error(w, "Missing access_token", http.StatusBadRequest)             
            errChan <- fmt.Errorf("callback did not return access_token")   
			return
		}

        expiresIn, _ := strconv.Atoi(expiresInStr)                                  
        expiry := time.Now().Add(time.Duration(expiresIn) * time.Second)            
                                                                                    
        token := &oauth2.Token{                                                     
            AccessToken:  accessToken,                                              
            RefreshToken: refreshToken,                                             
            TokenType:    tokenType,                                                
            Expiry:       expiry,                                                   
        }                                                                                                                                    

		fmt.Fprintln(w, "Authentication successful! You can close this window.")
		tokenChan <- token
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

	loginURL := fmt.Sprintf("%s?redirect_uri=%s", authServiceURL, redirectURI)
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
