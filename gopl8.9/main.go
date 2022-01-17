// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 251.

// The du4 command computes the disk usage of the files in a directory.
package main

// The du4 variant includes cancellation:
// it terminates quickly when the user hits return.

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//!+1
var done = make(chan struct{})

type rootDir struct {
	name string
	size int64
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

//!-1

func main() {
	// Determine the initial directories.
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	rootDirs := []rootDir{}
	for _, root := range roots {
		rootDirs = append(rootDirs, rootDir{name: root})
	}

	//!+2
	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()
	//!-2

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan rootDir)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes map[string]int64
loop:
	//!+3
	for {
		select {
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish.
			for range fileSizes {
				// Do nothing.
			}
			return
		case size, ok := <-fileSizes:
			// ...
			//!-3
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles[size.name]++
			nbytes[size.name] += size.size
		case <-tick:
			printDiskUsage(rootDirs, nfiles, nbytes)
		}
	}
	printDiskUsage(rootDirs, nfiles, nbytes) // final totals
}

func printDiskUsage(rootDirs []rootDir, nfiles, nbytes map[string]int64) {
	for _, rootdir := range rootDirs {
		name := rootdir.name
		fmt.Printf("[%s]%d files  %.1f GB\n", name, nfiles[name], float64(nbytes[name])/1e9)
	}
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+4
func walkDir(rootName, dir string, n *sync.WaitGroup, fileSizes chan<- rootDir) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		// ...
		//!-4
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(rootName, subdir, n, fileSizes)
		} else {
			fileSizes <- rootDir{name: rootName, size: entry.Size()}
		}
		//!+4
	}
}

//!-4

var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// dirents returns the entries of directory dir.
//!+5
func dirents(dir string) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil // cancelled
	}
	defer func() { <-sema }() // release token

	// ...read directory...
	//!-5

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0 => no limit; read all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// Don't return: Readdir may return partial results.
	}
	return entries
}
