package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	port := flag.String("port", "8000", "assigns port to server")
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c,
			time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
