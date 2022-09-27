package main

import (
	"fmt"
	"runtime"
	"time"

	// "github.com/aws/aws-lambda-go/lambda"
)
func main(){
//  lambda.Start(ttt)
	ttt()
}
func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func ttt() {
	runtime.GOMAXPROCS(1)
	go say("world")
	say("hello")
}
