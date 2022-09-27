package main

import (
	"bytes"
	"fmt"
	"sync"
)

//Confinement

func main() {
	// example1()
	// example2()
	ImposibleChannel()

}

// Ad hoc confinement is when you achieve confinement through a convention—
// Here’s an example of ad hoc confinement that demonstrates why:
func example1() {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}

// Lexical confinement involves using lexical scope to expose only the correct data and concurrency primitives for multiple concurrent processes to use
func example2(){
	chanOwner := func() <-chan int {
		results := make(chan int, 5) //1 Here we instantiate the channel 
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}
	
	consumer := func(results <-chan int) { //3 Here we receive a read-only copy of an int channel. 
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	
	results := chanOwner()        //2 Here we receive the read aspect 
	consumer(results)
}

// IMPORTANT Set up this way, it is impossible to utilize the channels in this small example. 
/*

This is a good lead-in to confinement, but probably not a very interesting example 
since channels are concurrent-safe. Let’s take a look at an example 
of confinement that uses a data structure which is not concurrent-safe, 
an instance of bytes.Buffer

Concurrent code that utilizes lexical confinement also has the benefit of usually 
being simpler to understand than concurrent code without lexically confined variables.

So what’s the point? Why pursue confinement if we have synchronization available to us? 
The answer is improved performance and reduced cognitive load on developers

Synchronization comes with a cost, and if you can avoid it you won’t have any critical sections

*/
func ImposibleChannel(){
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
	
		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}
	
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])     // 1 Here we pass in a slice containing the first three bytes in the data structure
	go printData(&wg, data[3:])     // 2 Here we pass in a slice containing the last three bytes in the data structure
	
	wg.Wait()
}