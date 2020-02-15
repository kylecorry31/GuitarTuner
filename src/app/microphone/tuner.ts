export interface TuneResult {
    note: string;
    noteNum: number;
    frequency: number;
    octave: number;
    offset: number;
};

export class Tuner {

    private fft: AnalyserNode;
    private context: AudioContext;
    private buffer: Float32Array;

    async initialize(): Promise<void> {
       if (!this.fft){
        navigator.getUserMedia = (navigator.getUserMedia ||
            navigator['webkitGetUserMedia'] ||
            navigator['mozGetUserMedia'] ||
            navigator['msGetUserMedia']);
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true, video: false });
        this.context = new (AudioContext || window['webkitAudioContext'])();
        const source = this.context.createMediaStreamSource(stream);
        this.fft = this.context.createAnalyser();
        this.fft.smoothingTimeConstant = 0;
        this.fft.fftSize = 2048;
        this.buffer = new Float32Array(this.fft.fftSize);
        source.connect(this.fft);
       } 
    }

    getTuneResult(): TuneResult | null {
        if (!this.fft || !window['getAudioFrequency']){
            return null;
        }

        this.fft.getFloatTimeDomainData(this.buffer);
        let freq = window['getAudioFrequency'](this.buffer, this.context.sampleRate);
        let offset = window['getNoteOffset'](freq);
        let note = window['getNoteName'](freq);
        let noteNum = window['getNoteNumber'](freq);
        let octave = window['getNoteOctave'](freq);

        return {frequency: freq, offset: offset, note: note, noteNum: noteNum, octave: octave};
    }


}