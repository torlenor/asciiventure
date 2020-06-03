package utils

func MaxInt32(x, y int32) int32 {
	if x < y {
		return y
	}
	return x
}

func MinInt32(x, y int32) int32 {
	if x > y {
		return y
	}
	return x
}

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}
