
package play

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
    "time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

var res = 36000

var targetSampleRate = beep.SampleRate(res)

var (
	speakerOnce    sync.Once
	speakerInitErr error

	playMu  sync.Mutex
	current *beep.Ctrl
)

func decode(path string, f *os.File) (beep.StreamSeekCloser, beep.Format, error) {
    ext := strings.ToLower(filepath.Ext(path))

    switch ext {
    case ".mp3":
        return mp3.Decode(f)

    case ".wav":
        return wav.Decode(f)

    case ".ogg", ".oga":
        return vorbis.Decode(f)

    case ".flac":
        return flac.Decode(f)

    default:
        return nil, beep.Format{}, fmt.Errorf(
            "формат %s не поддерживается библиотекой beep", ext,
        )
    }
}
func openFile(path string) (beep.StreamSeekCloser, beep.Format, *os.File, error,) {
    f, err := os.Open(path)
    if err != nil {
        return nil, beep.Format{}, nil, err
    }

    streamer, format, err := decode(path, f)
    if err != nil {
        f.Close()
        return nil, beep.Format{}, nil, err
    }

    return streamer, format, f, nil
}
func initSpeaker() error {
    speakerOnce.Do(func() {
        speakerInitErr = speaker.Init(
            targetSampleRate,
            targetSampleRate.N(time.Second/10),
        )
    })

    return speakerInitErr
}
func resampleAudio(streamer beep.Streamer, format beep.Format,) beep.Streamer {
    return beep.Resample(
        4,
        format.SampleRate,
        targetSampleRate,
        streamer,
    )
}
func wrapAudio(streamer beep.Streamer, decoder beep.StreamSeekCloser, file *os.File,) beep.Streamer {
    return beep.Seq(
        streamer,
        beep.Callback(func() {
            decoder.Close()
            file.Close()
        }),
    )
}
func startPlayback(streamer beep.Streamer) {
    playMu.Lock()
    defer playMu.Unlock()

    speaker.Lock()

    if current != nil {
        current.Streamer = nil
    }

    current = &beep.Ctrl{
        Streamer: streamer,
    }

    speaker.Unlock()

    speaker.Play(current)
}

func PlayAudio(path string) error {
    streamer, format, file, err := openFile(path)
    if err != nil {
        return fmt.Errorf("ошибка открытия аудио: %w", err)
    }

    if err := initSpeaker(); err != nil {
        streamer.Close()
        file.Close()
        return fmt.Errorf("ошибка speaker: %w", err)
    }

    resampled := resampleAudio(streamer, format)

    wrapped := wrapAudio(resampled, streamer, file)

	startPlayback(wrapped)

    return nil
}

func StopAudio() {
	playMu.Lock()
	defer playMu.Unlock()

	speaker.Lock()

	if current != nil {
		current.Streamer = nil
		current = nil
	}

	speaker.Unlock()
}

func PauseAudio() {
    playMu.Lock()
    defer playMu.Unlock()

    speaker.Lock()
    defer speaker.Unlock()

    if current != nil {
        current.Paused = true
    }
}

func ResumeAudio() {
    playMu.Lock()
    defer playMu.Unlock()

    speaker.Lock()
    defer speaker.Unlock()

    if current != nil {
        current.Paused = false
    }
}