package utils

// Mutex to protect concurrent read writes
type Mutex chan struct{}

// Lock the kraken
func (m Mutex) Lock() {
	<-m
}

// Unlock the kraken
func (m Mutex) Unlock() {
	m <- struct{}{}
}
