package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	// Init Data
	initData := Emptyjson()
	initList := Emptylist()
	InitMain(initData)

	//rest api init controller
	getDataController := NewDeviceController(initData, initList)

	// gin init
	gin.SetMode(gin.ReleaseMode)

	//router init
	router := gin.Default()
	router.GET("/", getDataController.GetData)

	router.Run(":5000")
}
