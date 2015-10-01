package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var new_links_chan = make(chan string, 1000000)
var visited = make(map[string]bool)
var counter = 0
var errcounter = 0
var mutexCounter = 0
var throttle = time.Tick(100 * time.Millisecond)
var mutex = &sync.Mutex{}
var mutexCount = &sync.Mutex{}
var wg1 = &sync.WaitGroup{}
var wg2 = &sync.WaitGroup{}

func MutexAdd() {
	mutexCount.Lock()
	mutexCounter++
	mutexCount.Unlock()
}

func MutexDone() {
	mutexCount.Lock()
	mutexCounter--
	if mutexCounter == 0 {
		mutexCount.Unlock()
		Info.Println("CLOSED")
		os.Exit(0)
	}
	mutexCount.Unlock()
}

/*
	MAIN
*/

func main() {

	/*
		FLAG PARAMETER
	*/
	linkPtr := flag.String("url", "example.de/", "site")
	numbWorkerPtr := flag.Int("con", 30, "connections")
	logLevelPtr := flag.Int("log", 2, "0-4")
	cpuprofilePtr := flag.String("cpu", "profile", "write cpu profile to file")

	flag.Parse()

	workers := *numbWorkerPtr
	link := *linkPtr
	logLevel := int32(*logLevelPtr)
	cpuprofile := *cpuprofilePtr

	/*
		CPU PROFILING
	*/

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*
		LOGGING
	*/
	file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error.Println("failed open file")
	}
	defer file.Close()

	/*
		START URL
	*/
	startPage := "http://www." + link
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host

	setLogLevel(logLevel, file)
	Info.Printf("Start: %s : %d workers : loglevel %d", startHost, workers, logLevel)

	Debug.Printf("added to chan: %s \n", startPage)
	Info.Printf("Counter: %-3d @ %s \n", counter, startUrl)
	counter++
	new_links_chan <- startPage

	/*
		CREATE WORKER
	*/
	for i := 1; i <= workers; i++ {
		go worker(startHost, mutex)
		Debug.Printf("worker %d created", i)
	}

}

func worker(startHost string, mutex *sync.Mutex) {
	for {
		select {
		case link := <-new_links_chan:
			<-throttle
			Debug.Printf("consumed from chan: %s \n", link)
			MutexAdd()
			Crawl(link, startHost, mutex)
		}
	}
}
