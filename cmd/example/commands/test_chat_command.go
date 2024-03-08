package commands

import (
	"log"

	"github.com/ynotnauk/go-twitch/entities"
)

type TestChatCommand struct{}

func (c *TestChatCommand) Execute(context *entities.ChatCommandContext) {
	log.Print("The test command has been called")
	context.Reply(context.Message, "Testing 1 2 3 4 5")
}
