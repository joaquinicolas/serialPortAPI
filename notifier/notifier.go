package notifier

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/joaquinicolas/Elca/Logger"
	"github.com/joaquinicolas/Elca/Novelty"
	"errors"
)

const (
	SERVER_URL       = "http://131.255.5.183:9091/"
	NOVELTY_RECEIVER = "News"
	ALIVE_RECEIVER   = "Alive"
	NETWORK_INTERFACE_NAME = "eth0"
)

var mac string

func NotifyNovelty(novelty *Novelty.Novelty) error{

	var jsonStr = []byte(fmt.Sprintf(`{"Mac":"%s", "Data":"%s"}`, mac, novelty.Text))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", SERVER_URL, NOVELTY_RECEIVER), bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return errors.New(fmt.Sprintf("Alive request response with status code: %d. Status: %s.\n Body: %s",resp.StatusCode, resp.Status, bodyString))
	}
	return nil
}

func NotifyAlive() error{

	var jsonStr = []byte(fmt.Sprintf(`{"mac":"%s"}`, mac))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", SERVER_URL, ALIVE_RECEIVER), bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return errors.New(fmt.Sprintf("Alive request response with status code: %d. Status: %s.\n Body: %s", resp.StatusCode, resp.Status, bodyString))
	}

	return nil
}

func init() {

	m, err := net.InterfaceByName(NETWORK_INTERFACE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	mac = m.HardwareAddr.String()
}
