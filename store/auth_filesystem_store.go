package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ynotnauk/go-twitch/entities"
)

var (
	ErrBlankStoreLocation error = errors.New("storeLocation cannot be blank")
)

type AuthFilesystemStore struct {
	storeLocation string
}

func (s *AuthFilesystemStore) GetAccessToken() (string, error) {
	store, err := s.readStore()
	if err != nil {
		return "", err
	}
	return store.AccessToken, nil
}

func (s *AuthFilesystemStore) readStore() (*entities.Auth, error) {
	// Read store
	fileContents, err := os.ReadFile(s.storeLocation)
	if err != nil {
		return nil, err
	}
	// Create struct for auth
	auth := &entities.Auth{}
	// Write store to struct
	err = json.Unmarshal(fileContents, auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (s *AuthFilesystemStore) updateStore(auth *entities.Auth) error {
	// Convert struct into a JSON byte array
	fileContents, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return err
	}
	// Write store
	err = os.WriteFile(s.storeLocation, fileContents, 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewAuthFilesystemStore(storeLocation string) (*AuthFilesystemStore, error) {
	if storeLocation == "" {
		return nil, ErrBlankStoreLocation
	}
	store := &AuthFilesystemStore{
		storeLocation: storeLocation,
	}
	return store, nil
}
