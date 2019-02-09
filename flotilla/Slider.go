package flotilla

import (
	"strconv"
)

// Slider report the value of the Slider
type Slider struct {
	Raw       uint16
	listeners []func(float64)
}

// Receive the value from the hardware
func (s *Slider) Receive(value string) error {
	number, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	s.Raw = uint16(number)
	v := s.Value()
	for _, listener := range s.listeners {
		listener(v)
	}
	return nil
}

// Value normalized to a value between 0 and 1
func (s *Slider) Value() float64 {
	return float64(s.Raw) / 1023
}

// OnChange listen to changes in value
func (s *Slider) OnChange(callback func(float64)) {
	s.listeners = append(s.listeners, callback)
}
