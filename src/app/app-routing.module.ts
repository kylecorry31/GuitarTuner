import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MicrophoneComponent } from './microphone/microphone.component';


const routes: Routes = [
  { path: "**", component: MicrophoneComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
