package main

import ("net/http"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/joaquinicolas/Elca/Ports"
	"encoding/json"
	"log"
)

type RequestError struct {
	ErrorString string
	HttpStatusCode int
}

// Error writes the error to http.ResponseWritter
func (re *RequestError) Error(w http.ResponseWriter)  {
	w.WriteHeader(re.HttpStatusCode)
	fmt.Fprintf(w,"{error: %q}",re.ErrorString)
}
type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		ListPorts,
	},
	Route{
		Name:"Port",
		Method:"POST",
		Pattern:"/port",
		HandlerFunc:OpenandRead,
	},
}


func IndexHandler(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()
	fmt.Fprintf(w,"Hola Mundo")
}

func ListPorts(w http.ResponseWriter,r *http.Request)  {
	ports,_ := Ports.List()

	defer r.Body.Close()
	json.NewEncoder(w).Encode(ports)
}



func OpenandRead(w http.ResponseWriter, r *http.Request)  {
	var port *Ports.Port
	var requestError *RequestError
	err := json.NewDecoder(r.Body).Decode(port)
	log.Println(port)
	log.Println(err)
	if err != nil {
		requestError = &RequestError{
			ErrorString:err.Error(),
			HttpStatusCode:http.StatusInternalServerError,
		}
		requestError.Error(w)
		return
	}

	defer r.Body.Close()

	dataChannel := make(chan string,1)
	errorChannel := make(chan error,1)
	go func() {
		for {
			d,e := port.Read()
			if e != nil {
				errorChannel <- e
			}else{
				dataChannel <- d
			}

		}
	}()


	select {
	case err := <- errorChannel:
		requestError = &RequestError{
			ErrorString: err.Error(),
			HttpStatusCode:http.StatusInternalServerError,
		}
		requestError.Error(w)
		return

	case data := <- dataChannel:
		fmt.Fprintf(w,"{data: %c}",data)
		return
	}
}

func CreateRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}