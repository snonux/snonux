package generator

import (
	"encoding/json"
	"html/template"
)

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
}

// themeSoundPresets maps CLI theme names to synth parameters (see themes.go registry).
var themeSoundPresets = map[string]themeSounds{
	"neon":      soundsNeon(),
	"terminal":  soundsTerminal(),
	"synthwave": soundsSynthwave(),
	"plasma":    soundsPlasma(),
	"brutalist": soundsBrutalist(),
	"volcano":   soundsVolcano(),
	"aurora":    soundsAurora(),
	"matrix":    soundsMatrix(),
	"ocean":     soundsOcean(),
	"dos":       soundsDos(),
	"retro":     soundsRetro(),
	"cosmos":    soundsCosmos(),
}

func soundsNeon() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.055, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 330, "square", 0.055, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 523.25, 1046.5, 0.13, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 880, 261.63, 0.16, 0.09
	return s
}

func soundsTerminal() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.11, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 800, "square", 0.045, 0.12
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 600, 1200, 0.12, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 900, 400, 0.14, 0.09
	return s
}

func soundsSynthwave() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.1, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 164.81, "triangle", 0.09, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 220, 440, 0.18, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 440, 110, 0.17, 0.09
	return s
}

func soundsPlasma() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{311.13, 415.3, 466.16, 622.25}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.095, "triangle"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 246.94, "sine", 0.085, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 349.23, 698.46, 0.15, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 523.25, 174.61, 0.17, 0.09
	return s
}

func soundsBrutalist() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{100, 150, 200, 120}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.14, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 120, "square", 0.07, 0.13
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 200, 400, 0.12, 0.11
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 400, 100, 0.14, 0.1
	return s
}

func soundsVolcano() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{196, 246.94, 293.66, 349.23}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.1, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 180, "sine", 0.09, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 261.63, 523.25, 0.16, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 392, 98, 0.17, 0.09
	return s
}

func soundsAurora() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{659.25, 880, 987.77, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 440, "sine", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 523.25, 880, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 704, 352, 0.18, 0.085
	return s
}

func soundsMatrix() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 587.33, 659.25, 698.46}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.05, 0.09, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 523.25, "square", 0.045, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 880, 1318.5, 0.11, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 880, 330, 0.13, 0.09
	return s
}

func soundsOcean() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{174.61, 196, 220, 246.94}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.095, 0.095, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 196, "triangle", 0.1, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 349.23, 523.25, 0.2, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 415.3, 246.94, 0.18, 0.085
	return s
}

func soundsDos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{800, 1000}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.08, 0.12, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1000, "square", 0.03, 0.1
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 400, 800, 0.1, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 800, 200, 0.1, 0.09
	return s
}

func soundsRetro() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{1046.5, 1318.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.12, 0.1, "square"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 1200, "square", 0.04, 0.11
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "square", 800, 1600, 0.14, 0.1
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "square", 1600, 400, 0.15, 0.09
	return s
}

func soundsCosmos() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{220, 277.18, 329.63, 392}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.09, 0.09, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 277.18, "sine", 0.09, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "triangle", 392, 587.33, 0.22, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 587.33, 196, 0.2, 0.085
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
