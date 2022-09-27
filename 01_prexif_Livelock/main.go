package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main(){

	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1*time.Millisecond) {
			cadence.Broadcast()
		}
	}()
	
	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}
	/*
	allows a person to attempt to move in a direction and returns whether or not they were successful. 
	Each direction is represented as a count of the number of people trying to move in that direction, dir
	*/
	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool { //1
		fmt.Fprintf(out, " %v", dirName)
		atomic.AddInt32(dir, 1) //2 we declare our intention to move in a direction by incrementing that direction by one
		takeStep() // 3 , each person must move at the same rate of speed, or cadence. takeStep simulates a constant cadence between all parties.
		if atomic.LoadInt32(dir) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(dir, -1) // 4 Here the person realizes they cannot go in this direction and gives up. We indicate this by decrementing that direction by one.
		return false
	}	
	
	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool { return tryDir("left", &left, out) }
	tryRight := func(out *bytes.Buffer) bool { return tryDir("right", &right, out) }


	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() { fmt.Println(out.String()) }()
		defer walking.Done()
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ { // 1 I placed an artificial limit on the number of attempts so that this program would end. In a program that has a livelock, there may be no such limit, which is why it’s a problem!
			if tryLeft(&out) || tryRight(&out) { //2 First, the person will attempt to step left, and if that fails, they will attempt to step right.
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup //3 	This variable provides a way for the program to wait until both people are either able to pass one another, or give up.
	peopleInHallway.Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Barbara")
	peopleInHallway.Wait()
}
