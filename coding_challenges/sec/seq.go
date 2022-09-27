package main

import "fmt"

func main() {
	a := []int{5, 1, 22, 25, 6, -1, 8, 10}
	b := []int{1, 6, -1, 10}
	// IsValidSukbsequence(a, b)
	fmt.Println(IsValidSubsequence(a, b))
}

func IsValidSubsequence(array []int, sequence []int) bool {
	seqIdx := 0

	for _, num := range(array) {
		fmt.Println("num: ",num, " seqIdx: ", seqIdx)
		fmt.Println("num: ",num, " sequence[seqIdx]: ", sequence[seqIdx])
		if num == sequence[seqIdx] {
			fmt.Println("**************")
			seqIdx++
		}

		if seqIdx == len(sequence) {
			return true
		}
	}

	return false
}
