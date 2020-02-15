export interface FrequencyInfo {
    note: string;
    frequency: number;
    octave: number;
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
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
        this.context = new (AudioContext || window['webkitAudioContext'])();
        const source = this.context.createMediaStreamSource(stream);
        this.fft = this.context.createAnalyser();
        this.fft.smoothingTimeConstant = 0;
        this.fft.fftSize = 2048;
        this.buffer = new Float32Array(this.fft.fftSize);
        source.connect(this.fft);
       } 
    }

    getFrequencyInfo(): FrequencyInfo | null {
        if (!this.fft || !window['getAudioFrequencyInfo']){
            return null;
        }

        this.fft.getFloatTimeDomainData(this.buffer);
        return window['getAudioFrequencyInfo'](this.buffer, this.context.sampleRate);
    }


}