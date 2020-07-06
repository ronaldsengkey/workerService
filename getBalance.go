package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/robfig/cron"
	"net"
	"os"
	"os/signal"
	"syscall"
	"github.com/takama/daemon"
)

const (
	// name of the service
	name        = "myservice"
	description = "My Echo Service"
	// port which daemon should be listen
	port = ":9970"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	fmt.Printf("Start Golang\n")
	c := cron.New()
	c.AddFunc("23 10 * * *", func() { 
		// log.Println("Start Generate Saldo")
		if checkLastDay() {
			generateSaldo()
		} 
	})
	log.Println("Start cron")
	c.Start()
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

type Customer struct {
	Id string `json:"id"`
	Nominal string `json:"nominal"`
	Periode string `json:"periode"`
}

type Response struct {
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Data []Customer
}

func getCustomer() Response{
	url := "http://localhost:8089/wallet/cronjob/getCustomer"
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("clientKey", "clientKey")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	// log.Println(res.ResponseMessage)
	return res
} 

func getSaldo(id string) Response{
	url := "http://localhost:8089/wallet/cronjob/getSaldo/" + id
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("clientKey", "clientKey")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	// log.Println(res.ResponseMessage)
	return res
}

func generateSaldo() {
	log.Printf("start generateSaldo")
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	response := getCustomer()
	// log.Println(response)
	if response.ResponseCode == "200" {
		for _, data := range response.Data {
			res := getSaldo(data.Id)
			if res.ResponseCode == "200" {
				_, err = conn.Do("HMSET", "customerSaldo:" + data.Id, "id", data.Id, "nominal", res.Data[0].Nominal, "periode", res.Data[0].Periode)
				if err != nil {
					log.Println(err)
				}
				// log.Println("HMSET", "customerSaldo:" + data.Id, "id", data.Id, "nominal", res.Data[0].Nominal, "periode", res.Data[0].Periode)
			} else {
				log.Println("error: ", res)
			}
		}
	} else {
		log.Println("error: ", response)
	}
	log.Printf("end generateSaldo")
}

func checkLastDay() bool{
	currentTime := time.Now().Format("2006-01-02")
	lastDay := getLastDay().Format("2006-01-02")
	log.Println("current: ", currentTime)
	log.Println("lastDay: ", lastDay)
	if currentTime == lastDay {
		log.Println("lastDay: true")
		return true
	} else {
		log.Println("lastDay: false")
		return false
	}
}

func getLastDay() time.Time{
	now := time.Now()
    currentYear, currentMonth, _ := now.Date()
    currentLocation := now.Location()

    firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
    lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

    // fmt.Println(firstOfMonth)
	// fmt.Println(lastOfMonth)
	
	return lastOfMonth
}