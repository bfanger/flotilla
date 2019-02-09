package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bfanger/flotilla/flotilla"
)

func main() {
	tty := "/dev/tty.usbmodem14101"
	d, err := flotilla.NewDock(tty)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		d.Write([]byte("e\r\n"))

		// time.Sleep(50 * time.Millisecond)
		// // d.Send(1, "1,0,0")
		// // d.Send(1, "1,0,0,0,0,0,0,1,0,0,0,0,0,0,1")
		// r := flotilla.NewRainbow(1, d)
		// r.Set(1, 0, 0)
		// time.Sleep(50 * time.Millisecond)
		// r.Set(1, 1, 1)
		// time.Sleep(1 * time.Second)
		// r.Set(0, 0, 1)
	}()
	go func() {
		if err := d.Pipe(os.Stdin); err != nil {
			log.Println(err)
		}
		// os.Exit(0)
	}()
	d.OnConnect(func(_ flotilla.Port, d flotilla.Device) {
		switch device := d.(type) {
		case *flotilla.Dial:
			fmt.Printf("Dial: %+v\n", device)
			go func() {
				for value := range device.Values() {
					fmt.Println("value:", value)
				}
			}()
		default:
			fmt.Printf("Device connected: %+v\n", d)
		}
	})
	if err := d.Listen(); err != nil {
		log.Fatal(err)
	}
}
