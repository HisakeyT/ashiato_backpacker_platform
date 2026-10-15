package handler

func MapSlice[T any, R any](items []T, convert func(T) R) []R {
	result := make([]R, 0, len(items))

	for _, item := range items {
		result = append(result, convert(item))
	}

	return result
}
