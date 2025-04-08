package main

import (
	"feedback-bot/storage"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

var admin tele.ChatID
var st *storage.Storage
var answerModeEnabled bool

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
	answerModeEnabled = os.Getenv("ENABLE_ANSWER_MODE") != ""

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
	st = storage.New()
	if st == nil {
		log.Fatalf("storage.New: %s", fmt.Errorf("storate is not initialized"))
	}

	bot.Handle("/start", onStart)
	bot.Handle(tele.OnText, onMessage)
	bot.Handle(tele.OnReply, onAdminAnswer, canAnswerOnReply)

	bot.Start()

}

func onStart(c tele.Context) error {
	return c.Send("Добрый день! Напишите все, что считаете нужным. Все, что вы напишете, я бережно передам администратору 🤖")
}

func onMessage(c tele.Context) error {
	//todo чтобы работал только в личном чате
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

	fwd, err := c.Bot().Forward(admin, c.Message())
	if err != nil {
		return fmt.Errorf("forward: %s", err)
	}
	msg := storage.Message{
		OriginalMessageID:  c.Message().ID,
		ForwardedMessageID: fwd.ID,
		ChatID:             c.Chat().ID,
		Text:               c.Message().Text,
		CreatedAt:          c.Message().Time(),
	}
	if err := st.SaveMessage(msg); err != nil {
		return fmt.Errorf("storage.SaveMessage: %s", err)
	}
	return nil

}

func onAdminAnswer(c tele.Context) error {
	if c.Message().ReplyTo == nil {
		return nil
	}
	r := c.Message().ReplyTo
	// skip if message is not forward
	if r.Origin == nil {
		return nil
	}
	//if user is hidden - OriginalChat and OriginalSender are empty
	//проверить, что это не ответ на сообщение из этого же чата
	if c.Chat().FirstName == r.OriginalSenderName {
		return nil
	}

	fb, err := st.GetMessageByForwardedID(r.ID)
	if err != nil {
		if err := c.Send("ошибка отправки ответа пользователю"); err != nil {
			return fmt.Errorf("onAdminAnswer.send: %w", err)
		}
		return fmt.Errorf("getMessageByForwardedID: %w", err)
	}
	var rc tele.Recipient
	if r.OriginalSender != nil {
		rc = r.OriginalSender
	} else if fb != nil {
		rc = tele.ChatID(fb.ChatID)
	} else {
		return c.Send("получатель скрыт")
	}

	var opts *tele.SendOptions
	if fb != nil {
		omsg := &tele.Message{ID: fb.OriginalMessageID}
		opts = &tele.SendOptions{ReplyTo: omsg}
	}
	answer := c.Message().Text

	_, err = c.Bot().Send(rc, answer, opts)
	if err != nil {
		return fmt.Errorf("onAdminAnswer.Send: %s", err)
	}
	// todo записывать ответ в историю тоже
	return c.Send("ответ отправлен пользователю")
}

func canAnswerOnReply(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if answerModeEnabled {
			return next(c)
		}
		return nil
	}
}
