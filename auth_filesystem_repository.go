package twitch

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrFileLocationBlank = errors.New("fileLocation cannot be blank")
)

type AuthFilesystemRepository struct {
	fileLocation string
}

func (r *AuthFilesystemRepository) Load() (*Auth, error) {
	fileContents, err := os.ReadFile(r.fileLocation)
	if err != nil {
		return nil, err
	}
	auth := &Auth{}
	err = json.Unmarshal(fileContents, auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (r *AuthFilesystemRepository) Save(auth *Auth) error {
	fileContents, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(r.fileLocation, fileContents, 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewAuthFilesystemRepository(fileLocation string) (*AuthFilesystemRepository, error) {
	if fileLocation == "" {
		return nil, ErrFileLocationBlank
	}
	repository := &AuthFilesystemRepository{
		fileLocation: fileLocation,
	}
	return repository, nil
}
