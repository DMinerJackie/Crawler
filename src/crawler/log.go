package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

var (
	Ever    *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

// turnOnLogging configures the logging writers.
func setLogLevel(logLevel int32, fileHandle io.Writer) {
	everHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard
	debugHandle := ioutil.Discard
	everHandle = os.Stdout

	if logLevel == 1 {
		infoHandle = os.Stdout
	}

	if logLevel == 2 {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
	}

	if logLevel == 3 {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 4 {
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 5 {
		errorHandle = os.Stderr
	}

	if logLevel == 6 {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
		debugHandle = os.Stdout
	}

	if fileHandle != nil && logLevel != -1 {
		if everHandle == os.Stdout {
			everHandle = io.MultiWriter(fileHandle, everHandle)
		}

		if debugHandle == os.Stdout {
			debugHandle = io.MultiWriter(fileHandle, debugHandle)
		}

		if infoHandle == os.Stdout {
			infoHandle = io.MultiWriter(fileHandle, infoHandle)
		}

		if warnHandle == os.Stdout {
			warnHandle = io.MultiWriter(fileHandle, warnHandle)
		}

		if errorHandle == os.Stderr {
			errorHandle = io.MultiWriter(fileHandle, errorHandle)
		}
	}

	Ever = log.New(everHandle, "LOG: ", log.Ldate|log.Ltime)
	Debug = log.New(debugHandle, "DEBUG: ", log.Ldate|log.Ltime)
	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime)
	Warning = log.New(warnHandle, "WARNING: ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime)

	atomic.StoreInt32(&logLevel, logLevel)
}
