package main

import (
	"container/heap"
)

type CallInitiation struct {
	time      float64
	speed     float64
	station   *BaseStation
	position  float64
	duration  float64
	direction int32
	index     int // for heap interface
}

type CallTermination struct {
	time  float64
	speed float64
	index int // for heap interface
}

type CallHandover struct {
	time      float64
	speed     float64
	station   *BaseStation
	duration  float64
	direction int32
	index     int // for heap interface
}

// getter methods are not required in golang but in order to have a FEL with
// different types in it, we use this interface
type Event interface {
	process() // process the event within this function
	get_time() float64
	set_time(float64)
	set_index(int) // get the index of the item in the priority queue
	get_index() int // set the index of the item in the priority queue
}

func (ci *CallInitiation) process_event() {

}

func (ci *CallTermination) process_event() {
	
}

func (ci *CallHandover) process_event() {
	
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

// TODO: verify if the declaring simulation_clock as
// a global variable here will allow it to be updated 
// in the main function 
var simulation_clock float64

// func handle_call_initiation(event *Event, FEL *PriorityQueue){
	
// 	simulation_clock =  event.get_time() // Required for later adding a new event to FEL
// 	current_station := ((*event).(*CallInitiation)).station //ptr to base_station
// 	num_blocked_calls := 0
// 	num_dropped_calls := 0

// }



