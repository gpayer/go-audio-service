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
