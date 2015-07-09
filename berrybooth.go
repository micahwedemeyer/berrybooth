package main

// #cgo CFLAGS: -I/Users/micah/Developer/include/gphoto2
// #cgo LDFLAGS: -L/Users/micah/Developer/lib -lgphoto2
// #include <gphoto2.h>
// #include <gphoto2/gphoto2-version.h>
import "C"
import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/micahwedemeyer/gphoto2go"
	"io"
	"log"
	"os"
	"time"
)

var bus *EventBus.EventBus

func initLogger() {
	f, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(f)
}

func handleNewImage(filename string) {
	fmt.Printf("New image: %s\n", filename)
}

func initEventBus() {
	bus = EventBus.New()
	log.Printf("Eventbus Initialized\n")
}

func initCamera() *gphoto2go.Camera {
	camera := new(gphoto2go.Camera)
	err := camera.Init()

	if err < 0 {
		log.Fatalf("No camera found. Exiting.\n")
	} else {
		bus.Publish("camera:init", camera)
	}

	return camera
}

func initCameraEventSource(camera *gphoto2go.Camera) {
	go func() {
		for {
			eventChan := camera.AsyncWaitForEvent(1000)
			evt := <-eventChan
			if evt.Type == gphoto2go.EVENT_FILE_ADDED {
				bus.Publish("photo:capture", camera, evt.Folder, evt.File)
			}
		}
	}()
}

func handleCaptureEvent(camera *gphoto2go.Camera, folder string, fileName string) {
	path := "/Users/micah/tmp/" + fileName
	reader := camera.FileReader(folder, fileName)
	fWriter, _ := os.Create(path)
	io.Copy(fWriter, reader)
	fWriter.Close()
	reader.Close()
	bus.Publish("photo:saved", path)
}

func main() {
	initLogger()
	log.Printf("Berrybooth Startup\n")
	initEventBus()

	bus.Subscribe("camera:init", func(camera *gphoto2go.Camera) {
		model, _ := camera.Model()
		log.Printf("Detected camera: %s\n", model)
	})
	bus.Subscribe("camera:init", initCameraEventSource)
	bus.Subscribe("photo:capture", handleCaptureEvent)
	bus.Subscribe("photo:saved", func(path string) {
		log.Printf("File saved: %s\n", path)
	})

	initCamera()

	/*
		folders := camera.RListFolders("/")
		for _, folder := range folders {
			fmt.Printf("Folder: %s\n", folder)

			files, _ := camera.ListFiles(folder)
			for _, fileName := range files {
				fmt.Printf("File: %s\n", folder+"/"+fileName)

				reader := camera.FileReader(folder, fileName)
				fWriter, _ := os.Create("/Users/micah/tmp/" + fileName)
				io.Copy(fWriter, reader)
				fWriter.Close()
				reader.Close()
			}
		}
	*/

	for {
		time.Sleep(time.Duration(1) * time.Second)
	}
}
