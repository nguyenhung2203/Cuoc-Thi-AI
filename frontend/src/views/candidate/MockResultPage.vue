<script setup>
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { ArrowLeft, MessageCircle, CheckCircle, AlertCircle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const id = route.params.id

// Mock data cho một kết quả chi tiết
const sessionData = {
  role: 'Frontend Developer',
  level: 'Middle',
  date: '20/10/2023',
  overallScore: 8.5,
  summaryFeedback: 'Bạn đã làm rất tốt trong việc giải thích các khái niệm cốt lõi của React. Tuy nhiên, phần System Design cần cụ thể hơn về cách scale ứng dụng.',
  questions: [
    {
      id: 1,
      question: 'Bạn hãy giải thích cơ chế Virtual DOM trong React và tại sao nó lại giúp tăng hiệu suất?',
      candidateAnswer: 'Virtual DOM là một bản copy của Real DOM. Khi state thay đổi, React tạo ra một Virtual DOM mới, so sánh với cái cũ (diffing), và chỉ cập nhật những node bị thay đổi lên Real DOM.',
      score: 9,
      feedback: 'Câu trả lời rất chính xác, ngắn gọn và đi thẳng vào trọng tâm. Bạn có thể bổ sung thêm về quá trình Reconciliation để đạt điểm tuyệt đối.',
      goodPoints: ['Hiểu rõ khái niệm bản copy', 'Nắm được quá trình diffing'],
      improvePoints: ['Thiếu key term Reconciliation']
    },
    {
      id: 2,
      question: 'Làm thế nào để tối ưu hóa hiệu suất (performance) của một ứng dụng React lớn?',
      candidateAnswer: 'Tôi thường dùng useMemo và useCallback để tránh re-render. Ngoài ra cũng dùng React.lazy để code splitting.',
      score: 7.5,
      feedback: 'Các ý chính đều đúng, tuy nhiên bạn cần giải thích rõ HƯỚNG áp dụng thực tế thay vì chỉ liệt kê hooks. Khi nào KHÔNG NÊN dùng useMemo cũng là một ý quan trọng.',
      goodPoints: ['Đề cập đúng các công cụ tối ưu (useMemo, React.lazy)'],
      improvePoints: ['Cần ví dụ thực tế', 'Thiếu cân nhắc trade-off khi lạm dụng useMemo']
    }
  ]
}
</script>

<template>
  <div>
    <div style="display: flex; align-items: center; gap: 16px; margin-bottom: 32px">
      <Button variant="ghost" @click="router.push('/mock-results')" style="padding: 8px">
        <ArrowLeft size="20" />
      </Button>
      <div>
        <h1 class="text-h1">Chi tiết kết quả luyện tập</h1>
        <p class="text-helper" style="margin-top: 4px">{{ sessionData.role }} - Cấp độ {{ sessionData.level }} - Ngày {{ sessionData.date }}</p>
      </div>
    </div>

    <div style="display: grid; grid-template-columns: 1fr 2fr; gap: 24px; margin-bottom: 24px">
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
