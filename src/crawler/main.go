package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
)

var input = make(chan string, 1000000)
var output = make(chan string, 1000000)
var visited = make(map[string]bool)
var counter = 0
var errcounter = 0
var throttle = time.Tick(100 * time.Millisecond)

/*
	Start
*/

func main() {
	start := time.Now()

	linkPtr := flag.String("url", "example.de", "site")
	numbWorkerPtr := flag.Int("w", 1, "connections")
	logPtr := flag.Int("log", 1, "0-4")

	flag.Parse()

	workers := *numbWorkerPtr
	link := *linkPtr
	logLevel := int32(*logPtr)

	// LOGGING
	file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error.Println("failed open file")
	}
	defer file.Close()

	// Waitgroup to know when all Goroutines are closed
	var wg sync.WaitGroup

	startPage := "http://www." + link
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host

	setLogLevel(logLevel, file)
	Info.Printf("Start: %s : %d workers : loglevel %d", startHost, workers, logLevel)

	output <- startPage

	// Create the number of workers
	for i := 0; i < workers; i++ {
		go worker(i, startHost, &wg, input, output)
	}

	go func() {
		for {
			select {
			case link := <-output:
				if visited[link] == false {
					visited[link] = true
					counter++
					//fmt.Printf("%-3d # %s \n", counter, link)
					if counter%1000 == 0 {
						//Info.Printf("Crawled: %-5d", counter)
					}
					input <- link
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		fmt.Println("CLOSED")
		close(input)
		close(output)
		elapsed := time.Since(start)
		Info.Printf("Stop: %d visited: %d failed: %f seconds", counter, errcounter, elapsed.Seconds())
		os.Exit(0)
		//ExportToCSV(startHost, visited)
		//fmt.Printf("\nCSV file created for %s\n", startHost)
	}()

	// keep console open
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func worker(i int, startHost string, wg *sync.WaitGroup, input, output chan string) {
	for {
		select {
		case link := <-input:
			wg.Add(1)
			<-throttle
			Crawl(link, startHost, wg, input, output)
		}
	}
}
