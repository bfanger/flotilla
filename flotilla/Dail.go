package flotilla

import "strconv"

// Dial report the value of the Dial
type Dial struct {
	Value        uint16
	dock         *Dock
	disconnected bool
	channels     []chan uint16
}

// NewDial creates a new Dial
func NewDial(d *Dock) *Dial {
	return &Dial{dock: d}
}

// Type returns "dial"
func (d *Dial) Type() string {
	return "dial"
}

// Update the value of the dail
func (d *Dial) Update(value string) error {
	number, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	d.Value = uint16(number)
	for _, c := range d.channels {
		c <- d.Value
	}
	return nil
}

// Disconnect closes all channels
func (d *Dial) Disconnect() {
	d.disconnected = true
	for _, c := range d.channels {
		close(c)
	}
}

// Values creates stream of value
func (d *Dial) Values() <-chan uint16 {
	c := make(chan uint16, 1)
	if d.disconnected {
		close(c)
		return c
	}
	c <- d.Value
	d.channels = append(d.channels, c)
	return c
}
