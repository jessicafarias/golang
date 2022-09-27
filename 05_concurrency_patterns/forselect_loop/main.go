package main

// Something you’ll see over and over again in Go programs is the for-select loop. It’s nothing more than something like this:

func main(){
	forselectloop()
}

func forselectloop(){
	for { // Either loop infinitely or range over something
		select {
		// Do some work with channels
		}
	}
}

// There are a couple of different scenarios where you’ll see this pattern pop up.
// NOTE : Sending iteration variables out on a channel
func IterationVariables(){
	var done chan string
	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}
}

// NOTE Looping infinitely waiting to be stopped
func Loopingindefinitely(){
	for {
		select {
		case <-done:
			return
		default:
		}
	
		// Do non-preemptable work
	}
}