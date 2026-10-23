package publisher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	connection   *amqp.Connection
	exchangeName string
	channel      *amqp.Channel
	outstandingConfirms sync.Map
}

func NewConsumer() (*Publisher, error){
	//open connection to rabbitmq
	conn, err:= connectToRabbitMQ()
	if err!=nil{
		return nil, err
	}

	chann,err:=conn.Channel()
	if err!=nil{
		return nil, err
	}

	//declare listener exchange
	exchangeName:=os.Getenv("EXCHANGE_NAME")
	if exchangeName==""{
		return nil, errors.New("Exchange name can't be empty")
	}

	err=chann.ExchangeDeclare(
		exchangeName, //exchange name
		"topic", //type
		true, //durable?
		false, //autodelete?
		false, //internal
		false, //noWait?
		nil, //arguments
	)
	if err!=nil{
		return nil, err
	}

	c:=&Publisher{
		connection: conn,
		exchangeName: exchangeName,
		channel: chann,
		outstandingConfirms: sync.Map{},
	}

	//enable Publisher Confirm
	err=c.enablePublisherConfirm()
	if err!=nil{
		return nil, err
	}
	

	return c, nil
}

func (p *Publisher) Close(){
	if p.channel!=nil{
		p.channel.Close()
	}
	if p.connection!=nil{
		p.connection.Close()
	}
}

func connectToRabbitMQ() (*amqp.Connection, error){
	noConnectionAmtempt, err:=strconv.Atoi(os.Getenv("NO_RABBITMQ_CONNECTION_AMTEMPT"))
	if err!=nil{
		return nil, err
	}
	coolDownTime, err:=strconv.Atoi(os.Getenv("RABBITMQ_CONNECTION_AMTEMPT_COOLDOWN_TIME"))
	if err!=nil{
		return nil, err
	}

	var noConnectionError error
	for i:=0; i<noConnectionAmtempt; i=i+1{
		conn,err:= amqp.Dial(os.Getenv("RABBITMQ_URI"))
		if err==nil{
			log.Println("Successfully connecting to rabbitmq")
			return conn, nil
		}else{
			log.Println(fmt.Sprintf("Connection fails, retry after %d second", coolDownTime))
			noConnectionError=err
			time.Sleep(time.Second*time.Duration(coolDownTime))
		}
	}

	return nil, errors.New(fmt.Sprintf("Connection fails after %d tries. Error: %s", noConnectionAmtempt,noConnectionError.Error()))
}

func (p *Publisher) enablePublisherConfirm() error{
	err:=p.channel.Confirm(false)
	if err!=nil{
		return err
	}

	go func(){
		confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 100))
		var confirm amqp.Confirmation
		for{
			confirm =<-confirms
			// Clean up map outstandingConfirms when confirms received
			//Log the message when its Publisher Confirm is NACK
			message, ok := p.outstandingConfirms.LoadAndDelete(confirm.DeliveryTag)
			if ok && !confirm.Ack {
				log.Printf("Message with body %s was nacked. Delivery tag: %d", message, confirm.DeliveryTag)
			}else if ok && confirm.Ack{
				log.Printf("Message with body %s was acked. Delivery tag: %d", message, confirm.DeliveryTag)
			}else if !ok{
				log.Printf("Message with body %s and sequence id %d aren't valid", message, confirm.DeliveryTag)
			}
		}
	}()

	return nil
}

func (p *Publisher) PublishMessage(messageBody string, routingKey string, ContentType string, isDurable bool) error{
	sequenceNumer:=p.channel.GetNextPublishSeqNo()
	var publishErr error
	if isDurable{
		publishErr= p.channel.Publish(
			p.exchangeName, 
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: ContentType,
				DeliveryMode: amqp.Persistent,
				Body: []byte(messageBody),
			},
		)
	}else{
		publishErr= p.channel.Publish(
			p.exchangeName, 
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: ContentType,
				Body: []byte(messageBody),
			},
		)
	}
	if publishErr!=nil{
		return publishErr
	}

	p.outstandingConfirms.Store(sequenceNumer, messageBody)
	return nil
}