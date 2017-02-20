package Ports

import (
	"errors"

	"github.com/mikepb/go-serial"
)

//PortReader is the interface that wraps the basic Port Read method.
//It returns the data read as string and any error encountred.
type PortReader interface {
	Read() (n string, err error)
}

//Port represents the serial port.
type Port struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Options     serial.Options `json:"options"`
	Opened      bool `json:"opened"`
}

//Read reads the incoming buffer in a specific port
// and configuration. If config options is nill
// read with serial.RawOptions.
func (p *Port) Read() (string, error) {

	if p == nil {
		return "", errors.New("configs cannot be nill")
	}

	if p.Options == (serial.Options{}) {
		p.Options = serial.RawOptions
	}

	p.Options.BitRate = 115200
	port, err := p.Options.Open(p.Name)

	if err != nil {
		return "", err
	}

	defer port.Close()

	buf := make([]byte, 1)
	if _, err := port.Read(buf); err != nil {
		return "", err
	}
	return string(buf), nil

}

//List the serial ports available on the system
func List() ([]*Port, error){

	var result []*Port
	if ports, err := serial.ListPorts(); err != nil {
		return nil, err
	}else {
		for _, value := range ports {

			result = append(result,&Port{Name:value.Name(),Description:value.Description()})

		}

	}

	return result, nil
}
