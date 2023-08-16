package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

func main() {
	url := flag.String("url", "", "The base url of the site")
	flag.Parse()

	if *url != "" {
		urls, error := webscrape(*url)
		if error != nil {
			fmt.Println(error.Error())
		} else {
			for _, url := range urls {
				fmt.Println(url)
			}
		}
	}
	return
}

func webscrape(url string) ([]string, error) {
	var urls []string
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Unable to fetch url")
	}
	defer response.Body.Close()

	tokeniser := html.NewTokenizer(response.Body)
TokenisationLoop:
	for {
		tokenType := tokeniser.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokeniser.Err()
			if err == io.EOF {
				break TokenisationLoop
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokeniser.Token()

			if token.Data == "a" || token.Data == "link" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						val := attr.Val
						if strings.HasPrefix(val, url) {
							if !slices.Contains(urls, val) {
								urls = append(urls, val)
							}
						}
					}
				}
			}
		}
	}
	return urls, nil
}
