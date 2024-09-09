package slicex

func IntersectionInt(x []int, y []int) []int {
	if len(x) == 0 || len(y) == 0 {
		return []int{}
	}

	var result []int
	xMap := make(map[int]struct{}, len(x))
	for _, v := range x {
		xMap[v] = struct{}{}
	}
	for _, v := range y {
		if _, ok := xMap[v]; ok {
			result = append(result, v)
		}
	}
	return result
}
