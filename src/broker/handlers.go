package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResp{
		Error:   false,
		Message: "hello from broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

	//out, _ := helpers.MarshalIndent(payload, "", "\t")
	//w.Header().Set("Content-Type", "application/helpers")
	//w.WriteHeader(http.StatusAccepted)
	//w.Write(out)
}
