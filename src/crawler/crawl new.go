package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
)

/*
	Crawler
*/
func CrawlNew(uri string, startHost string, queue chan<- string, client *http.Client, mutex *sync.RWMutex, wg *sync.WaitGroup) {
	mutex.Lock()
	if visited[uri] == false {
		visited[uri] = true
	} else {
		fmt.Println("VISITED ERROR")
		return
	}
	mutex.Unlock()

	defer wg.Done()

	resp, err := client.Get(uri)
	// if uri can't be resolved
	if err != nil {
		fmt.Println("URI ERROR: ", err)
		//mutex.Lock()
		//visited[uri] = false
		mutex.Unlock()
		return
	}
	defer resp.Body.Close()
	//CheckSitegroup(resp.Body)
	links := collectLinks(resp.Body)
	//fmt.Println(links)
	fmt.Printf("Status: %-10s Links: %-5d URL: %-10s Go: %-10d \n", resp.Status, len(links), uri, runtime.NumGoroutine())

	// check found Links for bad URLs, email adresses, Javascript, same Domain etc
	// if Link check doesn't fail it gets added to the queue
	for _, link := range links {
		//fmt.Println("LINK: ", link)
		absoluteUrl := FixUrl(&link, &uri)
		//fmt.Println(absoluteUrl)
		//fmt.Println("ABSOLUTE LINK: ", absoluteUrl)
		if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl, &startHost) {
			//fmt.Println(absoluteUrl)
			mutex.Lock()
			if visited[absoluteUrl] == false {
				queue <- absoluteUrl
			}
			mutex.Unlock()
			//temp := visited[absoluteUrl]
			//mutex.Unlock()
			//if temp == false {
			//fmt.Println("VALID")
			//mutex.Lock()
			//visitsOpen++
			//mutex.Unlock()
			//queue <- absoluteUrl
			/* wg.Add(1)
			go func() {
				queue <- absoluteUrl
				wg.Done()
				return
			}()
			*/
			//}
		}

	}
	return
}
