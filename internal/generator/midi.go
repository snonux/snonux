package generator

import (
	"bytes"
	"encoding/base64"
	"math"
)

func buildBackgroundMIDI(s themeSounds) ([]byte, error) {
	return nil, nil
}

func writeVLQ(buf *bytes.Buffer, val uint32) {
	var b [4]byte
	n := 4
	for i := 3; i >= 0; i-- {
		b[i] = byte(val & 0x7f)
		val >>= 7
		if val == 0 {
			n = i
			break
		}
	}
	for i := n; i < 3; i++ {
		b[i] |= 0x80
	}
	buf.Write(b[n:])
}

func freqToMIDI(freq float64) uint8 {
	if freq <= 0 {
		return 60
	}
	midi := int(math.Round(12*math.Log2(freq/440.0) + 69))
	if midi < 0 {
		midi = 0
	}
	if midi > 127 {
		midi = 127
	}
	return uint8(midi)
}

func midiToBase64(midi []byte) string {
	if midi == nil {
		return ""
	}
	return "data:audio/midi;base64," + base64.StdEncoding.EncodeToString(midi)
}
