package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func Crawl(link string, startHost string, wg *sync.WaitGroup, input, output chan string) {
	defer wg.Done()

	counter := 0
	//fmt.Printf("Unique: %-15d Goroutines %-5d \n", len(visited), runtime.NumGoroutine())

	transport := &http.Transport{}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	resp, err := client.Head(link)
	if err != nil {
		fmt.Println("HEAD ERROR: ", err)
		return
	}

	if strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")), "text") {

		resp, err := client.Get(link)
		if err != nil {
			fmt.Println("LINK ERROR: ", err)
			return
		}

		defer resp.Body.Close()

		links := collectLinks(resp.Body)
		//fmt.Printf("%-3d For Site: %s \n", len(links), link)

		for _, foundLink := range links {
			absoluteUrl := FixUrl(&foundLink, &link)
			if CheckHost(&absoluteUrl, &startHost) && CheckUrl(&absoluteUrl) {
				counter++
				//fmt.Printf("TOTAL: %-3d FROM: %-20s INPUT: %-20s \n", len(links), link, absoluteUrl)
				wg.Add(1)
				output <- absoluteUrl
			}
		}
		//fmt.Printf("T: %-3d R: %-3d AT: %s \n", len(links), counter, link)
		return
	}
}
