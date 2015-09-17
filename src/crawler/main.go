package main

import (
	//"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var visited = make(map[string]bool)
var visitsDone = 0
var visitsOpen = 0
var mutex = &sync.RWMutex{}
var throttle = time.Tick(1 * time.Millisecond)
var counter = 0

var useCpuProfile = true
var useRamProfile = false
var useThrottle = false

/*
	Start
*/

func main() {
	// CPU Profiler unter go tool pprof cpu.prof
	if useCpuProfile {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

	}

	// Memory Profiler unter http://localhost:6060/debug/pprof/
	if useRamProfile {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	start := time.Now()
	/*
		Read parameter and set URL for args[0]
	*/
	flag.Parse()

	/*
			args := flag.Args()
			baseUri := "http://www." + args[0]
			if len(args) < 1 {
		    fmt.Println("Please specify start page")  // if a starting page wasn't provided as an argument
		    os.Exit(1)                                // show a message and exit.
		  }                                           // Note that 'main' doesn't return anything.
	*/

	transport := &http.Transport{}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	// Waitgroup to know when all Goroutines are closed
	var wg sync.WaitGroup
	// Channel for the URLs
	queue := make(chan string)

	/*
		URL that gets fetched first
	*/
	startPage := "http://www.example.de/"
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host
	fmt.Printf("Crawling %s @ Host %s \n", startUrl, startHost)

	wg.Add(1)
	go Crawl(startPage, startHost, queue, &client, mutex, &wg)

	go func() {
		wg.Wait()
		close(queue)
	}()

	// Queue that spawns Goroutines for every URL in the queue
	for uri := range queue {
		if runtime.NumGoroutine() < 30 {
			wg.Add(1)
			go func() {
				<-throttle
				Crawl(uri, startHost, queue, &client, mutex, &wg)
				return
			}()
		}
	}

	// write fetched URLs to CSV file
	ExportToCSV(startHost, visited)
	fmt.Printf("\nCSV file created for %s\n", startHost)

	elapsed := time.Since(start)
	fmt.Printf("\n%d links in %f seconds\n", len(visited), elapsed.Seconds())

	// keep console open
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
}
