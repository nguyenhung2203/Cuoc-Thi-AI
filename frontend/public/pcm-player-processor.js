// pcm-player-processor.js
// AudioWorklet that plays a stream of PCM16 @24kHz chunks (Gemini Live output).
// The main thread posts Int16 PCM ArrayBuffers; this queues and plays them
// gaplessly. Posting the string 'flush' clears the queue (barge-in / interrupt).

class PCMPlayerProcessor extends AudioWorkletProcessor {
  constructor() {
    super();
    this._queue = [];      // array of Float32Array
    this._readIdx = 0;     // read offset into queue[0]
    this._playing = false;
    this.port.onmessage = (e) => {
      if (e.data === 'flush') {
        this._queue = [];
        this._readIdx = 0;
        return;
      }
      // Incoming is a PCM16 ArrayBuffer @24kHz.
      const pcm = new Int16Array(e.data);
      const f = new Float32Array(pcm.length);
      for (let i = 0; i < pcm.length; i++) f[i] = pcm[i] / 0x8000;
      this._queue.push(f);
    };
  }

  process(_inputs, outputs) {
    const out = outputs[0][0];
    if (!out) return true;
    let i = 0;
    const wasPlaying = this._playing;

    while (i < out.length) {
      if (this._queue.length === 0) {
        // Underrun: fill silence.
        out[i++] = 0;
        continue;
      }
      const cur = this._queue[0];
      out[i++] = cur[this._readIdx++];
      if (this._readIdx >= cur.length) {
        this._queue.shift();
        this._readIdx = 0;
      }
    }

    this._playing = this._queue.length > 0;
    // Notify main thread when playback starts/stops (drives avatar animation).
    if (this._playing !== wasPlaying) {
      this.port.postMessage(this._playing ? 'speaking' : 'idle');
    }
    return true;
  }
}

registerProcessor('pcm-player-processor', PCMPlayerProcessor);
