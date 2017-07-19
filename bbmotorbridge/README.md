# BBMotorBridge

BBMotorBridge is a Motor Bridge driver written in Go for the BeagleBone Motor Bridge Cape.

It is essentially a rewrite of the official package written in Python.

# Usage

    import	"Oon/bbmotorbridge"
    
    func main() {
      mb := bbmotorbridge.New("")
    	if mb == nil {
        log.Fatal("Failed to init motor bridge")
      }

      err := mb.EnableServo(1, true)
      if err != nil {
        log.Fatal(err)
      }

     	err = mb.SetServo(1, 10, 10)
      if err != nil {
        log.Fatal(err)
      }
    }

# FAQ

## Is it fully functional?

Implemented parts are, yes. Nevertheless this package is still under development, DC motors will come soon.

## Will it work on my BB?

This package was tested on BeagleBone Green but will probably work on Black edition too.

### It doesn't work on my BB-xxx, what should I try?

Try editing the i2cAddress, i2cLane and gpioPin consts and set appropriate values.

If you find something working with your BB model, please let me know!