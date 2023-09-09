package pubsub

import "sync"

type PubSub struct {
	mu       sync.RWMutex
	channels map[chan struct{}]struct{}
}

func NewPubSub() *PubSub {
	return &PubSub{
		channels: make(map[chan struct{}]struct{}),
	}
}

func (p *PubSub) Subscribe() (<-chan struct{}, func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c := make(chan struct{}, 1)
	p.channels[c] = struct{}{}

	fn := func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		delete(p.channels, c)
		close(c)
	}

	return c, fn
}

func (p *PubSub) Publish() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for channel := range p.channels {
		channel <- struct{}{}
	}
}
