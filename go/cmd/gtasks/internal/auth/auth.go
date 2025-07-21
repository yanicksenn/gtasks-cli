package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2_v2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

const (
	tokenFile         = "gtasks-token.json"
	credentialsFile   = "credentials.json"
	redirectURL       = "http://localhost:8080"
)

// ErrCredentialsNotFound is returned when the user's credentials are not found.
var ErrCredentialsNotFound = errors.New("credentials not found. Please run 'gtasks login'")

// Credentials represents the structure of the credentials.json file.
type Credentials struct {
	Installed struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"installed"`
}

// Authenticator handles the OAuth2 flow.
type Authenticator struct {
	config *oauth2.Config
}

// NewAuthenticator creates a new Authenticator.
func NewAuthenticator() (*Authenticator, error) {
	creds, err := loadCredentials()
	if err != nil {
		return nil, fmt.Errorf("could not load credentials: %w", err)
	}

	config := &oauth2.Config{
		ClientID:     creds.Installed.ClientID,
		ClientSecret: creds.Installed.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/tasks", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return &Authenticator{config: config}, nil
}

// GetClient returns an authenticated HTTP client for the given user.
// If the user has no token, it initiates the login flow.
func (a *Authenticator) GetClient(ctx context.Context, user string) (*http.Client, error) {
	cache, err := loadTokenCache()
	if err != nil {
		return nil, err
	}

	token, ok := cache.Tokens[user]
	if !ok {
		return nil, ErrCredentialsNotFound
	}

	return a.config.Client(ctx, token), nil
}

// NewClient performs the full OAuth2 flow to get a new token and returns an authenticated client.
func (a *Authenticator) NewClient(ctx context.Context) (*http.Client, string, error) {
	authURL := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	// Start a local server to listen for the authorization code.
	codeChan := make(chan string)
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintf(w, "Login successful! You can close this window.")
		codeChan <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Error running local server: %v\n", err)
		}
	}()

	code := <-codeChan
	server.Shutdown(ctx)

	token, err := a.config.Exchange(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	client := a.config.Client(ctx, token)
	svc, err := oauth2_v2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, "", fmt.Errorf("unable to create oauth2 service: %w", err)
	}

	userInfo, err := svc.Userinfo.Get().Do()
	if err != nil {
		return nil, "", fmt.Errorf("unable to retrieve user info: %w", err)
	}

	cache, err := loadTokenCache()
	if err != nil {
		return nil, "", err
	}
	if err := cache.save(userInfo.Email, token); err != nil {
		return nil, "", err
	}

	return client, userInfo.Email, nil
}

// TokenCache represents the structure of the credentials file.
type TokenCache struct {
	Tokens map[string]*oauth2.Token `json:"tokens"`
}

func getCredentialsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", credentialsFile), nil
}

func loadCredentials() (*Credentials, error) {
	path, err := getCredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}


func getTokenCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", tokenFile), nil
}

func loadTokenCache() (*TokenCache, error) {
	path, err := getTokenCachePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &TokenCache{Tokens: make(map[string]*oauth2.Token)}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache TokenCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func (tc *TokenCache) save(user string, token *oauth2.Token) error {
	tc.Tokens[user] = token
	path, err := getTokenCachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}