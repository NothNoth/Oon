package bbbmotorbridge

/*
Basically a port of https://github.com/Seeed-Studio/MotorBridgeCapeforBBG_BBB/blob/master/BBG_MotorBridgeCape/MotorBridge.py
*/
import (
	"fmt"

	"errors"

	bbhw "github.com/btittelbach/go-bbhw"
	i2c "github.com/d2r2/go-i2c"
)

const (
	i2cAddress = 0x4B
	i2cLane    = 2
	gpioPin    = 49 //Motor bridge PIN maps to P9_23 which is 49
)

type BBBMotorBridge struct {
	i2c *i2c.I2C
}

func (mb *BBBMotorBridge) writeHalfWord(reg byte, value uint16) error {
	var byteSeq []byte

	byteSeq = append(byteSeq, WRITE_MODE)       // Read/Write ?
	byteSeq = append(byteSeq, reg)              //Which register ?
	byteSeq = append(byteSeq, byte(value&0xFF)) //16 bits value
	byteSeq = append(byteSeq, byte((value>>8)&0xFF))
	_, err := mb.i2c.Write(byteSeq)

	if err != nil {
		fmt.Printf("Write failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (mb *BBBMotorBridge) writeByte(reg byte, value byte) error {
	var byteSeq []byte

	byteSeq = append(byteSeq, WRITE_MODE) // Read/Write ?
	byteSeq = append(byteSeq, reg)        //Which register ?
	byteSeq = append(byteSeq, value)      //8 bits value
	_, err := mb.i2c.Write(byteSeq)

	if err != nil {
		fmt.Printf("Write failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (mb *BBBMotorBridge) EnableServo(servo int, enable bool) error {
	_, _, enableReg, err := getRegisters(servo)
	if err != nil {
		return err
	}

	if enable == true {
		err = mb.writeByte(enableReg, 1)
	} else {
		err = mb.writeByte(enableReg, 0)
	}
	return err
}

func getRegisters(servo int) (freq byte, angle byte, enable byte, err error) {
	switch servo {
	case 1:
		freq = SVM1_FREQ
		angle = SVM1_ANGLE
		enable = SVM1_STATE
		break
	case 2:
		freq = SVM2_FREQ
		angle = SVM2_ANGLE
		enable = SVM2_STATE
		break
	case 3:
		freq = SVM3_FREQ
		angle = SVM3_ANGLE
		enable = SVM3_STATE
		break
	case 4:
		freq = SVM4_FREQ
		angle = SVM4_ANGLE
		enable = SVM4_STATE
		break
	case 5:
		freq = SVM5_FREQ
		angle = SVM5_ANGLE
		enable = SVM5_STATE
		break
	case 6:
		freq = SVM6_FREQ
		angle = SVM6_ANGLE
		enable = SVM6_STATE
		break
	default:
		freq = 0
		angle = 0
		enable = 0
		err = errors.New("Invalid servo id (1-6)")
	}
	return
}

func (mb *BBBMotorBridge) SetServo(servo int, angle uint16, speed uint16) error {
	speedReg, angleReg, _, err := getRegisters(servo)
	if err != nil {
		return err
	}
	//Set speed
	err = mb.writeHalfWord(speedReg, speed)
	if err != nil {
		return err
	}
	//Set angle
	err = mb.writeHalfWord(angleReg, angle)
	if err != nil {
		return err
	}
	return nil
}

func New(config string) *BBBMotorBridge {
	var mb BBBMotorBridge
	var err error

	//Setup GPIO / I2C
	reset := bbhw.NewMMappedGPIO(gpioPin, bbhw.OUT)
	reset.SetState(true)

	mb.i2c, err = i2c.NewI2C(i2cAddress, i2cLane)
	if err != nil {
		return nil
	}

	return &mb
}

func (mb *BBBMotorBridge) Destroy() {
	mb.i2c.Close()
}
