package generator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"sync"

	"codeberg.org/snonux/snonux/internal/generator/templates"
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
	File          string       `json:"file,omitempty"`
	Volume        float64      `json:"volume,omitempty"`
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

// soundCache holds unmarshalled themeSounds indexed by theme name. It is
// populated lazily on first access so repeated calls to loadThemeSounds are
// essentially free.
var (
	soundCache   map[string]themeSounds
	soundCacheMu sync.RWMutex
)

// initSoundCache eagerly loads every sounds.json from the embedded theme FS.
// It is called lazily under soundCacheMu.
func initSoundCache() {
	soundCache = make(map[string]themeSounds)
	for name := range getThemeSet() {
		b, err := templates.ThemeSounds(name)
		if err != nil {
			continue
		}
		var s themeSounds
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		soundCache[name] = s
	}
}

// loadThemeSounds returns the unmarshaled sound preset for a theme.
// On cache miss it falls back to the empty zero value.
func loadThemeSounds(themeName string) (themeSounds, error) {
	soundCacheMu.RLock()
	cached, ok := soundCache[themeName]
	soundCacheMu.RUnlock()
	if ok {
		return cached, nil
	}

	soundCacheMu.Lock()
	defer soundCacheMu.Unlock()
	if soundCache == nil {
		initSoundCache()
		if s, ok := soundCache[themeName]; ok {
			return s, nil
		}
	}
	return themeSounds{}, fmt.Errorf("no sounds.json for theme %q", themeName)
}

// defaultSounds returns the default theme (neon) sound preset.
func defaultSounds() themeSounds {
	s, _ := loadThemeSounds("neon")
	return s
}

// themeSoundsJSON returns a JS object literal for embedding in <script> (safe JSON).
func themeSoundsJSON(themeName string) template.JS {
	s, _ := loadThemeSounds(themeName)
	b, err := json.Marshal(s)
	if err != nil {
		b, _ = json.Marshal(themeSounds{})
	}
	return template.JS(b) //nolint:gosec // JSON from fixed structs
}
