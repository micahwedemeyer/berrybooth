package main

import "fmt"

// #cgo CFLAGS: -I/Users/micah/Developer/include/gphoto2
// #cgo LDFLAGS: -L/Users/micah/Developer/lib -lgphoto2
// #include <gphoto2.h>
// #include <gphoto2/gphoto2-version.h>
import "C"
import "unsafe"
import "github.com/asaskevich/EventBus"

func handleNewImage(filename string) {
	fmt.Printf("New image: %s\n", filename)
}

func initEventBus() EventBus.EventBus {
	var bus = EventBus.New()
	bus.Subscribe("image:new", handleNewImage)

	return *bus
}

func setupGphoto2(bus EventBus.EventBus) {
	context := initCameraContext()
	camera, err := initCamera(context)

	if err < 0 {
		fmt.Printf(cameraErrToString(err))
	}

	bus.Publish("camera:init", camera, context)

	var abilities C.CameraAbilities
	C.gp_camera_get_abilities(camera, &abilities)

	var model = C.GoBytes(unsafe.Pointer(&abilities.model), 255)
	fmt.Printf("Model: %s\n", model)

	var text C.CameraText
	C.gp_camera_get_about(camera, &text, context)

	var s = C.GoBytes(unsafe.Pointer(&text.text), 255)
	fmt.Printf("Summary: %s\n", s)

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

	C.gp_camera_trigger_capture(camera, context)
}

func handleCameraSetup(camera *C.Camera, context *C.GPContext) {
	model, err := cameraModel(camera)
	if err >= 0 {
		fmt.Printf("Model: %s\n", model)
	}
}

func main() {
	var bus = initEventBus()
	bus.Subscribe("camera:init", handleCameraSetup)
	setupGphoto2(bus)
}
