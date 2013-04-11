/* A simple parallel grabber
 * Note: It has scheduling flaws
 */
package main

import (
	"fmt"
	"net/http"
	"os"
)

func getUrl(url string) {
	fmt.Printf("Get %s\n", url)

	if resp, err := http.Get("http://example.com/"); err != nil {
		fmt.Printf("Error fetching %s: %s\n", url, err)

	} else {
		fmt.Printf("Done with %s: %d bytes\n", url, resp.ContentLength)
	}
}

func main() {
	urls := os.Args[1:]

	fmt.Printf("Getting %d urls: %v\n", len(urls), urls)

	for _, url := range urls {
		go getUrl(url)
	}

	fmt.Printf("Done\n")
}
