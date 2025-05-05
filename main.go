package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(os.Getenv("BOT_TELEGRAM_TOKEN"), bot.WithDefaultHandler(handler))
	if err != nil {
		log.Fatal(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	m := update.Message
	if m == nil || len(m.Photo) == 0 {
		return
	}

	// TODO: which one to download
	photo := m.Photo[len(m.Photo)-1]

	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: photo.FileID})
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := http.Get(b.FileDownloadLink(file))
	if err != nil {
		log.Print(err)
		return
	}

	_ = resp
}
