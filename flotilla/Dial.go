package flotilla

import (
	"errors"
	"strconv"
)

// Dial report the value of the Dial
type Dial struct {
	Raw uint16
	C   <-chan float64
	c   chan float64
}

// NewDial creates a new Dial
func NewDial() *Dial {
	c := make(chan float64, 1)
	return &Dial{C: c, c: c}
}

// Type returns "dial"
func (d *Dial) Type() string {
	return "dial"
}

// Input reads the value of the dail
func (d *Dial) Input(value string) error {
	number, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	d.Raw = uint16(number)
	d.c <- d.Value()
	return nil
}

// Output errors
func (d *Dial) Output() (string, error) {
	return "", errors.New("dial is an input device")
}

// Value normalized to a value between 0 and 1
func (d *Dial) Value() float64 {
	return float64(d.Raw) / 1023
}

// Disconnected closes all channels
func (d *Dial) Disconnected() {
	close(d.c)
}
