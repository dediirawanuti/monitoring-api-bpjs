package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/joho/godotenv"
)

// Struct untuk membaca konfigurasi dari environment variables
type Config struct {
	APIURL       string
	TelegramBot  string
	TelegramChat string
}

// Struct untuk mem-parsing respons JSON dari endpoint
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Fungsi untuk memuat konfigurasi dari file .env
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	return &Config{
		APIURL:       os.Getenv("API_URL"),
		TelegramBot:  os.Getenv("TELEGRAM_BOT"),
		TelegramChat: os.Getenv("TELEGRAM_CHAT"),
	}, nil
}

// Fungsi untuk memeriksa healthz endpoint
func checkHealthz(url string) error {
	log.Printf("Checking healthz endpoint: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call healthz: %w", err)
	}
	defer resp.Body.Close()

	// Periksa kode status HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Baca seluruh isi respons
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Response body: %s", string(bodyBytes))

	// Parse JSON ke struct
	var healthResp HealthResponse
	if err := json.Unmarshal(bodyBytes, &healthResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Validasi status dari JSON
	if healthResp.Status != "OK" {
		return fmt.Errorf("unexpected status value: %s", healthResp.Status)
	}

	return nil
}

// Fungsi untuk mengirim pesan error ke bot Telegram
func sendTelegramMessage(botToken, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to create request payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
	// resp, err := http.Post(url, "application/json", io.NopCloser(io.NewReader(payload)))
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Telegram: %d", resp.StatusCode)
	}

	return nil
}

// Fungsi utama
func main() {
	// Muat konfigurasi
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Loop untuk memantau endpoint setiap 60 detik
	for {
		err := checkHealthz(config.APIURL)
		if err != nil {
			log.Printf("Health check failed: %v", err)

			// Kirim pesan ke Telegram
			message := fmt.Sprintf("Health check failed: %v", err)
			if err := sendTelegramMessage(config.TelegramBot, config.TelegramChat, message); err != nil {
				log.Printf("Failed to send Telegram message: %v", err)
			}
		} else {
			log.Println("Health check passed.")
		}

		// Tunggu 60 detik sebelum pengecekan berikutnya
		time.Sleep(60 * time.Second)
	}
}
