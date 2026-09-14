package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestFallbackTokenFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials", "token.json")
	want := &oauth2.Token{
		AccessToken:  "access",
		RefreshToken: "refresh",
		Expiry:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	serialized, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if err := saveTokenFile(path, serialized); err != nil {
		t.Fatalf("save token: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("permissions = %o, want 600", info.Mode().Perm())
	}
	got, err := loadTokenFile(path)
	if err != nil {
		t.Fatalf("load token: %v", err)
	}
	if got.AccessToken != want.AccessToken || got.RefreshToken != want.RefreshToken || !got.Expiry.Equal(want.Expiry) {
		t.Error("loaded token does not match saved token")
	}
}

func TestLoadTokenFileReportsMissingCredentials(t *testing.T) {
	_, err := loadTokenFile(filepath.Join(t.TempDir(), "missing.json"))
	if !errors.Is(err, errTokenNotFound) {
		t.Fatalf("error = %v, want missing credentials", err)
	}
}
