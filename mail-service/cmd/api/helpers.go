package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (*App) readJSON(w http.ResponseWriter, r *http.Request, data any) error{
	maxBytes:= 1048576
	r.Body= http.MaxBytesReader(w, r.Body, int64(maxBytes))

	reqBody, err:= io.ReadAll(r.Body)
	if err!=nil{
		return err
	}
	err= json.Unmarshal(reqBody, data)
	if err!=nil{
		return err
	}

	return nil
}

func (*App) writeJSON(w http.ResponseWriter, data any, statusCode int, headers ...http.Header) error{
	marshalledData, err:= json.Marshal(data)
	if err!=nil{
		return err
	}

	w.WriteHeader(statusCode)
	
	//check if user did provide headers
	if len(headers)>0{
		for key,value:=range headers[0]{
			w.Header()[key]=value
		}
	}
	w.Header().Set("Content-Type","application/json")

	_,err= w.Write(marshalledData)
	if err!=nil{
		return err
	}

	return nil
}

func (app *App) errorJSON(w http.ResponseWriter, err error, statusCode ...int) error{
	var status int
	//check if user did provide statusCode
	if len(statusCode)>0{
		status=statusCode[0]
	} else{
		status=http.StatusBadRequest
	}

	jsonError:= jsonResponse{
		Error: true,
		Message: err.Error(),
	}
	
	err= app.writeJSON(w, &jsonError, status, nil)
	if err!=nil{
		return err
	}

	return nil
}