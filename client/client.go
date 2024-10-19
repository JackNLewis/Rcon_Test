package client

import (
	"errors"
	"log"
	"time"

	"github.com/gorcon/rcon"
)

const IP_ADDRESS = "0.0.0.0"
const PORT = "25576"
const PASSWORD = "password" // only using for local testing

var conn *rcon.Conn
var command_channel = make(chan *Request)

type Request struct {
	command string
	res     string
	done    chan bool
}

// Initiates connection to minecraft server
func Connect() error {
	var err error
	conn, err = rcon.Dial("0.0.0.0:25575", "password")
	if err != nil {
		return err
	}
	return nil
}

// Closes client connection
func ShutDown() {
	conn.Close()
}

// Starts reading commands to be executed on the server
func Start() {
	// for {
	// 	response, err := conn.Execute("list")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	time.Sleep(2 * time.Second)

	// 	fmt.Println(response)
	// }

	for {
		req := <-command_channel
		response, err := conn.Execute(req.command)
		if err != nil {
			log.Fatal(err)
		}
		req.res = response
		req.done <- true
	}

}

// sends a /list command to the get player channel
func GetPlayerCommand() (string, error) {
	request := &Request{
		command: "list",
		done:    make(chan bool),
	}

	command_channel <- request

	select {
	case <-request.done:
	case <-time.After(2 * time.Second):
		return "", errors.New("error: player list timeout")
	}
	return request.res, nil
}
