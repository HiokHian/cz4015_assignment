package main

import (
	"math"
	"math/rand"
)

// Chi square statistic for interarrival times:  111.66106610661065 Beta for exponential distribution:  1.3698169264765245
// Chi square:  Power_divergenceResult(statistic=111.66106610661068, pvalue=0.18115543265359552)
// Chi square statistic for durations:  97.02980298029806 Beta for exponential distribution:  99.83194913549607
// Chi square:  Power_divergenceResult(statistic=97.01999999999997, pvalue=0.5375175567180229)
// Chi square statistic for speed:  96.30973097309734 Mean for normal distribution:  120.07209801685805 Std for normal distribution:  9.01905789789691
// Chi square:  Power_divergenceResult(statistic=96.3, pvalue=0.558096218867612)
// KS test:  KstestResult(statistic=0.00642989259070903, pvalue=0.8002543274915996)
// KS test for uniformity:  KstestResult(statistic=0.055900000000000005, pvalue=1.332907138844621e-27)
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

