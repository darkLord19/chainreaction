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

// IsCorner checks if given position is corner or not
func IsCorner(m int, n int, x int, y int) bool {
	if (x == 0 && y == 0) ||
		(x == m-1 && y == 0) ||
		(x == 0 && y == n-1) ||
		(x == m-1 && y == n-1) {
		return true
	}
	return false
}

// IsOnEdge checks if given position is on edge or not
func IsOnEdge(m int, n int, x int, y int) bool {
	if (x == 0 && y != 0 && y != n-1) ||
		(x != 0 && x != m-1 && y == 0) ||
		(x != 0 && x != m-1 && y == n-1) ||
		(x == m-1 && y != 0 && y != n-1) {
		return true
	}
	return false
}
