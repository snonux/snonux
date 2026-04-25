package generator

import (
	"encoding/json"
	"html/template"
)

// ambientPreset describes a generative ambient background layer for a theme.
// All fields are optional at runtime; missing values should be treated as
// silence or safe defaults by the consumer.
type ambientPreset struct {
	BPM           float64   `json:"bpm,omitempty"`
	PulseInterval float64   `json:"pulseInterval,omitempty"`
	Gain          float64   `json:"gain,omitempty"`
	Wave          string    `json:"wave,omitempty"`
	DroneFreqs    []float64 `json:"droneFreqs,omitempty"`
	PulseFreqs    []float64 `json:"pulseFreqs,omitempty"`
	CutoffMin     float64   `json:"cutoffMin,omitempty"`
	CutoffMax     float64   `json:"cutoffMax,omitempty"`
	NoiseGain     float64   `json:"noiseGain,omitempty"`
	Attack        float64   `json:"attack,omitempty"`
	Release       float64   `json:"release,omitempty"`
	DetuneCents   float64   `json:"detuneCents,omitempty"`
	Rhythm        []float64 `json:"rhythm,omitempty"`
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

// ambient is a concise helper for building an ambientPreset with sensible defaults.
func ambient(gain, bpm float64, wave string, drone, pulse []float64) ambientPreset {
	return ambientPreset{
		Gain:       gain,
		BPM:        bpm,
		Wave:       wave,
		DroneFreqs: drone,
		PulseFreqs: pulse,
		Attack:     0.4,
		Release:    0.8,
	}
}

// themeSoundPresets maps CLI theme names to synth parameters (see themes.go registry).
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

func soundsNeon() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 330, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 523.25, 1046.5, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 880, 261.63, 0.16, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 180, 90, 0.12, 0.1
	s.Ambient.Normal = ambient(0.08, 80, "square", []float64{523.25, 659.25}, []float64{1046.5})
	s.Ambient.Wild = ambient(0.12, 140, "square", []float64{261.63, 523.25}, []float64{880, 1760})
	return s
}

func soundsTerminal() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.11, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 800, "square", 0.045, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 600, 1200, 0.12, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 900, 400, 0.14, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 200, 100, 0.1, 0.1
	s.Ambient.Normal = ambient(0.09, 72, "square", []float64{400, 600}, []float64{800})
	s.Ambient.Wild = ambient(0.13, 144, "square", []float64{200, 400}, []float64{600, 1200})
	return s
}

func soundsSynthwave() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 440, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 110, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 150, 75, 0.14, 0.09
	s.Ambient.Normal = ambient(0.07, 60, "sine", []float64{196, 293.66}, []float64{587.33})
	s.Ambient.Wild = ambient(0.11, 110, "triangle", []float64{130.81, 196}, []float64{440, 880})
	return s
}

func soundsPlasma() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{311.13, 415.3, 466.16, 622.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "sine", 0.085, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 349.23, 698.46, 0.15, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 174.61, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 200, 100, 0.13, 0.09
	s.Ambient.Normal = ambient(0.08, 90, "triangle", []float64{311.13, 466.16}, []float64{622.25})
	s.Ambient.Wild = ambient(0.12, 160, "square", []float64{155.56, 311.13}, []float64{622.25, 1244.5})
	return s
}

func soundsBrutalist() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{100, 150, 200, 120}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.14, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 120, "square", 0.07, 0.13
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 200, 400, 0.12, 0.11
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 400, 100, 0.14, 0.1
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 120, 60, 0.1, 0.12
	s.Ambient.Normal = ambient(0.10, 65, "square", []float64{100, 150}, []float64{200})
	s.Ambient.Wild = ambient(0.14, 130, "square", []float64{50, 100}, []float64{150, 300})
	return s
}

func soundsVolcano() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 180, "sine", 0.09, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 261.63, 523.25, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 98, 0.17, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 160, 80, 0.13, 0.1
	s.Ambient.Normal = ambient(0.08, 50, "sine", []float64{98, 146.83}, []float64{196})
	s.Ambient.Wild = ambient(0.12, 100, "triangle", []float64{65.41, 98}, []float64{196, 392})
	return s
}

func soundsAurora() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{659.25, 880, 987.77, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "sine", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 523.25, 880, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 704, 352, 0.18, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 220, 110, 0.15, 0.08
	s.Ambient.Normal = ambient(0.06, 55, "sine", []float64{440, 880}, []float64{1320})
	s.Ambient.Wild = ambient(0.10, 110, "triangle", []float64{220, 440}, []float64{880, 1760})
	return s
}

func soundsMatrix() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 587.33, 659.25, 698.46}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.05, 0.09, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 523.25, "square", 0.045, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 880, 1318.5, 0.11, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 880, 330, 0.13, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 260, 130, 0.1, 0.09
	s.Ambient.Normal = ambient(0.08, 85, "square", []float64{523.25, 659.25}, []float64{880})
	s.Ambient.Wild = ambient(0.12, 170, "square", []float64{261.63, 523.25}, []float64{880, 1760})
	return s
}

func soundsOcean() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 196, 220, 246.94}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.095, 0.095, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 349.23, 523.25, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 246.94, 0.18, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 140, 70, 0.15, 0.08
	s.Ambient.Normal = ambient(0.06, 45, "sine", []float64{174.61, 220}, []float64{349.23})
	s.Ambient.Wild = ambient(0.10, 90, "triangle", []float64{130.81, 174.61}, []float64{349.23, 698.46})
	return s
}

func soundsDos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{800, 1000}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.12, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1000, "square", 0.03, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 400, 800, 0.1, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 800, 200, 0.1, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 300, 150, 0.08, 0.1
	s.Ambient.Normal = ambient(0.09, 100, "square", []float64{400, 800}, []float64{1200})
	s.Ambient.Wild = ambient(0.13, 200, "square", []float64{200, 400}, []float64{800, 1600})
	return s
}

func soundsRetro() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{1046.5, 1318.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1200, "square", 0.04, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 800, 1600, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 1600, 400, 0.15, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 400, 200, 0.1, 0.1
	s.Ambient.Normal = ambient(0.08, 95, "square", []float64{800, 1200}, []float64{1600})
	s.Ambient.Wild = ambient(0.12, 190, "square", []float64{400, 800}, []float64{1600, 3200})
	return s
}

func soundsCosmos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{220, 277.18, 329.63, 392}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "sine", 0.09, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 392, 587.33, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 587.33, 196, 0.2, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 170, 85, 0.16, 0.08
	s.Ambient.Normal = ambient(0.06, 40, "sine", []float64{220, 440}, []float64{660})
	s.Ambient.Wild = ambient(0.10, 80, "triangle", []float64{110, 220}, []float64{440, 880})
	return s
}

func soundsRetrofuture() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 329.63, 440}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "triangle", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 330, 523.25, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 165, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 190, 95, 0.14, 0.09
	s.Ambient.Normal = ambient(0.07, 70, "triangle", []float64{220, 330}, []float64{523.25})
	s.Ambient.Wild = ambient(0.11, 140, "square", []float64{110, 220}, []float64{523.25, 1046.5})
	return s
}

// soundsSpaceage returns synth parameters for the Space Age theme.
// Clean teal-toned tones evoke retro-futuristic space-age electronics.
func soundsSpaceage() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25, 880}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 554.37, "triangle", 0.075, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 659.25, 330, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 240, 120, 0.13, 0.08
	s.Ambient.Normal = ambient(0.06, 75, "sine", []float64{440, 880}, []float64{1320})
	s.Ambient.Wild = ambient(0.10, 150, "triangle", []float64{220, 440}, []float64{880, 1760})
	return s
}

// soundsTropical returns synth parameters for the Tropical Beach theme.
// Warm sine arpeggio across a C-major pentatonic — steel drum impression on splash.
// Nav uses a breathy mid-freq sine "bird tone"; open/close sweep like breaking surf.
func soundsTropical() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.08, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "sine", 0.07, 0.07
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.08
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 660, 330, 0.17, 0.075
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 200, 100, 0.12, 0.07
	s.Ambient.Normal = ambient(0.06, 80, "sine", []float64{523.25, 659.25}, []float64{1046.5})
	s.Ambient.Wild = ambient(0.10, 160, "triangle", []float64{261.63, 523.25}, []float64{1046.5, 2093})
	return s
}

func soundsNoir() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 220, 261.63}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.08, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 392, 0.18, 0.085
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "triangle", 330, 165, 0.2, 0.08
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 130, 65, 0.14, 0.08
	s.Ambient.Normal = ambient(0.06, 50, "sine", []float64{174.61, 220}, []float64{330})
	s.Ambient.Wild = ambient(0.10, 100, "triangle", []float64{130.81, 174.61}, []float64{330, 660})
	return s
}

func soundsCathedral() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{293.66, 392, 587.33}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.14, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 293.66, "triangle", 0.11, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 392, 783.99, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 196, 0.22, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 180, 90, 0.18, 0.08
	s.Ambient.Normal = ambient(0.07, 40, "sine", []float64{293.66, 440}, []float64{587.33})
	s.Ambient.Wild = ambient(0.11, 80, "triangle", []float64{146.83, 293.66}, []float64{587.33, 1174.66})
	return s
}

func soundsSurveillance() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.06, 0.08, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "square", 0.04, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 660, 1320, 0.12, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 990, 330, 0.14, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "square", 280, 140, 0.09, 0.09
	s.Ambient.Normal = ambient(0.07, 70, "square", []float64{440, 660}, []float64{880})
	s.Ambient.Wild = ambient(0.11, 140, "square", []float64{220, 440}, []float64{880, 1760})
	return s
}

func soundsBiomech() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{164.81, 246.94, 311.13, 466.16}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 233.08, "sine", 0.085, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 220, 523.25, 0.18, 0.095
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 130.81, 0.2, 0.09
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "triangle", 160, 80, 0.14, 0.09
	s.Ambient.Normal = ambient(0.07, 65, "triangle", []float64{164.81, 246.94}, []float64{440})
	s.Ambient.Wild = ambient(0.11, 130, "square", []float64{82.41, 164.81}, []float64{440, 880})
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
