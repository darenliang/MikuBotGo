package framework

// Index returns index of string found in list
// return -1 if not found
func Index(str string, list []string) int {
	for i, val := range list {
		if val == str {
			return i
		}
	}
	return -1
}
