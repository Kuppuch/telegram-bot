package main

import (
	"cloud.google.com/go/firestore"
	"context"
	fb "firebase.google.com/go"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
	"strconv"
	"strings"
	//"google.golang.org/api/option"
	"log"
)

type MessageBody struct {
	id string
	MessageID string
	FirstName string
	Username string
	Text string
	Time string
}

func main() {
	fmt.Print("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	go telegram()

	ctx := context.Background()
	sa := option.WithCredentialsFile("/home/kuppuch/GoLangProject/test-e9d05-firebase-adminsdk-r03nj-409a98c114.json")
	conf := &fb.Config{ProjectID: "test-e9d05"}
	app, err := fb.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	body := MessageBody{}

	for {
		fmt.Println("1 - Вывести все записи")
		fmt.Println("2 - Редактировать запись")
		fmt.Println("3 - Создать запись")
		fmt.Println("4 - Выход")

		var a string
		fmt.Fscan(os.Stdin, &a)
		switch a {
		case "1":
			getAll(client, ctx)
		case "2":
			update(client, ctx, body)
		case "3":
			create(client, ctx, body)
		case "4":
			fmt.Println("Bye ...")
			return
		default:
			fmt.Println("Повторите ввод")
		}
	}

}

func getAll(client *firestore.Client, ctx context.Context) {

	iter := client.Collection("users").Documents(ctx)
	fmt.Printf("%T", iter, "\n")
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			fmt.Println()
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}


func update(client *firestore.Client, ctx context.Context, body MessageBody) {
	_, err := client.Collection("users/" + body.id +"/message").Doc(body.MessageID).Set(ctx, map[string]interface{} {
		"time":  body.Time,
		"text": body.Text,
	})

	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
}

func create(client *firestore.Client, ctx context.Context, body MessageBody) {
	_, err := client.Collection("users").Doc(body.id).Set(ctx, map[string]interface{} {
		"FirstName":  body.FirstName,
		"Username": body.Username,
	})

	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

}

func telegram() {
	bot, err := tg.NewBotAPI("1098776683:AAGk9n0Ux2nSBTf5xy8dTdYNR_jonZ0pp10")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	fmt.Printf("Авторизованно для аккаунта %s", bot.Self.UserName, "\n")

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates{
		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		body := MessageBody{}
		body.id = strconv.Itoa(update.Message.From.ID)
		body.MessageID = strconv.Itoa(update.Message.MessageID)
		body.FirstName = update.Message.From.FirstName
		body.Username = update.Message.From.UserName
		body.Text = update.Message.Text
		body.Time = strconv.Itoa(update.Message.Date)

		switch strings.ToLower(update.Message.Command()) {
		case "start":
			initCommand("create", body)
			reply = "Привет, данная версия бота написана с использованием библиотеки tgbotapi на golang"
		case "hello":
			initCommand("update", body)
			reply = "world"

		}

		switch strings.ToLower(update.Message.Text) {
		case "не знаю что сказать":
			initCommand("update", body)
			reply = "Да"
		case "не знаешь что сказать":
			initCommand("update", body)
			reply = "Именно"
		default:
			initCommand("update", body)
		}

		msg := tg.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

	}
}

func initCommand(command string, body MessageBody) {

	ctx := context.Background()
	sa := option.WithCredentialsFile("/home/kuppuch/GoLangProject/test-e9d05-firebase-adminsdk-r03nj-409a98c114.json")
	conf := &fb.Config{ProjectID: "test-e9d05"}
	app, err := fb.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	switch command {
	case "create":
		create(client, ctx, body)
		update(client, ctx, body)
	case "update":
		update(client, ctx, body)

	}

}
