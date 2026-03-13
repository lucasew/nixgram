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

// Runner holds the execution state for a single user command.
// It tracks the parsed command, arguments, and the ID of the requesting user.
type Runner struct {
    bot *NixGram
    command string
    args []string
    sender int
}

// NewRunner parses a Telegram text message into an executable command
// and arguments, associating them with a specific user session.
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

// Run executes the command in the foreground.
// Commands are dynamically prefixed with "nixgram-" and looked up in $PATH.
// Depending on output size, the result is sent as a Telegram message or attached as a file.
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
    if r.sendMessage(out.String()) != nil {
        r.sendTextFile(out)
    }
    if err != nil {
        r.sendMessage(fmt.Sprintf("Error: %s", err))
        return err
    }
    return nil
}

// PocSplitter tokenizes a text command string into a command and slice of arguments.
// It removes trailing/leading slashes and spaces before splitting by spaces.
//
// TODO: Write a better splitter
func PocSplitter(text string) ([]string, error) {
    trimmed := strings.Trim(text, " /")
    return strings.Split(trimmed, " "), nil
}
