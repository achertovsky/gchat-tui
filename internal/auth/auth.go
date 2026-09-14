// Package auth provides Google OAuth authentication for the application.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var Scopes = []string{
	"https://www.googleapis.com/auth/chat.spaces.readonly",
	"https://www.googleapis.com/auth/chat.messages.readonly",
	"https://www.googleapis.com/auth/chat.messages.create",
}

var ErrNotAuthenticated = errors.New("not authenticated; run `chat login` first")

// Provider exposes authenticated tokens without coupling callers to OAuth.
type Provider interface {
	Token(context.Context) (*oauth2.Token, error)
	HTTPClient(context.Context) *http.Client
	Logout() error
}

// AuthorizationPrompt tells the user where to continue the OAuth flow.
type AuthorizationPrompt func(string)

// OAuth is a user-authorized Google OAuth provider.
type OAuth struct {
	config oauth2.Config
	store  tokenStore
	prompt AuthorizationPrompt
}

// New loads desktop OAuth credentials and creates an authentication provider.
func New(clientConfigPath, dataDir string, prompt AuthorizationPrompt) (*OAuth, error) {
	contents, err := os.ReadFile(clientConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read OAuth client configuration: %w", err)
	}
	config, err := google.ConfigFromJSON(contents, Scopes...)
	if err != nil {
		return nil, fmt.Errorf("parse OAuth client configuration: %w", err)
	}
	return &OAuth{config: *config, store: newTokenStore(dataDir), prompt: prompt}, nil
}

// Login completes a browser-based OAuth authorization flow and saves the
// resulting credentials.
func (o *OAuth) Login(ctx context.Context) error {
	token, err := o.authorize(ctx)
	if err != nil {
		return err
	}
	if err := o.store.Save(token); err != nil {
		return fmt.Errorf("store OAuth credentials: %w", err)
	}
	return nil
}

// Token obtains a valid access token, refreshing it when necessary.
func (o *OAuth) Token(ctx context.Context) (*oauth2.Token, error) {
	token, err := o.store.Load()
	if err != nil {
		if errors.Is(err, errTokenNotFound) {
			return nil, ErrNotAuthenticated
		}
		return nil, fmt.Errorf("load OAuth credentials: %w", err)
	}

	refreshed, err := o.config.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, fmt.Errorf("refresh OAuth access token: %w", err)
	}
	if err := o.store.Save(refreshed); err != nil {
		return nil, fmt.Errorf("store refreshed OAuth credentials: %w", err)
	}
	return refreshed, nil
}

// HTTPClient returns an HTTP client that adds and refreshes OAuth credentials.
func (o *OAuth) HTTPClient(ctx context.Context) *http.Client {
	return oauth2.NewClient(ctx, oauth2.TokenSource(tokenSource{provider: o, context: ctx}))
}

// Logout removes all locally stored OAuth credentials.
func (o *OAuth) Logout() error {
	if err := o.store.Delete(); err != nil {
		return fmt.Errorf("delete OAuth credentials: %w", err)
	}
	return nil
}

// Logout removes credentials without requiring the OAuth client configuration.
func Logout(dataDir string) error {
	return (&OAuth{store: newTokenStore(dataDir)}).Logout()
}

func (o *OAuth) authorize(ctx context.Context) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("start OAuth callback listener: %w", err)
	}
	defer listener.Close()

	callbackURL := "http://" + listener.Addr().String() + "/oauth2/callback"
	config := o.config
	config.RedirectURL = callbackURL
	state, err := randomValue(32)
	if err != nil {
		return nil, fmt.Errorf("create OAuth state: %w", err)
	}
	verifier, err := randomValue(64)
	if err != nil {
		return nil, fmt.Errorf("create PKCE verifier: %w", err)
	}

	result := make(chan authorizationResult, 1)
	server := &http.Server{Handler: callbackHandler(state, result)}
	go server.Serve(listener)
	defer server.Close()

	url := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)
	_ = openBrowser(url)
	if o.prompt != nil {
		o.prompt(url)
	}

	select {
	case response := <-result:
		if response.err != nil {
			return nil, response.err
		}
		token, err := config.Exchange(ctx, response.code, oauth2.VerifierOption(verifier))
		if err != nil {
			return nil, fmt.Errorf("exchange OAuth authorization code: %w", err)
		}
		return token, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("OAuth authorization canceled: %w", ctx.Err())
	}
}

type authorizationResult struct {
	code string
	err  error
}

func callbackHandler(expectedState string, result chan<- authorizationResult) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/oauth2/callback" {
			http.NotFound(writer, request)
			return
		}
		if subtle.ConstantTimeCompare([]byte(expectedState), []byte(request.URL.Query().Get("state"))) != 1 {
			http.Error(writer, "OAuth state validation failed.", http.StatusBadRequest)
			result <- authorizationResult{err: errors.New("OAuth state validation failed")}
			return
		}
		if oauthError := request.URL.Query().Get("error"); oauthError != "" {
			http.Error(writer, "OAuth authorization was denied.", http.StatusForbidden)
			result <- authorizationResult{err: fmt.Errorf("OAuth authorization failed: %s", oauthError)}
			return
		}
		code := request.URL.Query().Get("code")
		if code == "" {
			http.Error(writer, "OAuth response did not include an authorization code.", http.StatusBadRequest)
			result <- authorizationResult{err: errors.New("OAuth response did not include an authorization code")}
			return
		}
		fmt.Fprintln(writer, "Authentication complete. You can return to gchat-tui.")
		result <- authorizationResult{code: code}
	})
}

func randomValue(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func openBrowser(url string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	return command.Start()
}

type tokenSource struct {
	provider *OAuth
	context  context.Context
}

func (s tokenSource) Token() (*oauth2.Token, error) {
	return s.provider.Token(s.context)
}

var _ Provider = (*OAuth)(nil)
