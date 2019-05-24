package subhandler

import (
	"net/http"
)

type WorkersChan struct {
	source chan interface{}
	quit   chan struct{}
}

func (workCh *WorkersChan) AddWorker(h func(http.ResponseWriter, *http.Request, chan<- interface{})) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, chan<- interface{}(workCh.source))
	}
}

func (workCh *WorkersChan) AddListener(l func(interface{})) {
	go func() {
		for {
			select {
			case msg := <-workCh.source:
				l(msg)
			case <-workCh.quit:
				return
			}
		}
	}()
}

func (workCh *WorkersChan) Init() {
	// FIXME what about channel being too full???
	workCh.source = make(chan interface{}, 10)
	workCh.quit = make(chan struct{})
}

func (workCh *WorkersChan) Close() {
	close(workCh.source)
	close(workCh.quit)
}
