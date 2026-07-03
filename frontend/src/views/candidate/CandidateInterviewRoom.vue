<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, MicOff, Video, VideoOff, MessageSquare, PhoneOff, CheckCircle, FileText } from 'lucide-vue-next'
import { useLiveKit } from '../../composables/useLiveKit'

const router = useRouter()
const route = useRoute()

const {
  isConnected, error, isMicOn, isCameraOn,
  localVideoEl, remoteVideoEl,
  connectToRoom, toggleMic, toggleCamera, disconnect
} = useLiveKit()

const activeTab = ref('info')

const showLeaveModal = ref(false)
const isLeaving = ref(false)

const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)
const interviewInfo = ref(history.state?.interviewInfo || null)

onMounted(async () => {
  if (history.state?.message) {
    window.history.replaceState({ interviewInfo: history.state.interviewInfo }, document.title)
  }
  
  // Logic kết nối LiveKit
  const token = route.query.token
  if (token) {
    try {
      const livekitUrl = import.meta.env.VITE_LIVEKIT_URL || 'ws://localhost:7880'
      await connectToRoom(livekitUrl, token)
    } catch (err) {
      console.error('Không thể vào phòng', err)
    }
  }
})

onUnmounted(() => {
  disconnect()
})

const handleLeave = () => {
  showLeaveModal.value = true
}

const confirmLeave = () => {
  showLeaveModal.value = false
  isLeaving.value = true
  setTimeout(() => {
    router.push({ path: '/home', state: { message: 'Rời phòng phỏng vấn thành công' } })
  }, 1500)
}
</script>

<template>
  <div style="height: 100vh; display: flex; flex-direction: column; background-color: var(--background)">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    <Toast v-if="isLeaving" type="info" message="Đang rời phòng phỏng vấn..." :duration="1500" />

    <!-- Header -->
    <div style="height: 64px; background-color: var(--surface); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 24px">
      <div style="display: flex; align-items: center; gap: 16px">
        <div>
          <h1 class="text-h2">Phỏng vấn: {{ interviewInfo?.job_title || 'Vị trí ứng tuyển' }}</h1>
          <p class="text-helper" style="margin-top: 4px">Công ty {{ interviewInfo?.company_name || 'Tuyển dụng' }}</p>
        </div>
      </div>
      <div style="display: flex; align-items: center; gap: 12px">
        <span class="text-helper" style="color: var(--success); display: flex; align-items: center; gap: 4px">
          <CheckCircle size="14" /> Đường truyền tốt
        </span>
        <span style="width: 1px; height: 24px; background-color: var(--border)"></span>
        <span class="text-body" style="font-weight: 600">00:15:32</span>
      </div>
    </div>

    <!-- Main Content -->
    <div style="display: flex; flex: 1; overflow: hidden">
      
      <!-- Left: Video Area (70%) -->
      <div style="flex: 7; display: flex; flex-direction: column; padding: 24px; gap: 16px">
        
        <!-- Notification Banner -->
        <div style="background-color: rgba(37, 99, 235, 0.05); border: 1px solid rgba(37, 99, 235, 0.2); padding: 12px 16px; border-radius: var(--radius); display: flex; align-items: center; gap: 8px; color: var(--primary); font-size: 13px">
          <span style="font-weight: 600">Lời nhắc:</span> Hãy trả lời tự nhiên. Nhà tuyển dụng sẽ dẫn dắt buổi phỏng vấn.
        </div>

        <!-- Video Grid -->
        <div style="flex: 1; display: flex; gap: 16px; position: relative">
          <!-- Recruiter Video (Large) -->
          <div style="flex: 1; background-color: #1E293B; border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden">
            <!-- Thẻ video thật cho LiveKit (Remote - Recruiter) -->
            <video ref="remoteVideoEl" autoplay playsinline style="width: 100%; height: 100%; object-fit: cover; position: absolute; top: 0; left: 0;"></video>
            
            <div style="text-align: center; color: white; position: relative; z-index: 1;">
              <div style="width: 80px; height: 80px; border-radius: 50%; background-color: rgba(255,255,255,0.1); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; font-size: 24px; font-weight: bold">R</div>
              <div style="font-size: 16px">Đang chờ nhà tuyển dụng...</div>
            </div>
            
            <div style="position: absolute; bottom: 16px; left: 16px; background-color: rgba(0,0,0,0.6); color: white; padding: 4px 12px; border-radius: var(--radius); font-size: 13px; z-index: 2">
              Nhà Tuyển Dụng
            </div>
          </div>
          
          <!-- Candidate Self View (Small/PiP or Side) -->
          <div style="width: 240px; background-color: #334155; border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; position: absolute; top: 16px; right: 16px; height: 160px; box-shadow: var(--shadow-md); border: 2px solid rgba(255,255,255,0.1); overflow: hidden; z-index: 10">
             <!-- Thẻ video thật cho LiveKit (Local) -->
             <video ref="localVideoEl" autoplay playsinline muted style="width: 100%; height: 100%; object-fit: cover; position: absolute; top: 0; left: 0; transform: scaleX(-1);"></video>
             <div style="color: white; font-size: 12px; position: relative; z-index: 1">Your Camera</div>
          </div>
        </div>

        <!-- Control Bar -->
        <div style="height: 72px; background-color: var(--surface); border-radius: var(--radius-lg); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; gap: 16px">
          <Button :variant="isMicOn ? 'secondary' : 'primary'" style="width: 48px; height: 48px; border-radius: 50%; padding: 0" @click="toggleMic">
            <Mic v-if="isMicOn" size="20" />
            <MicOff v-else size="20" />
          </Button>
          <Button :variant="isCameraOn ? 'secondary' : 'primary'" style="width: 48px; height: 48px; border-radius: 50%; padding: 0" @click="toggleCamera">
            <Video v-if="isCameraOn" size="20" />
            <VideoOff v-else size="20" />
          </Button>
          <div style="width: 1px; height: 32px; background-color: var(--border); margin: 0 8px"></div>
          <Button style="background-color: var(--danger); color: white; border: none; height: 48px; padding: 0 24px; border-radius: 24px" @click="handleLeave">
            <PhoneOff size="20" /> Rời phòng
          </Button>
        </div>
      </div>

      <!-- Right: Candidate Panel (30%) -->
      <div style="flex: 3; background-color: var(--surface); border-left: 1px solid var(--border); display: flex; flex-direction: column">
        <!-- Tabs -->
        <div style="display: flex; border-bottom: 1px solid var(--border)">
          <div :class="['nav-item', activeTab === 'info' ? 'active' : '']" style="flex: 1; justify-content: center; cursor: pointer; padding: 16px 0" @click="activeTab = 'info'">
            <FileText size="16" /> Thông tin JD
          </div>
          <div :class="['nav-item', activeTab === 'chat' ? 'active' : '']" style="flex: 1; justify-content: center; cursor: pointer; padding: 16px 0" @click="activeTab = 'chat'">
            <MessageSquare size="16" /> Chat
          </div>
        </div>

        <!-- Tab Content -->
        <div style="flex: 1; overflow-y: auto; padding: 24px">
          <div v-if="activeTab === 'info'" style="display: flex; flex-direction: column; gap: 24px">
            <div>
              <h3 class="text-body" style="font-weight: 600; margin-bottom: 8px">Mô tả công việc (Tóm tắt)</h3>
              <ul class="text-body" style="padding-left: 20px; color: var(--text-secondary); display: flex; flex-direction: column; gap: 8px">
                <li>Phát triển các tính năng Frontend sử dụng ReactJS.</li>
                <li>Tối ưu hóa hiệu suất ứng dụng web.</li>
                <li>Phối hợp với UI/UX designer và Backend developer.</li>
              </ul>
            </div>
            
            <div style="border-top: 1px solid var(--border); padding-top: 24px">
              <h3 class="text-body" style="font-weight: 600; margin-bottom: 8px">Ghi chú cá nhân của bạn</h3>
              <textarea 
                class="input-field" 
                placeholder="Viết nháp câu trả lời hoặc ghi chú tại đây (Người phỏng vấn không thấy)..."
                style="width: 100%; height: 200px; background-color: var(--surface-soft)"
              ></textarea>
            </div>
          </div>

          <div v-if="activeTab === 'chat'" style="height: 100%; display: flex; flex-direction: column">
            <div style="flex: 1; display: flex; flex-direction: column; justify-content: flex-end; padding-bottom: 16px">
              <p class="text-helper" style="text-align: center">Chưa có tin nhắn nào.</p>
            </div>
            <div style="display: flex; gap: 8px">
              <input type="text" class="input-field" placeholder="Nhập tin nhắn..." style="flex: 1" />
              <Button variant="primary">Gửi</Button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Modal :isOpen="showLeaveModal" @close="showLeaveModal = false" title="Xác nhận rời phòng">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn rời khỏi phòng phỏng vấn?</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showLeaveModal = false">Huỷ</Button>
        <Button variant="primary" @click="confirmLeave" style="background-color: var(--danger); color: white; border-color: var(--danger)">OK</Button>
      </div>
    </Modal>
  </div>
</template>
