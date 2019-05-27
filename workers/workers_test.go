package workers

import (
	"fmt"
	"sync"
	"testing"
)

func TestListeners(t *testing.T) {
	var arr []string
	workers := WorkersChan{}
	workers.Init()
	defer workers.Close()
	str1 := "Hello"
	str2 := "World"
	fun := workers.CreateWorker(func(ch chan<- interface{}, v ...interface{}) {
		reply := fmt.Sprint(v...)
		ch <- reply
	})
	wg := sync.WaitGroup{}
	wg.Add(2)
	workers.AddListener(func(str interface{}) {
		s := str.(string)
		arr = append(arr, s)
		wg.Done()
	})
	workers.AddListener(func(str interface{}) {
		s := str.(string)
		arr = append(arr, s)
		wg.Done()
	})
	fun(str1, str2)
	wg.Wait()
	if len(arr) != 2 {
		t.Errorf("Not all listeners takes all events from dispatchers. Result arr %v", arr)
	}
}
