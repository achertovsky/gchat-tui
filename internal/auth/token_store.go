package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	keyringService = "gchat-tui"
	keyringUser    = "oauth-token"
)

var errTokenNotFound = errors.New("OAuth credentials not found")

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
	serialized, err := keyring.Get(keyringService, keyringUser)
	if err == nil {
		return decodeToken([]byte(serialized))
	}

	token, fileErr := loadTokenFile(s.filePath)
	if fileErr == nil {
		return token, nil
	}
	if errors.Is(fileErr, errTokenNotFound) {
		return nil, errTokenNotFound
	}
	return nil, fmt.Errorf("read credential storage: %w", fileErr)
}

func (s secureTokenStore) Save(token *oauth2.Token) error {
	serialized, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("encode OAuth credentials: %w", err)
	}
	if err := keyring.Set(keyringService, keyringUser, string(serialized)); err == nil {
		return nil
	}
	return saveTokenFile(s.filePath, serialized)
}

func (s secureTokenStore) Delete() error {
	keyringErr := keyring.Delete(keyringService, keyringUser)
	fileErr := os.Remove(s.filePath)
	if errors.Is(fileErr, os.ErrNotExist) {
		fileErr = nil
	}
	if fileErr != nil {
		return fmt.Errorf("remove fallback credential file: %w", fileErr)
	}
	if keyringErr != nil && !errors.Is(keyringErr, keyring.ErrNotFound) {
		return fmt.Errorf("delete system credential: %w", keyringErr)
	}
	return nil
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
