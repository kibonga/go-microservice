package main

import (
	"fmt"
	"log"
	"modules/helpers"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	log.Println("Do i go into send mail handler")
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var payload mailMessage
	if err := helpers.ReadSingleJson(w, r, &payload); err != nil {
		log.Println("failed to read message from body")
		helpers.ErrorJson(w, err)
	}

	msg := Message{
		From:    payload.From,
		To:      payload.To,
		Subject: payload.Subject,
		Data:    payload.Message,
	}

	err := app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println("failed to send SMTP message...", err)
		helpers.ErrorJson(w, err)
		return
	}

	res := helpers.JsonResp{
		Error:   false,
		Message: fmt.Sprintf("email sent to %s", payload.To),
	}

	helpers.WriteJson(w, http.StatusAccepted, res)
}
