package snd

// Sample contains one sound sample as 32bit floats (range -1.0 to 1.0)
type Sample struct {
	L float32
	R float32
}

// Samples is a list of samples and their sample rate
type Samples struct {
	Frames     []Sample
	SampleRate uint32
}

// Add adds one or more Samples to a Samples struct
func (s *Samples) Add(samples ...Sample) {
	s.Frames = append(s.Frames, samples...)
}

func NewSamples(samplerate uint32, length int) *Samples {
	return &Samples{
		SampleRate: samplerate,
		Frames:     make([]Sample, length),
	}
}
