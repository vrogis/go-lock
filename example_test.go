package lock_test

import (
	"fmt"
	"github.com/vrogis/go-lock"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Locker struct {
	mtx sync.Mutex
}

type Locker2 struct {
	sync.Mutex
}

func (g *Locker) Lock() {
	g.mtx.Lock()
}

func (g *Locker) Unlock() {
	g.mtx.Unlock()
}

func Example() {
	var l1 sync.Mutex
	var l2, l3 Locker
	var l4, l5 Locker2

	unlock := lock.Lock(&l1, &l2)
	unlock()

	time.AfterFunc(time.Second*5, func() {
		stop.Store(true)
	})

	var wg sync.WaitGroup

	wg.Add(6)

	go lockCycle(&wg, "l1", &l3, &l1, &l2)
	go lockCycle(&wg, "l2", &l2, &l3, &l1, &l5, &l4)
	go lockCycle(&wg, "l3", &l4, &l1, &l3, &l2, &l5)
	go lockCycle(&wg, "l4", &l5, &l1, &l4, &l3, &l2)
	go lockCycle(&wg, "l5", &l1, &l4, &l5, &l3, &l2)

	lockCycle(&wg, "l6", &l5, &l3, &l2, &l1, &l4)

	wg.Wait()

	// Unordered Output:
	// l1
	// l2
	// l3
	// l4
	// l5
	// l6
}

var stop atomic.Bool

func lockCycle(wg *sync.WaitGroup, name string, locks ...sync.Locker) {
	for {
		unlock := lock.Lock(locks...)

		if stop.Load() {
			unlock()

			break
		}

		time.Sleep(time.Millisecond * time.Duration(200+rand.Int()%500))

		unlock()
	}

	fmt.Println(name)

	wg.Done()
}
