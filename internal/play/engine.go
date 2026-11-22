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

	tracked := &progressStreamer{
    source: playback,
	}

	ctrl := &beep.Ctrl{
    	Streamer: tracked,
	}

	e.current = audio
	e.ctrl = ctrl

	e.progress = tracked
	e.duration = audio.format.SampleRate.D(
    	audio.streamer.Len(),
	)
	e.ended = false
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

	e.progress = nil
	e.duration = 0
	e.ended = false

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

func (e *engine) getProgress() (time.Duration, time.Duration) {
    e.mu.Lock()
    defer e.mu.Unlock()

    if e.progress == nil {
        return 0, 0
    }

    duration := e.duration

    samples := atomic.LoadInt64(&e.progress.samples)
    position := targetSampleRate.D(int(samples))

    if e.ended {
        position = duration
    }

    if position > duration {
        position = duration
    }

    return position, duration
}