package bot

import (
	"log"

	"github.com/ynotnauk/go-twitch/chat"
	"github.com/ynotnauk/go-twitch/entities"
	"github.com/ynotnauk/go-twitch/interfaces"
)

type Bot struct {
	chat *chat.Client
}

func (b *Bot) ChatJoin(channel string) error {
	err := b.chat.Join(channel)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) OnChatPrivateMessage(handler func(*entities.ChatPrivateMessage)) error {
	b.chat.OnPrivateMessage(handler)
	return nil
}

func (b *Bot) OnTwitchChatConnect(handler func(*entities.ChatConnectMessage)) error {
	b.chat.OnConnect(handler)
	return nil
}

func (b *Bot) Start() {
	log.Println("Starting bot...")
	b.chat.Start()
}

func New(authProvider interfaces.AuthProvider) (*Bot, error) {
	// Create chat client
	chat, err := chat.NewClient(authProvider)
	if err != nil {
		return nil, err
	}
	// Create bot
	bot := &Bot{
		chat: chat,
	}
	return bot, nil
}
