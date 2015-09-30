package main

import (
	"net/http"
	"sync"
	"time"
)

func Crawl(link string, startHost string, wg *sync.WaitGroup, input, output chan string) {
	defer wg.Done()

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
		Error.Printf("  Connection Error for %s : %s", link, err)
		return
	}

	defer resp.Body.Close()

	links := collectLinks(resp.Body)

	for _, foundLink := range links {
		absoluteUrl := FixUrl(&foundLink, &link)
		if CheckUrl(&absoluteUrl) && CheckHost(&absoluteUrl, &startHost) {
			output <- absoluteUrl
		}
	}
	return
	//	} else {
	//		return
	//	}
}
