package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/albenik/go-serial/v2"
	"github.com/goburrow/modbus"
)

// Emptytoken ....
func Emptytoken() *Token {
	T := Token{Number: ""}
	return &T
}

// Emptylist ....
func Emptylist() *PortList {
	L := PortList{List: ""}
	return &L
}

// Emptyjson .....
func Emptyjson() *Harvester {
	S := 0
	I := 0
	C := 0
	PP := Converter{Voltage: 0, Current1: 0, Current2: 0, Power: 0, Temp: 0}
	IP := Inverter{Voltage: 0, Current: 0, Power: 0, Temp: 0, PF: 0, Freq: 0, Quality: 0}
	SP := [5]Pack{}
	for i := 0; i < 5; i++ {
		SP[i] = Pack{ID: i + 1, Cycle: 0, Cap: 0, Current: 0, Temp: 0,
			Data: [14]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}}
	}
	H := Harvester{
		Consumption:  I,
		Storage:      S,
		Production:   C,
		ProductionP:  PP,
		StorageP:     [5]Pack{SP[0], SP[1], SP[2], SP[3], SP[4]},
		ConsumptionP: IP,
	}
	return &H
}

// GetPack ........
func GetPack(id int, port string) (s Pack, vt float64, e error) {
	s = Pack{ID: id, Cycle: 0, Cap: 0, Current: 0, Temp: 0,
		Data: [14]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}}
	p, err := serial.Open(port,
		serial.WithBaudrate(9600),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(2000),
		serial.WithWriteTimeout(1000),
	)
	defer p.Close()
	if err != nil {
		e = err
		return
	}

	// Voltage
	_, err = p.Write([]byte("\xa5\x80\x95\x08\x00\x00\x00\x00\x00\x00\x00\x00\xc2"))
	if err != nil {
		e = err
		return
	}
	buf := make([]byte, 208)
	n, err := p.Read(buf)
	if err != nil {
		e = err
		return
	}
	if n == 0 {
		e = errors.New("EOF")
		return
	}
	for i := 0; i < 5; i++ {
		vsum := 0
		for emptyJSON := 0; emptyJSON < 12; emptyJSON++ {
			vsum = vsum + int(buf[13*i+emptyJSON])
		}
		if int(buf[13*i+12]) != vsum%256 {
			e = errors.New("Wrong checksum")
			return
		}
	}
	var v [14]float64
	vlist := []int{5, 7, 9, 18, 20, 22, 31, 33, 35, 44, 46, 48, 57, 59}
	for i := 0; i < 14; i++ {
		v[i] = float64((float64(buf[vlist[i-0]])*256 + float64(buf[vlist[i]])) / 1000)
	}
	s.Data = v

	// Current
	_, err = p.Write([]byte("\xa5\x80\x90\x08\x00\x00\x00\x00\x00\x00\x00\x00\xbd"))
	if err != nil {
		e = err
		return
	}
	n, err = p.Read(buf)
	if err != nil {
		e = err
		return
	}
	if n == 0 {
		e = errors.New("EOF")
		return
	}
	csum := 0
	for emptyJSON := 0; emptyJSON < 12; emptyJSON++ {
		csum = csum + int(buf[emptyJSON])
	}
	if int(buf[12]) != csum%256 {
		e = errors.New("Wrong checksum")
		return
	}
	c := float64(buf[8])*256 + float64(buf[9])
	s.Current = int((c - 30000.0) * 0.1)

	// Temperature
	_, err = p.Write([]byte("\xa5\x80\x96\x08\x00\x00\x00\x00\x00\x00\x00\x00\xc3"))
	if err != nil {
		e = err
		return
	}
	n, err = p.Read(buf)
	if err != nil {
		e = err
		return
	}
	if n == 0 {
		e = errors.New("EOF")
		return
	}
	for i := 0; i < 2; i++ {
		tsum := 0
		for emptyJSON := 0; emptyJSON < 12; emptyJSON++ {
			tsum = tsum + int(buf[13*i+emptyJSON])
		}
		if int(buf[13*i+12]) != tsum%256 {
			e = errors.New("Wrong checksum")
			return
		}
	}
	s.Temp = int(buf[5]) - 40
	var cap float64
	for i := 0; i < 14; i++ {
		cap = cap + v[i]
	}
	s.Cap = int(((cap - 45.0) / 13.8) * 3048.0)
	vt = (float64(s.Cap) * 100) / 3048
	e = nil
	return
}

// SetMcu .....
func SetMcu(port string) (e error) {
	p, err := serial.Open(port,
		serial.WithBaudrate(9600),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(2000),
		serial.WithWriteTimeout(1000),
	)
	defer p.Close()
	if err != nil {
		e = err
		return
	}

	// Voltage
	_, err = p.Write([]byte("b"))
	if err != nil {
		e = err
		return
	}
	e = nil
	return
}

func splitValue(value string) (result []string, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = errors.New("cannot split value")
			return
		}
	}()
	e = nil
	result = strings.Split(value, "#")
	return
}

// Block ...
type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

// Exception ...
type Exception interface{}

// Throw ....
func Throw(up Exception) {
	panic(up)
}

// Do ....
func (tcf Block) Do() {
	if tcf.Finally != nil {
		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

// NewSMBUS Factory method for SMBus
func NewSMBUS(bus uint, address byte) (*SMBus, error) {
	smb := &SMBus{bus: nil}
	err := smb.BusOpen(bus)
	if err != nil {
		return nil, err
	}
	err = smb.SetAddr(address)
	if err != nil {
		return nil, err
	}
	return smb, nil
}

// GetMcuI2c ...
func GetMcuI2c() (c Converter, e error) {
	con, err := NewSMBUS(1, 0x04)
	if err != nil {
		e = err
		os.Exit(1)
		return
	}

	var result []byte
	num, err := con.ReadI2cBlockData(0, result)
	if err != nil {
		e = err
		return
	}
	fmt.Println(result)
	if num != 6 {
		e = errors.New("data corrupt")
		return
	}

	voltage := float64(int(result[0])*255+int(result[1])) / 100
	current1 := float64((int(result[2])*255+int(result[3]))-510) * 0.122070313
	current2 := float64((int(result[4])*255+int(result[5]))-510) * 0.122070313
	power := voltage * (current1 + current2)

	c.Voltage = voltage
	c.Current1 = current1
	c.Current2 = current2
	c.Power = int(power)

	fmt.Println(voltage)
	fmt.Println(current1)
	fmt.Println(current2)
	fmt.Println(power)

	e = nil
	return
}

// GetMcu ......
func GetMcu(port string) (c Converter, e error) {
	// Converter
	c = Converter{Voltage: 0.0, Current1: 0.0, Current2: 0.0, Power: 0, Temp: 0}
	p, err := serial.Open(port,
		serial.WithBaudrate(9600),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(2000),
		serial.WithWriteTimeout(1000),
	)
	defer p.Close()
	if err != nil {
		e = err
		return
	}

	// Voltage
	_, err = p.Write([]byte("a"))
	if err != nil {
		e = err
		return
	}
	buf := make([]byte, 50)
	n, err := p.Read(buf)
	if err != nil {
		e = err
		return
	}
	if n == 0 {
		e = errors.New("EOF")
		return
	}

	data := string(buf)
	fmt.Println(data)
	splittedData := strings.Split(data, "#")
	if len(splittedData) != 7 {
		e = errors.New("data error")
		return
	}
	power, err := strconv.Atoi(splittedData[5])
	if err != nil {
		e = err
		return
	}
	voltage, err := strconv.ParseFloat(splittedData[2], 64)
	if err != nil {
		e = err
		return
	}
	current1, err := strconv.ParseFloat(splittedData[3], 64)
	if err != nil {
		e = err
		return
	}
	current2, err := strconv.ParseFloat(splittedData[4], 64)
	if err != nil {
		e = err
		return
	}
	c.Voltage = voltage
	c.Power = power
	c.Current1 = current1
	c.Current2 = current2

	e = nil
	return
}

// GetPzem .......
func GetPzem(port string) (i Inverter, e error) {
	// Inverter
	i = Inverter{Voltage: 0, Current: 0, Power: 0, Temp: 0, PF: 0, Freq: 0, Quality: 0}
	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 1 * time.Second
	err := handler.Connect()
	if err != nil {
		e = err
		return
	}
	defer handler.Close()
	c := modbus.NewClient(handler)
	rbuf, err := c.ReadInputRegisters(0, 9)
	if err != nil {
		e = err
		return
	}
	var r [9]int
	for i := 0; i < 9; i++ {
		r[i] = int(rbuf[i*2])*256 + int(rbuf[i*2+1])
	}
	i.Voltage = r[0] / 10
	i.Current = float64((r[1] + r[2]*256) / 1000)
	i.Power = (r[3] + r[4]*256) / 10
	i.Freq = r[7] / 10
	pf := r[8] / 100
	i.PF = pf
	ql := 0
	if pf == 0 {
		pf = 1
	}
	if pf == 100 {
		pf = 99
	}
	if pf >= 80 {
		ql = int(pf / 20)
	}
	i.Quality = ql
	e = nil
	return
}

// GetPortList ....
func GetPortList() (s []string, e error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		e = err
		return
	}
	if len(ports) <= 1 {
		e = errors.New("BMS not found")
	}

	s = ports
	e = nil
	return
}

// InitMain ......
func InitMain(emptyJSON *Harvester) {
	go func() {
		var lock sync.Mutex
		timer := time.NewTicker(time.Second)
		defer timer.Stop()
		timer2 := time.NewTicker(time.Millisecond * 200)
		defer timer2.Stop()
		for {
			select {
			case <-timer.C:
				go func() {
					lock.Lock()
					defer lock.Unlock()

					// Consumption
					con, err := GetPzem("/dev/ttyUSB0")
					if err == nil {
						emptyJSON.ConsumptionP = con
						emptyJSON.Consumption = con.Power
					}
				}()

			case <-timer2.C:
				go func() {
					lock.Lock()
					defer lock.Unlock()
					pro, err := GetMcuI2c()

					// Production
					fmt.Println(err)
					if err == nil {
						emptyJSON.ProductionP = pro
						emptyJSON.Production = pro.Power
					}
				}()
			}
		}
	}()

	go func() {
		persen := 0.0 // persentase daya tersimpan
		count := 0    // cycle count

		portList, err := GetPortList()
		if err == nil {
			for i := 1; i < len(portList); i++ {
				str, vt, err := GetPack(i, portList[i])
				if err == nil {
					count++
					persen = persen + vt
					emptyJSON.StorageP[i-1] = str
					emptyJSON.Storage = int(persen / float64(count))
				}
			}
		}
	}()
}
