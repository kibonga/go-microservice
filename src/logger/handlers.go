package main

import (
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

	logEntry := data.LogEntry{
		Name: reqPayload.Name,
		Data: reqPayload.Data,
	}

	err := app.Models.LogEntry.Insert(logEntry)
	if err != nil {
		helpers.ErrorJson(w, err)
		return
	}

	resp := helpers.JsonResp{
		Error:   false,
		Message: "successfully logged in",
	}

	helpers.WriteJson(w, http.StatusAccepted, resp)

}
