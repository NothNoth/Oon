package brain

import (
	"Oon/bbmotorbridge"
	"fmt"
)

const (
	stateIdle = iota
	stateSeek
	stateRotate
	stateAttack
	stateKill
)

type BrainHandler struct {
	currentState int
	mb           *bbmotorbridge.BBMotorBridge
}

func New(configFile string) *BrainHandler {
	var b BrainHandler
	b.currentState = stateIdle
	b.mb = bbmotorbridge.New(configFile)

	if b.mb == nil {
		fmt.Println("Failed to init motor bridge")
		return nil
	}
	return &b
}

func (b *BrainHandler) Destroy() {
	b.mb.Destroy()
}

func (b *BrainHandler) Start() {
	for {
		switch b.currentState {
		case stateIdle:
			//Check button press and goto seek
			break
		case stateSeek:
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
		break
	case stateRotate:
		//stop DC motors
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
	switch state {
	case stateIdle:
		//stop motors, put arm back to initial position, beep
		break
	case stateSeek:
		//beep^2, start DC motors
		break
	case stateRotate:
		//select random angle
		//compute duration from angle
		//start DC motors FW/BW
		break
	case stateAttack:
		// beep^3
		break
	case stateKill:
		// beep^10
		break
	}
}
