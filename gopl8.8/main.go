package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer (*wg).Done()
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func scan(c net.Conn, out chan<- string) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		out <- input.Text()
	}
	close(out)
}

func handleConn(c net.Conn) {
	var wg sync.WaitGroup

	scanCh := make(chan string)
	go scan(c, scanCh)

	loop := true
	for loop {
		select {
		case toEcho := <-scanCh:
			wg.Add(1)
			go echo(c, toEcho, 1*time.Second, &wg)
		case <-time.After(10 * time.Second):
			log.Println("connection timedout")
			loop = false
			break
		}
	}

	wg.Wait()
	c.(*net.TCPConn).CloseWrite()
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
