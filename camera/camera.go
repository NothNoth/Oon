package camera

import (
	"encoding/json"
	"fmt"
	"image"

	"bytes"

	"image/jpeg"

	"io/ioutil"

	"time"

	"sync"

	"github.com/blackjack/webcam"
)

type CameraConfig struct {
	Device   string
	Encoding string
	Format   uint32
	Width    uint32
	Height   uint32
}
type Camera struct {
	camH          *webcam.Webcam
	frameEncoding string
	stopGrabber   bool
	lastImg       *image.Image
	lastImgTs     time.Time
	lastLock      sync.Mutex
}

func New(config string) *Camera {
	var cam Camera
	var cameraConfig CameraConfig
	data, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = json.Unmarshal(data, &cameraConfig)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	cam.camH, err = webcam.Open(cameraConfig.Device) // Open webcam
	if err != nil {
		fmt.Println("Failed to open device " + cameraConfig.Device)
		return nil
	}

	_, _, _, err = cam.camH.SetImageFormat(webcam.PixelFormat(cameraConfig.Format), cameraConfig.Width, cameraConfig.Height)
	cam.frameEncoding = cameraConfig.Encoding

	//Capture
	err = cam.camH.StartStreaming()
	if err != nil {
		fmt.Println(err.Error())
		cam.camH.Close()
		return nil
	}
	return &cam
}

func (cam *Camera) Destroy() {
	cam.camH.Close()
}

func Detect() {
	idx := 0
	for {
		dev := fmt.Sprintf("/dev/video%d", idx)
		cam, err := webcam.Open(dev) // Open webcam
		if err != nil {
			break
		}
		defer cam.Close()
		fmt.Printf("Device: %s\n", dev)

		//List video formats
		formatDesc := cam.GetSupportedFormats()
		for format, encoding := range formatDesc {

			//For given video format, get frame sizes
			frames := cam.GetSupportedFrameSizes(format)

			for res := 0; res < len(frames); res++ {
				fmt.Printf("  Format: %d Encoding: %s Width: %4d Height: %4d\n", format, encoding, uint32(frames[res].MaxWidth), uint32(frames[res].MaxHeight))
			}
		}

		idx++
	}
	fmt.Printf("\n%d devices found.\n", idx)
}

//FrameGrabberGet returns the last fetched image with its timestamp
func (cam *Camera) FrameGrabberGet() (*image.Image, time.Time) {
	var img image.Image
	cam.lastLock.Lock()
	if cam.lastImg != nil {
		//Duplicate img here
		img = *cam.lastImg
	}
	ts := cam.lastImgTs
	cam.lastLock.Unlock()

	return &img, ts
}

//FrameGrabberStop stops the running frame grabber routine
func (cam *Camera) FrameGrabberStop() {
	cam.stopGrabber = true
}

//FrameGrabberStart starts frame grabber goroutine
func (cam *Camera) FrameGrabberStart() {
	cam.stopGrabber = false
	for {
		img := cam.GrabFrame()
		if img != nil {
			cam.lastLock.Lock()
			cam.lastImg = img
			cam.lastImgTs = time.Now()
			cam.lastLock.Unlock()
		}
		if cam.stopGrabber == true {
			break
		}
	}
}

func (cam *Camera) GrabFrameWithTimeout(timeout time.Duration) *image.Image {
	start := time.Now()
	for {
		frame := cam.GrabFrame()
		if frame != nil {
			return frame
		}
		if time.Since(start) > timeout {
			return nil
		}
	}
}

func (cam *Camera) GrabFrame() *image.Image {
	err := cam.camH.WaitForFrame(1)

	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		return nil
	default:
		return nil
	}

	frame, err := cam.camH.ReadFrame()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	switch cam.frameEncoding {
	case "MJPEG":
		//Decode JPEG
		rd := bytes.NewReader(frame)
		img, err := jpeg.Decode(rd)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return &img
	default:
		fmt.Println("Unknown encoding: " + cam.frameEncoding)
		return nil
	}
}
