<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Button from '../../components/common/AppButton.vue'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import { Mic, MicOff, Video, VideoOff, MonitorUp, MessageSquare, PhoneOff, Sparkles, CheckCircle, AlertTriangle } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()

const micOn = ref(true)
const cameraOn = ref(true)
const activeTab = ref('assistant')
const transcript = ref([
  { speaker: 'Candidate', text: 'Vâng, em đã sử dụng React khoảng 3 năm trong các dự án thực tế.', time: '14:02' }
])

const showEndModal = ref(false)
const isEnding = ref(false)
const entryToast = ref(history.state?.message ? { type: 'success', message: history.state.message } : null)

let timer = null

onMounted(() => {
  if (history.state?.message) {
    window.history.replaceState({}, document.title)
  }
  
  timer = setTimeout(() => {
    transcript.value.push({ speaker: 'AI', text: '[AI Phân tích] Câu trả lời khá tự tin, tuy nhiên chưa nêu rõ dự án cụ thể.', time: '14:03', isAi: true })
  }, 4000)
})

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})

const handleEndCall = () => {
  showEndModal.value = true
}

const confirmEndCall = () => {
  showEndModal.value = false
  isEnding.value = true
  setTimeout(() => {
    router.push({ path: '/reports', state: { message: 'Đã lưu kết quả phỏng vấn thành công' } })
  }, 1500)
}
</script>

<template>
  <div style="height: 100vh; display: flex; flex-direction: column; background-color: var(--background)">
    <Toast v-if="entryToast" :type="entryToast.type" :message="entryToast.message" @close="entryToast = null" />
    <Toast v-if="isEnding" type="success" message="Đã kết thúc phỏng vấn. Đang lưu kết quả..." :duration="1500" />

    <!-- Header -->
    <div style="height: 64px; background-color: var(--surface); border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 24px">
      <div style="display: flex; align-items: center; gap: 16px">
        <div>
          <h1 class="text-h2">Frontend Developer - Nguyễn Văn A</h1>
          <div style="display: flex; gap: 12px; margin-top: 4px">
            <span class="text-helper" style="color: var(--danger); display: flex; align-items: center; gap: 4px">
              <div style="width: 8px; height: 8px; border-radius: 50%; background-color: var(--danger)"></div> Đang ghi âm & Transcript
            </span>
            <span class="text-helper" style="color: var(--primary)">00:15:32</span>
          </div>
        </div>
      </div>
      <div style="display: flex; align-items: center; gap: 12px">
        <Badge type="info"><Sparkles size="12" style="margin-right: 4px" /> AI Active</Badge>
        <Button variant="secondary" @click="router.push('/dashboard')">Rời tạm thời</Button>
      </div>
    </div>

    <!-- Main Content -->
    <div style="display: flex; flex: 1; overflow: hidden">
      
      <!-- Left: Video Area (70%) -->
      <div style="flex: 7; display: flex; flex-direction: column; padding: 24px; gap: 16px">
        <!-- Video Grid -->
        <div style="flex: 1; display: flex; gap: 16px; position: relative">
          <!-- Candidate Video (Large) -->
          <div style="flex: 1; background-color: #1E293B; border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden">
            <div style="text-align: center; color: white">
              <div style="width: 80px; height: 80px; border-radius: 50%; background-color: rgba(255,255,255,0.1); display: flex; align-items: center; justify-content: center; margin: 0 auto 16px; font-size: 24px; font-weight: bold">N</div>
              <div style="font-size: 16px">Candidate Camera Feed</div>
            </div>
            <div style="position: absolute; bottom: 16px; left: 16px; background-color: rgba(0,0,0,0.6); color: white; padding: 4px 12px; border-radius: var(--radius); font-size: 13px">
              Nguyễn Văn A (Ứng viên)
            </div>
          </div>
          
          <!-- Recruiter Video (Small/PiP or Side) -->
          <div style="width: 240px; background-color: #334155; border-radius: var(--radius-lg); display: flex; align-items: center; justify-content: center; position: absolute; top: 16px; right: 16px; height: 160px; box-shadow: var(--shadow-md); border: 2px solid rgba(255,255,255,0.1)">
             <div style="color: white; font-size: 12px">Your Camera</div>
          </div>
        </div>

        <!-- Control Bar -->
        <div style="height: 72px; background-color: var(--surface); border-radius: var(--radius-lg); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; gap: 16px">
          <Button :variant="micOn ? 'secondary' : 'primary'" style="width: 48px; height: 48px; border-radius: 50%; padding: 0" @click="micOn = !micOn">
            <Mic v-if="micOn" size="20" />
            <MicOff v-else size="20" />
          </Button>
          <Button :variant="cameraOn ? 'secondary' : 'primary'" style="width: 48px; height: 48px; border-radius: 50%; padding: 0" @click="cameraOn = !cameraOn">
            <Video v-if="cameraOn" size="20" />
            <VideoOff v-else size="20" />
          </Button>
          <Button variant="secondary" style="width: 48px; height: 48px; border-radius: 50%; padding: 0">
            <MonitorUp size="20" />
          </Button>
          <Button variant="secondary" style="width: 48px; height: 48px; border-radius: 50%; padding: 0">
            <MessageSquare size="20" />
          </Button>
          <div style="width: 1px; height: 32px; background-color: var(--border); margin: 0 8px"></div>
          <Button style="background-color: var(--danger); color: white; border: none; height: 48px; padding: 0 24px; border-radius: 24px" @click="handleEndCall">
            <PhoneOff size="20" /> Kết thúc
          </Button>
        </div>
      </div>

      <!-- Right: AI Panel (30%) -->
      <div style="flex: 3; background-color: var(--surface); border-left: 1px solid var(--border); display: flex; flex-direction: column">
        <!-- Tabs -->
        <div style="display: flex; border-bottom: 1px solid var(--border)">
          <div :class="['nav-item', activeTab === 'assistant' ? 'active' : '']" style="flex: 1; justify-content: center; cursor: pointer; padding: 16px 0" @click="activeTab = 'assistant'">
            <Sparkles size="16" /> AI Assistant
          </div>
          <div :class="['nav-item', activeTab === 'rubric' ? 'active' : '']" style="flex: 1; justify-content: center; cursor: pointer; padding: 16px 0" @click="activeTab = 'rubric'">
            <CheckCircle size="16" /> Tiêu chí
          </div>
          <div :class="['nav-item', activeTab === 'transcript' ? 'active' : '']" style="flex: 1; justify-content: center; cursor: pointer; padding: 16px 0" @click="activeTab = 'transcript'">
            <MessageSquare size="16" /> Transcript
          </div>
        </div>

        <!-- Tab Content -->
        <div style="flex: 1; overflow-y: auto; padding: 24px">
          <div v-if="activeTab === 'assistant'" style="display: flex; flex-direction: column; gap: 24px">
            <Card style="background-color: rgba(8, 145, 178, 0.05); border-color: rgba(8, 145, 178, 0.2)">
              <div style="display: flex; gap: 8px; align-items: flex-start; color: var(--accent)">
                <Sparkles size="18" style="margin-top: 2px" />
                <div>
                  <h4 class="text-body" style="font-weight: 600; margin-bottom: 8px">Phân tích Real-time</h4>
                  <div style="display: flex; flex-direction: column; gap: 8px">
                    <div style="display: flex; justify-content: space-between">
                      <span class="text-helper">Mức độ tự tin:</span>
                      <span class="text-helper" style="font-weight: 600; color: var(--success)">Tốt (85%)</span>
                    </div>
                    <div style="display: flex; justify-content: space-between">
                      <span class="text-helper">Độ chi tiết:</span>
                      <span class="text-helper" style="font-weight: 600; color: var(--warning)">Cần đào sâu hơn</span>
                    </div>
                  </div>
                </div>
              </div>
            </Card>

            <div>
              <h3 class="text-body" style="font-weight: 600; margin-bottom: 16px">Gợi ý câu hỏi tiếp theo</h3>
              <div style="display: flex; flex-direction: column; gap: 12px">
                <div style="padding: 12px; border: 1px solid var(--border); border-radius: var(--radius); background-color: var(--surface)">
                  <p class="text-body" style="margin-bottom: 12px">"Bạn có thể kể một ví dụ cụ thể về việc tối ưu performance trong dự án React bạn vừa nhắc đến không?"</p>
                  <Button variant="secondary" style="width: 100%; font-size: 13px; height: 32px">Hỏi ngay</Button>
                </div>
                <div style="padding: 12px; border: 1px solid var(--border); border-radius: var(--radius); background-color: var(--surface)">
                  <p class="text-body" style="margin-bottom: 12px">"Trong dự án đó, bạn xử lý state management như thế nào? Dùng Redux hay Context?"</p>
                  <Button variant="secondary" style="width: 100%; font-size: 13px; height: 32px">Hỏi ngay</Button>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeTab === 'rubric'" style="display: flex; flex-direction: column; gap: 16px">
            <div>
              <div style="display: flex; justify-content: space-between; margin-bottom: 8px">
                <span class="text-body" style="font-weight: 500">Technical Knowledge (React)</span>
                <span class="text-body" style="color: var(--text-muted)">--/10</span>
              </div>
              <div style="height: 6px; background-color: var(--surface-soft); border-radius: 3px"></div>
              <p class="text-helper" style="margin-top: 4px">AI đang thu thập thêm bằng chứng...</p>
            </div>
            <div>
              <div style="display: flex; justify-content: space-between; margin-bottom: 8px">
                <span class="text-body" style="font-weight: 500">Communication</span>
                <span class="text-body" style="color: var(--success)">8/10</span>
              </div>
              <div style="height: 6px; background-color: var(--surface-soft); border-radius: 3px">
                <div style="height: 100%; width: 80%; background-color: var(--success); border-radius: 3px"></div>
              </div>
              <p class="text-helper" style="margin-top: 4px; color: var(--text-muted)">Bằng chứng: Trả lời lưu loát, không vấp váp.</p>
            </div>
          </div>

          <div v-if="activeTab === 'transcript'" style="display: flex; flex-direction: column; gap: 16px">
            <div v-for="(item, idx) in transcript" :key="idx" :style="{ padding: '12px', backgroundColor: item.isAi ? 'rgba(8, 145, 178, 0.05)' : 'var(--surface-soft)', borderRadius: 'var(--radius)', border: item.isAi ? '1px dashed rgba(8, 145, 178, 0.3)' : 'none' }">
              <div style="display: flex; justify-content: space-between; margin-bottom: 4px">
                <span class="text-helper" :style="{ fontWeight: 600, color: item.isAi ? 'var(--accent)' : 'var(--primary)' }">
                  <Sparkles v-if="item.isAi" size="12" style="margin-right: 4px; vertical-align: middle" />
                  {{ item.speaker }}
                </span>
                <span class="text-helper" style="font-size: 11px">{{ item.time }}</span>
              </div>
              <p class="text-body" :style="{ color: item.isAi ? 'var(--accent)' : 'var(--text-main)', fontStyle: item.isAi ? 'italic' : 'normal' }">
                {{ item.text }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Modal :isOpen="showEndModal" @close="showEndModal = false" title="Kết thúc phỏng vấn">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn kết thúc buổi phỏng vấn này?</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEndModal = false">Huỷ</Button>
        <Button variant="primary" @click="confirmEndCall" style="background-color: var(--danger); color: white; border-color: var(--danger)">OK</Button>
      </div>
    </Modal>
  </div>
</template>
