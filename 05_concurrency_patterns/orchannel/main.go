package main

import (
	"fmt"
	"time"
)

// At times you may find yourself wanting to combine one or more done channels into a single done channel that closes if any of its component channels close.

func main() {
	orchann()
}

func orchann() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	// 1 Here we have our function, or, which takes in a variadic slice of channels and returns a single channel.
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			/// 2 Since this is a recursive function, we must set up termination criteria.
			//  The first is that if the variadic slice is empty, we simply return a nil channel.
			// This is consistant with the idea of passing in no channels; we wouldn’t expect a composite channel to do anything.
			return nil
		case 1:
			// 3 Our second termination criteria states that if our variadic slice only contains one element, we just return that element.
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			// 4 Here is the main body of the function, and where the recursion happens.
			// We create a goroutine so that we can wait for messages on our channels without blocking.
			defer close(orDone)

			switch len(channels) {
			case 2:
				// 5 Because of how we’re recursing, every recursive call to or will at least have two channels.
				//  As an optimization to keep the number of goroutines constrained, we place a special case here for calls to or with only two channels.
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				// 6 Here we recursively create an or-channel from all the channels in our slice after the third index,
				// and then select from this. This recurrence relation will destructure the rest of the slice into or-channels
				// to form a tree from which the first signal will return. We also pass in the orDone channel so that when goroutines up the tree exit,
				//  goroutines down the tree also exit.
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
	sig := func(after time.Duration) <-chan interface{} { //1 This function simply creates a channel that will close when the time specified in the after elapses.
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	
	start := time.Now() //  2 Here we keep track of roughly when the channel from the or function begins to block.
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start)) // And here we print the time it took for the read to occur.
}

/*

Notice that despite placing several channels in our call to or that take various times to close, 
our channel that closes after one second causes the entire channel created by the call to or to close. 
This is because—despite its place in the tree the or function builds—it will always close first and thus the channels 
that depend on its closure will close as well.
*/