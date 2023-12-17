package main

import (
	"errors"
	"fmt"
	"modules/helpers"
	"net/http"
)

type authReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload authReq
	err := helpers.ReadSingleJson(w, r, &payload)
	if err != nil {
		helpers.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	// Validate user against database
	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		helpers.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Models.User.PasswordMatches(payload.Password)
	if err != nil || !valid {
		helpers.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	resp := helpers.JsonResp{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	helpers.WriteJson(w, http.StatusAccepted, resp)
}
