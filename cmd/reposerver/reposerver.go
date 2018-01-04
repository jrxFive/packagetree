package main

import (
	"github.com/jrxfive/packagetree/pkg/graph"
	"github.com/jrxfive/packagetree/pkg/logging"
	"github.com/jrxfive/packagetree/pkg/repomanager"
	"github.com/jrxfive/packagetree/pkg/server"
	"os"
	"strconv"
)

var PORT int = 8080
var MAX_HERD int = 10
var CONNECTION_TIMEOUT = 10
var logger = logging.GetLogger()

func init() {
	if envPort, ok := os.LookupEnv("PACKAGE_PORT"); ok {
		value, err := strconv.Atoi(envPort)
		if err != nil {

		} else {
			PORT = value
		}
	}

	if envMaxHerd, ok := os.LookupEnv("PACKAGE_MAX_HERD"); ok {
		value, err := strconv.Atoi(envMaxHerd)
		if err != nil {

		} else {
			MAX_HERD = value
		}
	}

	if envConnTimeout, ok := os.LookupEnv("PACKAGE_CONNECTION_TIMEOUT"); ok {
		value, err := strconv.Atoi(envConnTimeout)
		if err != nil {

		} else {
			CONNECTION_TIMEOUT = value
		}
	}

	logger.Printf("Starting new server on port:%v\n", PORT)
	logger.Printf("MAX_HERD set to:%v\n", MAX_HERD)
	logger.Printf("Connection timeout set to:%v\n", CONNECTION_TIMEOUT)
}

func main() {
	g, err := graph.NewGraph()
	if err != nil {
		os.Exit(1)
	}

	repo := repomanager.NewRepo(g)
	serverConfiguration := server.NewServerConfiguration(repo, PORT, MAX_HERD, CONNECTION_TIMEOUT)

	err = server.Listen(serverConfiguration)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
