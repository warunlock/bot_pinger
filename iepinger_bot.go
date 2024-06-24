package main

import (
	"encoding/json"
//	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"net"

	"github.com/tatsushid/go-fastping"

	//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramBotToken string `json:"TelegramBotToken"`
	TO               int    `json:"TO"`
	Hosts            []struct {
		Name    string `json:"Name"`
		IP      string `json:"IP"`
		Chat    int64  `json:"Chat"`
		CStatus int    `json:"CStatus"`
	} `json:"Hosts"`
}

func main() {
	runtime.GOMAXPROCS(2)

	//make log file
/*
	f, err := os.OpenFile("epinger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
*/
	//==========reading config file
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	file.Close()
	//==========reading config file

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() { //sending message form mts.msg to user

		/*msg := tgbotapi.NewInlineQueryResultArticleMarkdown("-635084032", "yes!", "yes!!")

		inlineConfig := tgbotapi.InlineConfig{
			InlineQueryID: "-635084032",
			IsPersonal:    true,
			CacheTime:     0,
			Results:       []interface{}{msg},
		}

		if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
			log.Println(err)
		}*/
		/*	msg := tgbotapi.NewMessage(-635084032, "yes2")
			if _, err = bot.Send(msg); err != nil {
				log.Println(err)
			}
		*/
		for {
			for inx, _ := range configuration.Hosts {
				host := &configuration.Hosts[inx]
				res := 0
				p := fastping.NewPinger()

				ra, err := net.ResolveIPAddr("ip4:icmp", host.IP)
				if err != nil {
					log.Println(err)
				}
				p.AddIPAddr(ra)
				p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
					res = res + 1
				}

				err = p.Run()
				err = p.Run()
				if err != nil {
					log.Println(err)
				}
				if res > 0 {
					res = 1
				}
				if host.CStatus == -1 {
					host.CStatus = res
					log.Println(host.Name, host.CStatus)
				} else {
					if host.CStatus != res {
						switch res {
						case 0:
							{
								log.Println(host.Name, "Свет выключили  :(")
								msg := tgbotapi.NewMessage(host.Chat, host.Name+" Світло вимкнули  :(    \xF0\x9F\x94\xA6")
								if _, err = bot.Send(msg); err != nil {
									log.Println(err)
								}

							}
						case 1:
							{
								log.Println(host.Name, "Свет включили!!! :)")
								msg := tgbotapi.NewMessage(host.Chat, host.Name+" Світло увімкнули!!!  	\xF0\x9F\x92\xA1")
								if _, err = bot.Send(msg); err != nil {
									log.Println(err)
								}
							}
						}
						host.CStatus = res
						log.Println(host.Name, "Status changed to: ", host.CStatus)

					} else {
						log.Println(host.Name+" no changes ", host.CStatus)
					}

				}

			}
			time.Sleep(time.Second * time.Duration(configuration.TO))
		}
	}()

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.

		if update.Message == nil && update.InlineQuery != nil {
			log.Println("1111")
			//log.Println(update.InlineQuery)
			var msgs []interface{}
			// код для inline режима
			msg := tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, "yes!", update.InlineQuery.Query)
			msgs = append(msgs, msg)

/*			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       []interface{}{msg},
			}

			if _, err := bot.AnswerInlineQuery(inlineConfig); err != nil {
				log.Println(err)
			}*/

		} else {

			log.Println("hmmmm........................")
			log.Println(update)

		}

	}
}
