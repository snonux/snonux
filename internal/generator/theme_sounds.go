package generator

import (
	"encoding/json"
	"html/template"
	"math"
)

// melodyNote is a single note in a looping ambient melody.
type melodyNote struct {
	Freq float64 `json:"freq"`
	Dur  float64 `json:"dur"`
	Step float64 `json:"step,omitempty"`
	Gain float64 `json:"gain,omitempty"`
}

// ambientPreset describes a generative ambient background layer for a theme.
type ambientPreset struct {
	BPM           float64      `json:"bpm,omitempty"`
	PulseInterval float64      `json:"pulseInterval,omitempty"`
	Gain          float64      `json:"gain,omitempty"`
	Wave          string       `json:"wave,omitempty"`
	DroneFreqs    []float64    `json:"droneFreqs,omitempty"`
	PulseFreqs    []float64    `json:"pulseFreqs,omitempty"`
	CutoffMin     float64      `json:"cutoffMin,omitempty"`
	CutoffMax     float64      `json:"cutoffMax,omitempty"`
	NoiseGain     float64      `json:"noiseGain,omitempty"`
	Attack        float64      `json:"attack,omitempty"`
	Release       float64      `json:"release,omitempty"`
	DetuneCents   float64      `json:"detuneCents,omitempty"`
	Rhythm        []float64    `json:"rhythm,omitempty"`
	Melody        []melodyNote `json:"melody,omitempty"`
}

// ambientSounds holds normal and wild-mode ambient presets for a theme.
type ambientSounds struct {
	Normal ambientPreset `json:"normal,omitempty"`
	Wild   ambientPreset `json:"wild,omitempty"`
}

// themeSounds is serialized into each page for Web Audio (splash + keyboard nav).
// Wave: "sine" | "triangle" | "square".
type themeSounds struct {
	Splash struct {
		Freqs   []float64 `json:"freqs"`
		Spacing float64   `json:"spacing"`
		Gain    float64   `json:"gain"`
		Wave    string    `json:"wave"`
	} `json:"splash"`
	Nav struct {
		Freq float64 `json:"freq"`
		Wave string  `json:"wave"`
		Dur  float64 `json:"dur"`
		Gain float64 `json:"gain"`
	} `json:"nav"`
	Open struct {
		Wave  string  `json:"wave"`
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Dur   float64 `json:"dur"`
		Gain  float64 `json:"gain"`
	} `json:"open"`
	Close struct {
		Wave  string  `json:"wave"`
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Dur   float64 `json:"dur"`
		Gain  float64 `json:"gain"`
	} `json:"close"`
	Bounce struct {
		Wave  string  `json:"wave"`
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Dur   float64 `json:"dur"`
		Gain  float64 `json:"gain"`
	} `json:"bounce"`
	Ambient ambientSounds `json:"ambient,omitempty"`
}

// ── helpers ────────────────────────────────────────────────────────

func n(freq, dur float64) melodyNote {
	return melodyNote{Freq: freq, Dur: dur}
}

// ns builds a note that rings for 'dur' but advances the sequencer by 'step'.
// When dur > step, the note overlaps the next one — creating chords on a
// monophonic step sequencer.
func ns(freq, dur, step float64) melodyNote {
	return melodyNote{Freq: freq, Dur: dur, Step: step}
}

// equal-temperament intervals
var (
	intMajor3rd = math.Pow(2, 4.0/12)
	intMinor3rd = math.Pow(2, 3.0/12)
	intPerf5th  = math.Pow(2, 7.0/12)
)

func major(freq float64) [3]float64 {
	return [3]float64{freq, freq * intMajor3rd, freq * intPerf5th}
}
func minor(freq float64) [3]float64 {
	return [3]float64{freq, freq * intMinor3rd, freq * intPerf5th}
}

// mNote returns one note in a measure.
func mNote(freq, dur, step float64) melodyNote {
	return ns(freq, dur, step)
}

// chordArp returns an up/down arpeggio of a triad; each note sustains for 'dur'
// while the sequencer only advances by 'step', so multiple voices overlap.
func chordArp(chord [3]float64, dur, step float64) []melodyNote {
	return []melodyNote{
		ns(chord[0], dur, step),
		ns(chord[1], dur, step),
		ns(chord[2], dur, step),
		ns(chord[1], dur, step),
	}
}

// concatNotes flattens many melody slices into one.
func concatNotes(groups ...[]melodyNote) []melodyNote {
	var out []melodyNote
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

// ── song builder ───────────────────────────────────────────────────

// buildMeasure creates one measure: bass + chord pad + melody top-line.
// bassFreq: audible bass root (130–260 Hz).
// chord:    triad in middle register (260–520 Hz).
// melody:   2–3 scalar notes in upper register (400–1000 Hz).
// Returns notes that sum to ~4.0 s total when stepNorm=0.5s.
func buildMeasure(bassFreq float64, chord [3]float64, melody []float64, step float64) []melodyNote {
	// Bass sustains for 7.5×step, giving long drone during the measure.
	bassDur := step * 7.5
	chordDur := step * 3.0
	melDur := step * 1.6

	// If bass is too low, shift it up one octave.
	if bassFreq < 120 {
		bassFreq *= 2
	}

	// Ensure chord is in audible middle register (260–520 Hz).
	if chord[0] < 200 {
		chord = [3]float64{chord[0] * 4, chord[1] * 4, chord[2] * 4}
	} else if chord[0] < 260 {
		chord = [3]float64{chord[0] * 2, chord[1] * 2, chord[2] * 2}
	}

	var notes []melodyNote
	notes = append(notes, ns(bassFreq, bassDur, step))
	notes = append(notes, chordArp(chord, chordDur, step)...)

	for _, f := range melody {
		ff := f
		if ff < 300 {
			ff *= 4
		} else if ff < 400 {
			ff *= 2
		}
		notes = append(notes, ns(ff, melDur, step))
	}
	return notes
}

// buildSong creates a 4-measure loop (~16 s normal, ~8 s wild) from a chord
// progression and matching melody lines.
func buildSong(bassFreqs [4]float64, chords [4][3]float64, melodies [4][]float64, step float64) []melodyNote {
	var out []melodyNote
	for i := 0; i < 4; i++ {
		out = append(out, buildMeasure(bassFreqs[i], chords[i], melodies[i], step)...)
	}
	return out
}

// ── theme registry ─────────────────────────────────────────────────

var themeSoundPresets = map[string]themeSounds{
	"neon":         soundsNeon(),
	"terminal":     soundsTerminal(),
	"synthwave":    soundsSynthwave(),
	"plasma":       soundsPlasma(),
	"brutalist":    soundsBrutalist(),
	"volcano":      soundsVolcano(),
	"aurora":       soundsAurora(),
	"matrix":       soundsMatrix(),
	"ocean":        soundsOcean(),
	"dos":          soundsDos(),
	"retro":        soundsRetro(),
	"cosmos":       soundsCosmos(),
	"retrofuture":  soundsRetrofuture(),
	"spaceage":     soundsSpaceage(),
	"tropicale":    soundsTropical(),
	"noir":         soundsNoir(),
	"cathedral":    soundsCathedral(),
	"surveillance": soundsSurveillance(),
	"biomech":      soundsBiomech(),
}

// ════════════════════════════════════════════════════════════════════
// NEON  –  bright synth-pop in C major (C → F → G → C)
// ════════════════════════════════════════════════════════════════════
func soundsNeon() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 330, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 523.25, 1046.5, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 880, 261.63, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.12, 0.1

	cm := major(65.41)
	fm := major(87.31)
	gm := major(98.00)

	bass := [4]float64{130.81, 174.61, 196.00, 130.81} // C3 F3 G3 C3
	chords := [4][3]float64{cm, fm, gm, cm}
	normalMel := [4][]float64{
		{523.25, 587.33, 659.25, 523.25},       // C D E C
		{698.46, 783.99, 880.00, 698.46},       // F G A F
		{783.99, 880.00, 1046.5, 783.99},      // G A C G
		{523.25, 659.25, 783.99, 1046.5},      // C E G C
	}
	wildMel := [4][]float64{
		{1046.5, 1174.66, 1318.5, 1046.5},
		{1396.92, 1567.98, 1760.00, 1396.92},
		{1567.98, 1760.00, 2093.0, 1567.98},
		{1046.5, 1318.5, 1567.98, 2093.0},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{130.81, 174.61, 196.00, 261.63},
		Attack: 0.6, Release: 1.5, CutoffMin: 800, CutoffMax: 3000,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{261.63, 349.23, 392.00, 523.25},
		Attack: 0.2, Release: 0.6, CutoffMin: 1500, CutoffMax: 6000, DetuneCents: 8,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TERMINAL  –  dark industrial, E minor (Em → Bm → Am → Em)
// ════════════════════════════════════════════════════════════════════
func soundsTerminal() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.11, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 800, "square", 0.045, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 600, 1200, 0.12, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 900, 400, 0.14, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1

	em := minor(82.41)
	bm := minor(61.74)
	am := minor(110.00)

	bass := [4]float64{164.81, 123.47, 220.00, 164.81} // E3 B2 A3 E3
	chords := [4][3]float64{em, bm, am, em}
	normalMel := [4][]float64{
		{329.63, 392.00, 440.00, 329.63},       // E3=164, so E4=329
		{246.94, 293.66, 329.63, 246.94},       // B2=123, B3=246
		{440.00, 523.25, 587.33, 440.00},       // A3=220
		{329.63, 392.00, 493.88, 329.63},       // E4=E4
	}
	wildMel := [4][]float64{
		{659.25, 783.99, 880.00, 659.25},
		{493.88, 587.33, 659.25, 493.88},
		{880.00, 1046.5, 1174.66, 880.00},
		{659.25, 783.99, 987.77, 659.25},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "square",
		DroneFreqs: []float64{82.41, 123.47, 164.81},
		Attack: 0.1, Release: 0.3,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 140, Wave: "square",
		DroneFreqs: []float64{164.81, 246.94, 329.63},
		Attack: 0.05, Release: 0.15,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SYNTHWAVE  –  retro-future, A minor (Am → F → C → G)
// ════════════════════════════════════════════════════════════════════
func soundsSynthwave() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 440, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 110, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 150, 75, 0.14, 0.09

	am := minor(55.00)
	fm := major(87.31)
	cm := major(65.41)
	gm := major(98.00)

	bass := [4]float64{110.00, 174.61, 130.81, 196.00} // A2 F3 C3 G3
	chords := [4][3]float64{am, fm, cm, gm}
	normalMel := [4][]float64{
		{220.00, 261.63, 293.66, 220.00},
		{349.23, 392.00, 440.00, 349.23},
		{261.63, 329.63, 392.00, 261.63},
		{392.00, 440.00, 493.88, 392.00},
	}
	wildMel := [4][]float64{
		{440.00, 523.25, 587.33, 440.00},
		{698.46, 783.99, 880.00, 698.46},
		{523.25, 659.25, 783.99, 523.25},
		{783.99, 880.00, 987.77, 783.99},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{110.00, 130.81, 174.61, 196.00},
		Attack: 0.8, Release: 1.5,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 130, Wave: "triangle",
		DroneFreqs: []float64{220.00, 261.63, 349.23, 392.00},
		Attack: 0.3, Release: 0.7, DetuneCents: 12,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// PLASMA  –  energetic, D minor (Dm → Bb → F → C)
// ════════════════════════════════════════════════════════════════════
func soundsPlasma() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{311.13, 415.3, 466.16, 622.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "sine", 0.085, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 349.23, 698.46, 0.15, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 174.61, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 200, 100, 0.13, 0.09

	dm := minor(73.42)
	bb := major(58.27)
	fm := major(87.31)
	cm := major(65.41)

	bass := [4]float64{146.83, 116.54, 174.61, 130.81}
	chords := [4][3]float64{dm, bb, fm, cm}
	normalMel := [4][]float64{
		{293.66, 349.23, 392.00, 293.66},
		{233.08, 293.66, 349.23, 233.08},
		{349.23, 440.00, 523.25, 349.23},
		{261.63, 329.63, 392.00, 261.63},
	}
	wildMel := [4][]float64{
		{587.33, 698.46, 783.99, 587.33},
		{466.16, 587.33, 698.46, 466.16},
		{698.46, 880.00, 1046.5, 698.46},
		{523.25, 659.25, 783.99, 523.25},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.035, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{146.83, 233.08, 349.23},
		Attack: 0.4, Release: 0.9, DetuneCents: 15,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 150, Wave: "square",
		DroneFreqs: []float64{293.66, 466.16, 698.46},
		Attack: 0.1, Release: 0.4, DetuneCents: 25,
		CutoffMin: 400, CutoffMax: 3000, NoiseGain: 0.02,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BRUTALIST  –  concrete minimalism, D major (D → A → D → A)
// ════════════════════════════════════════════════════════════════════
func soundsBrutalist() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{100, 150, 200, 120}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.14, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 120, "square", 0.07, 0.13
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 200, 400, 0.12, 0.11
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 400, 100, 0.14, 0.1
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 120, 60, 0.1, 0.12

	dm := major(73.42)
	am := major(55.00)

	bass := [4]float64{146.83, 110.00, 146.83, 110.00}
	chords := [4][3]float64{dm, am, dm, am}
	normalMel := [4][]float64{
		{293.66, 349.23, 392.00, 293.66},
		{220.00, 261.63, 293.66, 220.00},
		{293.66, 392.00, 440.00, 293.66},
		{220.00, 293.66, 349.23, 220.00},
	}
	wildMel := [4][]float64{
		{587.33, 698.46, 783.99, 587.33},
		{440.00, 523.25, 587.33, 440.00},
		{587.33, 783.99, 880.00, 587.33},
		{440.00, 587.33, 698.46, 440.00},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{73.42, 110.00, 146.83},
		Attack: 0.5, Release: 1.0,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 120, Wave: "square",
		DroneFreqs: []float64{146.83, 220.00, 293.66},
		Attack: 0.05, Release: 0.3,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// VOLCANO  –  molten, C major (C → D → F → G)
// ════════════════════════════════════════════════════════════════════
func soundsVolcano() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 180, "sine", 0.09, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 261.63, 523.25, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 98, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 160, 80, 0.13, 0.1

	cm := major(65.41)
	dm := major(73.42)
	fm := major(87.31)
	gm := major(98.00)

	bass := [4]float64{130.81, 146.83, 174.61, 196.00}
	chords := [4][3]float64{cm, dm, fm, gm}
	normalMel := [4][]float64{
		{261.63, 293.66, 329.63, 261.63},
		{293.66, 349.23, 392.00, 293.66},
		{349.23, 440.00, 523.25, 349.23},
		{392.00, 440.00, 493.88, 392.00},
	}
	wildMel := [4][]float64{
		{523.25, 587.33, 659.25, 523.25},
		{587.33, 698.46, 783.99, 587.33},
		{698.46, 880.00, 1046.5, 698.46},
		{783.99, 880.00, 987.77, 783.99},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{65.41, 130.81, 174.61},
		Attack: 1.0, Release: 2.0,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 110, Wave: "triangle",
		DroneFreqs: []float64{130.81, 261.63, 349.23},
		Attack: 0.2, Release: 0.6, NoiseGain: 0.015,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// AURORA  –  ethereal, E major (E → C#m → A → E)
// ════════════════════════════════════════════════════════════════════
func soundsAurora() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{659.25, 880, 987.77, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "sine", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 523.25, 880, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 704, 352, 0.18, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 220, 110, 0.15, 0.08

	em := major(82.41)
	csm := minor(69.30)
	am := major(55.00)

	bass := [4]float64{164.81, 138.59, 110.00, 164.81}
	chords := [4][3]float64{em, csm, am, em}
	normalMel := [4][]float64{
		{329.63, 369.99, 440.00, 329.63},
		{277.18, 329.63, 369.99, 277.18},
		{220.00, 261.63, 293.66, 220.00},
		{329.63, 369.99, 440.00, 329.63},
	}
	wildMel := [4][]float64{
		{659.25, 739.99, 880.00, 659.25},
		{554.37, 659.25, 739.99, 554.37},
		{440.00, 523.25, 587.33, 440.00},
		{659.25, 739.99, 880.00, 659.25},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{82.41, 110.00, 138.59, 164.81},
		Attack: 1.2, Release: 2.5,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 120, Wave: "sine",
		DroneFreqs: []float64{164.81, 220.00, 277.18, 329.63},
		Attack: 0.3, Release: 0.8, DetuneCents: 20,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// MATRIX  –  cyberpunk, E minor (Em → D → C → B)
// ════════════════════════════════════════════════════════════════════
func soundsMatrix() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 587.33, 659.25, 698.46}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.05, 0.09, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 523.25, "square", 0.045, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 880, 1318.5, 0.11, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 880, 330, 0.13, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 260, 130, 0.1, 0.09

	em := minor(82.41)
	dm := major(73.42)
	cm := major(65.41)
	bm := major(61.74)

	bass := [4]float64{164.81, 146.83, 130.81, 123.47}
	chords := [4][3]float64{em, dm, cm, bm}
	normalMel := [4][]float64{
		{329.63, 392.00, 440.00, 329.63},
		{293.66, 349.23, 392.00, 293.66},
		{261.63, 329.63, 392.00, 261.63},
		{246.94, 293.66, 349.23, 246.94},
	}
	wildMel := [4][]float64{
		{659.25, 783.99, 880.00, 659.25},
		{587.33, 698.46, 783.99, 587.33},
		{523.25, 659.25, 783.99, 523.25},
		{493.88, 587.33, 698.46, 493.88},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 164.81},
		Attack: 0.1, Release: 0.3,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 150, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 329.63},
		Attack: 0.05, Release: 0.15,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// OCEAN  –  flowing, G major (G → D → Em → C)
// ════════════════════════════════════════════════════════════════════
func soundsOcean() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 196, 220, 246.94}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.095, 0.095, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 349.23, 523.25, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 246.94, 0.18, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 140, 70, 0.15, 0.08

	gm := major(98.00)
	dm := major(73.42)
	em := minor(82.41)
	cm := major(65.41)

	bass := [4]float64{196.00, 146.83, 164.81, 130.81}
	chords := [4][3]float64{gm, dm, em, cm}
	normalMel := [4][]float64{
		{392.00, 440.00, 493.88, 392.00},
		{293.66, 349.23, 392.00, 293.66},
		{329.63, 392.00, 440.00, 329.63},
		{261.63, 329.63, 392.00, 261.63},
	}
	wildMel := [4][]float64{
		{783.99, 880.00, 987.77, 783.99},
		{587.33, 698.46, 783.99, 587.33},
		{659.25, 783.99, 880.00, 659.25},
		{523.25, 659.25, 783.99, 523.25},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{98.00, 130.81, 164.81, 196.00},
		Attack: 1.5, Release: 3.0, NoiseGain: 0.02,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{196.00, 261.63, 329.63, 392.00},
		Attack: 0.3, Release: 0.7, NoiseGain: 0.04,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// DOS  –  8-bit chip-tune, C major (C → G → Am → F)
// ════════════════════════════════════════════════════════════════════
func soundsDos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{800, 1000}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.12, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1000, "square", 0.03, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 400, 800, 0.1, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 800, 200, 0.1, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 300, 150, 0.08, 0.1

	cm := major(65.41)
	gm := major(98.00)
	am := minor(110.00)
	fm := major(87.31)

	bass := [4]float64{130.81, 196.00, 220.00, 174.61}
	chords := [4][3]float64{cm, gm, am, fm}
	normalMel := [4][]float64{
		{523.25, 587.33, 659.25, 523.25},
		{783.99, 880.00, 987.77, 783.99},
		{880.00, 1046.5, 1174.66, 880.00},
		{698.46, 783.99, 880.00, 698.46},
	}
	wildMel := [4][]float64{
		{1046.5, 1174.66, 1318.5, 1046.5},
		{1567.98, 1760.00, 1975.53, 1567.98},
		{1760.00, 2093.0, 2349.32, 1760.00},
		{1396.92, 1567.98, 1760.00, 1396.92},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 70, Wave: "square",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 220.00},
		Attack: 0.05, Release: 0.15,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "square",
		DroneFreqs: []float64{130.81, 261.63, 392.00, 440.00},
		Attack: 0.02, Release: 0.08,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETRO  –  warm nostalgia, C major (C → Am → F → G)
// ════════════════════════════════════════════════════════════════════
func soundsRetro() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{1046.5, 1318.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1200, "square", 0.04, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 800, 1600, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 1600, 400, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 400, 200, 0.1, 0.1

	cm := major(65.41)
	am := minor(55.00)
	fm := major(87.31)
	gm := major(98.00)

	bass := [4]float64{130.81, 110.00, 174.61, 196.00}
	chords := [4][3]float64{cm, am, fm, gm}
	normalMel := [4][]float64{
		{261.63, 329.63, 392.00, 261.63},
		{220.00, 261.63, 293.66, 220.00},
		{349.23, 392.00, 440.00, 349.23},
		{392.00, 440.00, 493.88, 392.00},
	}
	wildMel := [4][]float64{
		{523.25, 659.25, 783.99, 523.25},
		{440.00, 523.25, 587.33, 440.00},
		{698.46, 783.99, 880.00, 698.46},
		{783.99, 880.00, 987.77, 783.99},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "square",
		DroneFreqs: []float64{65.41, 110.00, 130.81, 174.61},
		Attack: 0.1, Release: 0.3,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "square",
		DroneFreqs: []float64{130.81, 220.00, 261.63, 349.23},
		Attack: 0.05, Release: 0.15,
		DetuneCents: 30, NoiseGain: 0.02,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// COSMOS  –  vast space, A major (A → D → C → G)
// ════════════════════════════════════════════════════════════════════
func soundsCosmos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{220, 277.18, 329.63, 392}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "sine", 0.09, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 392, 587.33, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 587.33, 196, 0.2, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 170, 85, 0.16, 0.08

	am := major(55.00)
	dm := major(73.42)
	cm := major(65.41)
	gm := major(98.00)

	bass := [4]float64{110.00, 146.83, 130.81, 196.00}
	chords := [4][3]float64{am, dm, cm, gm}
	normalMel := [4][]float64{
		{220.00, 261.63, 293.66, 220.00},
		{293.66, 349.23, 392.00, 293.66},
		{261.63, 329.63, 392.00, 261.63},
		{392.00, 440.00, 493.88, 392.00},
	}
	wildMel := [4][]float64{
		{440.00, 523.25, 587.33, 440.00},
		{587.33, 698.46, 783.99, 587.33},
		{523.25, 659.25, 783.99, 523.25},
		{783.99, 880.00, 987.77, 783.99},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{55.00, 110.00, 146.83, 196.00},
		Attack: 1.5, Release: 3.0,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{110.00, 220.00, 293.66, 392.00},
		Attack: 0.3, Release: 0.8,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETROFUTURE  –  retro sci-fi, E minor (Em → C → G → D)
// ════════════════════════════════════════════════════════════════════
func soundsRetrofuture() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 329.63, 440}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 330, 523.25, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 165, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 190, 95, 0.14, 0.09

	em := minor(82.41)
	cm := major(65.41)
	gm := major(98.00)
	dm := major(73.42)

	bass := [4]float64{164.81, 130.81, 196.00, 146.83}
	chords := [4][3]float64{em, cm, gm, dm}
	normalMel := [4][]float64{
		{329.63, 392.00, 440.00, 329.63},
		{261.63, 329.63, 392.00, 261.63},
		{392.00, 440.00, 493.88, 392.00},
		{293.66, 349.23, 392.00, 293.66},
	}
	wildMel := [4][]float64{
		{659.25, 783.99, 880.00, 659.25},
		{523.25, 659.25, 783.99, 523.25},
		{783.99, 880.00, 987.77, 783.99},
		{587.33, 698.46, 783.99, 587.33},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{82.41, 130.81, 164.81, 196.00},
		Attack: 0.8, Release: 1.8,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 130, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 329.63, 392.00},
		Attack: 0.05, Release: 0.15,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SPACEAGE  –  mid-century optimism, C major (C → Am → F → G)
// ════════════════════════════════════════════════════════════════════
func soundsSpaceage() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25, 880}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 554.37, "triangle", 0.075, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 659.25, 330, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 240, 120, 0.13, 0.08

	cm := major(65.41)
	am := minor(55.00)
	fm := major(87.31)
	gm := major(98.00)

	bass := [4]float64{130.81, 110.00, 174.61, 196.00}
	chords := [4][3]float64{cm, am, fm, gm}
	normalMel := [4][]float64{
		{261.63, 329.63, 392.00, 261.63},
		{220.00, 261.63, 293.66, 220.00},
		{349.23, 392.00, 440.00, 349.23},
		{392.00, 440.00, 493.88, 392.00},
	}
	wildMel := [4][]float64{
		{523.25, 659.25, 783.99, 523.25},
		{440.00, 523.25, 587.33, 440.00},
		{698.46, 783.99, 880.00, 698.46},
		{783.99, 880.00, 987.77, 783.99},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{65.41, 110.00, 130.81, 174.61},
		Attack: 0.8, Release: 1.5,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 130, Wave: "triangle",
		DroneFreqs: []float64{130.81, 220.00, 261.63, 349.23},
		Attack: 0.1, Release: 0.4, PulseInterval: 0.5,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TROPICALE  –  island vacation, C major (C → G → Am → F)
// ════════════════════════════════════════════════════════════════════
func soundsTropical() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.08, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "sine", 0.07, 0.07
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.08
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 660, 330, 0.17, 0.075
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 200, 100, 0.12, 0.07

	cm := major(65.41)
	gm := major(98.00)
	am := minor(110.00)
	fm := major(87.31)

	bass := [4]float64{130.81, 196.00, 220.00, 174.61}
	chords := [4][3]float64{cm, gm, am, fm}
	normalMel := [4][]float64{
		{261.63, 293.66, 329.63, 261.63},
		{392.00, 440.00, 493.88, 392.00},
		{440.00, 523.25, 587.33, 440.00},
		{349.23, 392.00, 440.00, 349.23},
	}
	wildMel := [4][]float64{
		{523.25, 587.33, 659.25, 523.25},
		{783.99, 880.00, 987.77, 783.99},
		{880.00, 1046.5, 1174.66, 880.00},
		{698.46, 783.99, 880.00, 698.46},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 65, Wave: "sine",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 220.00},
		Attack: 0.6, Release: 1.5, NoiseGain: 0.015,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "triangle",
		DroneFreqs: []float64{130.81, 261.63, 392.00, 440.00},
		Attack: 0.1, Release: 0.3, NoiseGain: 0.03,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// NOIR  –  smoky film noir, D minor (Dm → C → Bb → A)
// ════════════════════════════════════════════════════════════════════
func soundsNoir() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 220, 261.63}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.08, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 392, 0.18, 0.085
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "triangle", 330, 165, 0.2, 0.08
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.14, 0.08

	dm := minor(73.42)
	cm := major(65.41)
	bb := major(58.27)
	am := major(55.00)

	bass := [4]float64{146.83, 130.81, 116.54, 110.00}
	chords := [4][3]float64{dm, cm, bb, am}
	normalMel := [4][]float64{
		{293.66, 349.23, 392.00, 293.66},
		{261.63, 293.66, 329.63, 261.63},
		{233.08, 261.63, 293.66, 233.08},
		{220.00, 246.94, 277.18, 220.00},
	}
	wildMel := [4][]float64{
		{587.33, 698.46, 783.99, 587.33},
		{523.25, 587.33, 659.25, 523.25},
		{466.16, 523.25, 587.33, 466.16},
		{440.00, 493.88, 554.37, 440.00},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{73.42, 110.00, 130.81, 146.83},
		Attack: 1.0, Release: 2.0, NoiseGain: 0.015,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{146.83, 220.00, 261.63, 293.66},
		Attack: 0.3, Release: 0.7, PulseInterval: 1.2,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// CATHEDRAL  –  sacred organ, A minor (Am → C → G → D)
// ════════════════════════════════════════════════════════════════════
func soundsCathedral() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 392, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.14, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 293.66, "triangle", 0.11, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 392, 783.99, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 196, 0.22, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 180, 90, 0.18, 0.08

	am := minor(55.00)
	cm := major(65.41)
	gm := major(98.00)
	dm := major(73.42)

	bass := [4]float64{110.00, 130.81, 196.00, 146.83}
	chords := [4][3]float64{am, cm, gm, dm}
	normalMel := [4][]float64{
		{220.00, 261.63, 293.66, 220.00},
		{261.63, 329.63, 392.00, 261.63},
		{392.00, 440.00, 493.88, 392.00},
		{293.66, 349.23, 392.00, 293.66},
	}
	wildMel := [4][]float64{
		{440.00, 523.25, 587.33, 440.00},
		{523.25, 659.25, 783.99, 523.25},
		{783.99, 880.00, 987.77, 783.99},
		{587.33, 698.46, 783.99, 587.33},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{55.00, 110.00, 130.81, 196.00},
		Attack: 1.5, Release: 3.0,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{110.00, 220.00, 261.63, 392.00},
		Attack: 0.3, Release: 0.8,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SURVEILLANCE  –  spy tension, E major (E → C → D → B)
// ════════════════════════════════════════════════════════════════════
func soundsSurveillance() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "square", 0.04, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 660, 1320, 0.12, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 990, 330, 0.14, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 280, 140, 0.09, 0.09

	em := major(82.41)
	cm := major(65.41)
	dm := major(73.42)
	bm := major(61.74)

	bass := [4]float64{164.81, 130.81, 146.83, 123.47}
	chords := [4][3]float64{em, cm, dm, bm}
	normalMel := [4][]float64{
		{329.63, 369.99, 440.00, 329.63},
		{261.63, 293.66, 329.63, 261.63},
		{293.66, 349.23, 392.00, 293.66},
		{246.94, 277.18, 311.13, 246.94},
	}
	wildMel := [4][]float64{
		{659.25, 739.99, 880.00, 659.25},
		{523.25, 587.33, 659.25, 523.25},
		{587.33, 698.46, 783.99, 587.33},
		{493.88, 554.37, 622.25, 493.88},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 146.83, 164.81},
		Attack: 0.2, Release: 0.5,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 140, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 293.66, 329.63},
		Attack: 0.05, Release: 0.15, PulseInterval: 0.3,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BIOMECH  –  organic meets machine, C# major (C# → A → E → B)
// ════════════════════════════════════════════════════════════════════
func soundsBiomech() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{164.81, 246.94, 311.13, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 523.25, 0.18, 0.095
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	cs := major(69.30)
	am := major(55.00)
	emin := major(82.41)
	bm := major(61.74)

	bass := [4]float64{138.59, 110.00, 164.81, 123.47}
	chords := [4][3]float64{cs, am, emin, bm}
	normalMel := [4][]float64{
		{277.18, 329.63, 369.99, 277.18},
		{220.00, 261.63, 293.66, 220.00},
		{329.63, 392.00, 440.00, 329.63},
		{246.94, 293.66, 329.63, 246.94},
	}
	wildMel := [4][]float64{
		{554.37, 659.25, 739.99, 554.37},
		{440.00, 523.25, 587.33, 440.00},
		{659.25, 783.99, 880.00, 659.25},
		{493.88, 587.33, 659.25, 493.88},
	}

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{69.30, 110.00, 138.59, 164.81},
		Attack: 0.7, Release: 1.5, DetuneCents: 8,
		Melody: buildSong(bass, chords, normalMel, 0.5),
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "square",
		DroneFreqs: []float64{138.59, 220.00, 277.18, 329.63},
		Attack: 0.1, Release: 0.3, DetuneCents: 18,
		Melody: buildSong(bass, chords, wildMel, 0.25),
	}
	return s
}

func defaultSounds() themeSounds {
	return soundsNeon()
}

// themeSoundsJSON returns a JS object literal for embedding in <script> (safe JSON).
func themeSoundsJSON(themeName string) template.JS {
	p := defaultSounds()
	if x, ok := themeSoundPresets[themeName]; ok {
		p = x
	}
	b, err := json.Marshal(p)
	if err != nil {
		b, _ = json.Marshal(defaultSounds())
	}
	return template.JS(b) //nolint:gosec // JSON from fixed structs
}
