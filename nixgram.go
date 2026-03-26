package nixgram

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// NixGram orchestrates the core bot lifecycle. It holds the active Telegram
// connection, maintains the updates channel for inbound messages, and enforces
// security by isolating commands to a single authorized administrator ID.
type NixGram struct {
	Bot         *tgbotapi.BotAPI
	updatesChan tgbotapi.UpdatesChannel
	Adm         int
}

// NewNixGram initializes the Telegram bot client with the provided token and
// retrieves the asynchronous updates channel for polling new messages.
// Returns an error if the token is invalid or if the Telegram API is unreachable.
func NewNixGram(token string, adm int) (*NixGram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	updatesChan, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Limit:   1,
		Timeout: 10,
	})
	if err != nil {
		return nil, err
	}
	return &NixGram{
		Bot:         bot,
		Adm:         adm,
		updatesChan: updatesChan,
	}, nil
}

// Run starts the main polling loop, blocking until the provided context is canceled.
// It continuously consumes events from the updates channel and spawns command execution
// requests via handleMessage.
func (n *NixGram) Run(ctx context.Context) {
	for {
		select {
		case u := <-n.updatesChan:
			n.handleMessage(ctx, u)
		case <-ctx.Done():
			return
		}
	}
}

// handleMessage inspects inbound Telegram updates, immediately discarding anything
// not from the configured administrator ID. For authorized commands, it delegates
// execution to a new Runner instance on a separate goroutine.
func (n *NixGram) handleMessage(ctx context.Context, u tgbotapi.Update) {
	if u.Message == nil {
		return
	}
	if u.Message.From.ID != n.Adm {
		log.Printf("WARN: usuário não autorizado tentou usar o bot: %s (%d): %s", u.Message.From.UserName, u.Message.From.ID, u.Message.Text)
		return
	}
	r, err := NewRunner(n, u.Message.Text, u.Message.From.ID)
	if err != nil {
		return
	}
	go func() {
		err = r.Run(ctx)
		if err != nil {
			log.Printf("ERRO: (%d) %s", u.Message.From.ID, err.Error())
		}
	}()
}
