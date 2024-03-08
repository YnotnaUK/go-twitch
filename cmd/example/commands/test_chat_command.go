package commands

import (
	"log"

	"github.com/ynotnauk/go-twitch/entities"
)

type TestChatCommand struct{}

func (c *TestChatCommand) Execute(context *entities.ChatCommandContext) {
	log.Print("The test command has been called")
}
