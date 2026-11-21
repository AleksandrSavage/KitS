
package play

type EventType string

const (
    EventStateChanged EventType = "state_changed"
    EventTrackEnded   EventType = "track_ended"
    EventError        EventType = "error"
)

type Event struct {
    Type  EventType `json:"type"`
    State string    `json:"state,omitempty"`
    Error string    `json:"error,omitempty"`
}

type EventHandler func(Event)

func (p *Player) SetEventHandler(handler EventHandler) {
    p.mu.Lock()
    p.eventHandler = handler
    p.mu.Unlock()
}

func (p *Player) emit(event Event) {
    p.mu.RLock()
    handler := p.eventHandler
    p.mu.RUnlock()

    if handler != nil {
        handler(event)
    }
}

func (p *Player) emitError(err error) {
    if err == nil {
        return
    }

    p.emit(Event{
        Type:  EventError,
        State: StateError.String(),
        Error: err.Error(),
    })
}
