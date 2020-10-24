package semaphore

/**
* use buff channel implement semaphore
* @see http://www.golangpatterns.info/concurrency/Semaphores
 */

// Empty is empty interface type
type Empty interface{}

// Semaphore is empty type chan
type Semaphore chan Empty

// P used to acquire n resources
func (s Semaphore) P(n int) {
	e := new(Empty)
	for i := 0; i < n; i++ {
		s <- e
	}
}

// V used to release n resouces
func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

// Lock used to lock resource
func (s Semaphore) Lock() {
	s.P(1)
}

// Unlock used to unlock resource
func (s Semaphore) Unlock() {
	s.V(1)
}

// Wait used to wait signal
func (s Semaphore) Wait(n int) {
	s.P(n)
}

// Signal used to send signal
func (s Semaphore) Signal() {
	s.V(1)
}

// NewSemaphore return semaphore
func NewSemaphore(N int) Semaphore {
	return make(Semaphore, N)
}
