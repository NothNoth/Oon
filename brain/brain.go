package brain

import (
	"Oon/bbmotorbridge"
	"Oon/camera"
	"Oon/controls"
	"fmt"
	"math/rand"
	"time"
)

const (
	stateIdle = iota
	stateSeek
	stateRotate
	stateAttack
	stateKill
)

const (
	dcMotorDefaultDuty   = 1000
	defaultRotationSpeed = 5 * time.Second
)

type BrainHandler struct {
	currentState  int
	mb            *bbmotorbridge.BBMotorBridge
	ctrl          *controls.Controls
	cam           *camera.Camera
	rotationSpeed time.Duration
}

func New(motorBridgeConfig string, cameraConfig string) *BrainHandler {
	var b BrainHandler
	b.currentState = stateIdle
	fmt.Println("Setting up brain...")
	b.mb = bbmotorbridge.New(motorBridgeConfig)
	if b.mb == nil {
		fmt.Println("Failed to init motor bridge")
		return nil
	}
	fmt.Println("> Motors ready.")

	b.ctrl = controls.New()
	if b.ctrl == nil {
		b.mb.Destroy()
		fmt.Println("Failed to init controls")
		return nil
	}
	fmt.Println("> Controls ready.")

	b.cam = camera.New(cameraConfig)
	if b.cam == nil {
		fmt.Println("Failed to init camera")
		return nil
	}
	fmt.Println("> Camera ready.")

	return &b
}

func (b *BrainHandler) Destroy() {
	b.mb.Destroy()
	b.ctrl.Destroy()
}

func (b *BrainHandler) delayedStateSwitch(newState int, wait time.Duration) {
	time.Sleep(wait)
	b.stateSwitch(newState)
}

func (b *BrainHandler) Start() {
	var err error

	fmt.Println("Calibrating rotation...")
	b.rotationSpeed, err = b.calibrateRotation()
	if err != nil {
		fmt.Printf("Rotation calibration failed: %s, using default\n", err.Error())
		b.rotationSpeed = defaultRotationSpeed
	} else {
		fmt.Printf("Rotation calibration succeeded: %s\n", err.Error())
	}

	for {
		switch b.currentState {
		case stateIdle:
			//On button press, seek
			press, _ := b.ctrl.GetPressed()
			if press == true {
				fmt.Println("Button pressed")
				b.stateSwitch(stateSeek)
				time.Sleep(1 * time.Second) //leave enough time to release button
			}
			break
		case stateSeek:
			//On button press, idle
			press, _ := b.ctrl.GetPressed()
			if press == true {
				fmt.Println("Button pressed")
				b.stateSwitch(stateIdle)
				time.Sleep(1 * time.Second) //leave enough time to release button
			}

			//TODO: check frame
			// - grass: continue
			// - weed: goto attack
			// - other: goto rotate
			break
		case stateRotate:
			//if timeout: goto seek
			break
		case stateAttack:
			//Move arm, weed focus detect
			// If not perfect, continue
			// If perfect, goto kill
			// If timeout, goto seek
			break
		case stateKill:
			//Grab 3 times
			// weed focus detect
			// if NOT perfect goto seek
			// if perfect, continue
			// If timeout, goto seek
			break
		}

	}
}

func stateToStr(state int) string {
	switch state {
	case stateAttack:
		return "Attack"
	case stateIdle:
		return "Idle"
	case stateKill:
		return "Kill"
	case stateRotate:
		return "Rotate"
	case stateSeek:
		return "Seek"
	default:
		return "???"
	}
}

func (b *BrainHandler) stateSwitch(newState int) {
	b.endState(b.currentState)
	b.currentState = newState
	b.startState(newState)
}

func (b *BrainHandler) endState(state int) {
	switch state {
	case stateIdle:
		//nothing
		break
	case stateSeek:
		//stop DC motors
		b.mb.MoveDC(1, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
		b.mb.MoveDC(2, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
		break
	case stateRotate:
		//stop DC motors
		b.mb.MoveDC(1, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
		b.mb.MoveDC(2, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
		break
	case stateAttack:
		//nothing
		break
	case stateKill:
		//put arm back to initial position
		break
	}
}

func (b *BrainHandler) startState(state int) {
	fmt.Println("Entering state" + stateToStr(state))
	switch state {
	case stateIdle:
		//stop motors, put arm back to initial position, beep
		b.mb.MoveDC(1, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
		b.mb.MoveDC(2, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
		break
	case stateSeek:
		//beep^2, start DC motors
		b.mb.MoveDC(1, bbmotorbridge.TB_CW, dcMotorDefaultDuty)
		b.mb.MoveDC(2, bbmotorbridge.TB_CW, dcMotorDefaultDuty)
		break
	case stateRotate:
		b.mb.MoveDC(1, bbmotorbridge.TB_CW, dcMotorDefaultDuty)
		b.mb.MoveDC(2, bbmotorbridge.TB_CCW, dcMotorDefaultDuty)

		//generate random rotation duration
		r := rand.Uint32() % 8
		duration := b.rotationSpeed / 2
		duration = duration + duration/time.Duration(r)
		fmt.Printf("Will rotate for %s", duration.String())
		go b.delayedStateSwitch(stateSeek, duration)
		break
	case stateAttack:
		// beep^3
		break
	case stateKill:
		// beep^10
		break
	}
}
