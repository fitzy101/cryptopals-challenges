package main

import (
	"fmt"
)

var a = "this is a test"
var b = "wokka wokka!!!"

func main() {
	totalCount := 0
	for i := range a {
		totalCount += countBitDiff(a[i], b[i])
	}
	fmt.Println(totalCount)
}
