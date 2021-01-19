package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

const ()

var d *i2c.Dev

func newDev(addr uint64, b i2c.Bus) {
	d = &i2c.Dev{Addr: uint16(addr), Bus: b}
}

func main() {

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open the first available I2C bus
	fmt.Print("Opening first available I2C bus\n")

	b, err := i2creg.Open("")
	if err != nil {
		fmt.Printf("Failed to open: %v\n", err)
	}
	defer b.Close()

	fmt.Printf("Ready to bang the i2c...\n\n")
	fmt.Printf("[d]etect i2c devices\n")
	fmt.Printf("[a]ddress device\n")
	fmt.Printf("[q]uit\n")
	fmt.Printf("<->")

	for i := 1; i == 1; i += 0 {
		reader := bufio.NewReader(os.Stdin)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		switch cmd {
		case "d", "detect":

			break
		case "a", "address":
			fmt.Printf("ADDRESS(decimal): ")
			//Wait for an address
			in, _ := reader.ReadString('\n')
			in = strings.Replace(in, "\n", "", -1)
			//Receive it as a UInt 64 and cast to UInt16
			addr, _ := strconv.ParseUint(in, 10, 16)
			newDev(addr, b)
			fmt.Printf("Talking to %s\n", d.String())
			fmt.Printf("[w]rite\n")
			fmt.Printf("[r]ead\n")
			fmt.Printf("<%s>", d.String())
			break
		case "r", "read":
			fmt.Println("\nREAD MODE: ")
			fmt.Printf("<%s>[R] ", d.String())
			break
		case "w", "write":
			fmt.Println("\nWRITE MODE: [byte,byte]")
			fmt.Printf("<%s>[W] ", d.String())
			//Wait for an address
			in, _ := reader.ReadString('\n')
			in = strings.Replace(in, "\n", "", -1)
			ins := strings.Split(in, ",")
			//Receive it as an Int  and cast to byte
			ins0, _ := strconv.ParseInt(ins[0], 10, 8)
			ins1, _ := strconv.ParseInt(ins[1], 10, 8)
			var byte0 byte = byte(ins0)
			var byte1 byte = byte(ins1)
			bytes := []byte{byte0, byte1}
			//Try to call it
			d.Write(bytes)
			fmt.Printf("sent %v", bytes)
			fmt.Printf("<%s>[W] ", d.String())
			break
		case "q", "quit":
			fmt.Println("QUIT")
			i = 0
			b.Close()
			break
		}
	}

}
