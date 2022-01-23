package ping_pong

import (
	"fmt"
	"time"
)

func pingPong(ping, pong string) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go pp(ch1, ch2, ping)
	go pp(ch2, ch1, pong)
	ch1 <- ""

	<-time.Tick(2 * time.Second)
}

func pp(in, out chan string, respond string) {
	for {
		<-in
		fmt.Println(respond)
		out <- respond
	}
}

// TODO: run benchmarks

// Exercise 9.5: Write a program with two goroutines that send messages back and forth
// over two unbuffered channels in ping-pong fashion. How many communications per
// second can the program sustain?
