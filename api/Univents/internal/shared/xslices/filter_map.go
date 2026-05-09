package xslices

func FilterMap[T, U any](s []T, f func(T) (U, error)) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		if u, err := f(v); err == nil {
			out = append(out, u)
		}
	}
	return out
}
