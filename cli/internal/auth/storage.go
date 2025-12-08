package auth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	ServiceName = "luna-cli"
)

// TokenStore defines the interface for token storage
type TokenStore interface {
	Save(accessToken, refreshToken string) error
	Load(account string) (string, string, error)
	Clear(account string) error
}

// KeyringStore implements TokenStore using system keyring
type KeyringStore struct{}

func NewKeyringStore() *KeyringStore {
	return &KeyringStore{}
}

func (s *KeyringStore) Save(account, accessToken, refreshToken string) error {
	if err := keyring.Set(ServiceName, account+"_access", accessToken); err != nil {
		return err
	}
	if refreshToken != "" {
		if err := keyring.Set(ServiceName, account+"_refresh", refreshToken); err != nil {
			return err
		}
	}
	return nil
}

func (s *KeyringStore) Load(account string) (string, string, error) {
	accessToken, err := keyring.Get(ServiceName, account+"_access")
	if err != nil {
		return "", "", err
	}
	refreshToken, _ := keyring.Get(ServiceName, account+"_refresh")
	return accessToken, refreshToken, nil
}

func (s *KeyringStore) Clear(account string) error {
	_ = keyring.Delete(ServiceName, account+"_access")
	_ = keyring.Delete(ServiceName, account+"_refresh")
	return nil
}

// FileStore implements TokenStore using a local file (fallback)
type FileStore struct {
	Path string
}

func NewFileStore() *FileStore {
	home, _ := os.UserHomeDir()
	return &FileStore{
		Path: filepath.Join(home, ".luna", "credentials.json"),
	}
}

func (s *FileStore) Save(account, accessToken, refreshToken string) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0755); err != nil {
		return err
	}

	// Simple single-account implementation for file store
	data := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	bytes, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(s.Path, bytes, 0600)
}

func (s *FileStore) Load(account string) (string, string, error) {
	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return "", "", err
	}

	var data map[string]string
	if err := json.Unmarshal(bytes, &data); err != nil {
		return "", "", err
	}

	return data["access_token"], data["refresh_token"], nil
}

func (s *FileStore) Clear(account string) error {
	return os.Remove(s.Path)
}
