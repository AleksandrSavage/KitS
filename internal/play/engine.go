package play

import (
	"fmt"
	"sync"
	"time"
	"sync/atomic"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type engine struct {
	mu      sync.Mutex
	current *decodedAudio
	ctrl    *beep.Ctrl
	progress *progressStreamer
    duration time.Duration
    ended    bool
}

type progressStreamer struct {
    source  beep.Streamer
    samples int64
}

const (
	targetSampleRate      beep.SampleRate = 48000
	speakerBufferDuration                 = 100 * time.Millisecond
)

var (
	speakerOnce    sync.Once
	speakerInitErr error
)

func initSpeaker() error {
	speakerOnce.Do(func() {
		bufferSize := targetSampleRate.N(speakerBufferDuration)

		speakerInitErr = speaker.Init(
			targetSampleRate,
			bufferSize,
		)
	})

	return speakerInitErr
}

func (e *engine) start(
	audio *decodedAudio,
	playback beep.Streamer,
) error {
	if audio == nil || audio.streamer == nil {
		return fmt.Errorf("audio is nil")
	}

	if playback == nil {
		return fmt.Errorf("playback streamer is nil")
	}

	if err := initSpeaker(); err != nil {
		return fmt.Errorf("initialize speaker: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ctrl != nil {
		speaker.Lock()
		e.ctrl.Streamer = nil
		speaker.Unlock()
	}

	if e.current != nil {
		_ = e.current.streamer.Close()
	}

	ctrl := &beep.Ctrl{
		Streamer: playback,
	}

	e.current = audio
	e.ctrl = ctrl

	speaker.Play(ctrl)

	return nil
}

func (e *engine) pause() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ctrl == nil {
		return fmt.Errorf("nothing is playing")
	}

	speaker.Lock()
	e.ctrl.Paused = true
	speaker.Unlock()

	return nil
}

func (e *engine) resume() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ctrl == nil {
		return fmt.Errorf("nothing is playing")
	}

	speaker.Lock()
	e.ctrl.Paused = false
	speaker.Unlock()

	return nil
}

func (e *engine) stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ctrl != nil {
		speaker.Lock()
		e.ctrl.Streamer = nil
		speaker.Unlock()
	}

	audio := e.current

	e.ctrl = nil
	e.current = nil

	if audio != nil {
		return audio.streamer.Close()
	}

	return nil
}

func (s *progressStreamer) Stream(
    buffer [][2]float64,
) (int, bool) {
    n, ok := s.source.Stream(buffer)

    atomic.AddInt64(&s.samples, int64(n))

    return n, ok
}

func (s *progressStreamer) Err() error {
    return s.source.Err()
}