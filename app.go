
package main

import (
    "context"

    "kitS/internal/play"

    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx    context.Context
    player *play.Player
}

func NewApp() *App {
    return &App{
        player: play.NewPlayer(),
    }
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx

    a.player.SetEventHandler(func(event play.Event) {
        runtime.EventsEmit(
            a.ctx,
            string(event.Type),
            event,
        )
    })
}

func (a *App) shutdown(ctx context.Context) {
    if err := a.player.Close(); err != nil {
        runtime.LogError(ctx, err.Error())
    }
}

func (a *App) Play(path string) error {
    return a.player.Play(path)
}

func (a *App) Pause() error {
    return a.player.Pause()
}

func (a *App) Resume() error {
    return a.player.Resume()
}

func (a *App) Stop() error {
    return a.player.Stop()
}

func (a *App) GetState() string {
    return a.player.State().String()
}
