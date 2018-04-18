package goiffel

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/streadway/amqp"
)

type AmqpConfig struct {
	AmqpUrl   string `required:"true"`
	QueueName string `required:"true"`
}

type amqpData struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Messages   <-chan amqp.Delivery
}

func NewEiffelChannel(cfg AmqpConfig) (chn *EiffelChannel, err error) {
	channelData := &amqpData{
		Connection: nil,
		Channel:    nil,
	}
	echan := &EiffelChannel{
		ChannelData: channelData,
	}

	channelData.Connection, err = amqp.Dial(cfg.AmqpUrl)
	if err != nil {
		echan.CleanupChannel()
		return nil, err
	}

	channelData.Channel, err = channelData.Connection.Channel()
	if err != nil {
		echan.CleanupChannel()
		return nil, err
	}

	channelData.Queue, err = channelData.Channel.QueueDeclare(
		cfg.QueueName, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		echan.CleanupChannel()
		return nil, err
	}

	return echan, nil
}

func (echan EiffelChannel) RegisterOnEventCallback(cbacks EventCallbacks) (err error) {
	channelData, ok := echan.ChannelData.(*amqpData)
	if !ok {
		return errors.New("Not a Amqp based Eiffel channel")
	}

	channelData.Messages, err = channelData.Channel.Consume(
		channelData.Queue.Name, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range channelData.Messages {
			log.Printf("Received a message: %s", d.Body)
			var evt EiffelEvent
			json.Unmarshal(d.Body, &evt)

			cb := cbacks[evt.Meta.Type]
			if cb != nil {
				postReceiveParser(&evt)
				cb(evt)
			} else {
				cb = cbacks[DefaultEiffelEvent]
				if cb != nil {
					cb(evt)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	return nil
}

func (echan EiffelChannel) CleanupChannel() (err error) {
	channelData, ok := echan.ChannelData.(amqpData)
	if !ok {
		return errors.New("Not a Amqp based Eiffel channel")
	}

	if channelData.Channel != nil {
		channelData.Channel.Close()
	}

	if channelData.Connection != nil {
		channelData.Connection.Close()
	}

	return nil
}

func (echan EiffelChannel) TransmitEiffelEvent(evt EiffelEvent) (err error) {
	channelData, ok := echan.ChannelData.(*amqpData)
	if !ok {
		return errors.New("Not a Amqp based Eiffel channel")
	}

	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	err = channelData.Channel.Publish(
		"", // exchange
		channelData.Queue.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}
