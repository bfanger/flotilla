package flotilla

import (
	"errors"
	"strconv"
)

// Motor controls the speed of the Motor
type Motor struct {
	Raw int8
}

// Type returns "motor"
func (m *Motor) Type() string {
	return "motor"
}

// Input does nothing, is an output device
func (m *Motor) Input(value string) error {
	return errors.New("motor is an output device")
}

// Output errors
func (m *Motor) Output() (string, error) {
	return strconv.Itoa(int(m.Raw)), nil
}

// Disconnected motor
func (m *Motor) Disconnected() {
}

// Speed of the motor -1.0 to 1.0
func (m *Motor) Speed() float64 {
	return float64(m.Raw) / 63
}

// Forward speed of the motor
func (m *Motor) Forward(speed float64) {
	m.Raw = int8(speed * 63)
}
