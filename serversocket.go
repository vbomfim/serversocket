//Package serversocket offers an API to create server sockets TCP and UDP
package serversocket

import (
	"errors"
	"fmt"
	"net"
)

var ErrConnHandlerNil error = errors.New("ServerSocketTCP with nil connHandler not allowed")

type TCPConnHandler interface {
	Handle(conn *net.Conn)
}

type ServerSocketTCP struct {
	Addr        string
	TCPAddr     *net.TCPAddr
	connHandler TCPConnHandler
}

//NewServerSocketTCP is a helper function that returns an instance of ServerSocketTCP
//It is useful when you need to pass a custom Listener to the Serve method.
func NewServerSocketTCP(address string, connHandler TCPConnHandler) (*ServerSocketTCP, error) {
	if connHandler == nil {
		return nil, ErrConnHandlerNil
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("invalid address. %w", err)
	}
	return &ServerSocketTCP{
		Addr:        address,
		connHandler: connHandler,
		TCPAddr:     tcpAddr,
	}, nil
}

//TCPListenAndServe is a helper function that creates a listener and starts the server
//The address param follows the same directions as net.Listen
func TCPListenAndServe(address string, connHandler TCPConnHandler) error {
	srv, err := NewServerSocketTCP(address, connHandler)
	if err != nil {
		return fmt.Errorf("failed creating ServerSocketTCP with the params address: %s and connHandler: %v %w", address, connHandler, err)
	}

	ln, err := net.ListenTCP("tcp", srv.TCPAddr)
	if err != nil {
		return err
	}

	return srv.Serve(ln)
}

//Serve accepts a net.Listener param and starts a loop accepting connections and calling the server's connection handler in a go routine
func (srv *ServerSocketTCP) Serve(ln net.Listener) error {
	if srv.connHandler == nil {
		return ErrConnHandlerNil
	}
	//TODO:
	// 1 - failed to listen treatment
	// 2 - handle external shutdown signal
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go srv.connHandler.Handle(&conn)
	}
}
