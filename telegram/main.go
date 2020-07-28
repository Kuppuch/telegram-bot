package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	//"github.com/go-telegram-bot-api/telegram-bot-api"
)

type UpdateT struct {
	Ok bool `json:"ok"`
	Result []UpdateResultT `json:"result"`
}

type UpdateResultT struct {
	UpdateId int `json:"update_id"`
	Message UpdateResultMessageT `json:"message"`
}

type UpdateResultMessageT struct {
	MessageId int `json:"message_id"`
	From UpdateResultFromT `json:"from"`
	Chat UpdateResultChatT `json:"chat"`
	Date int `json:"date"`
	Text string `json:"text"`
	//Entities UpdateResultEntitiesT `json:"entities,omitempty"`
}

type UpdateResultFromT struct {
	Id int `json:"int"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	UserName string `json:"username"`
	Language string `json:"language_code"`
}

type UpdateResultChatT struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Type string `json:"private"`
}

type UpdateResultEntitiesT struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
	Type string `json:"bot_command"`
}

type SendMessageResponseT struct {
	Ok bool `json:"ok"`
	Result UpdateResultMessageT `json:"message"`
}

const baseTelegramUrl = "https://api.telegram.org"
const telegramToken = "1098776683:AAGk9n0Ux2nSBTf5xy8dTdYNR_jonZ0pp10"
const getUpdatesUri = "getUpdates"
const sendMessageUri = "sendMessage"

const keywordStart = "/start"

func main() {
	update, err := getUpdates()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(update)

	for _, item := range(update.Result) {
		fmt.Println(item.Message.Chat.Id, item.Message.From.FirstName, item.Message.Text)
		if strings.Contains(strings.ToLower(item.Message.Text), "пидр") {
			sendMessage(item.Message.Chat.Id, "Согласен, тот ещё " + item.Message.Text)
			continue
		}
		sendMessage(item.Message.Chat.Id, item.Message.From.FirstName + ", ты это написал(а)? " + "\"" + item.Message.Text + "\"")
	}

	result, err := sendMessage(538632285, "на тебе ответ")
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	if !result.Ok {
		fmt.Println("Сообщение не отправлено")
	}

	fmt.Println(result)
}

func getUpdates() (UpdateT, error) {
	url := baseTelegramUrl + "/bot" + telegramToken + "/" + getUpdatesUri
	response := getResponse(url)

	update := UpdateT{}
	err := json.Unmarshal(response, &update)
	if err != nil {
		return update, err
	}

	return  update, nil
}

func sendMessage(chatId int, text string) (SendMessageResponseT, error){
	url := baseTelegramUrl + "/bot" + telegramToken + "/" + sendMessageUri
	url = url + "?chat_id=" + strconv.Itoa(chatId) + "&text=" + text
	response := getResponse(url)

	fmt.Println(string(response))

	sendMessage := SendMessageResponseT{}
	err := json.Unmarshal(response, &sendMessage)
	if err != nil {
		return sendMessage, err
	}

	return sendMessage, nil
}

func getResponse(url string) []byte {
	response := make([]byte, 0)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

		return response
	}

	defer resp.Body.Close()

	for true {
		bs := make ([]byte, 1024)
		n, err := resp.Body.Read(bs)
		response = append(response, bs[:n]...)

		if n == 0 || err != nil {
			break
		}
	}

	return response
}
