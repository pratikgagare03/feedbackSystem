package helper

import (
	"math"
	"strconv"
)

func CalculatePercentage(curr, total int) float64 {
	if total == 0 {
		return 0
	}
	percentage := float64(curr) / float64(total) * 100
	return math.Round(percentage*100) / 100
}

func CalculateAverage(ratingsCount map[string]int) float64 {
	sum, count := 0, 0
	for rating, cnt := range ratingsCount {
		ratingInt, _ := strconv.Atoi(rating)
		sum += cnt * ratingInt
		count += cnt
	}
	if count == 0 {
		return 0
	}
	avg := float64(sum) / float64(count)
	return math.Round(avg*100) / 100
}
