package bbmotorbridge

/*
Basically a port of https://github.com/Seeed-Studio/MotorBridgeCapeforBBG_BB/blob/master/BBG_MotorBridgeCape/MotorBridge.py
*/
import (
	"fmt"

	"errors"

	"time"

	"encoding/json"

	"io/ioutil"

	bbhw "github.com/btittelbach/go-bbhw"
	i2c "github.com/d2r2/go-i2c"
)

const (
	i2cAddress         = 0x4B
	i2cLane            = 2
	gpioPin            = 49 //Motor bridge PIN maps to P9_23 which is 49
	defaultDCFrequency = 1000
	defaultCmdWait     = 10 * time.Millisecond
)

type ServoState struct {
	Enabled bool
	Angle   int
}

type DCState struct {
	Enabled   bool
	Direction byte
	Duty      uint32
	Inverted  bool
}

//InitialState holds the JSON configuration file
type InitialState struct {
	//Servos holds the state for the 6 servos
	ServoSpeed   uint16
	ServosStates []ServoState
	DCStates     []DCState
}

//BBMotorBridge Motor bridge handler
type BBMotorBridge struct {
	initialState InitialState
	i2c          *i2c.I2C
}

//New creates a new MotorBridge handler, returns nil on failure
func New(config string) *BBMotorBridge {
	var mb BBMotorBridge
	var err error

	mb.initialState.DCStates = make([]DCState, 4)
	mb.initialState.ServosStates = make([]ServoState, 6)

	//Setup GPIO / I2C
	reset := bbhw.NewMMappedGPIO(gpioPin, bbhw.OUT)
	reset.SetState(true)
	time.Sleep(defaultCmdWait)

	mb.i2c, err = i2c.NewI2C(i2cAddress, i2cLane)
	if err != nil {
		return nil
	}
	time.Sleep(defaultCmdWait)

	if len(config) != 0 {
		err = mb.loadInitialState(config)
		if err != nil {
			mb.Destroy()
			return nil
		}
	}

	return &mb
}

//Destroy cleanup resources
func (mb *BBMotorBridge) Destroy() {
	mb.i2c.Close()
}

func (mb *BBMotorBridge) loadInitialState(config string) error {
	data, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(data, &mb.initialState)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(mb.initialState.ServosStates) > 6 {
		return errors.New("Too many initial state for servos defined")
	}
	if len(mb.initialState.DCStates) > 4 {
		return errors.New("Too many initial state for DC defined")
	}

	//Setup servos
	for i := 0; i < len(mb.initialState.ServosStates); i++ {
		mb.EnableServo(i+1, mb.initialState.ServosStates[i].Enabled)
		mb.SetServo(i+1, uint16(mb.initialState.ServosStates[i].Angle), mb.initialState.ServoSpeed)
		fmt.Printf("Servo #%d> Enabled %t - Position: %d\n", i+1, mb.initialState.ServosStates[i].Enabled, mb.initialState.ServosStates[i].Angle)
	}

	//Setup DC
	for i := 0; i < len(mb.initialState.DCStates); i++ {
		mb.EnableDC(i+1, mb.initialState.DCStates[i].Enabled)
		mb.MoveDC(i+1, mb.initialState.DCStates[i].Direction, mb.initialState.DCStates[i].Duty)
		fmt.Printf("DC #%d> Enabled %t - Direction: %d - Duty: %d\n", i+1, mb.initialState.DCStates[i].Enabled, mb.initialState.DCStates[i].Direction, mb.initialState.DCStates[i].Duty)
	}

	return nil
}

//EnableServo switches the given servo on/off. Servo identifier must be within [1-6]
func (mb *BBMotorBridge) EnableServo(servo int, enable bool) error {
	_, _, enableReg, err := getServoRegisters(servo)
	if err != nil {
		return err
	}

	if enable == true {
		err = mb.writeByte(enableReg, 1)
	} else {
		err = mb.writeByte(enableReg, 0)
	}
	time.Sleep(defaultCmdWait)

	return err
}

func (mb *BBMotorBridge) DefaultServoSpeed() uint16 {
	return mb.initialState.ServoSpeed
}

//SetServo sets the given servo index at angle with given speed.
func (mb *BBMotorBridge) SetServo(servo int, angle uint16, speed uint16) error {
	speedReg, angleReg, _, err := getServoRegisters(servo)
	if err != nil {
		return err
	}
	//Set speed
	err = mb.writeHalfWord(speedReg, speed)
	if err != nil {
		return err
	}
	time.Sleep(defaultCmdWait)

	//Set angle
	err = mb.writeHalfWord(angleReg, angle)
	if err != nil {
		return err
	}
	time.Sleep(defaultCmdWait)

	return nil
}

func getServoRegisters(servo int) (freq byte, angle byte, enable byte, err error) {
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

func (mb *BBMotorBridge) writeWord(reg byte, value uint32) error {
	var byteSeq []byte

	byteSeq = append(byteSeq, WRITE_MODE)       // Read/Write ?
	byteSeq = append(byteSeq, reg)              //Which register ?
	byteSeq = append(byteSeq, byte(value&0xFF)) //32 bits value
	byteSeq = append(byteSeq, byte((value>>8)&0xFF))
	byteSeq = append(byteSeq, byte((value>>16)&0xFF))
	byteSeq = append(byteSeq, byte((value>>24)&0xFF))
	_, err := mb.i2c.Write(byteSeq)

	if err != nil {
		fmt.Printf("Write failed: %s\n", err.Error())
		return err
	}
	return nil
}

func (mb *BBMotorBridge) writeHalfWord(reg byte, value uint16) error {
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

func (mb *BBMotorBridge) writeByte(reg byte, value byte) error {
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

func getDCRegisters(dc int) (mode byte, direction byte, duty byte, err error) {
	switch dc {
	case 1:
		mode = TB_1A_MODE
		direction = TB_1A_DIR
		duty = TB_1A_DUTY
		break
	case 2:
		mode = TB_1B_MODE
		direction = TB_1B_DIR
		duty = TB_1B_DUTY
		break
	case 3:
		mode = TB_2A_MODE
		direction = TB_2A_DIR
		duty = TB_2A_DUTY
		break
	case 4:
		mode = TB_2B_MODE
		direction = TB_2B_DIR
		duty = TB_2B_DUTY
		break
	default:
		mode = 0
		direction = 0
		duty = 0
		err = errors.New("Invalid dc id (1-2)")
	}
	return
}

func (mb *BBMotorBridge) EnableDC(dc int, enable bool) error {
	modeReg, directionReg, _, err := getDCRegisters(dc)
	if err != nil {
		return err
	}

	mb.writeWord(CONFIG_TB_PWM_FREQ, defaultDCFrequency)
	time.Sleep(defaultCmdWait)
	mb.writeByte(modeReg, TB_DCM)
	time.Sleep(defaultCmdWait)
	mb.writeByte(directionReg, TB_STOP)
	time.Sleep(defaultCmdWait)

	return nil
}

func (mb *BBMotorBridge) MoveDC(dc int, direction byte, duty uint32) error {
	_, directionReg, dutyReg, err := getDCRegisters(dc)
	if err != nil {
		return err
	}
	if mb.initialState.DCStates[dc-1].Inverted == true {
		if direction == TB_CW {
			direction = TB_CCW
		} else if direction == TB_CCW {
			direction = TB_CW
		}
	}

	mb.writeByte(directionReg, direction)
	time.Sleep(defaultCmdWait)
	mb.writeWord(dutyReg, duty*10)
	time.Sleep(defaultCmdWait)
	return nil
}

func (mb *BBMotorBridge) StopDC(dc int) error {
	_, directionReg, _, err := getDCRegisters(dc)
	if err != nil {
		return err
	}

	mb.writeByte(directionReg, TB_STOP)
	time.Sleep(defaultCmdWait)
	return nil
}
