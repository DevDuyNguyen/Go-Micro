package main

import (
	"authentication/data"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	port:=flag.String("port", "80", "Port that server listens on")
	flag.Parse()

	db,err:=ConnectToDatabase()
	if err!=nil{
		log.Println(err)
		panic(err)
	}
	app:=App{
		Models: data.NewModels(db),
	}

	server:=&http.Server{
		Addr: fmt.Sprintf(":%s", *port),
		Handler: app.Routes(),
	}

	log.Println("Sever starts on port ", *port)
	err=server.ListenAndServe()
	if err!=nil{
		panic(err)
	}
}

func OpenDB(dsn string) ( *sql.DB, error) {
	db, err:=sql.Open("pgx", dsn)
	if err!=nil{
		return nil, err
	}
	if err:=db.Ping(); err!=nil{
		return nil, err
	}
	return db, nil
}
func ConnectToDatabase() (*sql.DB, error){
	dsn:=os.Getenv("DSN")
	// dsn:="host=localhost port=5432 user=postgres password=111 dbname=duy1 sslmode=disable connect_timeout=5"
	waitTimeInSecondEnv:=os.Getenv("WaitTimeInSecondEnv")
	repeatTimeEnv:=os.Getenv("RepeatTime")
	if waitTimeInSecondEnv==""{
		waitTimeInSecondEnv="5"
	}
	waitTime, err:=strconv.Atoi(waitTimeInSecondEnv)
	if err!=nil{
		return nil, err
	}

	if repeatTimeEnv==""{
		repeatTimeEnv="10"
	}
	repeatTime, err:=strconv.Atoi(repeatTimeEnv)
	
	if err!=nil{
		return nil, err
	}

	var db *sql.DB
	for ;repeatTime>0; repeatTime--{
		time.Sleep(time.Duration(waitTime)*time.Second)
		db, err=OpenDB(dsn)
		if db!=nil && err==nil{
			log.Println("Database connected!!!")
			return db, nil
		}else{
			log.Println("Can't connect to database, because of:")
			log.Println(err.Error())
			log.Println(fmt.Sprintf("Retry after %d seconds", waitTime))
		}
	}
	return nil, data.ErrCanNotConnectToDB
}