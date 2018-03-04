package bbmotorbridge_test

import (
	"Oon/bbmotorbridge"
	"fmt"
	"testing"
	"time"
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
	if mb == nil {
		t.Error("Valid config file is accepted")
		return
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

	for i := 1; i <= 6; i++ {
		err := mb.EnableServo(i, true)
		if err != nil {
			t.Errorf("Servo %d can be enabled", i)
		}

		err = mb.SetServo(i, 130, 10)
		if err != nil {
			t.Errorf("Set position for valid servo (%d) is accepted", i)
		}
		time.Sleep(1 * time.Second)
		err = mb.SetServo(i, 30, 10)
		if err != nil {
			t.Errorf("Set position for valid servo (%d) is accepted", i)
		}

		err = mb.EnableServo(i, false)
		if err != nil {
			t.Errorf("Servo %d can be disabled", i)
		}
	}
	mb.Destroy()

}

func TestDC(t *testing.T) {
	mb := bbmotorbridge.New("")
	err := mb.EnableDC(12, true)
	if err == nil {
		t.Error("Invalid DC motor is rejected")
	}

	for i := 1; i <= 4; i++ {
		fmt.Printf("Testing DC #%d\n", i)
		err = mb.EnableDC(i, true)
		if err != nil {
			t.Error("Valid DC motor can be enabled")
		}
		err = mb.MoveDC(i, 1, 50)
		if err != nil {
			t.Error("Valid motor can move")
		}
		time.Sleep(1 * time.Second)
		err = mb.StopDC(i)
		if err != nil {
			t.Error("Valid motor can be stopped")
		}
		err = mb.EnableDC(i, false)
		if err != nil {
			t.Error("Valid DC motor can be disabled")
		}
	}

	mb.Destroy()
}
