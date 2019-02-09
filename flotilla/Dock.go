package flotilla

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/tarm/serial"
)

// Port 1 t/m 8
type Port int

// Dock manages reading and writing to the hardware connected to the Flotilla Dock
type Dock struct {
	Debug        bool
	Inputs       [9]Receiver
	Outputs      [9]Sender
	serial       *serial.Port
	onConnect    []func(Device)
	onDisconnect []func(Device)
	mutex        sync.RWMutex
}

// NewDock opens a serial connection to the dock
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
		if d.Debug {
			fmt.Println(string(s.Bytes()))
		}
		line, err := d.parse(s.Bytes())
		if err != nil {
			return err
		}

		if line.Connected || line.Disconnected {
			if d.Inputs[line.Port] != nil {
				for _, listener := range d.onDisconnect {
					listener(d.Inputs[line.Port])
				}
			}
			d.Inputs[line.Port] = nil
			if d.Outputs[line.Port] != nil {
				for _, listener := range d.onDisconnect {
					listener(d.Outputs[line.Port])
				}
			}
			d.Outputs[line.Port] = nil

			if line.Connected {
				var device Device
				switch line.Device {
				case "motor":
					device = &Motor{}
				case "slider":
					device = &Slider{}
				case "dial":
					device = &Dial{}
				case "rainbow":
					device = &Rainbow{}
				}
				if device != nil {
					input, ok := device.(Receiver)
					if ok {
						d.Inputs[line.Port] = input
					}
					output, ok := device.(Sender)
					if ok {
						d.Outputs[line.Port] = output
					}
					for _, callback := range d.onConnect {
						callback(device)
					}
				}
			}
		}

		if line.Update && d.Inputs[line.Port] != nil {
			if err := d.Inputs[line.Port].Receive(line.Value); err != nil {
				return err
			}
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
		d.mutex.Lock()
		_, err := d.Write(line)
		d.mutex.Unlock()
		if err != nil {
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

// Update the hardware based on the Device values
func (d *Dock) Update(s Sender) error {
	var port Port
	for i, o := range d.Outputs {
		if s == o {
			port = Port(i)
		}
	}
	if port == 0 {
		return errors.New("output not connected")
	}
	value, err := s.Send()
	if err != nil {
		return err
	}
	command := "s " + strconv.Itoa(int(port)) + " " + value

	return d.Send(command)
}

// Send the command over the serial connection
func (d *Dock) Send(command string) error {
	if d.Debug {
		fmt.Println(command)
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, err := d.serial.Write([]byte(command + "\r")); err != nil {
		return err
	}
	return d.serial.Flush()
}

// Port of a device
func (d *Dock) Port(device Device) (Port, error) {
	if device == nil {
		return 0, errors.New("no device given")
	}
	for i, input := range d.Inputs {
		if input == device {
			return Port(i), nil
		}
	}
	for i, input := range d.Outputs {
		if input == device {
			return Port(i), nil
		}
	}
	return 0, errors.New("not connected")
}

// OnConnect is called when a (supported) device connects
func (d *Dock) OnConnect(callback func(Device)) {
	d.onConnect = append(d.onConnect, callback)
}

// OnDisconnect is called when a (supported) device is disconnected
func (d *Dock) OnDisconnect(callback func(Device)) {
	d.onDisconnect = append(d.onDisconnect, callback)
}

type line struct {
	Comment      bool
	Update       bool
	Connected    bool
	Disconnected bool
	Port         Port
	Device       string
	Value        string
}

func (d *Dock) parse(text []byte) (*line, error) {
	l := &line{}
	switch string(text[0]) {
	case "#":
		l.Comment = true
	case "u":
		l.Update = true
	case "c":
		l.Connected = true
	case "d":
		l.Disconnected = true
	}
	if l.Connected || l.Disconnected || l.Update {
		port, err := strconv.Atoi(string(text[2]))
		if err != nil {
			return nil, err
		}
		l.Port = Port(port)
		l.Device = string(text[4:])
		if l.Update {
			pos := strings.Index(l.Device, " ")
			l.Value = l.Device[pos+1:]
			l.Device = l.Device[:pos]
		}
	}
	return l, nil
}
