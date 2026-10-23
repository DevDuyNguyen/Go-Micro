package main

import (
	"log"
	"logger-service/data/models"
	"net/http"
)

func (app *App) WriteLog(w http.ResponseWriter, r *http.Request){
	var entry models.LogEntry

	err:= app.readJSON(w, r, &entry)
	if err!=nil{
		log.Println("here1")
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err=app.Models.LogEntry.Insert(&entry)
	if err!=nil{
		log.Println("here1")
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, jsonResponse{
		Error: false,
		Message: "Logged",
	}, http.StatusAccepted)

}