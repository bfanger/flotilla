package flotilla

import (
	"strconv"
)

// Dial report the value of the Dial
type Dial struct {
	Raw       uint16
	listeners []func(float64)
}

// Receive the value from the hardware
func (d *Dial) Receive(value string) error {
	number, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	d.Raw = uint16(number)
	v := d.Value()
	for _, listener := range d.listeners {
		listener(v)
	}
	return nil
}

// Value normalized to a value between 0 and 1
func (d *Dial) Value() float64 {
	return float64(d.Raw) / 1023
}

// OnChange listen to changes in value
func (d *Dial) OnChange(callback func(float64)) {
	d.listeners = append(d.listeners, callback)
}
