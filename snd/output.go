package snd

import (
	"encoding/binary"
	"fmt"

	"github.com/gen2brain/malgo"
)

// Output object for sound output
type Output struct {
	context    *malgo.AllocatedContext
	device     *malgo.Device
	samplerate uint32
	samplesize int
	bufferChan chan Sample
}

// NewOutput creates a new Output instance
func NewOutput(samplerate uint32, buffersize int) (*Output, error) {
	o := &Output{}
	var err error
	o.context, err = malgo.InitContext(nil, malgo.ContextConfig{}, func(_ string) {})
	if err != nil {
		return nil, err
	}

	deviceConfig := malgo.DefaultDeviceConfig()
	deviceConfig.Format = malgo.FormatS16
	deviceConfig.Channels = 2
	deviceConfig.SampleRate = samplerate

	o.samplerate = samplerate
	o.samplesize = malgo.SampleSizeInBytes(deviceConfig.Format)
	o.bufferChan = make(chan Sample, buffersize)

	onSendSamples := func(requestedSampleCount uint32, samples []byte) uint32 {
		// fmt.Println(requestedSampleCount, len(samples))
		var readCount uint32
		offset := 0
		for readCount < requestedSampleCount {
			select {
			case sample := <-o.bufferChan:
				l := floatToBytes(sample.L)
				samples[offset] = l[0]
				samples[offset+1] = l[1]
				r := floatToBytes(sample.R)
				samples[offset+2] = r[0]
				samples[offset+3] = r[1]
				offset += 4
				readCount++
			default:
				return readCount
			}
		}
		return readCount
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Send: malgo.SendProc(onSendSamples),
	}

	o.device, err = malgo.InitDevice(o.context.Context, malgo.Playback, nil, deviceConfig, deviceCallbacks)
	if err != nil {
		_ = o.context.Uninit()
		o.context.Free()
		return nil, err
	}
	o.samplerate = o.device.SampleRate()

	return o, nil
}

// Close closes the output
func (o *Output) Close() {
	o.device.Uninit()
	_ = o.context.Uninit()
	o.context.Free()
}

// Start starts playback on the output
func (o *Output) Start() (err error) {
	if !o.device.IsStarted() {
		err = o.device.Start()
	}
	return
}

// Stop stops playback
func (o *Output) Stop() (err error) {
	if o.device.IsStarted() {
		err = o.device.Stop()
	}
	return
}

// Write writes all given samples to the ouput buffer
func (o *Output) Write(samples *Samples) error {
	if samples.SampleRate != o.samplerate {
		return fmt.Errorf("wrong samplerate, device has %d, %d given", samples.SampleRate, o.samplerate)
	}
	for _, sample := range samples.Frames {
		o.bufferChan <- sample
	}
	return nil
}

func (o *Output) SetOutput(_ Filter) {}

func floatToBytes(f float32) []byte {
	var i int16
	if f > 0.0 {
		i = int16(f * 32767.0)
	} else {
		i = int16(f * 32768.0)
	}
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(i))
	return bs
}
