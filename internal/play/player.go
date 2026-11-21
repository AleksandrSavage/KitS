package play

import "sync"

type Player struct {
	mu     sync.RWMutex
	state  State
	engine *engine
	eventHandler EventHandler
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

    if p.state == state {
        p.mu.Unlock()
        return
    }

    p.state = state

    p.mu.Unlock()

    p.emit(Event{
        Type:  EventStateChanged,
        State: state.String(),
    })
}
