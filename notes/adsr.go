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
	on              bool
	releaseGain     float32 // current gain at point of release
	releaseTimecode uint32
}

func NewAdsr(attack, decay, sustain, release float32) *Adsr {
	return &Adsr{
		attack:  attack,
		decay:   decay,
		sustain: sustain,
		release: release,
		on:      false,
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

func (adsr *Adsr) calcRelease(timecode uint32) {
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
	adsr.d_release = -currentGain / (float32(adsr.samplerate) * adsr.release)
}

func (adsr *Adsr) SetReadable(r snd.Readable) {
	adsr.readable = r
}

func (adsr *Adsr) Read(samples *snd.Samples) {
	adsr.readable.Read(samples)
}

func (adsr *Adsr) ReadStateless(samples *snd.Samples, freq float32, timecode uint32, on bool) {
	if samples.SampleRate != adsr.samplerate {
		adsr.samplerate = samples.SampleRate
		adsr.calcParameters()
	}
	if adsr.on && !on {
		adsr.calcRelease(timecode)
	}
	adsr.on = on

	adsr.readable.ReadStateless(samples, freq, timecode, on)

	for i := 0; i < len(samples.Frames); i++ {
		currentTimecode := timecode + uint32(i)
		var gain float32
		if currentTimecode > adsr.t_sustain {
			gain = adsr.sustain
		} else if currentTimecode > adsr.t_decay {
			gain = 1.0 + float32(currentTimecode-adsr.t_decay)*adsr.d_decay
		} else {
			gain = float32(currentTimecode) * adsr.d_attack
		}

		samples.Frames[i].L *= gain
		samples.Frames[i].R *= gain
	}
}
