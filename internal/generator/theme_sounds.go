package generator

import (
	"encoding/json"
	"html/template"
)

// melodyNote is a single note in a looping ambient melody.
type melodyNote struct {
	Freq float64 `json:"freq"`
	Dur  float64 `json:"dur"`
	Gain float64 `json:"gain,omitempty"`
}

// ambientPreset describes a generative ambient background layer for a theme.
// All fields are optional at runtime; missing values should be treated as
// silence or safe defaults by the consumer.
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

// n is a concise helper for building a melodyNote.
func n(freq, dur float64) melodyNote {
	return melodyNote{Freq: freq, Dur: dur}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 50, Wave: "sine",
		DroneFreqs: []float64{523.25, 659.25, 783.99, 1046.5},
		Attack: 0.6, Release: 1.5, CutoffMin: 800, CutoffMax: 3000,
		Melody: []melodyNote{
			n(523.25, 0.15), n(659.25, 0.15), n(783.99, 0.15), n(1046.5, 0.15),
			n(783.99, 0.15), n(659.25, 0.15), n(523.25, 0.15), n(659.25, 0.15),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 120, Wave: "triangle",
		DroneFreqs: []float64{261.63, 523.25, 659.25},
		Attack: 0.2, Release: 0.6, CutoffMin: 1500, CutoffMax: 6000, DetuneCents: 8,
		Melody: []melodyNote{
			n(523.25, 0.08), n(659.25, 0.08), n(783.99, 0.08), n(1046.5, 0.08),
			n(1318.5, 0.08), n(1046.5, 0.08), n(783.99, 0.08), n(659.25, 0.08),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 40, Wave: "square",
		DroneFreqs: []float64{60, 120},
		Attack: 0.1, Release: 0.3, PulseInterval: 2.0,
		Melody: []melodyNote{
			n(196.00, 0.6), n(246.94, 0.3),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.05, BPM: 160, Wave: "square",
		DroneFreqs: []float64{50, 100},
		Attack: 0.05, Release: 0.15, PulseInterval: 0.25,
		Melody: []melodyNote{
			n(196.00, 0.15), n(246.94, 0.15), n(293.66, 0.15), n(392.00, 0.15),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 55, Wave: "sine",
		DroneFreqs: []float64{98, 196, 293.66},
		Attack: 0.8, Release: 1.5,
		Melody: []melodyNote{
			n(220.00, 0.3), n(261.63, 0.3), n(329.63, 0.3),
			n(293.66, 0.3), n(261.63, 0.3), n(246.94, 0.3), n(220.00, 0.3),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 110, Wave: "triangle",
		DroneFreqs: []float64{65.41, 130.81, 196},
		Attack: 0.3, Release: 0.7, DetuneCents: 12,
		Melody: []melodyNote{
			n(220.00, 0.15), n(261.63, 0.15), n(329.63, 0.15), n(392.00, 0.15),
			n(329.63, 0.15), n(261.63, 0.15), n(220.00, 0.15), n(261.63, 0.15),
			n(329.63, 0.15), n(440.00, 0.15),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.03, BPM: 70, Wave: "triangle",
		DroneFreqs: []float64{329.63, 466.16},
		Attack: 0.4, Release: 0.9, DetuneCents: 15,
		Melody: []melodyNote{
			n(261.63, 0.25), n(369.99, 0.25), n(261.63, 0.25), n(369.99, 0.25),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 140, Wave: "square",
		DroneFreqs: []float64{164.81, 329.63},
		Attack: 0.1, Release: 0.4, DetuneCents: 25,
		CutoffMin: 400, CutoffMax: 3000, NoiseGain: 0.02,
		Melody: []melodyNote{
			n(261.63, 0.12), n(369.99, 0.12), n(261.63, 0.12), n(369.99, 0.12),
			n(523.25, 0.12), n(369.99, 0.12), n(261.63, 0.12), n(369.99, 0.12),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 45, Wave: "triangle",
		DroneFreqs: []float64{80, 120},
		Attack: 0.5, Release: 1.0,
		Melody: []melodyNote{
			n(65.41, 1.2), n(98.00, 0.4), n(65.41, 1.2),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 100, Wave: "square",
		DroneFreqs: []float64{40, 80},
		Attack: 0.05, Release: 0.3,
		Melody: []melodyNote{
			n(65.41, 0.4), n(98.00, 0.2), n(65.41, 0.4), n(77.78, 0.2),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 35, Wave: "sine",
		DroneFreqs: []float64{82.41, 123.47},
		Attack: 1.0, Release: 2.0,
		Melody: []melodyNote{
			n(65.41, 0.5), n(77.78, 0.5), n(98.00, 0.5), n(130.81, 0.5),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 85, Wave: "triangle",
		DroneFreqs: []float64{65.41, 98},
		Attack: 0.2, Release: 0.6, NoiseGain: 0.015,
		Melody: []melodyNote{
			n(65.41, 0.25), n(77.78, 0.25), n(98.00, 0.25), n(130.81, 0.25),
			n(155.56, 0.25), n(196.00, 0.25),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 40, Wave: "sine",
		DroneFreqs: []float64{440, 880, 1320},
		Attack: 1.2, Release: 2.5,
		Melody: []melodyNote{
			n(659.25, 0.3), n(783.99, 0.3), n(987.77, 0.3), n(1318.5, 0.3),
			n(987.77, 0.3), n(783.99, 0.3), n(659.25, 0.3), n(783.99, 0.3),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.05, BPM: 95, Wave: "sine",
		DroneFreqs: []float64{220, 440, 880},
		Attack: 0.3, Release: 0.8, DetuneCents: 20,
		Melody: []melodyNote{
			n(659.25, 0.15), n(783.99, 0.15), n(987.77, 0.15), n(1318.5, 0.15),
			n(1567.98, 0.15), n(1318.5, 0.15), n(987.77, 0.15), n(783.99, 0.15),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 75, Wave: "square",
		DroneFreqs: []float64{523.25, 659.25},
		Attack: 0.1, Release: 0.3,
		Melody: []melodyNote{
			n(293.66, 0.2), n(440.00, 0.2), n(587.33, 0.2), n(739.99, 0.2),
			n(587.33, 0.2), n(440.00, 0.2), n(293.66, 0.2),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 160, Wave: "square",
		DroneFreqs: []float64{261.63, 523.25},
		Attack: 0.05, Release: 0.15,
		Melody: []melodyNote{
			n(293.66, 0.1), n(440.00, 0.1), n(587.33, 0.1), n(739.99, 0.1),
			n(880.00, 0.1), n(739.99, 0.1), n(587.33, 0.1), n(440.00, 0.1),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 30, Wave: "sine",
		DroneFreqs: []float64{130.81, 174.61},
		Attack: 1.5, Release: 3.0, NoiseGain: 0.02,
		Melody: []melodyNote{
			n(130.81, 0.35), n(164.81, 0.35), n(196.00, 0.35), n(164.81, 0.35),
			n(130.81, 0.35), n(196.00, 0.35), n(164.81, 0.35), n(130.81, 0.35),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 80, Wave: "triangle",
		DroneFreqs: []float64{98, 130.81},
		Attack: 0.3, Release: 0.7, NoiseGain: 0.04,
		Melody: []melodyNote{
			n(130.81, 0.18), n(164.81, 0.18), n(196.00, 0.18), n(261.63, 0.18),
			n(196.00, 0.18), n(164.81, 0.18), n(130.81, 0.18), n(98.00, 0.18),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 60, Wave: "square",
		DroneFreqs: []float64{200, 400},
		Attack: 0.05, Release: 0.15,
		Melody: []melodyNote{
			n(261.63, 0.18), n(329.63, 0.18), n(392.00, 0.18), n(523.25, 0.18),
			n(392.00, 0.18), n(329.63, 0.18), n(261.63, 0.18),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.065, BPM: 200, Wave: "square",
		DroneFreqs: []float64{100, 200},
		Attack: 0.02, Release: 0.08,
		Melody: []melodyNote{
			n(261.63, 0.08), n(329.63, 0.08), n(392.00, 0.08), n(523.25, 0.08),
			n(659.25, 0.08), n(523.25, 0.08), n(392.00, 0.08), n(329.63, 0.08),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 80, Wave: "square",
		DroneFreqs: []float64{523.25, 659.25},
		Attack: 0.1, Release: 0.3,
		Melody: []melodyNote{
			n(261.63, 0.25), n(349.23, 0.25), n(440.00, 0.25), n(523.25, 0.25),
			n(440.00, 0.25), n(349.23, 0.25), n(261.63, 0.25),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 155, Wave: "square",
		DroneFreqs: []float64{261.63, 523.25},
		Attack: 0.05, Release: 0.15,
		DetuneCents: 30, NoiseGain: 0.02,
		Melody: []melodyNote{
			n(261.63, 0.12), n(349.23, 0.12), n(440.00, 0.12), n(523.25, 0.12),
			n(698.46, 0.12), n(523.25, 0.12), n(440.00, 0.12), n(349.23, 0.12),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 30, Wave: "sine",
		DroneFreqs: []float64{220, 440, 660},
		Attack: 1.5, Release: 3.0,
		Melody: []melodyNote{
			n(130.81, 0.5), n(196.00, 0.5), n(261.63, 0.5), n(329.63, 0.5),
			n(261.63, 0.5), n(196.00, 0.5), n(130.81, 0.5),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.05, BPM: 75, Wave: "triangle",
		DroneFreqs: []float64{110, 220, 440},
		Attack: 0.3, Release: 0.8,
		Melody: []melodyNote{
			n(130.81, 0.2), n(196.00, 0.2), n(261.63, 0.2), n(329.63, 0.2),
			n(392.00, 0.2), n(329.63, 0.2), n(261.63, 0.2), n(196.00, 0.2),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 55, Wave: "triangle",
		DroneFreqs: []float64{220, 330, 440},
		Attack: 0.8, Release: 1.8,
		Melody: []melodyNote{
			n(261.63, 0.3), n(311.13, 0.3), n(392.00, 0.3), n(466.16, 0.3),
			n(392.00, 0.3), n(311.13, 0.3), n(261.63, 0.3),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 120, Wave: "square",
		DroneFreqs: []float64{110, 220},
		Attack: 0.05, Release: 0.15,
		Melody: []melodyNote{
			n(261.63, 0.12), n(311.13, 0.12), n(392.00, 0.12), n(466.16, 0.12),
			n(523.25, 0.12), n(466.16, 0.12), n(392.00, 0.12), n(311.13, 0.12),
		},
	}
	return s
}

// soundsSpaceage returns synth parameters for the Space Age theme.
func soundsSpaceage() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{440, 554.37, 659.25, 880}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.07, 0.085, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 554.37, "triangle", 0.075, 0.09
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.09
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 659.25, 330, 0.17, 0.085
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 240, 120, 0.13, 0.08
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 50, Wave: "sine",
		DroneFreqs: []float64{440, 880},
		Attack: 0.8, Release: 1.5,
		Melody: []melodyNote{
			n(261.63, 0.35), n(329.63, 0.35), n(392.00, 0.35), n(523.25, 0.35),
			n(392.00, 0.35), n(329.63, 0.35), n(261.63, 0.35),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.05, BPM: 110, Wave: "triangle",
		DroneFreqs: []float64{220, 440},
		Attack: 0.1, Release: 0.4, PulseInterval: 0.5,
		Melody: []melodyNote{
			n(261.63, 0.15), n(329.63, 0.15), n(392.00, 0.15), n(523.25, 0.15),
			n(659.25, 0.15), n(523.25, 0.15), n(392.00, 0.15), n(329.63, 0.15),
		},
	}
	return s
}

// soundsTropical returns synth parameters for the Tropical Beach theme.
func soundsTropical() themeSounds {
	var s themeSounds
	s.Splash.Freqs = []float64{523.25, 659.25, 783.99, 1046.5}
	s.Splash.Spacing, s.Splash.Gain, s.Splash.Wave = 0.085, 0.08, "sine"
	s.Nav.Freq, s.Nav.Wave, s.Nav.Dur, s.Nav.Gain = 880, "sine", 0.07, 0.07
	s.Open.Wave, s.Open.Start, s.Open.End, s.Open.Dur, s.Open.Gain = "sine", 440, 880, 0.18, 0.08
	s.Close.Wave, s.Close.Start, s.Close.End, s.Close.Dur, s.Close.Gain = "sine", 660, 330, 0.17, 0.075
	s.Bounce.Wave, s.Bounce.Start, s.Bounce.End, s.Bounce.Dur, s.Bounce.Gain = "sine", 200, 100, 0.12, 0.07
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 65, Wave: "sine",
		DroneFreqs: []float64{261.63, 329.63, 392, 523.25},
		Attack: 0.6, Release: 1.5, NoiseGain: 0.015,
		Melody: []melodyNote{
			n(261.63, 0.28), n(293.66, 0.28), n(329.63, 0.28), n(392.00, 0.28), n(440.00, 0.28),
			n(392.00, 0.28), n(329.63, 0.28), n(293.66, 0.28), n(261.63, 0.28),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 140, Wave: "triangle",
		DroneFreqs: []float64{130.81, 261.63, 329.63, 392},
		Attack: 0.1, Release: 0.3, NoiseGain: 0.03,
		Melody: []melodyNote{
			n(261.63, 0.12), n(293.66, 0.12), n(329.63, 0.12), n(392.00, 0.12), n(440.00, 0.12),
			n(523.25, 0.12), n(440.00, 0.12), n(392.00, 0.12), n(329.63, 0.12), n(293.66, 0.12),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.02, BPM: 40, Wave: "sine",
		DroneFreqs: []float64{174.61, 220, 261.63},
		Attack: 1.0, Release: 2.0, NoiseGain: 0.015,
		Melody: []melodyNote{
			n(130.81, 0.45), n(155.56, 0.45), n(196.00, 0.45), n(233.08, 0.45),
			n(196.00, 0.45), n(155.56, 0.45), n(130.81, 0.45),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.05, BPM: 90, Wave: "triangle",
		DroneFreqs: []float64{130.81, 174.61, 220},
		Attack: 0.3, Release: 0.7, PulseInterval: 1.2,
		Melody: []melodyNote{
			n(130.81, 0.2), n(155.56, 0.2), n(196.00, 0.2), n(233.08, 0.2),
			n(261.63, 0.2), n(233.08, 0.2), n(196.00, 0.2), n(155.56, 0.2),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 35, Wave: "sine",
		DroneFreqs: []float64{293.66, 440, 587.33},
		Attack: 1.5, Release: 3.0,
		Melody: []melodyNote{
			n(130.81, 0.5), n(164.81, 0.5), n(196.00, 0.5), n(261.63, 0.5), n(329.63, 0.5),
			n(261.63, 0.5), n(196.00, 0.5), n(164.81, 0.5), n(130.81, 0.5),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 70, Wave: "triangle",
		DroneFreqs: []float64{146.83, 293.66, 440},
		Attack: 0.3, Release: 0.8,
		Melody: []melodyNote{
			n(130.81, 0.22), n(164.81, 0.22), n(196.00, 0.22), n(261.63, 0.22), n(329.63, 0.22),
			n(392.00, 0.22), n(329.63, 0.22), n(261.63, 0.22), n(196.00, 0.22), n(164.81, 0.22),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 55, Wave: "square",
		DroneFreqs: []float64{440, 660},
		Attack: 0.2, Release: 0.5,
		Melody: []melodyNote{
			n(261.63, 0.35), n(349.23, 0.35), n(261.63, 0.35), n(349.23, 0.35),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.06, BPM: 140, Wave: "square",
		DroneFreqs: []float64{220, 440},
		Attack: 0.05, Release: 0.15, PulseInterval: 0.3,
		Melody: []melodyNote{
			n(261.63, 0.12), n(349.23, 0.12), n(261.63, 0.12), n(349.23, 0.12),
			n(523.25, 0.12), n(349.23, 0.12), n(261.63, 0.12), n(349.23, 0.12),
		},
	}
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
	s.Ambient.Normal = ambientPreset{
		Gain: 0.025, BPM: 50, Wave: "triangle",
		DroneFreqs: []float64{164.81, 246.94},
		Attack: 0.7, Release: 1.5, DetuneCents: 8,
		Melody: []melodyNote{
			n(130.81, 0.3), n(155.56, 0.3), n(185.00, 0.3), n(220.00, 0.3),
			n(185.00, 0.3), n(155.56, 0.3), n(130.81, 0.3),
		},
	}
	s.Ambient.Wild = ambientPreset{
		Gain: 0.055, BPM: 115, Wave: "square",
		DroneFreqs: []float64{82.41, 164.81},
		Attack: 0.1, Release: 0.3, DetuneCents: 18,
		Melody: []melodyNote{
			n(130.81, 0.13), n(155.56, 0.13), n(185.00, 0.13), n(220.00, 0.13),
			n(261.63, 0.13), n(220.00, 0.13), n(185.00, 0.13), n(155.56, 0.13),
		},
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
