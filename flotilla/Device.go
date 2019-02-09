package flotilla

// Device represents hardware
type Device interface{}

// Receiver receives values from the hardware (Input/Sensor)
type Receiver interface {
	Receive(string) error
}

// Sender send values to the hardware (Output/Display)
type Sender interface {
	Send() (string, error)
}
