package filters

import (
	"go-audio-service/snd"
	"math"
)

const CompressorMaxDelay int = 1024
const CompressorSpu int = 32
const CompressorSpacingDb float32 = 5.0

type CompressorState struct {
	// metergain can be read by the user after processing a chunk to see how much dB the
	// compressor would have liked to compress the sample; the meter values aren't used to shape the
	// sound in any way, only used for output if desired
	metergain float32

	meterrelease         float32
	threshold            float32
	knee                 float32
	linearpregain        float32
	linearthreshold      float32
	slope                float32
	attacksamplesinv     float32
	satreleasesamplesinv float32
	wet                  float32
	dry                  float32
	k                    float32
	kneedboffset         float32
	linearthresholdknee  float32
	mastergain           float32
	a                    float32 // adaptive release polynomial coefficients
	b                    float32
	c                    float32
	d                    float32
	detectoravg          float32
	compgain             float32
	maxcompdiffdb        float32
	delaybufsize         int
	delaywritepos        int
	delayreadpos         int
	delaybuf             [CompressorMaxDelay]snd.Sample
}

func NewCompressorState() *CompressorState {
	return &CompressorState{}
}

func (state *CompressorState) DefaultCompressor(rate int) {
	state.AdvancedCompressor(rate, 0.0, -24.0, 30.0, 12.0, 0.003, 0.250, 0.006, 0.090, 0.160, 0.420, 0.980, 0.0, 1.0)
}

func (state *CompressorState) SimpleCompressor(rate int, pregain, threshold, knee, ratio, attack, release float32) {
	state.AdvancedCompressor(rate, pregain, threshold, knee, ratio, attack, release, 0.006, 0.090, 0.160, 0.420, 0.980, 0.0, 1.0)
}

func db2lin(db float32) float32 {
	return float32(math.Pow(10.0, float64(0.05*db)))
}

func lin2db(lin float32) float32 {
	return 20.0 * float32(math.Log10(float64(lin)))
}

func kneecurve(x, k, linearthreshold float32) float32 {
	return linearthreshold + (1.0-float32(math.Exp(float64(-k*(x-linearthreshold)))))/k
}

func kneeslope(x, k, linearthreshold float32) float32 {
	return k * x / ((k*linearthreshold + 1.0) * float32(math.Exp(float64((k*(x-linearthreshold))-1))))
}

func compcurve(x, k, slope, linearthreshold, linearthresholdknee, threshold, knee, kneedboffset float32) float32 {
	if x < linearthreshold {
		return x
	}
	if knee <= 0.0 { // no knee in curve
		return db2lin(threshold + slope*(lin2db(x)-threshold))
	}
	if x < linearthresholdknee {
		return kneecurve(x, k, linearthreshold)
	}
	return db2lin(kneedboffset + slope*(lin2db(x)-threshold-knee))
}

func (state *CompressorState) AdvancedCompressor(rate int, pregain, threshold, knee, ratio, attack, release, predelay, releasezone1, releasezone2, releasezone3,
	releasezone4, postgain, wet float32) {
	// setup predelay buffer
	delaybufsize := int(float32(rate) * predelay)
	if delaybufsize < 1 {
		delaybufsize = 1
	} else if delaybufsize > CompressorMaxDelay {
		delaybufsize = CompressorMaxDelay
	}
	for i := 0; i < delaybufsize; i++ {
		state.delaybuf[i].L = 0.0
		state.delaybuf[i].R = 0.0
	}

	// useful values
	linearpregain := db2lin(pregain)
	linearthreshold := db2lin(threshold)
	slope := 1.0 / ratio
	attacksamples := float32(rate) * attack
	attacksamplesinv := 1.0 / attacksamples
	releasesamples := float32(rate) * release
	var satrelease float32 = 0.0025 // seconds
	satreleasesamplesinv := 1.0 / (float32(rate) * satrelease)
	dry := 1.0 - wet

	// metering values (not used in core algorithm, but used to output a meter if desired)
	var meterfalloff float32 = 0.325 // seconds
	var meterrelease float32 = 1.0 - float32(math.Exp(float64(-1.0/(float32(rate)*meterfalloff))))

	// calculate knee curve parameters
	var k float32 = 5.0 // initial guess
	var kneedboffset float32 = 0.0
	var linearthresholdknee float32 = 0.0
	if knee > 0.0 { // if a knee exists, search for a good k value
		xknee := db2lin(threshold + knee)
		var mink float32 = 0.1
		var maxk float32 = 10000.0
		// search by comparing the knee slope at the current k guess, to the ideal slope
		for i := 0; i < 15; i++ {
			if kneeslope(xknee, k, linearthreshold) < slope {
				maxk = k
			} else {
				mink = k
			}
			k = float32(math.Sqrt(float64(mink * maxk)))
		}
		kneedboffset = lin2db(kneecurve(xknee, k, linearthreshold))
		linearthresholdknee = db2lin(threshold + knee)
	}

	// calculate a master gain based on what sounds good
	fulllevel := compcurve(1.0, k, slope, linearthreshold, linearthresholdknee,
		threshold, knee, kneedboffset)
	mastergain := db2lin(postgain) * float32(math.Pow(float64(1.0/fulllevel), 0.6))

	// calculate the adaptive release curve parameters
	// solve a,b,c,d in `y = a*x^3 + b*x^2 + c*x + d`
	// interescting points (0, y1), (1, y2), (2, y3), (3, y4)
	freleasesamples := float32(releasesamples)
	y1 := freleasesamples * releasezone1
	y2 := freleasesamples * releasezone2
	y3 := freleasesamples * releasezone3
	y4 := freleasesamples * releasezone4
	a := (-y1 + 3.0*y2 - 3.0*y3 + y4) / 6.0
	b := y1 - 2.5*y2 + 2.0*y3 - 0.5*y4
	c := (-11.0*y1 + 18.0*y2 - 9.0*y3 + 2.0*y4) / 6.0
	d := y1

	// save everything
	state.metergain = 1.0 // large value overwritten immediately since it's always < 0
	state.meterrelease = meterrelease
	state.threshold = threshold
	state.knee = knee
	state.wet = wet
	state.linearpregain = linearpregain
	state.linearthreshold = linearthreshold
	state.slope = slope
	state.attacksamplesinv = attacksamplesinv
	state.satreleasesamplesinv = satreleasesamplesinv
	state.dry = dry
	state.k = k
	state.kneedboffset = kneedboffset
	state.linearthresholdknee = linearthresholdknee
	state.mastergain = mastergain
	state.a = a
	state.b = b
	state.c = c
	state.d = d
	state.detectoravg = 0.0
	state.compgain = 1.0
	state.maxcompdiffdb = -1.0
	state.delaybufsize = delaybufsize
	state.delaywritepos = 0
	if delaybufsize > 1 {
		state.delayreadpos = 1
	} else {
		state.delayreadpos = 0
	}
}

func adaptivereleasecurve(x, a, b, c, d float32) float32 {
	// a*x^3 + b*x^2 + c*x + d
	x2 := x * x
	return a*x2*x + b*x2 + c*x + d
}

func clampf(v, min, max float32) float32 {
	if v < min {
		return min
	} else if v > max {
		return max
	}
	return v
}

func absf(v float32) float32 {
	if v < 0.0 {
		return -v
	}
	return v
}

func fixf(v, def float32) float32 {
	// fix NaN and infinity values that sneak in... not sure why this is needed, but it is
	vd := float64(v)
	if math.IsNaN(vd) || math.IsInf(vd, 1) || math.IsInf(vd, -1) {
		return def
	}
	return v
}

func (state *CompressorState) Process(size int, input, output []snd.Sample) {
	metergain := state.metergain
	meterrelease := state.meterrelease
	threshold := state.threshold
	knee := state.knee
	linearpregain := state.linearpregain
	linearthreshold := state.linearthreshold
	slope := state.slope
	attacksamplesinv := state.attacksamplesinv
	satreleasesamplesinv := state.satreleasesamplesinv
	wet := state.wet
	dry := state.dry
	k := state.k
	kneedboffset := state.kneedboffset
	linearthresholdknee := state.linearthresholdknee
	mastergain := state.mastergain
	a := state.a
	b := state.b
	c := state.c
	d := state.d
	detectoravg := state.detectoravg
	compgain := state.compgain
	maxcompdiffdb := state.maxcompdiffdb
	delaybufsize := state.delaybufsize
	delaywritepos := state.delaywritepos
	delayreadpos := state.delayreadpos
	//sf_sample_st *delaybuf     = state->delaybuf;

	samplesperchunk := CompressorSpu
	chunks := size / samplesperchunk
	var ang90 float32 = math.Pi * 0.5
	var ang90inv float32 = 2.0 / math.Pi
	samplepos := 0
	spacingdb := CompressorSpacingDb

	for ch := 0; ch < chunks; ch++ {
		detectoravg = fixf(detectoravg, 1.0)
		desiredgain := detectoravg
		scaleddesiredgain := float32(math.Asin(float64(desiredgain))) * ang90inv
		compdiffdb := lin2db(compgain / scaleddesiredgain)

		// calculate envelope rate based on whether we're attacking or releasing
		var enveloperate float32
		if compdiffdb < 0.0 { // compgain < scaleddesiredgain, so we're releasing
			compdiffdb = fixf(compdiffdb, -1.0)
			maxcompdiffdb = -1 // reset for a future attack mode
			// apply the adaptive release curve
			// scale compdiffdb between 0-3
			x := (clampf(compdiffdb, -12.0, 0.0) + 12.0) * 0.25
			releasesamples := adaptivereleasecurve(x, a, b, c, d)
			enveloperate = db2lin(spacingdb / releasesamples)
		} else { // compresorgain > scaleddesiredgain, so we're attacking
			compdiffdb = fixf(compdiffdb, 1.0)
			if maxcompdiffdb == -1.0 || maxcompdiffdb < compdiffdb {
				maxcompdiffdb = compdiffdb
			}
			attenuate := maxcompdiffdb
			if attenuate < 0.5 {
				attenuate = 0.5
			}
			enveloperate = 1.0 - float32(math.Pow(float64(0.25/attenuate), float64(attacksamplesinv)))
		}

		// process the chunk
		for chi := 0; chi < samplesperchunk; chi++ {

			inputL := input[samplepos].L * linearpregain
			inputR := input[samplepos].R * linearpregain
			state.delaybuf[delaywritepos] = snd.Sample{L: inputL, R: inputR}

			inputL = absf(inputL)
			inputR = absf(inputR)
			var inputmax float32
			if inputL > inputR {
				inputmax = inputL
			} else {
				inputmax = inputR
			}

			var attenuation float32
			if inputmax < 0.0001 {
				attenuation = 1.0
			} else {
				inputcomp := compcurve(inputmax, k, slope, linearthreshold,
					linearthresholdknee, threshold, knee, kneedboffset)
				attenuation = inputcomp / inputmax
			}

			var rate float32
			if attenuation > detectoravg { // if releasing
				attenuationdb := -lin2db(attenuation)
				if attenuationdb < 2.0 {
					attenuationdb = 2.0
				}
				dbpersample := attenuationdb * satreleasesamplesinv
				rate = db2lin(dbpersample) - 1.0
			} else {
				rate = 1.0
			}

			detectoravg += (attenuation - detectoravg) * rate
			if detectoravg > 1.0 {
				detectoravg = 1.0
			}
			detectoravg = fixf(detectoravg, 1.0)

			if enveloperate < 1 { // attack, reduce gain
				compgain += (scaleddesiredgain - compgain) * enveloperate
			} else { // release, increase gain
				compgain *= enveloperate
				if compgain > 1.0 {
					compgain = 1.0
				}
			}

			// the final gain value!
			premixgain := float32(math.Sin(float64(ang90 * compgain)))
			gain := dry + wet*mastergain*premixgain

			// calculate metering (not used in core algo, but used to output a meter if desired)
			premixgaindb := lin2db(premixgain)
			if premixgaindb < metergain {
				metergain = premixgaindb // spike immediately
			} else {
				metergain += (premixgaindb - metergain) * meterrelease // fall slowly
			}

			// apply the gain
			output[samplepos] = snd.Sample{
				L: state.delaybuf[delayreadpos].L * gain,
				R: state.delaybuf[delayreadpos].R * gain,
			}

			samplepos++
			delayreadpos = (delayreadpos + 1) % delaybufsize
			delaywritepos = (delaywritepos + 1) % delaybufsize
		}
	}

	state.metergain = metergain
	state.detectoravg = detectoravg
	state.compgain = compgain
	state.maxcompdiffdb = maxcompdiffdb
	state.delaywritepos = delaywritepos
	state.delayreadpos = delayreadpos
}

type Compressor struct {
	state    *CompressorState
	buffer   []snd.Sample
	readable snd.Readable
	input    *snd.Samples
	output   []snd.Sample
}

func NewCompressor(samplerate uint32, state *CompressorState) *Compressor {
	return &Compressor{
		state:  state,
		input:  snd.NewSamples(samplerate, 128),
		output: make([]snd.Sample, 128),
	}
}

func (comp *Compressor) SetReadable(r snd.Readable) {
	comp.readable = r
}

func (comp *Compressor) Read(samples *snd.Samples) {
	comp.ReadStateless(samples, 0, 0, true)
}

func (comp *Compressor) ReadStateless(samples *snd.Samples, freq float32, timecode uint32, on bool) {
	length := len(samples.Frames)
	for len(comp.buffer) < length {
		comp.readable.ReadStateless(comp.input, freq, timecode+uint32(len(comp.buffer)), on)
		comp.state.Process(128, comp.input.Frames, comp.output)
		comp.buffer = append(comp.buffer, comp.output...)
	}
	length = copy(samples.Frames, comp.buffer[:length])
	comp.buffer = comp.buffer[length:]
}
