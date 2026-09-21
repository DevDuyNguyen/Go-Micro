package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data/models"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	mongoClient, err:= connectToMongoDB()
	if err!=nil{
		log.Panic(err)
	}
	defer func(){
		if err=mongoClient.Disconnect(context.TODO()); err!=nil{
			log.Panic(err)
		}
	}()

	app:= App{
		Models: models.New(mongoClient, os.Getenv("DatabaseName"), os.Getenv("CollectionName")),
	}
	
	port:=os.Getenv("Port")

	server := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s",port),
		Handler: app.Routes(),
	}

	err=server.ListenAndServe()
	if err!=nil{
		log.Panic(err)
	}
}

func connectToMongoDB() (*mongo.Client, error){
	clientOpts:=options.Client().ApplyURI(os.Getenv("MongoDBURI")).SetAuth(options.Credential{
		Username: os.Getenv("UserName"),
		Password: os.Getenv("Password"),
	})

	client,err:=mongo.Connect(clientOpts)
	if err!=nil{
		return nil, err
	}

	return client, nil
}