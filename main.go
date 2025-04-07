package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var admin tele.ChatID

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}
}

func main() {
	t, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is empty")
		return
	}

	id, ok := os.LookupEnv("ADMIN_ID")
	if !ok {
		return
	}
	adminID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatalf("strconv.ParseInt: %s", err)
		return
	}

	pref := tele.Settings{
		Token: t,
		Poller: &tele.LongPoller{
			Timeout: time.Second,
			AllowedUpdates: []string{
				"message",
				"edit_message", // todo implement
			},
		},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}
	admin = tele.ChatID(adminID)

	bot.Handle("/start", onStart)
	bot.Handle(tele.OnText, onMessage)

	bot.Start()

}

func onStart(c tele.Context) error {
	return c.Send("Добрый день! Напишите все, что считаете нужным. Все, что вы напишете, я бережно передам администратору 🤖")
}

func onMessage(c tele.Context) error {
	if err := c.Bot().React(
		c.Chat(),
		c.Message(),
		tele.ReactionOptions{
			Reactions: []tele.Reaction{{Emoji: "👌", Type: "emoji"}}},
	); err != nil {
		log.Printf("tele.Reaction: %s", err)
	}

	if c.Sender().Recipient() == admin.Recipient() {
		return nil
	}

	return c.ForwardTo(admin)
}
