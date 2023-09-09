package api

import "sync"

type Pool struct {
	mu      sync.Mutex
	players int
	max     int
}

func NewPool(max int) *Pool {
	p := Pool{
		players: 0,
		max:     max,
	}

	return &p
}

func (p *Pool) CanPlay() (bool, int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.players < p.max {
		p.players++

		return true, p.players
	}

	return false, p.players
}
