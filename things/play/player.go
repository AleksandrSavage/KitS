package play

type Player struct {
	// Здесь будут внутренние зависимости плеера.
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(path string) error {
	return nil
}

func (p *Player) Pause() error {
	return nil
}

func (p *Player) Resume() error {
	return nil
}

func (p *Player) Stop() error {
	return nil
}

func (p *Player) Close() error {
	return nil
}

func (p *Player) State() State {
	return StateIdle
}
