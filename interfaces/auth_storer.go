package interfaces

import "github.com/ynotnauk/go-twitch/entities"

type AuthStorer interface {
	GetByUserId(userId string) (*entities.AuthRecord, error)
	UpdateByUserId(auth *entities.AuthRecord) error
}
