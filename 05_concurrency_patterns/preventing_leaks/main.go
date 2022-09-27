package main

import (
	"fmt"
	"time"
)

// Preventing Goroutine Leaks
/*
As we covered in the section “Goroutines”, we know goroutines are cheap and easy to create;
it’s one of the things that makes Go such a productive language.
The runtime handles multiplexing the goroutines onto any number of operating
system threads so that we don’t often have to worry about that level of abstraction.

But they do cost resources, and goroutines are not garbage collected by the runtime,
so regardless of how small their memory footprint is,
we don’t want to leave them lying about our process. So how do we go about ensuring they’re cleaned up?


The goroutine has a few paths to termination:
    -When it has completed its work.
    -When it cannot continue its work due to an unrecoverable error.
    -When it’s told to stop working.

*/
func main(){
	// ExampleOfGoroutineLeak()
	// FixExamble()
	// newIssue()
	FixIssue2()
}

// the lifetime of the process is very short, but in a real program, goroutines could easily be 
//started at the beginning of a long-lived program. In the worst case, 
//the main goroutine could continue to spin up goroutines throughout its life, causing creep in memory utilization.
func ExampleOfGoroutineLeak(){
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				// Do something interesting
				fmt.Println(s)
			}
		}()
		return completed
	}
	
	doWork(nil)
	// Perhaps more work is done here
	fmt.Println("Done.")
}


// The way to successfully mitigate this is to establish 
// a signal between the parent goroutine and its children that allows 
// the parent to signal cancellation to its children. By convention, 
// this signal is usually a read-only channel named done. The parent goroutine passes this channel 
// to the child goroutine and then closes the channel when it wants to cancel the child goroutine. Here’s an example:
func FixExamble(){
	doWork := func(done <-chan interface{},
		strings <-chan string,
	  ) <-chan interface{} { // 1 Here we pass the done channel to the doWork function. As a convention, this channel is the first parameter.
		  terminated := make(chan interface{})
		  go func() {
			  defer fmt.Println("doWork exited.")
			  defer close(terminated)
			  for {
				  select {
				  case s := <-strings:
					  // Do something interesting
					  fmt.Println(s)
				  case <-done: // 2  On this line we see the ubiquitous for-select pattern in use. One of our case statements is checking whether our done channel has been signaled. If it has, we return from the goroutine.
					  return
				  }
			  }
		  }()
		  return terminated
	  }
	  
	  done := make(chan interface{})
	  terminated := doWork(done, nil)
	  
	  go func() { //3 ere we create another goroutine that will cancel the goroutine spawned in doWork if more than one second passes.
		  // Cancel the operation after 1 second.
		  time.Sleep(1 * time.Second)
		  fmt.Println("Canceling doWork goroutine...")
		  close(done)
	  }()
	  
	  <-terminated // 4 This is where we join the goroutine spawned from doWork with the main goroutine.
	  fmt.Println("Done.")
}

// You can see that despite passing in nil for our strings channel, 
//our goroutine still exits successfully. Unlike the example before it, 
// in this example we do join the two goroutines, and yet do not receive a deadlock. 
// This is because before we join the two goroutines, we create a third goroutine to cancel the goroutine 
// within doWork after a second. We have successfully eliminated our goroutine leak!




/*

The previous example handles the case for goroutines receiving on a channel nicely, 
but what if we’re dealing with the reverse situation: a goroutine blocked on attempting 
to write a value to a channel? Here’s a quick example to demonstrate the issue:*/
func newIssue(){
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // 1 Here we print out a message when the goroutine successfully terminates.
			defer close(randStream)
			for {
				randStream <- 1 // rand.Int()
			}
		}()
	
		return randStream
	}
	
	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

/*

You can see from the output that the deferred fmt.Println statement never gets run. 
// After the third iteration of our loop, our goroutine blocks trying to send the next random 
// integer to a channel that is no longer being read from. We have no way of telling the producer it can stop. 
// The solution, just like for the receiving case, 
// is to provide the producer goroutine with a channel informing it to exit:*/

func FixIssue2(){
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- 5:
				case <-done:
					return
				}
			}
		}()
	
		return randStream
	}
	
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	
	// Simulate ongoing work
	time.Sleep(1 * time.Second)
}