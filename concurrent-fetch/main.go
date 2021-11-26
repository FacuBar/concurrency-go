package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	f *os.File
)

func main() {
	var err error
	f, err = os.OpenFile("concurrent-fetch/test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic("couldn't open file")
	}
	defer f.Close()

	URLs := []string{"http://facebook.com", "https://github.com", "https://google.com"}
	fmt.Println("----Concurrent----")
	concurrent(URLs)
}

func concurrent(URLs []string) {
	start := time.Now()

	ch := make(chan string)
	for _, url := range URLs {
		go fetch(url, ch)
	}

	// Although go routines are executing during the stalled time
	// the flow of the go routine
	// is interrupted until the receiver channel is open

	fmt.Println("sleeping")
	time.Sleep(2 * time.Second)
	fmt.Println("awake")

	for range URLs {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	fmt.Println("hola")

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()

	var bytes bytes.Buffer
	nbytes, err := io.Copy(&bytes, resp.Body)
	f.Write(bytes.Bytes())

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	fmt.Println("defer adieu pre sender channel")

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)

	// if the sleep that is before the receiver channel is longer that
	// the slowest request is incertain wether 'adieu' will be printed
	// or not
	// as the program will most certainly end before the prints are executed -although that is incertain too-

	fmt.Println("adieu")
}
