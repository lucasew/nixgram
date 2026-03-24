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

// Runner encapsulates the execution context for a single command triggered by a user.
// It resolves the requested command against a restricted 'nixgram-' prefix in the PATH,
// executes it safely, and streams the standard output/error back to the user via Telegram.
type Runner struct {
    bot *NixGram
    command string
    args []string
    sender int
}

// NewRunner constructs a Runner instance by parsing the raw incoming Telegram message.
// The message is split into a base command and its arguments. The command must map to an
// executable prefixed with 'nixgram-' available in the system PATH.
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

// Run executes the resolved command within the provided context.
// It captures the combined stdout and stderr of the process. If the output fits within
// Telegram's message limits, it sends it directly; otherwise, it uploads the output as a text file.
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

// PocSplitter splits a raw command string into a slice of arguments.
// It handles basic space separation and trims leading/trailing whitespace and slashes.
// TODO: Implement a more robust parsing strategy (e.g., handling quoted arguments).
func PocSplitter(text string) ([]string, error) {
    trimmed := strings.Trim(text, " /")
    return strings.Split(trimmed, " "), nil
}
