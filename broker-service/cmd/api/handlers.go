package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
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
	Log LogPayload `json:"log,omitempty"`;
	Mail MailPayLoad `json:"mail,omitempty"`;
}
type AuthPayload struct{
	Email string `json:"email"`;
	Password string `json:"password"`;
}
type LogPayload struct{
	Name string `json:"name"`;
	Data  string `json:"data"`;
}
type MailPayLoad struct{
	MailFrom string
	MailTo string
	Headers map[string]string
	MessageBody string
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
		break
	case "log":
		// app.log(w, reqPayload.Log)
		app.logViaRabbitMQ(w, reqPayload.Log)
		break
	case "mail":
		app.sendMail(w, reqPayload.Mail)
		break
	default:
		app.errorJSON(w, errors.New("Unknown action"))
	}


}
func (app *App) authenticate(w http.ResponseWriter, authPayload AuthPayload){
	marshalledAuthPayload, err:= json.Marshal(authPayload)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	req,err:=http.NewRequest("POST", os.Getenv("AUTHENTICATING-URL"), bytes.NewReader(marshalledAuthPayload))
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

	resPayload,err:=io.ReadAll(res.Body)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resPayload)
}

func (app *App) log(w http.ResponseWriter, data LogPayload){
	marshalledData, err:=json.Marshal(data)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	client:=&http.Client{
		Timeout: time.Second*10,
	}
	req,err:=http.NewRequest("POST", os.Getenv("LOGGING-URL"), bytes.NewReader(marshalledData))
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	res, err:= client.Do(req)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	resPayload,err:=io.ReadAll(res.Body)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}	
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resPayload)
}
func (app *App) logViaRabbitMQ(w http.ResponseWriter, data LogPayload){
	marshalledData, err:=json.Marshal(data)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err=app.Publisher.PublishMessage(
		string(marshalledData),
		"log.info",
		"text/json",
		true,
	)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resPayload:=jsonResponse{
		Error: false,
		Message: "the log is pushed to rabbitmq",
	}
	textRespayload,err:= json.Marshal(resPayload)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-type", "text/json")
	w.Write(textRespayload)
}

func (app *App) sendMail(w http.ResponseWriter, data MailPayLoad){
	marshalledData,err:=json.Marshal(data)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	req,err:=http.NewRequest("POST", os.Getenv("MAIL_SERVICE_URL"), bytes.NewReader(marshalledData))
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	client:=&http.Client{}
	res,err:=client.Do(req)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	resPayload,err:=io.ReadAll(res.Body)
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(res.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resPayload)
}