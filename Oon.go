package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/streadway/amqp"
)

const (
	exchangeGPIOButton = "gpiobutton_events"
	exchangeDCMotors   = "bbdcmotors_ctrl"
)

type OonConfig struct {
	RmqServer   string
	MotorsSpeed int
}

type Oon struct {
	config          OonConfig
	conn            *amqp.Connection
	ch              *amqp.Channel
	gpioButtonQueue amqp.Queue
	killed          bool
}

func InitOon(configFile string) (*Oon, error) {
	var oon Oon
	var err error

	//Load config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &oon.config)
	if err != nil {
		return nil, err
	}

	//Setup AMQP
	oon.conn, err = amqp.Dial(oon.config.RmqServer)
	if err != nil {
		return nil, err
	}

	oon.ch, err = oon.conn.Channel()
	if err != nil {
		return nil, err
	}

	//Setup GPIOButton exchange & queue
	err = oon.ch.ExchangeDeclare(
		exchangeGPIOButton, // name
		"fanout",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, err
	}

	oon.gpioButtonQueue, err = oon.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	//Bind this queue to this exchange so that exchange will publish here
	err = oon.ch.QueueBind(
		oon.gpioButtonQueue.Name, // queue name
		"",                       // routing key
		exchangeGPIOButton,       // exchange
		false,
		nil)

	//Setup DCMotors exchange
	err = oon.ch.ExchangeDeclare(
		exchangeDCMotors, // name
		"fanout",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}

	oon.SetSpeed(oon.config.MotorsSpeed)
	return &oon, nil
}

func (oon *Oon) Destroy() {
	oon.MoveStop()
	oon.ch.Close()
	oon.conn.Close()
}

func (oon *Oon) SetSpeed(speed int) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(speed))
	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_speed",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to set speed")
	}
}

func (oon *Oon) MoveForward() {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(3))
	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_forward",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 3")
	}

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(4))
	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_backward",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 4")
	}
}

func (oon *Oon) MoveForwardDistance(distMm uint32) {

	distToTicks := 100 //TODO
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(3))
	binary.BigEndian.PutUint32(buf[4:8], uint32(distToTicks))

	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_forward_for_ticks",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 3")
	}

	buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(4))
	binary.BigEndian.PutUint32(buf[4:8], uint32(distToTicks))

	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_backward_for_ticks",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 4")
	}
}

func (oon *Oon) MoveBackwardDistance(distMm uint32) {

	distToTicks := 100 //TODO
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(3))
	binary.BigEndian.PutUint32(buf[4:8], uint32(distToTicks))

	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_backward_for_ticks",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 3")
	}

	buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(4))
	binary.BigEndian.PutUint32(buf[4:8], uint32(distToTicks))

	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_forward_for_ticks",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 4")
	}
}

func (oon *Oon) TurnBack() {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(3))
	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_forward",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 3")
	}

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(4))
	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_forward",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to start motor 4")
	}

	time.Sleep(1 * time.Second)
	oon.MoveStop()
}

func (oon *Oon) MoveStop() {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(3))
	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_stop",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to stop motor 3")
	}

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(4))
	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/dcmotor_stop",
			Body:        buf,
		})
	if err != nil {
		log.Println("Failed to stop motor 4")
	}
}

func (oon *Oon) Receive() error {
	msgs, err := oon.ch.Consume(
		oon.gpioButtonQueue.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			switch d.ContentType {
			case "application/button_press":
				oon.SetSpeed(30)
				oon.MoveForwardDistance(1000)
				time.Sleep(5 * time.Second)
				oon.MoveBackwardDistance(1000)
				time.Sleep(5 * time.Second)

			default:
				log.Printf("Received unexpected message: %s", d.Body)
			}
		}
	}()
	return nil
}

func (oon *Oon) Think() error {
	for {
		time.Sleep(100 * time.Millisecond)
		if oon.killed == true {
			break
		}
	}
	return nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " <config>")
		return
	}
	oon, err := InitOon(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig)
			oon.killed = true
		}
	}()

	oon.Receive()
	oon.Think()
	oon.Destroy()
}
