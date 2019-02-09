package flotilla

import "strconv"

// Motor controls the speed of the Motor
type Motor struct {
	Raw int8
}

// Send speed to motor
func (m *Motor) Send() (string, error) {
	return strconv.Itoa(int(m.Raw)), nil
}

// Speed of the motor -1.0 to 1.0
func (m *Motor) Speed() float64 {
	return float64(m.Raw) / 63
}

// SetSpeed of the motor -1.0 to 1.0
func (m *Motor) SetSpeed(v float64) {
	m.Raw = int8(v * 63)
}
