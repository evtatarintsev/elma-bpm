package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPrintResults(t *testing.T)  {
	// arrange
	results := make(chan Result, 3)
	results <- Result{"some url", 10, nil}
	results <- Result{"next url", 5, nil}
	results <- Result{"wrong url", 0, errors.New("some error")}
	close(results)

	const expectedOut = "Count for some url: 10\n" +
						"Count for next url: 5\n" +
						"Error fetch wrong url: some error\n" +
						"Total 15\n"

	tmpfile, _ := ioutil.TempFile("", "example")
	oldStdout := os.Stdout
	os.Stdout = tmpfile

	// act
	printResults(results)

	// assert
	os.Stdout.Seek(0, 0)
	output, _ := ioutil.ReadAll(os.Stdout)
	os.Stdout = oldStdout
	if string(output) != expectedOut {
		t.Error("Wrong output")
	}
}


func TestReadUrls(t *testing.T)  {
	// arrange
	urls := make(chan string)

	// act
	readUrls(urls)

	// assert
	_, notClosed := <- urls
	if notClosed {
		t.Error("readUrls must close channel")
	}
}

func TestFetchWordCountFromUrl(t *testing.T)  {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`some text Go Go Go other text`))
	}))

	// act
	result := fetchWordCountFromUrl(server.URL, "Go")

	// assert
	if result.count != 3 {
		t.Error("Key word counting error")
	}
}
