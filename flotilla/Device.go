package flotilla

// Device connected to the dock
type Device interface {
	Update(string) error
	Disconnect()
}
