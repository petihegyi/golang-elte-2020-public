package main

import (
"fmt"
"sync"
"time"
)

type limited struct {
	mu    sync.Mutex
	limit int
	count int
	max   int
}

func (l *limited) check() {
	if l.count-l.limit > l.max {
		l.max = l.count - l.limit
	}
}

func (l *limited) start() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.count++ // HL
	l.check()
}

func (l *limited) end() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.count-- // HL
	l.check()
}

func (l *limited) process(i int) {
	l.start()
	defer l.end()
	time.Sleep(100 * time.Millisecond)
}

func main() {
	var work []int
	for i := 0; i < 1000; i++ {
		work = append(work, i)
	}
	start := time.Now()
	// TODO: avoid overloading the limited resource
	//       and try minimizing the total processing time
	// START OMIT
	l := &limited{limit: 100} // HL
	sem := make(chan struct{}, l.limit)
	done := make(chan struct{})
	for _, i := range work {
		sem <- struct{}{}
		go func(i int) {
			l.process(i) // HL
			<-sem
			done <- struct{}{}
		}(i)
	}
	for range work {
		<-done
	}
	// END OMIT
	fmt.Println(time.Since(start))
	if l.max > 0 {
		fmt.Println("maximum overload: ", l.max)
	}
}
