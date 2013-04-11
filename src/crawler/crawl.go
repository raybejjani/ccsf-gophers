/* A simple parallel grabber
 * It uses channels to ensure coordination between the master goroutine and
 * url grabbers
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func getUrl(url string, duration chan time.Duration) {
	fmt.Printf("Get %s\n", url)

	start := time.Now()
	defer func() { duration <- time.Since(start) }()

	if resp, err := http.Get("http://example.com/"); err != nil {
		fmt.Printf("Error fetching %s: %s\n", url, err)

	} else {
		fmt.Printf("Done with %s: %d bytes and status: %s\n", url, resp.ContentLength, resp.Status)
	}
}

func main() {
	urls := os.Args[1:]
	var duration = make(chan time.Duration)

	fmt.Printf("Getting %d urls: %v\n", len(urls), urls)

	for _, url := range urls {
		go getUrl(url, duration)
	}

	for i := 0; i < len(urls); i++ {
		<-duration
	}

	fmt.Printf("Done\n")
}
