package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joaquinicolas/Elca/Ports"
)

var (
	r     *http.Request
	w     *httptest.ResponseRecorder
	ports []Ports.Port
)

func TestListPorts(t *testing.T) {
	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	CreateRouter().ServeHTTP(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(body, &ports)
	if err != nil {
		t.Error(err)
	}
	for _, p := range ports {
		fmt.Println(p)
	}

}

func TestOpenandRead(t *testing.T) {
	out, _ := json.Marshal(ports[0])
	reader := strings.NewReader(string(out))
	r, _ = http.NewRequest("POST", "/port", reader)
	CreateRouter().ServeHTTP(w, r)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(body))
}
