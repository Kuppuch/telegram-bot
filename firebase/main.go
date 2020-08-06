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
	FirstName string
	Username string
	Text string
}

func main() {

	go telegram()

	ctx := context.Background()
	sa := option.WithCredentialsFile("/home/kuppuch/GoLangProject/test-e9d05-firebase-adminsdk-r03nj-4e1eaf4351.json")
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
			update(client, ctx)
		case "3":
			create(client, ctx)
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

func update(client *firestore.Client, ctx context.Context) {
	getAll(client, ctx)
}

func updateT(client *firestore.Client, ctx context.Context, body MessageBody) {
	getAll(client, ctx)
}

func create(client *firestore.Client, ctx context.Context) {
	/*_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
		"FirstName":  "Alan",
		"Username": "Walker",
		"Text":   "Darkside",
		"id":   2000,
	})
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}*/
	_, err := client.Collection("users").Doc("LA").Set(ctx, map[string]interface{}{
		"FirstName":  "Some",
		"Username": "Pes",
		"Text":   "Awesome",
		"id":   1990,
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

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch strings.ToLower(update.Message.Command()) {
		case "start":
			reply = "Привет, данная версия бота написана с использованием библиотеки tgbotapi на golang"
		case "hello":
			reply = "world"

		}

		msg := tg.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

		body := MessageBody{}
		body.id = strconv.Itoa(update.Message.From.ID)
		body.FirstName = update.Message.From.FirstName
		body.Username = update.Message.From.UserName
		body.Text = update.Message.Text


		//log.Panic(body.id, "AAAAAAAAAAAAAAAAAAAAAAA")
	}
}