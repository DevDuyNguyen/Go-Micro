package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type payload struct{
	Id int `json:"id"`;
	Title string `json:"title"`;
	Price float32 `json:"price"`;
	Description string `json:"description"`;
	Category string `json:"category"`;
	Image string `json:"image"`;
}

func stopOnError(err error){
	if err!=nil{
		panic(err.Error())
	}
}

func main() {
	client := http.Client{
		Timeout: time.Second*10,
	}
	reqPayload:= payload{
		Id: 2,
		Title: "string",
		Price: 0.1,
		Description: "string",
		Category: "string",
		Image: "http://example.com",
	}

	marshalledReqPayload, err:= json.Marshal(reqPayload)
	stopOnError(err)

	req,err:= http.NewRequest("POST", "https://fakestoreapi.com/products", bytes.NewReader(marshalledReqPayload))
	req.Header.Set("Content-Type","application/json")
	stopOnError(err)

	res,err:= client.Do(req)
	defer res.Body.Close()
	stopOnError(err)

	
	marshalledResPayload,err:= io.ReadAll(res.Body)
	stopOnError(err)
	if res.StatusCode==http.StatusCreated{
		var resPayload payload
		err=json.Unmarshal(marshalledResPayload, &resPayload)
		fmt.Println(resPayload)
		stopOnError(err)
	}else{
		fmt.Println("Error with status:", res.Status)
		fmt.Println("Error:", string(marshalledResPayload))
	}
}