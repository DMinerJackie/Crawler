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

var new_links_chan = make(chan string, 1000000)
var visited = make(map[string]bool)
var counter = 0
var errcounter = 0
var throttle = time.Tick(100 * time.Millisecond)
var mutex = &sync.Mutex{}
var wg1 = &sync.WaitGroup{}
var wg2 = &sync.WaitGroup{}

/*
	Start
*/

func main() {
	start := time.Now()

	// console parameter
	linkPtr := flag.String("url", "example.de/", "site")
	numbWorkerPtr := flag.Int("con", 30, "connections")
	logLevelPtr := flag.Int("log", 1, "0-4")
	cpuprofilePtr := flag.String("cpu", "profile", "write cpu profile to file")

	flag.Parse()

	workers := *numbWorkerPtr
	link := *linkPtr
	logLevel := int32(*logLevelPtr)
	cpuprofile := *cpuprofilePtr

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// LOGGING
	file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error.Println("failed open file")
	}
	defer file.Close()

	// Waitgroup to know when all Goroutines are closed

	startPage := "http://www." + link
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host

	setLogLevel(logLevel, file)
	Info.Printf("Start: %s : %d workers : loglevel %d", startHost, workers, logLevel)

	Debug.Printf("added to chan: %s \n", startPage)
	Info.Printf("Counter: %-3d @ %s \n", counter, startUrl)
	counter++
	wg2.Add(1)
	new_links_chan <- startPage

	// Create the number of workers
	for i := 1; i <= workers; i++ {
		go worker(startHost, mutex)
		Debug.Printf("worker %d created", i)
	}

	go func() {
		wg1.Wait()
		wg2.Wait()
		close(new_links_chan)
		Info.Println("CLOSED")
		elapsed := time.Since(start)
		Info.Printf("Stop: %d visited: %d failed: %f seconds", counter, errcounter, elapsed.Seconds())
		os.Exit(0)
		//ExportToCSV(startHost, visited)
		//fmt.Printf("\nCSV file created for %s\n", startHost)
	}()

	//keep console open
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func worker(startHost string, mutex *sync.Mutex) {
	for {
		select {
		case link := <-new_links_chan:
			<-throttle
			Debug.Printf("consumed from chan: %s \n", link)
			wg1.Add(1)
			Crawl(link, startHost, mutex)
			wg2.Done()
		}
	}
}
