package flotilla

// Device connected to the dock
type Device interface {
	Type() string
	Input(string) error
	Output() (string, error)
	Disconnected()
}
