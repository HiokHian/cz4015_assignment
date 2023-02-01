package main

type CallInitiation struct {
	time      float64
	speed     float64
	station   *BaseStation
	position  float64
	duration  float64
	direction int32
}

type CallTermination struct {
	time  float64
	speed float64
}

type CallHandover struct {
	time      float64
	speed     float64
	station   *BaseStation
	duration  float64
	direction int32
}
