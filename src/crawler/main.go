package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"
)

/*
	GLOBAL PARAMETER
*/
var new_links_chan = make(chan string, 1000000)
var visited = make(map[string]bool)
var mutex = &sync.Mutex{}
var start = time.Now()
var multithreaded bool
var LinkCounter int32 = 0
var ErrCounter int32 = 0
var CounterA int32 = 0
var CounterB int32 = 0
var client = &http.Client{}
var startPage string
var startHost string

/*
	FLAG PARAMETER
*/
var linkPtr = flag.String("url", "http://www.example.de/", "webpage")
var workersPtr = flag.Int("con", 1, "connections")
var logLevelPtr = flag.Int("lvl", 1, "log level")
var logFilePtr = flag.Bool("log", true, "log file")
var cpuprofilePtr = flag.Bool("cpu", false, "cpu profiling")
var multiPtr = flag.Bool("exp", false, "experimental")

/*
@@@	MAIN START @@@
*/
func main() {
	flag.Parse()
	startPage = *linkPtr
	workers := *workersPtr
	logLevel := *logLevelPtr
	logFile := *logFilePtr
	cpuprofile := *cpuprofilePtr
	multithreaded = *multiPtr

	/*
		START URL
	*/
	startUrl, err := url.Parse(startPage)
	if err != nil {
		Error.Println(err)
		os.Exit(1)
	}
	startHost = startUrl.Host

	/*
		CPU PROFILING
	*/
	if cpuprofile == true {
		f, err := os.Create("bench.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*
		SET LOGGING FILE + LOGGING LEVEL
	*/
	if logFile == true {
		file, err := os.OpenFile(startHost+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			Error.Println("failed open file")
		}
		defer file.Close()
		setLogLevel(int32(logLevel), file)
	} else {
		setLogLevel(int32(logLevel), nil)
	}

	/*
		CREATE WORKER
	*/
	for i := 1; i <= workers; i++ {
		go worker(startHost)
	}

	/*
		START CRAWLING LOOP
	*/
	Ever.Printf("START \n %s @ %d worker(s) @ loglevel %d", startHost, workers, logLevel)
	AddLinkCount()
	Info.Printf(" Counter: %d @ %s \n", GetLinkCount(), startPage)
	visited[startPage] = true
	AddCountA()
	Crawl(startPage)

	//keep console open
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

/*
	WORKER FUNCTION
*/
func worker(startHost string) {
	for {
		select {
		case link := <-new_links_chan:
			Debug.Printf("consumed from chan: %s \n", link)
			AddCountA()
			Crawl(link)
			DoneCountB()
		}
	}
}

/*
	Atomic Add Counter for logging
*/
func AddLinkCount() {
	atomic.AddInt32(&LinkCounter, 1)
}
func GetLinkCount() int32 {
	return atomic.LoadInt32(&LinkCounter)
}
func AddErrCount() {
	atomic.AddInt32(&ErrCounter, 1)
}
func GetErrCount() int32 {
	return atomic.LoadInt32(&ErrCounter)
}

/*
	Atomic Add & Decrease Counter to test if the crawler has finished
*/
func AddCountA() {
	atomic.AddInt32(&CounterA, 1)
}
func AddCountB(x int) {
	atomic.AddInt32(&CounterB, int32(x))
}
func DoneCountA() {
	atomic.AddInt32(&CounterA, -1)
	if atomic.LoadInt32(&CounterA) == 0 && atomic.LoadInt32(&CounterB) == 0 {
		Close()
	}
}
func DoneCountB() {
	atomic.AddInt32(&CounterB, -1)
	if atomic.LoadInt32(&CounterB) == 0 && atomic.LoadInt32(&CounterA) == 0 {
		Close()
	}
}

/*
	CLOSE FUNCTION WHEN FINISHED
*/
func Close() {
	elapsed := time.Since(start).Seconds()
	timeType := "second(s)"
	if elapsed > 60 {
		elapsed = elapsed / 60
		timeType = "minute(s)"
	}
	Ever.Printf("STOP \n %d link(s) : %d error(s) : %f %s \n\n", GetLinkCount(), GetErrCount(), elapsed, timeType)
	os.Exit(0)
}
