package main 

import (
	"math/rand"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("594859196:AAFYQZiJ1dDlK_rsAvpLeWNargmSM1ki0Ag")
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	jokes := []string{"Assalamualaikum, \nyang nggak jawab PKI", "I don’t expect you to read this since comments are ignored by the compiler."}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		Message := update.Message.Text

		if Message == "/admin" {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@kokizzu @Dim_As_YP ")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

		}

		if Message == "/joke" {
			n := rand.Int() % len(jokes)
			
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, jokes[n])
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)

		}
	}
}