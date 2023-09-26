package rabbitmq

import (
	"github.com/micro/micro/v3/service/logger"
	"github.com/streadway/amqp"
	app "com.nandini.icici/config"
)

func ConsumerAndhraToIcici () {
    rabbitmqUrl :=app.GetVal("GO_MICRO_MESSAGE_BROKER")
	conn,err := amqp.Dial(rabbitmqUrl)

	exchangeName := "AndhraToIcici_message_exchange"
	routingKey := "AndhraToIcici_message_routingKey"
	queueName := "AndhraToIcici_message_queue"
	
	if err != nil {
		logger.Errorf("Failed Initializing Broker Connection")
	}
	channel, err := conn.Channel()
	if err != nil {
		logger.Errorf(err.Error())
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	err = channel.QueueBind(
		queueName,      
		routingKey,     
		exchangeName,   
		false,         
		nil,            
	)

	if err != nil {
		logger.Errorf(err.Error())
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			logger.Infof("Recieved Message: %s\n", d.Body)
		}
	}()
	<-forever
}