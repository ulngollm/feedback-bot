package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
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
	bot.Handle(tele.OnReply, onReply)

	bot.Start()

}

func onStart(c tele.Context) error {
	return c.Send("Добрый день! Напишите все, что считаете нужным. Все, что вы напишете, я бережно передам администратору 🤖")
}

func onMessage(c tele.Context) error {
	if err := c.Bot().React(
		c.Chat(),
		c.Message(),
		tele.Reactions{
			Reactions: []tele.Reaction{{Emoji: "👌", Type: "emoji"}}},
	); err != nil {
		log.Printf("tele.Reaction: %s", err)
	}

	if c.Sender().Recipient() == admin.Recipient() {
		return nil
	}

	return c.ForwardTo(admin)
}

func onReply(c tele.Context) error {
	if c.Message().ReplyTo == nil {
		return nil
	}
	if c.Message().ReplyTo.OriginalChat == c.Chat() {
		return nil
	}
	r := c.Message().ReplyTo
	msg := c.Message().Text

	// todo use ReplyTo (чтобы на той стороне было понятно, на какое сообщение ответили)
	// todo need originalMessageID . Если все сообщения будут писаться в бд - проблема решена
	// лог должен быть такой, чтобы по id of forwarded message в админ чате можно было получить id оригинального сообщения
	_, err := c.Bot().Send(r.OriginalSender, msg)
	if err != nil {
		return fmt.Errorf("onReply.Send: %s", err)
	}
	return c.Send("ответ отправлен пользователю")
}
