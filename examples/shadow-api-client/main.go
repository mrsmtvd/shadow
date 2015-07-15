package main // import "github.com/kihamo/shadow/examples/shadow-api-client"

import (
	"flag"
	"fmt"
	"log"

	"gopkg.in/jcelliott/turnpike.v2"
)

var (
	host  = flag.String("host", "0.0.0.0", "Service host")
	port  = flag.Int64("port", 8001, "Service port")
	debug = flag.Bool("debug", false, "Debug mode")
)

func main() {
	flag.Parse()

	if *debug {
		turnpike.Debug()
	}

	client, err := turnpike.NewWebsocketClient(turnpike.MSGPACK, fmt.Sprintf("ws://%s:%d/", *host, *port))
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = client.JoinRealm("api", turnpike.ALLROLES, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	response, err := client.Call("api.ping", nil, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Server response: %s\n", response.Arguments[0])

	response, err = client.Call("api.version", nil, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Server response: %s\n", response.ArgumentsKw)
}
