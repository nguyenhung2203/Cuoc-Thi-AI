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
  <div>
    <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
      <Button variant="ghost" @click="router.push('/mock-results')" style="padding: 8px">
        <ArrowLeft size="20" />
      </Button>
      <div>
        <h1 class="text-h1">Chi tiết kết quả luyện tập</h1>
        <p v-if="sessionData" class="text-helper" style="margin-top: 4px">{{ sessionData.role }} - Cấp độ {{ sessionData.level }} - Ngày {{ sessionData.date }}</p>
      </div>
    </div>

    <div v-if="loading" style="text-align: center; padding: 40px; color: var(--text-muted)">
      Đang tải kết quả...
    </div>
    <div v-else-if="!sessionData" style="text-align: center; padding: 40px; color: var(--danger)">
      Không thể tải kết quả bài luyện tập.
    </div>
    <div v-else style="display: grid; grid-template-columns: 1fr 2fr; gap: 24px; margin-bottom: 24px">
      <Card>
        <div style="text-align: center; margin-bottom: 24px">
          <div style="display: inline-flex; align-items: center; justify-content: center; width: 80px; height: 80px; border-radius: 50%; background-color: rgba(22, 163, 74, 0.1); color: var(--success); font-size: 28px; font-weight: bold; margin-bottom: 16px">
            {{ sessionData.overallScore }}
          </div>
          <h2 class="text-body" style="font-weight: 600">Điểm tổng kết</h2>
        </div>
        
        <div style="padding-top: 16px; border-top: 1px solid var(--border)">
          <h3 class="text-helper" style="font-weight: 600; margin-bottom: 8px">Nhận xét chung:</h3>
          <p class="text-body" style="color: var(--text-secondary)">{{ sessionData.summaryFeedback }}</p>
        </div>
      </Card>

      <Card title="Đánh giá chi tiết từng câu hỏi">
        <div style="display: flex; flex-direction: column; gap: 24px">
          <div v-for="(q, index) in sessionData.questions" :key="q.id" :style="{ paddingBottom: '24px', borderBottom: index < sessionData.questions.length - 1 ? '1px solid var(--border)' : 'none' }">
            <div style="display: flex; gap: 12px; margin-bottom: 12px">
              <MessageCircle size="20" color="var(--primary)" style="flex-shrink: 0; margin-top: 2px" />
              <div>
                <h4 class="text-body" style="font-weight: 600">Câu hỏi {{ index + 1 }}: {{ q.question }}</h4>
              </div>
            </div>

            <div style="margin-left: 32px; margin-bottom: 16px; padding: 12px; background-color: var(--surface-soft); border-radius: 8px">
              <p class="text-helper" style="font-weight: 600; margin-bottom: 4px; color: var(--text-main)">Câu trả lời của bạn:</p>
              <p class="text-body" style="color: var(--text-secondary)">"{{ q.candidateAnswer }}"</p>
            </div>

            <div style="margin-left: 32px">
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px">
                <p class="text-helper" style="font-weight: 600; color: var(--primary)">AI Feedback</p>
                <Badge :type="q.score >= 8 ? 'success' : 'warning'">Điểm: {{ q.score }}/10</Badge>
              </div>
              <p class="text-body" style="margin-bottom: 12px">{{ q.feedback }}</p>
              
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px">
                <div style="padding: 12px; border: 1px solid #bbf7d0; background-color: #f0fdf4; border-radius: 8px">
                  <p class="text-helper" style="color: #166534; font-weight: 600; display: flex; align-items: center; gap: 4px; margin-bottom: 8px"><CheckCircle size="14" /> Điểm tốt</p>
                  <ul class="text-body" style="padding-left: 16px; color: #166534; margin: 0">
                    <li v-for="(p, i) in q.goodPoints" :key="i">{{ p }}</li>
                  </ul>
                </div>
                <div style="padding: 12px; border: 1px solid #fef08a; background-color: #fefce8; border-radius: 8px">
                  <p class="text-helper" style="color: #854d0e; font-weight: 600; display: flex; align-items: center; gap: 4px; margin-bottom: 8px"><AlertCircle size="14" /> Cần cải thiện</p>
                  <ul class="text-body" style="padding-left: 16px; color: #854d0e; margin: 0">
                    <li v-for="(p, i) in q.improvePoints" :key="i">{{ p }}</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
