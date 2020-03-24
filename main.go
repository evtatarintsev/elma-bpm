package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Result struct {
	url string
	count int
	err error
}


func readUrls(urls chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urls <- scanner.Text()
	}
	close(urls)
}


func printResults(results chan Result)  {
	totalCount := 0
	for result := range results {
		totalCount += result.count
		if result.err != nil {
			fmt.Printf("Error fetch %s: %s\n", result.url, result.err)
			continue
		}
		fmt.Printf("Count for %s: %d\n", result.url, result.count)
	}
	fmt.Printf("Total %d\n", totalCount)
}


func fetchWordCountFromUrl(url string, word string) Result {
	result := Result{url: url, count: 0, err: nil}

	response, err := http.Get(url)
	if err != nil {
		result.err = err
		return result
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		result.err = err
		return result
	}
	result.count = strings.Count(string(body), word)
	return result
}


func fetchResults(urls chan string, word string, results chan Result, maxGoroutines int)  {
	var wg sync.WaitGroup
	goroutines := make(chan int, maxGoroutines)

	for url := range urls {
		goroutines <- 1
		wg.Add(1)

		go func() {
			results <- fetchWordCountFromUrl(url, word)
			<-goroutines
			wg.Done()
		}()
	}
	wg.Wait()
	close(results)
}


func main() {
	const maxGoroutines = 5
	const word = "Go"
	urls := make(chan string)
	results := make(chan Result)

	go readUrls(urls)
	go fetchResults(urls, word, results, maxGoroutines)

	printResults(results)
}
