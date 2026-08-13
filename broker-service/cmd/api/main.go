package main

import (
	"flag"
	"fmt"
	"net/http"
)

type App struct{}

func main() {
	addr := flag.String("addr","3000","Address that the server listens on")
	flag.Parse()

	app:=App{}

	server:= &http.Server{
		Addr: fmt.Sprintf(":%s", *addr),
		Handler: app.routes(),
	}

	fmt.Print("Server start on port ", *addr)

	err:=server.ListenAndServe()
	if err!=nil{
		panic(err.Error())
	}

}