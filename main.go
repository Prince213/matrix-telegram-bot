package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var c *mautrix.Client

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(os.Getenv("BOT_TELEGRAM_TOKEN"), bot.WithDefaultHandler(handler))
	if err != nil {
		log.Fatal(err)
	}

	c, err = mautrix.NewClient("matrix.org", "", "")
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Login(ctx, &mautrix.ReqLogin{
		Type:             mautrix.AuthTypePassword,
		DeviceID:         id.DeviceID(os.Getenv("BOT_MATRIX_DEVICEID")),
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: os.Getenv("BOT_MATRIX_USERID")},
		Password:         os.Getenv("BOT_MATRIX_PASSWORD"),
		StoreCredentials: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	b.Start(ctx)

	_, err = c.Logout(ctx)
	if err != nil {
		log.Fatal(err)
	}
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

	resp, err = http.Post("https://pb.nichi.co/", resp.Header.Get("Content-Type"), resp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	url, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}

	_, err = c.SendText(ctx, id.RoomID(os.Getenv("BOT_MATRIX_TARGET_ROOM")), string(url))
	if err != nil {
		log.Print(err)
	}
}
