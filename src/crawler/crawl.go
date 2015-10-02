package main

import (
	"net/http"
)

func Crawl(link string, startHost string) {
	defer DoneCountA()
	Debug.Printf("Begin crawl: %s \n", link)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		Error.Printf("  Connection Error for %s : %s \n", link, err)
		return
	}
	//	UNIQUE USER AGENT, NO 'BOT'
	req.Header.Set("User-Agent", "einmal gemischte Tüte ohne Lakritz")

	resp, err := client.Do(req)
	if err != nil {
		Error.Printf("  Connection Error for %s : %s \n", link, err)
		return
	}
	defer resp.Body.Close()

	links := collectLinks(&link, resp.Body)

	Debug.Printf("  Found %d links on %s \n", len(links), link)

	/*
		BUGGY FOR WHATEVER FUCKING REASON
	*/
	if multithread == true {
		AddCountB(len(links))
		for i, _ := range links {
			go test(link, startHost, links[i])
		}
	} else {
		// Works
		for _, foundLink := range links {
			absoluteUrl := FixUrl(&foundLink, &link)
			Debug.Println(absoluteUrl)
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

func test(link, startHost, foundLink string) {
	defer DoneCountB()
	absoluteUrl := FixUrl(&foundLink, &link)
	Debug.Println(absoluteUrl)
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
