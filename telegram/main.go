package main

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"strings"
)

func main(){
	bot, err := tg.NewBotAPI("1098776683:AAGk9n0Ux2nSBTf5xy8dTdYNR_jonZ0pp10")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	fmt.Printf("Авторизованно для аккаунта %s", bot.Self.UserName, "\n")

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	go fire()

	for update := range updates{
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch strings.ToLower(update.Message.Command()) {
		case "start":
			reply = "Привет, данная версия бота написана с использованием библиотеки tgbotapi на golang"
		case "hello":
			reply = "world"

		}

		msg := tg.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}

func fire() {

}
