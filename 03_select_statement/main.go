package main

import (
	"fmt"
	"time"

	"github.com/Millicom-MFS/kit-go/log"
)

func main() {
	// select1()
	// selectSimultaneouus()
	// Timeouts()

	DEfauktExample()
}

func select1() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c) //1 Here we close the channel after waiting five seconds.
	}()

	log.Info("Blocking on read...")
	select {
	case <-c: // 2 	Here we attempt a read on the channel. Note that as this code is written, we don’t require a select statement—we could simply write <-c—but we’ll expand on this example.
		log.Info("Unblocked %v later.\n", time.Since(start))
	}
}
func selectSimultaneouus() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func Timeouts() {
	start := time.Now()
	var c <-chan int
	select {
	case <-c: // his case statement will never become unblocked because we’re reading from a nil channel.
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}

}

func DEfauktExample() {
	done := make(chan interface{})
	go func() {
		time.Sleep(7 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

// Finally, there is a special case for empty select statements: select statements with no case clauses. These look like this:

// This statement will simply block forever.

