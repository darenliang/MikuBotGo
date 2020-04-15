package framework

func Index(str string, list []string) int {
	for i, val := range list {
		if val == str {
			return i
		}
	}
	return -1
}
