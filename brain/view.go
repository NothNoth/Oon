package brain

import (
	"Oon/bbmotorbridge"
	"errors"
	"fmt"
	"image"
	"math"
	"time"
)

const (
	rotateCalibrationMaxDelay = 10 * time.Second
)

func diffFrame(root *image.Image, current *image.Image) float64 {

	currentS := (*current).Bounds().Size()
	rootS := (*root).Bounds().Size()

	if rootS.Eq(currentS) == false {
		fmt.Println("Comparing images of different sizes")
		return 0.0
	}

	var totalDiff float64
	for x := 0; x < rootS.X; x++ {
		for y := 0; y < rootS.Y; y++ {
			rootR, rootG, rootB, _ := (*root).At(x, y).RGBA()
			curR, curG, curB, _ := (*current).At(x, y).RGBA()
			totalDiff += math.Abs(float64(rootR-curR)) + math.Abs(float64(rootG-curG)) + math.Abs(float64(rootB-curB))
		}
	}
	max := float64(rootS.X) * float64(rootS.Y) * 3.0 * 0xFFFF
	return (max - totalDiff) / max
}

func (b *BrainHandler) calibrateRotationWithLevel(calibrationLevel float64) (time.Duration, error) {
	if b.cam == nil {
		return 0.0, errors.New("Calibration failed (no camera)")
	}
	b.mb.MoveDC(1, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)
	b.mb.MoveDC(2, bbmotorbridge.TB_STOP, dcMotorDefaultDuty)

	rootFrame := b.cam.GrabFrameWithTimeout(5 * time.Second)
	if rootFrame == nil {
		return 0.0, errors.New("Calibration timeout (failed to retrieve root frame)")
	}
	tsStart := time.Now()
	b.mb.MoveDC(1, bbmotorbridge.TB_CW, dcMotorDefaultDuty)
	b.mb.MoveDC(2, bbmotorbridge.TB_CCW, dcMotorDefaultDuty)

	//Fetching new frame too early would lead to instant match
	time.Sleep(500 * time.Millisecond)
	for {
		curFrame := b.cam.GrabFrameWithTimeout(2 * time.Second)
		if curFrame == nil {
			return 0.0, errors.New("Calibration timeout (failed to retrieve diff frame)")
		}
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