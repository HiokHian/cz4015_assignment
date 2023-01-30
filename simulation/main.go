package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("here", rand.Float64())
	// r := rect{width: 3, height: 4}
	// c := circle{radius: 5}
	// u := UniformDistribution{left: 1, right: 5}
	// e := ExponentialDistribution{beta: 5}
	n := NormalDistribution{mean: 0, std:1}
	for i := 0; i < 10; i++ {
		// fmt.Println(u.sample())
		// fmt.Println(rand.Int63n(20))//0 to 20
		// fmt.Println(e.sample())
		fmt.Println(n.sample())

	}
	// fmt.Println(u.sample())
	// fmt.Println(u.sample())

	// measure(r)
	// measure(c)
}