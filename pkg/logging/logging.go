package logging

import (
	"log"
	"os"
)

var logger = &log.Logger{}
var instantiatedCheck bool = false

func GetLogger() *log.Logger {

	if instantiatedCheck == false {
		instantiatedCheck = true
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
		return logger
	} else {
		return logger
	}
}
