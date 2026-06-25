<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { ArrowLeft, Download, Share2, CheckCircle, AlertTriangle, FileText } from 'lucide-vue-next'

const router = useRouter()
const decision = ref('Chưa quyết định')
const toast = ref(null)

const handleSaveDecision = () => {
  toast.value = { type: 'success', message: 'Đã lưu quyết định tuyển dụng thành công!' }
}
</script>

<template>
  <div style="max-width: 900px; margin: 0 auto; padding-bottom: 64px">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 32px">
      <div style="display: flex; align-items: center; gap: 16px">
        <Button variant="ghost" @click="router.push('/interviews')" style="padding: 8px">
          <ArrowLeft size="20" />
        </Button>
        <div>
          <h1 class="text-h1">Báo cáo Phỏng vấn: Nguyễn Văn A</h1>
          <p class="text-helper" style="margin-top: 4px">Frontend Developer • 14/10/2023</p>
        </div>
      </div>
      <div style="display: flex; gap: 12px">
        <Button variant="secondary"><Share2 size="16" /> Chia sẻ</Button>
        <Button><Download size="16" /> Xuất PDF</Button>
      </div>
    </div>

    <!-- Top Overview -->
    <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 24px">
      <Card style="padding: 20px; text-align: center; background-color: var(--surface)">
        <p class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; font-weight: 600">Điểm tổng quan</p>
        <p style="font-size: 36px; font-weight: 700; color: var(--primary)">8.5<span style="font-size: 18px; color: var(--text-muted)">/10</span></p>
      </Card>
      <Card style="padding: 20px; text-align: center; background-color: var(--surface)">
        <p class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; font-weight: 600">Đề xuất từ AI</p>
        <Badge type="success" style="font-size: 14px; padding: 6px 16px; margin-top: 4px">Nên tuyển (Hire)</Badge>
      </Card>
      <Card style="padding: 20px; grid-column: span 2; background-color: rgba(37, 99, 235, 0.05); border: 1px solid rgba(37, 99, 235, 0.1)">
        <p class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; font-weight: 600; color: var(--primary)">Nhận xét cốt lõi</p>
        <p class="text-body" style="font-weight: 500">Ứng viên có kiến thức chuyên môn sâu về React, thái độ giao tiếp chuyên nghiệp. Phù hợp văn hóa công ty.</p>
      </Card>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 2fr; gap: 24px">
      
      <!-- Left: Rubric Scores -->
      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card title="Điểm chi tiết (Rubric)">
          <div style="display: flex; flex-direction: column; gap: 20px">
            <div v-for="item in [
              { name: 'Technical Knowledge', score: 9, color: 'var(--success)' },
              { name: 'Problem Solving', score: 8, color: 'var(--primary)' },
              { name: 'Communication', score: 9, color: 'var(--success)' },
              { name: 'Experience Relevance', score: 7, color: 'var(--warning)' },
              { name: 'Culture Fit', score: 8.5, color: 'var(--primary)' }
            ]" :key="item.name">
              <div style="display: flex; justify-content: space-between; margin-bottom: 8px">
                <span class="text-body" style="font-weight: 500">{{ item.name }}</span>
                <span class="text-body" style="font-weight: 600">{{ item.score }}/10</span>
              </div>
              <div style="height: 8px; background-color: var(--surface-soft); border-radius: 4px; overflow: hidden">
                <div :style="{ height: '100%', width: `${item.score * 10}%`, backgroundColor: item.color, borderRadius: '4px' }"></div>
              </div>
            </div>
          </div>
        </Card>
        
        <Card title="Quyết định của Bạn">
          <select class="input-field" v-model="decision" style="width: 100%; margin-bottom: 16px">
            <option value="Chưa quyết định">-- Chọn quyết định --</option>
            <option value="offer">Gửi Offer</option>
            <option value="reject">Từ chối</option>
            <option value="next_round">Phỏng vấn vòng sau</option>
          </select>
          <textarea class="input-field" placeholder="Nhập ghi chú HR..." style="width: 100%; height: 100px; margin-bottom: 16px"></textarea>
          <Button style="width: 100%" @click="handleSaveDecision">Lưu quyết định</Button>
        </Card>
      </div>

      <!-- Right: Insights & Evidence -->
      <div style="display: flex; flex-direction: column; gap: 24px">
        
        <Card title="Phân tích Điểm mạnh & Rủi ro" style="border-top: 4px solid var(--accent)">
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 24px">
            <div>
              <h4 class="text-body" style="display: flex; align-items: center; gap: 8px; color: var(--success); font-weight: 600; margin-bottom: 12px">
                <CheckCircle size="18" /> Điểm mạnh
              </h4>
              <ul class="text-body" style="padding-left: 24px; display: flex; flex-direction: column; gap: 12px">
                <li>Nắm rất vững vòng đời (Lifecycle) và Hooks trong React.</li>
                <li>Cách trình bày vấn đề có cấu trúc rõ ràng (STAR method).</li>
              </ul>
            </div>
            <div>
              <h4 class="text-body" style="display: flex; align-items: center; gap: 8px; color: var(--warning); font-weight: 600; margin-bottom: 12px">
                <AlertTriangle size="18" /> Điểm rủi ro (Cần lưu ý)
              </h4>
              <ul class="text-body" style="padding-left: 24px; display: flex; flex-direction: column; gap: 12px">
                <li>Chưa có nhiều kinh nghiệm lead team trực tiếp.</li>
                <li>Mức lương đề xuất hơi cao so với budget ban đầu.</li>
              </ul>
            </div>
          </div>
        </Card>

        <Card title="Trích xuất Transcript (Bằng chứng)">
          <div style="display: flex; flex-direction: column; gap: 16px">
            
            <div style="padding: 16px; background-color: var(--surface-soft); border-radius: var(--radius); border-left: 4px solid var(--success)">
              <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px">
                <Badge type="success">Bằng chứng: Technical</Badge>
                <span class="text-helper">14:15:20</span>
              </div>
              <p class="text-body" style="font-style: italic; color: var(--text-secondary)">
                "...Trong dự án đó, khi list data lên tới 10,000 items, em đã sử dụng kỹ thuật Virtualization (react-window) kết hợp với useMemo để tránh re-render những components không cần thiết, giúp FPS tăng từ 20 lên 60."
              </p>
            </div>

            <div style="padding: 16px; background-color: var(--surface-soft); border-radius: var(--radius); border-left: 4px solid var(--warning)">
              <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px">
                <Badge type="warning">Bằng chứng: Leadership</Badge>
                <span class="text-helper">14:30:10</span>
              </div>
              <p class="text-body" style="font-style: italic; color: var(--text-secondary)">
                "...Thường thì em sẽ nhận task từ Tech Lead và tự implement độc lập. Thỉnh thoảng em có review code cho các bạn Junior nhưng chưa chính thức lead dự án nào."
              </p>
            </div>

          </div>
          <Button variant="ghost" style="width: 100%; margin-top: 16px"><FileText size="16" /> Xem toàn bộ Transcript</Button>
        </Card>

      </div>
    </div>
  </div>
</template>
