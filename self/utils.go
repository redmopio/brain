package self

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
