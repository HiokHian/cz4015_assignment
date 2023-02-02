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

	// instantiate base stations
	base_stations := make([]BaseStation, 20)
	for i := range base_stations{
		base_stations[i] = BaseStation{num_free_channels: i}
	}

	//FEL needs to be a priority queue

	for i := 0; i < 2; i++ {
		// fmt.Println(u.sample())
		// fmt.Println(rand.Int63n(20))//0 to 20
		// fmt.Println(e.sample())
		fmt.Println(base_stations[i])
		bs := BaseStation{num_free_channels:10}
		bs.assign_channel()
		ci := CallInitiation{time: 2.14, speed: 0, station: &BaseStation{num_free_channels: 10}, position: 0, duration: 0, direction: 1}
		fmt.Println(bs.num_free_channels)
		bs.free_up_channel()
		fmt.Println(bs.num_free_channels)
		fmt.Println(bs.has_available_channel())
		fmt.Println(n.sample())
		fmt.Println(ci.time)

	}
	// fmt.Println(u.sample())
	// fmt.Println(u.sample())

	// measure(r)
	// measure(c)
}