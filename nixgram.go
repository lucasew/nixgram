package nixgram

import (
	"context"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lucasew/nixgram/pkg/errreporter"
)

type NixGram struct {
	Bot         *tgbotapi.BotAPI
	updatesChan tgbotapi.UpdatesChannel
	Adm         int
}

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

func (n *NixGram) handleMessage(ctx context.Context, u tgbotapi.Update) {
	if u.Message == nil {
		return
	}
	if u.Message.From.ID != n.Adm {
		errreporter.ReportError(
			fmt.Errorf("usuário não autorizado tentou usar o bot: %s", u.Message.Text),
			map[string]interface{}{
				"username": u.Message.From.UserName,
				"user_id":  u.Message.From.ID,
			},
		)
		return
	}
	r, err := NewRunner(n, u.Message.Text, u.Message.From.ID)
	if err != nil {
		errreporter.ReportError(err, map[string]interface{}{"context": "NewRunner"})
		return
	}
	go func() {
		err = r.Run(ctx)
		if err != nil {
			errreporter.ReportError(err, map[string]interface{}{"user_id": u.Message.From.ID})
		}
	}()
}
