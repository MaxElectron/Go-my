package keylock

import "sync"

type Lock struct {
	workerMap map[string]map[chan string]map[string]struct{}
	keyStatus map[string]bool
	rwMutex   *sync.RWMutex
}

func New() *Lock {
	return &Lock{
		rwMutex:   &sync.RWMutex{},
		keyStatus: map[string]bool{},
		workerMap: map[string]map[chan string]map[string]struct{}{},
	}
}

func (lock *Lock) bind(keys []string) {
	lock.rwMutex.Lock()
	defer lock.rwMutex.Unlock()

	for _, key := range keys {
		lock.keyStatus[key] = false
	}
}

func (lock *Lock) deleteWorkers(keys []string, updateChannel chan string, lockMutex bool) {
	if lockMutex {
		lock.rwMutex.Lock()
		defer lock.rwMutex.Unlock()
	}

	for _, key := range keys {
		delete(lock.workerMap[key], updateChannel)
	}
}

func (lock *Lock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	lock.rwMutex.Lock()

	for _, key := range keys {
		_, exists := lock.keyStatus[key]
		if !exists {
			lock.keyStatus[key] = false
			lock.workerMap[key] = make(map[chan string]map[string]struct{})
		}
	}

	updateChannel := make(chan string, 1)
	keyBindings := make(map[string]struct{})
	for _, key := range keys {
		keyBindings[key] = struct{}{}
		lock.workerMap[key][updateChannel] = keyBindings
	}

	lock.rwMutex.Unlock()

	for {
		lock.rwMutex.Lock()

		allUnbound := true
		for key := range keyBindings {
			if lock.keyStatus[key] {
				allUnbound = false
				break
			}
		}

		if !allUnbound {
			lock.rwMutex.Unlock()
		} else {
			break
		}

		select {
		case <-updateChannel:
			continue
		case <-cancel:
			lock.deleteWorkers(keys, updateChannel, true)
			return true, nil
		}
	}

	for _, key := range keys {
		lock.keyStatus[key] = true
	}

	lock.deleteWorkers(keys, updateChannel, false)
	lock.rwMutex.Unlock()

	return false, func() {
		lock.bind(keys)

		lock.rwMutex.RLock()
		for _, key := range keys {
			for updateChannel := range lock.workerMap[key] {
				select {
				case updateChannel <- key:
				default:
				}
			}
		}
		lock.rwMutex.RUnlock()
	}
}
