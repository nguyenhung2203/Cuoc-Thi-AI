<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { mockApi } from '../../utils/mockData'
import { ArrowLeft, Video, Copy, Calendar, Clock, User, Briefcase, Mail } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const interview = ref(null)
const loading = ref(true)
const copied = ref(false)
const toast = ref(null)

onMounted(async () => {
  const data = await mockApi.interviews.getById(id)
  interview.value = data || null
  loading.value = false
})

const copyLink = () => {
  if (interview.value && interview.value.link) {
    navigator.clipboard.writeText(interview.value.link)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  }
}

const handleSaveNotes = () => {
  toast.value = { type: 'success', message: 'Ghi chú đã được lưu thành công!' }
}

const dateObj = computed(() => interview.value ? new Date(interview.value.datetime) : null)
</script>

<template>
  <div style="max-width: 1000px; margin: 0 auto; padding-bottom: 40px">
    <Toast v-if="copied" type="success" message="Link phòng phỏng vấn đã được copy!" @close="copied = false" />
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <div v-if="loading" style="padding: 32px; text-align: center; color: var(--text-muted)">Đang tải thông tin phỏng vấn...</div>
    <div v-else-if="!interview" style="padding: 32px; text-align: center; color: var(--warning)">Không tìm thấy thông tin phỏng vấn.</div>
    <div v-else>
      <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
        <Button variant="ghost" @click="router.push('/interviews')" style="padding: 8px">
          <ArrowLeft size="20" />
        </Button>
        <div>
          <div style="display: flex; align-items: center; gap: 12px">
            <h1 class="text-h1">Chi tiết buổi phỏng vấn</h1>
            <Badge :type="interview.status === 'Scheduled' ? 'info' : 'neutral'">
              {{ interview.status === 'Scheduled' ? 'Đã lên lịch' : interview.status }}
            </Badge>
          </div>
          <p class="text-helper" style="margin-top: 4px">Lịch phỏng vấn > {{ interview.candidateName }}</p>
        </div>
      </div>

      <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
        <div style="display: flex; flex-direction: column; gap: 24px">
          <Card title="Thông tin chung">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px">
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><User size="14" /> Ứng viên</p>
                <p class="text-body" style="font-weight: 500">{{ interview.candidateName }}</p>
              </div>
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><Briefcase size="14" /> Vị trí ứng tuyển</p>
                <p class="text-body" style="font-weight: 500">{{ interview.jobTitle }}</p>
              </div>
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><Calendar size="14" /> Ngày phỏng vấn</p>
                <p class="text-body" style="font-weight: 500">{{ dateObj.toLocaleDateString('vi-VN') }}</p>
              </div>
              <div>
                <p class="text-helper" style="margin-bottom: 4px; display: flex; align-items: center; gap: 6px"><Clock size="14" /> Thời gian</p>
                <p class="text-body" style="font-weight: 500">{{ dateObj.toLocaleTimeString('vi-VN', {hour: '2-digit', minute:'2-digit'}) }}</p>
              </div>
            </div>
            
            <div style="margin-top: 24px; padding-top: 24px; border-top: 1px solid var(--border)">
              <p class="text-helper" style="margin-bottom: 12px; font-weight: 500; color: var(--text-main)">Hành động</p>
              <div style="display: flex; gap: 12px">
                <Button @click="router.push({ path: '/recruiter-room', state: { message: 'Vào phòng phỏng vấn thành công!' } })">
                  <Video size="16" /> Vào phòng phỏng vấn
                </Button>
                <Button variant="secondary" @click="copyLink">
                  <Copy size="16" /> Copy Link Invite
                </Button>
                <Button variant="ghost">
                  <Mail size="16" /> Gửi email nhắc nhở
                </Button>
              </div>
            </div>
          </Card>
          
          <Card title="Cấu hình AI Assistant" style="border-top: 4px solid var(--accent)">
            <ul class="text-body" style="padding-left: 20px; display: flex; flex-direction: column; gap: 12px">
              <li><strong>Phân tích realtime:</strong> Bật</li>
              <li><strong>Tạo transcript:</strong> Bật</li>
              <li><strong>Gợi ý câu hỏi (Rubric-based):</strong> Bật</li>
              <li><strong>Tự động chấm điểm:</strong> Bật</li>
            </ul>
          </Card>
        </div>

        <div>
          <Card title="Ghi chú nội bộ">
            <textarea 
              class="input-field" 
              rows="6"
              placeholder="Nhập ghi chú hoặc nhắc nhở trước buổi phỏng vấn (Chỉ recruiter xem được)..."
              style="width: 100%; margin-bottom: 16px"
            ></textarea>
            <Button variant="secondary" style="width: 100%" @click="handleSaveNotes">Lưu ghi chú</Button>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>
