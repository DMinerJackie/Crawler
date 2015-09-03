package main

import (
	//"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
)

var visited = make(map[string]bool)
var visitsDone = 0
var visitsOpen = 0
var mutex = &sync.RWMutex{}
var throttle = time.Tick(1000 * time.Millisecond)

/*
	Start
*/

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	/*
				Creates HTTP Client that skips SSL connections

				Timeout specifies a time limit for requests made by this Client
				The timeout includes connection time, any redirects, and reading the response body.
				The timer remains running after Get, Head, Post, or Do return and will interrupt
				reading of the Response.Body.
		    	A Timeout of zero means no timeout.
	*/
	timeout := time.Duration(5 * time.Second)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
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
	startPage := "http://golem.de"
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
		if runtime.NumGoroutine() < 150 {
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
	fmt.Println("\n CSV file created")

	elapsed := time.Since(start)
	fmt.Printf("\n%d links in %f seconds\n", len(visited), elapsed.Seconds())

	// keep console open
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
}