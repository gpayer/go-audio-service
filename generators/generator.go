package generators

import (
	"go-audio-service/snd"
)

type FreqModable interface {
	SetFreqMod(v float32)
}

type Generator interface {
	snd.Readable
	snd.WritableProvider
}
