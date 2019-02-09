package flotilla

// Device connected to the dock
type Device interface {
	Type() string
	Update(string) error
	Disconnect()
}
