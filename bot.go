package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Syfaro/telegram-bot-api"
)

var (
	// глобальная переменная в которой храним токен
	telegramBotToken string
)

func init() {
	// принимаем на входе флаг -telegrambottoken
	flag.StringVar(&telegramBotToken, "telegrambottoken", "391630050:AAFd2W9tic05iJmsl-dLWhjkpvnH7520jRU", "Telegram Bot Token")
	flag.Parse()

	// без него не запускаемся
	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
}

func main() {
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Запущен бот", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		// универсальный ответ на любое сообщение
		//var reply string
		reply := "Моя твоя не понимать"
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Println("[", update.Message.From.UserName, "]:", update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		chatID := update.Message.Chat.ID
		cid := strconv.Itoa(int(update.Message.Chat.ID))
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я телеграм-бот"
		case "hello":
			reply = "world"
		case "help":
			reply = "Доступные команды: /start, /hello, /chatid"
		case "chatid":
			reply = cid
		}
		switch update.Message.Text {
		case "Привет":
			reply = "Ну привет :)"
		case "привет":
			reply = "Дарова"
		}
		if update.Message.NewChatMembers != nil {
			newMem := *update.Message.NewChatMembers
			reply = "Привет, @" + newMem[0].UserName
		}
		if update.Message.LeftChatMember != nil {
			reply = "Пока пока, @" + update.Message.LeftChatMember.UserName
		}

		// создаем ответное сообщение клиенту
		//msg := tgbotapi.NewMessage(chatID, reply)
		// отправляем
		if chatID < 0 {
			log.Println("Сообщение пришло из группового чата с id", chatID)
		} else {
			log.Println("Сообщение пришло из приватного чата с id:", chatID)
		}
		msg := tgbotapi.NewMessage(chatID, reply)
		bot.Send(msg)
	}
}
