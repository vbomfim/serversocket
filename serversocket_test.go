package serversocket_test

import (
	"net"
	"testing"

	ss "github.com/vbomfim/serversocket"
)

type ConnHandlerStubFunc func(conn *net.Conn)

func (f ConnHandlerStubFunc) Handle(conn *net.Conn) {
	f(conn)
}

func TestServerSocket(t *testing.T) {

	address := ":5060"
	handler := ConnHandlerStubFunc(func(conn *net.Conn) {})
	t.Run("creating ServerSocketTCP", func(t *testing.T) {
		srv, err := ss.NewServerSocketTCP(address, handler)
		if err != nil {
			t.Fatalf("failed creating the TCP server socket - %v", err)
		}
		if address != srv.Addr {
			t.Fatalf("address not properly set. want %s got %s - %v", address, srv.Addr, err)
		}
	})

	t.Run("creating ServerSocketTCP with invalid handler", func(t *testing.T) {
		_, err := ss.NewServerSocketTCP(address, nil)
		if err == nil {
			t.Fatalf("failed NewServerTCPServer accepted nil handler - %v", err)
		}
	})

	t.Run("creating ServerSocketTCP with invalid address", func(t *testing.T) {
		_, err := ss.NewServerSocketTCP("abcd<-zzz", handler)
		if err == nil {
			t.Fatalf("failed NewServerTCPServer accepted invalid address - %v", err)
		}
	})

}
