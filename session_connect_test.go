package gocql

import (
	"golang.org/x/net/context"
	"net"
	"sync"
	"testing"
	"time"
)

type OneConnTestServer struct {
	Err     error
	Address string
	Host    string
	Port    int

	listener   net.Listener
	acceptChan chan struct{}
	mu         sync.Mutex
	closed     bool
}

func NewOneConnTestServer() (*OneConnTestServer, error) {
	lstn, err := net.Listen("tcp4", "localhost:0")
	if err != nil {
		return nil, err
	}
	host, port, _ := parseHostPort(lstn.Addr().String())
	return &OneConnTestServer{
		listener:   lstn,
		acceptChan: make(chan struct{}),
		Address:    lstn.Addr().String(),
		Host:       host,
		Port:       port,
	}, nil
}

func (c *OneConnTestServer) Accepted() chan struct{} {
	return c.acceptChan
}

func (c *OneConnTestServer) Close() {
	c.lockedClose()
}

func (c *OneConnTestServer) Serve() {
	conn, err := c.listener.Accept()
	c.Err = err
	if conn != nil {
		conn.Close()
	}
	c.lockedClose()
}

func (c *OneConnTestServer) lockedClose() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		close(c.acceptChan)
		c.listener.Close()
		c.closed = true
	}
}

func testConnErrorHandler(t *testing.T) ConnErrorHandler {
	return connErrorHandlerFn(func(conn *Conn, err error, closed bool) {
		t.Errorf("in connection handler: %v", err)
	})
}

func assertConnectionEventually(t *testing.T, wait time.Duration, srvr *OneConnTestServer) {
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			t.Errorf("waiting for connection: %v", ctx.Err())
		}
	case <-srvr.Accepted():
		if srvr.Err != nil {
			t.Errorf("accepting connection: %v", srvr.Err)
		}
	}
}

func TestSession_connect_WithNoTranslator(t *testing.T) {
	srvr, err := NewOneConnTestServer()
	assertNil(t, "error when creating tcp server", err)
	defer srvr.Close()

	session := createTestSession()
	defer session.Close()

	go srvr.Serve()

	session.connect(srvr.Address, testConnErrorHandler(t), &HostInfo{
		peer: srvr.Host,
		port: srvr.Port,
	})

	assertConnectionEventually(t, 500*time.Millisecond, srvr)
}

func TestSession_connect_WithTranslator(t *testing.T) {
	srvr, err := NewOneConnTestServer()
	assertNil(t, "error when creating tcp server", err)
	defer srvr.Close()

	session := createTestSession()
	defer session.Close()
	session.cfg.AddressTranslator = staticAddressTranslator(srvr.Host, srvr.Port)

	go srvr.Serve()

	// the provided address will be translated
	session.connect("10.10.10.10:5432", testConnErrorHandler(t), &HostInfo{
		peer: "10.10.10.10",
		port: 5432,
	})

	assertConnectionEventually(t, 500*time.Millisecond, srvr)
}
