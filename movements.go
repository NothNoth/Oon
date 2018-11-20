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
	oon.motorsSendCmd(3, "application/dcmotor_forward", 4, "application/dcmotor_backward")
}

func (oon *Oon) MoveForwardDistance(distMm uint32) {

	requiredTicks := oon.millimetersToTicks(distMm)
	if requiredTicks == 0 {
		return
	}
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
	if requiredTicks == 0 {
		return
	}
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
	oon.motorsSendCmdDist(3, oon.config.WheelDiameter, "application/dcmotor_backward_for_ticks", 4, oon.config.WheelDiameter, "application/dcmotor_backward_for_ticks")
}

func (oon *Oon) MoveStop() {
	oon.motorsSendCmd(3, "application/dcmotor_stop", 4, "application/dcmotor_stop")
}

func (oon *Oon) millimetersToTicks(distMm uint32) uint32 {

	if oon.motorsTicksPerRotation == 0 {
		log.Println("Didn't received motorsTicksPerRotation yet")
		return 0
	}

	return distMm * oon.motorsTicksPerRotation / oon.config.WheelDiameter
}

func (oon *Oon) motorsSendCmd(motorID1 uint32, cmd1 string, motorID2 uint32, cmd2 string) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(motorID1))
	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: cmd1,
			Body:        buf,
		})
	if err != nil {
		log.Printf("Failed to send command %s to motor %d", cmd1, motorID1)
	}

	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(motorID2))
	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: cmd2,
			Body:        buf,
		})
	if err != nil {
		log.Printf("Failed to send command %s to motor %d", cmd2, motorID2)
	}
}

func (oon *Oon) motorsSendCmdDist(motorID1 uint32, dist1 uint32, cmd1 string, motorID2 uint32, dist2 uint32, cmd2 string) {
	requiredTicks1 := oon.millimetersToTicks(dist1)
	if requiredTicks1 == 0 {
		return
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(motorID1))
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks1))

	err := oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: cmd1,
			Body:        buf,
		})
	if err != nil {
		log.Printf("Failed to run cmd %s to motor %d", cmd1, motorID1)
	}

	requiredTicks2 := oon.millimetersToTicks(dist2)
	if requiredTicks2 == 0 {
		return
	}
	buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf[0:4], uint32(motorID2))
	binary.BigEndian.PutUint32(buf[4:8], uint32(requiredTicks2))

	err = oon.ch.Publish(
		exchangeDCMotors, // exchange
		"",               // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: cmd2,
			Body:        buf,
		})
	if err != nil {
		log.Printf("Failed to run cmd %s to motor %d", cmd2, motorID2)
	}
}
