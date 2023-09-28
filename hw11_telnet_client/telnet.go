package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type BaseTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	connection net.Conn
}

func (c *BaseTelnetClient) Connect() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		// todo: wrap
		return err
	}

	c.connection = conn

	return nil
}

func (c *BaseTelnetClient) Close() error {
	err := c.connection.Close()
	if err != nil {
		// todo: wrap
	}
	// todo: what else?
	return nil
}

func (c *BaseTelnetClient) Send() error {
	if c.connection == nil {
		// todo:
	}

	scanner := bufio.NewScanner(c.in)
	for i := 0; i < 3; i++ { // todo: inf loop
		if !scanner.Scan() {
			// todo:
		}

		_, err := c.connection.Write([]byte(fmt.Sprintf("%s\n", scanner.Bytes())))
		if err != nil {
			// todo: wrap
			return err
		}
	}

	return nil
}

func (c *BaseTelnetClient) Receive() error {
	if c.connection == nil {
		// todo
	}

	for i := 0; i < 3; i++ { // todo: inf loop
		scanner := bufio.NewScanner(c.connection)
		if !scanner.Scan() {
			// todo
		}

		_, err := fmt.Fprintln(c.out, scanner.Text())
		if err != nil {
			// todo
		}
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &BaseTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
