<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import Toast from '../../components/common/AppToast.vue'
import { ArrowLeft, Save, Sparkles, AlertCircle, CheckCircle, Users } from 'lucide-vue-next'
import { jobService } from '../../services/job.service'
import { authStore } from '../../stores/auth.store'
import { JOB_STATUSES, JOB_LEVELS, EMPLOYMENT_TYPES } from '../../utils/constants.js'
import { isOneOf, isValidSalaryRange, maxLength, minLength, normalizeText, requiredTrim, validateForm } from '../../utils/validators.js'

const route = useRoute()
const router = useRouter()
const id = route.params.id
const isNew = ref(id === 'new')

const job = ref({
  title: '',
  location: '',
  department: '',
  level: '',
  employment_type: '',
  status: 'open',
  description: '',
  requirements: '',
  benefits: '',
  salary_min: '',
  salary_max: ''
})
const loading = ref(!isNew.value)
const saving = ref(false)
const aiAnalyzing = ref(false)
const rubric = ref(null)
const localToast = ref(null)
const errors = ref({})

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object' && 'String' in val) {
    return val.Valid ? val.String : ''
  }
  return val
}

onMounted(async () => {
  if (!isNew.value) {
    try {
      const companyId = authStore.user?.companies?.[0]?.id
      if (companyId) {
        const data = await jobService.getJob(companyId, id)
        const rawJob = data.data || data
        job.value = {
          ...rawJob,
          status: rawJob.status === 'open' ? 'Open' : (rawJob.status === 'closed' ? 'Closed' : 'Draft'),
          requirements: unwrap(rawJob.requirements),
          benefits: unwrap(rawJob.benefits)
        }
        
        if (rawJob.ai_analysis_json) {
          let result = typeof rawJob.ai_analysis_json === 'string' ? JSON.parse(rawJob.ai_analysis_json) : rawJob.ai_analysis_json
          aiResult.value = result
          rubric.value = (result?.suggested_rubric || []).map(r => ({
            criterion: r.name,
            weight: `${r.weight}%`
          }))
        }
      }
    } catch (error) {
      localToast.value = { type: 'error', message: 'Không thể tải chi tiết công việc' }
    } finally {
      loading.value = false
    }
  }
})

const handleSave = async (e) => {
  e.preventDefault()
  if (saving.value) return

  const normalizedJob = { ...job.value, title: normalizeText(job.value.title), location: normalizeText(job.value.location) }
  const validation = validateForm(normalizedJob, {
    title: [
      (value) => requiredTrim(value, 'Vui lòng nhập tiêu đề công việc.'),
      (value) => minLength(value, 2, 'Tiêu đề công việc phải có ít nhất 2 ký tự.'),
      (value) => maxLength(value, 255, 'Tiêu đề công việc không được vượt quá 255 ký tự.'),
    ],
    location: [
      (value) => requiredTrim(value, 'Vui lòng nhập vị trí làm việc.'),
      (value) => minLength(value, 2, 'Vị trí làm việc phải có ít nhất 2 ký tự.'),
      (value) => maxLength(value, 255, 'Vị trí làm việc không được vượt quá 255 ký tự.'),
    ],
    department: [(value) => maxLength(value, 255, 'Phòng ban không được vượt quá 255 ký tự.')],
    level: [(value) => isOneOf(String(value || '').toLowerCase(), JOB_LEVELS, 'Cấp độ công việc không hợp lệ.')],
    employment_type: [(value) => isOneOf(String(value || '').toLowerCase(), EMPLOYMENT_TYPES, 'Loại hình làm việc không hợp lệ.')],
    description: [
      (value) => requiredTrim(value, 'Vui lòng nhập mô tả công việc.'),
      (value) => minLength(value, 10, 'Mô tả công việc phải có ít nhất 10 ký tự.'),
      (value) => maxLength(value, 20000, 'Mô tả công việc không được vượt quá 20.000 ký tự.'),
    ],
    requirements: [(value) => maxLength(value, 10000, 'Yêu cầu không được vượt quá 10.000 ký tự.')],
    benefits: [(value) => maxLength(value, 10000, 'Quyền lợi không được vượt quá 10.000 ký tự.')],
    status: [(value) => isOneOf(String(value || '').toLowerCase(), JOB_STATUSES, 'Trạng thái công việc không hợp lệ.')],
  })
  const salaryError = isValidSalaryRange(job.value.salary_min, job.value.salary_max)
  errors.value = { ...validation.errors, ...(salaryError ? { salary: salaryError } : {}) }
  if (!validation.isValid || salaryError) return
  job.value.title = normalizedJob.title
  job.value.location = normalizedJob.location

  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const payload = {
      title: job.value.title,
      location: job.value.location,
      description: job.value.description,
      requirements: job.value.requirements,
      benefits: job.value.benefits,
      status: String(job.value.status || '').toLowerCase(),
      ...(job.value.department ? { department: normalizeText(job.value.department) } : {}),
      ...(job.value.level ? { level: String(job.value.level).toLowerCase() } : {}),
      ...(job.value.employment_type ? { employment_type: String(job.value.employment_type).toLowerCase() } : {}),
      ...(job.value.salary_min !== '' ? { salary_min: Number(job.value.salary_min) } : {}),
      ...(job.value.salary_max !== '' ? { salary_max: Number(job.value.salary_max) } : {}),
    }
    
    if (isNew.value) {
      await jobService.createJob(companyId, payload)
    } else {
      await jobService.updateJob(companyId, id, payload)
    }
    
    router.push({ path: '/jobs', state: { message: 'Lưu thông tin công việc thành công!' } })
  } catch (error) {
    localToast.value = { type: 'error', message: error.message || 'Lưu công việc thất bại' }
  } finally {
    saving.value = false
  }
}

const aiResult = ref(null)

const handleAiAnalyze = async () => {
  if (aiAnalyzing.value) return
  const titleError = requiredTrim(job.value.title, 'Vui lòng nhập tiêu đề công việc.') || minLength(job.value.title, 2, 'Tiêu đề công việc phải có ít nhất 2 ký tự.')
  const locationError = requiredTrim(job.value.location, 'Vui lòng nhập vị trí làm việc.') || minLength(job.value.location, 2, 'Vị trí làm việc phải có ít nhất 2 ký tự.')
  const descriptionError = requiredTrim(job.value.description, 'Vui lòng nhập mô tả công việc.') || minLength(job.value.description, 10, 'Mô tả công việc phải có ít nhất 10 ký tự.')
  errors.value = {
    ...(titleError ? { title: titleError } : {}),
    ...(locationError ? { location: locationError } : {}),
    ...(descriptionError ? { description: descriptionError } : {}),
  }
  if (Object.keys(errors.value).length) {
    localToast.value = { type: 'error', message: 'Vui lòng hoàn thiện JD hợp lệ trước khi yêu cầu AI phân tích.' }
    return
  }
  if (isNew.value || !id) {
    localToast.value = { type: 'error', message: 'Vui lòng lưu công việc trước khi yêu cầu AI phân tích.' }
    return
  }

  aiAnalyzing.value = true
  aiResult.value = null
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    // API_SPEC §4.6 — POST /companies/:company_id/jobs/:job_id/analyze
    await jobService.analyzeJD(companyId, id, false)
    
    // AIAnalysis is saved to the database, we need to fetch the job again
    const jobData = await jobService.getJob(companyId, id)
    const rawJob = jobData.data || jobData
    
    if (rawJob.ai_analysis_json) {
      let result = typeof rawJob.ai_analysis_json === 'string' ? JSON.parse(rawJob.ai_analysis_json) : rawJob.ai_analysis_json
      aiResult.value = result
      rubric.value = (result?.suggested_rubric || []).map(r => ({
        criterion: r.name,
        weight: `${r.weight}%`
      }))
    }
    
    localToast.value = { type: 'success', message: 'AI đã phân tích JD thành công!' }
  } catch (err) {
    localToast.value = { type: 'info', message: 'Tính năng AI phân tích đang chờ Backend của Khôi.' }
  } finally {
    aiAnalyzing.value = false
  }
}
</script>

<template>
  <div v-if="loading" class="flex flex-col items-center justify-center min-h-[400px] text-slate-500 dark:text-slate-400">
    <div class="w-10 h-10 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mb-4"></div>
    <span class="font-medium">Đang tải chi tiết công việc...</span>
  </div>
  <div v-else class="animate-fade-in space-y-6">
    <Toast v-if="localToast" :type="localToast.type" :message="localToast.message" @close="localToast = null" />

    <!-- Header -->
    <Card class="flex flex-col md:flex-row md:items-center gap-4 p-6 rounded-2xl shadow-sm">
      <button @click="router.push('/jobs')" class="p-2 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-500 dark:text-slate-400 rounded-lg transition-colors border border-transparent hover:border-slate-200 dark:hover:border-slate-600 shrink-0">
        <ArrowLeft size="20" />
      </button>
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-2xl font-bold text-slate-800 dark:text-slate-100">{{ isNew ? 'Tạo Job Mới' : job.title }}</h1>
          <span v-if="!isNew" class="px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider rounded-full border" :class="job.status === 'Open' ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20' : 'bg-slate-100 text-slate-600 border-slate-200 dark:bg-slate-700 dark:text-slate-300 dark:border-slate-600'">
            {{ job.status }}
          </span>
        </div>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">Jobs > {{ isNew ? 'New' : job.title }}</p>
      </div>
    </Card>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-6">
        <!-- Main Form Card -->
        <Card class="rounded-2xl shadow-sm p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-6">Thông tin chung</h2>
          <form @submit="handleSave" class="space-y-5">
            <!-- Job Title -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Tiêu đề công việc</label>
              <input
                v-model="job.title"
                required
                minlength="2"
                :class="['w-full px-4 py-2.5 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)] placeholder:text-[var(--text-secondary)]', { 'input-error': errors.title }]"
                placeholder="Ví dụ: Frontend Developer"
              />
              <span v-if="errors.title" class="error-text">{{ errors.title }}</span>
            </div>
            
            <!-- Job Location -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Vị trí làm việc</label>
              <input
                v-model="job.location"
                required
                placeholder="Ví dụ: Hà Nội, TP. Hồ Chí Minh hoặc Remote"
                :class="['w-full px-4 py-2.5 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)] placeholder:text-[var(--text-secondary)]', { 'input-error': errors.location }]"
              />
              <span v-if="errors.location" class="error-text">{{ errors.location }}</span>
            </div>

            <!-- Status -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Trạng thái</label>
              <select 
                v-model="job.status"
                :class="['w-full px-4 py-2.5 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)]', { 'input-error': errors.status }]"
              >
                <option value="Open">Đang mở (Open)</option>
                <option value="Closed">Đã đóng (Closed)</option>
              </select>
              <span v-if="errors.status" class="error-text">{{ errors.status }}</span>
            </div>

            <!-- Job Description -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Mô tả công việc (JD)</label>
              <textarea 
                rows="6"
                v-model="job.description"
                placeholder="Nhập mô tả tổng quan về công việc (tối thiểu 10 ký tự)..."
                required
                minlength="10"
                :class="['w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)] placeholder:text-[var(--text-secondary)] resize-y', { 'input-error': errors.description }]"
              ></textarea>
              <span v-if="errors.description" class="error-text">{{ errors.description }}</span>
            </div>

            <!-- Requirements -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Yêu cầu ứng viên <span class="text-red-500">*</span></label>
              <textarea 
                rows="6"
                v-model="job.requirements"
                placeholder="Nhập yêu cầu về kỹ năng, kinh nghiệm..."
                required
                :class="['w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)] placeholder:text-[var(--text-secondary)] resize-y', { 'input-error': errors.requirements }]"
              ></textarea>
              <span v-if="errors.requirements" class="error-text">{{ errors.requirements }}</span>
            </div>

            <!-- Benefits -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-slate-700 dark:text-slate-300">Quyền lợi <span class="text-red-500">*</span></label>
              <textarea 
                rows="6"
                v-model="job.benefits"
                placeholder="Nhập các quyền lợi, chế độ đãi ngộ..."
                required
                :class="['w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-[var(--text-primary)] placeholder:text-[var(--text-secondary)] resize-y', { 'input-error': errors.benefits }]"
              ></textarea>
              <span v-if="errors.benefits" class="error-text">{{ errors.benefits }}</span>
            </div>

            <div class="flex justify-end pt-4 gap-3 border-t border-slate-100 dark:border-slate-700">
              <Button type="button" variant="ghost" @click="router.push('/jobs')" class="text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700">Hủy</Button>
              <Button type="submit" :disabled="saving" class="bg-blue-600 hover:bg-blue-700 text-white border-none shadow-md shadow-blue-500/20">
                <Save size="16" class="mr-2" v-if="!saving" /> 
                <div v-else class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin mr-2"></div>
                {{ saving ? 'Đang lưu...' : 'Lưu thông tin' }}
              </Button>
            </div>
          </form>
        </Card>
        
        <Card v-if="!isNew" class="rounded-2xl shadow-sm p-6">
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-4">Danh sách ứng viên</h2>
          <div class="text-center text-slate-500 dark:text-slate-400 py-12">
            <Users size="48" class="mx-auto mb-4 text-slate-300 dark:text-slate-600" />
            <p class="font-medium text-sm">Chưa có ứng viên nào nộp đơn.</p>
          </div>
        </Card>
      </div>

      <!-- Right Column -->
      <div class="space-y-6">
        <!-- AI Analysis Card -->
        <Card class="rounded-2xl shadow-sm p-6 border-t-4 border-t-blue-500 relative overflow-hidden">
          <div class="absolute -right-6 -top-6 text-blue-500/10 pointer-events-none">
            <Sparkles size="100" />
          </div>
          <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 mb-3 relative z-10">AI Phân tích JD</h2>
          
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-5 relative z-10">
            AI sẽ phân tích JD và đề xuất kỹ năng, bộ tiêu chí (Rubric) và câu hỏi phù hợp nhất.
          </p>
          
          <button
            class="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-sm font-semibold transition-all relative z-10 border"
            :class="[
              aiAnalyzing || !job.description 
                ? 'bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-400 cursor-not-allowed'
                : 'bg-blue-50 hover:bg-blue-100 dark:bg-blue-500/10 dark:hover:bg-blue-500/20 border-blue-200 dark:border-blue-500/30 text-blue-600 dark:text-blue-400'
            ]"
            @click="handleAiAnalyze"
            :disabled="aiAnalyzing || !job.description"
          >
            <Sparkles size="16" v-if="!aiAnalyzing" />
            <div v-else class="w-4 h-4 border-2 border-blue-500/30 border-t-blue-500 rounded-full animate-spin"></div>
            {{ aiAnalyzing ? 'Đang phân tích...' : 'Phân tích JD bằng AI' }}
          </button>

          <!-- AI result: Summary -->
          <div v-if="aiResult" class="mt-6 pt-5 border-t border-slate-100 dark:border-slate-700 flex flex-col gap-5 relative z-10 animate-fade-in">
            <!-- Summary Box -->
            <div v-if="aiResult.summary" class="bg-blue-50/50 dark:bg-blue-500/5 rounded-xl p-4 border border-blue-100/50 dark:border-blue-500/10">
              <p class="text-xs font-bold text-blue-600 dark:text-blue-400 uppercase tracking-wider mb-2">Tóm tắt AI</p>
              <p class="text-sm text-slate-700 dark:text-slate-300 leading-relaxed">{{ aiResult.summary }}</p>
            </div>

            <!-- Required skills -->
            <div v-if="aiResult.required_skills?.length">
              <p class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-3">Kỹ năng bắt buộc</p>
              <div class="flex flex-wrap gap-2">
                <span v-for="s in aiResult.required_skills" :key="s"
                  class="px-3 py-1 rounded-full bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400 text-xs font-semibold border border-blue-100 dark:border-blue-500/20">
                  {{ s }}
                </span>
              </div>
            </div>

            <!-- Nice-to-have skills -->
            <div v-if="aiResult.nice_to_have_skills?.length">
              <p class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-3">Kỹ năng ưu tiên</p>
              <div class="flex flex-wrap gap-2">
                <span v-for="s in aiResult.nice_to_have_skills" :key="s"
                  class="px-3 py-1 rounded-full bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 text-xs font-semibold border border-emerald-100 dark:border-emerald-500/20">
                  {{ s }}
                </span>
              </div>
            </div>

            <!-- Rubric suggestions -->
            <div v-if="rubric?.length">
              <div class="flex items-center gap-2 mb-3">
                <CheckCircle size="14" class="text-emerald-500" />
                <p class="text-xs font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider">Tiêu chí đề xuất</p>
              </div>
              <div class="space-y-2">
                <div v-for="(r, i) in rubric" :key="i"
                  class="flex justify-between items-center p-3 bg-slate-50 dark:bg-slate-900/50 border border-slate-100 dark:border-slate-700/50 rounded-lg">
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-300">{{ r.criterion }}</span>
                  <span class="text-sm font-bold text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-500/10 px-2 py-0.5 rounded">{{ r.weight }}</span>
                </div>
              </div>
            </div>

            <!-- Suggested questions count -->
            <div v-if="aiResult.suggested_questions?.length" class="mt-2">
              <div class="flex justify-between items-center p-3.5 bg-violet-50 dark:bg-violet-500/10 border border-violet-100 dark:border-violet-500/20 rounded-xl">
                <p class="text-sm font-semibold text-violet-600 dark:text-violet-400">{{ aiResult.suggested_questions.length }} câu hỏi đã được đề xuất</p>
                <CheckCircle size="16" class="text-violet-500" />
              </div>
            </div>
          </div>

          <div v-if="!job.description && !aiResult" class="mt-5 flex items-center gap-2 text-amber-500 bg-amber-50 dark:bg-amber-500/10 p-3 rounded-lg border border-amber-100 dark:border-amber-500/20 text-sm relative z-10 font-medium">
            <AlertCircle size="16" class="shrink-0" /> Cần nhập JD để AI có thể phân tích.
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>
