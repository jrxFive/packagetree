package server

import (
	"fmt"
	"github.com/jrxfive/packagetree/pkg/logging"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var logger = logging.GetLogger()

type ConnectionHandler interface {
	Handle(conn net.Conn)
}

type Configuration struct {
	Handler           ConnectionHandler
	Port              int
	MaxHerd           int
	Timeout           time.Duration
	ConnectionChannel chan net.Conn
	SignalChannel     chan os.Signal
	mu                sync.RWMutex
}

//Creates and returns a new configuration that can has a net.Conn Handler. The
//port to bind the TCP server to, how large the buffer should be for simultaneous connections
//to the connection channel.
func NewServerConfiguration(handler ConnectionHandler, port, maxHerd, timeout int) *Configuration {

	signalChannel := make(chan os.Signal, 1)
	connChannel := make(chan net.Conn, maxHerd)

	return &Configuration{
		Handler:           handler,
		Port:              port,
		Timeout:           time.Duration(timeout),
		MaxHerd:           maxHerd,
		ConnectionChannel: connChannel,
		SignalChannel:     signalChannel,
		mu:                sync.RWMutex{},
	}
}

//Binds a TCP server to handle connections for anything that use the interface ConnectionHandler. Will
//close the server and channel on SIGINT or SIGTERM.
func Listen(sc *Configuration) error {

	signal.Notify(sc.SignalChannel, syscall.SIGINT, syscall.SIGTERM)
	channelBreaker := false

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", sc.Port))
	if err != nil {
		return err
	}

	go func(listener net.Listener, connChannel chan net.Conn, signalChannel chan os.Signal) {

		for {
			select {
			case s := <-signalChannel:
				logger.Printf("Signal:%v received closing channel and listener\n", s)
				close(connChannel)
				listener.Close()
			case conn, ok := <-connChannel:
				if ok {
					err = conn.SetDeadline(time.Now().Add(time.Second * sc.Timeout))
					if err != nil {
						conn.Close()
					} else {
						go sc.Handler.Handle(conn)
					}
				} else {
					sc.mu.Lock()
					channelBreaker = true
					sc.mu.Unlock()
				}
			}
		}

	}(listener, sc.ConnectionChannel, sc.SignalChannel)

	for {
		conn, err := listener.Accept()
		if err != nil {
			sc.mu.RLock()
			if channelBreaker {
				sc.mu.RUnlock()
				break
			}
			sc.mu.RUnlock()
		}

		sc.mu.RLock()
		if !channelBreaker {
			sc.ConnectionChannel <- conn
		}
		sc.mu.RUnlock()
	}

	return nil
}
