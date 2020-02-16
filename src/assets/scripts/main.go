package main

import (
	"syscall/js"

	"github.com/kylecorry31/gotar-tuner"
)

func getAudioFrequencyInfo(this js.Value, args []js.Value) interface{} {

	if len(args) < 2 {
		return js.Null()
	}

	if args[0].Length() == 0 {
		return js.Null()
	}

	if args[1].Float() <= 0 {
		return js.Null()
	}

	length := args[0].Length()
	dst := make([]float64, length)

	for i := 0; i < length; i++ {
		dst[i] = args[0].Index(i).Float()
	}

	sampleRate := args[1].Float()

	calculator := gotar.ZeroCrossingFrequencyCalculator{0.25}
	frequency := calculator.GetFrequency(dst, sampleRate)
	freqInfo := gotar.CreateFrequencyInfo(frequency)

	obj := make(map[string]interface{})
	obj["frequency"] = freqInfo.Frequency
	obj["note"] = freqInfo.Note
	obj["octave"] = freqInfo.Octave

	return js.ValueOf(obj)
}

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("getAudioFrequencyInfo", js.FuncOf(getAudioFrequencyInfo))
	<-wait
}
