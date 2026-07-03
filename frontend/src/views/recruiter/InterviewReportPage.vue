<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { ArrowLeft, Download, Share2, CheckCircle, AlertTriangle, FileText } from 'lucide-vue-next'
import { authStore } from '../../stores/auth.store'
import { reportService } from '../../services/report.service'

const router = useRouter()
const route = useRoute()
const decision = ref('Chưa quyết định')
const note = ref('')
const toast = ref(null)
const loading = ref(true)
const report = ref(null)

onMounted(async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const interviewId = route.params.id
    const reportData = await reportService.getReport(companyId, interviewId)
    report.value = reportData
    // Map backend JSON to frontend structure
    if (reportData && reportData.report_json) {
      const parsed = JSON.parse(reportData.report_json)
      report.value.overall_score = parsed.final_score || reportData.final_score
      report.value.ai_recommendation = parsed.recommendation
      report.value.core_feedback = parsed.summary
      report.value.strengths = parsed.strengths || []
      report.value.weaknesses = parsed.weaknesses || []
      report.value.transcript_highlights = parsed.evidence_json || []
      report.value.ai_reasoning_summary = parsed.ai_reasoning_summary
      // Create some fake rubric scores based on overall score, as AI doesn't generate them yet
      report.value.rubric_scores = [
        { name: 'Overall AI Score', score: report.value.overall_score, color: 'var(--primary)' }
      ]
    }
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi tải báo cáo' }
  } finally {
    loading.value = false
  }
})

const handleSaveDecision = async () => {
  const companyId = authStore.user?.companies?.[0]?.id
  const interviewId = route.params.id
  await reportService.saveDecision(companyId, interviewId, { decision: decision.value, comment: note.value })
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
          <h1 class="text-h1">Báo cáo Phỏng vấn: {{ report?.candidate_name || 'Đang tải...' }}</h1>
          <p class="text-helper" style="margin-top: 4px">{{ report?.job_title }} • {{ report ? new Date(report.date).toLocaleDateString('vi-VN') : '' }}</p>
        </div>
      </div>
      <div style="display: flex; gap: 12px">
        <Button variant="secondary"><Share2 size="16" /> Chia sẻ</Button>
        <Button><Download size="16" /> Xuất PDF</Button>
      </div>
    </div>

    <div v-if="loading" style="text-align: center; padding: 40px; color: var(--text-muted)">
      Đang tổng hợp báo cáo AI...
    </div>
    <div v-else-if="!report" style="text-align: center; padding: 40px; color: var(--danger)">
      Không thể tải báo cáo.
    </div>
    <div v-else>
      <!-- Top Overview -->
      <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 24px">
        <Card style="padding: 20px; text-align: center; background-color: var(--surface)">
          <p class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; font-weight: 600">Điểm tổng quan</p>
          <p style="font-size: 36px; font-weight: 700; color: var(--primary)">{{ report.overall_score }}<span style="font-size: 18px; color: var(--text-muted)">/10</span></p>
        </Card>
        <Card style="padding: 20px; text-align: center; background-color: var(--surface)">
          <p class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; font-weight: 600">Đề xuất từ AI</p>
          <Badge :type="report.ai_recommendation === 'hire' ? 'success' : 'warning'" style="font-size: 14px; padding: 6px 16px; margin-top: 4px">
            {{ report.ai_recommendation === 'hire' ? 'Nên tuyển (Hire)' : report.ai_recommendation }}
          </Badge>
        </Card>
        <Card style="padding: 20px; grid-column: span 2; background-color: rgba(37, 99, 235, 0.05); border: 1px solid rgba(37, 99, 235, 0.1)">
          <p class="text-helper" style="text-transform: uppercase; margin-bottom: 8px; font-weight: 600; color: var(--primary)">Nhận xét cốt lõi</p>
          <p class="text-body" style="font-weight: 500">{{ report.core_feedback }}</p>
        </Card>
      </div>

      <div style="display: grid; grid-template-columns: 1fr 2fr; gap: 24px">
        
        <!-- Left: Rubric Scores -->
        <div style="display: flex; flex-direction: column; gap: 24px">
          <Card title="Điểm chi tiết (Rubric)">
            <div style="display: flex; flex-direction: column; gap: 20px">
              <div v-for="item in report.rubric_scores" :key="item.name">
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
            <textarea class="input-field" v-model="note" placeholder="Nhập ghi chú HR..." style="width: 100%; height: 100px; margin-bottom: 16px"></textarea>
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
                  <li v-for="(str, idx) in report.strengths" :key="idx">{{ str }}</li>
                </ul>
              </div>
              <div>
                <h4 class="text-body" style="display: flex; align-items: center; gap: 8px; color: var(--warning); font-weight: 600; margin-bottom: 12px">
                  <AlertTriangle size="18" /> Điểm rủi ro (Cần lưu ý)
                </h4>
                <ul class="text-body" style="padding-left: 24px; display: flex; flex-direction: column; gap: 12px">
                  <li v-for="(weak, idx) in report.weaknesses" :key="idx">{{ weak }}</li>
                </ul>
              </div>
            </div>
          </Card>

          <Card title="Trích xuất Transcript (Bằng chứng)">
            <div style="display: flex; flex-direction: column; gap: 16px">
              
              <div v-for="(hl, idx) in report.transcript_highlights" :key="idx" 
                   :style="{ padding: '16px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)', borderLeft: `4px solid var(--${hl.type})` }">
                <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px">
                  <Badge :type="hl.type">{{ hl.badge }}</Badge>
                  <span class="text-helper">{{ hl.time }}</span>
                </div>
                <p class="text-body" style="font-style: italic; color: var(--text-secondary)">
                  "{{ hl.text }}"
                </p>
              </div>

            </div>
            <Button variant="ghost" style="width: 100%; margin-top: 16px"><FileText size="16" /> Xem toàn bộ Transcript</Button>
          </Card>

        </div>
      </div>
    </div>
  </div>
</template>
