package main

import (
	"log"

	"github.com/ynotnauk/go-twitch/auth"
	"github.com/ynotnauk/go-twitch/bot"
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
	// Handlers
	bot.OnTwitchChatConnect(func(message *entities.ChatConnectMessage) {
		bot.ChatJoin("ynotnauk")
	})
	// Start bot
	bot.Start()
}
