package main

// Inverter JSON init
type Inverter struct {
	Voltage int     `json:"voltage"`
	Current float64 `json:"current"`
	Power   int     `json:"power"`
	Freq    int     `json:"freq"`
	PF      int     `json:"pf"`
	Temp    int     `json:"temp"`
	Quality int     `json:"quality"`
}

// Converter JSON init
type Converter struct {
	Voltage  float64 `json:"voltage"`
	Current1 float64 `json:"current1"`
	Current2 float64 `json:"current2"`
	Power    int     `json:"power"`
	Temp     int     `json:"temp"`
}

// Pack battery init
type Pack struct {
	Data    [14]float64 `json:"data"`
	Current float64     `json:"current"`
	Cap     int         `json:"cap"`
	Cycle   int         `json:"cycle"`
	Temp    int         `json:"temp1"`
	ID      int         `json:"id"`
}

// Harvester init
type Harvester struct {
	ConverterS   int       `json:"converter"`
	InverterS    int       `json:"inverter"`
	SwitchS      int       `json:"switch"`
	Consumption  int       `json:"consumption"`
	Production   int       `json:"production"`
	Storage      int       `json:"storage"`
	ProductionP  Converter `json:"e-production"`
	ConsumptionP Inverter  `json:"e-consumption"`
	StorageP     [5]Pack   `json:"e-harvester"`
}
