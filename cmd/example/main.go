package main

import (
	"log"
	"time"

	"github.com/ynotnauk/go-twitch"
)

func main() {
	log.SetPrefix("[Twitch Bot Example] ")
	// Create auth store
	authStore, err := twitch.NewAuthFilesystemStore("data")
	if err != nil {
		panic(err)
	}
	// Create auth provider
	authProvider, err := twitch.NewRefreshingProvider(authStore, "142216347")
	if err != nil {
		panic(err)
	}
	// Create complete bot
	bot, err := twitch.NewBot(authProvider)
	if err != nil {
		panic(err)
	}
	// Add Chat Commands
	bot.OnChatCommand("test", &HelloChatCommand{})
	// Handlers
	bot.OnChatConnect(func(message *twitch.ChatConnectMessage) {
		bot.ChatJoin("ynotnauk")
	})
	bot.OnChatJoin(func(message *twitch.ChatJoinMessage) {
		log.Printf("[%s] %s has joined the channel",
			message.Channel,
			message.Username,
		)
	})
	bot.OnChatPart(func(message *twitch.ChatPartMessage) {
		log.Printf("[%s] %s has left the channel",
			message.Channel,
			message.Username,
		)
	})
	bot.OnChatPong(func(message *twitch.ChatPongMessage) {
		now := time.Now().Unix()
		latency := now - message.Timestamp
		log.Printf("Current Latency: %v ms", latency)
	})
	bot.OnChatPrivateMessage(func(message *twitch.ChatPrivateMessage) {
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
