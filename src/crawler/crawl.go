package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func Crawl(link string, startHost string, wg *sync.WaitGroup, input, output chan string) {
	defer wg.Done()
	//fmt.Println("For Site: " + link)
	//fmt.Printf("Unique: %-15d Goroutines %-5d \n", len(visited), runtime.NumGoroutine())

	transport := &http.Transport{}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	resp, err := client.Get(link)
	if err != nil {
		fmt.Println("LINK ERROR: ", err)
		return
	}

	defer resp.Body.Close()

	links := collectLinks(resp.Body)

	for _, foundLink := range links {
		absoluteUrl := FixUrl(&foundLink, &link)
		if CheckHost(&absoluteUrl, &startHost) && CheckUrl(&absoluteUrl) {
			fmt.Printf("TOTAL: %-3d FROM: %-20s INPUT: %-20s \n", len(links), link, absoluteUrl)
			output <- absoluteUrl
		}
	}
	return
}
