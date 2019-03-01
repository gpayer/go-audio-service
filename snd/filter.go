package snd

// Filter defines the interface for filtering components
type Filter interface {
	Write(samples *Samples) error
	SetOutput(out Filter)
}
