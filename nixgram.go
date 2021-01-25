package nixgram

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type NixGram struct {
    Bot *tgbotapi.BotAPI
    updatesChan tgbotapi.UpdatesChannel
    Adm int
}

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
