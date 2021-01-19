package main

import (
	"fmt"
	"log"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/pca9685"
	"periph.io/x/periph/host"
)

func lala(s string) {
	fmt.Printf(s)
}

func main() {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	var addr uint16 = 0x40
	pca, err := pca9685.NewI2C(bus, addr)
	if err != nil {
		log.Fatal(err)
	}

	if err := pca.SetPwmFreq(50 * physic.Hertz); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Done\n")

}
