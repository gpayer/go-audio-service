package notes

import (
	"go-audio-service/snd"
)

type Adsr struct {
	attack          float32 // time in s
	decay           float32 // time in s
	sustain         float32 // volume 0.0 .. 1.0
	release         float32 // time in s
	readable        snd.Readable
	samplerate      uint32
	d_attack        float32 // delta volume in attack phase
	t_decay         uint32  // timecode when decay starts
	d_decay         float32 // delta volume in decay phase
	t_sustain       uint32  // timecode when sustain starts
	d_release       float32 // delta volume in release phase
	t_end           uint32
	releaseGain     float32 // current gain at point of release
	releaseTimecode uint32
	ended           bool
	subnotestate    *snd.NoteState
}

func NewAdsr(attack, decay, sustain, release float32) *Adsr {
	return &Adsr{
		attack:       attack,
		decay:        decay,
		sustain:      sustain,
		release:      release,
		subnotestate: &snd.NoteState{On: true},
	}
}

func (adsr *Adsr) calcParameters() {
	ft_decay := float32(adsr.samplerate) * adsr.attack
	adsr.t_decay = uint32(ft_decay)
	adsr.d_attack = 1.0 / ft_decay
	ft_sustain := float32(adsr.samplerate) * (adsr.attack + adsr.decay)
	adsr.t_sustain = uint32(ft_sustain)
	adsr.d_decay = (adsr.sustain - 1.0) / (float32(adsr.samplerate) * adsr.decay)
}

func (adsr *Adsr) calcEnd() {
	if adsr.sustain == 0.0 {
		adsr.t_end = adsr.t_sustain
	} else {
		adsr.t_end = 300 * adsr.samplerate
	}
}

func (adsr *Adsr) calcRelease(timecode uint32) {
	if adsr.sustain == 0.0 {
		return
	}
	var currentGain float32
	if timecode > adsr.t_sustain {
		currentGain = adsr.sustain
	} else if timecode > adsr.t_decay {
		currentGain = 1.0 + float32(timecode-adsr.t_decay)*adsr.d_decay
	} else {
		currentGain = float32(timecode) * adsr.d_attack
	}
	adsr.releaseGain = currentGain
	adsr.releaseTimecode = timecode
	ft_decay := float32(adsr.samplerate) * adsr.release
	adsr.d_release = -currentGain / ft_decay
	adsr.t_end = adsr.releaseTimecode + uint32(ft_decay)
}

func (adsr *Adsr) SetReadable(r snd.Readable) {
	adsr.readable = r
}

func (adsr *Adsr) Read(samples *snd.Samples) {
	adsr.readable.Read(samples)
}

func (adsr *Adsr) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	if samples.SampleRate != adsr.samplerate {
		adsr.samplerate = samples.SampleRate
		adsr.calcParameters()
	}
	if !state.On {
		adsr.calcRelease(state.ReleaseTimecode)
	} else {
		adsr.calcEnd()
	}

	adsr.ended = state.Timecode > adsr.t_end

	adsr.subnotestate.Timecode = state.Timecode
	adsr.subnotestate.ReleaseTimecode = state.ReleaseTimecode
	adsr.subnotestate.Volume = state.Volume
	adsr.readable.ReadStateless(samples, freq, adsr.subnotestate)

	for i := 0; i < len(samples.Frames); i++ {
		currentTimecode := state.Timecode + uint32(i)
		var gain float32
		if state.On || adsr.sustain == 0.0 {
			if currentTimecode > adsr.t_sustain && adsr.sustain > 0 {
				gain = adsr.sustain
			} else if currentTimecode > adsr.t_decay {
				gain = 1.0 + float32(currentTimecode-adsr.t_decay)*adsr.d_decay
			} else {
				gain = float32(currentTimecode) * adsr.d_attack
			}
		} else {
			gain = adsr.releaseGain + float32(currentTimecode-adsr.releaseTimecode)*adsr.d_release
		}
		if gain < 0.0 {
			gain = 0.0
		}

		samples.Frames[i].L *= gain
		samples.Frames[i].R *= gain
	}
}

func (adsr *Adsr) NoteEnded() bool {
	return adsr.ended
}
