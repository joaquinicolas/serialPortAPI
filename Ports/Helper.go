package Ports

import "github.com/joaquinicolas/Elca/Store"

type PortStorer interface {

	Store.Storer
}

func ReadAndStore()  {

	for  {
		select {
		case s := S:
			
		}
	}
}