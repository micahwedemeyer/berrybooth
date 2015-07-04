package main

// #cgo CFLAGS: -I/Users/micah/Developer/include/gphoto2
// #cgo LDFLAGS: -L/Users/micah/Developer/lib -lgphoto2
// #include <gphoto2.h>
// #include <gphoto2/gphoto2-version.h>
import "C"
import "unsafe"

func initCameraContext() *C.GPContext {
	context := C.gp_context_new()
	return context
}

func initCamera(context *C.GPContext) (*C.Camera, int) {
	var camera *C.Camera

	C.gp_camera_new(&camera)
	err := C.gp_camera_init(camera, context)

	return camera, int(err)
}

func cameraAbilities(camera *C.Camera) (C.CameraAbilities, int) {
	var abilities C.CameraAbilities
	err := C.gp_camera_get_abilities(camera, &abilities)
	return abilities, int(err)
}

func cameraModel(camera *C.Camera) (string, int) {
	abilities, err := cameraAbilities(camera)
	modelBytes := C.GoBytes(unsafe.Pointer(&abilities.model), 255)
	model := string(modelBytes[:255])

	return model, err
}

func isCameraConnected(*C.Camera, *C.GPContext) bool {
	return false
}

func cameraErrToString(err int) string {
	return C.GoString(C.gp_result_as_string(C.int(err)))
}
