
package play

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/gopxl/beep/v2"
    "github.com/gopxl/beep/v2/flac"
    "github.com/gopxl/beep/v2/mp3"
    "github.com/gopxl/beep/v2/vorbis"
    "github.com/gopxl/beep/v2/wav"
)

type decodedAudio struct {
    streamer beep.StreamSeekCloser
    format   beep.Format
}

func openAudio(path string) (*decodedAudio, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open audio file: %w", err)
    }

    ext := strings.ToLower(filepath.Ext(path))

    streamer, format, err := decodeFile(ext, file)
    if err != nil {
        _ = file.Close()
        return nil, fmt.Errorf("decode audio file: %w", err)
    }

    return &decodedAudio{
        streamer: streamer,
        format:   format,
    }, nil
}

func decodeFile(
    ext string,
    file *os.File,
) (beep.StreamSeekCloser, beep.Format, error) {

    switch ext {
    case ".mp3":
        return mp3.Decode(file)

    case ".wav":
        return wav.Decode(file)

    case ".ogg", ".oga":
        return vorbis.Decode(file)

    case ".flac":
        return flac.Decode(file)

    default:
        return nil, beep.Format{}, fmt.Errorf(
            "unsupported audio format: %s", ext,
        )
    }
}
