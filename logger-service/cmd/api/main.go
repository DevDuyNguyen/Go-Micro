package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"logger-service/data/models"
	"net/http"
	"os"
	"strconv"
	"time"

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
		Models: models.New(mongoClient, os.Getenv("DATABASE_NAME"), os.Getenv("COLLECTION_NAME")),
	}
	
	port:=os.Getenv("PORT")

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
	log.Println("Trying to connect to MongoDB")

	noConnectionAmtempt, err:=strconv.Atoi(os.Getenv("NO_MONGODB_CONNECTION_AMTEMPT"))
	if err!=nil{
		return nil, err
	}
	coolDownTime, err:=strconv.Atoi(os.Getenv("MONGODB_CONNECTION_AMTEMPT_COOLDOWN_TIME"))
	if err!=nil{
		return nil, err
	}
	
	clientOpts:=options.Client().ApplyURI(os.Getenv("MONGO_DB_URI")).SetAuth(options.Credential{
		Username: os.Getenv("USER_NAME"),
		Password: os.Getenv("PASSWORD"),
	})

	var noConnectionError error
	for i:=0; i<noConnectionAmtempt; i++{
		client,err:=mongo.Connect(clientOpts)
		noConnectionError=err

		if err!=nil{
			noConnectionError=err
		}else{
			pingErr:=client.Ping(context.TODO(), nil)
			if pingErr==nil{
				log.Println("Sucessfully connect to mongodb")
				return client, nil
			}else{
				noConnectionError=pingErr
			}
		}

			log.Println(fmt.Sprintf("Connection fails, retry after %d times", noConnectionAmtempt))
			time.Sleep(time.Second*time.Duration(coolDownTime))
	}
	
	return nil, errors.New(fmt.Sprintf("Connection fails after %d tries. Error: %s", noConnectionAmtempt, noConnectionError.Error()))
}