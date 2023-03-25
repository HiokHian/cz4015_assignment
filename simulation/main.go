package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

// TODO: verify if the declaring simulation_clock as
// a global variable here will allow it to be updated
// in the main function

// TODO: verify that base stations are edited in place with the subtraction of the channels
func main() {
	// define experiment settings
	var using_reserved_fca_scheme bool = true
	var num_reserved_channels int = 5

	// var write_running_pct_to_csv bool = true 
	// var write_final_pct_to_csv bool = false
	var num_expts int = 1

	var running_pct_file *os.File
	// var final_pct_file *os.File
	// set up file to write running pct to csv
	// for showing steady state

	running_pct_file, err := os.Create("logs/reservation_" + strconv.FormatBool(using_reserved_fca_scheme) + "_num_" + strconv.Itoa(num_reserved_channels) + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	defer running_pct_file.Close()
	_, err2 := running_pct_file.WriteString("pct_blocked_calls" + "," +  "pct_dropped_calls"  +"\n")
	if err2 != nil {
		log.Fatal(err2)
	}


	// final_pct_file, err := os.Create("results/reservation_" + strconv.FormatBool(using_reserved_fca_scheme) + "_num_" + strconv.Itoa(num_reserved_channels) + ".csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer final_pct_file.Close()
	// _, err2 := final_pct_file.WriteString("final_pct_blocked_calls" + "," +  "final_pct_dropped_calls"  +"\n")
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }


	for i := 0; i < num_expts; i++{
		fmt.Println("Running expt ", strconv.Itoa(i))
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
				num_free_channels: 10 - num_reserved_channels,
				num_free_reserved_channels: num_reserved_channels,
				using_reserved_fca_scheme: using_reserved_fca_scheme,
			}
		}

		//FEL needs to be a priority queue
		var total_num_initiation_calls = 10000 // same number of calls as given in csv
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
		var computation_interval int = 50
		var blocked_calls  float64
		var dropped_calls float64

		// FOR DEBUGGING 
		var verbose bool = false

		// //if need to write running pct to csv file
		// if (write_running_pct_to_csv){
		// 	_, err2 := running_pct_file.WriteString("pct_blocked_calls" + "," +  "pct_dropped_calls"  +"\n")
		// 	if err2 != nil {
		// 		log.Fatal(err2)
		// 	}
		// }

		for {
			// fmt.Println("Number of events left", len(FEL))
			if (len(FEL) == 0){ //  Stop iterating when FEL is empty
				break
			}

			total_num_calls += 1
			calls_since_previous_computation += 1
			
			if event, ok := heap.Pop(&FEL).(*Event); ok{
				if call_initiation, ok := (*event).(*CallInitiation); ok { //dereference to an event then cast to cal initiation pointer
					if (verbose){
						fmt.Println("CallInitiation", call_initiation)
					}
					blocked_calls, dropped_calls = (call_initiation).process_event(&FEL, base_stations) 
					// fmt.Println("CallInitiation: ", blocked_calls, dropped_calls)
				} else if call_handover, ok := (*event).(*CallHandover); ok{
					if (verbose){
						fmt.Println("Callhandover", call_handover)
					}
					blocked_calls, dropped_calls = (call_handover).process_event(&FEL, base_stations)
					// fmt.Println("Callhandover: ", blocked_calls, dropped_calls)
				} else if call_termination, ok := (*event).(*CallTermination); ok{
					if (verbose){
						fmt.Println("Calltermination", call_termination)
					}
					blocked_calls, dropped_calls = (call_termination).process_event(&FEL, base_stations) 
					// fmt.Println("Calltermination: ", blocked_calls, dropped_calls)
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

			// check if warm up is done but check it at an interval of every 50 calls
		
			_, err2 := running_pct_file.WriteString(strconv.FormatFloat(num_blocked_calls/total_num_calls * 100.0, 'E', -1, 64) + "," + strconv.FormatFloat(num_dropped_calls/total_num_calls * 100.0, 'E', -1, 64) +"\n")
			if err2 != nil {
				log.Fatal(err2)
			}


			if (!done_with_warmup) && (calls_since_previous_computation == computation_interval){
				previous_pct_blocked_calls = current_pct_blocked_calls
				previous_pct_dropped_calls = current_pct_dropped_calls
				current_pct_blocked_calls = num_blocked_calls/total_num_calls * 100
				current_pct_dropped_calls = num_dropped_calls/total_num_calls * 100

				// if performance metrics have converged => reset counters
				if ( (total_num_calls > 300) && 
					(current_pct_blocked_calls > 0) &&
					(current_pct_dropped_calls > 0) &&
					(math.Abs((float64) (current_pct_blocked_calls - previous_pct_blocked_calls)) < 0.1) && 
					(math.Abs((float64) (current_pct_dropped_calls - previous_pct_dropped_calls)) < 0.1) ){
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
		fmt.Println("total num of blocked calls", num_blocked_calls)
		fmt.Println("total num of dropped calls", num_dropped_calls)

		fmt.Println("pct of blocked calls", num_blocked_calls/total_num_calls * 100)
		fmt.Println("pct of dropped calls", num_dropped_calls/total_num_calls * 100)

		// _, err2 := final_pct_file.WriteString(strconv.FormatFloat(num_blocked_calls/total_num_calls * 100.0, 'E', -1, 64) + "," + strconv.FormatFloat(num_dropped_calls/total_num_calls * 100.0, 'E', -1, 64) +"\n")
		// if err2 != nil {
		// 	log.Fatal(err2)
		// }

	}
}