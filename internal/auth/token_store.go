package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	keyringService = "gchat-tui"
	keyringUser    = "oauth-token"
	keyringTimeout = 2 * time.Second
)

var (
	errTokenNotFound      = errors.New("OAuth credentials not found")
	errKeyringUnavailable = errors.New("system credential storage is unavailable")
)

type tokenStore interface {
	Load() (*oauth2.Token, error)
	Save(*oauth2.Token) error
	Delete() error
}

type secureTokenStore struct {
	filePath string
}

func newTokenStore(dataDir string) tokenStore {
	return secureTokenStore{filePath: filepath.Join(dataDir, "token.json")}
}

func (s secureTokenStore) Load() (*oauth2.Token, error) {
	token, fileErr := loadTokenFile(s.filePath)
	if fileErr == nil {
		return token, nil
	}
	if !errors.Is(fileErr, errTokenNotFound) {
		return nil, fmt.Errorf("read fallback credential file: %w", fileErr)
	}

	serialized, err := keyringGet()
	if err == nil {
		return decodeToken([]byte(serialized))
	}
	return nil, errTokenNotFound
}

func (s secureTokenStore) Save(token *oauth2.Token) error {
	serialized, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("encode OAuth credentials: %w", err)
	}
	if err := keyringSet(string(serialized)); err == nil {
		return nil
	}
	return saveTokenFile(s.filePath, serialized)
}

func (s secureTokenStore) Delete() error {
	keyringErr := keyringDelete()
	fileErr := os.Remove(s.filePath)
	if errors.Is(fileErr, os.ErrNotExist) {
		fileErr = nil
	}
	if fileErr != nil {
		return fmt.Errorf("remove fallback credential file: %w", fileErr)
	}
	if keyringErr != nil && !errors.Is(keyringErr, keyring.ErrNotFound) && !errors.Is(keyringErr, errKeyringUnavailable) {
		return fmt.Errorf("delete system credential: %w", keyringErr)
	}
	return nil
}

func keyringGet() (string, error) {
	type result struct {
		value string
		err   error
	}
	results := make(chan result, 1)
	go func() {
		value, err := keyring.Get(keyringService, keyringUser)
		results <- result{value: value, err: err}
	}()
	select {
	case result := <-results:
		return result.value, result.err
	case <-time.After(keyringTimeout):
		return "", errKeyringUnavailable
	}
}

func keyringSet(value string) error {
	results := make(chan error, 1)
	go func() {
		results <- keyring.Set(keyringService, keyringUser, value)
	}()
	select {
	case err := <-results:
		return err
	case <-time.After(keyringTimeout):
		return errKeyringUnavailable
	}
}

func keyringDelete() error {
	results := make(chan error, 1)
	go func() {
		results <- keyring.Delete(keyringService, keyringUser)
	}()
	select {
	case err := <-results:
		return err
	case <-time.After(keyringTimeout):
		return errKeyringUnavailable
	}
}

func loadTokenFile(path string) (*oauth2.Token, error) {
	contents, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, errTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return decodeToken(contents)
}

func decodeToken(contents []byte) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal(contents, &token); err != nil {
		return nil, fmt.Errorf("decode OAuth credentials: %w", err)
	}
	return &token, nil
}

func saveTokenFile(path string, contents []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create credential directory: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".token-*")
	if err != nil {
		return err
	}
	tempPath := file.Name()
	defer os.Remove(tempPath)
	if err := file.Chmod(0o600); err != nil {
		file.Close()
		return err
	}
	if _, err := file.Write(contents); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}
