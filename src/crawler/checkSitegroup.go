package main

import (
	"code.google.com/p/go.net/html"
	//"fmt"
	"io"
	//"regexp"
)

func CheckSitegroup(httpBody io.Reader) []string {
	sitegroups := make([]string, 0)
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		//fmt.Println("TokenType:", tokenType)
		// check if HTML file has ended
		if tokenType == html.ErrorToken {
			return sitegroups
		}
		token := page.Token()
		//fmt.Println("Token:", token)
		if tokenType == html.StartTagToken && token.DataAtom.String() == "script" {
			for _, attr := range token.Attr {
				//fmt.Println("ATTR.KEY:", attr.Key)
				sitegroups = append(sitegroups, attr.Val)
			}

		}
	}
}
