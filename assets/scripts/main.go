package main

import (
	"math"
	"syscall/js"
	"time"
)

type note struct {
	name      rune
	frequency float64
}

var notes = []note{
	note{'E', 82.4069},
	note{'A', 110},
	note{'D', 146.832},
	note{'G', 195.998},
	note{'B', 246.932},
	note{'e', 329.628},
}

var assessStringsUntilTime int64 = 0
var assessedStringsInLastFrame = false

const rmsThreshold = 0.006

var lastRms = 0.0

var differences = make([]float64, len(notes))
var lastMinDifference = 0

func jsGetFrequency(this js.Value, args []js.Value) interface{} {
	length := args[0].Length()
	dst := make([]float64, length)

	for i := 0; i < length; i++ {
		dst[i] = args[0].Index(i).Float()
	}

	// js.CopyBytesToGo(dst, args[0])
	sampleRate := args[1].Float()
	return getFrequency(dst, sampleRate)
}

func jsGetNote(this js.Value, args []js.Value) interface{} {
	frequency := args[0].Float()
	return string(getNote(frequency).name)
}

func jsGetNoteNumber(this js.Value, args []js.Value) interface{} {
	frequency := args[0].Float()
	return getNoteNumber(frequency)
}

func jsGetOctave(this js.Value, args []js.Value) interface{} {
	frequency := args[0].Float()
	return getOctave(frequency)
}

func jsGetOffset(this js.Value, args []js.Value) interface{} {
	frequency := args[0].Float()
	return getOffset(frequency)
}

func getFrequency(waveform []float64, sampleRate float64) float64 {

	searchSize := len(waveform) / 2

	tolerance := 0.001
	rms := 0.0
	rmsMin := 0.008

	prevAssessedStrings := assessedStringsInLastFrame

	for _, amplitude := range waveform {
		rms += amplitude * amplitude
	}

	rms = math.Sqrt(rms / float64(len(waveform)))

	if rms < rmsMin {
		return 0
	}

	time := (time.Now().UnixNano() / 1000000)

	if rms > lastRms+rmsThreshold {
		assessStringsUntilTime = time + 250
	}

	if time < assessStringsUntilTime {
		assessedStringsInLastFrame = true

		for i, note := range notes {
			offset := int(math.Floor(sampleRate / note.frequency))
			difference := 0.0

			if !prevAssessedStrings {
				differences[i] = 0
			}

			for j := 0; j < searchSize; j++ {
				currentAmp := waveform[j]
				offsetAmp := waveform[j+offset]
				difference += math.Abs(currentAmp - offsetAmp)
			}

			difference /= float64(searchSize)

			differences[i] += difference * float64(offset)
		}
	} else {
		assessedStringsInLastFrame = false
	}

	if !assessedStringsInLastFrame && prevAssessedStrings {
		lastMinDifference = argmin(differences)
	}

	assumedString := notes[lastMinDifference]
	searchRange := 10
	actualFrequency := int(math.Round(sampleRate / assumedString.frequency))
	searchStart := actualFrequency - searchRange
	searchEnd := actualFrequency + searchRange
	smallestDifference := math.Inf(1)

	for i := searchStart; i < searchEnd; i++ {
		difference := 0.0

		for j := 0; j < searchSize; j++ {
			currentAmp := float64(waveform[j]) / 255.0
			offsetAmp := float64(waveform[j+i]) / 255.0
			difference += math.Abs(currentAmp - offsetAmp)
		}

		difference /= float64(searchSize)

		if difference < smallestDifference {
			smallestDifference = difference
			actualFrequency = i
		}

		if difference < tolerance {
			actualFrequency = i
			break
		}

	}

	lastRms = rms

	return sampleRate / float64(actualFrequency)
}

func linearizeFreq(frequency float64) float64 {
	return math.Log2(frequency / 440.0)
}

func getOffset(frequency float64) float64 {
	note := getNote(frequency)
	actualNumber := linearizeFreq(frequency)
	desiredNumber := linearizeFreq(note.frequency)

	// semitonesFromA4 := 12 * actualNumber
	// octave := 4 + ((9 + semitonesFromA4) / 12)
	// octave = math.Floor(octave)
	// n := (12 + (int(math.Round(semitonesFromA4)) % 12)) % 12

	// println(octave, n)

	return actualNumber - desiredNumber
}

func getOctave(frequency float64) float64 {
	semitonesFromA4 := 12 * linearizeFreq(frequency)
	octave := 4 + ((9 + semitonesFromA4) / 12)
	octave = math.Floor(octave)
	return octave
}

func getNoteNumber(frequency float64) int {
	semitonesFromA4 := 12 * linearizeFreq(frequency)
	return (12 + (int(math.Round(semitonesFromA4)) % 12)) % 12
}

func getNote(frequency float64) note {
	min := 0

	for i, note := range notes {
		minDiff := math.Abs(notes[min].frequency - frequency)
		diff := math.Abs(note.frequency - frequency)
		if diff < minDiff {
			min = i
		}
	}

	return notes[min]
}

func argmin(arr []float64) int {
	min := 0

	for i, value := range arr {
		if value < arr[min] {
			min = i
		}
	}

	return min
}

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("getAudioFrequency", js.FuncOf(jsGetFrequency))
	js.Global().Set("getNoteOffset", js.FuncOf(jsGetOffset))
	js.Global().Set("getNoteName", js.FuncOf(jsGetNote))
	js.Global().Set("getNoteNumber", js.FuncOf(jsGetNoteNumber))
	js.Global().Set("getNoteOctave", js.FuncOf(jsGetOctave))
	<-wait
}
