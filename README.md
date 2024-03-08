# go-twitch
Complete twitch bot written in Go, this is a personal project to learn the Go language while trying to create something that will be useful to me.

**WARNING!:** This bot is in a very alpha state and I would not reconmend using this in production

## Example
There will be an example bot in the cmd/example folder which can be run. If you would like to run your own then the below boilerplate will get you started

```go
package main

import (
	"github.com/ynotnauk/go-twitch/auth"
	"github.com/ynotnauk/go-twitch/bot"
	"github.com/ynotnauk/go-twitch/entities"
	"github.com/ynotnauk/go-twitch/store"
)

func main() {
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
```

## In the wild
This bot will be running in the following channels

[YnotnaUK](https://twitch.tv/YnotnaUK "My twitch channel")

## Running tests

```
go test ./... -cover
```
