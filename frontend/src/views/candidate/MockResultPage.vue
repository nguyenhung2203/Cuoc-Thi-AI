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

const demoMockReports = {
  'mock-01': {
    role: 'Senior AI Engineer (LLM / RAG Pipeline)',
    level: 'Senior Tier 2',
    date: '11/07/2026',
    overallScore: '9.4',
    summaryFeedback: 'Ứng viên thể hiện hiểu biết sâu sắc về kiến trúc LLM hiện đại, làm chủ hoàn toàn luồng RAG Pipeline và có tư duy tối ưu chi phí hạ tầng xuất sắc. Sẵn sàng cho vị trí Senior AI Engineer tại các hệ thống Enterprise.',
    questions: [
      {
        id: 1,
        question: 'Giải thích sự khác biệt cốt lõi giữa Fine-tuning LLM và RAG Pipeline? Trong kịch bản doanh nghiệp nào bạn sẽ ưu tiên chọn RAG thay vì Fine-tuning?',
        candidateAnswer: 'RAG giúp bổ sung kiến thức động từ cơ sở dữ liệu bên ngoài mà không cần huấn luyện lại mô hình, rất tốt cho dữ liệu doanh nghiệp thay đổi liên tục hàng ngày. Fine-tuning giúp thay đổi văn phong, định dạng đầu ra JSON hoặc kiến thức miền chuyên sâu cố định, nhưng tốn chi phí và không cập nhật real-time được.',
        score: 9.5,
        feedback: 'Trả lời chính xác, phân định rõ ràng trade-off giữa chi phí huấn luyện và tính thời gian thực của dữ liệu.',
        goodPoints: [
          'Nắm vững nguyên lý hoạt động của Vector DB và Embedding retrieval pipeline.',
          'Đề cập chính xác đến bài toán chi phí (Cost optimization) và độ trễ (Latency).'
        ],
        improvePoints: [
          'Có thể bổ sung thêm kỹ thuật Hybrid Search (BM25 + Semantic Search) để tối ưu độ chính xác truy xuất từ khóa đặc thù.'
        ]
      },
      {
        id: 2,
        question: 'Làm thế nào để giảm thiểu tối đa hiện tượng Hallucination (ảo giác) trong các hệ thống RAG quy mô lớn?',
        candidateAnswer: 'Sử dụng kỹ thuật Re-ranking bằng Cohere/BGE reranker sau khi retrieve từ Vector DB, áp dụng Prompt Engineering yêu cầu mô hình chỉ trả lời dựa trên context được cung cấp (Grounding), và thiết lập Self-Check / Guardrails để kiểm duyệt đầu ra.',
        score: 9.3,
        feedback: 'Đưa ra giải pháp đa tầng (Multi-stage RAG) đạt chuẩn hệ thống cấp doanh nghiệp lớn.',
        goodPoints: [
          'Khả năng áp dụng Re-ranking và Guardrails để kiểm soát chất lượng đầu ra chặt chẽ.',
          'Tư duy phân lớp bảo mật và độ tin cậy của AI response.'
        ],
        improvePoints: [
          'Nên định lượng thêm bộ chỉ số kiểm thử tự động RAGAs (Faithfulness, Answer Relevance) trong quy trình CI/CD.'
        ]
      }
    ]
  },
  'mock-02': {
    role: 'Backend Tech Lead (High-concurrency)',
    level: 'Tech Lead',
    date: '08/07/2026',
    overallScore: '8.7',
    summaryFeedback: 'Nền tảng hệ thống phân tán (Distributed Systems) và tư duy thiết kế Microservices vững vàng. Có kinh nghiệm thực chiến với tải cao, cần làm rõ thêm các chi tiết khóa phân tán trên Redis Cluster.',
    questions: [
      {
        id: 1,
        question: 'Thiết kế kiến trúc hệ thống xử lý 10,000 requests/giây (RPS) cho tính năng đặt vé flash-sale mà không bị sập cơ sở dữ liệu?',
        candidateAnswer: 'Sử dụng Redis Cluster để giữ chỗ (reservation) trong bộ nhớ tạm thời, áp dụng rate limiting bằng Token Bucket tại API Gateway, và đưa giao dịch ghi DB vào Kafka message queue để xử lý bất đồng bộ (Asynchronous processing).',
        score: 8.8,
        feedback: 'Cấu trúc giải pháp rất tốt, phân tách luồng Synchronous và Asynchronous hiệu quả nhằm bảo vệ Database.',
        goodPoints: [
          'Sử dụng Redis giữ chỗ trong RAM giúp giảm 95% tải trực tiếp vào Relational DB.',
          'Áp dụng Message Queue (Kafka) để điều tiết lưu lượng (Traffic shaping / Peak clipping).'
        ],
        improvePoints: [
          'Cần giải thích sâu hơn cách xử lý Race-condition khi 2 user cùng tranh 1 vé cuối cùng (sử dụng Lua script trên Redis hoặc Distributed Lock).'
        ]
      }
    ]
  },
  'mock-03': {
    role: 'Fullstack Systems Architect (Vue / Go)',
    level: 'Principal',
    date: '03/07/2026',
    overallScore: '8.2',
    summaryFeedback: 'Tư duy tổng thể vững vàng từ Frontend đến Backend. Trình bày cấu trúc STAR rõ ràng, logic phản biện sắc bén và hiểu rõ vòng đời phát triển phần mềm Agile.',
    questions: [
      {
        id: 1,
        question: 'Làm thế nào để tối ưu Core Web Vitals (LCP, INP, CLS) cho một ứng dụng Vue 3 quy mô lớn có hàng ngàn component?',
        candidateAnswer: 'Áp dụng Lazy loading cho các route và heavy component, sử dụng Vite dynamic code-splitting, tối ưu hình ảnh sang định dạng WebP/AVIF và tránh thay đổi layout đột ngột (CLS) bằng cách giữ chỗ (skeleton/placeholder) trước khi tải dữ liệu.',
        score: 8.5,
        feedback: 'Nắm vững các kỹ thuật tối ưu Frontend hiệu năng cao và cấu hình Vite/Vue Router hiện đại.',
        goodPoints: [
          'Đưa ra các giải pháp thực tế có thể tích hợp ngay vào quy trình build Vite.',
          'Hiểu rõ cơ chế Virtual DOM và reactive performance của Vue 3 Composition API.'
        ],
        improvePoints: [
          'Có thể bổ sung thêm kiến trúc Server-Side Rendering (SSR) hoặc Nuxt 3 để giảm tối đa thời gian Time-to-First-Byte (TTFB).'
        ]
      }
    ]
  }
}

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
          if (m.question_type) lastQ = m.content
        } else if (m.sender_type === 'candidate') {
          const nextAi = messages.find(m2 => m2.sender_type === 'ai' && m2.created_at > m.created_at)
          let score = 0, feedback = ''
          if (nextAi && nextAi.score_json) {
             const sj = JSON.parse(nextAi.score_json)
             score = sj.score
             feedback = nextAi.content
          }
          if (lastQ) {
            formattedQuestions.push({
              id: formattedQuestions.length + 1,
              question: lastQ,
              candidateAnswer: m.content,
              score: score,
              feedback: feedback,
              goodPoints: ['Trình bày logic mạch lạc', 'Bám sát trọng tâm câu hỏi'],
              improvePoints: ['Có thể định lượng thêm các chỉ số hiệu năng cụ thể']
            })
          }
        }
      }
    }

    let finalScore = data.final_score
    if (finalScore == null && formattedQuestions.length > 0) {
       finalScore = formattedQuestions.reduce((sum, q) => sum + q.score, 0) / formattedQuestions.length
    }
    
    sessionData.value = {
      role: data.target_role || 'Candidate Interview',
      level: data.target_level || 'Senior',
      date: data.created_at ? new Date(data.created_at).toLocaleDateString('vi-VN') : '11/07/2026',
      overallScore: finalScore ? Number(finalScore).toFixed(1) : '8.5',
      summaryFeedback: data.feedback_json || 'Nhận xét tổng hợp tự động từ Trí tuệ Nhân tạo.',
      questions: formattedQuestions.length > 0 ? formattedQuestions : (demoMockReports['mock-01'].questions)
    }
  } catch (err) {
    console.warn('Sử dụng dữ liệu mô phỏng AI Báo cáo chi tiết cho ID:', id)
    sessionData.value = demoMockReports[id] || demoMockReports['mock-01']
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

    <div v-else-if="!sessionData" class="mr-error">
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
              <circle class="ring-fg" stroke-width="8" stroke-linecap="round" cx="50" cy="50" r="40" fill="transparent" :stroke-dasharray="251.2" :stroke-dashoffset="251.2 - (251.2 * (parseFloat(sessionData.overallScore) * 10)) / 100"></circle>
            </svg>
            <div class="score-center">
              <span class="score-val">{{ sessionData.overallScore }}</span>
              <span class="score-max">/10</span>
            </div>
          </div>

          <div class="score-note">
            <h3 class="score-note-title"><Sparkles :size="18" /> Nhận xét chung</h3>
            <p class="score-note-text">{{ sessionData.summaryFeedback }}</p>
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
