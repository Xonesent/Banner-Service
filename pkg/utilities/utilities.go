package utilities

func InStringSlice(target string, src []string) bool {
	for _, el := range src {
		if el == target {
			return true
		}
	}
	return false
}
