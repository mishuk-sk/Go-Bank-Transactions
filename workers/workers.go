package workers

import (
	"net/http"
	"sync"
)

type WorkersChan struct {
	source    chan interface{}
	listeners []func(interface{})
	sync.RWMutex
	quit chan struct{}
}

func (workCh *WorkersChan) CreateHttpWorker(h func(http.ResponseWriter, *http.Request, chan<- interface{})) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, chan<- interface{}(workCh.source))
	}
}

func (workCh *WorkersChan) CreateWorker(f func(chan<- interface{}, ...interface{})) func(...interface{}) {
	return func(v ...interface{}) {
		f(chan<- interface{}(workCh.source), v...)
	}
}

// FIXME multiply listeners!!! https://stackoverflow.com/questions/28527038/go-one-channel-with-multiple-listeners
func (workCh *WorkersChan) AddListener(f func(interface{})) {
	workCh.Lock()
	workCh.listeners = append(workCh.listeners, f)
	workCh.Unlock()
}

func (workCh *WorkersChan) Init() {
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

func (workCh *WorkersChan) Close() {
	close(workCh.source)
	close(workCh.quit)
}
