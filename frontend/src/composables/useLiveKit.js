import { ref, shallowRef } from 'vue';
import { Room, RoomEvent, createLocalVideoTrack, createLocalAudioTrack } from 'livekit-client';

export function useLiveKit() {
  const room = shallowRef(null);
  
  // Trạng thái kết nối
  const isConnected = ref(false);
  const error = ref(null);
  
  // Trạng thái thiết bị local
  const isMicOn = ref(true);
  const isCameraOn = ref(true);

  // References cho video elements
  const localVideoEl = ref(null);
  const remoteVideoEl = ref(null);

  // Tracks
  let localVideoTrack = null;
  let localAudioTrack = null;

  const connectToRoom = async (url, token) => {
    try {
      error.value = null;
      
      const newRoom = new Room({
        adaptiveStream: true,
        dynacast: true,
        videoCaptureDefaults: {
          resolution: { width: 640, height: 480, frameRate: 30 },
        },
      });

      // Lắng nghe sự kiện remote track
      newRoom.on(RoomEvent.TrackSubscribed, (track, publication, participant) => {
        if (track.kind === 'video' && remoteVideoEl.value) {
          track.attach(remoteVideoEl.value);
        } else if (track.kind === 'audio') {
          // Audio track tự động attach vào thẻ audio ảo hoặc có thể tạo html element
          const audioElement = document.createElement('audio');
          track.attach(audioElement);
        }
      });

      newRoom.on(RoomEvent.TrackUnsubscribed, (track) => {
        track.detach();
      });

      newRoom.on(RoomEvent.Disconnected, () => {
        isConnected.value = false;
        cleanUp();
      });

      await newRoom.connect(url, token);
      room.value = newRoom;
      isConnected.value = true;

      // Publish local camera & mic
      await publishLocalTracks();

    } catch (err) {
      console.error('Lỗi kết nối LiveKit:', err);
      error.value = err.message || 'Không thể kết nối vào phòng.';
      cleanUp();
    }
  };

  const publishLocalTracks = async () => {
    if (!room.value) return;

    try {
      if (isCameraOn.value) {
        localVideoTrack = await createLocalVideoTrack();
        await room.value.localParticipant.publishTrack(localVideoTrack);
        if (localVideoEl.value) {
          localVideoTrack.attach(localVideoEl.value);
        }
      }

      if (isMicOn.value) {
        localAudioTrack = await createLocalAudioTrack();
        await room.value.localParticipant.publishTrack(localAudioTrack);
      }
    } catch (err) {
      console.error('Lỗi lấy thiết bị (Camera/Mic):', err);
      error.value = 'Không thể truy cập Camera hoặc Microphone.';
    }
  };

  const toggleMic = async () => {
    if (!room.value) return;
    try {
      if (isMicOn.value) {
        // Tắt mic
        if (localAudioTrack) {
          await room.value.localParticipant.unpublishTrack(localAudioTrack, true);
          localAudioTrack.stop();
          localAudioTrack = null;
        }
      } else {
        // Bật mic
        localAudioTrack = await createLocalAudioTrack();
        await room.value.localParticipant.publishTrack(localAudioTrack);
      }
      isMicOn.value = !isMicOn.value;
    } catch (err) {
      console.error('Lỗi toggle Mic:', err);
    }
  };

  const toggleCamera = async () => {
    if (!room.value) return;
    try {
      if (isCameraOn.value) {
        // Tắt camera
        if (localVideoTrack) {
          await room.value.localParticipant.unpublishTrack(localVideoTrack, true);
          localVideoTrack.detach();
          localVideoTrack.stop();
          localVideoTrack = null;
        }
      } else {
        // Bật camera
        localVideoTrack = await createLocalVideoTrack();
        await room.value.localParticipant.publishTrack(localVideoTrack);
        if (localVideoEl.value) {
          localVideoTrack.attach(localVideoEl.value);
        }
      }
      isCameraOn.value = !isCameraOn.value;
    } catch (err) {
      console.error('Lỗi toggle Camera:', err);
    }
  };

  const disconnect = () => {
    if (room.value) {
      room.value.disconnect();
    }
    cleanUp();
  };

  const cleanUp = () => {
    if (localVideoTrack) {
      localVideoTrack.detach();
      localVideoTrack.stop();
      localVideoTrack = null;
    }
    if (localAudioTrack) {
      localAudioTrack.stop();
      localAudioTrack = null;
    }
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
    disconnect
  };
}
