package main

import (
	"net/http"
	"sync"
)

func Crawl(link string, startHost string, mutex *sync.Mutex) {
	defer DoneCountA()
	Debug.Printf("Begin crawl: %s \n", link)

	client := http.Client{}

	resp, err := client.Get(link)
	if err != nil {
		AddErrCount()
		Error.Printf("  Connection Error for %s : %s \n", link, err)
		return
	}
	defer resp.Body.Close()

	links := collectLinks(&link, resp.Body)
	Debug.Printf("  Found %d links on %s \n", len(links), link)

	for _, foundLink := range links {
		absoluteUrl := FixUrl(&foundLink, &link)
		if absoluteUrl != "" {
			if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl, &startHost) {
				Debug.Printf("     *** Tests passed: %s \n", absoluteUrl)
				mutex.Lock()
				if visited[absoluteUrl] == false {
					visited[absoluteUrl] = true
					AddLinkCount()
					Info.Printf(" Counter: %d @ %s \n", GetLinkCount(), absoluteUrl)
					mutex.Unlock()
					Debug.Printf("added to channel: %s \n", absoluteUrl)
					AddCountB()
					new_links_chan <- absoluteUrl
				} else {
					mutex.Unlock()
					Debug.Printf("DUPLICATE VISIT: %s \n", absoluteUrl)
				}
			} else {
				Debug.Printf("  - CheckUrl not passed: %s \n", absoluteUrl)
			}
		} else {
			Debug.Printf("  - AbsoluteUrl not passed: %s \n", absoluteUrl)

		}

	}
}
