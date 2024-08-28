import { HttpClient } from '@angular/common/http';
import {
  Component,
  ElementRef,
  inject,
  signal,
  ViewChild,
} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { take, tap } from 'rxjs';
import { ReplaceSTPipe } from './replace-st.pipe';
import { StreamService } from './stream-service.service';

enum commitState {
  Idle = 0,
  Success = 1,
  Error = 2,
}

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, ReplaceSTPipe],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  private http = inject(HttpClient);
  @ViewChild('input') input!: ElementRef<HTMLTextAreaElement>;
  output = signal('');
  commitState = signal(commitState.Idle); // Initialize with the Idle state

  ICommitState = commitState;

  constructor(private streamService: StreamService) {
    this.getMessage();
  }

  getMessage() {
    this.http
      .get<string>('/message')
      .pipe(
        take(1),
        tap((message) => {
          if (this.input && message) {
            this.input.nativeElement.value = message.trim();
          }
        })
      )
      .subscribe(console.log);
  }

  commit() {
    if (!this.input.nativeElement.value) {
      console.log('No input value');
      return;
    }

    this.streamService.streamData(
      '/commit',
      this.input.nativeElement.value,
      (chunk) => {
        if (
          chunk
            .toLowerCase()
            .includes('Command executed successfully'.toLowerCase())
        ) {
          console.log('Success');

          this.commitState.set(commitState.Success);
        }
        if (
          chunk.toLowerCase().includes('Command execution error'.toLowerCase())
        ) {
          console.log('Error');
          this.commitState.set(commitState.Error);
        }

        chunk.trim() &&
          this.output.set(this.output() + chunk + '<br>' + '<br>');
      },
      (res) => {},
      (error) => {
        console.error('Error receiving stream:', error);
      }
    );
  }

  close() {
    this.http.post<string>('/close', {}).pipe(take(1)).subscribe();
  }

  replaceStdoutWithNewline(input: string): string {
    return input.replace(/STDOUT:\s*/g, '\n').trim();
  }
}
