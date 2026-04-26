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
	// Each slot is one 16th-note at BPM. Patterns loop.
	Drums []string `json:"drums,omitempty"`
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
func ns(freq, dur, step float64) melodyNote {
	return melodyNote{Freq: freq, Dur: dur, Step: step}
}

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

// chordArp returns an up/down arpeggio of a triad; each note sustains for 'dur'
// while the sequencer only advances by 'step', so multiple voices overlap.
func chordArp(chord [3]float64, dur, step float64) []melodyNote {
	return []melodyNote{
		ns(chord[0], dur, step), ns(chord[1], dur, step),
		ns(chord[2], dur, step), ns(chord[1], dur, step),
	}
}

func bass(freq, dur, step float64) melodyNote {
	return ns(freq, dur, step)
}

func concatNotes(groups ...[]melodyNote) []melodyNote {
	var out []melodyNote
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

// buildMeasure creates one measure: bass + chord pad + melody top-line.
func buildMeasure(bassFreq float64, chord [3]float64, melody []melodyNote, step float64) []melodyNote {
	bassDur := step * 7.5
	chordDur := step * 3.0
	if bassFreq < 120 {
		bassFreq *= 2
	}
	if chord[0] < 200 {
		chord = [3]float64{chord[0] * 4, chord[1] * 4, chord[2] * 4}
	} else if chord[0] < 260 {
		chord = [3]float64{chord[0] * 2, chord[1] * 2, chord[2] * 2}
	}
	var notes []melodyNote
	notes = append(notes, bass(bassFreq, bassDur, step))
	notes = append(notes, chordArp(chord, chordDur, step)...)
	for _, nn := range melody {
		notes = append(notes, nn)
	}
	return notes
}

// buildSong creates a 4-measure loop from a chord progression and melodies.
func buildSong(bassFreqs [4]float64, chords [4][3]float64, melodies [4][]melodyNote, step float64) []melodyNote {
	var out []melodyNote
	for i := 0; i < 4; i++ {
		out = append(out, buildMeasure(bassFreqs[i], chords[i], melodies[i], step)...)
	}
	return out
}

// ── swing melody builder ────────────────────────────────────────────

// buildSwingMelody converts a list of freqs into a swing-pattern melody:
// even-index notes (on-beat) hold longer; odd-index notes (off-beat) are shorter.
func buildSwingMelody(freqs []float64, baseStep float64) []melodyNote {
	var out []melodyNote
	for i, freq := range freqs {
		if i%2 == 0 {
			out = append(out, ns(freq, baseStep*1.8, baseStep*1.0))
		} else {
			out = append(out, ns(freq, baseStep*0.8, baseStep*0.5))
		}
	}
	return out
}

// ── drum-pattern helpers ────────────────────────────────────────────

func pat(names ...string) []string {
	return names
}

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
// NEON  –  bright synth-pop in C major
// ════════════════════════════════════════════════════════════════════
func soundsNeon() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 330, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 523.25, 1046.5, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 880, 261.63, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.12, 0.1

	c, f, g1, c2 := major(65.41), major(87.31), major(98.00), major(130.81)

	normalMel := buildSong(
		[4]float64{130.81, 174.61, 196.00, 130.81},
		[4][3]float64{c, f, g1, c2},
		[4][]melodyNote{
			buildSwingMelody([]float64{523.25, 587.33, 659.25, 523.25}, 0.5),
			buildSwingMelody([]float64{698.46, 783.99, 880.00, 698.46}, 0.5),
			buildSwingMelody([]float64{783.99, 880.00, 1046.5, 783.99}, 0.5),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 1046.5}, 0.5),
		}, 0.5)

	wildMel := buildSong(
		[4]float64{261.63, 349.23, 392.00, 261.63},
		[4][3]float64{c, f, g1, c2},
		[4][]melodyNote{
			buildSwingMelody([]float64{1046.5, 1174.66, 1318.5, 1046.5}, 0.25),
			buildSwingMelody([]float64{1396.92, 1567.98, 1760.00, 1396.92}, 0.25),
			buildSwingMelody([]float64{1567.98, 1760.00, 2093.0, 1567.98}, 0.25),
			buildSwingMelody([]float64{1046.5, 1318.5, 1567.98, 2093.0}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 70, Wave: "sine",
		DroneFreqs: []float64{130.81, 174.61, 196.00, 261.63},
		Attack: 0.3, Release: 0.8, CutoffMin: 800, CutoffMax: 3000,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 130, Wave: "triangle",
		DroneFreqs: []float64{261.63, 349.23, 392.00, 523.25},
		Attack: 0.1, Release: 0.4, CutoffMin: 1500, CutoffMax: 6000, DetuneCents: 8,
		Drums: pat("kick", "hat", "hat", "_", "kick", "hat", "snare", "hat", "kick", "hat", "hat", "_", "kick", "hat", "snare", "hat"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TERMINAL  –  dark industrial, E minor
// ════════════════════════════════════════════════════════════════════
func soundsTerminal() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.11, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 800, "square", 0.045, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 600, 1200, 0.12, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 900, 400, 0.14, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1

	e, b, a := minor(82.41), minor(61.74), minor(110.00)

	normalMel := buildSong(
		[4]float64{164.81, 123.47, 220.00, 164.81}, [4][3]float64{e, b, a, e},
		[4][]melodyNote{
			buildSwingMelody([]float64{329.63, 392.00, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{246.94, 293.66, 329.63, 246.94}, 0.5),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.5),
			buildSwingMelody([]float64{329.63, 392.00, 493.88, 329.63}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{329.63, 246.94, 440.00, 329.63}, [4][3]float64{e, b, a, e},
		[4][]melodyNote{
			buildSwingMelody([]float64{659.25, 783.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{493.88, 587.33, 659.25, 493.88}, 0.25),
			buildSwingMelody([]float64{880.00, 1046.5, 1174.66, 880.00}, 0.25),
			buildSwingMelody([]float64{659.25, 783.99, 987.77, 659.25}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 65, Wave: "square",
		DroneFreqs: []float64{82.41, 123.47, 164.81},
		Attack: 0.1, Release: 0.3,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "clap", "snare", "kick", "hat", "kick"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 140, Wave: "square",
		DroneFreqs: []float64{164.81, 246.94, 329.63},
		Attack: 0.05, Release: 0.15,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "clap", "snare", "kick", "hat", "kick"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SYNTHWAVE  –  retro-future, A minor
// ════════════════════════════════════════════════════════════════════
func soundsSynthwave() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 440, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 110, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 150, 75, 0.14, 0.09

	am, f, c, g := minor(55.00), major(87.31), major(65.41), major(98.00)

	normalMel := buildSong(
		[4]float64{110.00, 174.61, 130.81, 196.00}, [4][3]float64{am, f, c, g},
		[4][]melodyNote{
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{349.23, 392.00, 440.00, 349.23}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{220.00, 349.23, 261.63, 392.00}, [4][3]float64{am, f, c, g},
		[4][]melodyNote{
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{698.46, 783.99, 880.00, 698.46}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 70, Wave: "sine",
		DroneFreqs: []float64{110.00, 130.81, 174.61, 196.00},
		Attack: 0.3, Release: 0.8,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 130, Wave: "triangle",
		DroneFreqs: []float64{220.00, 261.63, 349.23, 392.00},
		Attack: 0.1, Release: 0.4, DetuneCents: 12,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// PLASMA  –  energetic, D minor
// ════════════════════════════════════════════════════════════════════
func soundsPlasma() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{311.13, 415.3, 466.16, 622.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "sine", 0.085, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 349.23, 698.46, 0.15, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 174.61, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 200, 100, 0.13, 0.09

	dm, bb, f, c := minor(73.42), major(58.27), major(87.31), major(65.41)

	normalMel := buildSong(
		[4]float64{146.83, 116.54, 174.61, 130.81}, [4][3]float64{dm, bb, f, c},
		[4][]melodyNote{
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{233.08, 293.66, 349.23, 233.08}, 0.5),
			buildSwingMelody([]float64{349.23, 440.00, 523.25, 349.23}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{293.66, 233.08, 349.23, 261.63}, [4][3]float64{dm, bb, f, c},
		[4][]melodyNote{
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{466.16, 587.33, 698.46, 466.16}, 0.25),
			buildSwingMelody([]float64{698.46, 880.00, 1046.5, 698.46}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.035, BPM: 70, Wave: "triangle",
		DroneFreqs: []float64{146.83, 233.08, 349.23},
		Attack: 0.4, Release: 0.9, DetuneCents: 15,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 150, Wave: "square",
		DroneFreqs: []float64{293.66, 466.16, 698.46},
		Attack: 0.1, Release: 0.4, DetuneCents: 25,
		CutoffMin: 400, CutoffMax: 3000, NoiseGain: 0.02,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BRUTALIST  –  concrete minimalism, D major
// ════════════════════════════════════════════════════════════════════
func soundsBrutalist() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{100, 150, 200, 120}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.14, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 120, "square", 0.07, 0.13
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 200, 400, 0.12, 0.11
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 400, 100, 0.14, 0.1
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 120, 60, 0.1, 0.12

	d, a := major(73.42), major(55.00)

	normalMel := buildSong(
		[4]float64{146.83, 110.00, 146.83, 110.00}, [4][3]float64{d, a, d, a},
		[4][]melodyNote{
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{293.66, 392.00, 440.00, 293.66}, 0.5),
			buildSwingMelody([]float64{220.00, 293.66, 349.23, 220.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{293.66, 220.00, 293.66, 220.00}, [4][3]float64{d, a, d, a},
		[4][]melodyNote{
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{587.33, 783.99, 880.00, 587.33}, 0.25),
			buildSwingMelody([]float64{440.00, 587.33, 698.46, 440.00}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{73.42, 110.00, 146.83},
		Attack: 0.5, Release: 1.0,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "clap", "snare", "kick", "hat", "kick"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 120, Wave: "square",
		DroneFreqs: []float64{146.83, 220.00, 293.66},
		Attack: 0.05, Release: 0.3,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "clap", "snare", "kick", "hat", "kick"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// VOLCANO  –  molten, C major
// ════════════════════════════════════════════════════════════════════
func soundsVolcano() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 180, "sine", 0.09, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 261.63, 523.25, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 98, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 160, 80, 0.13, 0.1

	c, d, f, g := major(65.41), major(73.42), major(87.31), major(98.00)

	normalMel := buildSong(
		[4]float64{130.81, 146.83, 174.61, 196.00}, [4][3]float64{c, d, f, g},
		[4][]melodyNote{
			buildSwingMelody([]float64{261.63, 293.66, 329.63, 261.63}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{349.23, 440.00, 523.25, 349.23}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{261.63, 293.66, 349.23, 392.00}, [4][3]float64{c, d, f, g},
		[4][]melodyNote{
			buildSwingMelody([]float64{523.25, 587.33, 659.25, 523.25}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{698.46, 880.00, 1046.5, 698.46}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{65.41, 130.81, 174.61},
		Attack: 1.0, Release: 2.0,
		Drums: pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "_", "kick", "hat", "_", "_", "snare", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 110, Wave: "triangle",
		DroneFreqs: []float64{130.81, 261.63, 349.23},
		Attack: 0.2, Release: 0.6, NoiseGain: 0.015,
		Drums: pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "_", "kick", "hat", "_", "_", "snare", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// AURORA  –  ethereal, E major
// ════════════════════════════════════════════════════════════════════
func soundsAurora() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{659.25, 880, 987.77, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "sine", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 523.25, 880, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 704, 352, 0.18, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 220, 110, 0.15, 0.08

	e, bcm, a := major(82.41), minor(69.30), major(55.00)

	normalMel := buildSong(
		[4]float64{164.81, 138.59, 110.00, 164.81}, [4][3]float64{e, bcm, a, e},
		[4][]melodyNote{
			buildSwingMelody([]float64{329.63, 369.99, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{277.18, 329.63, 369.99, 277.18}, 0.5),
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{329.63, 369.99, 440.00, 329.63}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{329.63, 277.18, 220.00, 329.63}, [4][3]float64{e, bcm, a, e},
		[4][]melodyNote{
			buildSwingMelody([]float64{659.25, 739.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{554.37, 659.25, 739.99, 554.37}, 0.25),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{659.25, 739.99, 880.00, 659.25}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{82.41, 110.00, 138.59, 164.81},
		Attack: 1.2, Release: 2.5,
		Drums: pat("kick", "_", "_", "hat", "_", "_", "hat", "_", "kick", "_", "_", "hat", "_", "snare", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 120, Wave: "sine",
		DroneFreqs: []float64{164.81, 220.00, 277.18, 329.63},
		Attack: 0.3, Release: 0.8, DetuneCents: 20,
		Drums: pat("kick", "_", "_", "hat", "_", "_", "hat", "_", "kick", "_", "_", "hat", "_", "snare", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// MATRIX  –  cyberpunk, E minor
// ════════════════════════════════════════════════════════════════════
func soundsMatrix() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 587.33, 659.25, 698.46}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.05, 0.09, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 523.25, "square", 0.045, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 880, 1318.5, 0.11, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 880, 330, 0.13, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 260, 130, 0.1, 0.09

	e, dd, c, b := minor(82.41), major(73.42), major(65.41), major(61.74)

	normalMel := buildSong(
		[4]float64{164.81, 146.83, 130.81, 123.47}, [4][3]float64{e, dd, c, b},
		[4][]melodyNote{
			buildSwingMelody([]float64{329.63, 392.00, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{246.94, 293.66, 349.23, 246.94}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{329.63, 293.66, 261.63, 246.94}, [4][3]float64{e, dd, c, b},
		[4][]melodyNote{
			buildSwingMelody([]float64{659.25, 783.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{493.88, 587.33, 698.46, 493.88}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 70, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 164.81},
		Attack: 0.1, Release: 0.3,
		Drums: pat("kick", "_", "hat", "_", "kick", "snare", "hat", "_", "kick", "_", "hat", "_", "snare", "hat", "kick", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 150, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 329.63},
		Attack: 0.05, Release: 0.15,
		Drums: pat("kick", "_", "hat", "_", "kick", "snare", "hat", "_", "kick", "_", "hat", "_", "snare", "hat", "kick", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// OCEAN  –  flowing, G major
// ════════════════════════════════════════════════════════════════════
func soundsOcean() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 196, 220, 246.94}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.095, 0.095, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 349.23, 523.25, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 246.94, 0.18, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 140, 70, 0.15, 0.08

	g, d, em, c := major(98.00), major(73.42), minor(82.41), major(65.41)

	normalMel := buildSong(
		[4]float64{196.00, 146.83, 164.81, 130.81}, [4][3]float64{g, d, em, c},
		[4][]melodyNote{
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{329.63, 392.00, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{392.00, 293.66, 329.63, 261.63}, [4][3]float64{g, d, em, c},
		[4][]melodyNote{
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{659.25, 783.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{98.00, 130.81, 164.81, 196.00},
		Attack: 1.5, Release: 3.0, NoiseGain: 0.02,
		Drums: pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "kick", "_", "hat", "_", "_", "snare", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{196.00, 261.63, 329.63, 392.00},
		Attack: 0.3, Release: 0.7, NoiseGain: 0.04,
		Drums: pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "kick", "_", "hat", "_", "_", "snare", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// DOS  –  8-bit chip-tune, C major
// ════════════════════════════════════════════════════════════════════
func soundsDos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{800, 1000}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.12, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1000, "square", 0.03, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 400, 800, 0.1, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 800, 200, 0.1, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 300, 150, 0.08, 0.1

	c, gg, am, ff := major(65.41), major(98.00), minor(110.00), major(87.31)

	normalMel := buildSong(
		[4]float64{130.81, 196.00, 220.00, 174.61}, [4][3]float64{c, gg, am, ff},
		[4][]melodyNote{
			buildSwingMelody([]float64{523.25, 587.33, 659.25, 523.25}, 0.5),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.5),
			buildSwingMelody([]float64{880.00, 1046.5, 1174.66, 880.00}, 0.5),
			buildSwingMelody([]float64{698.46, 783.99, 880.00, 698.46}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{261.63, 392.00, 440.00, 349.23}, [4][3]float64{c, gg, am, ff},
		[4][]melodyNote{
			buildSwingMelody([]float64{1046.5, 1174.66, 1318.5, 1046.5}, 0.25),
			buildSwingMelody([]float64{1567.98, 1760.00, 1975.53, 1567.98}, 0.25),
			buildSwingMelody([]float64{1760.00, 2093.0, 2349.32, 1760.00}, 0.25),
			buildSwingMelody([]float64{1396.92, 1567.98, 1760.00, 1396.92}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 75, Wave: "square",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 220.00},
		Attack: 0.05, Release: 0.15,
		Drums: pat("kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "square",
		DroneFreqs: []float64{130.81, 261.63, 392.00, 440.00},
		Attack: 0.02, Release: 0.08,
		Drums: pat("kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETRO  –  warm nostalgia, C major
// ════════════════════════════════════════════════════════════════════
func soundsRetro() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{1046.5, 1318.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1200, "square", 0.04, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 800, 1600, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 1600, 400, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 400, 200, 0.1, 0.1

	c, am, f, ggg := major(65.41), minor(55.00), major(87.31), major(98.00)

	normalMel := buildSong(
		[4]float64{130.81, 110.00, 174.61, 196.00}, [4][3]float64{c, am, f, ggg},
		[4][]melodyNote{
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{349.23, 392.00, 440.00, 349.23}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{261.63, 220.00, 349.23, 392.00}, [4][3]float64{c, am, f, ggg},
		[4][]melodyNote{
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{698.46, 783.99, 880.00, 698.46}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 70, Wave: "square",
		DroneFreqs: []float64{65.41, 110.00, 130.81, 174.61},
		Attack: 0.1, Release: 0.3,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "square",
		DroneFreqs: []float64{130.81, 220.00, 261.63, 349.23},
		Attack: 0.05, Release: 0.15, DetuneCents: 30, NoiseGain: 0.02,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// COSMOS  –  vast space, A major
// ════════════════════════════════════════════════════════════════════
func soundsCosmos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{220, 277.18, 329.63, 392}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "sine", 0.09, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 392, 587.33, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 587.33, 196, 0.2, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 170, 85, 0.16, 0.08

	a, ddd, ccc, g := major(55.00), major(73.42), major(65.41), major(98.00)

	normalMel := buildSong(
		[4]float64{110.00, 146.83, 130.81, 196.00}, [4][3]float64{a, ddd, ccc, g},
		[4][]melodyNote{
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{220.00, 293.66, 261.63, 392.00}, [4][3]float64{a, ddd, ccc, g},
		[4][]melodyNote{
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{55.00, 110.00, 146.83, 196.00},
		Attack: 1.5, Release: 3.0,
		Drums: pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "kick", "_", "hat", "_", "_", "snare", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{110.00, 220.00, 293.66, 392.00},
		Attack: 0.3, Release: 0.8,
		Drums: pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "kick", "_", "hat", "_", "_", "snare", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETROFUTURE  –  retro sci-fi, E minor
// ════════════════════════════════════════════════════════════════════
func soundsRetrofuture() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 329.63, 440}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 330, 523.25, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 165, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 190, 95, 0.14, 0.09

	em, ccc, gg, ddd := minor(82.41), major(65.41), major(98.00), major(73.42)

	normalMel := buildSong(
		[4]float64{164.81, 130.81, 196.00, 146.83}, [4][3]float64{em, ccc, gg, ddd},
		[4][]melodyNote{
			buildSwingMelody([]float64{329.63, 392.00, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{329.63, 261.63, 392.00, 293.66}, [4][3]float64{em, ccc, gg, ddd},
		[4][]melodyNote{
			buildSwingMelody([]float64{659.25, 783.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{82.41, 130.81, 164.81, 196.00},
		Attack: 0.8, Release: 1.8,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 130, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 329.63, 392.00},
		Attack: 0.05, Release: 0.15,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SPACEAGE  –  mid-century optimism, C major
// ════════════════════════════════════════════════════════════════════
func soundsSpaceage() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25, 880}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 554.37, "triangle", 0.075, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 659.25, 330, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 240, 120, 0.13, 0.08

	c, am, ff, ggg := major(65.41), minor(55.00), major(87.31), major(98.00)

	normalMel := buildSong(
		[4]float64{130.81, 110.00, 174.61, 196.00}, [4][3]float64{c, am, ff, ggg},
		[4][]melodyNote{
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{349.23, 392.00, 440.00, 349.23}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{261.63, 220.00, 349.23, 392.00}, [4][3]float64{c, am, ff, ggg},
		[4][]melodyNote{
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{698.46, 783.99, 880.00, 698.46}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{65.41, 110.00, 130.81, 174.61},
		Attack: 0.8, Release: 1.5,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "hat", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 130, Wave: "triangle",
		DroneFreqs: []float64{130.81, 220.00, 261.63, 349.23},
		Attack: 0.1, Release: 0.4,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "hat", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TROPICALE  –  island vacation, C major
// ════════════════════════════════════════════════════════════════════
func soundsTropical() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.08, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "sine", 0.07, 0.07
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.08
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 660, 330, 0.17, 0.075
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 200, 100, 0.12, 0.07

	c, gg, am, ff := major(65.41), major(98.00), minor(110.00), major(87.31)

	normalMel := buildSong(
		[4]float64{130.81, 196.00, 220.00, 174.61}, [4][3]float64{c, gg, am, ff},
		[4][]melodyNote{
			buildSwingMelody([]float64{261.63, 293.66, 329.63, 261.63}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.5),
			buildSwingMelody([]float64{349.23, 392.00, 440.00, 349.23}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{261.63, 392.00, 440.00, 349.23}, [4][3]float64{c, gg, am, ff},
		[4][]melodyNote{
			buildSwingMelody([]float64{523.25, 587.33, 659.25, 523.25}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
			buildSwingMelody([]float64{880.00, 1046.5, 1174.66, 880.00}, 0.25),
			buildSwingMelody([]float64{698.46, 783.99, 880.00, 698.46}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 70, Wave: "sine",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 220.00},
		Attack: 0.6, Release: 1.5, NoiseGain: 0.015,
		Drums: pat("kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "clap", "kick", "hat", "snare", "hat"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "triangle",
		DroneFreqs: []float64{130.81, 261.63, 392.00, 440.00},
		Attack: 0.1, Release: 0.3, NoiseGain: 0.03,
		Drums: pat("kick", "hat", "snare", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "snare", "clap", "kick", "hat", "snare", "hat"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// NOIR  –  smoky film noir, D minor
// ════════════════════════════════════════════════════════════════════
func soundsNoir() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 220, 261.63}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.08, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 392, 0.18, 0.085
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "triangle", 330, 165, 0.2, 0.08
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.14, 0.08

	dm, ccc, bb, aa := minor(73.42), major(65.41), major(58.27), major(55.00)

	normalMel := buildSong(
		[4]float64{146.83, 130.81, 116.54, 110.00}, [4][3]float64{dm, ccc, bb, aa},
		[4][]melodyNote{
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{261.63, 293.66, 329.63, 261.63}, 0.5),
			buildSwingMelody([]float64{233.08, 261.63, 293.66, 233.08}, 0.5),
			buildSwingMelody([]float64{220.00, 246.94, 277.18, 220.00}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{293.66, 261.63, 233.08, 220.00}, [4][3]float64{dm, ccc, bb, aa},
		[4][]melodyNote{
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{523.25, 587.33, 659.25, 523.25}, 0.25),
			buildSwingMelody([]float64{466.16, 523.25, 587.33, 466.16}, 0.25),
			buildSwingMelody([]float64{440.00, 493.88, 554.37, 440.00}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{73.42, 110.00, 130.81, 146.83},
		Attack: 1.0, Release: 2.0, NoiseGain: 0.015,
		Drums: pat("kick", "_", "_", "hat", "_", "snare", "_", "hat", "kick", "_", "_", "hat", "_", "snare", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{146.83, 220.00, 261.63, 293.66},
		Attack: 0.3, Release: 0.7,
		Drums: pat("kick", "_", "_", "hat", "_", "snare", "_", "hat", "kick", "_", "_", "hat", "_", "snare", "hat", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// CATHEDRAL  –  sacred organ, A minor
// ════════════════════════════════════════════════════════════════════
func soundsCathedral() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 392, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.14, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 293.66, "triangle", 0.11, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 392, 783.99, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 196, 0.22, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 180, 90, 0.18, 0.08

	am, ccc, gg, ddd := minor(55.00), major(65.41), major(98.00), major(73.42)

	normalMel := buildSong(
		[4]float64{110.00, 130.81, 196.00, 146.83}, [4][3]float64{am, ccc, gg, ddd},
		[4][]melodyNote{
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{261.63, 329.63, 392.00, 261.63}, 0.5),
			buildSwingMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{220.00, 261.63, 392.00, 293.66}, [4][3]float64{am, ccc, gg, ddd},
		[4][]melodyNote{
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{523.25, 659.25, 783.99, 523.25}, 0.25),
			buildSwingMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "sine",
		DroneFreqs: []float64{55.00, 110.00, 130.81, 196.00},
		Attack: 1.5, Release: 3.0,
		Drums: pat("kick", "_", "_", "_", "_", "snare", "_", "_", "kick", "_", "_", "_", "_", "snare", "_", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{110.00, 220.00, 261.63, 392.00},
		Attack: 0.3, Release: 0.8,
		Drums: pat("kick", "_", "_", "_", "_", "snare", "_", "_", "kick", "_", "_", "_", "_", "snare", "_", "_"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SURVEILLANCE  –  spy tension, E major
// ════════════════════════════════════════════════════════════════════
func soundsSurveillance() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "square", 0.04, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 660, 1320, 0.12, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 990, 330, 0.14, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 280, 140, 0.09, 0.09

	e, ccc, ddd, bbb := major(82.41), major(65.41), major(73.42), major(61.74)

	normalMel := buildSong(
		[4]float64{164.81, 130.81, 146.83, 123.47}, [4][3]float64{e, ccc, ddd, bbb},
		[4][]melodyNote{
			buildSwingMelody([]float64{329.63, 369.99, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{261.63, 293.66, 329.63, 261.63}, 0.5),
			buildSwingMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.5),
			buildSwingMelody([]float64{246.94, 277.18, 311.13, 246.94}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{329.63, 261.63, 293.66, 246.94}, [4][3]float64{e, ccc, ddd, bbb},
		[4][]melodyNote{
			buildSwingMelody([]float64{659.25, 739.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{523.25, 587.33, 659.25, 523.25}, 0.25),
			buildSwingMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.25),
			buildSwingMelody([]float64{493.88, 554.37, 622.25, 493.88}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 70, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 146.83, 164.81},
		Attack: 0.2, Release: 0.5,
		Drums: pat("kick", "_", "hat", "kick", "snare", "_", "hat", "kick", "kick", "_", "hat", "clap", "snare", "_", "hat", "kick"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 140, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 293.66, 329.63},
		Attack: 0.05, Release: 0.15,
		Drums: pat("kick", "_", "hat", "kick", "snare", "_", "hat", "kick", "kick", "_", "hat", "clap", "snare", "_", "hat", "kick"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BIOMECH  –  organic meets machine, C# major
// ════════════════════════════════════════════════════════════════════
func soundsBiomech() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{164.81, 246.94, 311.13, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 523.25, 0.18, 0.095
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	cs, am, em, bm := major(69.30), major(55.00), major(82.41), major(61.74)

	normalMel := buildSong(
		[4]float64{138.59, 110.00, 164.81, 123.47}, [4][3]float64{cs, am, em, bm},
		[4][]melodyNote{
			buildSwingMelody([]float64{277.18, 329.63, 369.99, 277.18}, 0.5),
			buildSwingMelody([]float64{220.00, 261.63, 293.66, 220.00}, 0.5),
			buildSwingMelody([]float64{329.63, 392.00, 440.00, 329.63}, 0.5),
			buildSwingMelody([]float64{246.94, 293.66, 329.63, 246.94}, 0.5),
		}, 0.5)
	wildMel := buildSong(
		[4]float64{277.18, 220.00, 329.63, 246.94}, [4][3]float64{cs, am, em, bm},
		[4][]melodyNote{
			buildSwingMelody([]float64{554.37, 659.25, 739.99, 554.37}, 0.25),
			buildSwingMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.25),
			buildSwingMelody([]float64{659.25, 783.99, 880.00, 659.25}, 0.25),
			buildSwingMelody([]float64{493.88, 587.33, 659.25, 493.88}, 0.25),
		}, 0.25)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 60, Wave: "triangle",
		DroneFreqs: []float64{69.30, 110.00, 138.59, 164.81},
		Attack: 0.7, Release: 1.5, DetuneCents: 8,
		Drums: pat("kick", "hat", "kick", "snare", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "snare", "kick", "hat", "snare", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "square",
		DroneFreqs: []float64{138.59, 220.00, 277.18, 329.63},
		Attack: 0.1, Release: 0.3, DetuneCents: 18,
		Drums: pat("kick", "hat", "kick", "snare", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "snare", "kick", "hat", "snare", "clap"),
		Melody: wildMel,
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
