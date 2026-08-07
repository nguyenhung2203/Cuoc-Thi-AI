<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { ArrowLeft, MessageCircle, CheckCircle, AlertCircle, Sparkles, CheckCircle2 } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { apiService } from '../../services/api.service'

const route = useRoute()
const router = useRouter()
const id = route.query.mock_id || route.params.id

const loading = ref(true)
const sessionData = ref(null)
const loadError = ref(false)

onMounted(async () => {
  if (!id) {
    loading.value = false
    return
  }
  try {
    const data = await apiService.get(`/mock-interviews/${id}`)
    const messages = await mockService.getMessages(id)

    let formattedQuestions = []
    let lastQ = null
    if (Array.isArray(messages)) {
      for (let m of messages) {
        if (m.sender_type === 'ai') {
          lastQ = m.content
        } else if (m.sender_type === 'candidate') {
          const nextAi = messages.find(m2 => m2.sender_type === 'ai' && m2.created_at > m.created_at)
          let score = 0, feedback = ''
          let goodPoints = [], improvePoints = []
          if (nextAi && nextAi.score_json) {
             const sj = typeof nextAi.score_json === 'string' ? JSON.parse(nextAi.score_json) : nextAi.score_json
             score = sj.score
             feedback = nextAi.content
             if (Array.isArray(sj.good_points)) goodPoints = sj.good_points
             if (Array.isArray(sj.improve_points)) improvePoints = sj.improve_points
          }
          if (lastQ) {
            formattedQuestions.push({
              id: formattedQuestions.length + 1,
              question: lastQ,
              candidateAnswer: m.content,
              score: score,
              feedback: feedback,
              goodPoints: goodPoints,
              improvePoints: improvePoints
            })
          }
        }
      }
    }

    let finalScore = data.final_score
    if (finalScore == null && formattedQuestions.length > 0) {
       finalScore = formattedQuestions.reduce((sum, q) => sum + q.score, 0) / formattedQuestions.length
    }

    let feedbackSummary = data.feedback_json || ''
    if (typeof feedbackSummary === 'string') {
      try { feedbackSummary = JSON.parse(feedbackSummary) } catch (_) { /* legacy plain text */ }
    }
    const allFeedback = formattedQuestions.reduce((acc, q) => ({
      strengths: [...acc.strengths, ...q.goodPoints],
      improvements: [...acc.improvements, ...q.improvePoints],
    }), { strengths: [], improvements: [] })

    sessionData.value = {
      role: data.target_role || 'Phỏng vấn thử',
      level: data.target_level || '',
      date: data.created_at ? new Date(data.created_at).toLocaleDateString('vi-VN') : '',
      overallScore: finalScore != null ? Number(finalScore).toFixed(1) : null,
      summaryFeedback: typeof feedbackSummary === 'string' ? feedbackSummary : '',
      strengths: [...new Set(allFeedback.strengths)],
      improvements: [...new Set(allFeedback.improvements)],
      questions: formattedQuestions
    }
  } catch (err) {
    console.error('Lỗi tải kết quả phỏng vấn thử:', err)
    loadError.value = true
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="space-y-8 pb-12 max-w-6xl mx-auto">
    <div class="flex items-center gap-4 mb-2 animate-rise">
      <button @click="router.push('/mock-results')" class="back-btn"><ArrowLeft :size="20" /></button>
      <div>
        <h1 class="page-title">Báo cáo kết quả</h1>
        <p v-if="sessionData" class="result-tags">
          <span class="tag">{{ sessionData.role }}</span>
          <span class="tag">Cấp độ: {{ sessionData.level }}</span>
          <span class="tag">{{ sessionData.date }}</span>
        </p>
      </div>
    </div>

    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="mr-spinner mb-4"></div>
      <p class="text-helper">Đang tải kết quả đánh giá...</p>
    </div>

    <div v-else-if="loadError || !sessionData" class="mr-error">
      <AlertCircle :size="42" />
      <p>Không thể tải kết quả bài luyện tập.</p>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">

      <!-- Overall Score Sidebar -->
      <div class="lg:col-span-1 space-y-6 sticky top-6 animate-rise">
        <Card class="card-elevate score-card">
          <h2 class="score-title">Điểm tổng kết</h2>

          <div class="score-ring">
            <svg class="score-svg" viewBox="0 0 100 100">
              <circle class="ring-bg" stroke-width="8" cx="50" cy="50" r="40" fill="transparent"></circle>
              <circle class="ring-fg" stroke-width="8" stroke-linecap="round" cx="50" cy="50" r="40" fill="transparent" :stroke-dasharray="251.2" :stroke-dashoffset="251.2 - (251.2 * (parseFloat(sessionData.overallScore || 0) * 10)) / 100"></circle>
            </svg>
            <div class="score-center">
              <span class="score-val">{{ sessionData.overallScore != null ? sessionData.overallScore : '—' }}</span>
              <span class="score-max">/10</span>
            </div>
          </div>

          <div class="score-note">
            <h3 class="score-note-title"><Sparkles :size="18" /> Nhận xét chung</h3>
            <p class="score-note-text">{{ sessionData.summaryFeedback || 'Chưa có nhận xét tổng hợp cho bài luyện tập này.' }}</p>
          </div>
          <div v-if="sessionData.strengths?.length || sessionData.improvements?.length" class="mt-5 space-y-4 text-left">
            <div v-if="sessionData.strengths?.length">
              <h3 class="font-bold text-emerald-600">Điểm mạnh</h3>
              <ul class="mt-2 list-disc pl-5 text-sm text-[var(--text-secondary)]"><li v-for="(item, i) in sessionData.strengths" :key="'s' + i">{{ item }}</li></ul>
            </div>
            <div v-if="sessionData.improvements?.length">
              <h3 class="font-bold text-amber-600">Gợi ý cải thiện</h3>
              <ul class="mt-2 list-disc pl-5 text-sm text-[var(--text-secondary)]"><li v-for="(item, i) in sessionData.improvements" :key="'i' + i">{{ item }}</li></ul>
            </div>
          </div>
        </Card>
      </div>

      <!-- Detailed Questions List -->
      <div class="lg:col-span-2 space-y-6">
        <Card class="card-elevate qa-card animate-rise">
          <h2 class="qa-head">
            <MessageCircle :size="22" />
            Đánh giá chi tiết từng câu hỏi
          </h2>

          <div class="space-y-10">
            <div v-if="sessionData.questions.length === 0" class="mr-empty-q">
              <MessageCircle :size="34" />
              <p>Bài luyện tập này chưa có câu hỏi và câu trả lời nào được ghi nhận.</p>
            </div>
            <div v-for="(q, index) in sessionData.questions" :key="q.id">
              <div class="flex items-start gap-4 mb-4">
                <div class="q-badge">Q{{ index + 1 }}</div>
                <div class="flex-1">
                  <h4 class="q-text">{{ q.question }}</h4>
                  <div class="q-answer">
                    <div class="q-answer-label">Câu trả lời của bạn</div>
                    <p class="q-answer-text">"{{ q.candidateAnswer }}"</p>
                  </div>
                </div>
              </div>

              <div class="qa-body">
                <div class="flex justify-between items-center mb-4">
                  <h5 class="fb-title"><CheckCircle2 :size="18" /> AI Feedback</h5>
                  <div class="flex items-center gap-3">
                    <div class="w-24 h-2 bg-slate-100 rounded-full overflow-hidden">
                      <div class="h-full rounded-full transition-all duration-1000" :class="q.score >= 8 ? 'bg-emerald-500' : q.score >= 6 ? 'bg-amber-500' : 'bg-rose-500'" :style="`width: ${q.score * 10}%`"></div>
                    </div>
                    <span class="font-bold text-sm" :class="q.score >= 8 ? 'text-emerald-600' : q.score >= 6 ? 'text-amber-600' : 'text-rose-600'">
                      {{ q.score }}/10
                    </span>
                  </div>
                </div>

                <p class="fb-text">{{ q.feedback }}</p>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div class="pt-box is-good">
                    <p class="pt-head"><CheckCircle :size="18" /> Điểm tốt</p>
                    <ul>
                      <li v-for="(p, i) in q.goodPoints" :key="i"><span class="pt-dot is-good"></span>{{ p }}</li>
                    </ul>
                  </div>
                  <div class="pt-box is-improve">
                    <p class="pt-head"><AlertCircle :size="18" /> Cần cải thiện</p>
                    <ul>
                      <li v-for="(p, i) in q.improvePoints" :key="i"><span class="pt-dot is-improve"></span>{{ p }}</li>
                    </ul>
                  </div>
                </div>
              </div>

              <hr v-if="index < sessionData.questions.length - 1" class="qa-divider" />
            </div>
          </div>
        </Card>
      </div>

    </div>
  </div>
</template>

<style scoped>
.back-btn { display: inline-flex; align-items: center; justify-content: center; width: 42px; height: 42px; border-radius: var(--radius); background: var(--surface); border: 1px solid var(--border); color: var(--text-secondary); cursor: pointer; box-shadow: var(--shadow-sm); transition: all 0.2s ease; flex-shrink: 0; }
.back-btn:hover { color: var(--primary); border-color: var(--primary-light); }
.result-tags { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 8px; }
.tag { background: var(--surface-soft); color: var(--text-secondary); padding: 3px 12px; border-radius: var(--radius); font-size: 13px; font-weight: 500; }

.mr-spinner { width: 44px; height: 44px; border-radius: 50%; border: 3px solid var(--primary-light); border-top-color: var(--primary); animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.mr-error { display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 40px; text-align: center; background: rgba(220,38,38,0.06); border: 1px solid rgba(220,38,38,0.2); border-radius: var(--radius-lg); color: var(--danger); font-weight: 600; }
.mr-empty-q { display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 40px 16px; text-align: center; color: var(--text-secondary); }
.mr-empty-q :deep(svg) { color: var(--text-muted); }

.score-card { padding: 28px; text-align: center; }
.score-title { font-size: 16px; font-weight: 700; color: var(--text-main); margin-bottom: 22px; }
.score-ring { position: relative; width: 160px; height: 160px; margin: 0 auto 8px; }
.score-svg { width: 100%; height: 100%; transform: rotate(-90deg); }
.ring-bg { stroke: var(--surface-soft); }
.ring-fg { stroke: var(--success); transition: stroke-dashoffset 1s cubic-bezier(0.2,0.8,0.2,1); }
.score-center { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.score-val { font-size: 40px; font-weight: 800; color: var(--success); line-height: 1; }
.score-max { font-size: 14px; font-weight: 600; color: var(--text-muted); }
.score-note { margin-top: 22px; padding-top: 22px; border-top: 1px solid var(--border); text-align: left; }
.score-note-title { display: flex; align-items: center; gap: 8px; font-weight: 700; color: var(--text-main); margin-bottom: 10px; }
.score-note-title :deep(svg) { color: var(--accent); }
.score-note-text { color: var(--text-secondary); font-size: 14px; line-height: 1.6; }

.qa-card { padding: 28px; }
.qa-head { display: flex; align-items: center; gap: 10px; font-size: 18px; font-weight: 700; color: var(--text-main); padding-bottom: 18px; margin-bottom: 24px; border-bottom: 1px solid var(--border); }
.qa-head :deep(svg) { color: var(--primary); }
.q-badge { display: flex; align-items: center; justify-content: center; width: 40px; height: 40px; border-radius: var(--radius); background: var(--primary-light); color: var(--primary); font-weight: 700; flex-shrink: 0; margin-top: 2px; }
.q-text { font-size: 16px; font-weight: 700; color: var(--text-main); line-height: 1.4; margin-bottom: 12px; }
.q-answer { position: relative; background: var(--surface-soft); border: 1px solid var(--border); border-radius: var(--radius); padding: 18px; }
.q-answer-label { position: absolute; top: -10px; left: 14px; background: var(--surface); padding: 0 8px; font-size: 11px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.04em; }
.q-answer-text { color: var(--text-secondary); font-style: italic; }
.qa-body { margin-left: 56px; }
.fb-title { display: flex; align-items: center; gap: 8px; font-weight: 700; color: var(--text-main); }
.fb-title :deep(svg) { color: var(--accent); }
.fb-text { color: var(--text-secondary); line-height: 1.6; margin: 16px 0 20px; }
.pt-box { padding: 18px; border-radius: var(--radius); }
.pt-box.is-good { background: rgba(22,163,74,0.08); border: 1px solid rgba(22,163,74,0.2); }
.pt-box.is-improve { background: rgba(217,119,6,0.08); border: 1px solid rgba(217,119,6,0.2); }
.pt-head { display: flex; align-items: center; gap: 8px; font-weight: 700; margin-bottom: 12px; }
.pt-box.is-good .pt-head { color: var(--success); }
.pt-box.is-improve .pt-head { color: var(--warning); }
.pt-box ul { display: flex; flex-direction: column; gap: 8px; }
.pt-box li { display: flex; align-items: flex-start; gap: 8px; font-size: 14px; color: var(--text-secondary); }
.pt-dot { width: 6px; height: 6px; border-radius: 50%; margin-top: 7px; flex-shrink: 0; }
.pt-dot.is-good { background: var(--success); }
.pt-dot.is-improve { background: var(--warning); }
.qa-divider { margin: 28px 0; border: none; border-top: 1px solid var(--border); }
</style>
