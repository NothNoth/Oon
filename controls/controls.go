package controls

import bbhw "github.com/btittelbach/go-bbhw"

//Mapping between Px_y notation and port numbers can be found here:
// https://github.com/adafruit/adafruit-beaglebone-io-python/blob/master/source/common.c
// P9_22 is 2
// P8_8 is 67
// etc.

const (
	i2cAddress = 0x4B
	i2cLane    = 2
)

type Controls struct {
	button *bbhw.MMappedGPIO
}

func New() *Controls {
	var ctrl Controls
	ctrl.button = bbhw.NewMMappedGPIO(2, bbhw.IN) // Right grove port is P9_22

	return &ctrl
}

func (ctrl *Controls) Destroy() {
}

func (ctrl *Controls) GetPressed() (bool, error) {
	st, err := ctrl.button.GetState()
	if err != nil {
		return false, err
	}

	return st, nil
}
