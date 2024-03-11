package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	ErrBlankStoreLocation error = errors.New("storeLocation cannot be blank")
)

type AuthFilesystemStore struct {
	storeLocation string
}

func (s *AuthFilesystemStore) GetByUserId(userId string) (*AuthRecord, error) {
	storeFilePath := fmt.Sprintf("%s/auth.%s.json", s.storeLocation, userId)
	// Read store
	fileContents, err := os.ReadFile(storeFilePath)
	if err != nil {
		return nil, err
	}
	log.Printf("loaded file: %s", storeFilePath)
	// Create struct for auth
	authRecord := &AuthRecord{}
	// Write store to struct
	err = json.Unmarshal(fileContents, authRecord)
	if err != nil {
		return nil, err
	}
	return authRecord, nil
}

func (s *AuthFilesystemStore) UpdateByUserId(auth *AuthRecord) error {
	// Convert struct into a JSON byte array
	fileContents, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return err
	}
	storeFilePath := fmt.Sprintf("%s/auth.%s.json", s.storeLocation, auth.UserId)
	// Write store
	err = os.WriteFile(storeFilePath, fileContents, 0644)
	if err != nil {
		return err
	}
	log.Printf("wrote file: %s", storeFilePath)
	return nil
}

func NewAuthFilesystemStore(storeLocation string) (*AuthFilesystemStore, error) {
	// Ensure store location is not blank
	if storeLocation == "" {
		return nil, ErrBlankStoreLocation
	}
	// Ensure the path does not have a trailing slash
	storeLocation = strings.TrimSuffix(storeLocation, "/")
	// Create store
	store := &AuthFilesystemStore{
		storeLocation: storeLocation,
	}
	// Return store
	return store, nil
}
