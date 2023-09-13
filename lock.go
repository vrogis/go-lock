package lock

import (
	"sort"
	"sync"
	"unsafe"
)

func Lock(lockers ...sync.Locker) (unlock func()) {
	count := len(lockers)

	if 0 == count {
		return func() {}
	}

	sort.Slice(lockers, func(i, j int) bool {
		return lockerToPtr(&lockers[i]) < lockerToPtr(&lockers[j])
	})

	for _, locker := range lockers {
		locker.Lock()
	}

	return func() {
		for i := count - 1; i >= 0; i-- {
			lockers[i].Unlock()
		}
	}
}

func lockerToPtr(locker *sync.Locker) uintptr {
	return (*struct{ uintptr, value uintptr })(unsafe.Pointer(locker)).value
}
