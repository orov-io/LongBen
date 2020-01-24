package feature

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/briandowns/spinner"
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
		}
		fmt.Println("Leaving server living")
	})
}

const dockerComposeTimeOut = 120
const dockerCompose = "docker-compose"
const useCustomFile = "-f"
const customFile = "../docker-compose.yml"
const up = "up"
const detachedMode = "-d"
const logs = "logs"
const streamMode = "-f"
const top = "top"

func serviceAlreadyRunning() bool {
	top := exec.Command(dockerCompose, useCustomFile, customFile, top)
	output, err := top.Output()
	if err != nil {
		log.Fatal("Can't run docker-compose ps: ", err)
	}

	return len(output) > 0
}

func upServer() {

	fmt.Println("Starting service...")
	startSpinner := showSpinner()
	upService()
	serviceOutput := getServiceStreamOutput()

	waitToServiceAlive(serviceOutput)

	startSpinner.Stop()
	fmt.Println("Service is running :)")
}

func showSpinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
	s.Start()
	return s
}

func getServiceStreamOutput() io.ReadCloser {
	service := exec.Command(dockerCompose, useCustomFile, customFile, logs, streamMode)
	serviceOutput, err := service.StdoutPipe()
	if err != nil {
		log.Fatalf("Can't open stream with service logs: %v\n", err)
	}
	err = service.Start()
	if err != nil {
		log.Fatalf("Can't do 'make up logs': %v\n", err)
	}
	return serviceOutput
}

func waitToServiceAlive(stream io.ReadCloser) {
	serviceBuild := make(chan bool)
	// This also can be achieved without go routines
	go waitForDockerCompose(stream, serviceBuild)
	isRunning := <-serviceBuild
	if !isRunning {
		downServer()
		log.Fatal("Unable to start the service")
	}
}

func upService() {
	serviceUp := exec.Command(dockerCompose, useCustomFile, customFile, up, detachedMode)
	err := serviceUp.Run()
	if err != nil {
		log.Fatalf("Can't up docker images: %v\n", err)
	}
}

func waitForDockerCompose(serviceOutput io.ReadCloser, serviceBuild chan<- bool) {

	reader := bufio.NewReader(serviceOutput)
	go shutdownIfTimeout(serviceBuild)

	for {
		checkServiceAliveness(reader, serviceBuild)
		return
	}
}

func checkServiceAliveness(reader *bufio.Reader, serviceBuild chan<- bool) {
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			serviceBuild <- false
			close(serviceBuild)
			log.Println("EOF reached")
			return
		}

		if err != nil {
			log.Fatal("Read Error:", err)
			return
		}
		if strings.Contains(line, "Build Failed") {
			serviceBuild <- false
			close(serviceBuild)
			return
		}
		if strings.Contains(line, "Running...") {
			time.Sleep(2 * time.Second)
			serviceBuild <- true
			close(serviceBuild)
			return
		}
	}
}
func shutdownIfTimeout(serviceBuild chan<- bool) {
	time.Sleep(dockerComposeTimeOut * time.Second)
	fmt.Printf("Service don't initialized in first %v seconds\n", dockerComposeTimeOut)
	downServer()
	log.Fatal("Unable to start service")
}

func downServer() {
	fmt.Println("Shutting down service...")
	serviceDown := exec.Command("docker-compose", "-f", "../docker-compose.yml", "down")
	err := serviceDown.Run()
	if err != nil {
		fmt.Println("Can't shutdown service: ", err)
	}
	fmt.Println("Service is down")
}
