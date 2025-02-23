package player

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func PlaySong(filename string) error {
	filepath := fmt.Sprintf("/home/helix/code/go/firefly/backend/cmd/client/songs/%s", filename)
	f, err := os.Open(filepath)
	if err != nil {
		log.Printf(" Could not open file %s: %v", filepath, err)
		return err
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Printf(" Could not decode file %s: %v", filename, err)
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	log.Printf(" Playing: %s", filename)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		log.Printf("Finished playing: %s", filename)
	})))

	select {}
}
