import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class StreamService {
  constructor() {}

  streamData(
    url: string,
    message: string,
    onChunk: (chunk: string) => void,
    onComplete: (res: any) => void,
    onError: (error: any) => void
  ): void {
    fetch(url, {
      method: 'POST',
      body: JSON.stringify({ message }),
    })
      .then((response) => {
        if (!response.body) {
          throw new Error('ReadableStream not supported in this browser.');
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder('utf-8');

        const readChunk = () => {
          reader
            .read()
            .then(({ done, value }) => {
              if (done) {
                onComplete(value);
                return;
              }

              const chunk = decoder.decode(value, { stream: true });
              onChunk(chunk);
              readChunk(); // Read the next chunk
            })
            .catch(onError);
        };

        readChunk(); // Start reading the first chunk
      })
      .catch(onError);
  }
}
