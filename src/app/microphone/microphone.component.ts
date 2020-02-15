import { Component, OnInit } from '@angular/core';
import { MicrophoneService } from './microphone.service';
import { Tuner } from './tuner';

@Component({
  selector: 'app-microphone',
  templateUrl: './microphone.component.html',
  styleUrls: ['./microphone.component.scss']
})
export class MicrophoneComponent implements OnInit {

  frequency: number;
  offset: number;
  note: string;
  noteNum: number;
  octave: number;
  
  private tuner: Tuner;

  constructor(private micService: MicrophoneService) {}

  ngOnInit() {
    document.addEventListener('click', async () => {
      this.tuner = await this.micService.getTuner();
      requestAnimationFrame(this.dataLoop.bind(this));
    });
  }

  private dataLoop(){
    var tuneResult = this.tuner.getTuneResult();
    if (tuneResult != null && tuneResult.frequency){
      this.frequency = tuneResult.frequency;
      this.offset = tuneResult.offset;
      this.note = tuneResult.note;
      this.noteNum = tuneResult.noteNum;
      this.octave = tuneResult.octave;
    }
    requestAnimationFrame(this.dataLoop.bind(this));
  }
}
