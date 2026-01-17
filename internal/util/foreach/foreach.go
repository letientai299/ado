package foreach

func Map[I, O any](in []I, fn func(I) O) []O {
	result := make([]O, 0, len(in))
	for _, item := range in {
		result = append(result, fn(item))
	}
	return result
}
