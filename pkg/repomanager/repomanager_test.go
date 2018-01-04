package repomanager

import (
	"errors"
	"fmt"
	"github.com/jrxfive/packagetree/pkg/server"
	"math/rand"
	"net"
	"testing"
	"time"
)

type MockBackend struct {
}

func (mb *MockBackend) Exists(name string) bool {
	if name == "generate-error" {
		return false
	}

	return true
}

func (mb *MockBackend) Add(name string, edges ...string) error {
	if name == "generate-error" {
		return errors.New("error")
	}

	return nil
}

func (mb *MockBackend) Remove(name string) error {
	if name == "generate-error" {
		return errors.New("error")
	}

	return nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func TestValidatePackage(t *testing.T) {
	var tests = []struct {
		input    string
		expected bool
	}{
		{"cless", false},
		{"clang-omp", false},
		{"gnome-doc-utils", false},
		{"emacs=elisp", true},
		{"emacs-elisp++", false},
		{"emacs+elisp", true},
		{"emacs elisp", true},
		{"dvd+rw-tools", false},
		{"g++", false},
	}

	for _, test := range tests {
		err := validatePackage(test.input)
		var errResult bool

		if err == nil {
			errResult = false
		} else {
			errResult = true
		}

		if errResult != test.expected {
			t.Errorf("validatePackage(%q) = %v", test.input, err)
		}

	}

}

func TestValidateAndCreateOperator(t *testing.T) {

	r := NewRepo(&MockBackend{})

	var tests = []struct {
		instruction     *instruction
		repo            *Repo
		expectedCommand string
	}{
		{&instruction{"INDEX", "emacs=elisp", []string{}}, r, "ERROR"},
		{&instruction{"INDEX", "g++", []string{}}, r, "INDEX"},
		{&instruction{"INDEX", "g++", []string{"ba+r"}}, r, "ERROR"},
		{&instruction{"INDE", "g++", []string{}}, r, "ERROR"},
		{&instruction{"QUERY", "g++", []string{}}, r, "QUERY"},
		{&instruction{"QUERY", "g++", []string{"bar"}}, r, "QUERY"},
		{&instruction{"QUERY", "emacs=elisp", []string{}}, r, "ERROR"},
		{&instruction{"QRY", "g++", []string{}}, r, "ERROR"},
		{&instruction{"REMOVE", "g++", []string{}}, r, "REMOVE"},
		{&instruction{"REMOVE", "g++", []string{"bar"}}, r, "REMOVE"},
		{&instruction{"REMOV", "emacs=elisp", []string{}}, r, "ERROR"},
		{&instruction{"REM", "g++", []string{}}, r, "ERROR"},
	}

	for _, test := range tests {
		if output := validateAndCreateOperator(test.instruction, test.repo); output.GetCommand() != test.expectedCommand {
			t.Errorf("validateAndCreateOperator(%#v) = %s, expected:%s", test.instruction, output.GetCommand(), test.expectedCommand)
		}
	}
}

func TestCreateInstructionSet(t *testing.T) {

	var instructions = []struct {
		rawInstruction           []byte
		expectedCMD              string
		expectedPackageName      string
		expectedDependencyLength int
		expectedErrResult        bool
	}{
		{[]byte("INDEX|boo|\n"), "INDEX", "boo", 0, false},
		{[]byte("INDEX|boo|foo\n"), "INDEX", "boo", 1, false},
		{[]byte("INDEX|boo|foo,bar\n"), "INDEX", "boo", 2, false},
		{[]byte("QUERY"), "", "", 0, true},
	}

	for _, instruct := range instructions {
		var errResult bool
		i, err := createInstructionSet(instruct.rawInstruction)

		if err == nil {
			errResult = false
		} else {
			errResult = true
		}

		if errResult != instruct.expectedErrResult {
			t.Errorf("createInstructionSet(%q) = %q, expected:%v", i, err, instruct.expectedErrResult)
		}

		if i.cmd != instruct.expectedCMD {
			t.Errorf("createInstructionSet(%q) = %v, expected:%v", i, string(i.cmd), instruct.expectedCMD)
		}

		if i.packageName != instruct.expectedPackageName {
			t.Errorf("createInstructionSet(%q) = %v, expected:%v", i, string(i.packageName), instruct.expectedPackageName)
		}

		if len(i.packageDependencies) != instruct.expectedDependencyLength {
			t.Errorf("createInstructionSet(%q) = %v, expected:%v", i, len(i.packageDependencies), instruct.expectedDependencyLength)
		}
	}
}

func TestHandle(t *testing.T) {

	port := random(1025, 9999)

	go func() {
		r := NewRepo(&MockBackend{})
		configuration := server.NewServerConfiguration(r, port, 1, 10)

		err := server.Listen(configuration)
		if err != nil {
			fmt.Println(err)
			t.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	var connections = []struct {
		rawCommand              string
		expectedReturnValue     string
		shouldCloseBeforeReturn bool
	}{
		{"INDEX|boo|", OK, false},
		{"QUERY|boo|", OK, false},
		{"REMOVE|boo|", OK, false},
		{"REMOVE|boo|extra|", ERROR, false},
		{"INDEX|generate-error|", FAIL, false},
		{"QUERY|generate-error|", FAIL, false},
		{"REMOVE|generate-error|", FAIL, false},
		{"REMOVE|boo|", "", true},
	}

	c, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		t.Fatalf("Failed to establish connection to repo server")
	}
	defer c.Close()

	for _, connection := range connections {

		_, err = fmt.Fprintln(c, connection.rawCommand)
		if err != nil {
			t.Error(err)
		}
		if connection.shouldCloseBeforeReturn {
			c.Close()
			continue
		}

		c.SetReadDeadline(time.Now().Add(time.Second * 3))
		buffer := make([]byte, 1024)

		bufferLength, err := c.Read(buffer)
		if err != nil {
			c.Close()
			t.Error(err)
		}

		if string(buffer[:bufferLength]) != fmt.Sprintf("%s\n", connection.expectedReturnValue) {
			t.Errorf("Command:%s Expected:%s Got:%s", connection.rawCommand, connection.expectedReturnValue, string(buffer[:bufferLength]))
		}
	}

	c.Close()
}
