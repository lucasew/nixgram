package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/lucasew/nixgram"
)

func main() {
    token := os.Getenv("NIXGRAM_TOKEN")
    if (token == "") {
        panic("Missing NIXGRAM_TOKEN")
    }
    adm := os.Getenv("NIXGRAM_ADM")
    if (adm == "") {
        panic("Missing NIXGRAM_ADM")
    }
    admNum, err := strconv.Atoi(adm)
    if (err != nil) {
        panic(fmt.Sprintf("NIXGRAM_ADM: %s is not a number", adm))
    }
    bot, err := nixgram.NewNixGram(token, admNum)
    if err != nil {
        panic(err)
    }
    bot.Run(context.Background())
}
