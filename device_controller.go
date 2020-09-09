package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//DeviceController is interface to controller device
type DeviceController interface {
	GetData(ctx *gin.Context)
}

//NewDeviceController is function to make new implements struct from Interface
func NewDeviceController(j *Harvester) DeviceController {
	Impl := ImplDeviceController{
		harvester: j,
	}
	return &Impl
}

//ImplDeviceController is implement struct from DeviceController
type ImplDeviceController struct {
	harvester *Harvester
}

//GetData from MCU
func (d *ImplDeviceController) GetData(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, d.harvester)
	return
}

//GetPortList from Python
func (d *ImplDeviceController) GetPortList(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
	return
}

// GetControllSignal from MCU
func GetControllSignal(ctx *gin.Context) {

	return
}
