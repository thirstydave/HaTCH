package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

var refTemp physic.Temperature
var refHumidity physic.RelativeHumidity
var refPressure physic.Pressure

func createEnv() {
	refTemp = 18
	// refHumidity = 999
	// refPressure = 999
}

func setTempTrigger(e *physic.Env, t physic.Temperature) {
	e.Temperature = t
}

func setHumidityTrigger(e *physic.Env, h physic.RelativeHumidity) {
	e.Humidity = h
}

func setPressureTrigger(e *physic.Env, p physic.Pressure) {
	e.Pressure = p
}

func monitorEnvironment(d *bmxx80.Dev, out chan int) {

	e := physic.Env{}
	for {
		if err := d.Sense(&e); err != nil {
			log.Fatal(err)
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {

	createEnv()
	refEnv := physic.Env{}
	setTempTrigger(&refEnv, refTemp)
	setHumidityTrigger(&refEnv, refHumidity)
	setPressureTrigger(&refEnv, refPressure)

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer b.Close()

	d, err := bmxx80.NewI2C(b, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}

	sensorData := make(chan int)
	// defer close(sensorData)
	go monitorEnvironment(d, sensorData)
	for {
		data := <-sensorData
		fmt.Printf("Temp %d\n", data)
	}
}
