package gphoto

import (
	"github.com/aqiank/go-gphoto2"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type Cam struct {
	context *gp.Context
	camera *gp.Camera
}

var (
	cam Cam
	cb chan []byte
)

func init(){
	initCamera()
	cb = make(chan []byte)
	go emitPreview(cb)
}

func CameraHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("New %s request\n", vars["type"])
	switch vars["type"] {
	case "preview":
		w.Write(<- cb)
	case "photo":
		w.Write(GetPhoto())
	default:
		return
	}
}

func emitPreview(c chan []byte) {
	for  {
		c <- GetPreview()
	}
}

func initCamera() {
	context := gp.NewContext()
	camera, err := gp.NewCamera()
	if err != nil {
		fmt.Printf("Error NewCamera: %s", err)
	}

	err = camera.Init(context)
	if err != nil {
		fmt.Printf("Error Init: %s", err)
	}
	cam.context = context
	cam.camera = camera
}

func Close() {
	err := cam.camera.Free()
	if err != nil {
		fmt.Printf("Error Free: %s", err)
	}
	cam.context.Free()
}

func GetPreview () []byte {
	img, err := cam.camera.CapturePreview(cam.context)
	if err != nil {
		fmt.Printf("Capture error: %s\n", err)
	}

	return img
}

func GetPhoto() []byte {
	img, err := cam.camera.CaptureImage(cam.context)
	if err != nil {
		fmt.Printf("Capture error: %s\n", err)
	}

	return img
}
