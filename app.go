package main

import (
	"context"
	"kitS/things/play"

)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context)  {
	a.ctx = ctx
}


func (a *App) PlayAudio(path string) error {
	return play.PlayAudio(path)
}

func (a *App) StopAudio() {
    play.StopAudio()
}

func (a *App) PauseAudio() {
	play.PauseAudio()
}

func (a *App) ResumeAudio() {
	play.ResumeAudio()
}


