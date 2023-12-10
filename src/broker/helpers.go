package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResp struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// max size of data payload - 1MB
const maxBytes = 1048576

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	// sets the max size of bytes per request, if exceeded error returned
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	// check if single json value is sent, this is done by checking if we reached EOF
	// tries to decode the next JSON-encoded value from decoder into an empty struct
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		// if err is anything other than EOF we need to return it
		return errors.New("body must have only one JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for _, header := range headers {
			for k, v := range header {
				w.Header()[k] = v
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	resp := &jsonResp{
		Error:   true,
		Message: err.Error(),
		Data:    &struct{}{},
	}

	return app.writeJSON(w, statusCode, resp)
}
