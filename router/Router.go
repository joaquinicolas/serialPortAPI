package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinicolas/Elca/Ports"
	"github.com/joaquinicolas/Elca/Store"
	"github.com/joaquinicolas/Elca/Store/models"

	"io/ioutil"
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
	b, err := ioutil.ReadAll(r.Body)

	json.Unmarshal(b, &port)

	if err != nil {
		requestError = &RequestError{
			ErrorString:    err.Error(),
			HttpStatusCode: http.StatusInternalServerError,
		}
		requestError.Error(w)
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
		return
	}

	msgCh := make(chan string, 5)
	errCh := make(chan error)
	go func() {
		for {
			Ports.Read(serialPort, msgCh, errCh)
		}
	}()

	go func() {
		for {
			select {
			case msg := <-msgCh:
				n := &models.News{
					Text: msg,
				}
				store, _ := Store.GetStore("sqlite3")
				store.StoreNews(n)
			case err := <-errCh:
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
