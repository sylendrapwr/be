package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// GetToken ..
func GetToken() (token string, e error) {
	// LOGIN
	urlLogin := "http://103.25.209.137/login"
	type serialNumberStruct struct {
		SerialNumber int `json:"serial_number"`
	}
	var serialNumber serialNumberStruct
	type loginResp struct {
		Code int `json:"code"`
		Data struct {
			SerialNumber int    `json:"serial_number"`
			Token        string `json:"token"`
		} `json:"data"`
		Msg string `json:"message"`
	}
	var loginInformation loginResp
	content, err := ioutil.ReadFile(".sn.txt")
	if err != nil {
		e = err
		return
	}
	serialNumber.SerialNumber, err = strconv.Atoi(string(content))
	serialNumberJSON, err := json.Marshal(serialNumber)
	if err != nil {
		e = err
		return
	}
	resp, err := http.Post(urlLogin,
		"application/json", bytes.NewBuffer(serialNumberJSON))
	if err != nil {
		e = err
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e = err
		return
	}
	err = json.Unmarshal(body, &loginInformation)
	if err != nil {
		e = err
		return
	}

	// CEK TOKEN
	var cekTokenResp struct {
		Code int    `json:"code"`
		Data bool   `json:"data"`
		Msg  string `json:"message"`
	}
	urlCekToken := "http://103.25.209.137/auth/validate-token"
	var bearer = "Bearer " + loginInformation.Data.Token

	req, err := http.NewRequest("GET", urlCekToken, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		e = err
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &cekTokenResp)
	if err != nil {
		e = err
		return
	}
	if cekTokenResp.Data != true {
		e = errors.New("Salah token")
		return
	}
	e = nil
	token = loginInformation.Data.Token
	return
}

// GetDeviceInformation ...
func GetDeviceInformation(token string) (descrption string, e error) {
	type deviceInformationResp struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data struct {
			SerialNumber int    `json:"serial_number"`
			Desc         string `json:"description"`
			Name         string `json:"name"`
			BP           int    `json:"battery_pack"`
		}
	}
	var deviceInformation deviceInformationResp
	urlCekToken := "http://103.25.209.137/auth/device-information"
	var bearer = "Bearer " + token

	req, err := http.NewRequest("GET", urlCekToken, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e = err
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &deviceInformation)
	if err != nil {
		e = err
		return
	}
	descrption = deviceInformation.Data.Name
	e = nil
	return
}

// SendData ...
func SendData(hvs Harvester, token string) (e error) {
	type SendDataResp struct {
		Code int    `json:"code"`
		Msg  string `json:"success"`
		Data string `json:"data"`
	}
	var SendData SendDataResp
	url := "http://103.25.209.137/device-data/relay-endpoint"
	var bearer = "Bearer " + token

	hvsJSON, _ := json.Marshal(hvs)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(hvsJSON))
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e = err
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &SendData)
	if err != nil {
		e = err
		return
	}
	fmt.Println(string(body))
	if SendData.Code != 0 {
		e = errors.New("Gagal")
		return
	}
	e = nil
	return
}
