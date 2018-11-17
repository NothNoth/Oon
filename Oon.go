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
	exchangeGPIOButton     = "gpiobutton_events"
	exchangeDCMotors       = "bbdcmotors_ctrl"
	exchangeDCMotorsEvents = "bbdcmotors_events"
)

type OonConfig struct {
	RmqServer     string
	MotorsSpeed   int
	WheelDiameter uint32
	Button1Name   string
	Button2Name   string
}

type Oon struct {
	config                 OonConfig
	conn                   *amqp.Connection
	ch                     *amqp.Channel
	gpioButtonQueue        amqp.Queue
	dcMotorsQueue          amqp.Queue
	motorsTicksPerRotation uint32
	killed                 bool
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

	//Setup DCMotors exchange (OUT)
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

	//Setup DCMotors exchange (IN)
	err = oon.ch.ExchangeDeclare(
		exchangeDCMotorsEvents, // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, err
	}

	oon.dcMotorsQueue, err = oon.ch.QueueDeclare(
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
		oon.dcMotorsQueue.Name, // queue name
		"",                     // routing key
		exchangeDCMotorsEvents, // exchange
		false,
		nil)

	oon.SetSpeed(oon.config.MotorsSpeed)
	return &oon, nil
}

func (oon *Oon) Destroy() {
	oon.MoveStop()
	oon.ch.Close()
	oon.conn.Close()
}

func (oon *Oon) ReceiveGPIOButton() error {
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
	button1ContentType := fmt.Sprintf("application/button_press_%s", oon.config.Button1Name)
	go func() {
		for d := range msgs {
			switch d.ContentType {
			case button1ContentType:
				oon.SetSpeed(30)
				oon.MoveForwardDistance(1000)
				time.Sleep(10 * time.Second)
				oon.MoveBackwardDistance(1000)
				time.Sleep(10 * time.Second)
				oon.TurnBack()
				time.Sleep(10 * time.Second)

			default:
				log.Printf("Received unexpected message: %s", d.Body)
			}
		}
	}()
	return nil
}

func (oon *Oon) ReceiveDCMotors() error {
	msgs, err := oon.ch.Consume(
		oon.dcMotorsQueue.Name, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			switch d.ContentType {
			case "application/dcmotor_ticks_per_rotation":
				oon.motorsTicksPerRotation = binary.BigEndian.Uint32(d.Body)
				log.Printf("Setting ticks per rotation to: %d\n", oon.motorsTicksPerRotation)
			case "application/dcmotor_autostop":
				stopID := binary.BigEndian.Uint32(d.Body)
				log.Printf("Motor %d has stopped", stopID)
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

	oon.ReceiveGPIOButton()
	oon.ReceiveDCMotors()
	oon.Think()
	oon.Destroy()
}
