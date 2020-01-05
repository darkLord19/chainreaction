package utils

// Contains checks if a slice containse given value
func Contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}
