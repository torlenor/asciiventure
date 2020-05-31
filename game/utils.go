package game

func max(x, y int32) int32 {
	if x < y {
		return y
	}
	return x
}

func min(x, y int32) int32 {
	if x > y {
		return y
	}
	return x
}
