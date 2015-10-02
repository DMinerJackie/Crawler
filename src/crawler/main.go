package main

import (
	"bufio"
	"flag"
	"log"
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
var LinkCounter int32 = 0
var ErrCounter int32 = 0
var CounterA int32 = 0
var CounterB int32 = 0

/*
	FLAG PARAMETER
*/
var linkPtr = flag.String("url", "http://www.golem.de/", "webpage")
var workersPtr = flag.Int("con", 500, "connections")
var logLevelPtr = flag.Int("log", 1, "log level")
var cpuprofilePtr = flag.Bool("cpu", true, "cpu profiling")

/*
	MAIN START
*/
func main() {
	flag.Parse()
	link := *linkPtr
	workers := *workersPtr
	logLevel := *logLevelPtr
	cpuprofile := *cpuprofilePtr

	/*
		START URL
	*/
	startPage := link
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host

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
	if logLevel != -1 {
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
		go worker(startHost, mutex)
		//Debug.Printf("worker %d created", i)
	}

	/*
		START CRAWLING LOOP
	*/
	Ever.Printf("START \n %s @ %d worker(s) @ loglevel %d", startHost, workers, logLevel)
	AddLinkCount()
	Info.Printf(" Counter: %d @ %s \n", GetLinkCount(), startPage)
	visited[startPage] = true
	AddCountA()
	Crawl(startPage, startHost, mutex)

	//keep console open
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

/*
	WORKER FUNCTION
*/
func worker(startHost string, mutex *sync.Mutex) {
	for {
		select {
		case link := <-new_links_chan:
			Debug.Printf("consumed from chan: %s \n", link)
			AddCountA()
			Crawl(link, startHost, mutex)
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
func AddCountB() {
	atomic.AddInt32(&CounterB, 1)
}
func DoneCountA() {
	atomic.AddInt32(&CounterA, -1)
	if atomic.LoadInt32(&CounterA) == 0 && atomic.LoadInt32(&CounterB) == 0 {
		Close()
	}
}
func DoneCountB() {
	atomic.AddInt32(&CounterB, -1)
	if atomic.LoadInt32(&CounterA) == 0 && atomic.LoadInt32(&CounterB) == 0 {
		Close()
	}
}

/*
	CLOSE FUNCTION WHEN FINISHED
*/
func Close() {
	elapsed := time.Since(start)
	Ever.Printf("STOP \n %d link(s) : %d error(s) : %f seconds \n\n", GetLinkCount(), GetErrCount(), elapsed.Seconds())
	os.Exit(0)
}
