package main

import (
	"flag"
	"fmt"

	"bitbucket.org/advbet/rng"
)

func main() {
	var n int
	var count int

	flag.IntVar(&n, "range", 256, "range to draw random number from, value will be from 0 to n-1")
	flag.IntVar(&count, "count", 1, "limit on how many numbers wil be drawn")
	flag.Parse()

	for i := 0; i < count; i++ {
		fmt.Println(rng.Intn(n))
	}
}
