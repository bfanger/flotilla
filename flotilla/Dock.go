package flotilla

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacobsa/go-serial/serial"
)

// Dock manages reading and writing to the hardware connected to the Flotilla Dock
type Dock struct {
	Ports      [8]Device
	connection io.ReadWriteCloser
}

// NewDock connect to a
func NewDock(tty string) (*Dock, error) {
	c, err := serial.Open(serial.OpenOptions{
		PortName:        tty,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 8,
	})
	if _, err = c.Write([]byte("v\r")); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("can't create dock: %v", err)
	}
	return &Dock{connection: c}, nil
}

// Close the serial connection to the dock
func (d *Dock) Close() error {
	return d.connection.Close()
}

// Listen to the serial connection
func (d *Dock) Listen() error {
	s := bufio.NewScanner(d.connection)
	for s.Scan() {
		line := string(s.Bytes())

		fmt.Println(line)
	}
	if s.Err() != nil {
		return s.Err()
	}
	return nil
}

// Pipe an input stream such as io.Stdin
func (d *Dock) Pipe(r io.Reader) error {
	s := bufio.NewScanner(r)
	eol := []byte("\r")[0]
	for s.Scan() {
		line := s.Bytes()
		line = append(line, eol)
		if _, err := d.Write(line); err != nil {
			return err
		}
	}
	if s.Err() != nil {
		return s.Err()
	}
	return nil
}

func (d *Dock) Write(p []byte) (int, error) {
	return d.connection.Write(p)
}

// Send the value to a port
func (d *Dock) Send(port int, value string) error {
	command := "s " + strconv.Itoa(port) + " " + value + "\r\n"
	fmt.Print(command)
	_, err := d.connection.Write([]byte(command))
	return err
}
