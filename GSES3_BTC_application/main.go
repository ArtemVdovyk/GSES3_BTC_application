package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type Subscription struct {
	Email string `json:"email"`
}

var subscriptions []Subscription

var dataFilePath string

type CoinGeckoResponse struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

func main() {
	dataFilePath = "subscriptions.json"
	loadEnvVariables()
	createDataFileIfNotExist()
	getSubscriptionsFromFile()

	http.HandleFunc("/api/rate", getRateHandler)
	http.HandleFunc("/api/subscribe", subscribeHandler)
	http.HandleFunc("/api/sendEmails", sendEmailsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadEnvVariables() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func createDataFileIfNotExist() {
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		file, err := os.Create(dataFilePath)
		if err != nil {
			log.Fatal(err)
		}
		content := []byte("[]")
		err_ := ioutil.WriteFile(dataFilePath, content, 0644)
		if err_ != nil {
			log.Fatal(err_)
		}
		defer file.Close()

		log.Println("File was created:", dataFilePath)
	} else {
		log.Println("File exists:", dataFilePath)
	}
}

func getSubscriptionsFromFile() {
	fileBytes, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(fileBytes, &subscriptions)
	if err != nil {
		log.Fatal(err)
	}
}

func getRateHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data CoinGeckoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rate := data.Bitcoin.UAH

	response := map[string]float64{
		"rate": rate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	for _, sub := range subscriptions {
		if sub.Email == email {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	sub := Subscription{
		Email: email,
	}

	subscriptions = append(subscriptions, sub)

	updateSubData, err := json.Marshal(subscriptions)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dataFilePath, updateSubData, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

func sendEmailsHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"
	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var data CoinGeckoResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rate := data.Bitcoin.UAH

	for _, subscriber := range subscriptions {
		err := sendEmail(subscriber.Email, rate)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", subscriber, err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func sendEmail(recipientEmail string, rate float64) error {
	fromEmail := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	smtpAddress := fmt.Sprintf("%s:%s", smtpServer, smtpPort)

	auth := smtp.PlainAuth("", fromEmail, password, smtpServer)
	message := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Поточний курс BTC до UAH\r\n"+
			"\r\n"+
			"Поточний курс BTC до UAH: %.2f.\r\n", recipientEmail, rate))

	err := smtp.SendMail(smtpAddress, auth, fromEmail, []string{recipientEmail}, message)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", recipientEmail, err)
		return err
	}
	log.Printf("Email sent to %s with rate %.2f", recipientEmail, rate)
	return nil
}
