package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var host, port string
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connect timeout")
	flag.Parse()

	host = flag.Arg(0)
	port = flag.Arg(1)

	host = "localhost"
	port = "4241"

	client := NewTelnetClient(
		host+":"+port,
		timeout,
		os.Stdin,
		os.Stdout,
	)

	sigintChan := make(chan os.Signal)
	signal.Notify(sigintChan, syscall.SIGINT)

	connectErr := client.Connect()
	if connectErr != nil {
		// todo:
		fmt.Println(connectErr)
		os.Exit(1)
	}
	defer client.Close()

	receiveChan := make(chan error)
	go func(doneCh chan<- error) {
		receiveErr := client.Receive()
		if receiveErr != nil {
			if receiveErr == io.EOF {
				close(doneCh)
			}
			// todo:?
		}
	}(receiveChan)

	sendChan := make(chan error)
	go func(doneCh chan<- error) {
		sendErr := client.Send()
		if sendErr != nil {
			if errors.Is(sendErr, io.EOF) {
				close(doneCh)
			}
			// todo
		}
	}(sendChan)

	for {
		select {
		case <-receiveChan:
			fmt.Fprintln(os.Stderr, "Connection was closed by server")

			client.Close()
			os.Exit(1)
		case <-sendChan:
			fmt.Fprintln(os.Stderr, "Connection was closed by peer")

			client.Close()
			os.Exit(1)
		case <-sigintChan:
			fmt.Fprintln(os.Stderr, "SIGINT. Closing.")

			client.Close()
			os.Exit(1)
		}
	}

	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
