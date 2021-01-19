	package main

import (
    "fmt"
    "log"
    "time"

    "periph.io/x/periph/conn/i2c/i2creg"
    "periph.io/x/periph/conn/physic"
    "periph.io/x/periph/devices/bmxx80"
    "periph.io/x/periph/host"
)

func readSensor(dev *bmxx80.Dev) {
for {
    // Read temperature from the sensor:
    var env physic.Env
    if err := dev.Sense(&env); err != nil {
        log.Fatal(err)
    }

        fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
        time.Sleep(1000 * time.Millisecond)
    }
}

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
   
    readSensor(dev)
    defer dev.Halt()
}
