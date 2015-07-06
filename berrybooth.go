package main

import "fmt"

// #cgo CFLAGS: -I/Users/micah/Developer/include/gphoto2
// #cgo LDFLAGS: -L/Users/micah/Developer/lib -lgphoto2
// #include <gphoto2.h>
// #include <gphoto2/gphoto2-version.h>
import "C"
import "github.com/asaskevich/EventBus"
import "time"
import "github.com/micahwedemeyer/gphoto2go"
import "io"
import "os"

func handleNewImage(filename string) {
	fmt.Printf("New image: %s\n", filename)
}

func initEventBus() EventBus.EventBus {
	var bus = EventBus.New()
	bus.Subscribe("image:new", handleNewImage)

	return *bus
}

func setupGphoto2(bus EventBus.EventBus) *gphoto2.Camera {
	camera := new(gphoto2.Camera)
	err := camera.Init()

	if err < 0 {
		fmt.Printf(gphoto2.CameraResultToString(err))
	}

	bus.Publish("camera:init", camera)

	/*
		var rootFolder = C.CString("/")
		var folderList *C.CameraList
		C.gp_list_new(&folderList)

		var r = C.gp_camera_folder_list_folders(camera, rootFolder, folderList, context)
		if r < 0 {
			var err = C.GoString(C.gp_result_as_string(r))
			fmt.Printf("Error: %s", err)
		}
		var folderCount = int(C.gp_list_count(folderList))
		fmt.Printf("There are %d files at the root.\n", folderCount)

		//C.gp_camera_trigger_capture(camera, context)

		return camera, context
	*/
	return camera
}

func handleCameraSetup(camera *gphoto2.Camera) {
	model, err := camera.Model()
	if err >= 0 {
		fmt.Printf("Model: %s\n", model)
	}
}

func main() {
	var bus = initEventBus()
	bus.Subscribe("camera:init", handleCameraSetup)
	camera := setupGphoto2(bus)

	/*
		go func() {
			for {
				handler := func(eventType int, data string) {
					if eventType == C.GP_EVENT_FILE_ADDED {
						fmt.Printf("File added!\n")
						fmt.Printf("File: %s\n", data)
					}
				}
				camera.WaitForCameraEvent(1000, handler)
			}
		}()
		fmt.Printf("Waiting for event\n")
	*/

	/*
		err := camera.TriggerCapture()
		if err < 0 {
			fmt.Printf("Capture Error: %s\n", gphoto2.CameraResultToString(err))
		}
	*/

	time.Sleep(time.Duration(1) * time.Second)

	folders := camera.RListFolders("/")
	readers := make([]io.Reader, 20)
	for _, folder := range folders {
		fmt.Printf("Folder: %s\n", folder)

		files, _ := camera.ListFiles(folder)
		for _, fileName := range files {
			fmt.Printf("File: %s\n", folder+"/"+fileName)
			reader := camera.BufferedFileReader(folder, fileName)
			readers = append(readers, reader)

			fWriter, _ := os.Create("/Users/micah/tmp/" + fileName)
			io.Copy(fWriter, reader)
			fWriter.Close()

			fmt.Printf("File read/write complete\n")
		}
	}

	fmt.Printf("Done\n")
}
