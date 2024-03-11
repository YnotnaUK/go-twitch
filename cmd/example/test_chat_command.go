package main

import (
	"github.com/ynotnauk/go-twitch"
)

type HelloChatCommand struct{}

func (c *HelloChatCommand) Execute(context *twitch.ChatCommandContext) {
	context.Reply(context.Message, "Hello!")
}
