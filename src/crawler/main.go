package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
)

var workers = 1
var input = make(chan string, 10000)
var output = make(chan string, 1)
var visited = make(map[string]bool)
var counter = 1
var throttle = time.Tick(5000 * time.Millisecond)

/*
	Start
*/

func main() {

	start := time.Now()

	// Waitgroup to know when all Goroutines are closed
	var wg sync.WaitGroup

	startPage := "http://www.ebay.de/"
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host
	fmt.Printf("Crawling %s @ Host %s \n", startUrl, startHost)

	wg.Add(1)
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
					fmt.Printf("%-3d # %s \n", counter, link)
					counter++
					input <- link
					wg.Done()
				} else {
					wg.Done()
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(input)
		elapsed := time.Since(start)
		fmt.Printf("\n%d links in %f seconds\n", len(visited), elapsed.Seconds())
		os.Exit(0)
		//ExportToCSV(startHost, visited)
		//fmt.Printf("\nCSV file created for %s\n", startHost)
	}()

	// keep console open
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func worker(i int, startHost string, wg *sync.WaitGroup, input, output chan string) {
	//fmt.Printf("Worker started: %d \n", i)
	for {
		select {
		case link := <-input:
			//<-throttle
			wg.Add(1)
			Crawl(link, startHost, wg, input, output)
		}
	}
}
