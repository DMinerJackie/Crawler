package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

/*
	GLOBAL PARAMETER
*/
var new_links_chan = make(chan string, 1000000)
var visited = make(map[string]bool)
var counter = 0
var mutexErrorCounter = 0
var mutexCounter1 = 0
var mutexCounter2 = 0
var mutex = &sync.Mutex{}
var mutexCount1 = &sync.Mutex{}
var mutexCount2 = &sync.Mutex{}
var mutexErrorCount = &sync.Mutex{}
var start = time.Now()

/*
	MAIN START
*/
func main() {

	/*
		FLAG PARAMETER
	*/
	linkPtr := flag.String("url", "roller.com", "site")
	numberOfWorkersPtr := flag.Int("con", 30, "connections")
	logLevelPtr := flag.Int("log", 2, "0-4")
	cpuprofilePtr := flag.String("cpu", "profile", "write cpu profile to file")

	flag.Parse()

	workers := *numberOfWorkersPtr
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
	startPage := "http://www." + link + "/"
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host

	setLogLevel(logLevel, file)
	//new_links_chan <- startPage

	/*
		CREATE WORKER
	*/
	for i := 1; i <= workers; i++ {
		go worker(startHost, mutex)
		Debug.Printf("worker %d created", i)
	}
	Info.Printf("Start: %s : %d workers : loglevel %d", startHost, workers, logLevel)
	MutexAdd1()
	counter++
	Info.Printf("Counter: %-3d @ %s \n", counter, startPage)
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
			MutexAdd1()
			Crawl(link, startHost, mutex)
			MutexDone2()
		}
	}
}

/*
	MUTEX ADD + MUTEX DONE with check if mutexCounter == 0
*/
func MutexAdd1() {
	mutexCount1.Lock()
	mutexCounter1++
	mutexCount1.Unlock()
}

func MutexAdd2() {
	mutexCount2.Lock()
	mutexCounter2++
	mutexCount2.Unlock()
}

func MutexErrorAdd() {
	mutexErrorCount.Lock()
	mutexErrorCounter++
	mutexErrorCount.Unlock()
}

func MutexDone1() {
	mutexCount1.Lock()
	mutexCount2.Lock()
	mutexCounter1--
	if mutexCounter1 == 0 && mutexCounter2 == 0 {
		mutexCount1.Unlock()
		mutexCount2.Unlock()
		Close()
	}
	mutexCount1.Unlock()
	mutexCount2.Unlock()
}

func MutexDone2() {
	mutexCount1.Lock()
	mutexCount2.Lock()
	mutexCounter2--
	if mutexCounter1 == 0 && mutexCounter2 == 0 {
		mutexCount1.Unlock()
		mutexCount2.Unlock()
		Close()
	}
	mutexCount1.Unlock()
	mutexCount2.Unlock()
}

func Close() {
	//close(new_links_chan)
	Info.Println("CLOSED")
	elapsed := time.Since(start)
	Info.Printf("%d link(s) : %d error(s) : %f seconds\n", counter, mutexErrorCounter, elapsed.Seconds())
	os.Exit(0)
}
