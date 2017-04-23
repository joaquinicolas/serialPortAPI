package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	logger "github.com/joaquinicolas/Elca/Logger"
	"github.com/joaquinicolas/Elca/Novelty"
	"github.com/joaquinicolas/Elca/Ports"
	"github.com/joaquinicolas/Elca/Store"

)

type RequestError struct {
	ErrorString    string
	HttpStatusCode int
	Excepton       error
}

// Error writes the error to http.ResponseWritter
func (re *RequestError) Error(w http.ResponseWriter) {
	w.WriteHeader(re.HttpStatusCode)
	fmt.Fprintf(w, "{error: %q, StackTrace: %q}", re.ErrorString, re.Excepton)
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
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
		Name:        "Port",
		Method:      "POST",
		Pattern:     "/port",
		HandlerFunc: OpenandRead,
	},
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintf(w, "Hola Mundo")
}

func ListPorts(w http.ResponseWriter, r *http.Request) {
	ports, _ := Ports.List()

	json.NewEncoder(w).Encode(ports)
}

func OpenandRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var port Ports.Port
	var requestError *RequestError
	err := json.NewDecoder(r.Body).Decode(&port)

	if err != nil {
		requestError = &RequestError{
			ErrorString:    err.Error(),
			HttpStatusCode: http.StatusInternalServerError,
		}
		requestError.Error(w)
		logger.Error.Println(err)
		return
	}

	defer r.Body.Close()

	serialPort, err := port.Open()
	if err != nil {
		requestError = &RequestError{
			ErrorString:    err.Error(),
			HttpStatusCode: http.StatusInternalServerError,
			Excepton:       err,
		}
		requestError.Error(w)
		logger.Error.Println(err)
		return
	}

	msgCh := make(chan string, 5)
	errCh := make(chan error)
	go func() {
		for {
			Ports.Read(serialPort, msgCh, errCh)
			select {
			case msg := <-msgCh:
				n := &Novelty.Novelty{
					Text: msg,
				}
				db, err := Store.GetStore("sqlite3")
				if err != nil {
					logger.Error.Println(err)
				}

				db.StoreNovelty(n)

			case err := <-errCh:
				logger.Error.Println(err)

			}
		}
	}()
	fmt.Fprintf(w, "{data: %s change status to open}", port.Name)
	return
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
