package Ports

import (
	"errors"

	"github.com/mikepb/go-serial"
)

//PortReader is the interface that wraps the basic Port Read method.
//It returns the data read as string and any error encountred.
type PortReader interface {
	Read(port *serial.Port) (string, error)
}

//Port represents the serial port.
type Port struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Options     serial.Options `json:"options"`
	Opened      bool `json:"opened"`
}

//Open open a port and return and instance of serial.Port
func (p *Port) Open()(*serial.Port,error) {

	if p == nil {
		return nil, errors.New("configs cannot be nill")
	}

	if p.Options == (serial.Options{}) {
		p.Options = serial.RawOptions
	}

	p.Options.BitRate = 115200
	port, err := p.Options.Open(p.Name)

	if err != nil {
		return nil, err
	}

	defer p.Close(port)
	p.Opened = true
	return p.Options.Open(p.Name)
}

func (p *Port) Close(port *serial.Port)  {
	p.Opened = false
	port.Close()
}

//Read reads the incoming buffer in a specific port
func Read(port *serial.Port) (string, error) {

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
