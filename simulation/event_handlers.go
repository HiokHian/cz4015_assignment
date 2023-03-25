package main

import (
	"container/heap"
)

// Each call initiation will have no "previous_staion", only current station
// During call initiation process_event has to check if there is a
type CallInitiation struct {
	time      float64
	speed     float64
	current_station   *BaseStation
	position  float64 // define as the percentage of the way to the right end of the station
	duration  float64
	direction int
	index     int // for heap interface
}

type CallTermination struct {
	time  float64
	previous_station   *BaseStation
	was_reserved_channel bool
	index int // for heap interface
}

type CallHandover struct {
	time      float64
	speed     float64
	previous_station *BaseStation //maintain ptr to previous to free up occupied channel
	current_station   *BaseStation
	duration  float64
	direction int
	was_reserved_channel bool // whether the previous channel occurpied was a reserved channel
	index     int // for heap interface
}

// getter methods are not required in golang but in order to have a FEL with
// different types in it, we use this interface
type Event interface {
	process_event(*PriorityQueue, []BaseStation) (float64, float64) // process the event within this function
	get_time() float64
	set_time(float64)
	set_index(int) // get the index of the item in the priority queue
	get_index() int // set the index of the item in the priority queue
}

func (ci *CallInitiation) process_event(FEL *PriorityQueue, all_base_stations []BaseStation) (float64, float64) {
	// return the number of blocked calls, dropped calls and the simulation clock
	// https://stackoverflow.com/questions/39993688/are-slices-passed-by-value
	// we can just get base station as is from this 
	simulation_clock := ci.time
	current_station := ci.current_station
	var next_base_station *BaseStation 
	var was_reserved_channel bool
	var num_blocked_calls float64 = 0.0
	var num_dropped_calls float64 = 0.0
	var distance_remaining_in_current_station float64

	if current_station.has_available_non_reserved_channel(){
		
		current_station.assign_non_reserved_channel() // call initiation will always take up a non-reserved channel
		was_reserved_channel = false // current channel assigned to call initiation is non-reserved
		
		if ci.direction > 0 { // car is heading to the right
			distance_remaining_in_current_station = (1 - ci.position) * 2 // each base station is 2 km
		} else { // car is heading to the left
			distance_remaining_in_current_station = ci.position * 2
		}

		time_left_in_station := distance_remaining_in_current_station / ci.speed // in hours
		time_left_in_station = time_left_in_station * 60* 60

		if (current_station.id != 0) && (current_station.id != 19) && (ci.duration > time_left_in_station) {
			if ci.direction > 0{
				next_base_station = &(all_base_stations[current_station.id + 1])
			} else {
				next_base_station = &(all_base_stations[current_station.id - 1])
			}
			new_event := (Event)(&CallHandover{
				time: simulation_clock + time_left_in_station,
				speed: ci.speed,
				previous_station: current_station,
				current_station: next_base_station,
				duration: ci.duration - time_left_in_station,
				direction: ci.direction,
				was_reserved_channel: was_reserved_channel,
				index: len(*FEL),
			})
			heap.Push(FEL, &new_event)
		} else {
			// if current station is at either ends, there will not be any handover, it will simply terminate
			// if the amount of time left in the station is <= the duratin of the call, it will simply terminate
			new_event :=(Event)(&CallTermination{
				time: simulation_clock + time_left_in_station,
				previous_station: current_station,
				was_reserved_channel: was_reserved_channel,
				index: len(*FEL),
			})
			
			heap.Push(FEL, &new_event)
		}
	} else { //if no more channel in current station, drop the call
		num_blocked_calls += 1
	}
	return num_blocked_calls, num_dropped_calls
}

func (ct *CallTermination) process_event(FEL *PriorityQueue, all_base_stations []BaseStation) (float64, float64) {
	//call termination events correspond to either:
	// 1. calls naturally finishing while the car is still in one of the 20 base stations
	// 2. car exiting the area of all 20 stations through either the left or right end
	previous_station := ct.previous_station
	var num_blocked_calls float64 = 0.0
	var num_dropped_calls float64 = 0.0
	if ct.was_reserved_channel {
		previous_station.free_up_reserved_channel()
	} else {
		previous_station.free_up_non_reserved_channel()
	}
	return num_blocked_calls, num_dropped_calls
}

func (ch *CallHandover) process_event(FEL *PriorityQueue, all_base_stations []BaseStation) (float64, float64) {
	simulation_clock := ch.time
	previous_station := ch.previous_station
	current_station := ch.current_station
	var num_blocked_calls float64 = 0.0
	var num_dropped_calls float64 = 0.0
	var next_base_station *BaseStation 
	var new_channel_is_reserved_channel bool
	var distance_remaining_in_current_station float64
	var current_station_has_available_channel bool

	// free up channel previously occupied in the previous station
	if (ch.was_reserved_channel) {
		previous_station.free_up_reserved_channel()
	} else {
		previous_station.free_up_non_reserved_channel()
	}

	// for the current station we want to handover to, check whether to handover
	// to a reserved or non reserved channel
	if (current_station.using_reserved_fca_scheme) && (!current_station.has_available_non_reserved_channel()){
		// allow usage of reserved channels if using reservation scheme and no more non-reserved channels
		new_channel_is_reserved_channel = true
	} else {
		new_channel_is_reserved_channel = false
	}

	// if we intend to use a reserved channel, then we check if there are available reserved channels
	if (new_channel_is_reserved_channel) {
		if (current_station.has_available_reserved_channel()){
			current_station_has_available_channel = true
		} else {
			current_station_has_available_channel = false
		}
	} else { // if we intend to use a non reserved channel, then we check if there are available non-reserved channels
		if (current_station.has_available_non_reserved_channel()){
			current_station_has_available_channel = true
		} else {
			current_station_has_available_channel = false
		}
	}

	if (current_station_has_available_channel){
		distance_remaining_in_current_station = 2.0 // full base station distance in km
		time_left_in_station := distance_remaining_in_current_station / ch.speed //in hours since speed is km/hr
		// need to convert hours to seconds
		time_left_in_station = time_left_in_station * 60* 60

		if (new_channel_is_reserved_channel){
			current_station.assign_reserved_channel()
		} else {
			current_station.assign_non_reserved_channel()
		}

		if (current_station.id != 0) && (current_station.id != 19) && (ch.duration > time_left_in_station) {
			if ch.direction > 0{
				next_base_station = &(all_base_stations[current_station.id + 1])
			} else {
				next_base_station = &(all_base_stations[current_station.id - 1])
			}

			new_event :=(Event) (&CallHandover{
				time: simulation_clock + time_left_in_station,
				speed: ch.speed,
				previous_station: current_station,
				current_station: next_base_station,
				duration: ch.duration - time_left_in_station,
				direction: ch.direction,
				was_reserved_channel: new_channel_is_reserved_channel,
				index: len(*FEL),
			})
			heap.Push(FEL, &new_event)
		} else {
			// if current station is at either ends, there will not be any handover, it will simply terminate
			// if the amount of time left in the station is <= the duratin of the call, it will simply terminate
			new_event := (Event) (&CallTermination{
				time: simulation_clock + time_left_in_station,
				previous_station: current_station,
				was_reserved_channel: new_channel_is_reserved_channel,
				index: len(*FEL),
			})
			
			heap.Push(FEL, &new_event)
		}
	} else { //no available channels to handover to, so we drop the call
		num_dropped_calls += 1
	}

	return num_blocked_calls, num_dropped_calls
}


func (ci *CallInitiation) get_time() float64 {
	return ci.time
}

func (ct *CallTermination) get_time() float64 {
	return ct.time
}

func (ch *CallHandover) get_time() float64 {
	return ch.time
}

func (ci *CallInitiation) get_index() int {
	return ci.index
}

func (ct *CallTermination) get_index() int {
	return ct.index
}

func (ch *CallHandover) get_index() int {
	return ch.index
}

func (ci *CallInitiation) set_index(idx int) {
	ci.index = idx
}


func (ct *CallTermination) set_index(idx int) {
	ct.index = idx
}

func (ch *CallHandover) set_index(idx int) {
	ch.index = idx
}

func (ci *CallInitiation) set_time(time float64) {
	ci.time = time
}

func (ct *CallTermination) set_time(time float64) {
	ct.time = time
}

func (ch *CallHandover) set_time(time float64) {
	ch.time = time
}

type PriorityQueue []*Event

// each struct implements Event but not the pointer
// each element in the pq is a pointer to an Event not an Event
// each operation the pq expects a 
// https://pkg.go.dev/container/heap
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return (*pq[i]).get_time() < (*pq[j]).get_time()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	(*pq[i]).set_index(i)
	(*pq[j]).set_index(j)
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Event)
	(*item).set_index(n)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil        // avoid memory leak
	(*item).set_index(-1) // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Event, time float64) {
	(*item).set_time(time)
	heap.Fix(pq, (*item).get_index())
}

// https://www.reddit.com/r/golang/comments/dyicuk/slices_take_in_multiple_types/
// https://go.dev/play/p/330yHpjRMLG
// func process_event(event_ptr *Event, FEL *PriorityQueue) int32, int32{
// 	switch event_type := (*event_ptr).(type) {
// 		case *CallInitiation:
// 			return handle_call_initiation(event_ptr, FEL)
// 		case *CallHandover:
// 			return handle_call_handover(event_ptr, FEL)
// 		case *CallTermination:
// 			return handle_call_termination(event_ptr, FEL)
// 	}
// }





// func handle_call_initiation(event *Event, FEL *PriorityQueue){
	
// 	simulation_clock =  event.get_time() // Required for later adding a new event to FEL
// 	current_station := ((*event).(*CallInitiation)).station //ptr to base_station
// 	num_blocked_calls := 0
// 	num_dropped_calls := 0

// }



