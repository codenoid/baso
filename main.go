package main

import (
	"os"
	"log"
	"io/ioutil"
	"math/rand"
	"gopkg.in/yaml.v2"
)

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	botConfig struct {
		Token string `yaml:"token"`
		Debug bool `yaml:"debug"`
		Timeout int `yaml:"timeout"`
	}
)

func (c *botConfig) getConf() *botConfig {

    yamlFile, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
        os.Exit(1)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
        os.Exit(1)
    }

    return c
}

func main() {

	var c botConfig
    c.getConf()

	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = c.Debug

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

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@admin")
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
