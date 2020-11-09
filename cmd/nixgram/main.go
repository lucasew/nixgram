package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/lucasew/nixgram"
)

var token string
var adm int
var err error

func loadEnvironment() error {
    token = os.Getenv("NIXGRAM_TOKEN")
    if (token == "") {
        return fmt.Errorf("Missing NIXGRAM_TOKEN")
    }
    admStr := os.Getenv("NIXGRAM_ADM")
    if (admStr == "") {
        return fmt.Errorf("Missing NIXGRAM_ADM")
    }
    adm, err = strconv.Atoi(admStr)
    if (err != nil) {
        return fmt.Errorf("NIXGRAM_ADM: %s is not a number", admStr)
    }
    return nil
}

func main() {
    err := loadEnvironment()
    if err != nil {
        panic(err)
    }
    bot, err := nixgram.NewNixGram(token, adm)
    if err != nil {
        panic(err)
    }
    bot.Run(context.Background())
}
