package main

import (
	"encoding/binary"
	"log"

	"github.com/streadway/amqp"
)

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

	requiredTicks := oon.millimetersToTicks(distMm)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:], uint32(3))
	binary.BigEndian.PutUint32(buf[4:], uint32(requiredTicks))
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
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks))

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

	requiredTicks := oon.millimetersToTicks(distMm)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(3))
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks))

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
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks))

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
	//One full rotation on each wheel means 180°
	requiredTicks := oon.millimetersToTicks(oon.config.WheelDiameter)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(3))
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks))

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
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks))

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
