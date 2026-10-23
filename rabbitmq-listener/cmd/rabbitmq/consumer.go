package MessagingManager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	connection *amqp.Connection;
	queue amqp.Queue;
	exchangeName string;
	channel *amqp.Channel;
}

type Payload struct{
	Name string `json:"name"`;
	Data string `json:"data"`;
}

type PayloadHandler func(Payload, amqp.Delivery)

func NewConsumer() (*Consumer, error){
	//open connection to rabbitmq
	conn, err:= connectToRabbitMQ()
	if err!=nil{
		return nil, err
	}

	chann,err:=conn.Channel()
	if err!=nil{
		return nil, err
	}

	//broker can only send 100 messages without acknowledgement
	err=chann.Qos(100, 0, false)
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

	queueName:=os.Getenv("QUEUE_NAME")
	if queueName==""{
		return nil, errors.New("Queue name can't be empty")
	}
	//declare listener queue
	q,err:=chann.QueueDeclare(
		queueName, //queue name
		true, //durable
		false, //auto delete
		true, //exclusive?
		false, //no wait
		nil, //arguments
	)
	if err!=nil{
		return nil, err
	}

	c:=&Consumer{
		connection: conn,
		queue: q,
		exchangeName: exchangeName,
		channel: chann,
	}

	//subcribe to topics under interest
	err=c.subscribeTopicsOfListenerExchange([]string{"log.*"})
	if err!=nil{
		return nil, err
	}

	err=c.registerToListenerQueue(handlePayLoad)
	if err!=nil{
		return nil, err
	}

	return c, nil
}

func (c *Consumer) Close(){
	if c.channel!=nil{
		c.channel.Close()
	}
	if c.connection!=nil{
		c.connection.Close()
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
			log.Println(fmt.Sprintf("Connection fails, retry after %d seconds", coolDownTime))
			noConnectionError=err
			time.Sleep(time.Second*time.Duration(coolDownTime))
		}
	}

	return nil, errors.New(fmt.Sprintf("Connection fails after %d tries. Error: %s", noConnectionAmtempt,noConnectionError.Error()))
}

func (c *Consumer) registerToListenerQueue(handler PayloadHandler) error{
	if c.channel==nil{
		return errors.New("Channel is not opened yet")
	}

	//register to the queue
	messages, err:=c.channel.Consume(
		c.queue.Name, //queue name
		"",
		false, //auto ack?
		true, //exclusive?
		false, //no local?
		false, //no wait?
		nil, // arguments
	)
	if err!=nil{
		return err
	}

	fmt.Printf("Listenning for messages from queue %s", c.queue.Name)
	fmt.Println()
	go func(){
		var messageBody Payload
		for message:= range messages{
			err= json.Unmarshal(message.Body, &messageBody)
			if err!=nil{
				log.Println(fmt.Sprintf("Message \"%s\" can't be unmarshalled into Payload", string(message.Body)))
				continue
			}else{
				go handler(messageBody, message)
			}
		}
	}()
	
	return nil
}

//subscribe to topics of the listener's topic exchange
func (c *Consumer) subscribeTopicsOfListenerExchange (topics []string) error{
	if c.channel==nil{
		return errors.New("Channel is not opened yet")
	}
	
	//[not done]: use RabbitMQ transaction for this one
	for _,topic := range topics{
		err:= c.channel.QueueBind(
			c.queue.Name,
			topic,
			c.exchangeName,
			false,
			nil,
		)
		if err!=nil{
			return err
		}
	}

	return nil
}



func handlePayLoad(payload Payload, message amqp.Delivery){
	log.Println("here here")

	switch payload.Name {
	case "log", "event":
		err:=handleLogPayLoad(payload)
		if err!=nil{
			log.Println(err.Error())
		}else{
			log.Println("the following payload is handle: {name:%s, data:%s}", payload.Name, payload.Data)
		}
		message.Ack(false)
	}
}


func handleLogPayLoad(payload Payload) error{
	marshalledData, err:=json.Marshal(payload)
	if err!=nil{
		return err
	}

	client:=&http.Client{
		Timeout: time.Second*10,
	}
	req,err:=http.NewRequest("POST", os.Getenv("LOGGING_URL"), bytes.NewReader(marshalledData))
	if err!=nil{
		return err
	}
	res, err:= client.Do(req)
	if err!=nil{
		return err
	}
	defer res.Body.Close()

	resBody, err:=io.ReadAll(res.Body)
	if err!=nil{
		return err
	}
		
	if res.StatusCode!=http.StatusAccepted{
		return errors.New(
			fmt.Sprintf("Log entry \"%s\" is not recorded by logger. Error message %s", marshalledData, resBody),
		)
	}

	return nil
}

