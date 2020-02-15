import { Injectable } from '@angular/core';
import { Tuner } from './tuner';

@Injectable({
  providedIn: 'root'
})
export class MicrophoneService {

  constructor() { }

  async getTuner(): Promise<Tuner> {
    var tuner = new Tuner();
    await tuner.initialize();
    return tuner;
  }

}
