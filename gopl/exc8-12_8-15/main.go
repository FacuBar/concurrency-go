// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

//!+broadcaster
type client struct {
	ch   chan<- string // an outgoing message channel
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages

	iddleLimit = 15 * time.Second
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				select {
				case cli.ch <- msg:
					// mssg sent successfully
				default:
					// mssg could not be sent, not blocking
				}

			}

		case cli := <-entering:
			clients[cli] = true
			actClients := strings.Builder{}
			for client := range clients {
				actClients.WriteString(fmt.Sprintf("%s, ", client.name))
			}
			cli.ch <- fmt.Sprintf("active clients: %s", actClients.String())

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	msgSent := make(chan struct{})
	go disconnectIddle(conn, ch, msgSent)

	input := bufio.NewScanner(conn)
	fmt.Fprint(conn, "insert your name: ")
	var name string
	if input.Scan() {
		name = input.Text()
	} else {
		fmt.Fprintln(conn, "name not insert or invalid")
		conn.Close()
		return
	}

	ch <- "You are " + name
	messages <- name + " has arrived"

	cli := client{ch: ch, name: name}
	entering <- cli

	for input.Scan() {
		messages <- name + ": " + input.Text()
		msgSent <- struct{}{}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messages <- name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

func disconnectIddle(conn net.Conn, clich chan<- string, msgSentCh <-chan struct{}) {
	for {
		select {
		case <-msgSentCh:
			continue
		case <-time.After(iddleLimit):
			conn.Close()
		}
	}
}

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
