package main

import (
	"log"
	"logger/data"
	"modules/helpers"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) AddLog(w http.ResponseWriter, r *http.Request) {
	var reqPayload JSONPayload
	_ = helpers.ReadSingleJson(w, r, &reqPayload)

	log.Println("request payload", reqPayload)
	logEntry := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	log.Println("adding log entry...")
	log.Println("log entry to be added:", logEntry)
	err := app.Models.LogEntry.Insert(logEntry)
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}
	log.Println("log entry added")

	resp := helpers.JsonResp{
		Error:   false,
		Message: "successfully logged in",
	}

	helpers.WriteJson(w, http.StatusAccepted, resp)
}
