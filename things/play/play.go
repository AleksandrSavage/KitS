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

var res int = 36000

var targetSampleRate = beep.SampleRate(res)

var (
    speakerOnce    sync.Once
    speakerInitErr error
)

func PlayAudio(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("не удалось открыть файл: %w", err)
    }
    defer f.Close()

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
        return fmt.Errorf("формат %s не поддерживается библиотекой beep", ext)
    }

    if err != nil {
        return fmt.Errorf("ошибка декодирования файла: %w", err)
    }
    defer streamer.Close()

    speakerOnce.Do(func() {
        speakerInitErr = speaker.Init(targetSampleRate, targetSampleRate.N(time.Second/10))
    })
    if speakerInitErr != nil {
        return fmt.Errorf("ошибка инициализации аудиовыхода: %w", speakerInitErr)
    }

    resampled := beep.Resample(4, format.SampleRate, targetSampleRate, streamer)

    done := make(chan struct{})
    speaker.Play(beep.Seq(resampled, beep.Callback(func() {
        close(done)
    })))

    <-done
    return nil
}