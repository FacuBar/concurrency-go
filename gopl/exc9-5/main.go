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

// Exercise 9.4: Construct a pipeline that connects an arbitrary number of goroutines
// with channels. What is the maximum number of pipeline stages you can create
// without running out of memory? How long does a value take to transit the entire
// pipeline?
