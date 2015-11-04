package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

var (
	Pure    *log.Logger
	Ever    *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

// turnOnLogging configures the logging writers.
func setLogLevel(logLevel int32, fileHandle io.Writer) {
	pureHandle := ioutil.Discard
	everHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard
	debugHandle := ioutil.Discard

	if logLevel == 1 {
		everHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 2 {
		everHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel == 3 {
		everHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
		debugHandle = os.Stdout
	}

	if logLevel == 4 {
		everHandle = os.Stdout
		pureHandle = os.Stdout
	}

	if logLevel == 5 {
		everHandle = os.Stdout
	}

	if fileHandle != nil && logLevel != -1 {
		if pureHandle == ioutil.Discard {
			pureHandle = io.MultiWriter(fileHandle, pureHandle)
		}
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

	Pure = log.New(pureHandle, "", 0)
	Ever = log.New(everHandle, "", log.Ldate|log.Ltime)
	Debug = log.New(debugHandle, "", log.Ldate|log.Ltime)
	Info = log.New(infoHandle, "", log.Ldate|log.Ltime)
	Warning = log.New(warnHandle, "WARNING: ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "", log.Ldate|log.Ltime)

	atomic.StoreInt32(&logLevel, logLevel)
}
