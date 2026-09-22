package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LogEntry struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

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

	err=app.logRequest(LogEntry{
		Name:"Authentication Service",
		Data:fmt.Sprintf("%s is logged in", reqPayload.Email),
	})
	if err!=nil{
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resPayload:=jsonResponse{
		Error: false,
		Message: fmt.Sprintf("User with email %s is logged in", reqPayload.Email),
		Data: *user,
	}
	app.writeJSON(w, resPayload, http.StatusAccepted)
}

func (*App) logRequest(entry LogEntry) error{
	marshalledEntry,err:= json.Marshal(entry)
	if err!=nil{
		return err
	}

	req,err:=http.NewRequest("POST", os.Getenv("LOGGING-URL"), bytes.NewReader(marshalledEntry))
	if err!=nil{
		return err
	}

	client:=&http.Client{
		Timeout: time.Second*10,
	}
	res,err:=client.Do(req)
	if err!=nil{
		return err
	}
	defer res.Body.Close()

	if res.StatusCode!=http.StatusAccepted{
		return errors.New("Log fails")
	}
	
	return nil
}