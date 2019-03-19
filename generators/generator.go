package generators

import (
	"go-audio-service/snd"
)

type Generator interface {
	snd.Readable
	snd.WritableProvider
}
