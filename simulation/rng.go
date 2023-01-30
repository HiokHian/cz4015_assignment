package main

import (
	"math"
	"math/rand"
)

type RNG interface {
    sample() float64
}

type UniformDistribution struct {
    left, right float64
}

type ExponentialDistribution struct {
    beta float64
}

type NormalDistribution struct {
    mean, std float64
}

func (rng UniformDistribution) sample() float64 {
    return rand.Float64() * (rng.right - rng.left) + rng.left
}

func (rng ExponentialDistribution) sample() float64 {
    return -rng.beta*math.Log(1- rand.Float64())
}

func (rng NormalDistribution) sample() float64 {
    sum := 0.0
    num_samples := 100
    for i := 0; i < num_samples; i++ {
        sum += rand.Float64()
    }
    //Uniform(0, 1) have mean = 0.5 and variance 1/12

    //by CLT, the sample mean should be approx ~N(mu, variance/n)
    sample_mean := sum/float64(num_samples)
    z := (sample_mean-0.5)/math.Sqrt((1.0/12.0)/ float64(num_samples))

    return z * rng.std + rng.mean // to change to larger sample size
}

// func (r rect) perim() float64 {
//     return 2*r.width + 2*r.height
// }

// func (c circle) area() float64 {
//     return math.Pi * c.radius * c.radius
// }
// func (c circle) perim() float64 {
//     return 2 * math.Pi * c.radius
// }

// func measure(g geometry) {
//     fmt.Println(g)
//     fmt.Println(g.area())
//     fmt.Println(g.perim())
// }

