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

func (s *AuthFilesystemStore) GetByUserId(userId string) (*entities.AuthRecord, error) {
	// Read store
	fileContents, err := os.ReadFile(s.storeLocation)
	if err != nil {
		return nil, err
	}
	// Create struct for auth
	authRecord := &entities.AuthRecord{}
	// Write store to struct
	err = json.Unmarshal(fileContents, authRecord)
	if err != nil {
		return nil, err
	}
	return authRecord, nil
}

func (s *AuthFilesystemStore) UpdateByUserId(auth *entities.AuthRecord) error {
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
	// Ensure store location is not blank
	if storeLocation == "" {
		return nil, ErrBlankStoreLocation
	}
	// TODO: ensure store location is a json file
	store := &AuthFilesystemStore{
		storeLocation: storeLocation,
	}
	return store, nil
}
