package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"modules/helpers"
	"net/http"
)

const (
	AUTH_URL = "http://auth-service:420/auth"
	LOG_URL  = "http://logger-service:7070/log"
	MAIL_URL = "http://mailer-service/send"
)

type ReqPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := helpers.JsonResp{
		Error:   false,
		Message: "hello from broker",
	}

	_ = helpers.WriteJson(w, http.StatusOK, payload)

	//out, _ := helpers.MarshalIndent(payload, "", "\t")
	//w.Header().Set("Content-Type", "application/helpers")
	//w.WriteHeader(http.StatusAccepted)
	//w.Write(out)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	log.Println("handling submission")
	var reqPayload ReqPayload

	err := helpers.ReadSingleJson(w, r, &reqPayload)
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.authenticate(w, reqPayload.Auth)
	case "log":
		app.logAction(w, reqPayload.Log)
	case "mail":
		log.Println("mail handler, send mail from broker switch...")
		app.sendMail(w, reqPayload.Mail)
	default:
		helpers.ErrorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) logAction(w http.ResponseWriter, lp LogPayload) {
	log.Println("logging action from broker...")

	// Convert payload to json
	jsonData, _ := json.MarshalIndent(lp, "", "\t")

	// Create new request
	req, err := http.NewRequest("POST", LOG_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}

	// Create new HTTP client
	client := &http.Client{}

	log.Println("making http req...")
	// Make HTTP request
	resp, err := client.Do(req)
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}
	log.Println("made http req")

	// Defer closing BODY
	defer resp.Body.Close()

	log.Println("checking valid status")
	// Check if response is valid
	if resp.StatusCode != http.StatusAccepted {
		helpers.ErrorJson(w, errors.New("failed to log action"))
		return
	}

	log.Println("decoding resp")
	// Decode response
	var jsonResp helpers.JsonResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jsonResp)
	if err != nil {
		helpers.ErrorJson(w, errors.New("failed to decode response"))
		return
	}

	log.Println("internal error")
	// Check for Internal errors
	if jsonResp.Error {
		helpers.ErrorJson(w, errors.New(jsonResp.Message))
		return
	}

	log.Println("writing log response")
	// Send back response
	helpers.WriteJson(w, http.StatusAccepted, &helpers.JsonResp{
		Error:   false,
		Message: "action logged",
		Data:    jsonResp.Data,
	})
	log.Println("action logged successfully...")
}

func (app *Config) authenticate(w http.ResponseWriter, ap AuthPayload) {
	log.Println("authenticating from broker...")
	// Convert auth payload to json
	jsonData, _ := json.MarshalIndent(ap, "", "\t")

	// Call the auth service - Authenticate
	req, err := http.NewRequest("POST", AUTH_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}

	// Create new instance of HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		helpers.ErrorJson(w, errors.New("invalid credentials"))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		helpers.ErrorJson(w, errors.New("error calling auth service"))
		return
	}

	// Read response into variable
	var jsonResp helpers.JsonResp
	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}

	// Check for internal errors
	if jsonResp.Error {
		helpers.ErrorJson(w, err, http.StatusUnauthorized)
		return
	}

	helpers.WriteJson(w, http.StatusAccepted, &helpers.JsonResp{
		Error:   false,
		Message: "authenticated",
		Data:    jsonResp.Data,
	})
}

func (app *Config) sendMail(w http.ResponseWriter, mp MailPayload) {
	log.Println("sending mail from broker service...")

	jsonData, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		log.Println("failed marshaling mail payload", err)
		helpers.ErrorJson(w, err)
		return
	}

	req, err := http.NewRequest("POST", MAIL_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("failed creating mailer request...", err)
		helpers.ErrorJson(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("failed to make mailer request...", err)
		helpers.ErrorJson(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		helpers.ErrorJson(w, err)
		return
	}

	payload := helpers.JsonResp{
		Error:   false,
		Message: fmt.Sprintf("email sent to %s", mp.To),
	}

	helpers.WriteJson(w, http.StatusAccepted, payload)
}
