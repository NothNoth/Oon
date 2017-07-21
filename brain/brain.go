package brain

import (
	"Oon/bbmotorbridge"
	"Oon/controls"
	"fmt"
)

const (
	stateIdle = iota
	stateSeek
	stateRotate
	stateAttack
	stateKill
)

const (
	dcMotorDefaultDuty = 1000
)

type BrainHandler struct {
	currentState int
	mb           *bbmotorbridge.BBMotorBridge
	ctrl         *controls.Controls
}

func New(configFile string) *BrainHandler {
	var b BrainHandler
	b.currentState = stateIdle
	b.mb = bbmotorbridge.New(configFile)

	if b.mb == nil {
		fmt.Println("Failed to init motor bridge")
		return nil
	}

	b.ctrl = controls.New()
	if b.ctrl == nil {
		b.mb.Destroy()
		fmt.Println("Failed to init controls")
		return nil
	}

	return &b
}

func (b *BrainHandler) Destroy() {
	b.mb.Destroy()
	b.ctrl.Destroy()
}

func (b *BrainHandler) Start() {
	for {
		switch b.currentState {
		case stateIdle:
			press, _ := b.ctrl.GetPressed()
			if press == true {
				fmt.Println("Button pressed")
				b.stateSwitch(stateSeek)
			}
			break
		case stateSeek:
			press, _ := b.ctrl.GetPressed()
			if press == true {
				fmt.Println("Button pressed")
				b.stateSwitch(stateIdle)
			}
			//Check button press and goto Idle
			//Check frame
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
		//select random angle
		//compute duration from angle
		//start DC motors FW/BW
		b.mb.MoveDC(1, bbmotorbridge.TB_CW, dcMotorDefaultDuty)
		b.mb.MoveDC(2, bbmotorbridge.TB_CCW, dcMotorDefaultDuty)
		break
	case stateAttack:
		// beep^3
		break
	case stateKill:
		// beep^10
		break
	}
}
