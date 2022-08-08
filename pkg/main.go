package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
	fmt.Println("")
	fmt.Println("TO START")
	fmt.Print("PRESS ANY KEY")
	input := bufio.NewScanner(os.Stdin)
    input.Scan()
    fmt.Println(input.Text())
}