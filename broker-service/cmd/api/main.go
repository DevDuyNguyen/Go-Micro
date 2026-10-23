package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"duy.broker.server/cmd/publisher"
)

type App struct{
	Publisher *publisher.Publisher
}

func main() {
	addr := flag.String("addr","80","Address that the server listens on")
	flag.Parse()

	publisher, err:=publisher.NewConsumer()
	if err!=nil{
		panic(err.Error())
	}

	app:=App{
		Publisher: publisher,
	}

	server:= &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", *addr),
		Handler: app.routes(),
	}

	log.Print("Server start on port ", *addr)

	err=server.ListenAndServe()
	if err!=nil{
		panic(err.Error())
	}

}