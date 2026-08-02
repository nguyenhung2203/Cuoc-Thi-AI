<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { ArrowLeft, Download, Share2, CheckCircle, AlertTriangle, FileText, Sparkles, Save } from 'lucide-vue-next'
import { authStore } from '../../stores/auth.store'
import { reportService } from '../../services/report.service'
import { transcriptService } from '../../services/transcript.service'
import { isOneOf, maxLength, minLength, normalizeText, requiredTrim } from '../../utils/validators.js'

const router = useRouter()
const route = useRoute()
const decision = ref('Chưa quyết định')
const note = ref('')
const toast = ref(null)
const loading = ref(true)
const report = ref(null)

const showTranscriptModal = ref(false)
const fullTranscripts = ref([])
const loadingTranscripts = ref(false)
const retryingReport = ref(false)

let reportPollTimer = null
const loadReport = async () => {
  const companyId = authStore.user?.companies?.[0]?.id
  const interviewId = route.params.id
  const reportData = await reportService.getReport(companyId, interviewId)
  report.value = reportData
  if (reportData?.report_json) {
    const parsed = JSON.parse(reportData.report_json)
    report.value.overall_score = parsed.final_score || reportData.final_score
    report.value.ai_recommendation = parsed.recommendation
    report.value.core_feedback = parsed.summary
    report.value.strengths = parsed.strengths || []
    report.value.weaknesses = parsed.weaknesses || []
    report.value.transcript_highlights = parsed.evidence_json || []
    report.value.ai_reasoning_summary = parsed.ai_reasoning_summary
    report.value.rubric_scores = parsed.scores || []
  }
  return reportData
}

onMounted(async () => {
  try {
    const data = await loadReport()
    if (data?.status === 'generating' || data?.status === 'pending') {
      reportPollTimer = setInterval(async () => {
        try {
          const latest = await loadReport()
          if (latest?.status === 'ready' || latest?.status === 'failed') {
            clearInterval(reportPollTimer)
            reportPollTimer = null
          }
        } catch (error) {
          console.error('Lỗi cập nhật trạng thái báo cáo', error)
        }
      }, 3000)
    }
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi tải báo cáo' }
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (reportPollTimer) clearInterval(reportPollTimer)
})

const savingDecision = ref(false)
const handleSaveDecision = async () => {
  if (savingDecision.value) return
  const allowed = ['offer', 'reject', 'next_round']
  const normalizedNote = normalizeText(note.value)
  let validationError = requiredTrim(decision.value === 'Chưa quyết định' ? '' : decision.value, 'Vui lòng chọn quyết định tuyển dụng.')
  validationError ||= isOneOf(decision.value, allowed, 'Quyết định tuyển dụng không hợp lệ.')
  validationError ||= maxLength(normalizedNote, 5000, 'Ghi chú không được vượt quá 5.000 ký tự.')
  if (decision.value === 'reject') validationError ||= minLength(normalizedNote, 10, 'Khi từ chối, ghi chú phải có ít nhất 10 ký tự.')
  if (validationError) {
    toast.value = { type: 'error', message: validationError }
    return
  }
  savingDecision.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const interviewId = route.params.id
    await reportService.saveDecision(companyId, interviewId, { decision: decision.value, comment: normalizedNote })
    note.value = normalizedNote
    toast.value = { type: 'success', message: 'Đã lưu quyết định tuyển dụng thành công!' }
  } catch (error) {
    toast.value = { type: 'error', message: error?.message || 'Không thể lưu quyết định tuyển dụng.' }
  } finally {
    savingDecision.value = false
  }
}

const handleViewTranscripts = async () => {
  showTranscriptModal.value = true
  if (fullTranscripts.value.length === 0) {
    loadingTranscripts.value = true
    try {
      const companyId = authStore.user?.companies?.[0]?.id
      const interviewId = route.params.id
      const data = await transcriptService.getTranscripts(companyId, interviewId)
      fullTranscripts.value = data || []
    } catch (err) {
      toast.value = { type: 'error', message: 'Lỗi tải transcript' }
    } finally {
      loadingTranscripts.value = false
    }
  }
}

const handleShare = async () => {
  const url = window.location.href
  const title = `Báo cáo phỏng vấn: ${report.value?.candidate_name || ''}`
  try {
    if (navigator.share) {
      await navigator.share({ title, url })
      return
    }
    await navigator.clipboard.writeText(url)
    toast.value = { type: 'success', message: 'Đã sao chép liên kết báo cáo vào clipboard!' }
  } catch (err) {
    // Người dùng hủy hộp thoại chia sẻ — không coi là lỗi
    if (err && err.name !== 'AbortError') {
      toast.value = { type: 'error', message: 'Không thể chia sẻ báo cáo.' }
    }
  }
}

const handleExportPDF = () => {
  // Dùng cơ chế in của trình duyệt (Save as PDF) — không cần thư viện ngoài.
  window.print()
}

const handleRetryReport = async () => {
  retryingReport.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const interviewId = route.params.id
    await reportService.retryReport(companyId, interviewId)
    toast.value = { type: 'success', message: 'Đang gửi yêu cầu tạo lại báo cáo...' }
    report.value = { ...(report.value || {}), status: 'generating' }
    if (reportPollTimer) clearInterval(reportPollTimer)
    reportPollTimer = setInterval(async () => {
      const latest = await loadReport()
      if (latest?.status === 'ready' || latest?.status === 'failed') {
        clearInterval(reportPollTimer)
        reportPollTimer = null
        retryingReport.value = false
      }
    }, 3000)
  } catch (err) {
    toast.value = { type: 'error', message: 'Không thể tạo lại báo cáo.' }
    retryingReport.value = false
  }
}
</script>

<template>
  <div class="max-w-5xl mx-auto pb-16 animate-fade-in space-y-6">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
    
    <!-- Page Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm mt-2">
      <div class="flex items-center gap-4">
        <button @click="router.push('/interviews')" class="p-2 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-500 dark:text-slate-400 rounded-lg transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-600 shrink-0">
          <ArrowLeft size="20" />
        </button>
        <div>
          <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Báo cáo Phỏng vấn: <span class="text-blue-600 dark:text-blue-400">{{ report?.candidate_name || 'Đang tải...' }}</span></h1>
          <p class="text-slate-500 dark:text-slate-400 mt-1 text-sm font-medium">{{ report?.job_title }} • {{ report ? new Date(report.date).toLocaleDateString('vi-VN') : '' }}</p>
        </div>
      </div>
      <div class="flex gap-3 no-print">
        <Button variant="secondary" @click="handleShare" class="bg-[var(--surface)] border-[var(--border)] text-[var(--text-primary)] hover:bg-[var(--surface-hover)]">
          <Share2 size="16" class="mr-1.5" /> Chia sẻ
        </Button>
        <Button @click="handleExportPDF" class="bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20">
          <Download size="16" class="mr-1.5" /> Xuất PDF
        </Button>
      </div>
    </Card>

    <!-- Loading State -->
    <div v-if="loading" class="flex flex-col items-center justify-center min-h-[400px] text-slate-500 dark:text-slate-400">
      <div class="w-12 h-12 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mb-4"></div>
      <span class="font-medium text-lg">Đang tổng hợp báo cáo AI...</span>
    </div>
    
    <!-- Error State -->
    <div v-else-if="!report" class="bg-rose-50 dark:bg-rose-500/10 border border-rose-200 dark:border-rose-500/20 rounded-2xl p-8 text-center text-rose-600 dark:text-rose-400 font-medium flex flex-col items-center justify-center gap-2">
      <AlertTriangle size="32" />
      <span>Không thể tải báo cáo hoặc AI gặp lỗi khi phân tích.</span>
      <Button @click="handleRetryReport" :disabled="retryingReport" class="mt-4 bg-rose-600 hover:bg-rose-700 text-white border-none shadow-md shadow-rose-500/20">
        <div v-if="retryingReport" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2"></div>
        {{ retryingReport ? 'Đang gửi...' : 'Thử tạo lại báo cáo bằng AI' }}
      </Button>
    </div>
    
    <div v-else class="space-y-6">
      <!-- Top Overview Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 md:gap-6">
        <!-- Score Card -->
        <Card class="rounded-2xl shadow-sm p-6 flex flex-col justify-center items-center">
          <p class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2">Điểm tổng quan</p>
          <div class="flex items-baseline gap-1">
            <span class="text-[28px] font-extrabold text-blue-600 dark:text-blue-400 leading-none">{{ report.overall_score }}</span>
            <span class="text-xl font-bold text-slate-400 dark:text-slate-500">/10</span>
          </div>
        </Card>
        
        <!-- Recommendation Card -->
        <Card class="rounded-2xl shadow-sm p-6 flex flex-col justify-center items-center">
          <p class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-3">Đề xuất từ AI</p>
          <span 
            class="px-4 py-1.5 text-sm font-bold uppercase tracking-wider rounded-full border"
            :class="[
              report.ai_recommendation === 'hire' ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20' : 
              'bg-amber-50 text-amber-600 border-amber-200 dark:bg-amber-500/10 dark:text-amber-400 dark:border-amber-500/20'
            ]"
          >
            {{ report.ai_recommendation === 'hire' ? 'Nên tuyển (Hire)' : report.ai_recommendation }}
          </span>
        </Card>
        
        <!-- Core Feedback Card -->
        <div class="md:col-span-2 bg-blue-50/50 dark:bg-blue-500/5 rounded-2xl shadow-sm border border-blue-100/50 dark:border-blue-500/10 p-6">
          <p class="text-xs font-bold text-blue-600 dark:text-blue-400 uppercase tracking-wider mb-3 flex items-center gap-1.5"><Sparkles size="14" /> Nhận xét cốt lõi</p>
          <p class="text-slate-700 dark:text-slate-300 font-medium leading-relaxed">{{ report.core_feedback }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <!-- Left Column: Scores & Decision -->
        <div class="space-y-6">
          <Card class="rounded-2xl shadow-sm p-6">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-5">Điểm chi tiết (Rubric)</h2>
            <div class="space-y-5">
              <div v-if="!report.rubric_scores || report.rubric_scores.length === 0" class="text-center py-8 text-slate-500 dark:text-slate-400">
                <p class="text-sm font-medium">AI chưa tạo rubric scores. Vui lòng chờ hoặc tạo lại báo cáo.</p>
              </div>
              <div v-for="item in report.rubric_scores" :key="item.name" class="space-y-2">
                <div class="flex justify-between items-center text-sm">
                  <span class="font-semibold text-slate-700 dark:text-slate-300">{{ item.name }}</span>
                  <span class="font-bold text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-500/10 px-2 py-0.5 rounded">{{ item.score }}/10</span>
                </div>
                <div class="h-2.5 bg-slate-100 dark:bg-slate-700/50 rounded-full overflow-hidden border border-slate-200/50 dark:border-slate-600/30">
                  <div class="h-full rounded-full transition-all duration-500" :style="{ width: `${item.score * 10}%`, backgroundColor: item.color || '#6366f1' }"></div>
                </div>
              </div>
            </div>
          </Card>
          
          <Card class="rounded-2xl shadow-sm p-6">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-5">Quyết định của Bạn</h2>
            <div class="space-y-4">
              <div class="space-y-2">
                <select 
                  v-model="decision" 
                  class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-slate-700 dark:text-slate-200"
                >
                  <option value="Chưa quyết định">-- Chọn quyết định --</option>
                  <option value="offer">Gửi Offer</option>
                  <option value="reject">Từ chối</option>
                  <option value="next_round">Phỏng vấn vòng sau</option>
                </select>
              </div>
              <div class="space-y-2">
                <textarea 
                  v-model="note" 
                  placeholder="Nhập ghi chú HR..." 
                  class="w-full px-4 py-3 bg-slate-50 dark:bg-slate-900 border border-slate-300 dark:border-slate-600 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-slate-700 dark:text-slate-200 placeholder:text-slate-400 resize-none h-28"
                ></textarea>
              </div>
              <Button :disabled="savingDecision" class="w-full justify-center bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20 py-2.5" @click="handleSaveDecision">
                <Save size="16" class="mr-1.5" /> Lưu quyết định
              </Button>
            </div>
          </Card>
        </div>

        <!-- Right Column: Insights & Evidence -->
        <div class="lg:col-span-2 space-y-6">
          
          <!-- Strengths & Weaknesses -->
          <Card class="rounded-2xl shadow-sm p-6 border-t-4 border-t-blue-500">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-6">Phân tích Điểm mạnh & Rủi ro</h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Strengths -->
              <div class="bg-emerald-50/50 dark:bg-emerald-500/5 rounded-xl p-5 border border-emerald-100/50 dark:border-emerald-500/10">
                <h4 class="flex items-center gap-2 text-emerald-600 dark:text-emerald-400 font-bold mb-4">
                  <CheckCircle size="18" /> Điểm mạnh
                </h4>
                <ul class="space-y-3">
                  <li v-for="(str, idx) in report.strengths" :key="idx" class="flex items-start gap-2.5 text-sm text-slate-700 dark:text-slate-300">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 mt-1.5 shrink-0"></span>
                    <span class="leading-relaxed">{{ str }}</span>
                  </li>
                </ul>
              </div>
              
              <!-- Weaknesses -->
              <div class="bg-rose-50/50 dark:bg-rose-500/5 rounded-xl p-5 border border-rose-100/50 dark:border-rose-500/10">
                <h4 class="flex items-center gap-2 text-rose-600 dark:text-rose-400 font-bold mb-4">
                  <AlertTriangle size="18" /> Điểm rủi ro (Cần lưu ý)
                </h4>
                <ul class="space-y-3">
                  <li v-for="(weak, idx) in report.weaknesses" :key="idx" class="flex items-start gap-2.5 text-sm text-slate-700 dark:text-slate-300">
                    <span class="w-1.5 h-1.5 rounded-full bg-rose-500 mt-1.5 shrink-0"></span>
                    <span class="leading-relaxed">{{ weak }}</span>
                  </li>
                </ul>
              </div>
            </div>
          </Card>

          <!-- Transcript Highlights -->
          <Card class="rounded-2xl shadow-sm p-6">
            <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-6 flex items-center gap-2">
              <FileText size="18" class="text-blue-500" /> Trích xuất Transcript (Bằng chứng)
            </h2>
            
            <div class="space-y-4">
              <div v-for="(hl, idx) in report.transcript_highlights" :key="idx" 
                   class="p-4 bg-slate-50 dark:bg-slate-900 rounded-xl border border-slate-100 dark:border-slate-700/50 relative overflow-hidden group">
                <div class="absolute left-0 top-0 bottom-0 w-1 bg-slate-300 dark:bg-slate-600" :style="`background-color: var(--${hl.type})`"></div>
                <div class="flex items-center gap-3 mb-2 ml-2">
                  <span class="px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider rounded border" :style="`color: var(--${hl.type}); border-color: var(--${hl.type})`">{{ hl.badge }}</span>
                  <span class="text-xs font-semibold text-slate-500 dark:text-slate-400">{{ hl.time }}</span>
                </div>
                <p class="text-sm font-medium text-slate-600 dark:text-slate-300 italic ml-2">
                  "{{ hl.text }}"
                </p>
              </div>
            </div>
            
            <Button @click="handleViewTranscripts" variant="ghost" class="w-full mt-6 text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-500/10 hover:bg-blue-100 dark:hover:bg-blue-500/20 border-none font-semibold transition-colors">
              <FileText size="16" class="mr-1.5" /> Xem toàn bộ Transcript
            </Button>
          </Card>

        </div>
      </div>
    </div>
    
    <!-- Transcripts Modal -->
    <div v-if="showTranscriptModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4 animate-fade-in">
      <div class="bg-[var(--surface)] rounded-2xl shadow-xl w-full max-w-2xl max-h-[80vh] flex flex-col overflow-hidden border border-[var(--border)]">
        <div class="flex items-center justify-between p-4 border-b border-slate-100 dark:border-slate-700 bg-slate-50/50 dark:bg-slate-800/50">
          <h3 class="font-bold text-slate-800 dark:text-slate-100 flex items-center gap-2">
            <FileText size="18" class="text-blue-500" /> Toàn bộ Transcript
          </h3>
          <button @click="showTranscriptModal = false" class="p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>
        <div class="p-4 flex-1 overflow-y-auto space-y-4 bg-slate-50 dark:bg-slate-900/20">
          <div v-if="loadingTranscripts" class="text-center text-slate-500 py-10 flex flex-col items-center gap-3">
             <div class="w-8 h-8 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin"></div>
             Đang tải lịch sử trò chuyện...
          </div>
          <div v-else-if="!fullTranscripts.length" class="text-center text-slate-500 py-10">
             Chưa có dữ liệu transcript nào.
          </div>
          <div v-else v-for="t in fullTranscripts" :key="t.id" class="p-4 rounded-xl border border-[var(--border)] bg-[var(--surface)] shadow-sm relative">
            <div class="flex items-center justify-between mb-2">
               <span class="font-bold text-sm" :class="t.speaker_type === 'recruiter' ? 'text-blue-600' : 'text-emerald-600'">
                 {{ t.speaker_name }}
               </span>
               <span class="text-xs font-semibold text-slate-400">
                 {{ new Date(t.created_at).toLocaleTimeString('vi-VN', {hour: '2-digit', minute:'2-digit', second:'2-digit'}) }}
               </span>
            </div>
            <p class="text-slate-700 dark:text-slate-300 text-sm leading-relaxed">{{ t.content }}</p>
          </div>
        </div>
        <div class="p-4 border-t border-[var(--border)] bg-[var(--surface)] flex justify-end">
          <Button @click="showTranscriptModal = false" variant="ghost" class="text-slate-600 dark:text-slate-300">Đóng</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  .no-print { display: none !important; }
}
</style>
