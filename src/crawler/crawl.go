package main

import (
	"net/http"
	"time"
)

var maxRetries = 4
var resp *http.Response = nil
var err error

func Crawl(link string, workerID int) {
	defer DoneCountA()
	Debug.Printf("Begin crawl: %s \n", link)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		Error.Printf("ERROR \t REQ Connection Error for workerID %d : %s : %s \n", workerID, link, err)
		AddErrCount()
		return
	}
	//Info.Println(req)
	//	UNIQUE USER AGENT, NO 'BOT'
	req.Header.Set("User-Agent", "WebSpider")
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		Error.Printf("ERROR \t RESP Connection Error for %s : %s \n", link, err)
	//		AddErrCount()
	//		return
	//	}

	// Retry requests if err != nil, sleep between each retry 'i' seconds
	for i := 0; i <= maxRetries; i++ {
		resp, err = client.Do(req)
		if err != nil {
			if i == maxRetries {
				Error.Printf("ERROR \t RESP Connection Error for workerID %d : %s : %s \n", workerID, link, err)
				AddErrCount()
				mutex.Lock()
				visited[link] = true 
				mutex.Unlock()
				return
			} else {
				time.Sleep(time.Duration(i+1) * time.Second)
				Error.Printf("Sleep for %d seconds @ workerID %d @ %s \n", i+1, workerID, link)
				continue
			}
		} else {
			//Info.Println(resp.Body)
			break
		}
	}
	defer resp.Body.Close()

	links := collectLinks(link, resp.Body)
	Debug.Printf("DEBUG \t %s contains: %s \n", link, links)
	Debug.Printf("DEBUG \t Found %d link(s) on %s \n", len(links), link)
	/*
		BUGGY FOR WHATEVER FUCKING REASON
	*/
	if multithreaded == true {
		AddCountB(len(links))
		for i, _ := range links {
			go test(link, links[i])
		}
	} else {
		// Works
		for _, foundLink := range links {
			absoluteUrl := FixUrl(&foundLink, &link)
			if absoluteUrl != "" {
				if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl) {
					Debug.Printf("DEBUG \t Tests passed: %s \n", absoluteUrl)
					mutex.Lock()
					if visited[absoluteUrl] == false {
						visited[absoluteUrl] = true
						AddLinkCount()
						//Info.Printf("INFO\t%d @ workerID %d @ %s \n", GetLinkCount(), workerID, absoluteUrl)
						//Pure.Println(absoluteUrl)
						mutex.Unlock()
						Debug.Printf("added to channel: %s \n", absoluteUrl)
						AddCountB(1)
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

}

func test(link, foundLink string) {
	defer DoneCountB()
	absoluteUrl := FixUrl(&foundLink, &link)
	Debug.Println(absoluteUrl)
	if absoluteUrl != "" {
		if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl) {
			Debug.Printf("     *** Tests passed: %s \n", absoluteUrl)
			mutex.Lock()
			if visited[absoluteUrl] == false {
				visited[absoluteUrl] = true
				AddLinkCount()
				//Info.Printf(" %d @ %s \n", GetLinkCount(), absoluteUrl)
				mutex.Unlock()
				Debug.Printf("added to channel: %s \n", absoluteUrl)
				AddCountB(1)
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
	return
}
