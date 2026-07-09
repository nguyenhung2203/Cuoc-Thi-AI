import { ref, shallowRef } from 'vue';

/**
 * useLiveKit — Composable quản lý Camera/Mic cho phòng phỏng vấn.
 * 
 * Chiến lược kết nối 2 lớp:
 * 1. Thử kết nối LiveKit server (production).
 * 2. Nếu không có server → dùng native getUserMedia (dev/demo fallback).
 * 
 * Trong cả 2 trường hợp, camera/mic đều hoạt động bình thường với các nút Toggle.
 */
export function useLiveKit() {
  const room = shallowRef(null);
  
  // Trạng thái kết nối
  const isConnected = ref(false);
  const error = ref(null);
  
  // Trạng thái thiết bị local
  const isMicOn = ref(true);
  const isCameraOn = ref(true);

  // References cho video elements (được bind bởi ref="localVideoEl" trong template)
  const localVideoEl = ref(null);
  const remoteVideoEl = ref(null);

  // LiveKit tracks (khi dùng server)
  let localVideoTrack = null;
  let localAudioTrack = null;

  // Native stream (khi dùng fallback)
  let nativeStream = null;

  // ─── HELPER: Khởi động native camera/mic ────────────────────────────────

  const startNativeMedia = async () => {
    try {
      // Dừng stream cũ nếu có
      if (nativeStream) {
        nativeStream.getTracks().forEach(t => t.stop());
      }
      nativeStream = await navigator.mediaDevices.getUserMedia({
        video: isCameraOn.value,
        audio: isMicOn.value,
      });

      // Gắn vào video element
      if (localVideoEl.value) {
        localVideoEl.value.srcObject = nativeStream;
        localVideoEl.value.muted = true;
        await localVideoEl.value.play().catch(() => {});
      }
    } catch (err) {
      console.warn('[useLiveKit] getUserMedia failed:', err.name, err.message);
      error.value = 'Không thể truy cập Camera/Mic. Hãy kiểm tra quyền trình duyệt.';
    }
  };

  // ─── CONNECT ─────────────────────────────────────────────────────────────

  const connectToRoom = async (url, token) => {
    error.value = null;

    // Thử kết nối LiveKit với timeout 3s
    try {
      const { Room, RoomEvent, createLocalVideoTrack: clvt, createLocalAudioTrack: clat } =
        await import('livekit-client');

      const newRoom = new Room({
        adaptiveStream: true,
        dynacast: true,
        videoCaptureDefaults: { resolution: { width: 640, height: 480, frameRate: 30 } },
      });

      newRoom.on(RoomEvent.TrackSubscribed, (track) => {
        if (track.kind === 'video' && remoteVideoEl.value) track.attach(remoteVideoEl.value);
        else if (track.kind === 'audio') {
          const el = document.createElement('audio');
          document.body.appendChild(el);
          track.attach(el);
        }
      });
      newRoom.on(RoomEvent.TrackUnsubscribed, (track) => track.detach());
      newRoom.on(RoomEvent.Disconnected, () => { isConnected.value = false; _cleanLiveKit(); });

      // Timeout 3 giây để không chờ mãi
      await Promise.race([
        newRoom.connect(url, token),
        new Promise((_, reject) => setTimeout(() => reject(new Error('timeout')), 3000)),
      ]);

      room.value = newRoom;
      isConnected.value = true;

      // Publish local tracks qua LiveKit
      if (isCameraOn.value) {
        localVideoTrack = await clvt();
        await newRoom.localParticipant.publishTrack(localVideoTrack);
        if (localVideoEl.value) localVideoTrack.attach(localVideoEl.value);
      }
      if (isMicOn.value) {
        localAudioTrack = await clat();
        await newRoom.localParticipant.publishTrack(localAudioTrack);
      }

      console.info('[useLiveKit] Connected via LiveKit server.');
      return; // thành công, không cần fallback
    } catch (err) {
      console.warn('[useLiveKit] LiveKit unavailable, using native getUserMedia fallback.', err.message);
      _cleanLiveKit();
    }

    // Fallback: native getUserMedia
    await startNativeMedia();
  };

  // ─── TOGGLE MIC ──────────────────────────────────────────────────────────

  const toggleMic = async () => {
    if (room.value) {
      // LiveKit mode
      try {
        if (isMicOn.value && localAudioTrack) {
          await room.value.localParticipant.unpublishTrack(localAudioTrack, true);
          localAudioTrack.stop();
          localAudioTrack = null;
        } else {
          const { createLocalAudioTrack: clat } = await import('livekit-client');
          localAudioTrack = await clat();
          await room.value.localParticipant.publishTrack(localAudioTrack);
        }
      } catch (e) { console.warn('[useLiveKit] toggleMic LiveKit error', e); }
    } else if (nativeStream) {
      // Native fallback mode
      nativeStream.getAudioTracks().forEach(t => { t.enabled = !isMicOn.value; });
    }
    isMicOn.value = !isMicOn.value;
  };

  // ─── TOGGLE CAMERA ───────────────────────────────────────────────────────

  const toggleCamera = async () => {
    if (room.value) {
      // LiveKit mode
      try {
        if (isCameraOn.value && localVideoTrack) {
          await room.value.localParticipant.unpublishTrack(localVideoTrack, true);
          localVideoTrack.detach();
          localVideoTrack.stop();
          if (localVideoTrack.mediaStreamTrack) {
            localVideoTrack.mediaStreamTrack.stop();
          }
          localVideoTrack = null;
          if (localVideoEl.value) localVideoEl.value.srcObject = null;
        } else {
          const { createLocalVideoTrack: clvt } = await import('livekit-client');
          localVideoTrack = await clvt();
          await room.value.localParticipant.publishTrack(localVideoTrack);
          if (localVideoEl.value) localVideoTrack.attach(localVideoEl.value);
        }
      } catch (e) { console.warn('[useLiveKit] toggleCamera LiveKit error', e); }
    } else if (nativeStream) {
      // Native fallback mode
      const videoTracks = nativeStream.getVideoTracks();
      if (isCameraOn.value) {
        // Tắt camera và HOÀN TOÀN NGẮT track khỏi hardware để đèn LED camera tắt
        videoTracks.forEach(t => { 
          t.enabled = false;
          t.stop(); // Ngắt hẳn hardware
          nativeStream.removeTrack(t);
        });
        if (localVideoEl.value) localVideoEl.value.srcObject = null;
      } else {
        // Bật lại camera — xin lại track mới từ hardware vì track cũ đã stop()
        try {
          const newCamStream = await navigator.mediaDevices.getUserMedia({ video: true });
          newCamStream.getVideoTracks().forEach(t => nativeStream.addTrack(t));
          if (localVideoEl.value) {
            localVideoEl.value.srcObject = nativeStream;
            await localVideoEl.value.play().catch(() => {});
          }
        } catch (err) {
          console.error('[useLiveKit] Không thể bật lại camera:', err);
        }
      }
    }
    isCameraOn.value = !isCameraOn.value;
  };

  // ─── DISCONNECT ──────────────────────────────────────────────────────────

  const disconnect = () => {
    if (room.value) { room.value.disconnect(); }
    _cleanLiveKit();
    if (nativeStream) {
      nativeStream.getTracks().forEach(t => t.stop());
      nativeStream = null;
    }
    if (localVideoEl.value) localVideoEl.value.srcObject = null;
  };

  const _cleanLiveKit = () => {
    if (localVideoTrack) { localVideoTrack.detach(); localVideoTrack.stop(); localVideoTrack = null; }
    if (localAudioTrack) { localAudioTrack.stop(); localAudioTrack = null; }
    isConnected.value = false;
    room.value = null;
  };

  return {
    room,
    isConnected,
    error,
    isMicOn,
    isCameraOn,
    localVideoEl,
    remoteVideoEl,
    connectToRoom,
    toggleMic,
    toggleCamera,
    disconnect,
  };
}
