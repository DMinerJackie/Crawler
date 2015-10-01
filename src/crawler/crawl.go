package main

import (
	"net/http"
	"sync"
	"time"
	//"fmt"
)

func Crawl(link string, startHost string, mutex *sync.Mutex) {
	defer wg1.Done()
	Debug.Printf("Begin crawl: %s \n", link)

	transport := &http.Transport{}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	//	resp, err := client.head(link)
	//	if err != nil {
	//		fmt.printf("head error: %s %s \n", err, link)
	//		return
	//	}

	//	if strings.hasprefix(strings.tolower(resp.header.get("content-type")), "text/plain") {

	resp, err := client.Get(link)
	if err != nil {
		errcounter++
		Error.Printf("  Connection Error for %s : %s \n", link, err)
		return
	}
	defer resp.Body.Close()

	links := collectLinks(resp.Body)
	Debug.Printf("  Found %d links on %s \n", len(links), link)

	for i, foundLink := range links {
		Debug.Printf("  # Begin range %d for %s \n", i, foundLink)
		absoluteUrl := FixUrl(&foundLink, &link)
		Debug.Printf("  + FixUrl %d is: %s \n", i, absoluteUrl)
		if absoluteUrl != "" {
			Debug.Printf("  + AbsoluteUrl %d is : %s \n", i, absoluteUrl)
			if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl, &startHost) {
				Debug.Printf("  + CheckUrl %d is okay: %s \n", i, absoluteUrl)
				Debug.Printf("     *** Tests for item %d passed: %s \n", i, absoluteUrl)
				mutex.Lock()
				Debug.Printf("STATUS visit:%t for %s \n", visited[absoluteUrl], absoluteUrl)
				if visited[absoluteUrl] == false {
					Debug.Printf("SET %s to TRUE \n", absoluteUrl)
					visited[absoluteUrl] = true
					Info.Printf("Counter: %-3d @ %s \n", counter, absoluteUrl)
					//fmt.Println(absoluteUrl)
					counter++
					mutex.Unlock()

					//					if counter%1000 == 0 {
					//						//Info.Printf("Crawled: %-5d", counter)
					//					}
					wg2.Add(1)
					new_links_chan <- absoluteUrl
					Debug.Printf("added to channel: %s \n", absoluteUrl)
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
	//	} else {
	//		return
	//	}
}
