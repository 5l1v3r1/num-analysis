package kahan

func Sum64(nums []float64) float64 {
	var sum float64
	var compensation float64

	for _, num := range nums {
		num -= compensation
		roundedSum := sum + num
		compensation = (roundedSum - sum) - num
		sum = roundedSum
	}

	return sum
}
