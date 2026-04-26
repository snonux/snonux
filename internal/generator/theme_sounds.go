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

// power returns a root+5th power chord (only two voices — aggressive and clear).
func power(freq float64) [3]float64 {
	return [3]float64{freq, freq * intPerf5th, 0}
}

// chordArp returns an up/down arpeggio of a triad; each note sustains for 'dur'
// while the sequencer only advances by 'step', so multiple voices overlap.
func chordArp(chord [3]float64, dur, step float64) []melodyNote {
	return []melodyNote{
		ns(chord[0], dur, step), ns(chord[1], dur, step),
		ns(chord[2], dur, step), ns(chord[1], dur, step),
	}
}

// powerArp arpeggiates a power chord (root→5th→root→5th).
func powerArp(chord [3]float64, dur, step float64) []melodyNote {
	return []melodyNote{
		ns(chord[0], dur, step), ns(chord[1], dur, step),
		ns(chord[0], dur, step), ns(chord[1], dur, step),
	}
}

func bassNote(freq, dur, step float64) melodyNote {
	return ns(freq, dur, step)
}

// ── song builders ──────────────────────────────────────────────────

// buildRockSong creates a 4-bar loop with power-chord pads, sub-bass, and a
// driving melody line. BPM is fast; every slot is a 16th-note.
func buildRockSong(
	bassOct [4]float64, chords [4][3]float64, melody [4][]melodyNote,
	step float64,
) []melodyNote {
	var out []melodyNote
	for i := 0; i < 4; i++ {
		// Sub-bass: low octave, long sustain so drone + bass overlap
		out = append(out, bassNote(bassOct[i]/2, step*7.9, step*8))
		// Main bass hits on each beat
		out = append(out, bassNote(bassOct[i], step*3.8, step*4))
		// Power chord stab (sustains across the measure)
		out = append(out, powerArp(chords[i], step*3.5, step)...)
		// Melody top-line
		out = append(out, melody[i]...)
	}
	return out
}

// buildDriveMelody: straight eighth-note drive, no swing. Each note is a short
// punch stab so the line stays rhythmic and aggressive.
func buildDriveMelody(freqs []float64, step float64) []melodyNote {
	var out []melodyNote
	for _, freq := range freqs {
		out = append(out, ns(freq, step*0.85, step))
	}
	return out
}

// buildStabMelody: short 16th-note stabs with gaps for air. Perfect for fast
// drum-and-bass, hardstyle, or industrial grooves.
func buildStabMelody(freqs []float64, step float64) []melodyNote {
	var out []melodyNote
	for i, freq := range freqs {
		// On-beat stabs hit hard; off-beat ghosts are tighter
		dur := step * 0.4
		if i%2 == 0 {
			dur = step * 0.8
		}
		out = append(out, ns(freq, dur, step))
	}
	return out
}

// ── drum-pattern helpers ────────────────────────────────────────────

func pat(names ...string) []string {
	return names
}

// ════════════════════════════════════════════════════════════════════
// NEON  –  synth-pop banger, C major, 150/220
// ════════════════════════════════════════════════════════════════════
func soundsNeon() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 330, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 523.25, 1046.5, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 880, 261.63, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.12, 0.1

	c, f, g1, c2 := power(65.41), power(87.31), power(98.00), power(130.81)

	normalMel := buildRockSong(
		[4]float64{130.81, 174.61, 196.00, 130.81},
		[4][3]float64{c, f, g1, c2},
		[4][]melodyNote{
			buildDriveMelody([]float64{523.25, 587.33, 659.25, 783.99}, 0.25),
			buildDriveMelody([]float64{698.46, 783.99, 880.00, 1046.5}, 0.25),
			buildDriveMelody([]float64{783.99, 880.00, 1046.5, 783.99}, 0.25),
			buildDriveMelody([]float64{523.25, 659.25, 783.99, 1046.5}, 0.25),
		}, 0.25)

	wildMel := buildRockSong(
		[4]float64{261.63, 349.23, 392.00, 261.63},
		[4][3]float64{c, f, g1, c2},
		[4][]melodyNote{
			buildStabMelody([]float64{1046.5, 1174.66, 1318.5, 1567.98}, 0.125),
			buildStabMelody([]float64{1396.92, 1567.98, 1760.00, 2093.0}, 0.125),
			buildStabMelody([]float64{1567.98, 1760.00, 2093.0, 1567.98}, 0.125),
			buildStabMelody([]float64{1046.5, 1318.5, 1567.98, 2093.0}, 0.125),
		}, 0.125)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 150, Wave: "square",
		DroneFreqs: []float64{130.81, 196.00, 261.63},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 220, Wave: "square",
		DroneFreqs: []float64{130.81, 196.00, 261.63, 392.00},
		Attack: 0.02, Release: 0.08, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 12,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TERMINAL  –  dark industrial, E minor, 160/240
// ════════════════════════════════════════════════════════════════════
func soundsTerminal() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.11, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 800, "square", 0.045, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 600, 1200, 0.12, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 900, 400, 0.14, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1

	e, b, a := power(82.41), power(61.74), power(110.00)

	normalMel := buildRockSong(
		[4]float64{164.81, 123.47, 220.00, 164.81}, [4][3]float64{e, b, a, e},
		[4][]melodyNote{
			buildDriveMelody([]float64{329.63, 392.00, 440.00, 523.25}, 0.233),
			buildDriveMelody([]float64{246.94, 293.66, 329.63, 392.00}, 0.233),
			buildDriveMelody([]float64{440.00, 523.25, 587.33, 440.00}, 0.233),
			buildDriveMelody([]float64{329.63, 392.00, 493.88, 329.63}, 0.233),
		}, 0.233)
	wildMel := buildRockSong(
		[4]float64{329.63, 246.94, 440.00, 329.63}, [4][3]float64{e, b, a, e},
		[4][]melodyNote{
			buildStabMelody([]float64{659.25, 783.99, 880.00, 1046.5}, 0.117),
			buildStabMelody([]float64{493.88, 587.33, 659.25, 783.99}, 0.117),
			buildStabMelody([]float64{880.00, 1046.5, 1174.66, 880.00}, 0.117),
			buildStabMelody([]float64{659.25, 783.99, 987.77, 659.25}, 0.117),
		}, 0.117)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 160, Wave: "square",
		DroneFreqs: []float64{82.41, 123.47, 164.81},
		Attack: 0.02, Release: 0.1, CutoffMin: 1000, CutoffMax: 5000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "snare", "kick", "kick", "hat", "kick", "snare", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 240, Wave: "square",
		DroneFreqs: []float64{164.81, 246.94, 329.63, 440.00},
		Attack: 0.02, Release: 0.05, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 15,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "snare", "kick", "kick", "hat", "kick", "snare", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SYNTHWAVE  –  retro-future racing, A minor, 140/200
// ════════════════════════════════════════════════════════════════════
func soundsSynthwave() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 440, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 110, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 150, 75, 0.14, 0.09

	am, f, c, g := minor(55.00), power(87.31), power(65.41), power(98.00)

	normalMel := buildRockSong(
		[4]float64{110.00, 174.61, 130.81, 196.00}, [4][3]float64{am, f, c, g},
		[4][]melodyNote{
			buildDriveMelody([]float64{220.00, 261.63, 293.66, 349.23}, 0.268),
			buildDriveMelody([]float64{349.23, 392.00, 440.00, 523.25}, 0.268),
			buildDriveMelody([]float64{261.63, 329.63, 392.00, 440.00}, 0.268),
			buildDriveMelody([]float64{220.00, 293.66, 349.23, 392.00}, 0.268),
		}, 0.268)
	wildMel := buildRockSong(
		[4]float64{220.00, 349.23, 261.63, 392.00}, [4][3]float64{am, f, c, g},
		[4][]melodyNote{
			buildStabMelody([]float64{440.00, 523.25, 587.33, 698.46}, 0.134),
			buildStabMelody([]float64{698.46, 783.99, 880.00, 1046.5}, 0.134),
			buildStabMelody([]float64{523.25, 659.25, 783.99, 880.00}, 0.134),
			buildStabMelody([]float64{440.00, 587.33, 698.46, 783.99}, 0.134),
		}, 0.134)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 140, Wave: "square",
		DroneFreqs: []float64{110.00, 174.61, 220.00},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 200, Wave: "square",
		DroneFreqs: []float64{220.00, 349.23, 440.00, 523.25},
		Attack: 0.02, Release: 0.06, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 12,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "kick", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// PLASMA  –  white-hot fusion drive, F# minor, 170/250
// ════════════════════════════════════════════════════════════════════
func soundsPlasma() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{369.99, 493.88, 587.33, 739.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 370, 740, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 740, 246.94, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 120, 60, 0.1, 0.1

	fsm, bm, cs, d := minor(92.50), minor(61.74), power(69.30), power(73.42)

	normalMel := buildRockSong(
		[4]float64{185.00, 246.94, 138.59, 146.83}, [4][3]float64{fsm, bm, cs, d},
		[4][]melodyNote{
			buildDriveMelody([]float64{369.99, 440.00, 493.88, 554.37}, 0.22),
			buildDriveMelody([]float64{493.88, 554.37, 587.33, 369.99}, 0.22),
			buildDriveMelody([]float64{277.18, 329.63, 369.99, 440.00}, 0.22),
			buildDriveMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.22),
		}, 0.22)
	wildMel := buildRockSong(
		[4]float64{369.99, 493.88, 277.18, 293.66}, [4][3]float64{fsm, bm, cs, d},
		[4][]melodyNote{
			buildStabMelody([]float64{739.99, 880.00, 987.77, 1108.73}, 0.11),
			buildStabMelody([]float64{987.77, 1108.73, 1174.66, 739.99}, 0.11),
			buildStabMelody([]float64{554.37, 659.25, 739.99, 880.00}, 0.11),
			buildStabMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.11),
		}, 0.11)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 170, Wave: "square",
		DroneFreqs: []float64{92.50, 138.59, 185.00},
		Attack: 0.02, Release: 0.12, CutoffMin: 1400, CutoffMax: 5500,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "snare", "kick", "kick", "hat", "kick", "snare", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 250, Wave: "square",
		DroneFreqs: []float64{92.50, 138.59, 185.00, 277.18},
		Attack: 0.02, Release: 0.04, CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 15,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "snare", "kick", "kick", "hat", "kick", "snare", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BRUTALIST  –  concrete demolition, C minor, 140/210
// ════════════════════════════════════════════════════════════════════
func soundsBrutalist() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 329.63, 392.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.06, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 440, 880, 0.15, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 660, 220, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1

	cm, gm, dm := minor(65.41), minor(98.00), minor(73.42)
	normalMel := buildRockSong(
		[4]float64{130.81, 196.00, 146.83, 130.81}, [4][3]float64{cm, gm, dm, cm},
		[4][]melodyNote{
			buildDriveMelody([]float64{261.63, 311.13, 392.00, 523.25}, 0.268),
			buildDriveMelody([]float64{392.00, 415.30, 523.25, 392.00}, 0.268),
			buildDriveMelody([]float64{293.66, 349.23, 392.00, 293.66}, 0.268),
			buildDriveMelody([]float64{261.63, 311.13, 392.00, 523.25}, 0.268),
		}, 0.268)
	wildMel := buildRockSong(
		[4]float64{261.63, 196.00, 146.83, 261.63}, [4][3]float64{cm, gm, dm, cm},
		[4][]melodyNote{
			buildStabMelody([]float64{523.25, 622.25, 783.99, 1046.5}, 0.134),
			buildStabMelody([]float64{783.99, 830.61, 1046.5, 783.99}, 0.134),
			buildStabMelody([]float64{587.33, 698.46, 783.99, 587.33}, 0.134),
			buildStabMelody([]float64{523.25, 622.25, 783.99, 1046.5}, 0.134),
		}, 0.134)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 140, Wave: "square",
		DroneFreqs: []float64{65.41, 98.00, 130.81},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 210, Wave: "square",
		DroneFreqs: []float64{65.41, 98.00, 130.81, 196.00},
		Attack: 0.02, Release: 0.05, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 18,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// VOLCANO  –  molten hardstyle, D minor, 160/240
// ════════════════════════════════════════════════════════════════════
func soundsVolcano() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 369.99, 440.00, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 293.66, 587.33, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 146.83, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.12, 0.09

	dm, am, gm := minor(73.42), power(110.00), power(98.00)
	normalMel := buildRockSong(
		[4]float64{146.83, 220.00, 196.00, 146.83}, [4][3]float64{dm, am, gm, dm},
		[4][]melodyNote{
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 493.88}, 0.233),
			buildDriveMelody([]float64{440.00, 493.88, 587.33, 440.00}, 0.233),
			buildDriveMelody([]float64{392.00, 440.00, 493.88, 392.00}, 0.233),
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 587.33}, 0.233),
		}, 0.233)
	wildMel := buildRockSong(
		[4]float64{293.66, 440.00, 392.00, 293.66}, [4][3]float64{dm, am, gm, dm},
		[4][]melodyNote{
			buildStabMelody([]float64{587.33, 698.46, 880.00, 987.77}, 0.117),
			buildStabMelody([]float64{880.00, 987.77, 1174.66, 880.00}, 0.117),
			buildStabMelody([]float64{783.99, 880.00, 987.77, 783.99}, 0.117),
			buildStabMelody([]float64{587.33, 698.46, 880.00, 1174.66}, 0.117),
		}, 0.117)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 160, Wave: "square",
		DroneFreqs: []float64{73.42, 110.00, 146.83},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "kick", "hat", "snare", "kick", "kick", "hat", "kick", "snare", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 240, Wave: "square",
		DroneFreqs: []float64{73.42, 110.00, 146.83, 196.00},
		Attack: 0.02, Release: 0.05, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 12,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "snare", "kick", "kick", "hat", "kick", "snare", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// AURORA  –  arctic storm surge, G major, 130/190
// ════════════════════════════════════════════════════════════════════
func soundsAurora() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{392.00, 493.88, 587.33, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 261.63, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 390, 780, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523, 196, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 180, 90, 0.12, 0.09

	gm, d, em, c := power(98.00), power(73.42), power(82.41), power(65.41)
	normalMel := buildRockSong(
		[4]float64{196.00, 146.83, 164.81, 130.81}, [4][3]float64{gm, d, em, c},
		[4][]melodyNote{
			buildDriveMelody([]float64{392.00, 440.00, 493.88, 587.33}, 0.288),
			buildDriveMelody([]float64{293.66, 349.23, 392.00, 440.00}, 0.288),
			buildDriveMelody([]float64{329.63, 392.00, 440.00, 523.25}, 0.288),
			buildDriveMelody([]float64{261.63, 329.63, 392.00, 523.25}, 0.288),
		}, 0.288)
	wildMel := buildRockSong(
		[4]float64{196.00, 146.83, 164.81, 130.81}, [4][3]float64{gm, d, em, c},
		[4][]melodyNote{
			buildStabMelody([]float64{783.99, 880.00, 987.77, 1174.66}, 0.144),
			buildStabMelody([]float64{587.33, 698.46, 783.99, 880.00}, 0.144),
			buildStabMelody([]float64{659.25, 783.99, 880.00, 1046.5}, 0.144),
			buildStabMelody([]float64{523.25, 659.25, 783.99, 1046.5}, 0.144),
		}, 0.144)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 130, Wave: "square",
		DroneFreqs: []float64{98.00, 146.83, 196.00},
		Attack: 0.02, Release: 0.18, CutoffMin: 1000, CutoffMax: 4500,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 190, Wave: "square",
		DroneFreqs: []float64{98.00, 146.83, 196.00, 261.63},
		Attack: 0.02, Release: 0.08, CutoffMin: 1800, CutoffMax: 6000, DetuneCents: 12,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// MATRIX  –  hard-grid cyberpunk, C minor, 160/240
// ════════════════════════════════════════════════════════════════════
func soundsMatrix() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 311.13, 392.00, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 440, 880, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 660, 220, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.1, 0.1

	cm, fm, gm, cm2 := minor(65.41), power(87.31), power(98.00), minor(130.81)
	normalMel := buildRockSong(
		[4]float64{130.81, 196.00, 196.00, 130.81}, [4][3]float64{cm, fm, gm, cm2},
		[4][]melodyNote{
			buildDriveMelody([]float64{523.25, 587.33, 659.25, 783.99}, 0.233),
			buildDriveMelody([]float64{698.46, 783.99, 880.00, 523.25}, 0.233),
			buildDriveMelody([]float64{783.99, 880.00, 932.33, 587.33}, 0.233),
			buildDriveMelody([]float64{523.25, 659.25, 783.99, 932.33}, 0.233),
		}, 0.233)
	wildMel := buildRockSong(
		[4]float64{261.63, 392.00, 392.00, 261.63}, [4][3]float64{cm, fm, gm, cm2},
		[4][]melodyNote{
			buildStabMelody([]float64{1046.5, 1174.66, 1318.5, 1567.98}, 0.117),
			buildStabMelody([]float64{1396.92, 1567.98, 1760.00, 1046.5}, 0.117),
			buildStabMelody([]float64{1567.98, 1760.00, 1864.66, 1174.66}, 0.117),
			buildStabMelody([]float64{1046.5, 1318.5, 1567.98, 1864.66}, 0.117),
		}, 0.117)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 160, Wave: "square",
		DroneFreqs: []float64{65.41, 130.81, 196.00},
		Attack: 0.02, Release: 0.12, CutoffMin: 1200, CutoffMax: 5000, DetuneCents: 6,
		Drums: pat("kick", "_", "hat", "_", "kick", "snare", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 240, Wave: "square",
		DroneFreqs: []float64{130.81, 196.00, 261.63, 392.00},
		Attack: 0.02, Release: 0.04, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 18,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// OCEAN  –  deep-current surge, E minor, 130/200
// ════════════════════════════════════════════════════════════════════
func soundsOcean() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{329.63, 392.00, 493.88, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 220, "sine", 0.08, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 329.63, 659.25, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 493.88, 164.81, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 140, 70, 0.12, 0.09

	em, c, d, e2 := minor(82.41), power(65.41), power(73.42), minor(164.81)
	normalMel := buildRockSong(
		[4]float64{164.81, 130.81, 146.83, 329.63}, [4][3]float64{em, c, d, e2},
		[4][]melodyNote{
			buildDriveMelody([]float64{329.63, 392.00, 440.00, 493.88}, 0.288),
			buildDriveMelody([]float64{261.63, 311.13, 392.00, 440.00}, 0.288),
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 523.25}, 0.288),
			buildDriveMelody([]float64{329.63, 392.00, 493.88, 659.25}, 0.288),
		}, 0.288)
	wildMel := buildRockSong(
		[4]float64{329.63, 261.63, 293.66, 659.25}, [4][3]float64{em, c, d, e2},
		[4][]melodyNote{
			buildStabMelody([]float64{659.25, 783.99, 880.00, 987.77}, 0.144),
			buildStabMelody([]float64{523.25, 622.25, 783.99, 880.00}, 0.144),
			buildStabMelody([]float64{587.33, 698.46, 880.00, 1046.5}, 0.144),
			buildStabMelody([]float64{659.25, 783.99, 987.77, 1318.5}, 0.144),
		}, 0.144)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 130, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 164.81},
		Attack: 0.02, Release: 0.18, CutoffMin: 1000, CutoffMax: 4500,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 200, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 164.81, 196.00},
		Attack: 0.02, Release: 0.08, CutoffMin: 1800, CutoffMax: 6500, DetuneCents: 12,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// DOS  –  chiptune blast, C major, 180/260
// ════════════════════════════════════════════════════════════════════
func soundsDos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 329.63, 392.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.05, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 440, 880, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 660, 220, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.1, 0.1

	c, f, g1, c2 := power(65.41), power(87.31), power(98.00), power(130.81)
	normalMel := buildRockSong(
		[4]float64{130.81, 174.61, 196.00, 130.81}, [4][3]float64{c, f, g1, c2},
		[4][]melodyNote{
			buildDriveMelody([]float64{523.25, 587.33, 659.25, 783.99}, 0.2),
			buildDriveMelody([]float64{698.46, 783.99, 880.00, 1046.5}, 0.2),
			buildDriveMelody([]float64{783.99, 880.00, 1046.5, 783.99}, 0.2),
			buildDriveMelody([]float64{523.25, 659.25, 783.99, 1046.5}, 0.2),
		}, 0.2)
	wildMel := buildRockSong(
		[4]float64{261.63, 349.23, 392.00, 261.63}, [4][3]float64{c, f, g1, c2},
		[4][]melodyNote{
			buildStabMelody([]float64{1046.5, 1174.66, 1318.5, 1567.98}, 0.1),
			buildStabMelody([]float64{1396.92, 1567.98, 1760.00, 2093.0}, 0.1),
			buildStabMelody([]float64{1567.98, 1760.00, 2093.0, 1567.98}, 0.1),
			buildStabMelody([]float64{1046.5, 1318.5, 1567.98, 2093.0}, 0.1),
		}, 0.1)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 180, Wave: "square",
		DroneFreqs: []float64{130.81, 196.00, 261.63},
		Attack: 0.02, Release: 0.1, CutoffMin: 1500, CutoffMax: 6000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "_", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 260, Wave: "square",
		DroneFreqs: []float64{130.81, 196.00, 261.63, 392.00},
		Attack: 0.02, Release: 0.03, CutoffMin: 2500, CutoffMax: 7500, DetuneCents: 8,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETRO  –  CRT scanline banger, A minor, 140/200
// ════════════════════════════════════════════════════════════════════
func soundsRetro() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{220.00, 261.63, 329.63, 440.00}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 440, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 330, 110, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 150, 75, 0.12, 0.09

	am, dm, em, am2 := minor(55.00), minor(73.42), minor(82.41), minor(110.00)
	normalMel := buildRockSong(
		[4]float64{110.00, 146.83, 164.81, 110.00}, [4][3]float64{am, dm, em, am2},
		[4][]melodyNote{
			buildDriveMelody([]float64{220.00, 293.66, 329.63, 440.00}, 0.268),
			buildDriveMelody([]float64{293.66, 349.23, 392.00, 440.00}, 0.268),
			buildDriveMelody([]float64{329.63, 392.00, 440.00, 493.88}, 0.268),
			buildDriveMelody([]float64{220.00, 261.63, 329.63, 440.00}, 0.268),
		}, 0.268)
	wildMel := buildRockSong(
		[4]float64{220.00, 146.83, 164.81, 220.00}, [4][3]float64{am, dm, em, am2},
		[4][]melodyNote{
			buildStabMelody([]float64{440.00, 587.33, 659.25, 880.00}, 0.134),
			buildStabMelody([]float64{587.33, 698.46, 783.99, 880.00}, 0.134),
			buildStabMelody([]float64{659.25, 783.99, 880.00, 987.77}, 0.134),
			buildStabMelody([]float64{440.00, 523.25, 659.25, 880.00}, 0.134),
		}, 0.134)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 140, Wave: "square",
		DroneFreqs: []float64{55.00, 110.00, 164.81},
		Attack: 0.02, Release: 0.18, CutoffMin: 1000, CutoffMax: 4500,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 200, Wave: "square",
		DroneFreqs: []float64{110.00, 220.00, 261.63, 329.63},
		Attack: 0.02, Release: 0.06, CutoffMin: 1800, CutoffMax: 6500, DetuneCents: 10,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// COSMOS  –  supernova impact, D minor, 150/220
// ════════════════════════════════════════════════════════════════════
func soundsCosmos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 369.99, 440.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 220, "triangle", 0.08, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 293.66, 587.33, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 146.83, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 150, 75, 0.12, 0.09

	dm, am, gm, dm2 := minor(73.42), power(110.00), power(98.00), minor(146.83)
	normalMel := buildRockSong(
		[4]float64{146.83, 220.00, 196.00, 146.83}, [4][3]float64{dm, am, gm, dm2},
		[4][]melodyNote{
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 493.88}, 0.25),
			buildDriveMelody([]float64{440.00, 493.88, 587.33, 293.66}, 0.25),
			buildDriveMelody([]float64{392.00, 440.00, 493.88, 349.23}, 0.25),
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 587.33}, 0.25),
		}, 0.25)
	wildMel := buildRockSong(
		[4]float64{293.66, 440.00, 392.00, 293.66}, [4][3]float64{dm, am, gm, dm2},
		[4][]melodyNote{
			buildStabMelody([]float64{587.33, 698.46, 880.00, 987.77}, 0.125),
			buildStabMelody([]float64{880.00, 987.77, 1174.66, 587.33}, 0.125),
			buildStabMelody([]float64{783.99, 880.00, 987.77, 698.46}, 0.125),
			buildStabMelody([]float64{587.33, 698.46, 880.00, 1174.66}, 0.125),
		}, 0.125)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 150, Wave: "square",
		DroneFreqs: []float64{73.42, 110.00, 146.83},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "_", "hat", "snare", "kick", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 220, Wave: "square",
		DroneFreqs: []float64{73.42, 146.83, 196.00, 293.66},
		Attack: 0.02, Release: 0.05, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 12,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETROFUTURE  –  atomic twilight stomp, E major, 140/210
// ════════════════════════════════════════════════════════════════════
func soundsRetrofuture() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{329.63, 415.30, 493.88, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 329.63, 659.25, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 493.88, 164.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	em, am, bm, e2 := power(82.41), power(55.00), power(61.74), power(164.81)
	normalMel := buildRockSong(
		[4]float64{164.81, 110.00, 123.47, 329.63}, [4][3]float64{em, am, bm, e2},
		[4][]melodyNote{
			buildDriveMelody([]float64{329.63, 392.00, 440.00, 493.88}, 0.268),
			buildDriveMelody([]float64{220.00, 261.63, 293.66, 329.63}, 0.268),
			buildDriveMelody([]float64{246.94, 293.66, 329.63, 369.99}, 0.268),
			buildDriveMelody([]float64{329.63, 392.00, 493.88, 659.25}, 0.268),
		}, 0.268)
	wildMel := buildRockSong(
		[4]float64{329.63, 220.00, 246.94, 329.63}, [4][3]float64{em, am, bm, e2},
		[4][]melodyNote{
			buildStabMelody([]float64{659.25, 783.99, 880.00, 987.77}, 0.134),
			buildStabMelody([]float64{440.00, 523.25, 587.33, 659.25}, 0.134),
			buildStabMelody([]float64{493.88, 587.33, 659.25, 739.99}, 0.134),
			buildStabMelody([]float64{659.25, 783.99, 987.77, 1318.5}, 0.134),
		}, 0.134)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 140, Wave: "square",
		DroneFreqs: []float64{82.41, 110.00, 164.81},
		Attack: 0.02, Release: 0.18, CutoffMin: 1000, CutoffMax: 4500,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 210, Wave: "square",
		DroneFreqs: []float64{82.41, 164.81, 220.00, 329.63},
		Attack: 0.02, Release: 0.06, CutoffMin: 1800, CutoffMax: 6500, DetuneCents: 10,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SPACEAGE  –  booster-ignition, G major, 130/200
// ════════════════════════════════════════════════════════════════════
func soundsSpaceage() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{392.00, 493.88, 587.33, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 261.63, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 392, 784, 0.17, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 587, 196, 0.19, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.12, 0.09

	gm, em, cm, gm2 := power(98.00), power(82.41), power(65.41), power(196.00)
	normalMel := buildRockSong(
		[4]float64{196.00, 164.81, 130.81, 392.00}, [4][3]float64{gm, em, cm, gm2},
		[4][]melodyNote{
			buildDriveMelody([]float64{392.00, 440.00, 493.88, 587.33}, 0.288),
			buildDriveMelody([]float64{329.63, 369.99, 392.00, 440.00}, 0.288),
			buildDriveMelody([]float64{261.63, 293.66, 329.63, 392.00}, 0.288),
			buildDriveMelody([]float64{392.00, 440.00, 587.33, 783.99}, 0.288),
		}, 0.288)
	wildMel := buildRockSong(
		[4]float64{392.00, 329.63, 261.63, 392.00}, [4][3]float64{gm, em, cm, gm2},
		[4][]melodyNote{
			buildStabMelody([]float64{783.99, 880.00, 987.77, 1174.66}, 0.144),
			buildStabMelody([]float64{659.25, 739.99, 783.99, 880.00}, 0.144),
			buildStabMelody([]float64{523.25, 587.33, 659.25, 783.99}, 0.144),
			buildStabMelody([]float64{783.99, 880.00, 1174.66, 1567.98}, 0.144),
		}, 0.144)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 130, Wave: "square",
		DroneFreqs: []float64{98.00, 130.81, 196.00},
		Attack: 0.02, Release: 0.18, CutoffMin: 1000, CutoffMax: 4500,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 200, Wave: "square",
		DroneFreqs: []float64{98.00, 196.00, 261.63, 392.00},
		Attack: 0.02, Release: 0.06, CutoffMin: 1800, CutoffMax: 6500, DetuneCents: 10,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TROPICAL  –  typhoon rush, F major, 170/260
// ════════════════════════════════════════════════════════════════════
func soundsTropical() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{349.23, 440.00, 523.25, 698.46}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 349.23, 698.46, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 174.61, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 170, 85, 0.12, 0.09

	fm, dm, gm, c := power(87.31), minor(73.42), power(98.00), power(65.41)
	normalMel := buildRockSong(
		[4]float64{174.61, 146.83, 196.00, 130.81}, [4][3]float64{fm, dm, gm, c},
		[4][]melodyNote{
			buildDriveMelody([]float64{349.23, 440.00, 523.25, 587.33}, 0.22),
			buildDriveMelody([]float64{293.66, 349.23, 392.00, 440.00}, 0.22),
			buildDriveMelody([]float64{392.00, 440.00, 493.88, 523.25}, 0.22),
			buildDriveMelody([]float64{261.63, 349.23, 392.00, 523.25}, 0.22),
		}, 0.22)
	wildMel := buildRockSong(
		[4]float64{349.23, 293.66, 196.00, 261.63}, [4][3]float64{fm, dm, gm, c},
		[4][]melodyNote{
			buildStabMelody([]float64{698.46, 880.00, 1046.5, 1174.66}, 0.11),
			buildStabMelody([]float64{587.33, 698.46, 783.99, 880.00}, 0.11),
			buildStabMelody([]float64{783.99, 880.00, 987.77, 1046.5}, 0.11),
			buildStabMelody([]float64{523.25, 698.46, 783.99, 1046.5}, 0.11),
		}, 0.11)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 170, Wave: "square",
		DroneFreqs: []float64{87.31, 174.61, 196.00},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 260, Wave: "square",
		DroneFreqs: []float64{87.31, 174.61, 196.00, 261.63},
		Attack: 0.02, Release: 0.04, CutoffMin: 2000, CutoffMax: 7500, DetuneCents: 12,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// NOIR  –  midnight chase, D minor, 130/190
// ════════════════════════════════════════════════════════════════════
func soundsNoir() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 349.23, 440.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 293.66, 587.33, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 146.83, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 130, 65, 0.12, 0.09

	dm, am, c, dm2 := minor(73.42), minor(110.00), power(65.41), minor(146.83)
	normalMel := buildRockSong(
		[4]float64{146.83, 220.00, 130.81, 146.83}, [4][3]float64{dm, am, c, dm2},
		[4][]melodyNote{
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 523.25}, 0.288),
			buildDriveMelody([]float64{440.00, 493.88, 587.33, 293.66}, 0.288),
			buildDriveMelody([]float64{261.63, 311.13, 392.00, 261.63}, 0.288),
			buildDriveMelody([]float64{293.66, 349.23, 440.00, 587.33}, 0.288),
		}, 0.288)
	wildMel := buildRockSong(
		[4]float64{293.66, 440.00, 261.63, 293.66}, [4][3]float64{dm, am, c, dm2},
		[4][]melodyNote{
			buildStabMelody([]float64{587.33, 698.46, 880.00, 1046.5}, 0.144),
			buildStabMelody([]float64{880.00, 987.77, 1174.66, 587.33}, 0.144),
			buildStabMelody([]float64{523.25, 622.25, 783.99, 523.25}, 0.144),
			buildStabMelody([]float64{587.33, 698.46, 880.00, 1174.66}, 0.144),
		}, 0.144)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.028, BPM: 130, Wave: "square",
		DroneFreqs: []float64{73.42, 110.00, 146.83},
		Attack: 0.02, Release: 0.18, CutoffMin: 1000, CutoffMax: 4500,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 190, Wave: "square",
		DroneFreqs: []float64{73.42, 146.83, 196.00, 293.66},
		Attack: 0.02, Release: 0.08, CutoffMin: 1800, CutoffMax: 6500, DetuneCents: 12,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// CATHEDRAL  –  bell-tower assault, C major, 120/180
// ════════════════════════════════════════════════════════════════════
func soundsCathedral() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 329.63, 392.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 261.63, 523.25, 0.2, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.22, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	cm, fm, gm, c2 := power(65.41), power(87.31), power(98.00), power(130.81)
	normalMel := buildRockSong(
		[4]float64{130.81, 174.61, 196.00, 130.81}, [4][3]float64{cm, fm, gm, c2},
		[4][]melodyNote{
			buildDriveMelody([]float64{523.25, 587.33, 659.25, 783.99}, 0.31),
			buildDriveMelody([]float64{698.46, 783.99, 880.00, 1046.5}, 0.31),
			buildDriveMelody([]float64{783.99, 880.00, 1046.5, 783.99}, 0.31),
			buildDriveMelody([]float64{523.25, 659.25, 783.99, 1046.5}, 0.31),
		}, 0.31)
	wildMel := buildRockSong(
		[4]float64{261.63, 349.23, 392.00, 261.63}, [4][3]float64{cm, fm, gm, c2},
		[4][]melodyNote{
			buildStabMelody([]float64{1046.5, 1174.66, 1318.5, 1567.98}, 0.155),
			buildStabMelody([]float64{1396.92, 1567.98, 1760.00, 2093.0}, 0.155),
			buildStabMelody([]float64{1567.98, 1760.00, 2093.0, 1567.98}, 0.155),
			buildStabMelody([]float64{1046.5, 1318.5, 1567.98, 2093.0}, 0.155),
		}, 0.155)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 120, Wave: "square",
		DroneFreqs: []float64{65.41, 130.81, 196.00},
		Attack: 0.02, Release: 0.2, CutoffMin: 800, CutoffMax: 4000,
		Drums: pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "_"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 180, Wave: "square",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 261.63},
		Attack: 0.02, Release: 0.08, CutoffMin: 1600, CutoffMax: 6000, DetuneCents: 10,
		Drums: pat("kick", "hat", "kick", "hat", "snare", "hat", "kick", "hat", "kick", "hat", "snare", "hat", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SURVEILLANCE  –  spy chase, E major, 150/200
// ════════════════════════════════════════════════════════════════════
func soundsSurveillance() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "square", 0.04, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 660, 1320, 0.12, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 990, 330, 0.14, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 280, 140, 0.09, 0.09

	em, cm, dm, bm := power(82.41), power(65.41), power(73.42), power(61.74)
	normalMel := buildRockSong(
		[4]float64{164.81, 130.81, 146.83, 123.47}, [4][3]float64{em, cm, dm, bm},
		[4][]melodyNote{
			buildDriveMelody([]float64{329.63, 369.99, 440.00, 493.88}, 0.25),
			buildDriveMelody([]float64{261.63, 311.13, 329.63, 392.00}, 0.25),
			buildDriveMelody([]float64{293.66, 349.23, 392.00, 440.00}, 0.25),
			buildDriveMelody([]float64{246.94, 293.66, 329.63, 369.99}, 0.25),
		}, 0.25)
	wildMel := buildRockSong(
		[4]float64{329.63, 261.63, 293.66, 246.94}, [4][3]float64{em, cm, dm, bm},
		[4][]melodyNote{
			buildStabMelody([]float64{659.25, 739.99, 880.00, 987.77}, 0.125),
			buildStabMelody([]float64{523.25, 622.25, 659.25, 783.99}, 0.125),
			buildStabMelody([]float64{587.33, 698.46, 783.99, 880.00}, 0.125),
			buildStabMelody([]float64{493.88, 587.33, 659.25, 739.99}, 0.125),
		}, 0.125)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 150, Wave: "square",
		DroneFreqs: []float64{82.41, 130.81, 164.81},
		Attack: 0.02, Release: 0.15, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "_", "hat", "snare", "kick", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 200, Wave: "square",
		DroneFreqs: []float64{164.81, 261.63, 329.63, 493.88},
		Attack: 0.02, Release: 0.04, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 18,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
		Melody: wildMel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BIOMECH  –  organic stomper, C# major, 160/240
// ════════════════════════════════════════════════════════════════════
func soundsBiomech() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{164.81, 246.94, 311.13, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 523.25, 0.18, 0.095
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	cs, am, emin, bm := power(69.30), power(55.00), power(82.41), power(61.74)
	normalMel := buildRockSong(
		[4]float64{138.59, 110.00, 164.81, 123.47}, [4][3]float64{cs, am, emin, bm},
		[4][]melodyNote{
			buildDriveMelody([]float64{277.18, 329.63, 369.99, 440.00}, 0.233),
			buildDriveMelody([]float64{220.00, 261.63, 293.66, 329.63}, 0.233),
			buildDriveMelody([]float64{329.63, 392.00, 440.00, 493.88}, 0.233),
			buildDriveMelody([]float64{246.94, 293.66, 329.63, 369.99}, 0.233),
		}, 0.233)
	wildMel := buildRockSong(
		[4]float64{277.18, 220.00, 329.63, 246.94}, [4][3]float64{cs, am, emin, bm},
		[4][]melodyNote{
			buildStabMelody([]float64{554.37, 659.25, 739.99, 880.00}, 0.117),
			buildStabMelody([]float64{440.00, 523.25, 587.33, 659.25}, 0.117),
			buildStabMelody([]float64{659.25, 783.99, 880.00, 987.77}, 0.117),
			buildStabMelody([]float64{493.88, 587.33, 659.25, 739.99}, 0.117),
		}, 0.117)

	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 160, Wave: "square",
		DroneFreqs: []float64{69.30, 110.00, 138.59},
		Attack: 0.02, Release: 0.12, CutoffMin: 1200, CutoffMax: 5000,
		Drums: pat("kick", "_", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "_", "hat", "snare", "kick", "_", "hat", "clap"),
		Melody: normalMel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 240, Wave: "square",
		DroneFreqs: []float64{69.30, 138.59, 220.00, 277.18},
		Attack: 0.02, Release: 0.04, CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 18,
		Drums: pat("kick", "kick", "hat", "kick", "snare", "kick", "hat", "kick", "kick", "kick", "hat", "snare", "kick", "kick", "hat", "clap"),
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
