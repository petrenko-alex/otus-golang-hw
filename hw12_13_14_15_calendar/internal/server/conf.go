package server

import "time"

type Options struct {
	GRPC GRPCOptions
	HTTP HTTPOptions
}

type GRPCOptions struct {
	Host, Port     string
	ConnectTimeout time.Duration
}

type HTTPOptions struct {
	Host, Port                string
	ReadTimeout, WriteTimeout time.Duration
}
