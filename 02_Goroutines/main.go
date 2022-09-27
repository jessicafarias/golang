package main

import (
	"fmt"
	"sync"
)

func main() {
	PrintTwo()
}

func PrintOne() {
	var wg sync.WaitGroup
	salutation := "hello"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome" //1 Here we see the goroutine modifying the value of the variable salutation.
	}()
	wg.Wait()
	fmt.Println(salutation)
}

func PrintTwo() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation) //1
		}(salutation)
	}
	wg.Wait()
}
