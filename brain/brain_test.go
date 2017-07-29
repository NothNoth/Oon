package brain_test

import (
	"Oon/brain"
	"Oon/camera"
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func regenerateImages() {
	camera.Detect()
	cam := camera.New("./camera/example.conf")
	img1 := cam.GrabFrameWithTimeout(5 * time.Second)
	time.Sleep(2 * time.Second)
	img2 := cam.GrabFrameWithTimeout(5 * time.Second)

	fo, _ := os.Create("./data/img1.jpg")
	w := bufio.NewWriter(fo)
	jpeg.Encode(w, *img1, nil)
	fo.Close()

	fo, _ = os.Create("./data/img2.jpg")
	w = bufio.NewWriter(fo)
	jpeg.Encode(w, *img2, nil)
	fo.Close()

}

func TestView(t *testing.T) {
	buf1, err := ioutil.ReadFile("./data/img1.jpg")
	if err != nil {
		panic("Test image not found")
	}
	buf2, err := ioutil.ReadFile("./data/img2.jpg")
	if err != nil {
		panic("Test image not found")
	}
	r := bytes.NewReader(buf1)
	img1, err := jpeg.Decode(r)
	if err != nil {
		panic("Test image not found")
	}

	r = bytes.NewReader(buf2)
	img2, err := jpeg.Decode(r)
	if err != nil {
		panic("Test image not found")
	}

	dif := brain.DiffFrame(&img1, &img1)
	if dif < 0.999 {
		t.Error("Auto image diff is 1.0")
	}

	dif = brain.DiffFrame(&img1, &img2)
	if dif < 0.75 {
		t.Error("Close pictures are close to 1.0")
	}
	fmt.Printf("diff %f\n", dif)
}
