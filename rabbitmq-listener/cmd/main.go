package main

import (
	"log"
	MessagingManager "rabbitmq-listener/cmd/rabbitmq"
)

func main() {
	consumer, err := MessagingManager.NewConsumer()
	if err!=nil{
		log.Println(err.Error())
		return
	}
	defer consumer.Close()

	forever := make(chan int)
	<-forever
}