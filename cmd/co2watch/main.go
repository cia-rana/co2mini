package main

import (
	"log"
	"time"

	"github.com/cia-rana/co2mini"
)

func main() {
	cm, err := co2mini.NewCO2Mini()
	if err != nil {
		log.Println(err)
		return
	}

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			data, err := cm.GetData()
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("CO2: %.0f, Temperature: %.1f", data.CO2, data.Temperature)
		}
	}
}
