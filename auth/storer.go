package auth

import (
	"github.com/ynotnauk/go-twitch/entities"
)

type Storer interface {
	GetByUserId(userId string) (*entities.AuthRecord, error)
	UpdateByUserId(auth *entities.AuthRecord) error
}
