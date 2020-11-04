package nixgram

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Runner struct {
    bot *NixGram
    command string
    args []string
    sender int
}

func NewRunner(bot *NixGram, message string, sender int) (*Runner, error) {
    params, err := PocSplitter(message)
    return &Runner{
        bot: bot,
        command: params[0],
        args: params[1:],
        sender: sender,
    }, err
}

func (r *Runner) Run(ctx context.Context) error {
    cmdname := fmt.Sprintf("nixgram-%s", r.command)
    log.Printf("Command %d: %s [ %s ]", r.sender, cmdname, strings.Join(r.args, ", "))
    fullCmdName, err := exec.LookPath(cmdname)
    if err != nil {
        r.bot.Bot.Send(tgbotapi.NewMessage(int64(r.sender), fmt.Sprintf("Comando não encontrado: %s", r.command)))
        return err
    }
    out := bytes.NewBuffer([]byte{})
    cmd := exec.CommandContext(ctx, fullCmdName, r.args...)
    cmd.Stdout = out
    cmd.Stderr = out
    err = cmd.Run()
    if err != nil {
        return err
    }
    _, err = r.bot.Bot.Send(tgbotapi.NewMessage(int64(r.sender), out.String()))
    if err != nil {
        _, err = r.bot.Bot.Send(tgbotapi.NewDocumentUpload(int64(r.sender), tgbotapi.FileReader{
            Name: "out.txt",
            Reader: out,
            Size: int64(out.Len()),
        }))
        if err != nil {
            return err
        }
    }
    return nil
}

func PocSplitter(text string) ([]string, error) {
    trimmed := strings.Trim(text, " /")
    return strings.Split(trimmed, " "), nil
}
