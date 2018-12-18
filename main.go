package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"encoding/json"
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

	JadwalSholat struct {
		JadwalData `json:"data"`
	}

	JadwalData struct {
		Subuh string `json:"Fajr"`
		Dzuhur string `json:"Dhuhr"`
		Ashar string `json:"Asr"`
		Maghrib string `json:"Maghrib"`
		Isha string `json:"Isha"`
	}
)

func currentDir() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    return dir
}

func InBetween(i, min, max int) bool { return (i >= min) && (i <= max) }

func (c *botConfig) getConf() *botConfig {
    yamlFile, err := ioutil.ReadFile(currentDir() + "/config.yml") // bad
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
	u.Timeout = c.Timeout

	updates, err := bot.GetUpdatesChan(u)

	content, err := ioutil.ReadFile(currentDir() + "/jokes.txt")
	if err != nil {
	    panic(err)
	}
	jokes := strings.Split(string(content), "\n")

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

		if Message == "/next-sholat" {
			resp, err := http.Get("https://time.siswadi.com/pray/?address=Jakarta")
			if err == nil {

			    body, err := ioutil.ReadAll(resp.Body)
			    if err != nil {
			        panic(err)
			    }

				var j JadwalSholat
				Err := json.Unmarshal(body, &j)

				if Err == nil {

					OniChannnn := "Jadwal Sholat tak Ditemukan :O"

					Hour := time.Now().Hour()

					if InBetween(Hour, 21, 6) {
						OniChannnn = "Kurang Lebih Anda akan Sholat Subuh pada Jam : " + j.JadwalData.Subuh
					}

					if InBetween(Hour, 6, 13) {
						OniChannnn = "Kurang Lebih Anda akan Sholat Dzuhur pada Jam : " + j.JadwalData.Dzuhur
					}

					if InBetween(Hour, 13, 16) {
						OniChannnn = "Kurang Lebih Anda akan Sholat Ashar pada Jam : " + j.JadwalData.Ashar
					}

					if InBetween(Hour, 16, 18) {
						OniChannnn = "Kurang Lebih Anda akan Sholat Maghrib pada Jam : " + j.JadwalData.Maghrib
					}

					if InBetween(Hour, 18, 21) {
						OniChannnn = "Kurang Lebih Anda akan Sholat Isha pada Jam : " + j.JadwalData.Isha
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, OniChannnn)
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}

			}
		}

		if Message == "/sholat" {
			resp, err := http.Get("https://time.siswadi.com/pray/?address=Jakarta")
			if err == nil {

			    body, err := ioutil.ReadAll(resp.Body)
			    if err != nil {
			        panic(err)
			    }

				var j JadwalSholat
				Err := json.Unmarshal(body, &j)

				if Err == nil {

					OniChannnn := fmt.Sprintf(`   Jadwal Sholat Untuk DKI Jakarta dan Sekitarnya
					Subuh: %v,
					Dzuhur: %v,
					Ashar: %v,
					Maghrib: %v,
					Isha: %v
					`, string(j.JadwalData.Subuh), j.JadwalData.Dzuhur, j.JadwalData.Ashar, j.JadwalData.Maghrib, j.JadwalData.Isha)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, OniChannnn)
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}

			}
		}

		// Jadwal Misa
		// http://www.imankatolik.or.id/kaj.html
		// API Nya belum di temukan

	}
}
