package main

import (
	"net/url"
	"strings"
)

var badChars = []string{"#", "@", ";"}
var badFileEndings = []string{".gif", ".jpg", ".jpeg", ".svg", ".js", ".png", ".pdf", ".swf"}

/*
	Parses String to URL and parses them to absolute URLs if necessary
*/
func FixUrl(href, base *string) string {
	uri, err := url.Parse(*href)
	if err != nil {
		MutexErrorAdd()
		Error.Printf("  FixUrl() - Parsing Url failed: \n", err)
		return ""
	}
	baseUrl, err := url.Parse(*base)
	if err != nil {
		MutexErrorAdd()
		Error.Printf("  BaseURL ERROR: %s \n", err)
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

/*
	filters URLs for email adresses, Javascript, images, PDFs etc
*/
func CheckUrl(uri *string) bool {
	for _, str := range badChars {
		if strings.Contains(*uri, str) {
			Debug.Printf("  Bad Char %s for %s", str, *uri)
			return false
		}
	}
	for _, str := range badFileEndings {
		if strings.HasSuffix(*uri, str) {
			Debug.Printf("  Bad File Ending %s for %s", str, *uri)
			return false
		}
	}

	return true
}

/*
	checks if the found URL and the start URL have the same domain = link to the same page = don't leave the start page (startPage)
*/
func CheckHost(uri, startHost *string) bool {
	uriUrl, err := url.Parse(*uri)
	if err != nil {
		MutexErrorAdd()
		Error.Printf("  CheckHost() - Url parsing failed: %s", err)
		return false
	}
	//fmt.Printf("Site Host: %s ## START HOST: %s \n", uriUrl.Host, *startHost)
	if uriUrl.Host == *startHost {
		return true
	} else {
		Debug.Printf("  Bad Host: %s for %s \n", uriUrl.Host, uriUrl)
		return false
	}

}
