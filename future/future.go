package future

import "sync"

func New() *Future {
	return &Future{mutex: sync.RWMutex{}}
}

type Future struct {
	resolved bool
	waiters  []chan bool
	mutex    sync.RWMutex
}

func (f *Future) Done() {
	f.mutex.Lock()
	if !f.resolved {
		f.resolved = true
		for _, w := range f.waiters {
			w <- true
		}
		f.waiters = nil
	}
	f.mutex.Unlock()
}

func (f *Future) Then() {
	f.mutex.RLock()
	resolved := f.resolved
	f.mutex.RUnlock()

	if resolved {
		return
	}

	ch := make(chan bool)

	f.mutex.Lock()
	f.waiters = append(f.waiters, ch)
	f.mutex.Unlock()

	<-ch
}
