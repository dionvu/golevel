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
	minVolume = -2.0
	maxVolume = 2.0
)

type Player struct {
	streamer  beep.StreamCloser
	seeker    beep.StreamSeeker
	format    beep.Format
	ctrl      *beep.Ctrl
	resampler *beep.Resampler
	sound     *effects.Volume
}

func New(file *os.File) (*Player, error) {
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}

	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	speed := beep.ResampleRatio(4, 1, volume)

	p := &Player{
		streamer:  streamer,
		seeker:    streamer,
		ctrl:      ctrl,
		format:    format,
		resampler: speed,
		sound:     volume,
	}

	return p, nil
}

func (p *Player) Forward(seconds int) {
	speaker.Lock()
	defer speaker.Unlock()

	samplesToSkip := p.format.SampleRate.N(time.Second) * seconds
	targetPositionSamples := p.seeker.Position() + samplesToSkip

	if targetPositionSamples > p.seeker.Len() {
		targetPositionSamples = p.seeker.Len()
	}

	p.seeker.Seek(targetPositionSamples)
}

func (p *Player) Backward(seconds int) {
	speaker.Lock()
	defer speaker.Unlock()

	samplesToSkip := p.format.SampleRate.N(time.Second) * seconds
	targetPositionSamples := p.seeker.Position() - samplesToSkip

	if targetPositionSamples < 0 {
		targetPositionSamples = 0
	}

	p.seeker.Seek(targetPositionSamples)
}

func (p *Player) VolumeUp(num float64) {
	num = math.Abs(num)

	volume := &p.sound.Volume

	if *volume >= maxVolume {
		return
	}

	speaker.Lock()
	*volume += num
	speaker.Unlock()
}

func (p *Player) Total() time.Duration {
	return p.format.SampleRate.D(p.seeker.Len())
}

func (p *Player) VolumeDown(num float64) {
	num = -math.Abs(num)

	volume := &p.sound.Volume

	if *volume <= minVolume {
		return
	}

	speaker.Lock()
	*volume += num
	speaker.Unlock()
}

func (p *Player) Volume() float64 {
	return p.sound.Volume
}

func (p *Player) Start() {
	speaker.Play(p.resampler)
}

func (p *Player) Close() {
	p.streamer.Close()
}
