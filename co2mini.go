package co2mini

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/zserge/hid"
)

type operation int

const (
	opCO2         operation = 0x50
	opTemperature operation = 0x42
)

var key = []byte{0x86, 0x41, 0xc9, 0xa8, 0x7f, 0x41, 0x3c, 0xac}

type Data struct {
	CO2         float64
	Temperature float64
}

type CO2Mini interface {
	GetCO2() (float64, error)
	GetTemperature() (float64, error)
	GetData() (*Data, error)
}

type co2Mini struct {
	device  hid.Device
	running bool
	data    Data
	mu      sync.RWMutex
}

func NewCO2Mini() (CO2Mini, error) {
	cm := &co2Mini{}

	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		if "04d9:a052:0100:00" != fmt.Sprintf(
			"%04x:%04x:%04x:%02x",
			info.Vendor,
			info.Product,
			info.Revision,
			info.Interface,
		) {
			return
		}
		cm.device = device
	})

	if cm.device == nil {
		return nil, errors.New("co2mini is not found")
	}

	if err := cm.device.Open(); err != nil {
		return nil, err
	}

	cm.running = true

	go func() {
		defer func() {
			cm.running = false
			cm.device.Close()
		}()

		if err := cm.device.SetReport(0, key); err != nil {
			return
		}

		for {
			encreptedData, err := cm.device.Read(-1, 5*time.Second) // CO2 and Temperature are updated every 15 seconds and 5 seconds respectively, so the timeout is set to 5 seconds.
			if err != nil {
				log.Println(err)
				continue
			}

			if len(encreptedData) != 8 {
				log.Println("the size of encrepted data is not 8")
				continue
			}

			plainData := decrypt(encreptedData, key)

			if !checksum(plainData) {
				log.Println("checksum error")
				continue
			}

			value := uint16(plainData[1])<<8 | uint16(plainData[2])

			switch op := operation(plainData[0]); op {
			case opCO2:
				cm.mu.RLock()
				cm.data.CO2 = float64(value)
				cm.mu.RUnlock()
			case opTemperature:
				cm.mu.RLock()
				cm.data.Temperature = math.RoundToEven(convertFtoC(float64(value))*10.0) / 10.0
				cm.mu.RUnlock()
			default:
				log.Println(encreptedData)
			}
		}
	}()

	return cm, nil
}

func (cm co2Mini) GetCO2() (float64, error) {
	data, err := cm.GetData()
	if err != nil {
		return 0, err
	}

	return data.CO2, nil
}

func (cm co2Mini) GetTemperature() (float64, error) {
	data, err := cm.GetData()
	if err != nil {
		return 0, err
	}

	return data.Temperature, nil
}

func (cm co2Mini) GetData() (*Data, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if !cm.running {
		return nil, errors.New("co2mini is not running")
	}

	return &Data{
		CO2:         cm.data.CO2,
		Temperature: cm.data.Temperature,
	}, nil
}
