package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

var dataStream chan interface{} // Here we declare a channel. We say it is “of type” interface{} since the type we’ve declared is the empty interface.

func main() {
	// main2()
	// main3()
	// buggerChannelExample()
	// selects()
	// nochanneIsReady()
	// noChannelReady2()
	ExampleAllowsExitWithoutBlocking()
	// dataStream2 := make(chan interface{}) // Here we instantiate the channel using the built-in make function.

	// fmt.Println(dataStream2)
	// // Valid statements:
	// stringStream := make(chan string)
	// go func() {
	// 	stringStream <- "Hello channels!"
	// }()
	// salutation, ok := <-stringStream //1 Here we receive both a string, salutation, and a boolean, ok.
	// fmt.Printf("(%v): %v", ok, salutation)
}

//Closed channel can be read
func main2() {

	intStream := make(chan int)
	go func() {
		defer close(intStream) // Here we ensure that the channel is closed before we exit the goroutine. This is a very common pattern.
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream { //2 Here we range over intStream.
		fmt.Printf("%v ", integer)
	}
}

// Closing a channel is also one of the ways you can signal multiple goroutines simultaneously.
func main3() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin //1 Here the goroutine waits until it is told it can continue.
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin) //2 Here we close the channel, thus unblocking all the goroutines simultaneously.
	wg.Wait()
}

// You can see that none of the goroutines begin to run until after we close the begin channel:
/* output: -------first line Unblocking goroutines...
4 has begun
2 has begun
3 has begun
0 has begun
1 has begun */

// IMPORTANT:We can also create buffered channels
// buffered channel is simply a buffered channel created with a capacity of 0 make(chan int, 0)
func bufferedChannles() {
	// var dataStream chan interface{}
	dataStream = make(chan interface{}, 4) //1 Here we create a buffered channel with a capacity of four. This means that we can place four things onto the channel regardless of whether it’s being read from.
}

func buggerChannelExample() {
	var stdoutBuff bytes.Buffer         //1 Here we create an in-memory buffer to help mitigate the nondeterministic nature of the output. It doesn’t give us any guarantees, but it’s a little faster than writing to stdout directly.
	defer stdoutBuff.WriteTo(os.Stdout) // 2 Here we ensure that the buffer is written out to stdout before the process exits.<

	intStream := make(chan int, 4) //3 Here we create a buffered channel with a capacity of one.
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
}

/*
Let’s begin with channel owners. The goroutine that owns a channel should:

Instantiate the channel.

Perform writes, or pass ownership to another goroutine.

Close the channel.

Ecapsulate the previous three things in this list and expose them via a reader channel.

By assigning these responsibilities to channel owners, a few things happen:

Because we’re the one initializing the channel, we remove the risk of deadlocking by writing to a nil channel.

Because we’re the one initializing the channel, we remove the risk of panicing by closing a nil channel.

Because we’re the one who decides when the channel gets closed, we remove the risk of panicing by writing to a closed channel.

Because we’re the one who decides when the channel gets closed, we remove the risk of panicing by closing a channel more than once.

We wield the type checker at compile time to prevent improper writes to our channel.


select statements can help safely bring channels together with concepts like cancellations, timeouts, waiting, and default values.

*/

func selects() {
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

// What about the second question: what happens if there are never any channels that become ready?
func nochanneIsReady() {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second): // The time.After function takes in a time.Duration argument and returns a channel that will send the current time after the duration you provide it.
		fmt.Println("Timed out.")
	}
}

// no channel is ready, and we need to do something in the meantime?
func noChannelReady2() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start)) // almost instantaneously
	}
}

func ExampleAllowsExitWithoutBlocking() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
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
	// Achieved 6 cycles of work before signalled to stop.
}


// Finally, there is a special case for empty select statements: select statements with no case clauses. These look like this: