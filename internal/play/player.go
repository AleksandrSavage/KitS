package play

import ("sync"
        "time"
)

type Player struct {
	mu     sync.RWMutex
	state  State
	engine *engine
	eventHandler EventHandler
	opMu sync.Mutex 
}

type PlaybackProgress struct {
    PositionMs int64 `json:"positionMs"`
    DurationMs int64 `json:"durationMs"`
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


func (p *Player) Play(path string) error {
    p.opMu.Lock()
    defer p.opMu.Unlock()

    if err := p.engine.stop(); err != nil {
        p.emitError(err)
        p.setState(StateError)
        return err
    }

    p.setState(StateLoading)

    audio, err := openAudio(path)
    if err != nil {
        p.emitError(err)
        p.setState(StateError)
        return err
    }

    playback, err := resampleAudio(
        audio.streamer,
        audio.format.SampleRate,
        targetSampleRate,
    )
    if err != nil {
        _ = audio.streamer.Close()
        p.emitError(err)
        p.setState(StateError)
        return err
    }

    if err := p.engine.start(audio, playback); err != nil {
        _ = audio.streamer.Close()
        p.emitError(err)
        p.setState(StateError)
        return err
    }

    p.setState(StatePlaying)

    return nil
}


func (p *Player) Pause() error {
    p.opMu.Lock()
    defer p.opMu.Unlock()

    if err := p.engine.pause(); err != nil {
        return err
    }

    p.setState(StatePaused)

    return nil
}


func (p *Player) Resume() error {
    p.opMu.Lock()
    defer p.opMu.Unlock()

    if err := p.engine.resume(); err != nil {
        return err
    }

    p.setState(StatePlaying)

    return nil
}


func (p *Player) Stop() error {
    p.opMu.Lock()
    defer p.opMu.Unlock()

    if err := p.engine.stop(); err != nil {
        p.emitError(err)
        p.setState(StateError)
        return err
    }

    p.setState(StateStopped)

    return nil
}


func (p *Player) Close() error {
    return p.Stop()
}

func (p *Player) Progress() PlaybackProgress {
    position, duration := p.engine.getProgress()

    return PlaybackProgress{
        PositionMs: position.Milliseconds(),
        DurationMs: duration.Milliseconds(),
    }
}

func (p *Player) Seek(position time.Duration) error {
    p.opMu.Lock()
    defer p.opMu.Unlock()

    if err := p.engine.seek(position); err != nil {
        return err
    }

    return nil
}