package utils

import (
	"fmt"

	"golang.org/x/net/html"
)

func Contains(strings []string, target string) bool {
	for _, value := range strings {
		fmt.Println(value)
		if value == target {
			return true
		}
	}
	return false
}

func FindUrls(node *html.Node) []string {
	urls := []string{}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				urls = append(urls, a.Val)
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		others := FindUrls(c)
		for i := range others {
			urls = append(urls, others[i])
		}
	}
	return urls
}
