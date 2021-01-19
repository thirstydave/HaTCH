package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/pca9685"
	"periph.io/x/periph/host"
)

var interval time.Duration = 1
var enable bool = false

// SetPcaDuty for PCA9865
func SetPcaDuty(pca *pca9685.Dev, i uint16) {
	if err := pca.SetPWMMan(0xFA, 0, i); err != nil {
		log.Fatal(err)
	}
}

func runPcaSample(pca *pca9685.Dev) {
	for true {
		for i := 0; i <= 1024; i++ {
			SetPcaDuty(pca, uint16(i))
			time.Sleep(2 * (interval * time.Millisecond))
		}
		for i := 1024; i >= 0; i-- {
			SetPcaDuty(pca, uint16(i))
			time.Sleep(interval * time.Millisecond)
		}
	}
}

// Stop stops the pin.
func Stop(pin gpio.PinIO) {
	if err := pin.Halt(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	pca, err := pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		log.Fatal(err)
	}
	pca.SetPwmFreq(250 * physic.Hertz)
	go runPcaSample(pca)

	// Use gpioreg GPIO pin registry to find a GPIO pin by name.
	pin := gpioreg.ByName("GPIO4")
	pin.Out(gpio.Low)

	// Consonle commands
	fmt.Printf("Hatching...\n")
	for i := 1; i == 1; i += 0 {
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		switch char {
		case 'o':
			fmt.Println("SLOWER")
			interval += 2
			break
		case 'p':
			fmt.Println("FASTER")
			if interval >= 3 {
				interval -= 2
			}
			break
		case 'q':
			fmt.Println("QUIT")
			i = 0
			break
		case '1':
			fmt.Println("OUTNE: 00b")
			pca.SetOUTNE(20)
			break
		case '2':
			fmt.Println("OUTNE: 01b")
			pca.SetOUTNE(21)
			break

		case '3':
			fmt.Println("OUTNE: 10b")
			pca.SetOUTNE(22)
			break

		case '4':
			fmt.Println("OUTNE: 11b")
			pca.SetOUTNE(23)
			break
		case 'z':
			if enable {
				pin.Out(gpio.Low)
				enable = false
				fmt.Println("GPIO LOW")
			} else {
				pin.Out(gpio.High)
				enable = true
				fmt.Println("GPIO HIGH")
			}
			break
		}
	}

	pin.Out(gpio.High)
	Stop(pin)
	time.Sleep(1 * time.Second)
	fmt.Println("Stopped")
}
