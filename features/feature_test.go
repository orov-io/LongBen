package feature

import (
	"fmt"
	"net/http"

	"github.com/DATA-DOG/godog"
	"github.com/orov-io/LongBen/client"
	"github.com/orov-io/LongBen/models"
)

var lastError error
var pong *models.Pong
var resp *http.Response

func iHaveAPingCall() error {
	pong, lastError = client.Ping()
	return nil
}

func iReceiveTheResponse() error {
	if lastError != nil {
		return fmt.Errorf("Call fails with error: %v", lastError)
	}
	return nil
}

func iShouldReceiveAPongResponse() error {
	if pong.Status == "" || pong.Message == "" {
		return fmt.Errorf("Pong is empty")
	}
	return nil
}

func iHaveAnInvalidCall() error {
	resp, lastError = http.Get("http://localhost:8080/v1/longBen/invalid")
	return nil
}

func codeIs404() error {
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Unexpected response code: %v", resp.StatusCode)
	}
	return nil
}
func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a ping call$`, iHaveAPingCall)
	s.Step(`^I receive the response$`, iReceiveTheResponse)
	s.Step(`^I should receive a pong response$`, iShouldReceiveAPongResponse)

	s.Step(`^I have an invalid call$`, iHaveAnInvalidCall)
	s.Step(`^Code should be a Not Found HTTP Code$`, codeIs404)

	/* 	s.BeforeSuite(func() {
	   		upServer()
	   	})

	   	s.AfterSuite(func() {
	   		downServer()
	   	}) */
}
