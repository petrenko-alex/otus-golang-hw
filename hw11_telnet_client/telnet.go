package main

import (
	"bufio"
	"errors"
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
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
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
	for {
		if !scanner.Scan() {
			return io.EOF
		}

		_, err := c.connection.Write([]byte(fmt.Sprintf("%s\n", scanner.Bytes())))
		if err != nil {
			// todo: wrap
			return err
		}
	}
}

func (c *BaseTelnetClient) Receive() error {
	if c.connection == nil {
		// todo
	}

	scanner := bufio.NewScanner(c.connection)
	for {
		if !scanner.Scan() {
			return io.EOF
		}

		_, err := fmt.Fprintln(c.out, scanner.Text())
		if err != nil {
			return errors.New("error printing received msg")
		}
	}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &BaseTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
