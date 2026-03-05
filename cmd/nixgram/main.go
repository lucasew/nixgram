package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/lucasew/nixgram"
	"github.com/lucasew/nixgram/pkg/errreporter"
)

type Config struct {
    Token string
    Adm   int
}

func loadEnvironment() (*Config, error) {
    token := os.Getenv("NIXGRAM_TOKEN")
    if (token == "") {
        return nil, fmt.Errorf("Missing NIXGRAM_TOKEN")
    }
    admStr := os.Getenv("NIXGRAM_ADM")
    if (admStr == "") {
        return nil, fmt.Errorf("Missing NIXGRAM_ADM")
    }
    adm, err := strconv.Atoi(admStr)
    if (err != nil) {
        return nil, fmt.Errorf("NIXGRAM_ADM: %s is not a number", admStr)
    }
    return &Config{Token: token, Adm: adm}, nil
}

func main() {
    config, err := loadEnvironment()
    if err != nil {
        errreporter.ReportError(err, "failed to load environment")
        os.Exit(1)
    }
    bot, err := nixgram.NewNixGram(config.Token, config.Adm)
    if err != nil {
        errreporter.ReportError(err, "failed to initialize bot")
        os.Exit(1)
    }
    bot.Run(context.Background())
}
