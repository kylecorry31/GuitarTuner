import { TestBed } from '@angular/core/testing';

import { MicrophoneService } from './microphone.service';

describe('MicrophoneService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: MicrophoneService = TestBed.get(MicrophoneService);
    expect(service).toBeTruthy();
  });
});
