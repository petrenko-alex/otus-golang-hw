package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
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

	sigintChan := make(chan os.Signal)
	signal.Notify(sigintChan, syscall.SIGINT)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	connectErr := client.Connect()
	if connectErr != nil {
		log.Fatal(connectErr)
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
			return
		case <-sendChan:
			fmt.Fprintln(os.Stderr, "Connection was closed by peer")
			return
		case <-sigintChan:
			fmt.Fprintln(os.Stderr, "SIGINT. Closing.")
			return
		}
	}

	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
