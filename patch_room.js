const fs = require('fs');
const path = require('path');

const file = path.join(__dirname, 'frontend/src/views/recruiter/InterviewRoomPage.vue');
let content = fs.readFileSync(file, 'utf8');

// 1. Add imports
content = content.replace(
  "import { useRoom } from '../../composables/useRoom'",
  `import { useRoom } from '../../composables/useRoom'\nimport { aiService } from '../../services/ai.service'\nimport { transcriptService } from '../../services/transcript.service'\nimport { interviewService } from '../../services/interview.service'`
);

// 2. Add state
content = content.replace(
  "const activeTab = ref('assistant')",
  `const activeTab = ref('assistant')\nconst aiSuggestions = ref([])\nconst aiScores = ref([])\nconst isSuggesting = ref(false)\nconst isScoring = ref(false)\nconst criteria = ref([])\nconst interviewDetails = ref(null)`
);

// 3. Update onMounted to fetch interview
content = content.replace(
  "const response = await roomService.getRoomToken(companyId, interviewId)",
  `try {
          const intv = await interviewService.getInterview(companyId, interviewId)
          interviewDetails.value = intv
          if (intv.job && intv.job.rubric) {
            criteria.value = intv.job.rubric.criteria || intv.job.rubric.rubric_criteria || []
          }
        } catch(e) { console.error('Lỗi lấy chi tiết phỏng vấn', e) }
        
        const response = await roomService.getRoomToken(companyId, interviewId)`
);

// 4. Add AI functions before handleEndCall
const aiFunctions = `
const getAiSuggestions = async () => {
  if (!interviewId || isSuggesting.value) return
  isSuggesting.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const res = await aiService.suggestFollowUp(companyId, interviewId, { focus: 'Chung' })
    if (res && res.suggested_question) {
      aiSuggestions.value.unshift(res)
    }
  } catch(err) {
    // toast.value = { type: 'error', message: 'Lỗi lấy gợi ý AI' }
    console.error('Lỗi lấy gợi ý', err)
  } finally {
    isSuggesting.value = false
  }
}

const scoreCurrentAnswer = async () => {
  if (!interviewId || isScoring.value || criteria.value.length === 0) {
    // toast.value = { type: 'warning', message: 'Không có tiêu chí chấm điểm' }
    return
  }
  isScoring.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const transcripts = await transcriptService.getTranscripts(companyId, interviewId)
    const tIds = transcripts.slice(-5).map(t => t.id).filter(id => id)
    if (tIds.length === 0) {
      // toast.value = { type: 'warning', message: 'Chưa có transcript hợp lệ để chấm điểm' }
      isScoring.value = false
      return
    }
    const cIds = criteria.value.map(c => c.id)
    const res = await aiService.scoreAnswer(companyId, interviewId, { transcript_ids: tIds, criterion_ids: cIds })
    if (res && res.scores) {
      aiScores.value = res.scores
      // toast.value = { type: 'success', message: 'Đã cập nhật điểm AI' }
    }
  } catch(err) {
    // toast.value = { type: 'error', message: 'Lỗi chấm điểm AI' }
    console.error('Lỗi chấm điểm', err)
  } finally {
    isScoring.value = false
  }
}

const handleEndCall`;

content = content.replace("const handleEndCall", aiFunctions);

// 5. Update confirmEndCall
content = content.replace(
  `const confirmEndCall = () => {
  showEndModal.value = false
  isEnding.value = true
  setTimeout(() => {
    router.push({ path: '/reports', state: { message: 'Đã lưu kết quả phỏng vấn thành công' } })
  }, 1500)
}`,
  `const confirmEndCall = async () => {
  showEndModal.value = false
  isEnding.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await interviewService.endInterview(companyId, interviewId, { generate_report: true })
    await aiService.generateReport(companyId, interviewId, false)
  } catch(e) {
    console.error('Lỗi kết thúc phỏng vấn', e)
  }
  setTimeout(() => {
    router.push({ path: '/reports', state: { message: 'Đã lưu kết quả phỏng vấn thành công' } })
  }, 1500)
}`
);

// 6. Update Header info
content = content.replace(
  `<h1 class="text-h2">Frontend Developer - Nguyễn Văn A</h1>`,
  `<h1 class="text-h2">{{ interviewDetails?.job?.title || 'Đang tải...' }} - {{ interviewDetails?.candidate?.name || 'Đang tải...' }}</h1>`
);

// 7. Update UI: Assistant tab
const mockAssistant = `<div>
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
          </div>`;

const realAssistant = `<div>
                  <h4 class="text-body" style="font-weight: 600; margin-bottom: 8px">Trợ lý AI</h4>
                  <p class="text-helper">AI đang nghe và sẵn sàng phân tích ngữ cảnh để đưa ra câu hỏi gợi ý.</p>
                </div>
              </div>
            </Card>

            <div>
              <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
                <h3 class="text-body" style="font-weight: 600">Gợi ý câu hỏi tiếp theo</h3>
                <Button variant="secondary" style="height: 32px; padding: 0 12px" @click="getAiSuggestions" :disabled="isSuggesting">
                  <span v-if="isSuggesting">Đang lấy...</span>
                  <span v-else>Lấy gợi ý</span>
                </Button>
              </div>
              <div style="display: flex; flex-direction: column; gap: 12px">
                <div v-if="aiSuggestions.length === 0" style="text-align: center; color: var(--text-muted); padding: 20px 0; font-size: 13px">
                  Bấm "Lấy gợi ý" để AI đề xuất câu hỏi dựa trên hội thoại gần nhất.
                </div>
                <div v-for="(sug, i) in aiSuggestions" :key="i" style="padding: 12px; border: 1px solid var(--border); border-radius: var(--radius); background-color: var(--surface)">
                  <p class="text-body" style="margin-bottom: 12px; font-weight: 500">"{{ sug.suggested_question }}"</p>
                  <div style="margin-bottom: 12px; font-size: 12px; color: var(--text-muted)">
                    Lý do: {{ sug.reason }}
                  </div>
                  <Button variant="secondary" style="width: 100%; font-size: 13px; height: 32px" @click="aiSuggestions.splice(i, 1)">Đã hỏi</Button>
                </div>
              </div>
            </div>
          </div>`;

content = content.replace(mockAssistant, realAssistant);

// 8. Update UI: Rubric tab
const mockRubric = `<div v-if="activeTab === 'rubric'" style="display: flex; flex-direction: column; gap: 16px">
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
          </div>`;

const realRubric = `<div v-if="activeTab === 'rubric'" style="display: flex; flex-direction: column; gap: 16px">
            <div style="display: flex; justify-content: flex-end; margin-bottom: 8px">
              <Button variant="secondary" style="height: 32px; padding: 0 12px" @click="scoreCurrentAnswer" :disabled="isScoring">
                <span v-if="isScoring">Đang chấm...</span>
                <span v-else>AI Chấm điểm hiện tại</span>
              </Button>
            </div>
            <div v-if="criteria.length === 0" style="text-align: center; padding: 20px 0; color: var(--text-muted); font-size: 13px">
              Không có tiêu chí đánh giá cho Job này.
            </div>
            <div v-for="c in criteria" :key="c.id" style="padding-bottom: 12px; border-bottom: 1px solid var(--border)">
              <div style="display: flex; justify-content: space-between; margin-bottom: 8px">
                <span class="text-body" style="font-weight: 500">{{ c.name }}</span>
                <span class="text-body" style="color: var(--primary); font-weight: 600">
                  {{ aiScores.find(s => s.criterion_id === c.id)?.score || '--' }}/{{ c.max_score || 5 }}
                </span>
              </div>
              <div style="height: 6px; background-color: var(--surface-soft); border-radius: 3px; overflow: hidden">
                <div 
                  :style="{
                    height: '100%', 
                    width: (aiScores.find(s => s.criterion_id === c.id)?.score / (c.max_score || 5) * 100) + '%',
                    backgroundColor: 'var(--primary)'
                  }">
                </div>
              </div>
              <p class="text-helper" style="margin-top: 8px; font-style: italic" v-if="aiScores.find(s => s.criterion_id === c.id)?.reasoning">
                "{{ aiScores.find(s => s.criterion_id === c.id)?.reasoning }}"
              </p>
              <p class="text-helper" style="margin-top: 4px; color: var(--text-muted)" v-else>
                Đang chờ dữ liệu...
              </p>
            </div>
          </div>`;

content = content.replace(mockRubric, realRubric);

fs.writeFileSync(file, content);
console.log('Done modifying InterviewRoomPage.vue');
