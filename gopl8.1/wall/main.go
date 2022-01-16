package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	errch := make(chan error, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		conn, err := net.Dial("tcp", strings.Split(arg, "=")[1])
		if err != nil {
			log.Printf("couldnt connect to %s\n", arg)
		}
		defer conn.Close()
		go mustCopy(os.Stdout, conn, strings.Split(arg, "=")[0], errch)
	}

	<-errch
}

func mustCopy(dst io.Writer, src io.Reader, city string, errch chan error) {
	scan := bufio.NewScanner(src)
	for scan.Scan() {
		if err := scan.Err(); err != nil {
			errch <- err
		}
		fmt.Fprintf(dst, "%15s: %s\n", city, scan.Text())
	}
}
