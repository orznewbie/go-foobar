package main

import (
	"fmt"
	"time"
)

const (
	DefaultHost    = "127.0.0.1"
	DefaultPort    = "80"
	DefaultTimeout = 30 * time.Second
)

type Server struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func New(opts ...Options) *Server {
	var server = &Server{
		Host:    DefaultHost,
		Port:    DefaultPort,
		Timeout: DefaultTimeout,
	}
	for _, opt := range opts {
		opt(server)
	}
	return server
}

type Options func(*Server)

func WithHost(host string) Options {
	return func(server *Server) {
		server.Host = host
	}
}

func WithPort(port string) Options {
	return func(server *Server) {
		server.Port = port
	}
}
func WithTimeout(timeout time.Duration) Options {
	return func(server *Server) {
		server.Timeout = timeout
	}
}

func main() {
	s := New(WithPort("8080"), WithTimeout(60*time.Second))
	fmt.Println(s)
}
