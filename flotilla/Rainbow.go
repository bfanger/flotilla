package flotilla

import (
	"errors"
	"fmt"
	"strconv"
)

// Rainbow 5 Colors
type Rainbow struct {
	Colors [5]Color
}

// Type returns "motor"
func (mr *Rainbow) Type() string {
	return "rainbow"
}

// Set all pixels to the same color
func (r *Rainbow) Set(red, g, b uint8) {
	for i, _ := range r.Colors {
		r.Colors[i].R = red
		r.Colors[i].G = g
		r.Colors[i].B = b
	}
}

// Input errors
func (r *Rainbow) Input(value string) error {
	return errors.New("rainbow is an output device")
}
func (r *Rainbow) Output() (string, error) {
	return fmt.Sprintf("%s,%s,%s,%s,%s", r.Colors[0].String(), r.Colors[1].String(), r.Colors[2].String(), r.Colors[3].String(), r.Colors[4].String()), nil
}

// Disconnected rainbow
func (r *Rainbow) Disconnected() {
}

// Color a RGB value
type Color struct {
	R uint8
	G uint8
	B uint8
}

// Set a color
func (c *Color) Set(r, g, b uint8) {
	c.R = r
	c.G = g
	c.B = b
}

func (c *Color) String() string {
	return strconv.Itoa(int(c.R)) + "," + strconv.Itoa(int(c.G)) + "," + strconv.Itoa(int(c.B))
}
