package player

import (
	"math"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

const (
	minVol = -2.0
	maxVol = 2.0
)

type Player struct {
	Streamer  beep.StreamCloser
	format    beep.Format
	Ctrl      *beep.Ctrl
	Resampler *beep.Resampler
	sound     *effects.Volume
}

func New(file *os.File) (*Player, error) {
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}

	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	speed := beep.ResampleRatio(4, 1, volume)

	return &Player{
		Streamer:  streamer,
		Ctrl:      ctrl,
		Resampler: speed,
		sound:     volume,
	}, nil
}

func (p *Player) VolumeUp(num float64) {
	num = math.Abs(num)

	vol := &p.sound.Volume

	if *vol >= maxVol {
		return
	}

	speaker.Lock()
	*vol += num
	speaker.Unlock()
}

func (p *Player) Volume() float64 {
	return p.sound.Volume
}
