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

type Player struct {
	streamer  beep.StreamCloser
	seeker    beep.StreamSeeker
	format    beep.Format
	ctrl      *beep.Ctrl
	resampler *beep.Resampler
	sound     *effects.Volume
	MinVolume float64
	MaxVolume float64
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
		MaxVolume: 2,
		MinVolume: -10,
	}

	return p, nil
}

func (p *Player) PlayPause() {
	speaker.Lock()
	defer speaker.Unlock()

	p.ctrl.Paused = !p.ctrl.Paused
}

func (p *Player) Paused() bool {
	return p.ctrl.Paused
}

func (p *Player) Forward(seconds int) {
	speaker.Lock()
	defer speaker.Unlock()

	samplesToSkip := p.format.SampleRate.N(time.Second) * seconds
	targetPositionSamples := p.seeker.Position() + samplesToSkip

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
	speaker.Lock()
	defer speaker.Unlock()

	num = math.Abs(num)

	volume := &p.sound.Volume

	if *volume < p.MaxVolume {
		*volume += num
	}
}

func (p *Player) VolumeDown(num float64) {
	speaker.Lock()
	defer speaker.Unlock()

	volume := &p.sound.Volume

	if *volume > p.MinVolume {
		*volume -= math.Abs(num)
	}
}

// The total time length of the media currently playing.
func (p *Player) Total() time.Duration {
	return p.format.SampleRate.D(p.seeker.Len())
}

// The current time position of the media currently playing.
func (p *Player) Current() time.Duration {
	return time.Duration(p.format.SampleRate.D(p.seeker.Position()))
}

// The current volume. 0 = Default
func (p *Player) Volume() float64 {
	return math.Round(p.sound.Volume*10) / 10
}

// Starts audio in speaker.
func (p *Player) Start() {
	speaker.Play(p.resampler)
}

// closes the player streamer.
func (p *Player) Close() {
	p.streamer.Close()
}
