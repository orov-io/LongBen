package feature

import (
	"fmt"

	"github.com/DATA-DOG/godog"
)

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a ping call$`, iHaveAPingCall)
	s.Step(`^I receive the response$`, iReceiveTheResponse)
	s.Step(`^I should receive a pong response$`, iShouldReceiveAPongResponse)

	s.Step(`^I have an invalid call$`, iHaveAnInvalidCall)
	s.Step(`^Code should be a Not Found HTTP Code$`, codeIs404)

	wasRunning := false
	s.BeforeSuite(func() {
		if !serviceAlreadyRunning() {
			upServer()
			return
		}
		fmt.Print("Server was running before run tests!")
		wasRunning = true
	})

	s.AfterSuite(func() {
		if !wasRunning {
			downServer()
			return
		}
		fmt.Println("Leaving server living")
	})
}
