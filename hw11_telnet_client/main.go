package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	var timeout, host, port string
	var wg sync.WaitGroup
	flag.StringVar(&timeout, "timeout", "10s", "connect timeout")
	flag.Parse()

	host = flag.Arg(0)
	port = flag.Arg(1)

	host = "localhost"
	port = "4241"

	timeoutObj, err := time.ParseDuration(timeout)
	if err != nil {
		// todo:
	}

	client := NewTelnetClient(
		host+":"+port,
		timeoutObj,
		os.Stdin,
		os.Stdout,
	)

	connectErr := client.Connect()
	if connectErr != nil {
		// todo:
		fmt.Println(connectErr)
		os.Exit(1)
	}
	defer client.Close()

	wg.Add(2)
	go func() {
		receiveErr := client.Receive()
		if receiveErr != nil {

		}

		wg.Done()
	}()

	go func() {
		sendErr := client.Send()
		if sendErr != nil {
			// todo
		}

		wg.Done()
	}()

	wg.Wait()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
