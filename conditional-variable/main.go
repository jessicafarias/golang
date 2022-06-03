package main

import (
	"fmt"
	"sync"
	"time"
)

//Use this condituional variable to coordinate a single producer and multiple consumer
var SharedRsc = make(map[string]interface{})

// There are two goroutines waiting for differnt conditions
func main() {
	var wg sync.WaitGroup

	// We are using mutex bc we are accesing to sharing resource
	mu := sync.Mutex{}
	condition := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()
		//Susped goroutine until SharedRsc is populated
		condition.L.Lock()

		// Check if there are less than 1
		for len(SharedRsc) < 1 {
			condition.Wait()
			//time.Sleep(1 * time.Microsecond)
		}
		fmt.Println(SharedRsc["rsc1"])

		condition.L.Unlock()

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		//Suspend goroutine until sharedRd is populated
		condition.L.Lock()
		// Check if there are less than 2
		for len(SharedRsc) < 2 {
			condition.Wait()
			//time.Sleep(1 * time.Microsecond)
		}
		fmt.Println(SharedRsc["rsc2"])
		condition.L.Unlock()

	}()

	condition.L.Lock()

	//writes changes to SaredRsc
	SharedRsc["rsc1"] = "rsc1"
	time.Sleep(3 * time.Second)
	SharedRsc["rsc2"] = "rsc2"

	// Brocast method indicates all gouroutines the condition is met
	condition.Broadcast()
	condition.L.Unlock()
	wg.Wait()

}
