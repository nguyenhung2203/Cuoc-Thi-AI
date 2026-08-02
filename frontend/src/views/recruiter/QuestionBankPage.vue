<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import Input from '../../components/common/AppInput.vue'
import { Plus, Search, Filter, Layers, Edit, Trash2, Sparkles, RefreshCw } from 'lucide-vue-next'
import { questionBankService } from '../../services/questionBank.service'
import { jobService } from '../../services/job.service'
import { authStore } from '../../stores/auth.store'
import { hasDuplicateNormalized, isOneOf, maxLength, minLength, normalizeText, requiredTrim, validateForm } from '../../utils/validators.js'

const questions = ref([])
const loading = ref(true)
const aiGenerating = ref(false)
const searchKeyword = ref('')
const filterLevel = ref('')
const filterType = ref('')
const jobs = ref([])
const selectedJobId = ref('')

const columns = [
  { header: 'Câu hỏi', key: 'text' },
  { header: 'Vị trí (Role)', key: 'role' },
  { header: 'Cấp độ', key: 'level' },
  { header: 'Loại', key: 'type' },
  { header: 'Nguồn', key: 'source' },
  { header: 'Hành động', key: 'action' }
]

const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const showFilterModal = ref(false)
const deletingId = ref(null)
const editingId = ref(null)
const saving = ref(false)
const toast = ref(null)
const formErrors = ref({})
const generateCount = ref(10)
const newQuestion = ref({
  text: '',
  role: 'All',
  level: 'Fresher',
  type: 'Technical',
  expected_signals: ''
})

// Map API response → display format
const mapQuestion = (q) => ({
  id: q.id,
  text: q.question_text || q.text || '(Không có nội dung)',
  role: Array.isArray(q.skill_tags) ? q.skill_tags.join(', ') : (q.role || '—'),
  level: q.level || '—',
  type: q.question_type || q.type || '—',
  source: q.is_ai_generated ? 'AI' : 'Thủ công',
  _raw: q
})

const loadQuestions = async () => {
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) { loading.value = false; return }
    const params = {}
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (filterLevel.value) params.level = filterLevel.value
    if (filterType.value) params.question_type = filterType.value
    const data = await questionBankService.getQuestions(companyId, params)
    questions.value = (Array.isArray(data) ? data : []).map(mapQuestion)
  } catch (err) {
    console.error('Lỗi tải kho câu hỏi', err)
    questions.value = [] // Empty state — Backend chưa sẵn sàng
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const companyId = authStore.user?.companies?.[0]?.id
  await loadQuestions()
  if (!companyId) return
  try {
    const data = await jobService.getJobs(companyId, { status: 'active', page_size: 100 })
    jobs.value = Array.isArray(data) ? data : (data?.items || data?.data || [])
    selectedJobId.value = jobs.value[0]?.id || ''
  } catch (error) {
    console.error('Lỗi tải danh sách công việc', error)
  }
})

// Open the shared modal in edit mode, prefilled from a row.
const openEdit = (row) => {
  editingId.value = row.id
  const q = row._raw || row
  const cap = (s) => s ? s.charAt(0).toUpperCase() + s.slice(1) : ''
  newQuestion.value = {
    text: q.question_text || row.text || '',
    role: (Array.isArray(q.skill_tags) && q.skill_tags[0]) || 'All',
    level: cap(q.level) || 'Fresher',
    type: cap(q.question_type) || 'Technical',
    expected_signals: Array.isArray(q.expected_signals) ? q.expected_signals.join(', ') : '',
  }
  showCreateModal.value = true
}

const openCreate = () => {
  editingId.value = null
  newQuestion.value = { text: '', role: 'All', level: 'Fresher', type: 'Technical', expected_signals: '' }
  showCreateModal.value = true
}

const handleCreate = async () => {
  if (saving.value) return
  const values = { ...newQuestion.value, text: normalizeText(newQuestion.value.text) }
  const validation = validateForm(values, {
    text: [
      (value) => requiredTrim(value, 'Vui lòng nhập nội dung câu hỏi.'),
      (value) => minLength(value, 10, 'Câu hỏi phải có ít nhất 10 ký tự.'),
      (value) => maxLength(value, 2000, 'Câu hỏi không được vượt quá 2.000 ký tự.'),
    ],
    level: [(value) => isOneOf(value, ['Fresher', 'Junior', 'Middle', 'Senior'], 'Cấp độ không hợp lệ.')],
    type: [(value) => isOneOf(value, ['Technical', 'Behavioral', 'System Design', 'Custom'], 'Loại câu hỏi không hợp lệ.')],
  })
  if (!editingId.value && hasDuplicateNormalized([...questions.value.map(item => item.text), values.text])) {
    validation.errors.text = 'Câu hỏi này đã tồn tại trong kho.'
    validation.isValid = false
  }
  formErrors.value = validation.errors
  if (!validation.isValid) return
  newQuestion.value.text = values.text
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const payload = {
      question_text: newQuestion.value.text,
      question_type: newQuestion.value.type.toLowerCase(),
      level: newQuestion.value.level.toLowerCase(),
      skill_tags: newQuestion.value.role !== 'All' ? [newQuestion.value.role] : [],
      expected_signals: newQuestion.value.expected_signals
        ? newQuestion.value.expected_signals.split(',').map(s => s.trim())
        : []
    }
    if (editingId.value) {
      await questionBankService.updateQuestion(companyId, editingId.value, payload)
      toast.value = { type: 'success', message: 'Đã cập nhật câu hỏi!' }
    } else {
      await questionBankService.createQuestion(companyId, payload)
      toast.value = { type: 'success', message: 'Đã thêm câu hỏi vào kho!' }
    }
    showCreateModal.value = false
    editingId.value = null
    newQuestion.value = { text: '', role: 'All', level: 'Fresher', type: 'Technical', expected_signals: '' }
    await loadQuestions()
  } catch (err) {
    toast.value = { type: 'error', message: 'Lưu thất bại. Vui lòng thử lại.' }
  } finally {
    saving.value = false
  }
}

const confirmDelete = async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await questionBankService.deleteQuestion(companyId, deletingId.value)
    showDeleteModal.value = false
    toast.value = { type: 'success', message: 'Đã xóa câu hỏi khỏi kho!' }
    await loadQuestions()
  } catch (err) {
    toast.value = { type: 'error', message: 'Xóa thất bại. Vui lòng thử lại.' }
  }
}

const handleGenerateAI = async () => {
  const companyId = authStore.user?.companies?.[0]?.id
  if (!companyId) return
  const count = Number(generateCount.value)
  if (!Number.isInteger(count) || count < 1 || count > 20) {
    toast.value = { type: 'error', message: 'Số câu hỏi AI phải nằm trong khoảng từ 1 đến 20.' }
    return
  }
  if (!selectedJobId.value) {
    toast.value = { type: 'error', message: 'Vui lòng chọn công việc để AI dựa vào JD tạo câu hỏi.' }
    return
  }
  aiGenerating.value = true
  try {
    await questionBankService.generateWithAI(companyId, selectedJobId.value, { count, level: filterLevel.value || 'middle' })
    toast.value = { type: 'success', message: 'AI đã tạo thêm câu hỏi vào kho!' }
    await loadQuestions()
  } catch (err) {
    toast.value = { type: 'error', message: 'Tạo câu hỏi bằng AI thất bại. Vui lòng thử lại.' }
  } finally {
    aiGenerating.value = false
  }
}
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Kho câu hỏi</h1>
        <p class="text-helper" style="margin-top: 4px">Quản lý ngân hàng câu hỏi dùng chung cho các buổi phỏng vấn.</p>
      </div>
      <Button @click="openCreate"><Plus size="16" /> Thêm câu hỏi</Button>
    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <Card>
      <div style="display: flex; gap: 16px; margin-bottom: 24px">
        <div style="position: relative; flex: 1; max-width: 400px">
          <Search size="16" style="position: absolute; left: 12px; top: 12px; color: var(--text-muted)" />
          <input
            type="text"
            placeholder="Tìm kiếm nội dung câu hỏi..."
            class="input-field"
            style="width: 100%; padding-left: 36px"
            v-model="searchKeyword"
            @keyup.enter="loadQuestions"
          />
        </div>
        <select class="input-field" style="width: 220px" v-model="selectedJobId">
          <option value="">Chọn công việc cho AI</option>
          <option v-for="job in jobs" :key="job.id" :value="job.id">{{ job.title || job.name }}</option>
        </select>
        <select class="input-field" style="width: 180px" v-model="filterType" @change="loadQuestions">
          <option value="">Tất cả loại</option>
          <option value="technical">Technical</option>
          <option value="behavioral">Behavioral</option>
          <option value="system design">System Design</option>
        </select>
        <select class="input-field" style="width: 180px" v-model="filterLevel" @change="loadQuestions">
          <option value="">Tất cả cấp độ</option>
          <option value="fresher">Fresher</option>
          <option value="junior">Junior</option>
          <option value="middle">Middle</option>
          <option value="senior">Senior</option>
        </select>
        <Button variant="secondary" @click="loadQuestions"><Search size="16" /> Tìm</Button>
        <Button variant="secondary" @click="showFilterModal = true"><Filter size="16" /> Lọc nâng cao</Button>
      </div>

      <Table :columns="columns" :data="questions">
        <template #text="{ row }">
          <div style="font-weight: 500; color: var(--text-main); display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; text-overflow: ellipsis; max-width: 400px;" :title="row.text">
            {{ row.text }}
          </div>
        </template>
        <template #level="{ row }">
          <span class="badge badge-neutral">{{ row.level }}</span>
        </template>
        <template #action="{ row }">
          <div style="display: flex; gap: 8px">
            <Button variant="ghost" style="padding: 4px" @click="openEdit(row)"><Edit size="16" /></Button>
            <Button variant="ghost" style="padding: 4px; color: var(--danger)" @click="deletingId = row.id; showDeleteModal = true"><Trash2 size="16" /></Button>
          </div>
        </template>
      </Table>
    </Card>
    
    <div style="margin-top: 24px; padding: 16px; background-color: rgba(37, 99, 235, 0.05); border-radius: var(--radius); border: 1px solid rgba(37, 99, 235, 0.2); display: flex; align-items: flex-start; gap: 12px">
      <Layers size="20" color="var(--primary)" style="flex-shrink: 0; margin-top: 2px" />
      <div>
        <p class="text-body" style="font-weight: 600; color: var(--primary)">AI Suggestion Engine đang hoạt động</p>
        <p class="text-helper" style="color: var(--text-secondary)">Trong lúc phỏng vấn, AI sẽ tự động tìm kiếm các câu hỏi liên quan trong Kho câu hỏi này dựa trên ngữ cảnh để gợi ý cho bạn.</p>
      </div>
    </div>

    <Modal :isOpen="showCreateModal" @close="showCreateModal = false" :title="editingId ? 'Chỉnh sửa câu hỏi' : 'Thêm câu hỏi mới'">
      <Input label="Nội dung câu hỏi" v-model="newQuestion.text" placeholder="Nhập câu hỏi..." style="margin-bottom: 16px" />
      
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 24px">
        <div class="input-group">
          <label class="input-label">Vị trí (Role)</label>
          <select class="input-field" v-model="newQuestion.role">
            <option value="All">Tất cả vị trí</option>
            <option value="Frontend">Frontend</option>
            <option value="Backend">Backend</option>
            <option value="Fullstack">Fullstack</option>
          </select>
        </div>
        
        <div class="input-group">
          <label class="input-label">Cấp độ (Level)</label>
          <select class="input-field" v-model="newQuestion.level">
            <option value="Any">Bất kỳ</option>
            <option value="Fresher">Fresher</option>
            <option value="Junior">Junior</option>
            <option value="Middle">Middle</option>
            <option value="Senior">Senior</option>
          </select>
        </div>
        
        <div class="input-group" style="grid-column: 1 / -1">
          <label class="input-label">Loại câu hỏi (Type)</label>
          <select class="input-field" v-model="newQuestion.type">
            <option value="Technical">Technical</option>
            <option value="Behavioral">Behavioral</option>
            <option value="System Design">System Design</option>
            <option value="Custom">Khác (Custom)</option>
          </select>
        </div>
      </div>

      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showCreateModal = false">Hủy</Button>
        <Button variant="primary" :disabled="saving" @click="handleCreate">{{ saving ? 'Đang lưu...' : (editingId ? 'Cập nhật' : 'Lưu câu hỏi') }}</Button>
      </div>
    </Modal>

    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn xóa câu hỏi này? Hành động này không thể hoàn tác.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showDeleteModal = false">Hủy</Button>
        <Button variant="primary" style="background-color: var(--danger); border-color: var(--danger)" @click="confirmDelete">Xóa câu hỏi</Button>
      </div>
    </Modal>

    <Modal :isOpen="showFilterModal" @close="showFilterModal = false" title="Bộ lọc nâng cao">
      <div style="display: flex; flex-direction: column; gap: 16px; margin-bottom: 24px">
        <div class="input-group">
          <label class="input-label">Loại câu hỏi (Type)</label>
          <select class="input-field" v-model="filterType">
            <option value="">Tất cả</option>
            <option value="technical">Technical</option>
            <option value="behavioral">Behavioral</option>
            <option value="system design">System Design</option>
          </select>
        </div>
        <div class="input-group">
          <label class="input-label">Cấp độ (Level)</label>
          <select class="input-field" v-model="filterLevel">
            <option value="">Tất cả cấp độ</option>
            <option value="fresher">Fresher</option>
            <option value="junior">Junior</option>
            <option value="middle">Middle</option>
            <option value="senior">Senior</option>
          </select>
        </div>
        <div class="input-group">
          <label class="input-label">Từ khóa</label>
          <Input v-model="searchKeyword" placeholder="Nhập từ khóa nội dung câu hỏi..." />
        </div>
      </div>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="searchKeyword = ''; filterType = ''; filterLevel = ''; showFilterModal = false; loadQuestions()">Xóa bộ lọc</Button>
        <Button variant="primary" @click="showFilterModal = false; loadQuestions()">Áp dụng</Button>
      </div>
    </Modal>
  </div>
</template>
