package main

import (
	"errors"
	"fmt"
	"log"
	"modules/helpers"
	"net/http"
)

type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthReq

	// Read request body(stream) into json variable
	err := helpers.ReadSingleJson(w, r, &req)
	if err != nil {
		helpers.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	// Get the user by email from json variable
	user, err := app.Models.User.GetByEmail(req.Email)
	if err != nil {
		log.Println("failed to get by email")
		helpers.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// Validate if passwords match
	valid, err := user.PasswordMatches(req.Password)
	if err != nil || !valid {
		log.Println("failed to match passwords")
		helpers.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// Send back response
	helpers.WriteJson(w, http.StatusAccepted, &helpers.JsonResp{
		Error:   false,
		Message: fmt.Sprintf("logged in user %v", user.Email),
		Data:    user,
	})
}
