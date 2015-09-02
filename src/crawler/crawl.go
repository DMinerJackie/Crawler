package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/jackdanger/collectlinks"
)

/*
	Crawler
*/
func Crawl(uri string, startHost string, queue chan<- string, client *http.Client, mutex *sync.RWMutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mutex.Lock()
	visited[uri] = true

	/*	visitsDone++
		if visitsOpen > 0 {	visitsOpen-- }
	*/
	fmt.Printf("Unique: %-15d Goroutines %-5d \n", len(visited), runtime.NumGoroutine())

	mutex.Unlock()

	resp, err := client.Get(uri)
	// if uri can't be resolved
	if err != nil {
		//fmt.Println("URI ERROR: ", err)
		return
	}

	defer resp.Body.Close()

	/*
	// Foo should ideally be a struct that matches structure of that json
	type CampaignID struct {
		x int
	}
	var campID CampaignID
	err := json.Unmarshal([]byte(campaignID), &campID)
*/
	
	//CheckSitegroup(resp.Body)
	links := collectlinks.All(resp.Body)

	// check found Links for bad URLs, email adresses, Javascript, same Domain etc
	// if Link check doesn't fail it gets added to the queue
	for _, link := range links {
		//fmt.Println("LINK: ", link)
		absoluteUrl := FixUrl(&link, &uri)
		//fmt.Println("ABSOLUTE LINK: ", absoluteUrl)
		if uri != "" && CheckUrl(&absoluteUrl) && CheckTLD(&absoluteUrl, &startHost) {
			//fmt.Println("VALID")
			mutex.RLock()
			temp := visited[absoluteUrl]
			mutex.RUnlock()
			if temp == false {
				//fmt.Println("VALID")
				//mutex.Lock()
				//visitsOpen++
				//mutex.Unlock()
				queue <- absoluteUrl
				/* wg.Add(1)
				go func() {
					queue <- absoluteUrl
					wg.Done()
					return
				}()
				*/
			}
		}

	}
	return
}
