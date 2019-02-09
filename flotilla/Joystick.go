package flotilla

import (
	"strconv"
	"strings"
)

// Joystick reports the position of the stick
type Joystick struct {
	Raw struct {
		Button bool
		X      int
		Y      int
	}
	listeners []func(x, y float64, pressed bool)
}

// Receive the value from the hardware
func (j *Joystick) Receive(value string) error {
	parts := strings.Split(value, ",")
	j.Raw.Button = parts[0] == "1"
	var err error
	j.Raw.X, err = strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	j.Raw.Y, err = strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	for _, listener := range j.listeners {
		// @todo Deadzone & Button events
		x := float64(j.Raw.X)/512 - 1
		y := float64(j.Raw.Y)/512 - 1
		listener(x, y, j.Raw.Button)
	}
	return nil
}

// OnChange listen to joystick movements
func (j *Joystick) OnChange(callback func(x, y float64, pressed bool)) {
	j.listeners = append(j.listeners, callback)
}
