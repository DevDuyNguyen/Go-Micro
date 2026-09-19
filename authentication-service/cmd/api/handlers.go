package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *App) Authenticate(w http.ResponseWriter, r *http.Request){
	var reqPayload struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err:=app.readJSON(w, r, &reqPayload)
	if err!=nil{
		app.errorJSON(w, err)
		return
	}

	user,err:= app.Models.User.GetByEmail(reqPayload.Email)
	if err!=nil{
		app.errorJSON(w, err)
		return
	}else if user==nil{
		app.errorJSON(w, errors.New("Invalid Credentials"), http.StatusUnauthorized)
		return
	}

	valid,err:=user.PasswordMatches(reqPayload.Password)
	if err==nil && valid!=false{
		app.errorJSON(w, errors.New("Invalid Credentials"), http.StatusUnauthorized)
		return
	} else if err!=nil{
		app.errorJSON(w, err)
		return
	}

	resPayload:=jsonResponse{
		Error: false,
		Message: fmt.Sprintf("User with email %s is logged in", reqPayload.Email),
		Data: *user,
	}

	app.writeJSON(w, resPayload, http.StatusAccepted)
}