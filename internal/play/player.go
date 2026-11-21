package play

import "sync"

type Player struct {
	mu     sync.RWMutex
	state  State
	engine *engine
}

func NewPlayer() *Player {
	return &Player{
		state:  StateIdle,
		engine: &engine{},
	}
}

func (p *Player) State() State {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.state
}

func (p *Player) setState(state State) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.state = state
}
