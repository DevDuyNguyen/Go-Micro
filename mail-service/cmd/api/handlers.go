package main

import (
	"fmt"
	"mail-service/cmd/mailer"
	"net/http"
)


func (app *App) SendMail(w http.ResponseWriter, r *http.Request){
	var mess mailer.SMTPMessage
	err:= app.readJSON(w,r,&mess)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err=app.Mailer.SendSMTPMAil(mess)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, jsonResponse{
		Error: false,
		Message: fmt.Sprint("Mail is sent successfully to %s", mess.MailTo),
	}, http.StatusAccepted)
}