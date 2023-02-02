package main

import "container/heap"

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
	get_time() float64
	set_time(float64)
	set_index(int) // get the index of the item in the priority queue
	get_index() int // get the index of the item in the priority queue

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
