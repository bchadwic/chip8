package speaker

import "sync"

type Speaker interface {
	IsActive() bool
	Set(bool)
}

type speaker struct {
	active bool
	mu     sync.Mutex
}

func Create() Speaker {
	return &speaker{}
}

func (sp *speaker) IsActive() bool {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.active
}

func (sp *speaker) Set(active bool) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.active = active
}
