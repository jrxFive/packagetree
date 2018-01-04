# Coding Assignment

# Requirements
- Not to use any library apart from your chosen runtime's standard library
- Code must pass the supplied test harness using different random seeds and concurrency factor up to 100

The system you are going to write keeps track of package dependencies. Clients will connect to your server and inform which packages should be indexed, and which dependencies they might have on other packages. We want to keep our index consistent, so your server must not index any package until all of its dependencies have been indexed first. The server should also not remove a package if any other packages depend on it.

The server will open a TCP socket on port 8080. It must accept connections from multiple clients at the same time, all trying to add and remove items to the index concurrently. Clients are independent of each other, and it is expected that they will send repeated or contradicting messages. New clients can connect and disconnect at any moment, and sometimes clients can behave badly and try to send broken messages.

Messages from clients follow this pattern:

```
<command>|<package>|<dependencies>\n
```

Where:
* `<command>` is mandatory, and is either `INDEX`, `REMOVE`, or `QUERY`
* `<package>` is mandatory, the name of the package referred to by the command, e.g. `mysql`, `openssl`, `pkg-config`, `postgresql`, etc.
* `<dependencies>` is optional, and if present it will be a comma-delimited list of packages that need to be present before `<package>` is installed. e.g. `cmake,sphinx-doc,xz`
* The message always ends with the character `\n`

Here are some sample messages:
```
INDEX|cloog|gmp,isl,pkg-config\n
INDEX|ceylon|\n
REMOVE|cloog|\n
QUERY|cloog|\n
```

For each message sent, the client will wait for a response code from the server. Possible response codes are `OK\n`, `FAIL\n`, or `ERROR\n`. After receiving the response code, the client can send more messages.

The response code returned should be as follows:
* For `INDEX` commands, the server returns `OK\n` if the package can be indexed. It returns `FAIL\n` if the package cannot be indexed because some of its dependencies aren't indexed yet and need to be installed first. If a package already exists, then its list of dependencies is updated to the one provided with the latest command.
* For `REMOVE` commands, the server returns `OK\n` if the package could be removed from the index. It returns `FAIL\n` if the package could not be removed from the index because some other indexed package depends on it. It returns `OK\n` if the package wasn't indexed.
* For `QUERY` commands, the server returns `OK\n` if the package is indexed. It returns `FAIL\n` if the package isn't indexed.
* If the server doesn't recognize the command or if there's any problem with the message sent by the client it should return `ERROR\n`.

# Package Server

Utilizes a simple graph to index packages. TCP Server can be configured using the ENV_VAR `PACKAGE_PORT`
to listen for connections, defaults to 8080. Utilizes channels and go-routines to handle long living connections. Connections
that do not write will be closed by the default amount of time (10seconds), can be configured using ENV_VAR `PACKAGE_CONNECTION_TIMEOUT`. The TCP Server
will attempt to exit cleanly upon receiving SIGINT or SIGTERM.

To use a different backend you must have a type that adheres to:

```go
type Backend interface {
	Exists(name string) bool
	Add(name string, edges ...string) error
	Remove(name string) error
}
```

The repository manager will handle read/write locking of the backend so you do not to directly make it go routine safe. The 
backend must be supplied to the server as part of a configuration.

##Add additional commands
To add additional commands you must adhere to the the interface in repomanager.go:
```go
type Operation interface {
	Run() (string, error)
	GetCommand() string
}
```
And update:
`func validateAndCreateOperator(instruction *instruction, repo *Repo) Operation {}` which returns an Operator type
to perform the corresponding action. Read and write locks should be use whenever interacting with the backend type.

Example with Query:
```go
const QUERY string = "QUERY"

type QueryOperator struct {
	instruction *instruction
	repo        *Repo
}

func NewQueryOperator(instruction *instruction, repo *Repo) *QueryOperator {
	return &QueryOperator{
		instruction: instruction,
		repo:        repo,
	}
}

func (o QueryOperator) Run() (string, error) {

	o.repo.mu.RLock()
	exists := o.repo.backend.Exists(o.instruction.packageName)
	o.repo.mu.RUnlock()

	if exists {
		return OK, nil
	}
	return FAIL, nil

}

func (o QueryOperator) GetCommand() string {
	return o.instruction.cmd
}

```


## Build
```
make
```

## Test Unit + Integration
```
make test
```

## Docker
```
make build.docker
```

## Built with
 - Golang 1.8.3
 - Docker 17.06.0-ce
 - golang:1.8.3-alpine (Docker build step)
 - alpine:3.6 (Docker artifact step)