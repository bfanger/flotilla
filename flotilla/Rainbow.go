package flotilla

import (
	"strconv"
)

// Rainbow 5 Colors
type Rainbow struct {
	Colors [5]*Color
	port   int
	d      *Dock
}

// NewRainbow creates a new rainbow device
func NewRainbow(port int, d *Dock) *Rainbow {
	r := &Rainbow{port: port, d: d}
	for i := range r.Colors {
		r.Colors[i] = &Color{}
	}
	return r
}

func (r *Rainbow) Set(red, green, blue uint8) error {
	for _, c := range r.Colors {
		c.Red = red
		c.Green = green
		c.Blue = blue
	}
	return r.d.Send(r.port, r.Colors[0].String())
}
func (r *Rainbow) Flush() error {
	// @todo Update the leds based on all Color values
	return nil
}

// Update
func (r *Rainbow) Update(value string) {
	// implement Device interface
}

func (r *Rainbow) Disconnect() {
	r.port = 0
}

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (c *Color) Set(red, green, blue uint8) {
	c.Red = red
	c.Green = green
	c.Blue = blue
}

func (c *Color) String() string {
	return strconv.Itoa(int(c.Red)) + "," + strconv.Itoa(int(c.Green)) + "," + strconv.Itoa(int(c.Blue))
}
