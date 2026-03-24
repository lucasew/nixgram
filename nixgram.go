package nixgram

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// NixGram orchestrates the Telegram bot connection and message routing.
// It maintains the API client, an updates channel for incoming messages,
// and enforces authorization by restricting access to a single configured administrator.
type NixGram struct {
    Bot *tgbotapi.BotAPI
    updatesChan tgbotapi.UpdatesChannel
    Adm int
}

// NewNixGram initializes a new bot instance with the provided Telegram token.
// It also establishes the updates channel configuration and binds the administrator ID.
// Returns an error if the token is invalid or if the API connection fails.
func NewNixGram(token string, adm int) (*NixGram, error) {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }
    updatesChan, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
        Limit: 1,
        Timeout: 10,
    })
    if err != nil {
        return nil, err
    }
    return &NixGram{
        Bot: bot,
        Adm: adm,
        updatesChan: updatesChan,
    }, nil
}

// Run starts the main event loop, consuming messages from the updates channel.
// It blocks until the provided context is canceled, dispatching each valid update
// to the message handler for execution.
func (n* NixGram) Run(ctx context.Context) {
    for {
        select {
            case u := <- n.updatesChan:
                n.handleMessage(ctx, u)
            case <- ctx.Done():
                return
        }
    }
}

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
        if (err != nil) {
            log.Printf("ERRO: (%d) %s", u.Message.From.ID, err.Error())
        }
    }()
}
