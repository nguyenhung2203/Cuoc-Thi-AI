<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { ArrowLeft, MessageCircle, CheckCircle, AlertCircle } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'
import { apiService } from '../../services/api.service'

const route = useRoute()
const router = useRouter()
const id = route.query.mock_id || route.params.id

const loading = ref(true)
const sessionData = ref(null)

onMounted(async () => {
  if (!id) {
    loading.value = false
    return
  }
  try {
    const data = await apiService.get(`/mock-interviews/${id}`)
    const messages = await mockService.getMessages(id)
    
    // Parse questions from messages
    let formattedQuestions = []
    let lastQ = null
    for (let m of messages) {
      if (m.sender_type === 'ai') {
        if (m.question_type) lastQ = m.content
      } else if (m.sender_type === 'candidate') {
        // find the ai score next
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
            feedback: feedback
          })
        }
      }
    }

    // fallback score logic if final_score null
    let finalScore = data.final_score
    if (finalScore == null && formattedQuestions.length > 0) {
       finalScore = formattedQuestions.reduce((sum, q) => sum + q.score, 0) / formattedQuestions.length
    }
    
    sessionData.value = {
      role: data.target_role,
      level: data.target_level || 'Junior',
      date: new Date(data.created_at).toLocaleDateString('vi-VN'),
      overallScore: finalScore ? finalScore.toFixed(1) : '0',
      summaryFeedback: data.feedback_json || 'Không có nhận xét chung',
      questions: formattedQuestions.map(q => ({
        ...q,
        goodPoints: ['Đã trả lời câu hỏi'],
        improvePoints: ['Có thể cung cấp thêm ví dụ']
      }))
    }
  } catch (err) {
    console.error('Lỗi tải kết quả mock', err)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto">
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div class="flex items-center gap-4">
        <button @click="router.push('/mock-results')" class="p-2.5 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 hover:border-gray-300 transition-colors shadow-sm text-gray-600">
          <ArrowLeft class="w-5 h-5" />
        </button>
        <div>
          <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-emerald-600 to-teal-600">Báo cáo kết quả</h1>
          <p v-if="sessionData" class="text-gray-500 mt-2 font-medium flex items-center gap-2">
            <span class="bg-gray-100 text-gray-700 px-2.5 py-0.5 rounded-lg text-sm">{{ sessionData.role }}</span>
            <span class="bg-gray-100 text-gray-700 px-2.5 py-0.5 rounded-lg text-sm">Cấp độ: {{ sessionData.level }}</span>
            <span class="bg-gray-100 text-gray-700 px-2.5 py-0.5 rounded-lg text-sm">{{ sessionData.date }}</span>
          </p>
        </div>
      </div>
    </div>

    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="w-12 h-12 border-4 border-emerald-200 border-t-emerald-600 rounded-full animate-spin mb-4"></div>
      <p class="text-gray-500 font-medium animate-pulse">Đang tải kết quả đánh giá...</p>
    </div>

    <div v-else-if="!sessionData" class="bg-rose-50 border border-rose-200 rounded-2xl p-8 text-center text-rose-600 flex flex-col items-center justify-center shadow-sm">
      <AlertCircle class="w-12 h-12 mb-3 opacity-80" />
      <p class="font-bold text-lg">Không thể tải kết quả bài luyện tập.</p>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8 items-start">
      
      <!-- Overall Score Sidebar -->
      <div class="lg:col-span-1 space-y-6 sticky top-6">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-8 relative overflow-hidden text-center">
          <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-emerald-100 to-transparent rounded-bl-full opacity-50"></div>
          
          <h2 class="text-lg font-bold text-gray-800 mb-6 relative z-10">Điểm tổng kết</h2>
          
          <!-- Circular Progress (CSS based) -->
          <div class="relative w-40 h-40 mx-auto mb-6 z-10 drop-shadow-md">
            <svg class="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
              <circle class="text-gray-100 stroke-current" stroke-width="8" cx="50" cy="50" r="40" fill="transparent"></circle>
              <circle class="text-emerald-500 stroke-current drop-shadow" stroke-width="8" stroke-linecap="round" cx="50" cy="50" r="40" fill="transparent" :stroke-dasharray="251.2" :stroke-dashoffset="251.2 - (251.2 * (parseFloat(sessionData.overallScore) * 10)) / 100"></circle>
            </svg>
            <div class="absolute inset-0 flex flex-col items-center justify-center">
              <span class="text-4xl font-black text-emerald-600">{{ sessionData.overallScore }}</span>
              <span class="text-sm font-bold text-gray-400">/10</span>
            </div>
          </div>
          
          <div class="pt-6 border-t border-gray-100 relative z-10 text-left">
            <h3 class="font-bold text-gray-800 mb-3 flex items-center gap-2">
              <Sparkles class="w-5 h-5 text-amber-500" /> Nhận xét chung
            </h3>
            <p class="text-gray-600 text-sm leading-relaxed">{{ sessionData.summaryFeedback }}</p>
          </div>
        </Card>
      </div>

      <!-- Detailed Questions List -->
      <div class="lg:col-span-2 space-y-6">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl p-8">
          <h2 class="text-xl font-bold text-gray-800 mb-6 flex items-center gap-2 border-b border-gray-100 pb-4">
            <MessageCircle class="w-6 h-6 text-emerald-600" />
            Đánh giá chi tiết từng câu hỏi
          </h2>
          
          <div class="space-y-12">
            <div v-for="(q, index) in sessionData.questions" :key="q.id" class="group">
              <div class="flex items-start gap-4 mb-4">
                <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 font-bold flex items-center justify-center shrink-0 border border-emerald-100 mt-1 shadow-sm">
                  Q{{ index + 1 }}
                </div>
                <div class="flex-1">
                  <h4 class="text-lg font-bold text-gray-800 leading-snug mb-3">{{ q.question }}</h4>
                  
                  <div class="bg-slate-50 border border-slate-200 rounded-2xl p-5 mb-5 relative">
                    <div class="absolute -top-3 left-4 bg-white px-2 text-xs font-bold text-slate-500 uppercase tracking-wider">Câu trả lời của bạn</div>
                    <p class="text-gray-700 italic">"{{ q.candidateAnswer }}"</p>
                  </div>
                </div>
              </div>

              <div class="ml-14">
                <div class="flex justify-between items-center mb-4">
                  <h5 class="font-bold text-gray-800 flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                    AI Feedback
                  </h5>
                  <div :class="q.score >= 8 ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'" class="px-3 py-1 rounded-lg font-bold text-sm shadow-sm border" :style="q.score >= 8 ? 'border-color: #a7f3d0;' : 'border-color: #fde68a;'">
                    Điểm: {{ q.score }}/10
                  </div>
                </div>
                
                <p class="text-gray-600 leading-relaxed mb-6">{{ q.feedback }}</p>
                
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <!-- Good Points Card -->
                  <div class="bg-emerald-50 border border-emerald-100 rounded-xl p-5 hover:shadow-md transition-shadow">
                    <p class="text-emerald-800 font-bold flex items-center gap-2 mb-3">
                      <CheckCircle class="w-5 h-5 text-emerald-600" /> Điểm tốt
                    </p>
                    <ul class="space-y-2">
                      <li v-for="(p, i) in q.goodPoints" :key="i" class="flex items-start gap-2 text-emerald-700 text-sm">
                        <span class="w-1.5 h-1.5 bg-emerald-400 rounded-full mt-1.5 shrink-0"></span>
                        {{ p }}
                      </li>
                    </ul>
                  </div>
                  
                  <!-- Improve Points Card -->
                  <div class="bg-amber-50 border border-amber-100 rounded-xl p-5 hover:shadow-md transition-shadow">
                    <p class="text-amber-800 font-bold flex items-center gap-2 mb-3">
                      <AlertCircle class="w-5 h-5 text-amber-600" /> Cần cải thiện
                    </p>
                    <ul class="space-y-2">
                      <li v-for="(p, i) in q.improvePoints" :key="i" class="flex items-start gap-2 text-amber-700 text-sm">
                        <span class="w-1.5 h-1.5 bg-amber-400 rounded-full mt-1.5 shrink-0"></span>
                        {{ p }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
              
              <hr v-if="index < sessionData.questions.length - 1" class="my-8 border-gray-100" />
            </div>
          </div>
        </Card>
      </div>
      
    </div>
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}
</style>
