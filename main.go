package main

import (
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
)

var adminUser tele.Recipient
var bot *tele.Bot

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}

	botToken := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:     botToken,
		ParseMode: tele.ModeMarkdown,
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}
	adminUser, err = bot.ChatByUsername(os.Getenv("ADMIN_USER"))
	if err != nil {
		log.Fatalf("admin user is not found: %s", err)
		return
	}
}

func main() {
	bot.Handle(tele.OnText, handler)
	//todo use в другом проекте
	bot.Handle(tele.OnUserJoined, joinHandler)

	bot.Start()
}

func handler(c tele.Context) error {
	if isAdmin(c) {
		return sendResponse(c)
	}
	return sendRequest(c)

}

func joinHandler(c tele.Context) error {
	return c.Send("Welcome!")

}

func isAdmin(c tele.Context) bool {
	return adminUser.Recipient() == c.Sender().Recipient()
}

func sendResponse(c tele.Context) error {
	if c.Message().ReplyTo == nil {
		return c.Send("Выберите сообщение, на которое хотите ответить")
	}
	sender := c.Message().ReplyTo.OriginalSender
	_, err := bot.Send(sender, c.Text())
	return err
}

func sendRequest(c tele.Context) error {
	return c.ForwardTo(adminUser)
}
