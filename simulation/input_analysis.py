import pandas as pd
import numpy as np
from scipy import stats

df = pd.read_excel("PCS_TEST_DETERMINSTIC.xls")

interarrival_times = df["Arrival time (sec)"].iloc[1:].reset_index(drop=True) - df[
    "Arrival time (sec)"
].iloc[:-1].reset_index(drop=True)
interarrival_times = interarrival_times.tolist()


""" 
Here the number of samples s is 10,000, so the recommended k number of intervals is between 100 and 20,000
according to the lecture slides.
The DOF = k - s - 1 where s is the number of parameters in the hypothesized distribution
So s = 1 for exponential distribution and s = 2 for normal distribution.
the probabilities for each interval = 1/k where k is the number of intervals
Here, we can choose k = 1000, which will be the middle ground between the 2 bounds for k

We hypothesise the following:
1) Interarrival times are exponentially distributed with mean = mean of the interarrival times
2) Call durations (left shifted to by the minimum of the call durations) => normal?
3) Speeds are normally distributed with a mean and std equal to those of the data
4) Base stations are uniformly distributed as stated in the question.
"""


# Interarrival time
# since we let k = 1000, p_j = 1/1000 => n*p_j = 10
# k - s -1 = 998, and we choose the significance to be 0.95
k = 100
n = len(interarrival_times)
interval_prob = 1 / k
expected_counts_per_interval = n * interval_prob
beta = sum(interarrival_times) / len(interarrival_times)
interval_end_pts = [beta * np.log(1 / (1 - i * interval_prob)) for i in range(k)] + [
    np.inf
]
chi_square_statistic = 0
observed_frequencies = []
for left, right in zip(interval_end_pts[:-1], interval_end_pts[1:]):
    filtered = [time for time in interarrival_times if left <= time < right]
    interval_actual_count = len(filtered)
    observed_frequencies.append(interval_actual_count)
    interval_chi_square_statistic = (
        interval_actual_count - expected_counts_per_interval
    ) ** 2 / expected_counts_per_interval
    chi_square_statistic += interval_chi_square_statistic
print(
    "Chi square statistic for interarrival times: ",
    chi_square_statistic,
    "Beta for exponential distribution: ",
    beta,
)
print("Chi square: ", stats.chisquare(observed_frequencies))


# Call duration
min_call_duration = df["Call duration (sec)"].min()
df["Call duration (sec)"] = df["Call duration (sec)"] - min_call_duration
beta = df["Call duration (sec)"].mean()
interval_end_pts = [beta * np.log(1 / (1 - i * interval_prob)) for i in range(k)] + [
    np.inf
]
call_durations = df["Call duration (sec)"]
chi_square_statistic = 0
observed_frequencies = []
for left, right in zip(interval_end_pts[:-1], interval_end_pts[1:]):
    filtered = [duration for duration in call_durations if left <= duration < right]
    interval_actual_count = len(filtered)
    observed_frequencies.append(interval_actual_count)
    interval_chi_square_statistic = (
        interval_actual_count - expected_counts_per_interval
    ) ** 2 / expected_counts_per_interval
    chi_square_statistic += interval_chi_square_statistic
print(
    "Chi square statistic for durations: ",
    chi_square_statistic,
    "Beta for exponential distribution: ",
    beta,
)
print("Chi square: ", stats.chisquare(observed_frequencies))

# Speed
mean = df["velocity (km/h)"].mean()
std = df["velocity (km/h)"].std()
interval_end_pts = [
    stats.norm.ppf(i * interval_prob, loc=mean, scale=std) for i in range(k)
] + [np.inf]
speeds = df["velocity (km/h)"].tolist()
chi_square_statistic = 0
observed_frequencies = []
for left, right in zip(interval_end_pts[:-1], interval_end_pts[1:]):
    filtered = [speed for speed in speeds if left <= speed < right]
    interval_actual_count = len(filtered)
    observed_frequencies.append(interval_actual_count)
    interval_chi_square_statistic = (
        interval_actual_count - expected_counts_per_interval
    ) ** 2 / expected_counts_per_interval
    chi_square_statistic += interval_chi_square_statistic
print(
    "Chi square statistic for speed: ",
    chi_square_statistic,
    "Beta for normal distribution: ",
    beta,
)
ks = stats.kstest(speeds, lambda x: stats.norm.cdf(x, loc=mean, scale=std))
print("Chi square: ", stats.chisquare(observed_frequencies))
print("KS test: ", ks)

# Base station
base_stations = df["Base station "].tolist()
ks = stats.kstest(base_stations, lambda x: stats.uniform.cdf(x, loc=0, scale=20))
print("KS test for uniformity: ", ks)
