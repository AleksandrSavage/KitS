
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

func PlayAudio(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}

	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
	)

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg", ".oga":
		streamer, format, err = vorbis.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	default:
		f.Close()
		return fmt.Errorf("формат %s не поддерживается", ext)
	}

	if err != nil {
		f.Close()
		return fmt.Errorf("ошибка декодирования: %w", err)
	}

	speakerOnce.Do(func() {
		speakerInitErr = speaker.Init(targetSampleRate, targetSampleRate.N(time.Second/10))
	})

	if speakerInitErr != nil {
		streamer.Close()
		f.Close()
		return fmt.Errorf("ошибка инициализации speaker: %w", speakerInitErr)
	}

	
	resampled := beep.Resample(
		4,
		format.SampleRate,
		targetSampleRate,
		streamer,
	)

	
	wrapped := beep.Seq(
		resampled,
		beep.Callback(func() {
			streamer.Close()
			f.Close()
		}),
	)

	playMu.Lock()
	defer playMu.Unlock()

	speaker.Lock()


	if current != nil {
		current.Streamer = nil
	}


	current = &beep.Ctrl{
		Streamer: wrapped,
	}

	speaker.Unlock()

	speaker.Play(current)

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
