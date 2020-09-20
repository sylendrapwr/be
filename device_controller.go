package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//DeviceController is interface to controller device
type DeviceController interface {
	GetData(ctx *gin.Context)
	GetPortList(ctx *gin.Context)
}

//NewDeviceController is function to make new implements struct from Interface
func NewDeviceController(j *Harvester, l *PortList) DeviceController {
	Impl := ImplDeviceController{
		harvester: j,
		list:      l,
	}
	return &Impl
}

//ImplDeviceController is implement struct from DeviceController
type ImplDeviceController struct {
	harvester *Harvester
	list      *PortList
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

// PostData ...
func PostData(ctx *gin.Context) {
	signal := ctx.PostForm("turnoff")
	if signal == "1" {
		SetMcu("/dev/ttyUSB0")
	}
	ctx.Status(http.StatusOK)
	return
}
