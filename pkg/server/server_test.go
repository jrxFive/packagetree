package server

import (
	"fmt"
	"math/rand"
	"net"
	"syscall"
	"testing"
	"time"
)

type MockConnectionHandler struct {
}

func (mch *MockConnectionHandler) Handle(conn net.Conn) {
	conn.Close()
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix() + 1)
	return rand.Intn(max-min) + min
}

func serverSetup(port int, t *testing.T) *Configuration {
	mch := &MockConnectionHandler{}
	configuration := NewServerConfiguration(mch, port+1, 1, 10)

	go func() {
		err := Listen(configuration)
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	return configuration
}

func TestListenSignal(t *testing.T) {

	port := random(1025, 9999)
	configuration := serverSetup(port, t)

	configuration.SignalChannel <- syscall.SIGINT
}

func TestListenClose(t *testing.T) {

	port := random(1025, 9999)
	configuration := serverSetup(port, t)

	close(configuration.ConnectionChannel)
}

func TestListenPortTaken(t *testing.T) {

	port := random(1025, 9999)
	configuration := serverSetup(port, t)

	err := Listen(configuration)
	if err == nil {
		t.Fatal("Port should be in use and server retuning error")
	}
}

func TestListenMessagesClose(t *testing.T) {

	port := random(1025, 9999)
	serverSetup(port, t)

	var connections = []struct {
		rawCommand string
	}{
		{"QUERY|boo|"},
	}

	c, err := net.Dial("tcp", fmt.Sprintf(":%v", port+1))
	if err != nil {
		t.Fatalf("Failed to establish connection to repo server")
	}
	defer c.Close()

	for _, connection := range connections {
		_, err = fmt.Fprintln(c, connection.rawCommand)
	}
}
