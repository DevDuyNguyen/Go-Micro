package main

import (
	"fmt"
	"log"
	"mail-service/cmd/mailer"
	"net/http"
	"os"
	"strconv"
)

type App struct{
	Mailer mailer.Mailer
}

func main() {
	port,err:=strconv.Atoi(os.Getenv("MAIL_SERVER_PORT"))
	if err!=nil{
		log.Panic(err)
	}
	mailer,err:=mailer.New(
		os.Getenv("MAIL_SERVER_HOST"),
		port,
		os.Getenv("HAS_SMTP_AUTH")=="1",
		os.Getenv("HAS_TLS_POLICY")=="1",
		os.Getenv("USER_NAME"),
		os.Getenv(("PASSWORD")),
	)
	if err!=nil{
		log.Panic(err)
	}
	app:=App{
		Mailer: *mailer,
	}

	server := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		Handler: app.Routes(),
	}
	fmt.Println("Listen on port %s", os.Getenv("HTTP_PORT"))
	err=server.ListenAndServe()
	if err!=nil{
		log.Panic(err)
	}
}