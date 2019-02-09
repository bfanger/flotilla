package flotilla

import (
	"errors"
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

// Set all pixels to the same color
func (r *Rainbow) Set(red, green, blue uint8) error {
	if r.port == 0 {
		return errors.New("rainbow disconnected")
	}
	for _, c := range r.Colors {
		c.Red = red
		c.Green = green
		c.Blue = blue
	}
	return r.d.Send(r.port, r.Colors[0].String())
}

// Update does nothing, Rainbow is an output device
func (r *Rainbow) Update(value string) {
}

// Disconnect prevents future writes
func (r *Rainbow) Disconnect() {
	r.port = 0
}

// Color a RGB value
type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// Set a color
func (c *Color) Set(red, green, blue uint8) {
	c.Red = red
	c.Green = green
	c.Blue = blue
}

func (c *Color) String() string {
	return strconv.Itoa(int(c.Red)) + "," + strconv.Itoa(int(c.Green)) + "," + strconv.Itoa(int(c.Blue))
}
