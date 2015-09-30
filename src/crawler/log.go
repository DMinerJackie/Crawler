package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// turnOnLogging configures the logging writers.
func setLogLevel(logLevel int32, fileHandle io.Writer) {
	debugHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard

	if logLevel == 1 {
		debugHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 2 {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 3 {
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 4 {
		errorHandle = os.Stderr
	}

	if fileHandle != nil && logLevel != 0 {
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

	Debug = log.New(debugHandle, "DEBUG: ", log.Ltime)
	Info = log.New(infoHandle, "INFO: ", log.Ltime)
	Warning = log.New(warnHandle, "WARNING: ", log.Ltime)
	Error = log.New(errorHandle, "ERROR: ", log.Ltime)

	atomic.StoreInt32(&logLevel, logLevel)
}
