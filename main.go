package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/monitoring-api-bpjs/config"
)

func main() {
	config := config.LoadConfig()

	for {
		// Cek status API
		err := checkAPIStatus(config.APIUrl)
		if err != nil {
			// Kirim pesan error ke bot Telegram
			sendTelegramMessage(config.TelegramBot, config.TelegramChat, err.Error())
		}

		// Tunggu sebelum memeriksa ulang
		time.Sleep(5 * time.Minute)
	}
}

// checkAPIStatus checks the status of the API
func checkAPIStatus(apiUrl string) error {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("Failed to reach API: %v", err)
	}
	defer resp.Body.Close()

	// Jika status bukan 200, maka dianggap error
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API Error: %s, Response: %s", resp.Status, string(body))
	}
	return nil
}

// sendTelegramMessage sends a message to the Telegram bot
func sendTelegramMessage(botToken, chatID, message string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data := fmt.Sprintf(`{"chat_id": "%s", "text": "%s"}`, chatID, message)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Printf("Failed to create Telegram request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send Telegram message: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Telegram API Error: %s\n", string(body))
	}
}
