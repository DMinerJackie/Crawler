package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	//"runtime"
	"sync"
)

var absoluteCounter = 0

/*
	Crawler
*/
func Crawl(uri string, startHost string, queue chan<- string, client *http.Client, mutex *sync.RWMutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mutex.Lock()
	if visited[uri] == false {
		visited[uri] = true
		//fmt.Printf("Unique: %-15d Goroutines %-5d \n", len(visited), runtime.NumGoroutine())
		mutex.Unlock()

		resp, err := client.Get(uri)
		// if uri can't be resolved
		if err != nil {
			fmt.Println("URI    ERROR: ", err)
			return
		}

		defer resp.Body.Close()

		links := collectLinks(resp.Body)

		// check found Links for bad URLs, email adresses, Javascript, same Domain etc
		// if Link check doesn't fail it gets added to the queue
		for _, link := range links {
			absoluteUrl := FixUrl(&link, &uri)
			if CheckHost(&absoluteUrl, &startHost) && CheckUrl(&absoluteUrl) {
				mutex.Lock()
				if visited[absoluteUrl] == false {
					fmt.Println("AbsoluteURL: " + absoluteUrl)
					queue <- absoluteUrl
				}
				mutex.Unlock()
			}
		}
	} else {
		//fmt.Println("DOUBLE: " + uri)
		mutex.Unlock()
		return
	}
	return
}
