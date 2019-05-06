# go-audio-service

**Important note: this is under heavy development and anything might change anytime!**

*go-audio-service* aims to be a flexible audio synthesizer framework. At the moment the following features are either implemented or planned:
* Generators
  * Rectangle generator
  * Sine generator
  * Sound samples from standard file formats: MP3, WAV
  * *Planned: sawtooth and triangle generators*
* Filters
  * Compressor
  * Lowpass
  * *Planned: highpass, bandpass, reverb etc.*
* Modulation capabilities
  * Generators have inputs for AM and FM modulation
* Pluggable architecture
  * Components can be connected to other components to create complex instruments, e.g. connect a sine generator to the FM modulation input of another sine generator
* ADSR Envelope
* Multi-channel mixer
* Output sink provided by `gen2brain/malgo`
* *Planned: LFO as filter*
* *Planned: high level sound and music player*

The inspiration for this library is the CSound library. However instead of writing instruments and notes in a custom language I wanted this to be used programmatically. This way it can be utilized more flexibly.

## Requirements

(TODO: test installation procedure and create go1.12 dependencies file)

Install the following package:
> go get -u github.com/gen2brain/malgo

For the main example you need additional packages:
* `libportmidi` and `rakyll/portmidi`
* `gpayer/pixelext` (depends on `faiface/pixel`, see https://github.com/faiface/pixel#requirements for additional requirements)

#### Ubuntu
> sudo apt install libportmidi0
#### OpenSUSE
> zypper in libportmidi0

Finally install the golang packages:

> go get -u github.com/rakyll/portmidi\
> go get -u github.com/faiface/pixel\
> go get -u github.com/gpayer/pixelext

## Examples
Go to into directory `examples` and run
> go run examples.go

Not all examples might work at the moment, especially old singular examples in subdirectories.
In the main `examples` executable select an example by pressing its number. 

The most fleshed out example is the *Keyboard* example. It lets you play with two connected oscillators. Input is read from the keyboard (keys *A* to *K*) or from a connected MIDI input. The active MIDI input can be selected from a dropdown in the top-right corner.

## Attributions
Code in package _filters_ is translated from C to Go based on [github.com/voidqk/sndfilter](https://github.com/voidqk/sndfilter) by Sean Connelly (@voidqk).