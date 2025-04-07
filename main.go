package main

import (
	"fmt"
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
	admin = tele.ChatID(adminID)

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

	bot.Handle("/start", onStart)
	bot.Handle(tele.OnText, onMessage)

	bot.Start()

}

func onStart(c tele.Context) error {
	return c.Send("Добрый день! Напишите все, что считаете нужным. Все, что вы напишете, я бережно передам администратору 🤖")
}

func onMessage(c tele.Context) error {
	_, err := c.Bot().Forward(admin, c.Message())
	if err != nil {
		return fmt.Errorf("forward: %s", err)
	}
	return c.Bot().React(
		c.Chat(),
		c.Message(),
		tele.ReactionOptions{
			Reactions: []tele.Reaction{{Emoji: "👌", Type: "emoji"}}},
	)
}
