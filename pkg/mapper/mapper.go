package mapper

func MapSliceWithFunc[R, D any](rows []R, fn func(R) D) []D {
	out := make([]D, 0, len(rows))
	for _, row := range rows {
		out = append(out, fn(row))
	}
	return out
}

func MapSlice[T any](rows []T) []T {
	out := make([]T, 0, len(rows))
	for _, row := range rows {
		out = append(out, row)
	}
	return out
}
