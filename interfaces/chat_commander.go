package interfaces

import "github.com/ynotnauk/go-twitch/entities"

type ChatCommander interface {
	Execute(message *entities.ChatCommandContext)
}
