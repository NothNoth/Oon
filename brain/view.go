package brain

import "Oon/bbmotorbridge"
import "time"
import "image"
import "errors"

const (
	rotateCalibrationMaxDelay = 10 * time.Second
)

func diffFrame(root *image.Image, current *image.Image) float64 {
	return 0.0
}

func (b *BrainHandler) calibrateRotationWithLevel(calibrationLevel float64) (time.Duration, error) {
	b.mb.MoveDC(1, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
	b.mb.MoveDC(2, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)

	rootFrame := b.cam.GrabFrameWithTimeout(100 * time.Millisecond)

	tsStart := time.Now()
	b.mb.MoveDC(1, bbmotorbridge.TB_CW, dcMotorDefaultDuty)
	b.mb.MoveDC(2, bbmotorbridge.TB_CCW, dcMotorDefaultDuty)

	//Fetching new frame too early would lead to instant match
	time.Sleep(500 * time.Millisecond)
	for {
		curFrame := b.cam.GrabFrameWithTimeout(100 * time.Millisecond)
		diff := diffFrame(rootFrame, curFrame)

		if diff > calibrationLevel {
			b.mb.MoveDC(1, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
			b.mb.MoveDC(2, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
			break
		}

		if time.Since(tsStart) > rotateCalibrationMaxDelay {
			return 0.0, errors.New("Calibration timeout")
		}
	}

	return time.Since(tsStart), nil
}

func (b *BrainHandler) calibrateRotation() (time.Duration, error) {
	calibrationLevel := 1.0

	for {
		t, err := b.calibrateRotationWithLevel(calibrationLevel)
		if err == nil {
			return t, nil
		}
		calibrationLevel -= 0.1
		if calibrationLevel < 0.0 {
			return 0.0, errors.New("Calibration failed")
		}
	}
}
