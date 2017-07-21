package bbmotorbridge_test

import (
	"Oon/bbmotorbridge"
	"testing"
)

func TestInit(t *testing.T) {
	mb := bbmotorbridge.New("")
	if mb == nil {
		t.Error("Init without config is accepted")
	}
	mb.Destroy()

	mb = bbmotorbridge.New("foobar")
	if mb != nil {
		t.Error("Missing config file is not accepted")
	}

	mb = bbmotorbridge.New("example.conf")
	if mb != nil {
		t.Error("Valid config file is accepted")
	}
	mb.Destroy()

}

func TestServos(t *testing.T) {

	mb := bbmotorbridge.New("")

	err := mb.EnableServo(1, true)
	if err != nil {
		t.Error("Servo 1 can be enabled")
	}

	err = mb.EnableServo(0, true)
	if err == nil {
		t.Error("Servo 0 does not exist and cannot be enabled")
	}

	err = mb.EnableServo(7, true)
	if err == nil {
		t.Error("Servo 7 does not exist and cannot be enabled")
	}

	err = mb.SetServo(1, 10, 10)
	if err != nil {
		t.Error("Set position for valid servo is accepted")
	}
	mb.Destroy()

}

func TestDC(t *testing.T) {
	mb := bbmotorbridge.New("")
	err := mb.EnableDC(12, true)
	if err != nil {
		t.Error("Invalid DC motor is rejected")
	}

	err = mb.EnableDC(1, true)
	if err != nil {
		t.Error("Valid DC motor can be enabled")
	}

	err = mb.MoveDC(1, 1, 50)
	if err != nil {
		t.Error("Valid motor can move")
	}

	mb.Destroy()
}
