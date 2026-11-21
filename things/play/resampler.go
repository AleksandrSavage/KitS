package play

import (
    "fmt"

    "github.com/gopxl/beep/v2"
)

const resampleQuality = 4

func resampleAudio(
    streamer beep.Streamer,
    sourceRate beep.SampleRate,
    targetRate beep.SampleRate,
) (beep.Streamer, error) {

    if streamer == nil {
        return nil, fmt.Errorf("audio streamer is nil")
    }

    if sourceRate <= 0 || targetRate <= 0 {
        return nil, fmt.Errorf(
            "invalid sample rate: source=%d, target=%d",
            sourceRate,
            targetRate,
        )
    }

    if sourceRate == targetRate {
        return streamer, nil
    }

    return beep.Resample(
        resampleQuality,
        sourceRate,
        targetRate,
        streamer,
    ), nil
}
