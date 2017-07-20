package camera_test

import "testing"
import "Oon/camera"
import "time"

func TestDump(t *testing.T) {
	camera.Detect()
}

func TestInit(t *testing.T) {
	cam := camera.New("foo")
	if cam != nil {
		t.Error("Empty config is rejected")
	}

	cam = camera.New("example.conf")
	if cam == nil {
		t.Error("Valid config is allowed")
	}
	cam.Destroy()
}

func TestGrab(t *testing.T) {
	cam := camera.New("example.conf")

	frame := cam.GrabFrameWithTimeout(5 * time.Second)
	if frame == nil {
		t.Error("Proper config allows frame grabbing within 5s")
	}
}
