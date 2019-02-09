package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bfanger/flotilla/flotilla"
)

func main() {
	tty := "/dev/tty.usbmodem14201"
	dock, err := flotilla.NewDock(tty)
	if err != nil {
		log.Fatal(err)
	}
	defer dock.Close()
	dock.Debug = true
	// go func() {
	// 	dock.Send("s 3 1,0,0,0,0,0,0,1,0,0,0,0,0,0,1")
	// }()
	go func() {
		if err := dock.Pipe(os.Stdin); err != nil {
			log.Println(err)
		}
		// os.Exit(0)
	}()
	var motor *flotilla.Motor
	var rainbow *flotilla.Rainbow
	dock.OnConnect(func(d flotilla.Device) {
		switch device := d.(type) {
		case *flotilla.Dial:
			fmt.Printf("Dial: %+v\n", device)
			go func() {
				for value := range device.C {
					// fmt.Printf("value: %.2f\n", value)
					if motor != nil {
						motor.Forward(value)
						if err := dock.Update(motor); err != nil {
							log.Fatal(err)
						}
					}
					if rainbow != nil {
						v := uint8(value * 20)
						rainbow.Colors[0].Set(v, v, 0)
						// cmd, err = rainbow.
						if err := dock.Update(rainbow); err != nil {
							log.Fatal(err)
						}
					}
				}
			}()
		case *flotilla.Motor:
			fmt.Printf("Motor: %+v\n", device)
			motor = device
		case *flotilla.Rainbow:
			fmt.Printf("Rainbow: %+v\n", device)
			rainbow = device
			dock.Update(rainbow)
		default:
			fmt.Printf("Device connected: %+v\n", device)
		}
	})
	dock.OnDisconnect(func(device flotilla.Device) {
		fmt.Printf("Device disconnected: %+v\n", device)
		if motor == device {
			motor = nil
		}
		if rainbow == device {
			rainbow = nil
		}
	})
	if err := dock.Listen(); err != nil {
		log.Fatal(err)
	}
}
