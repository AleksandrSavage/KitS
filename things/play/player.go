package play

import "sync"

type Player struct {
	// Здесь будут внутренние зависимости плеера.
	mu sync.RWMutex
	state State
}

func NewPlayer() *Player {
	return &Player{
		state: StateIdle,
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