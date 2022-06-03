package main

import (
	"fmt"
	"runtime"
	"sync"
)

// We need to use mutex for the shared variable don't get corrupted

func main() {
	runtime.GOMAXPROCS(3)
	var balance int
	var wg sync.WaitGroup

	var mu sync.Mutex

	deposit := func(amount int) {
		mu.Lock()
		balance += amount
		mu.Unlock()
	}

	whitdrawal := func(amount int) {
		mu.Lock()
		balance -= amount
		mu.Unlock()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			deposit(1)
		}()
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			whitdrawal(1)
		}()
	}
	wg.Wait()
	fmt.Println(balance)
}
