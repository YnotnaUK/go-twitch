package main

import (
	"log"

	"github.com/ynotnauk/go-twitch/auth"
	"github.com/ynotnauk/go-twitch/bot"
	"github.com/ynotnauk/go-twitch/cmd/example/commands"
	"github.com/ynotnauk/go-twitch/entities"
	"github.com/ynotnauk/go-twitch/store"
)

func main() {
	log.SetPrefix("[Twitch Bot Example] ")
	// Create auth store
	authStore, err := store.NewAuthFilesystemStore("data")
	if err != nil {
		panic(err)
	}
	// Create auth provider
	authProvider, err := auth.NewRefreshingProvider(authStore, "142216347")
	if err != nil {
		panic(err)
	}
	// Create complete bot
	bot, err := bot.New(authProvider)
	if err != nil {
		panic(err)
	}
	// Add Chat Commands
	bot.OnChatCommand("test", &commands.TestChatCommand{})
	// Handlers
	bot.OnChatJoin(func(message *entities.ChatJoinMessage) {
		log.Printf("[%s] %s has joined the channel",
			message.Channel,
			message.Username,
		)
	})
	bot.OnChatPart(func(message *entities.ChatPartMessage) {
		log.Printf("[%s] %s has left the channel",
			message.Channel,
			message.Username,
		)
	})
	bot.OnTwitchChatConnect(func(message *entities.ChatConnectMessage) {
		bot.ChatJoin("ynotnauk")
	})
	bot.OnChatPrivateMessage(func(message *entities.ChatPrivateMessage) {
		log.Printf(
			"[%s] <%s:%s> %s",
			message.Channel,
			message.Tags["user-id"],
			message.Tags["display-name"],
			message.Message,
		)
	})
	// Start bot
	bot.Start()
}
