/* A simple parallel grabber
 * It uses channels to ensure coordination between the master goroutine and
 * url grabbers.
 * It tracks durations and whether the responses included Vary headers.
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type VarySender interface {
	ICanHazVary() bool
}

type concreteVarySender struct {
	Url     string
	hazVary bool
}

func (ccs *concreteVarySender) ICanHazVary() bool {
	return ccs.hazVary
}

func getUrl(url string, duration chan time.Duration, varies chan VarySender) {
	fmt.Printf("Get %s\n", url)

	start := time.Now()
	defer func() { duration <- time.Since(start) }()

	if resp, err := http.Get("http://example.com/"); err != nil {
		fmt.Printf("Error fetching %s: %s\n", url, err)

	} else {
		if resp.Header.Get(http.CanonicalHeaderKey("Vary")) != "" {
			varies <- &concreteVarySender{
				Url:     url,
				hazVary: true}
		}

		fmt.Printf("Done with %s: %d bytes and status: %s\n", url, resp.ContentLength, resp.Status)
	}
}

func main() {
	var (
		urls     = os.Args[1:]
		duration = make(chan time.Duration)
		varies   = make(chan VarySender, 128)
	)

	fmt.Printf("Version %s\n", version)
	fmt.Printf("Getting %d urls: %v\n", len(urls), urls)

	for _, url := range urls {
		go getUrl(url, duration, varies)
	}

	for i := 0; i < len(urls); i++ {
		<-duration
	}

	close(varies)

	for vary := range varies {
		if vary.ICanHazVary() {
			if ccs, ok := vary.(*concreteVarySender); ok {
				fmt.Printf("%s hazes the varies\n", ccs.Url)

			} else {
				fmt.Printf("something hazes the varies\n")
			}
		}
	}

	fmt.Printf("Done\n")
}
