package repomanager

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

const (
	plusSignCharacterValue = 43
	OK                     = "OK"
	FAIL                   = "FAIL"
	ERROR                  = "ERROR"
)

var allowedCharacters = regexp.MustCompile(`([a-zA-z\_\-\+\.\d]+)`)

type Backend interface {
	Exists(name string) bool
	Add(name string, edges ...string) error
	Remove(name string) error
}

type Repo struct {
	backend Backend
	mu      *sync.RWMutex
}

type Operation interface {
	Run() (string, error)
	GetCommand() string
}

type instruction struct {
	cmd                 string
	packageName         string
	packageDependencies []string
}

func NewRepo(graph Backend) *Repo {

	return &Repo{
		backend: graph,
		mu:      &sync.RWMutex{},
	}
}

func createInstructionSet(input []byte) (*instruction, error) {

	var CMD_INDEX int = 0
	var PACKAGE_NAME_INDEX int = 1
	var DEPENDENCY_INDEX int = 2

	trimmedInput := bytes.Trim(input, "\n")
	delimit := bytes.Split(trimmedInput, []byte{'|'})

	if len(delimit) == 3 {
		return &instruction{
			cmd:                 string(delimit[CMD_INDEX]),
			packageName:         string(delimit[PACKAGE_NAME_INDEX]),
			packageDependencies: getDependencyPackages(delimit[DEPENDENCY_INDEX]),
		}, nil
	} else {
		return &instruction{}, errors.New(fmt.Sprintf("Invalid Protocol Specification"))
	}
}

func getDependencyPackages(dependencies []byte) []string {
	var dep []string

	if len(dependencies) > 0 {
		r := strings.SplitN(string(dependencies), ",", -1)

		dep = r
	}

	return dep
}

func validatePackage(packageName string) error {

	packageNameLength := len(packageName)
	allowedCharactersMatches := allowedCharacters.FindAllString(packageName, -1)

	if len(allowedCharactersMatches) != 1 {
		return errors.New("Package can only contain a-zA-Z9+-_.")
	}

	if check := strings.Contains(packageName, "+"); check {

		for idx, char := range packageName {

			if unicode.IsSymbol(char) {
				if char == plusSignCharacterValue {

					if idx == packageNameLength-1 {
						break
					}

					if idx < packageNameLength-1 {
						nextChar := packageName[idx+1]

						if unicode.IsLetter(int32(nextChar)) {

							if !strings.Contains(packageName, "-") {
								return errors.New("Package name should be split by '-'")
							} else {
								return nil
							}

						}
					}
				}

			}
		}
	}

	return nil

}

func validateDependencies(dependencies []string) error {

	for _, dependency := range dependencies {
		return validatePackage(dependency)
	}
	return nil
}

func validateAndCreateOperator(instruction *instruction, repo *Repo) Operation {

	err := validatePackage(instruction.packageName)
	if err != nil {
		return NewUnknownOperator()
	}

	err = validateDependencies(instruction.packageDependencies)
	if err != nil {
		return NewUnknownOperator()
	}

	switch instruction.cmd {
	case INDEX:
		return NewIndexOperator(instruction, repo)
	case REMOVE:
		return NewRemoveOperator(instruction, repo)
	case QUERY:
		return NewQueryOperator(instruction, repo)
	default:
		return NewUnknownOperator()
	}
}

//Handles TCP connections for the repo manager. Will set a read dead line before
//serving the connection to avoid stale connections. The buffer is set to 1024 bytes
//this could have been more dynamic by using bufio instead of a static limit. After
//receiving and creating an instruction set it will run the corresponding command and
//return the value based on the backend.
func (r *Repo) Handle(conn net.Conn) {

	for {
		//err := conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(r.readTimeout)))
		//if err != nil {
		//	conn.Close()
		//	break
		//}

		buffer := make([]byte, 1024)

		bufferLength, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			break
		}

		instruction, err := createInstructionSet(buffer[:bufferLength])
		if err != nil {
			fmt.Fprintln(conn, ERROR)
			continue
		}

		operator := validateAndCreateOperator(instruction, r)
		output, err := operator.Run()
		_, err = fmt.Fprintln(conn, output)
		if err != nil {
			conn.Close()
		}

	}
}
