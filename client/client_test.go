package client_test

import (
	"testing"

	"github.com/orov-io/LongBen/client"
	. "github.com/smartystreets/goconvey/convey"
)

// convey phrases
const (
	givenAClient         = "Given a LongBen client"
	callHandlerByService = "When call is handler by the service"
	responseShouldBeOK   = "Then response should be OK"
)

func TestPing(t *testing.T) {
	Convey(givenAClient, t, func() {
		serviceClient := client.NewWithDefaults()
		Convey(callHandlerByService, func() {
			pong, err := serviceClient.Ping()
			Convey(responseShouldBeOK, func() {

				So(err, ShouldBeNil)
				So(pong.Message, ShouldNotBeEmpty)
			})
		})
	})
}
