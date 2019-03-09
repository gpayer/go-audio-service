package notes

type Adsr struct {
	attack  float32 // time in s
	decay   float32 // time in s
	sustain float32 // volume 0.0 .. 1.0
	release float32 // time in s
}

func NewAdsr(attack, decay, sustain, release float32) *Adsr {
	return &Adsr{
		attack:  attack,
		decay:   decay,
		sustain: sustain,
		release: release,
	}
}
