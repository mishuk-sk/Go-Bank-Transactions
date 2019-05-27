package workers

import (
	"net/http"
	"sync"
)

// Observer provides way to attach multiple listeners to single channel, where lots of dispatchers writes
type Observer struct {
	source chan interface{}
	// event listeners
	listeners []func(interface{})
	sync.RWMutex
	quit chan struct{}
}

//CreateHTTPWorker allows to add new worker with http.HandlerFunc signature
func (workCh *Observer) CreateHTTPWorker(h func(http.ResponseWriter, *http.Request, chan<- interface{})) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, chan<- interface{}(workCh.source))
	}
}

//CreateWorker allows to create any worker
func (workCh *Observer) CreateWorker(f func(chan<- interface{}, ...interface{})) func(...interface{}) {
	return func(v ...interface{}) {
		f(chan<- interface{}(workCh.source), v...)
	}
}

//AddListener adds listener to all events in this Observer instance. Listener should accept some parameters and mustn't return anything
func (workCh *Observer) AddListener(f func(interface{})) {
	workCh.Lock()
	workCh.listeners = append(workCh.listeners, f)
	workCh.Unlock()
}

//Init initializes Observer for future work. !!!Observer must be closed by calling Observer.Close() after you're done with it
func (workCh *Observer) Init() {
	// FIXME what about channel being too full???
	workCh.source = make(chan interface{}, 10)
	workCh.quit = make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-workCh.source:
				for _, f := range workCh.listeners {
					f(msg)
				}
			case <-workCh.quit:
				return
			}
		}
	}()
}

//Close closes Observer correctly, preventing future resources waste
func (workCh *Observer) Close() {
	close(workCh.source)
	close(workCh.quit)
}
