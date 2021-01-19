package main

import (
    "fmt"
    "log"

    "periph.io/x/periph/conn/i2c/i2creg"
    "periph.io/x/periph/devices/bmxx80"
    "periph.io/x/periph/host"
)


func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

    // Open a handle to the first available I²C bus:
    bus, err := i2creg.Open("")
    if err != nil {
        log.Fatal(err)
    }
    defer bus.Close()

    // Open a handle to a bme280/bmp280 connected on the I²C bus using default
    // settings:
    dev, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
    if err != nil {
        log.Fatal(err)
    }
    defer dev.Halt()
    fmt.Printf("%T", dev)

}
