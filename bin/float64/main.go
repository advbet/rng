package main

import (
	"flag"
	"fmt"

	"bitbucket.org/advbet/rng"
)

func main() {
	var count int

	flag.IntVar(&count, "count", 1, "limit on how many numbers wil be drawn")
	flag.Parse()

	for i := 0; i < count; i++ {
		fmt.Println(rng.Float64())
	}
}
