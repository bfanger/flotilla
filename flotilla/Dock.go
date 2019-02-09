package flotilla

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

// Port 1 t/m 8
type Port int

// Dock manages reading and writing to the hardware connected to the Flotilla Dock
type Dock struct {
	Debug     bool
	Devices   [8]Device
	serial    *serial.Port
	onConnect []func(Port, Device)
}

// NewDock connect to a
func NewDock(tty string) (*Dock, error) {
	s, err := serial.OpenPort(&serial.Config{Name: tty, Baud: 115200})
	if err != nil {
		return nil, fmt.Errorf("can't create dock: %v", err)
	}
	return &Dock{serial: s}, nil
}

// Close the serial connection to the dock
func (d *Dock) Close() error {
	return d.serial.Close()
}

// Listen to the serial connection
func (d *Dock) Listen() error {
	if _, err := d.serial.Write([]byte("e\r")); err != nil {
		return err
	}
	d.serial.Flush()
	s := bufio.NewScanner(d.serial)
	for s.Scan() {
		line := s.Bytes()
		if d.Debug {
			fmt.Println(string(line))
		}
		r, err := d.parse(line)
		if err != nil {
			return err
		}
		device := d.Devices[r.Port-1]
		if (r.Connected || r.Update) && (device == nil || device.Type() != r.Device) {
			if device != nil {
				device.Disconnect()
				device = nil
			}
			switch r.Device {
			case "dial":
				device = NewDial(d)

			}
			if device != nil {
				d.Devices[r.Port-1] = device
				for _, callback := range d.onConnect {
					callback(r.Port, device)
				}
			}
		}
		if r.Update && device != nil {
			device.Update(r.Value)
		}
		if r.Disconnected && device != nil {
			device.Disconnect()
			d.Devices[r.Port-1] = nil
		}
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
	return d.serial.Write(p)
}

// Send the value to a port
func (d *Dock) Send(port int, value string) error {
	command := "s " + strconv.Itoa(port) + " " + value + "\r\n"
	if d.Debug {
		fmt.Print(command)
	}
	_, err := d.serial.Write([]byte(command))
	return err
}

// OnConnect is called when a (supported) device connects
func (d *Dock) OnConnect(callback func(Port, Device)) {
	d.onConnect = append(d.onConnect, callback)
}

type result struct {
	Comment      bool
	Update       bool
	Connected    bool
	Disconnected bool
	Port         Port
	Device       string
	Value        string
}

func (d *Dock) parse(line []byte) (*result, error) {
	r := &result{}
	switch string(line[0]) {
	case "#":
		r.Comment = true
	case "u":
		r.Update = true
	case "c":
		r.Connected = true
	case "d":
		r.Disconnected = true
	}
	if r.Connected || r.Disconnected || r.Update {
		port, err := strconv.Atoi(string(line[2]))
		if err != nil {
			return nil, err
		}
		r.Port = Port(port)
		r.Device = string(line[4:])
		if r.Update {
			pos := strings.Index(r.Device, " ")
			r.Value = r.Device[pos+1:]
			r.Device = r.Device[:pos]
		}
	}
	return r, nil
}
