package lib

import "math"

func calculateReservation(mups []memoryUtilizationPoint) int64 {

	hours := make(map[int][]float64)

	for _, mup := range mups {
		hours[mup.time.Hour()] = append(hours[mup.time.Hour()], *mup.value)
	}

	avgHours := make([]float64, 24)

	for i, hour := range hours {
		avgHours[i] = Average(hour)
	}

	max := Max(avgHours)

	return int64(math.Round(max / .8))
}

func Max(array []float64) float64 {
	var max float64
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func Average(array []float64) float64 {
	sum := 0.0
	for _, value := range array {
		sum += value
	}
	return sum / float64(len(array))
}
