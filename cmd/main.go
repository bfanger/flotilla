package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bfanger/flotilla/flotilla"
)

func main() {
	tty := "/dev/tty.usbmodem14101"
	dock, err := flotilla.NewDock(tty)
	if err != nil {
		log.Fatal(err)
	}
	defer dock.Close()
	// dock.Debug = true

	go func() {
		if err := dock.Pipe(os.Stdin); err != nil {
			log.Println(err)
		}
	}()

	dock.OnConnect(func(d flotilla.Device) {
		fmt.Printf("connected: %T\n", d)
		switch device := d.(type) {

		case *flotilla.Motor:
			speed := 0.1
			fmt.Printf("motor: %0.3f\n", speed)
			device.SetSpeed(speed)
			dock.Update(device)

		case *flotilla.Slider:
			device.OnChange(func(v float64) {
				fmt.Printf("slider: %0.3f\n", v)
			})

		case *flotilla.Dial:
			device.OnChange(func(v float64) {
				fmt.Printf("dial: %0.3f\n", v)
			})
		case *flotilla.Joystick:
			device.OnChange(func(x, y float64, pressed bool) {

				fmt.Printf("joystick: (x: %0.3f, y: %0.3f, b: %v)\n", x, y, pressed)
			})

		}
	})
	dock.OnDisconnect(func(d flotilla.Device) {
		fmt.Printf("disconnected: %T\n", d)
	})

	if err := dock.Listen(); err != nil {
		log.Fatal(err)
	}
}
