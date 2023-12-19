package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"modules/helpers"
	"net/http"
)

type ReqPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	default:
		helpers.ErrorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, ap AuthPayload) {
	log.Println("trying to authenticate broker")
	// Convert auth payload to json
	jsonData, _ := json.MarshalIndent(ap, "", "\t")

	// Call the auth service - Authenticate
	req, err := http.NewRequest("POST", "http://auth-service:420/auth", bytes.NewBuffer(jsonData))
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

	log.Println("finally authenticated")
	helpers.WriteJson(w, http.StatusAccepted, &helpers.JsonResp{
		Error:   false,
		Message: "authenticated",
		Data:    jsonResp.Data,
	})
}
