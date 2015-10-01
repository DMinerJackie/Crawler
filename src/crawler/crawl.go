package main

import (
	"net/http"
	"sync"
	"time"
)

func Crawl(link string, startHost string, mutex *sync.Mutex) {
	defer MutexDone1()
	Debug.Printf("Begin crawl: %s \n", link)

	transport := &http.Transport{}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	resp, err := client.Get(link)
	if err != nil {
		MutexErrorAdd()
		Error.Printf("  Connection Error for %s : %s \n", link, err)
		return
	}
	defer resp.Body.Close()

	links := collectLinks(resp.Body)
	Debug.Printf("  Found %d links on %s \n", len(links), link)

	for i, foundLink := range links {
		absoluteUrl := FixUrl(&foundLink, &link)
		if absoluteUrl != "" {
			if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl, &startHost) {
				Debug.Printf("     *** Tests for item %d passed: %s \n", i, absoluteUrl)
				mutex.Lock()
				if visited[absoluteUrl] == false {
					Debug.Printf("SET %s to TRUE \n", absoluteUrl)
					visited[absoluteUrl] = true
					Info.Printf("Counter: %-3d @ %s \n", counter, absoluteUrl)
					counter++
					mutex.Unlock()
					Debug.Printf("added to channel: %s \n", absoluteUrl)
					MutexAdd2()
					new_links_chan <- absoluteUrl
				} else {
					Debug.Printf("DUPLICATE VISIT: %s \n", absoluteUrl)
					mutex.Unlock()
				}
			} else {
				Debug.Printf("  - CheckUrl not passed: %s \n", absoluteUrl)
			}
		} else {
			Debug.Printf("  - AbsoluteUrl not passed: %s \n", absoluteUrl)

		}

	}
}
