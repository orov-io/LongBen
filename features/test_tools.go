package feature

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

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
