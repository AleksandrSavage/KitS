
package main

import (
	"log"
	"os"
	"time"


	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

func main() {
	f, err := os.Open("./audio-test.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
	defer streamer.Close()

	err = speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)
	if err != nil {
		log.Fatal(err)
	}

	speaker.PlayAndWait(streamer)
}
