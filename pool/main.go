package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

//Create aand make available pool of things for use

func log(w io.Writer, debug string) {
	//The issue with the log function is each call to the log will result in the creation of a new bytes.buffer
	// If the log function is called from thousand of places in your aplicatin , then that could lead to a lot of stale memory a lot of garbage collector
	// To solve this issue we can use pool, to hold the pool of bytes.Buffer that can be reused instead of creating new instances on each call
	var b bytes.Buffer

	b.WriteString(time.Now().Format("15:05:05"))
	b.WriteString(":")
	b.WriteString(debug)
	b.WriteString("\n")

	w.Write(b.Bytes())
}

// Create a instance of the sync.Pool
var bufPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("Allocating new bytes.Buffer")
		return new(bytes.Buffer)
	}, //Holds the reference of the function that will be called
}

func log2(w io.Writer, debug string) {
	b := bufPool.Get().(*bytes.Buffer) //Get buffpooll casting with bytes.Buffer

	b.Reset() // Reset the content of the buffer

	b.WriteString(time.Now().Format("15:05:05"))
	b.WriteString(":")
	b.WriteString(debug)
	b.WriteString("\n")

	w.Write(b.Bytes())

	bufPool.Put(b)
}

func main() {
	log2(os.Stdout, "debug-srting1")
	log2(os.Stdout, "debug-strng2")
}

// With log2 allocate bytes.Buffer
// The second time
// The same instance of bytes.Buffer is reused to process the debug strig
// Reater than create instance each time
