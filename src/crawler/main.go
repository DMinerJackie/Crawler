package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"
)

var workers = 10
var input = make(chan string, 100000)
var output = make(chan string, 100000)
var visited = make(map[string]bool)

/*
	Start
*/

func main() {
	//visited[`https://www.roller.de/mhs-order-status`] = true

	start := time.Now()

	// Waitgroup to know when all Goroutines are closed
	var wg sync.WaitGroup

	startPage := "http://www.roller.de/"
	startUrl, _ := url.Parse(startPage)
	startHost := startUrl.Host
	fmt.Printf("Crawling %s @ Host %s \n", startUrl, startHost)

	// Create the number of workers
	for i := 0; i < workers; i++ {
		//fmt.Printf("Worker created: %d \n", i)
		//wg.Add(1)
		go worker(i, startHost, &wg, input, output)
	}

	wg.Add(1)
	go Crawl(startPage, startHost, &wg, input, output)

	go func() {
		for link := range output {
			if visited[link] == false {
				visited[link] = true
				//fmt.Println("OUTPUT: " + link)
				input <- link
			}
		}
	}()

	go func() {
		wg.Wait()
		//close(input)
		elapsed := time.Since(start)
		fmt.Printf("\n%d links in %f seconds\n", len(visited), elapsed.Seconds())
		//os.Exit(0)
		//ExportToCSV(startHost, visited)
		//fmt.Printf("\nCSV file created for %s\n", startHost)
	}()

	// keep console open
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func worker(i int, startHost string, wg *sync.WaitGroup, input, output chan string) {
	//fmt.Printf("Worker started: %d \n", i)
	for link := range input {
		wg.Add(1)
		Crawl(link, startHost, wg, input, output)
	}
}
