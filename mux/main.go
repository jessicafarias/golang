package main

import (
	"fmt"
	"runtime"
	"sync"
	"github.com/aws/aws-lambda-go/lambda"
)

func main(){	
	lambda.Start(test)

}

// We need to use mutex for the shared variable don't get corrupted
func getGOMAXPROCS() int {
    return runtime.GOMAXPROCS(0)
}

func test() {
	fmt.Printf("GOMAXPROCS is %d\n", getGOMAXPROCS())
	runtime.GOMAXPROCS(2)
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
