package bbbmotorbridge_test

import (
	"Oon/bbbmotorbridge"
	"testing"
)

func TestServos(t *testing.T) {
	mb := bbbmotorbridge.New("")
	if mb == nil {
		t.Fail()
	}

	err := mb.EnableServo(1, true)
	if err != nil {
		t.Fail()
	}

	err = mb.EnableServo(0, true)
	if err == nil {
		t.Fail()
	}

	err = mb.EnableServo(7, true)
	if err == nil {
		t.Fail()
	}

	err = mb.SetServo(1, 10, 10)
	if err != nil {
		t.Fail()
	}

}
