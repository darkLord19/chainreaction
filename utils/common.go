package utils

import "sync"

// Pair struct to
type Pair struct {
	X int
	Y int
}

// Queue the queue of int pairs
type Queue struct {
	items []Pair
	lock  sync.RWMutex
}

// NewQueue creates a new Queue
func NewQueue() Queue {
	s := Queue{}
	s.items = []Pair{}
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *Queue) Enqueue(t Pair) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes a Pair from the start of the queue
func (s *Queue) Dequeue() (int, int) {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return item.X, item.Y
}

// IsEmpty returns true if the queue is empty
func (s *Queue) IsEmpty() bool {
	return len(s.items) == 0
}

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
