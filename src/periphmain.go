package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func lala(s string) {
	fmt.Printf(s)
}

func readtemp(c chan float64, d *i2c.Dev) {
	fmt.Printf("Beginning readtemp go routine")
	write := []byte{0x0} //0x0-Request Ambient Temperature Register
	read := make([]byte, 2)
	for true {
		if err := d.Tx(write, read); err != nil {
			log.Fatal(err)
		}
		temperature := ((uint16(read[0]) << 8) | uint16(read[1]))
		temperature = temperature >> 4
		tempC := 1.0 * float64(temperature) * 0.0625
		c <- tempC
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func notmain() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Enumerate all I2C buses available and the corresponding pins.
	fmt.Print("I2C buses available:\n")
	for _, ref := range i2creg.All() {
		fmt.Printf("- %s\n", ref.Name)
		if ref.Number != -1 {
			fmt.Printf("  %d\n", ref.Number)
		}
		if len(ref.Aliases) != 0 {
			fmt.Printf("  %s\n", strings.Join(ref.Aliases, " "))
		}

		b, err := ref.Open()
		if err != nil {
			fmt.Printf("  Failed to open: %v\n", err)
		}
		if p, ok := b.(i2c.Pins); ok {
			fmt.Printf("  SDA: %s\n", p.SDA())
			fmt.Printf("  SCL: %s\n", p.SCL())
		}
		d := &i2c.Dev{Addr: 72, Bus: b}
		fmt.Println(d.String())
		config := []byte{0x01, 0x64} // 0x01-Address to Configuration register // 0x64-ADC resolution: 12-bit,
		d.Write(config)

		tempchan := make(chan float64)
		go readtemp(tempchan, d)
		for i := range tempchan {
			fmt.Printf("Temp: %v\u2103", i)
		}

	}
}
