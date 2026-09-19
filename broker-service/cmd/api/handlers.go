package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

func (app *App) Broker(w http.ResponseWriter, r *http.Request){
	response:=jsonResponse{
		Error: false,
		Message: "Hit the broker",
	}
	app.writeJSON(w, response, http.StatusAccepted)
}

type RequestPayload struct{
	Action string `json:"action"`;
	Auth AuthPayload `json:"auth,omitempty"`;
}
type AuthPayload struct{
	Email string `json:"email"`;
	Password string `json:"password"`;
}

func (app *App) HandleSubmit(w http.ResponseWriter, r *http.Request){
	var reqPayload RequestPayload
	err:=app.readJSON(w, r, &reqPayload)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.authenticate(w, reqPayload.Auth)
	default:
		app.errorJSON(w, errors.New("Unknown action"))
	}


}
func (app *App) authenticate(w http.ResponseWriter, authPayload AuthPayload){
	log.Println("Here1")
	marshalledAuthPayload, err:= json.Marshal(authPayload)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	req,err:=http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewReader(marshalledAuthPayload))
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	client:=&http.Client{
		Timeout: time.Second*10,
	}
	res,err:=client.Do(req)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	marshalledResPayload,err:=io.ReadAll(res.Body)
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshalledResPayload)
}