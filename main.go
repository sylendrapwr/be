package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init Data
	initData := Emptyjson()
	InitMain(initData)

	//rest api init controller
	getDataController := NewDeviceController(initData)

	// gin init
	gin.SetMode(gin.ReleaseMode)
	token, err := GetToken()
	if err != nil {
		fmt.Println(err)
	}
	desc, err := GetDeviceInformation(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(desc)
	err = SendData(*initData, token)
	if err != nil {
		fmt.Println(err)
	}
	//router init
	router := gin.Default()
	router.GET("/", getDataController.GetData)
	router.POST("/", GetControllSignal)

	router.Run(":5000")
}
