package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ynotnauk/go-twitch/entities"
)

const (
	tokenRefreshEndpoint    string = "https://id.twitch.tv/oauth2/token"
	tokenValidationEndpoint string = "https://id.twitch.tv/oauth2/validate"
	userAgent               string = "TwitchBot v1.0"
)

var (
	ErrBlankAccessToken  error = errors.New("accessToken cannot be blank")
	ErrBlankClientId     error = errors.New("clientId cannot be blank")
	ErrBlankClientSecret error = errors.New("clientSecret cannot be blank")
	ErrBlankRefreshToken error = errors.New("refreshToken cannot be blank")
	ErrNilAuthStore      error = errors.New("authStore cannot be nil")
	ErrBlankUserId       error = errors.New("userId cannot be blank")
)

type RefreshingAuthProvider struct {
	authStore  Storer
	httpClient *http.Client
	userId     string
}

func (a *RefreshingAuthProvider) GetAccessToken() (string, error) {
	authRecord, err := a.getAuthRecord()
	if err != nil {
		return "", err
	}
	return authRecord.AccessToken, nil
}

func (a *RefreshingAuthProvider) GetLoginAndAccessToken() (string, string, error) {
	authRecord, err := a.getAuthRecord()
	if err != nil {
		return "", "", err
	}
	return authRecord.Login, authRecord.AccessToken, nil
}

func (a *RefreshingAuthProvider) getAuthRecord() (*entities.AuthRecord, error) {
	// Get current record
	currentAuthRecord, err := a.authStore.GetByUserId(a.userId)
	if err != nil {
		return nil, err
	}
	// Ensure that the current access token is valid
	_, err = a.validateAccessToken(currentAuthRecord.AccessToken)
	if err != nil {
		// Get new tokens
		refreshTokenSuccess, err := a.refreshAccessToken(
			currentAuthRecord.ClientId,
			currentAuthRecord.ClientSecret,
			currentAuthRecord.RefreshToken,
		)
		if err != nil {
			return nil, err
		}
		// Validate the new tokens, this is so we can get the login name and id
		validateTokenSuccess, err := a.validateAccessToken(refreshTokenSuccess.AccessToken)
		if err != nil {
			return nil, err
		}
		// Create new auth record
		newAuthRecord := &entities.AuthRecord{
			AccessToken:  refreshTokenSuccess.AccessToken,
			ClientId:     validateTokenSuccess.ClientId,
			ClientSecret: currentAuthRecord.ClientSecret,
			ExpiresIn:    refreshTokenSuccess.ExpiresIn,
			Login:        validateTokenSuccess.Login,
			RefreshToken: refreshTokenSuccess.RefreshToken,
			Scope:        refreshTokenSuccess.Scope,
			TokenType:    refreshTokenSuccess.TokenType,
			UserId:       validateTokenSuccess.UserId,
		}
		// Save new auth record to store
		err = a.authStore.UpdateByUserId(newAuthRecord)
		if err != nil {
			return nil, err
		}
		// Return updated record
		return newAuthRecord, nil
	} else {
		// Return current record
		return currentAuthRecord, nil
	}
}

func (a *RefreshingAuthProvider) refreshAccessToken(
	clientId string,
	clientSecret string,
	refreshToken string,
) (*RefreshTokenSuccess, error) {
	if clientId == "" {
		return nil, ErrBlankClientId
	}
	if clientSecret == "" {
		return nil, ErrBlankClientSecret
	}
	if refreshToken == "" {
		return nil, ErrBlankRefreshToken
	}
	// Create body to send to the server
	formBody := []byte(fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s",
		refreshToken,
		clientId,
		clientSecret,
	))
	// Create body reader
	bodyReader := bytes.NewReader(formBody)
	// request a new request
	request, err := http.NewRequest(http.MethodPost, tokenRefreshEndpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	// Set request headers
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", userAgent)
	// Process request
	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// Check if refresh was successful
	if response.StatusCode != 200 {
		refreshTokenFailed := &RefreshTokenFailed{}
		err := json.NewDecoder(response.Body).Decode(refreshTokenFailed)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(refreshTokenFailed.Message)
	} else {
		refreshTokenSuccess := &RefreshTokenSuccess{}
		err := json.NewDecoder(response.Body).Decode(refreshTokenSuccess)
		if err != nil {
			return nil, err
		}
		return refreshTokenSuccess, nil
	}
}

func (a *RefreshingAuthProvider) validateAccessToken(accessToken string) (*ValidateTokenSuccess, error) {
	if accessToken == "" {
		return nil, ErrBlankAccessToken
	}
	// Create request
	request, err := http.NewRequest(http.MethodGet, tokenValidationEndpoint, nil)
	if err != nil {
		return nil, err
	}
	// Set request headers
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("User-Agent", userAgent)
	// Send request
	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// If the response is not 200 then the token is invalid
	if response.StatusCode != 200 {
		validateTokenFailed := &ValidateTokenFailed{}
		err := json.NewDecoder(response.Body).Decode(validateTokenFailed)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(validateTokenFailed.Message)
	}
	// Decode response
	validateTokenSuccess := &ValidateTokenSuccess{}
	err = json.NewDecoder(response.Body).Decode(validateTokenSuccess)
	if err != nil {
		return nil, err
	}
	return validateTokenSuccess, nil
}

func NewRefreshingProvider(authStore Storer, userId string) (*RefreshingAuthProvider, error) {
	if authStore == nil {
		return nil, ErrNilAuthStore
	}
	if userId == "" {
		return nil, ErrBlankUserId
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	provider := &RefreshingAuthProvider{
		authStore:  authStore,
		httpClient: httpClient,
		userId:     userId,
	}
	return provider, nil
}
