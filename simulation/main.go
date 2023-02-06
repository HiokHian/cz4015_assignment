package main

import (
	"container/heap"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("here", rand.Float64())
	// r := rect{width: 3, height: 4}
	// c := circle{radius: 5}
	// u := UniformDistribution{left: 1, right: 5}
	// e := ExponentialDistribution{beta: 5}
	// n := NormalDistribution{mean: 0, std:1}

	// instantiate base stations
	base_stations := make([]BaseStation, 20)
	for i := range base_stations{
		base_stations[i] = BaseStation{num_free_channels: i}
	}

	//FEL needs to be a priority queue
	FEL:= make(PriorityQueue, 10)

	for i := 0; i < 10; i = i + 2 {

		// fmt.Println(u.sample())
		// fmt.Println(rand.Int63n(20))//0 to 20
		// fmt.Println(e.sample())
		// fmt.Println(base_stations[i])
		// bs := BaseStation{num_free_channels:10}
		// bs.assign_channel()
		
		// TODO: rewrite this into idiomatic golang using chatgpt
		i_float := (float64)(i)
		ci := (Event) (&CallInitiation{time: 2.14-i_float, speed: 0, station: &BaseStation{num_free_channels: 10}, position: 0, duration: 0, direction: 1})
		// ci_event := ((Event) (&ci)) // cast to Event 
		// FEL.Push(&ci)
		FEL[i] = &ci
		
		
		ct := (Event) (&CallTermination{speed: 0.5, time: 5.0 + i_float})
		// FEL.Push(&ct)
		FEL[i+1] = &ct
		// for j:= 0; j<i; j ++{

		// 	fmt.Println((*FEL[10]).get_time())
		// 	fmt.Println((*FEL[9]).get_time())
		// }

		// fmt.Println(bs.num_free_channels)
		// bs.free_up_channel()
		// fmt.Println(bs.num_free_channels)
		// fmt.Println(bs.has_available_channel())
		// fmt.Println(n.sample())
		// fmt.Println(ci.time)

	}
	fmt.Println(FEL)
	heap.Init(&FEL)
	ct := (Event) (&CallTermination{speed: 0.5, time: -500.0})
	heap.Push(&FEL, &ct) // NOTE: remember to use head.Push NOT FEL.Push!
	// FEL.Push(&ct)
	// fmt.Println("size of pq: ", len(FEL))
	for j := 0; j < 11; j ++ {
		// fmt.Println((*FEL[len(FEL)-1]).get_time(), " ", (*FEL[len(FEL)-1]).get_index())
		popped := heap.Pop(&FEL).(*Event) // Assert that the type is an Event pointer
		fmt.Println((*popped).get_time(), " ", )
	}

}