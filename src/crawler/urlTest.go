package main

import (
	"fmt"
	"net/url"
	"strings"
)

var badChars = []string{"@", ";"}
var badFileEndings = []string{".gif", ".jpg", ".jpeg", ".js", ".png", ".pdf", ".swf"}

/*
	Parses String to URL and parses them to absolute URLs if necessary
*/
func FixUrl(href, base *string) string {
	uri, err := url.Parse(*href)
	if err != nil {
		//fmt.Printf("PARSE ERROR: \n", err)
		return ""
	}
	baseUrl, err := url.Parse(*base)
	if err != nil {
		//fmt.Printf("BaseURL ERROR: \n", err)
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	//fmt.Println("FIXED URL: ", uri)
	return uri.String()
}

/*
	filters URLs for email adresses, Javascript, images, PDFs etc
*/
func CheckUrl(uri *string) bool {
	for _, str := range badChars {
		if strings.Contains(*uri, str) {
			//fmt.Printf("BADCHAR: %-80s  ERR: %-80s \n", *uri, str)
			return false
		}
	}
	for _, str := range badFileEndings {
		if strings.HasSuffix(*uri, str) {
			return false
		}
	}

	for u := range badFileEndings {
		str := string(u)
		if strings.HasSuffix(*uri, str) {
			return false
		}
	}

	return true
}

/*
	checks if the found URL and the start URL have the same domain = link to the same page = don't leave the start page (startPage)
*/
func CheckTLD(uri, startHost *string) bool {
	uriUrl, err := url.Parse(*uri)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Site Host: %s ## START HOST: %s \n", uriUrl.Host, *startHost)
	if uriUrl.Host == *startHost {
		return true
	} else {
		//fmt.Printf("WRONG HOST: %-30s INSTEAD OF %-30s \n", uriUrl.Host, *startHost)
		return false
	}

}
