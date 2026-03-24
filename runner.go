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

// Runner encapsulates a single system command execution triggered by a Telegram
// message. It manages input arguments, parses execution paths, and handles sending
// stdout/stderr results back to the bot administrator.
type Runner struct {
	bot     *NixGram
	command string
	args    []string
	sender  int
}

// NewRunner parses raw Telegram message text to construct a new execution Runner.
// It uses PocSplitter to separate the command name from its arguments.
// The resulting Runner assumes the command targets the system $PATH.
func NewRunner(bot *NixGram, message string, sender int) (*Runner, error) {
	params, err := PocSplitter(message)
	return &Runner{
		bot:     bot,
		command: params[0],
		args:    params[1:],
		sender:  sender,
	}, err
}

// getCommand resolves the raw command string against the system $PATH.
// Crucially, it prefixes all requested commands with "nixgram-" to ensure
// isolated execution and prevent arbitrary binaries (like /bin/rm) from
// being invoked without an explicit nixgram wrapper.
func (r *Runner) getCommand() (string, bool) {
	cmdname := fmt.Sprintf("nixgram-%s", r.command)
	fullCmd, err := exec.LookPath(cmdname)
	if err != nil {
		return "", false
	}
	return fullCmd, true
}

// sendMessage delivers a plain-text response directly to the chat associated
// with the command sender. Useful for short stdout responses or status updates.
func (r *Runner) sendMessage(msg string) error {
	_, err := r.bot.Bot.Send(tgbotapi.NewMessage(int64(r.sender), msg))
	return err
}

// sendTextFile uploads execution output as a distinct text file attachment.
// Used as a fallback when standard bot messages exceed Telegram's length limits.
func (r *Runner) sendTextFile(b *bytes.Buffer) error {
	_, err := r.bot.Bot.Send(tgbotapi.NewDocumentUpload(int64(r.sender), tgbotapi.FileReader{
		Name:   "out.txt",
		Reader: b,
		Size:   int64(b.Len()),
	}))
	return err
}

// handleCommand spawns the resolved system binary with the provided arguments,
// piping both stdout and stderr into the provided io.Writer for capture.
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

// Run performs the end-to-end command sequence: verifies the executable path,
// invokes the shell command, and flushes output buffers back to the user chat.
// Automatically falls back to file uploads if the text payload is too large.
func (r *Runner) Run(ctx context.Context) error {
	log.Printf("Command %d: %s [ %s ]", r.sender, r.command, strings.Join(r.args, ", "))
	_, ok := r.getCommand()
	if !ok {
		err := fmt.Errorf("comando %s não encontrado", r.command)
		_ = r.sendMessage(err.Error())
		return err
	}
	out := bytes.NewBuffer([]byte{})
	_ = r.sendMessage("Running...")
	err := r.handleCommand(ctx, out)
	if r.sendMessage(out.String()) != nil {
		_ = r.sendTextFile(out)
	}
	if err != nil {
		_ = r.sendMessage(fmt.Sprintf("Error: %s", err))
		return err
	}
	return nil
}

// PocSplitter breaks the raw Telegram message string into an executable name
// and arguments by performing a naive space-delimited split.
// TODO: Implement proper shell-like tokenization to support quoted arguments.
func PocSplitter(text string) ([]string, error) {
	trimmed := strings.Trim(text, " /")
	return strings.Split(trimmed, " "), nil
}
