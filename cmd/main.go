package main

import (
	"log"
	"os"
	"time"

	"github.com/bfanger/flotilla/flotilla"
)

func main() {
	tty := "/dev/tty.usbmodem14201"
	d, err := flotilla.NewDock(tty)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		time.Sleep(50 * time.Millisecond)
		// d.Send(1, "1,0,0")
		// d.Send(1, "1,0,0,0,0,0,0,1,0,0,0,0,0,0,1")
		r := flotilla.NewRainbow(1, d)
		r.Set(1, 0, 0)
		time.Sleep(50 * time.Millisecond)
		r.Set(1, 1, 1)
		time.Sleep(1 * time.Second)
		r.Set(0, 0, 1)
	}()
	go func() {
		if err := d.Pipe(os.Stdin); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	if err := d.Listen(); err != nil {
		log.Fatal(err)
	}
}
