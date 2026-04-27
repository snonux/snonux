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

// ── basic note + chord helpers ─────────────────────────────────────

// ns builds a note that rings for 'dur' but advances the sequencer by 'step',
// so multiple voices appended back-to-back can overlap when dur > step.
func ns(freq, dur, step float64) melodyNote {
	return melodyNote{Freq: freq, Dur: dur, Step: step}
}

var (
	intMajor3rd = math.Pow(2, 4.0/12)
	intMinor3rd = math.Pow(2, 3.0/12)
	intPerf5th  = math.Pow(2, 7.0/12)
	intMinor7th = math.Pow(2, 10.0/12)
)

func major(freq float64) [3]float64 {
	return [3]float64{freq, freq * intMajor3rd, freq * intPerf5th}
}
func minor(freq float64) [3]float64 {
	return [3]float64{freq, freq * intMinor3rd, freq * intPerf5th}
}

// ── voice builders ────────────────────────────────────────────────
// Each builder returns a slice of notes appended back-to-back into a melody.
// Because the engine advances time by note.Step (not note.Dur), notes with
// dur > step ring through the next entries — that is how distinct "voices"
// (bass, pad, lead) are stacked into one flat array.

// hook turns explicit (freq, beats) pairs into a melodic phrase. A pair with
// freq <= 0 OR beats <= 0 is treated as a rest, and its duration is folded
// into the preceding note's step so the next real note triggers later. We
// can't emit a zero-frequency rest entry because the JS engine substitutes
// 440 Hz for falsy freq values, which would turn rests into audible tones.
//   beat = seconds-per-quarter-note. dur multiplier of 1.0 = quarter note.
func hook(beat float64, pairs ...float64) []melodyNote {
	out := make([]melodyNote, 0, len(pairs)/2)
	leadIn := 0.0
	for i := 0; i+1 < len(pairs); i += 2 {
		f, b := pairs[i], pairs[i+1]
		if f <= 0 || b <= 0 {
			silence := math.Abs(b) * beat
			if len(out) == 0 {
				leadIn += silence
			} else {
				out[len(out)-1].Step += silence
			}
			continue
		}
		step := b * beat
		// 0.92 leaves a small gap so repeated notes re-trigger cleanly.
		out = append(out, ns(f, step*0.92, step))
	}
	// A leading rest in a looped phrase is equivalent to a trailing delay on
	// the last note (the loop wraps), so push it onto the final note's step.
	if leadIn > 0 && len(out) > 0 {
		out[len(out)-1].Step += leadIn
	}
	return out
}

// padHold lays a chord triad as a long sustained pad. The three notes are
// emitted back-to-back but each rings for almost the full duration, so they
// pile into a held chord.
func padHold(chord [3]float64, totalBeats, beat float64) []melodyNote {
	dur := totalBeats * beat
	return []melodyNote{
		ns(chord[0], dur*0.96, dur*0.34),
		ns(chord[1], dur*0.62, dur*0.33),
		ns(chord[2], dur*0.30, dur*0.33),
	}
}

// octaveBass: pumping synth-pop bassline — root, octave, root, fifth — over 4
// quarter-notes. Iconic Outrun / "Take On Me" feel.
func octaveBass(root, beat float64) []melodyNote {
	fifth := root * intPerf5th
	return []melodyNote{
		ns(root, beat*0.85, beat),
		ns(root*2, beat*0.85, beat),
		ns(root, beat*0.85, beat),
		ns(fifth, beat*0.85, beat),
	}
}

// palmMute: 8 chugging eighth-notes alternating root and fifth — heavy-metal
// palm-mute feel (think Doom Eternal / Megadeth).
func palmMute(root, beat float64) []melodyNote {
	out := make([]melodyNote, 0, 8)
	half := beat * 0.5
	fifth := root * intPerf5th
	for i := 0; i < 8; i++ {
		f := root
		if i == 4 || i == 6 {
			f = fifth
		}
		out = append(out, ns(f, half*0.45, half))
	}
	return out
}

// arpUpDown: chord up-and-down across 8 sixteenths (1-3-5-8-5-3-1-3) — Mario
// chiptune feel.
func arpUpDown(chord [3]float64, beat float64) []melodyNote {
	s := beat * 0.5
	return []melodyNote{
		ns(chord[0], s*0.85, s), ns(chord[1], s*0.85, s),
		ns(chord[2], s*0.85, s), ns(chord[0]*2, s*0.85, s),
		ns(chord[2], s*0.85, s), ns(chord[1], s*0.85, s),
		ns(chord[0], s*0.85, s), ns(chord[1], s*0.85, s),
	}
}

// walkBass: jazz quarter-note walking pattern root → 3rd → 5th → 3rd.
func walkBass(root, beat float64) []melodyNote {
	return []melodyNote{
		ns(root, beat*0.92, beat),
		ns(root*intMajor3rd, beat*0.92, beat),
		ns(root*intPerf5th, beat*0.92, beat),
		ns(root*intMajor3rd, beat*0.92, beat),
	}
}

// chromDesc: 8-note chromatic descent from 'top' — Bond / Pink Panther feel.
func chromDesc(top, beat float64) []melodyNote {
	out := make([]melodyNote, 0, 8)
	step := beat * 0.5
	for i := 0; i < 8; i++ {
		f := top * math.Pow(2, -float64(i)/12.0)
		out = append(out, ns(f, step*0.6, step))
	}
	return out
}

// fanfareStab: a punchy orchestral chord-hit held across 'beats' beats.
// Three voices land together because each rings for most of the duration.
func fanfareStab(chord [3]float64, beats, beat float64) []melodyNote {
	dur := beats * beat
	return []melodyNote{
		ns(chord[0], dur*0.9, dur*0.34),
		ns(chord[1], dur*0.7, dur*0.33),
		ns(chord[2], dur*0.5, dur*0.33),
	}
}

// reverseKick: hardstyle reverse-bass effect — a sub-octave thump followed by
// a short octave-up snap (the "wub" before the kick).
func reverseKick(root, beat float64) []melodyNote {
	return []melodyNote{
		ns(root/2, beat*0.85, beat*0.5),
		ns(root, beat*0.45, beat*0.5),
	}
}

// bossaBass: bossa nova bass — dotted-quarter root, eighth fifth, repeat.
func bossaBass(root, beat float64) []melodyNote {
	fifth := root * intPerf5th
	return []melodyNote{
		ns(root, beat*1.4, beat*1.5),
		ns(fifth, beat*0.45, beat*0.5),
		ns(root, beat*1.4, beat*1.5),
		ns(fifth, beat*0.45, beat*0.5),
	}
}

// dubSwell: very long sustained bass — held across the whole bar. Caller
// passes the actual played frequency (no implicit octave-down) so we don't
// drop below ~30 Hz, where most speakers reproduce nothing but clicks.
func dubSwell(freq, beats, beat float64) []melodyNote {
	return []melodyNote{ns(freq, beats*beat*0.95, beats*beat)}
}

// alienSlide: 4 sliding atonal arpeggio notes spaced by minor-3rds and tritones.
func alienSlide(root, beat float64) []melodyNote {
	tritone := math.Pow(2, 6.0/12)
	s := beat * 0.5
	return []melodyNote{
		ns(root, s*0.85, s),
		ns(root*intMinor3rd, s*0.85, s),
		ns(root*tritone, s*0.85, s),
		ns(root*intMinor7th, s*0.85, s),
	}
}

// concat flattens several voices into one note slice.
func concat(parts ...[]melodyNote) []melodyNote {
	var out []melodyNote
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

// ── drum patterns ─────────────────────────────────────────────────
// Each is a 16-step bar (one entry per 16th-note at the preset BPM).

func pat(names ...string) []string { return names }

var (
	// Synth-pop punch — kick on 1+4+8, claps on 16.
	drumPop = pat("kick", "_", "hat", "kick", "snare", "_", "hat", "_", "kick", "kick", "hat", "_", "snare", "_", "hat", "clap")
	// Four-on-the-floor dance — kick every quarter, hat offbeat.
	drumFour = pat("kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_", "kick", "_", "hat", "_")
	// Heavy rock — straight 4/4 with offbeat hats (Duke Nukem feel).
	drumRock = pat("kick", "_", "hat", "_", "snare", "_", "hat", "_", "kick", "_", "hat", "_", "snare", "_", "hat", "kick")
	// Industrial pulse — hammering double-kick, syncopated snare.
	drumPulse = pat("kick", "kick", "_", "hat", "snare", "_", "kick", "hat", "kick", "_", "_", "hat", "snare", "kick", "_", "hat")
	// Hardstyle — kick on every beat plus the offbeats (gabber).
	drumHardstyle = pat("kick", "_", "_", "kick", "kick", "_", "_", "kick", "kick", "_", "_", "kick", "kick", "kick", "_", "kick")
	// DnB amen-break — kick on 1+11, snare with ghosted "skips" on 4.
	drumDnb = pat("kick", "_", "hat", "kick", "snare", "hat", "_", "snare", "_", "kick", "hat", "_", "snare", "_", "hat", "kick")
	// Bossa nova clave (3-2 son) — Latin syncopation.
	drumBossa = pat("kick", "_", "hat", "_", "_", "_", "hat", "snare", "_", "kick", "hat", "_", "snare", "_", "hat", "_")
	// Boom-bap — hip-hop / film-noir groove.
	drumBoom = pat("kick", "_", "_", "hat", "snare", "_", "_", "hat", "_", "_", "kick", "hat", "snare", "_", "kick", "hat")
	// Dub — kick on 1, snare on 9, hi-hat on 16. Spacious.
	drumDub = pat("kick", "_", "_", "_", "_", "_", "_", "_", "snare", "_", "_", "_", "_", "_", "_", "hat")
	// Ambient — single soft kick + sparse clap. Mostly silence.
	drumAmbient = pat("kick", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "_", "clap", "_", "_", "_")
	// March — kick/snare alternate, no hat (cosmos military feel).
	drumMarch = pat("kick", "_", "_", "_", "snare", "_", "_", "_", "kick", "_", "_", "_", "snare", "_", "_", "_")
	// Bell tolls — only deep kicks, very sparse (cathedral).
	drumBell = pat("kick", "_", "_", "_", "_", "_", "_", "_", "kick", "_", "_", "_", "_", "_", "_", "_")
	// Latin funk — busy clave with conga-like accents.
	drumLatin = pat("kick", "hat", "_", "hat", "snare", "kick", "hat", "_", "kick", "hat", "_", "hat", "snare", "_", "hat", "kick")
	// Chiptune — bouncy NES drum-machine: kick+snare on every other 16th.
	drumChip = pat("kick", "snare", "kick", "snare", "kick", "snare", "kick", "snare", "kick", "snare", "kick", "snare", "kick", "snare", "kick", "snare")
	// Spy — sneaky offbeat snare with double-kick pickup.
	drumSpy = pat("kick", "_", "hat", "_", "_", "snare", "hat", "_", "kick", "kick", "hat", "_", "snare", "_", "hat", "snare")
	// Organic / biomech — irregular limb-like accents.
	drumOrganic = pat("kick", "_", "hat", "kick", "_", "snare", "kick", "_", "_", "kick", "hat", "_", "snare", "_", "kick", "hat")
	// Metal chug — double-kick galloping metal rhythm (plasma normal).
	drumChug = pat("kick", "kick", "_", "kick", "snare", "_", "kick", "_", "kick", "kick", "_", "kick", "snare", "kick", "kick", "kick")
	// Gated 80s electro — tight punchy clap on 5+13 (synthwave).
	drumElectro = pat("kick", "_", "_", "_", "clap", "_", "hat", "_", "kick", "kick", "_", "_", "clap", "_", "hat", "snare")
	// Wild rock — extra-busy version of drumRock for wild mode.
	drumRockWild = pat("kick", "kick", "hat", "kick", "snare", "_", "hat", "kick", "kick", "kick", "hat", "kick", "snare", "_", "hat", "kick")
)

// ── theme songs ───────────────────────────────────────────────────
//
// Each theme builds two ambient presets (normal + wild). Wild is the same
// musical character as normal but faster, denser, brighter — never a
// different song, so the listener always hears "the theme" intensifying.

// ════════════════════════════════════════════════════════════════════
// NEON  –  synth-pop banger, C major, Outrun pumping octave bass
// ════════════════════════════════════════════════════════════════════
func soundsNeon() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 330, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 523.25, 1046.5, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 880, 261.63, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.12, 0.1

	// Hook: sparkly major arpeggio that resolves down — C-E-G-C-A-G-E-C feel.
	beat := 0.4 // 150 BPM quarter
	mel := concat(
		octaveBass(130.81, beat), padHold(major(261.63), 4, beat),
		hook(beat, 1046.5, 0.5, 783.99, 0.5, 659.25, 0.5, 1046.5, 0.5,
			880.00, 0.5, 783.99, 0.5, 659.25, 0.5, 523.25, 0.5),
		octaveBass(174.61, beat), padHold(major(349.23), 4, beat),
		hook(beat, 1396.92, 0.5, 1046.5, 0.5, 880.00, 0.5, 1396.92, 0.5,
			1174.66, 0.5, 1046.5, 0.5, 880.00, 0.5, 698.46, 0.5),
	)
	wbeat := 0.273 // 220 BPM
	wmel := concat(
		octaveBass(261.63, wbeat), padHold(major(523.25), 4, wbeat),
		hook(wbeat, 1046.5, 0.25, 1318.5, 0.25, 1567.98, 0.25, 2093.0, 0.25,
			1567.98, 0.25, 1318.5, 0.25, 1046.5, 0.25, 783.99, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 150, Wave: "square",
		DroneFreqs: []float64{130.81, 196.00, 261.63}, Attack: 0.02, Release: 0.15,
		CutoffMin: 1200, CutoffMax: 5000, Drums: drumPop, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 220, Wave: "square",
		DroneFreqs: []float64{130.81, 261.63, 392.00, 523.25}, Attack: 0.02, Release: 0.08,
		CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 12, Drums: drumPop, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TERMINAL  –  industrial techno, E minor, hammering pulse
// ════════════════════════════════════════════════════════════════════
func soundsTerminal() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.11, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 800, "square", 0.045, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 600, 1200, 0.12, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 900, 400, 0.14, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1

	// Hook: dissonant rising tritone stab — Hotline Miami / industrial menace.
	beat := 0.375 // 160 BPM
	mel := concat(
		palmMute(82.41, beat), padHold(minor(164.81), 4, beat),
		hook(beat, 329.63, 0.5, 392.00, 0.5, 466.16, 0.5, 329.63, 0.5,
			466.16, 0.5, 622.25, 0.5, 466.16, 0.5, 329.63, 0.5),
		palmMute(110.00, beat), padHold(minor(220.00), 4, beat),
		hook(beat, 440.00, 0.5, 523.25, 0.5, 622.25, 0.5, 440.00, 0.5,
			622.25, 0.5, 783.99, 0.5, 622.25, 0.5, 440.00, 0.5),
	)
	wbeat := 0.353 // 170 BPM, but denser stabs
	wmel := concat(
		palmMute(164.81, wbeat),
		hook(wbeat, 659.25, 0.25, 783.99, 0.25, 932.33, 0.25, 659.25, 0.25,
			932.33, 0.25, 1244.51, 0.25, 932.33, 0.25, 659.25, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 160, Wave: "sawtooth",
		DroneFreqs: []float64{82.41, 123.47, 164.81}, Attack: 0.02, Release: 0.1,
		CutoffMin: 1000, CutoffMax: 5000, Drums: drumPulse, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 170, Wave: "sawtooth",
		DroneFreqs: []float64{82.41, 164.81, 246.94, 329.63}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2000, CutoffMax: 7000, DetuneCents: 18, Drums: drumPulse, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SYNTHWAVE  –  Carpenter Brut "Turbo Killer" anthem rock, A minor
// ════════════════════════════════════════════════════════════════════
func soundsSynthwave() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 440, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 110, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 150, 75, 0.14, 0.09

	// Driving Am-F-C-G with a sawtooth-stab lead riff. Action-rock not chill.
	beat := 0.4 // 150 BPM — punchy synth-anthem
	mel := concat(
		octaveBass(110.00, beat), padHold(minor(220.00), 4, beat),
		hook(beat, 880.00, 0.25, 1046.5, 0.25, 1318.5, 0.5, 1046.5, 0.25, 880.00, 0.25,
			659.25, 0.5, 880.00, 0.5, 1046.5, 1.0),
		octaveBass(87.31, beat), padHold(major(174.61), 4, beat),
		hook(beat, 698.46, 0.25, 880.00, 0.25, 1046.5, 0.5, 880.00, 0.25, 698.46, 0.25,
			523.25, 0.5, 698.46, 0.5, 880.00, 1.0),
		octaveBass(130.81, beat), padHold(major(261.63), 4, beat),
		hook(beat, 1046.5, 0.25, 1318.5, 0.25, 1567.98, 0.5, 1318.5, 0.25, 1046.5, 0.25,
			783.99, 0.5, 1046.5, 0.5, 1318.5, 1.0),
		octaveBass(98.00, beat), padHold(major(196.00), 4, beat),
		hook(beat, 783.99, 0.25, 987.77, 0.25, 1174.66, 0.5, 987.77, 0.25, 783.99, 0.25,
			587.33, 0.5, 783.99, 0.5, 987.77, 1.0),
	)
	wbeat := 0.316 // 190 BPM
	wmel := concat(
		octaveBass(220.00, wbeat),
		hook(wbeat, 1760.00, 0.25, 2093.0, 0.25, 2637.0, 0.25, 2093.0, 0.25,
			1760.00, 0.25, 1318.5, 0.25, 1760.00, 0.25, 2093.0, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "sawtooth",
		DroneFreqs: []float64{110.00, 174.61, 220.00}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumElectro, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 190, Wave: "sawtooth",
		DroneFreqs: []float64{220.00, 349.23, 440.00, 523.25}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 14, Drums: drumElectro, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// PLASMA  –  Doom Eternal djent, F# minor, drop-tuned palm-mute chug
// ════════════════════════════════════════════════════════════════════
func soundsPlasma() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{369.99, 493.88, 587.33, 739.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 370, 740, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 740, 246.94, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 120, 60, 0.1, 0.1

	// Djent: low palm-mute chugs against a screaming high lead.
	beat := 0.4 // 150 BPM — heavy not blistering
	mel := concat(
		palmMute(46.25, beat), padHold(minor(92.50), 4, beat),
		hook(beat, 369.99, 0.5, 0, 0.5, 369.99, 0.5, 440.00, 0.5,
			554.37, 0.5, 440.00, 0.5, 369.99, 1.0),
		palmMute(46.25, beat), padHold(minor(92.50), 4, beat),
		hook(beat, 369.99, 0.5, 440.00, 0.5, 554.37, 0.5, 739.99, 0.5,
			659.25, 0.5, 554.37, 0.5, 440.00, 0.5, 369.99, 0.5),
	)
	wbeat := 0.333 // 180 BPM thrash
	wmel := concat(
		palmMute(92.50, wbeat),
		hook(wbeat, 739.99, 0.25, 880.00, 0.25, 1108.73, 0.25, 1480.00, 0.25,
			1108.73, 0.25, 880.00, 0.25, 739.99, 0.25, 554.37, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "square",
		DroneFreqs: []float64{46.25, 92.50, 138.59}, Attack: 0.01, Release: 0.08,
		CutoffMin: 800, CutoffMax: 5500, Drums: drumChug, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "square",
		DroneFreqs: []float64{46.25, 92.50, 138.59, 277.18}, Attack: 0.01, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 8000, DetuneCents: 22, Drums: drumChug, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BRUTALIST  –  Akira-style minimalist techno, C minor, concrete pulse
// ════════════════════════════════════════════════════════════════════
func soundsBrutalist() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 329.63, 392.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.06, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 440, 880, 0.15, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 660, 220, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1

	// Hook: pulse-step bass with metallic stab leads. Cm-Ab-Gm-Fm.
	beat := 0.429 // 140 BPM driving
	mel := concat(
		palmMute(65.41, beat), padHold(minor(130.81), 4, beat),
		hook(beat, 523.25, 0.5, 0, 0.5, 622.25, 0.5, 0, 0.5,
			783.99, 1.0, 622.25, 1.0),
		palmMute(51.91, beat), padHold(major(103.83), 4, beat),
		hook(beat, 415.30, 0.5, 0, 0.5, 622.25, 0.5, 0, 0.5,
			830.61, 1.0, 622.25, 1.0),
	)
	wbeat := 0.375 // 160 BPM
	wmel := concat(
		palmMute(130.81, wbeat),
		hook(wbeat, 1046.5, 0.25, 1244.5, 0.25, 1567.98, 0.25, 1046.5, 0.25,
			1244.5, 0.25, 1046.5, 0.25, 830.61, 0.25, 622.25, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 140, Wave: "sawtooth",
		DroneFreqs: []float64{65.41, 98.00, 130.81}, Attack: 0.02, Release: 0.1,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumPulse, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.075, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{65.41, 98.00, 130.81, 196.00}, Attack: 0.02, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 20, Drums: drumPulse, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// VOLCANO  –  hardstyle / gabber, D minor, reverse-bass kicks
// ════════════════════════════════════════════════════════════════════
func soundsVolcano() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 369.99, 440.00, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 293.66, 587.33, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 146.83, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.12, 0.09

	// Reverse-bass + screech lead. Dm-Bb-F-C — classic hardstyle.
	beat := 0.4 // 150 BPM
	mel := concat(
		reverseKick(73.42, beat), reverseKick(73.42, beat),
		reverseKick(73.42, beat), reverseKick(73.42, beat),
		hook(beat, 587.33, 0.5, 698.46, 0.5, 880.00, 0.5, 1174.66, 0.5,
			880.00, 0.5, 698.46, 0.5, 587.33, 1.0),
		reverseKick(58.27, beat), reverseKick(58.27, beat),
		reverseKick(87.31, beat), reverseKick(65.41, beat),
		hook(beat, 466.16, 0.5, 587.33, 0.5, 698.46, 0.5, 880.00, 0.5,
			698.46, 0.5, 587.33, 0.5, 466.16, 1.0),
	)
	wbeat := 0.333 // 180 BPM
	wmel := concat(
		reverseKick(146.83, wbeat), reverseKick(146.83, wbeat),
		hook(wbeat, 1174.66, 0.25, 1396.92, 0.25, 1760.00, 0.25, 2349.32, 0.25,
			1760.00, 0.25, 1396.92, 0.25, 1174.66, 0.25, 880.00, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "sawtooth",
		DroneFreqs: []float64{73.42, 110.00, 146.83}, Attack: 0.01, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumHardstyle, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{73.42, 110.00, 146.83, 220.00}, Attack: 0.01, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 8000, DetuneCents: 18, Drums: drumHardstyle, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// AURORA  –  arctic anthem rock, G major, soaring stadium lead
// ════════════════════════════════════════════════════════════════════
func soundsAurora() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{392.00, 493.88, 587.33, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 261.63, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 390, 780, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523, 196, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 180, 90, 0.12, 0.09

	// Driving G-major anthem (G-Em-C-D) with a soaring fist-pump lead hook.
	beat := 0.462 // 130 BPM
	mel := concat(
		octaveBass(98.00, beat), padHold(major(196.00), 4, beat),
		hook(beat, 783.99, 0.5, 880.00, 0.5, 987.77, 0.5, 1174.66, 0.5,
			987.77, 0.5, 880.00, 1.0, 783.99, 0.5),
		octaveBass(82.41, beat), padHold(minor(164.81), 4, beat),
		hook(beat, 659.25, 0.5, 783.99, 0.5, 880.00, 0.5, 987.77, 0.5,
			880.00, 0.5, 783.99, 1.0, 659.25, 0.5),
		octaveBass(65.41, beat), padHold(major(130.81), 4, beat),
		hook(beat, 783.99, 0.5, 880.00, 0.5, 987.77, 0.5, 1174.66, 0.5,
			1318.5, 1.0, 1174.66, 1.0),
		octaveBass(73.42, beat), padHold(major(146.83), 4, beat),
		hook(beat, 880.00, 0.5, 987.77, 0.5, 1174.66, 0.5, 1318.5, 0.5,
			1174.66, 0.5, 987.77, 1.0, 880.00, 0.5),
	)
	wbeat := 0.353 // 170 BPM thrash anthem
	wmel := concat(
		octaveBass(196.00, wbeat),
		hook(wbeat, 1567.98, 0.25, 1760.00, 0.25, 1975.53, 0.25, 2349.32, 0.25,
			1975.53, 0.25, 1760.00, 0.25, 1567.98, 0.25, 1318.5, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 130, Wave: "triangle",
		DroneFreqs: []float64{98.00, 146.83, 196.00}, Attack: 0.02, Release: 0.1,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumPop, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 170, Wave: "sawtooth",
		DroneFreqs: []float64{98.00, 196.00, 293.66, 392.00}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2200, CutoffMax: 7500, DetuneCents: 12, Drums: drumPop, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// MATRIX  –  liquid drum-and-bass, C minor, sub-bass + amen-break stabs
// ════════════════════════════════════════════════════════════════════
func soundsMatrix() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 311.13, 392.00, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 440, 880, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 660, 220, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.1, 0.1

	// Reese-bass + glitchy stabs. Cm-Ab-Bb-Cm.
	beat := 0.353 // 170 BPM DnB
	mel := concat(
		dubSwell(41.20, 4, beat), padHold(minor(130.81), 4, beat),
		hook(beat, 783.99, 0.25, 0, 0.25, 783.99, 0.5, 932.33, 0.5,
			0, 0.5, 1046.5, 0.5, 932.33, 1.0),
		dubSwell(34.65, 4, beat), padHold(major(103.83), 4, beat),
		hook(beat, 622.25, 0.25, 0, 0.25, 622.25, 0.5, 783.99, 0.5,
			0, 0.5, 932.33, 0.5, 783.99, 1.0),
	)
	wbeat := 0.333 // 180 BPM neurofunk
	wmel := concat(
		palmMute(65.41, wbeat),
		hook(wbeat, 1046.5, 0.25, 1244.5, 0.25, 1567.98, 0.25, 1244.5, 0.25,
			1864.66, 0.25, 1567.98, 0.25, 1244.5, 0.25, 1046.5, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 170, Wave: "sawtooth",
		DroneFreqs: []float64{32.70, 65.41, 130.81}, Attack: 0.02, Release: 0.1,
		CutoffMin: 1000, CutoffMax: 5500, DetuneCents: 6, Drums: drumDnb, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 261.63}, Attack: 0.02, Release: 0.04,
		CutoffMin: 2200, CutoffMax: 8000, DetuneCents: 18, Drums: drumDnb, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// OCEAN  –  surf rock, E minor, Dick Dale tremolo guitar
// ════════════════════════════════════════════════════════════════════
func soundsOcean() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{329.63, 392.00, 493.88, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 220, "sine", 0.08, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 329.63, 659.25, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 493.88, 164.81, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 140, 70, 0.12, 0.09

	// "Misirlou"-style tremolo-picked surf riff over driving rock drums.
	// Em-D-C-B7 with a chromatic descending phrygian lead.
	beat := 0.429 // 140 BPM
	mel := concat(
		palmMute(82.41, beat), padHold(minor(164.81), 4, beat),
		hook(beat, 659.25, 0.25, 698.46, 0.25, 783.99, 0.25, 698.46, 0.25,
			659.25, 0.25, 622.25, 0.25, 587.33, 0.25, 523.25, 0.25,
			493.88, 0.5, 587.33, 0.5, 659.25, 1.0),
		palmMute(73.42, beat), padHold(major(146.83), 4, beat),
		hook(beat, 587.33, 0.25, 622.25, 0.25, 698.46, 0.25, 622.25, 0.25,
			587.33, 0.25, 523.25, 0.25, 493.88, 0.25, 440.00, 0.25,
			493.88, 0.5, 587.33, 0.5, 659.25, 1.0),
	)
	wbeat := 0.353 // 170 BPM
	wmel := concat(
		palmMute(164.81, wbeat),
		hook(wbeat, 1318.5, 0.125, 1396.92, 0.125, 1567.98, 0.125, 1396.92, 0.125,
			1318.5, 0.125, 1244.5, 0.125, 1174.66, 0.125, 1046.5, 0.125,
			987.77, 0.25, 1174.66, 0.25, 1318.5, 0.5),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 140, Wave: "triangle",
		DroneFreqs: []float64{82.41, 130.81, 164.81}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1500, CutoffMax: 6000, Drums: drumRock, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 170, Wave: "sawtooth",
		DroneFreqs: []float64{82.41, 164.81, 246.94, 329.63}, Attack: 0.02, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 14, Drums: drumRockWild, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// DOS  –  Duke Nukem 3D "Grabbag" lick, E minor, palm-muted metal
// ════════════════════════════════════════════════════════════════════
func soundsDos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 329.63, 392.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.05, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 440, 880, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 660, 220, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.1, 0.1

	// THE Lee Jackson Grabbag riff:  E5 E5 G5 A5 | E5 E5 G5 A5 Bb5 A5 G5 E5
	// then drop down  B4 D5 E5 ...  Iconic 1996 PC-DOS metal.
	beat := 0.462 // 130 BPM
	const (
		E2, E3, B3 = 82.41, 164.81, 246.94
		E4, G4, A4 = 329.63, 392.00, 440.00
		Bb4, B4    = 466.16, 493.88
	)
	riff := []float64{
		E4, 0.5, E4, 0.5, G4, 0.5, A4, 0.5, // bar 1a
		E4, 0.5, E4, 0.5, G4, 0.5, A4, 0.5, // bar 1b
		B4, 0.5, A4, 0.5, G4, 0.5, E4, 0.5, // bar 2 — the climbs-then-drops hook
		Bb4, 0.5, A4, 0.5, G4, 0.5, E4, 0.5,
	}
	mel := concat(
		palmMute(E2, beat), padHold([3]float64{E3, B3, E4}, 4, beat),
		hook(beat, riff...),
	)
	wbeat := 0.375 // 160 BPM thrash
	wmel := concat(
		palmMute(E3, wbeat),
		hook(wbeat, E4*2, 0.25, E4*2, 0.25, G4*2, 0.25, A4*2, 0.25,
			E4*2, 0.25, E4*2, 0.25, G4*2, 0.25, A4*2, 0.25,
			B4*2, 0.25, A4*2, 0.25, G4*2, 0.25, E4*2, 0.25,
			Bb4*2, 0.25, A4*2, 0.25, G4*2, 0.25, E4*2, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 130, Wave: "square",
		DroneFreqs: []float64{82.41, 164.81, 246.94}, Attack: 0.01, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumRock, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 160, Wave: "square",
		DroneFreqs: []float64{82.41, 164.81, 246.94, 329.63}, Attack: 0.01, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 8000, DetuneCents: 12, Drums: drumRockWild, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETRO  –  8-bit chiptune Mario bounce, A minor, NES arpeggios
// ════════════════════════════════════════════════════════════════════
func soundsRetro() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{220.00, 261.63, 329.63, 440.00}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 440, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 330, 110, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 150, 75, 0.12, 0.09

	// Bouncy Am-F-G-Em arp lead — Mega Man action-platformer pace.
	beat := 0.4 // 150 BPM
	mel := concat(
		octaveBass(110.00, beat), arpUpDown(minor(220.00), beat),
		hook(beat, 880.00, 0.25, 1046.5, 0.25, 1318.5, 0.25, 1760.00, 0.25,
			1318.5, 0.25, 1046.5, 0.25, 880.00, 0.25, 659.25, 0.25),
		octaveBass(87.31, beat), arpUpDown(major(174.61), beat),
		hook(beat, 698.46, 0.25, 880.00, 0.25, 1046.5, 0.25, 1396.92, 0.25,
			1046.5, 0.25, 880.00, 0.25, 698.46, 0.25, 523.25, 0.25),
		octaveBass(98.00, beat), arpUpDown(major(196.00), beat),
		hook(beat, 783.99, 0.25, 987.77, 0.25, 1174.66, 0.25, 1567.98, 0.25,
			1174.66, 0.25, 987.77, 0.25, 783.99, 0.25, 587.33, 0.25),
		octaveBass(82.41, beat), arpUpDown(minor(164.81), beat),
		hook(beat, 659.25, 0.25, 783.99, 0.25, 987.77, 0.25, 1318.5, 0.25,
			987.77, 0.25, 783.99, 0.25, 659.25, 0.25, 493.88, 0.25),
	)
	wbeat := 0.316 // 190 BPM
	wmel := concat(
		octaveBass(220.00, wbeat),
		hook(wbeat, 1760.00, 0.125, 2093.0, 0.125, 2637.0, 0.125, 3520.0, 0.125,
			2637.0, 0.125, 2093.0, 0.125, 1760.00, 0.125, 1318.5, 0.125,
			880.00, 0.125, 1046.5, 0.125, 1318.5, 0.125, 1760.00, 0.125,
			1318.5, 0.125, 1046.5, 0.125, 880.00, 0.125, 659.25, 0.125),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "square",
		DroneFreqs: []float64{55.00, 110.00, 220.00}, Attack: 0.01, Release: 0.05,
		CutoffMin: 1500, CutoffMax: 6500, Drums: drumChip, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 190, Wave: "square",
		DroneFreqs: []float64{110.00, 220.00, 329.63, 440.00}, Attack: 0.005, Release: 0.03,
		CutoffMin: 2500, CutoffMax: 8000, DetuneCents: 8, Drums: drumChip, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// COSMOS  –  Imperial March stomp, D minor, action fanfare
// ════════════════════════════════════════════════════════════════════
func soundsCosmos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 369.99, 440.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 220, "triangle", 0.08, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 293.66, 587.33, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 146.83, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 150, 75, 0.12, 0.09

	// Imperial-March-style stomping brass fanfare at action tempo.
	beat := 0.4 // 150 BPM driving
	mel := concat(
		palmMute(73.42, beat), padHold(minor(146.83), 4, beat),
		hook(beat, 587.33, 0.75, 587.33, 0.25, 587.33, 0.5, 880.00, 0.5,
			698.46, 0.25, 659.25, 0.25, 622.25, 0.25, 587.33, 0.25, 880.00, 1.0),
		palmMute(58.27, beat), padHold(major(116.54), 4, beat),
		hook(beat, 466.16, 0.75, 466.16, 0.25, 466.16, 0.5, 698.46, 0.5,
			554.37, 0.25, 523.25, 0.25, 493.88, 0.25, 466.16, 0.25, 698.46, 1.0),
		palmMute(65.41, beat), padHold(minor(130.81), 4, beat),
		hook(beat, 523.25, 0.75, 523.25, 0.25, 523.25, 0.5, 783.99, 0.5,
			622.25, 0.25, 587.33, 0.25, 554.37, 0.25, 523.25, 0.25, 783.99, 1.0),
	)
	wbeat := 0.333 // 180 BPM
	wmel := concat(
		palmMute(146.83, wbeat),
		hook(wbeat, 1174.66, 0.25, 1174.66, 0.25, 1396.92, 0.5, 1760.00, 0.5,
			2349.32, 0.5, 2093.0, 0.25, 1760.00, 0.25, 1396.92, 0.5, 1174.66, 0.5),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "sawtooth",
		DroneFreqs: []float64{73.42, 110.00, 146.83}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumMarch, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{73.42, 146.83, 220.00, 293.66}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 16, Drums: drumChug, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// RETROFUTURE  –  Tron / Daft Punk filter-square, E major
// ════════════════════════════════════════════════════════════════════
func soundsRetrofuture() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{329.63, 415.30, 493.88, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 329.63, 659.25, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 493.88, 164.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	// Filtered square bass + locked-grid analog lead.  E-A-B-E.
	beat := 0.4 // 150 BPM driving Daft Punk action-disco
	mel := concat(
		octaveBass(82.41, beat), padHold(major(164.81), 4, beat),
		hook(beat, 659.25, 0.5, 659.25, 0.5, 783.99, 0.5, 987.77, 0.5,
			987.77, 0.5, 783.99, 0.5, 659.25, 0.5, 493.88, 0.5),
		octaveBass(110.00, beat), padHold(major(220.00), 4, beat),
		hook(beat, 880.00, 0.5, 880.00, 0.5, 1046.5, 0.5, 1318.5, 0.5,
			1046.5, 0.5, 880.00, 0.5, 659.25, 0.5, 493.88, 0.5),
	)
	wbeat := 0.316 // 190 BPM
	wmel := concat(
		octaveBass(164.81, wbeat),
		hook(wbeat, 1318.5, 0.25, 1567.98, 0.25, 1864.66, 0.25, 1318.5, 0.25,
			1567.98, 0.25, 1864.66, 0.25, 1318.5, 0.25, 1046.5, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "square",
		DroneFreqs: []float64{82.41, 164.81, 246.94}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumFour, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 190, Wave: "sawtooth",
		DroneFreqs: []float64{82.41, 164.81, 329.63, 493.88}, Attack: 0.01, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 14, Drums: drumFour, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SPACEAGE  –  Star Trek action theme, G major, hero-fanfare drive
// ════════════════════════════════════════════════════════════════════
func soundsSpaceage() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{392.00, 493.88, 587.33, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 261.63, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 392, 784, 0.17, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 587, 196, 0.19, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.12, 0.09

	// Heroic ascending fanfare at action tempo — "to boldly go" energy.
	beat := 0.4 // 150 BPM
	mel := concat(
		octaveBass(98.00, beat), padHold(major(196.00), 4, beat),
		hook(beat, 783.99, 0.5, 987.77, 0.5, 1174.66, 0.5, 1567.98, 0.5,
			1318.5, 0.5, 1174.66, 0.5, 987.77, 0.5, 783.99, 0.5),
		octaveBass(82.41, beat), padHold(minor(164.81), 4, beat),
		hook(beat, 659.25, 0.5, 783.99, 0.5, 987.77, 0.5, 1318.5, 0.5,
			1174.66, 0.5, 987.77, 0.5, 783.99, 0.5, 659.25, 0.5),
		octaveBass(110.00, beat), padHold(minor(220.00), 4, beat),
		hook(beat, 880.00, 0.5, 1046.5, 0.5, 1318.5, 0.5, 1760.00, 0.5,
			1396.92, 0.5, 1318.5, 0.5, 1046.5, 0.5, 880.00, 0.5),
		octaveBass(73.42, beat), padHold(major(146.83), 4, beat),
		hook(beat, 587.33, 0.5, 880.00, 0.5, 1108.73, 0.5, 1318.5, 0.5,
			1567.98, 0.5, 1318.5, 0.5, 1108.73, 0.5, 880.00, 0.5),
	)
	wbeat := 0.333 // 180 BPM
	wmel := concat(
		octaveBass(196.00, wbeat),
		hook(wbeat, 1567.98, 0.25, 1975.53, 0.25, 2349.32, 0.25, 3135.96, 0.25,
			2637.0, 0.25, 2349.32, 0.25, 1975.53, 0.25, 1567.98, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "triangle",
		DroneFreqs: []float64{98.00, 130.81, 196.00}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumPop, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{98.00, 196.00, 293.66, 392.00}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2200, CutoffMax: 7500, DetuneCents: 12, Drums: drumPop, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// TROPICALE  –  Latin afro-funk, F major, conga clave + brass stabs
// ════════════════════════════════════════════════════════════════════
func soundsTropical() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{349.23, 440.00, 523.25, 698.46}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 349.23, 698.46, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 174.61, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 170, 85, 0.12, 0.09

	// Tito Puente brass-stab energy, F-Bb-C-F.
	beat := 0.5 // 120 BPM
	mel := concat(
		bossaBass(87.31, beat), padHold(major(174.61), 4, beat),
		hook(beat, 698.46, 0.5, 0, 0.25, 880.00, 0.25, 1046.5, 0.5, 880.00, 0.5,
			698.46, 0.5, 880.00, 0.5, 1046.5, 1.0),
		bossaBass(116.54, beat), padHold(major(233.08), 4, beat),
		hook(beat, 932.33, 0.5, 0, 0.25, 1174.66, 0.25, 1396.92, 0.5, 1174.66, 0.5,
			932.33, 0.5, 1174.66, 0.5, 1396.92, 1.0),
		bossaBass(130.81, beat), padHold(major(261.63), 4, beat),
		hook(beat, 1046.5, 0.5, 0, 0.25, 1318.5, 0.25, 1567.98, 0.5, 1318.5, 0.5,
			1046.5, 0.5, 1318.5, 0.5, 1567.98, 1.0),
	)
	wbeat := 0.353 // 170 BPM
	wmel := concat(
		bossaBass(174.61, wbeat),
		hook(wbeat, 1396.92, 0.25, 1760.00, 0.25, 2093.0, 0.25, 1760.00, 0.25,
			2349.32, 0.25, 2093.0, 0.25, 1760.00, 0.25, 1396.92, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 120, Wave: "sawtooth",
		DroneFreqs: []float64{87.31, 174.61, 261.63}, Attack: 0.02, Release: 0.15,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumLatin, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 170, Wave: "sawtooth",
		DroneFreqs: []float64{87.31, 174.61, 261.63, 349.23}, Attack: 0.02, Release: 0.06,
		CutoffMin: 2200, CutoffMax: 7500, DetuneCents: 12, Drums: drumLatin, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// NOIR  –  rockabilly chase, D minor, slap-bass + chromatic guitar
// ════════════════════════════════════════════════════════════════════
func soundsNoir() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 349.23, 440.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 293.66, 587.33, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 146.83, 0.18, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 130, 65, 0.12, 0.09

	// Rockabilly chase — Stray Cats slap-bass walking + Pulp-Fiction-style
	// surf-rock chromatic guitar lead at action tempo.  Dm-Gm-A7-Dm.
	beat := 0.4 // 150 BPM rockabilly
	mel := concat(
		walkBass(73.42, beat), padHold(minor(146.83), 4, beat),
		chromDesc(880.00, beat),
		walkBass(98.00, beat), padHold(minor(196.00), 4, beat),
		chromDesc(739.99, beat),
		walkBass(110.00, beat), padHold(major(220.00), 4, beat),
		chromDesc(659.25, beat),
		walkBass(73.42, beat), padHold(minor(146.83), 4, beat),
		chromDesc(587.33, beat),
	)
	wbeat := 0.316 // 190 BPM
	wmel := concat(
		walkBass(146.83, wbeat),
		chromDesc(1318.5, wbeat),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "triangle",
		DroneFreqs: []float64{73.42, 110.00, 146.83}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumBoom, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 190, Wave: "sawtooth",
		DroneFreqs: []float64{73.42, 146.83, 196.00, 293.66}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2200, CutoffMax: 7500, DetuneCents: 14, Drums: drumBoom, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// CATHEDRAL  –  gothic metal organ, C minor, Type-O-Negative drive
// ════════════════════════════════════════════════════════════════════
func soundsCathedral() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{261.63, 329.63, 392.00, 523.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 261.63, 523.25, 0.2, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.22, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	// Gothic-metal organ pedal + driven Toccata-style lead. Cm-Gm-Ab-Bb.
	beat := 0.4 // 150 BPM driving gothic
	mel := concat(
		palmMute(65.41, beat), padHold(minor(130.81), 4, beat),
		hook(beat, 523.25, 0.5, 622.25, 0.5, 783.99, 0.5, 622.25, 0.5,
			523.25, 0.5, 466.16, 0.5, 415.30, 0.5, 391.99, 0.5),
		palmMute(98.00, beat), padHold(minor(196.00), 4, beat),
		hook(beat, 783.99, 0.5, 932.33, 0.5, 1174.66, 0.5, 932.33, 0.5,
			783.99, 0.5, 698.46, 0.5, 622.25, 0.5, 587.33, 0.5),
		palmMute(51.91, beat), padHold(major(103.83), 4, beat),
		hook(beat, 415.30, 0.5, 523.25, 0.5, 622.25, 0.5, 523.25, 0.5,
			415.30, 0.5, 391.99, 0.5, 369.99, 0.5, 311.13, 0.5),
		palmMute(58.27, beat), padHold(major(116.54), 4, beat),
		hook(beat, 466.16, 0.5, 587.33, 0.5, 698.46, 0.5, 587.33, 0.5,
			466.16, 0.5, 415.30, 0.5, 391.99, 0.5, 349.23, 0.5),
	)
	wbeat := 0.333 // 180 BPM
	wmel := concat(
		palmMute(130.81, wbeat),
		hook(wbeat, 1046.5, 0.25, 1244.5, 0.25, 1567.98, 0.25, 1244.5, 0.25,
			1046.5, 0.25, 932.33, 0.25, 783.99, 0.25, 622.25, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "triangle",
		DroneFreqs: []float64{65.41, 98.00, 130.81}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumChug, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.07, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{65.41, 130.81, 196.00, 261.63}, Attack: 0.02, Release: 0.05,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 14, Drums: drumChug, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// SURVEILLANCE  –  John Barry spy theme, E minor, chromatic sneak
// ════════════════════════════════════════════════════════════════════
func soundsSurveillance() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "square", 0.04, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 660, 1320, 0.12, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 990, 330, 0.14, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 280, 140, 0.09, 0.09

	// Bond "dum-da-da-dum" + chromatic descend lead at chase tempo.
	beat := 0.4 // 150 BPM action
	mel := concat(
		walkBass(82.41, beat), padHold(minor(164.81), 4, beat),
		hook(beat, 329.63, 0.5, 392.00, 0.5, 493.88, 0.5, 392.00, 0.5,
			659.25, 1.0, 622.25, 1.0),
		walkBass(110.00, beat), padHold(major(220.00), 4, beat),
		chromDesc(987.77, beat),
		walkBass(82.41, beat), padHold(minor(164.81), 4, beat),
		hook(beat, 329.63, 0.5, 392.00, 0.5, 493.88, 0.5, 392.00, 0.5,
			740.00, 1.0, 698.46, 1.0),
		walkBass(123.47, beat), padHold(major(246.94), 4, beat),
		chromDesc(880.00, beat),
	)
	wbeat := 0.316 // 190 BPM chase
	wmel := concat(
		walkBass(164.81, wbeat),
		chromDesc(1318.5, wbeat),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "triangle",
		DroneFreqs: []float64{82.41, 123.47, 164.81}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumSpy, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.075, BPM: 190, Wave: "sawtooth",
		DroneFreqs: []float64{164.81, 246.94, 329.63, 493.88}, Attack: 0.02, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 18, Drums: drumSpy, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// BIOMECH  –  alien organic, C# minor, atonal sliding arpeggios
// ════════════════════════════════════════════════════════════════════
func soundsBiomech() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{164.81, 246.94, 311.13, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 523.25, 0.18, 0.095
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09

	// Alien arpeggios with tritone slide — Giger-style organic action.
	beat := 0.4 // 150 BPM
	mel := concat(
		alienSlide(69.30, beat), padHold(minor(138.59), 4, beat),
		alienSlide(138.59, beat), alienSlide(277.18, beat),
		alienSlide(82.41, beat), padHold(minor(164.81), 4, beat),
		alienSlide(164.81, beat), alienSlide(329.63, beat),
		alienSlide(58.27, beat), padHold(major(116.54), 4, beat),
		alienSlide(116.54, beat), alienSlide(233.08, beat),
		alienSlide(73.42, beat), padHold(minor(146.83), 4, beat),
		alienSlide(146.83, beat), alienSlide(293.66, beat),
	)
	wbeat := 0.316 // 190 BPM
	wmel := concat(
		alienSlide(138.59, wbeat), alienSlide(277.18, wbeat),
		alienSlide(554.37, wbeat), alienSlide(1108.73, wbeat),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.032, BPM: 150, Wave: "sawtooth",
		DroneFreqs: []float64{69.30, 110.00, 138.59}, Attack: 0.02, Release: 0.08,
		CutoffMin: 1200, CutoffMax: 5500, Drums: drumOrganic, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.075, BPM: 190, Wave: "sawtooth",
		DroneFreqs: []float64{69.30, 138.59, 220.00, 277.18}, Attack: 0.02, Release: 0.04,
		CutoffMin: 2400, CutoffMax: 7500, DetuneCents: 22, Drums: drumOrganic, Melody: wmel,
	}
	return s
}

// ════════════════════════════════════════════════════════════════════
// NUKEM  –  action-hero hard rock, E minor, Grabbag-inspired power riffs
// ════════════════════════════════════════════════════════════════════
func soundsNukem() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{164.81, 329.63, 440.00, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.1, "sawtooth"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "square", 0.05, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sawtooth", 329.63, 659.25, 0.14, 0.11
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 659.25, 164.81, 0.16, 0.1
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sawtooth", 220, 110, 0.1, 0.1

	// THE Grabbag riff (Lee Jackson) — E5 E5 G5 A5 | E5 E5 G5 A5 Bb5 A5 G5 E5
	// Heavier sawtooth rendition at 140 BPM with palm-muted bass chugging.
	beat := 0.429 // 140 BPM hard rock
	const (
		E2, E3, B3 = 82.41, 164.81, 246.94
		E4, G4, A4 = 329.63, 392.00, 440.00
		Bb4, B4    = 466.16, 493.88
	)
	riff := []float64{
		E4, 0.5, E4, 0.5, G4, 0.5, A4, 0.5, // bar 1a
		E4, 0.5, E4, 0.5, G4, 0.5, A4, 0.5, // bar 1b
		B4, 0.5, A4, 0.5, G4, 0.5, E4, 0.5, // bar 2 — climbs then drops
		Bb4, 0.5, A4, 0.5, G4, 0.5, E4, 0.5,
	}
	mel := concat(
		palmMute(E2, beat), padHold([3]float64{E3, B3, E4}, 4, beat),
		hook(beat, riff...),
	)
	wbeat := 0.333 // 180 BPM thrash
	wmel := concat(
		palmMute(E3, wbeat),
		hook(wbeat, E4*2, 0.25, E4*2, 0.25, G4*2, 0.25, A4*2, 0.25,
			E4*2, 0.25, E4*2, 0.25, G4*2, 0.25, A4*2, 0.25,
			B4*2, 0.25, A4*2, 0.25, G4*2, 0.25, E4*2, 0.25,
			Bb4*2, 0.25, A4*2, 0.25, G4*2, 0.25, E4*2, 0.25),
	)
	s.Ambient.Normal = ambientPreset{
		Gain: 0.035, BPM: 140, Wave: "sawtooth",
		DroneFreqs: []float64{82.41, 164.81, 246.94}, Attack: 0.01, Release: 0.08,
		CutoffMin: 1400, CutoffMax: 6000, Drums: drumRock, Melody: mel,
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.075, BPM: 180, Wave: "sawtooth",
		DroneFreqs: []float64{82.41, 164.81, 329.63, 440.00}, Attack: 0.01, Release: 0.04,
		CutoffMin: 2800, CutoffMax: 8500, DetuneCents: 14, Drums: drumRockWild, Melody: wmel,
	}
	return s
}

// ── public entry points ───────────────────────────────────────────

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
	"nukem":        soundsNukem(),
}
