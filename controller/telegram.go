/*
Copyright 2018 Davi Garcia (davivcgarcia@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/davivcgarcia/awx-telegram-bot/config"
	"github.com/davivcgarcia/awx-telegram-bot/model"
	"gopkg.in/tucnak/telebot.v2"
)

// global Telebot bot instance
var mainBot *telebot.Bot

// init initialize the Telebot bot object
func init() {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.EnvVars["TELEGRAM_API_TOKEN"],
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalln("[FATAL] Unable to establish connection to Telegram Bot API")
	}
	mainBot = bot
}

// getTelegramUsername validates if user has username configured on Telegram
func getTelegramUsername(message *telebot.Message) (string, error) {
	if message.Sender.Username != "" {
		return message.Sender.Username, nil
	}
	return "", fmt.Errorf("Telegram user didn't configured an username")
}

// validChatID validates if the Telegram Chat ID is the one allowed to receive commands
func validChatID(message *telebot.Message) bool {
	chatID, _ := strconv.ParseInt(config.EnvVars["TELEGRAM_CHAT_ID"], 10, 64)
	if message.Chat.ID != chatID {
		log.Printf("[ERROR] The chatID %d is not authorized to run commands.", message.Chat.ID)
		mainBot.Leave(message.Chat)
		return false
	}
	return true
}

// checkInHandler responds to /checkin command and does checkout for the user
func checkInHandler(message *telebot.Message) {
	if !validChatID(message) {
		return
	}
	username, err := getTelegramUsername(message)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você não configurou seu Telegram Username.")
		log.Printf("[ERROR] %s", err)
		return
	}
	log.Printf("[INFO] User %s triggered Telegram command /checkin", username)
	err = model.CheckInRegistry(username)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você já fez check-in.")
		log.Printf("[ERROR] %s", err)
		return
	}
	mainBot.Send(message.Chat, "Check-In processado!")
}

// checkOutHandler responds to /checkout command and does checkout for the user
func checkOutHandler(message *telebot.Message) {
	if !validChatID(message) {
		return
	}
	username, err := getTelegramUsername(message)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você não configurou seu Telegram Username.")
		log.Printf("[ERROR] %s", err)
		return
	}
	log.Printf("[INFO] User %s triggered Telegram command /checkout", username)
	err = model.CheckOutRegistry(username)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Você não tem registro de check-in ativo.")
		log.Printf("[ERROR] %s", err)
		return
	}
	mainBot.Send(message.Chat, "Check-Out processado!")
}

// statusHandler responds to /status command and lists all active registries in the database
func statusHandler(message *telebot.Message) {
	if !validChatID(message) {
		return
	}
	username, err := getTelegramUsername(message)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você não configurou seu Telegram Username.")
		log.Printf("[ERROR] %s", err)
		return
	}
	log.Printf("[INFO] User %s triggered Telegram command /status", username)
	registries := model.GetAllActiveRegistries()
	if len(registries) == 0 {
		mainBot.Send(message.Chat, "Não existem registros ativos!")
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("Registros Ativos: ")
	for _, registry := range registries {
		buffer.WriteString(fmt.Sprintf("\n- @%s", registry.Username))
	}
	mainBot.Send(message.Chat, buffer.String())
}

// clearHandler responds to hidden /clear command and checkout all active registries in the database
func clearHandler(message *telebot.Message) {
	if !validChatID(message) {
		return
	}
	username, err := getTelegramUsername(message)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você não configurou seu Telegram Username.")
		log.Printf("[ERROR] %s", err)
		return
	}
	log.Printf("[INFO] User %s triggered Telegram command /clear", username)
	model.CheckOutAllRegistries()
	mainBot.Send(message.Chat, "Limpeza de registros ativos processada! Você não deveria usar isso...")
}

// startHandler responds to /start command and run Ansible Playbook to turn on the lab
func startHandler(message *telebot.Message) {
	if !validChatID(message) {
		return
	}
	username, err := getTelegramUsername(message)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você não configurou seu Telegram Username.")
		log.Printf("[ERROR] %s", err)
		return
	}
	log.Printf("[INFO] User %s triggered Telegram command /start", username)
	if len(model.GetAllActiveRegistries()) > 0 {
		mainBot.Send(message.Chat, "Acho que lab já está ligado!")
		statusHandler(message)
		log.Printf("[INFO] /start command aborted")
		return
	}
	token, err := awxAuthenticate()
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Tivemos problemas com o seu comando... =(")
		log.Println("[ERROR] Problem authenticating with Ansible Tower API")
		return
	}
	awxLaunchJob(token, config.EnvVars["AWX_JOB_TEMPLATE_START"])
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Tivemos problemas com o seu comando... =(")
		log.Println("[ERROR] Problem launching a job with Ansible Tower API")
		return
	}
	checkInHandler(message)
	mainBot.Send(message.Chat, "Acionando o Ansible Tower para ligar o lab!")
}

// stopHandler responds to /stop command and run Ansible Playbook to turn off the lab
func stopHandler(message *telebot.Message) {
	if !validChatID(message) {
		return
	}
	username, err := getTelegramUsername(message)
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Parece que você não configurou seu Telegram Username.")
		log.Printf("[ERROR] %s", err)
		return
	}
	log.Printf("[INFO] User %s triggered Telegram command /stop", username)
	checkOutHandler(message)
	if len(model.GetAllActiveRegistries()) > 0 {
		mainBot.Send(message.Chat, "Não posso desligar pois o ambiente está sendo usado!")
		statusHandler(message)
		log.Printf("[INFO] /stop command aborted")
		return
	}
	token, err := awxAuthenticate()
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Tivemos problemas com o seu comando... =(")
		log.Println("[ERROR] Problem authenticating with Ansible Tower API")
		return
	}
	awxLaunchJob(token, config.EnvVars["AWX_JOB_TEMPLATE_STOP"])
	if err != nil {
		mainBot.Send(message.Chat, "Ops! Tivemos problemas com o seu comando... =(")
		log.Println("[ERROR] Problem launching a job with Ansible Tower API")
		return
	}
	mainBot.Send(message.Chat, "Acionando o Ansible Tower para desligar o lab!")
}

// StartTelegramController does Telegram Controller initialization
func StartTelegramController(exit chan bool) {
	log.Println("[INFO] Starting Telegram Controller")
	mainBot.Handle("/checkin", checkInHandler)
	mainBot.Handle("/checkout", checkOutHandler)
	mainBot.Handle("/status", statusHandler)
	mainBot.Handle("/clear", clearHandler)
	mainBot.Handle("/start", startHandler)
	mainBot.Handle("/stop", stopHandler)
	mainBot.Start()
	exit <- true
}
