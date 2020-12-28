package nixgram

import (
	"bytes"
	"context"
	"fmt"
    "io"
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

func (r *Runner) getCommand() (string, bool) {
    cmdname := fmt.Sprintf("nixgram-%s", r.command)
    fullCmd, err := exec.LookPath(cmdname)
    if err != nil {
        return "", false
    }
    return fullCmd, true
}

func (r *Runner) sendMessage(msg string) error {
    _, err := r.bot.Bot.Send(tgbotapi.NewMessage(int64(r.sender), msg))
    return err
}

func (r *Runner) sendTextFile(b *bytes.Buffer) error {
    _, err := r.bot.Bot.Send(tgbotapi.NewDocumentUpload(int64(r.sender), tgbotapi.FileReader{
        Name: "out.txt",
        Reader: b,
        Size: int64(b.Len()),
    }))
    return err
}

func (r *Runner) handleCommand(ctx context.Context, b io.Writer) error {
    cmdPath, ok := r.getCommand()
    if !ok {
        return fmt.Errorf("comando %s não encontrado", r.command)
    }
    cmd := exec.CommandContext(ctx, cmdPath, r.args...)
    cmd.Stdout = b
    cmd.Stderr = b
    return cmd.Run()
}

func (r *Runner) Run(ctx context.Context) error {
    log.Printf("Command %d: %s [ %s ]", r.sender, r.command, strings.Join(r.args, ", "))
    _, ok := r.getCommand()
    if !ok {
        err := fmt.Errorf("comando %s não encontrado", r.command)
        r.sendMessage(err.Error())
        return err
    }
    out := bytes.NewBuffer([]byte{})
    r.sendMessage("Running...")
    err := r.handleCommand(ctx, out)
    if err != nil {
        r.sendMessage(fmt.Sprintf("Error: %s", err))
        return err
    }
    if r.sendMessage(out.String()) != nil {
        return r.sendTextFile(out)
    }
    return nil
}

//TODO: Write a better splitter
func PocSplitter(text string) ([]string, error) {
    trimmed := strings.Trim(text, " /")
    return strings.Split(trimmed, " "), nil
}
