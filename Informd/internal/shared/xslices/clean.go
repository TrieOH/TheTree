package xslices

import "strings"

func Clean(ss []string) []string {
	if len(ss) == 0 {
		return nil
	}

	out := ss[:0]

	for _, s := range ss {
		if v := strings.TrimSpace(s); v != "" {
			out = append(out, v)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
