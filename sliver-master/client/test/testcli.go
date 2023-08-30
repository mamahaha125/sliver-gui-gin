package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "cat"
	b := strings.SplitAfterN(a, " ", 2)
	fmt.Println(len(b))
	for _, v := range b {
		fmt.Println(v)
	}

}
