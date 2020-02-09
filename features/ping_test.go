package feature

import (
	"fmt"
	"net/http"

	"github.com/orov-io/LongBen/client"
	"github.com/orov-io/LongBen/models"
)

var lastError error
var pong *models.Pong
var resp *http.Response

func iHaveAPingCall() error {
	serviceClient := client.NewWithDefaults()
	pong, lastError = serviceClient.Ping()
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
