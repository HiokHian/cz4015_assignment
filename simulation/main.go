package main

import (
	"container/heap"
	"fmt"
	"math"
)

// TODO: verify if the declaring simulation_clock as
// a global variable here will allow it to be updated
// in the main function

// TODO: verify that base stations are edited in place with the subtraction of the channels
func main() {
	base_station_distribution := UniformDistribution{left: 0, right: 19} 
	direction_distribution := UniformDistribution{left: 0, right: 1} 
	call_interarrival_time_distribution := ExponentialDistribution{beta: 1.3698169264765245} //in seconds
	call_duration_distribution := ExponentialDistribution{beta:  99.83194913549607}
	speed_distribution := NormalDistribution{mean:  120.07209801685805, std: 9.01905789789691} //in km/h
	position_distribution := UniformDistribution{left: 0, right: 1} 
	// instantiate base stations
	base_stations := make([]BaseStation, 20)
	for i := range base_stations{
		base_stations[i] = BaseStation{
			id: i,
			num_free_channels: 10,
			num_free_reserved_channels: 0,
			using_reserved_fca_scheme: false,
		}
	}

	//FEL needs to be a priority queue
	var total_num_initiation_calls = 100000
	var cur_time = 0.0
	FEL:= make(PriorityQueue, total_num_initiation_calls)

	for i := 0; i < total_num_initiation_calls; i = i + 1 {
		// sample from random distributions for each of the parameters
		arrival_time := call_interarrival_time_distribution.sample() + cur_time
		cur_time = arrival_time

		call_duration := call_duration_distribution.sample()
		speed := speed_distribution.sample()
		
		station_number := (int) (math.Floor(base_station_distribution.sample()))
		direction_distribution := (int)(math.Round(direction_distribution.sample()))
		position := position_distribution.sample()
		// TODO: rewrite this into idiomatic golang using chatgpt
		ci := (Event) (&CallInitiation{
			time: arrival_time, 
			speed: speed, 
			current_station: &base_stations[station_number], 
			position: position, 
			duration: call_duration, 
			direction: direction_distribution,
		})
		FEL[i] = &ci // here we push a POINTER to an event => it is a pointer to a pointer to a call event
	}

	heap.Init(&FEL)  // Initialise the heap
	var num_blocked_calls float64 = 0.0
	var num_dropped_calls float64 = 0.0
	var total_num_calls float64 = 0.0
	// var warm_up_num_calls float64 = 200.0
	var done_with_warmup bool = false
	var current_pct_blocked_calls float64
	var current_pct_dropped_calls float64 
	var previous_pct_blocked_calls float64
	var previous_pct_dropped_calls float64 
	var calls_since_previous_computation int = 0
	var computation_interval int = 0
	var blocked_calls  float64
	var dropped_calls float64

	// FOR DEBUGGING 
	var verbose bool = false
	for {
		fmt.Println("Number of events left", len(FEL))
		if (len(FEL) == 0){ //  Stop iterating when FEL is empty
			break
		}

		total_num_calls += 1
		calls_since_previous_computation += 1
		
		if event, ok := heap.Pop(&FEL).(*Event); ok{
			fmt.Println("Event" , event)
			if call_initiation, ok := (*event).(*CallInitiation); ok { //dereference to an event then cast to cal initiation pointer
				fmt.Println("CallInitiation", call_initiation)
				blocked_calls, dropped_calls = (call_initiation).process_event(&FEL, base_stations) 
				fmt.Println("CallInitiation: ", blocked_calls, dropped_calls)
			} else if call_handover, ok := (*event).(*CallHandover); ok{
				fmt.Println("Callhandover", call_handover)
				blocked_calls, dropped_calls = (call_handover).process_event(&FEL, base_stations)
				fmt.Println("Callhandover: ", blocked_calls, dropped_calls)
			} else if call_termination, ok := (*event).(*CallTermination); ok{
				fmt.Println("Calltermination", call_termination)
				blocked_calls, dropped_calls = (call_termination).process_event(&FEL, base_stations) 
				fmt.Println("Calltermination: ", blocked_calls, dropped_calls)
			} else {
				fmt.Println("Unable to cast to event type")
			}
		} else {
			fmt.Println("ERROR, not an event")
		}
		if (verbose) {
			for i := range base_stations{
				fmt.Println(base_stations[i])
			}
		}

		// blocked_calls, dropped_calls := event.process_event(&FEL, base_stations) //now we dereference it to an Event

		num_blocked_calls += blocked_calls
		num_dropped_calls += dropped_calls

		// check if warm up is done but check it at an interval of every 10 calls
		if (!done_with_warmup) && (calls_since_previous_computation == computation_interval){
			previous_pct_blocked_calls = current_pct_blocked_calls
			previous_pct_dropped_calls = current_pct_dropped_calls
			current_pct_blocked_calls = num_blocked_calls/total_num_calls * 100
			current_pct_dropped_calls = num_dropped_calls/total_num_calls * 100

			// if performance metrics have converged => reset counters
			if ( //(total_num_calls > 100) && 
				(math.Abs((float64) (current_pct_blocked_calls - previous_pct_blocked_calls)) < 1) && 
				(math.Abs((float64) (current_pct_dropped_calls - previous_pct_dropped_calls)) < 1) ){
				fmt.Println("Done with warm up period! Resetting statistical counters.")
				done_with_warmup = true
				num_blocked_calls = 0
				num_dropped_calls = 0
				total_num_calls = 0
			}

			calls_since_previous_computation = 0
		}
		// determine warm up period by using convergence of performance metrics
	} 

	fmt.Println("pct of blocked calls", num_blocked_calls/total_num_calls * 100)
	fmt.Println("pct of dropped calls", num_dropped_calls/total_num_calls * 100)



	// ct := (Event) (&CallTermination{speed: 0.5, time: -500.0})
	// heap.Push(&FEL, &ct) // NOTE: remember to use head.Push NOT FEL.Push!
	// // FEL.Push(&ct)
	// // fmt.Println("size of pq: ", len(FEL))
	// for j := 0; j < 11; j ++ {
	// 	// fmt.Println((*FEL[len(FEL)-1]).get_time(), " ", (*FEL[len(FEL)-1]).get_index())
	// 	popped := heap.Pop(&FEL).(*Event) // Assert that the type is an Event pointer
	// 	fmt.Println((*popped).get_time(), " ", )
	// }
	// fmt.Println(FEL)

}