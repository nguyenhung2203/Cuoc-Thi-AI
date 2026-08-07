// mic-capture-processor.js
// AudioWorklet that downsamples the mic (usually 48kHz float32) to
// PCM16 mono @16kHz and posts ArrayBuffer chunks back to the main thread.
// Gemini Live expects 16kHz little-endian signed 16-bit PCM for input.

class MicCaptureProcessor extends AudioWorkletProcessor {
  constructor() {
    super();
    this.targetRate = 16000;
    this._buf = [];
    this._bufLen = 0;
    // ~100ms per posted chunk at 16kHz => 1600 samples.
    this._chunkSamples = 1600;
  }

  // Linear-interpolation resample from sampleRate -> targetRate.
  _resample(input) {
    const ratio = sampleRate / this.targetRate;
    const outLen = Math.floor(input.length / ratio);
    const out = new Float32Array(outLen);
    for (let i = 0; i < outLen; i++) {
      const idx = i * ratio;
      const i0 = Math.floor(idx);
      const i1 = Math.min(i0 + 1, input.length - 1);
      const frac = idx - i0;
      out[i] = input[i0] * (1 - frac) + input[i1] * frac;
    }
    return out;
  }

  process(inputs) {
    const input = inputs[0];
    if (!input || !input[0]) return true;
    const resampled = this._resample(input[0]);
    this._buf.push(resampled);
    this._bufLen += resampled.length;

    if (this._bufLen >= this._chunkSamples) {
      // Flatten collected float samples.
      const flat = new Float32Array(this._bufLen);
      let off = 0;
      for (const b of this._buf) { flat.set(b, off); off += b.length; }
      this._buf = [];
      this._bufLen = 0;

      // Convert float [-1,1] -> PCM16.
      const pcm = new Int16Array(flat.length);
      for (let i = 0; i < flat.length; i++) {
        let s = Math.max(-1, Math.min(1, flat[i]));
        pcm[i] = s < 0 ? s * 0x8000 : s * 0x7fff;
      }
      this.port.postMessage(pcm.buffer, [pcm.buffer]);
    }
    return true;
  }
}

registerProcessor('mic-capture-processor', MicCaptureProcessor);
