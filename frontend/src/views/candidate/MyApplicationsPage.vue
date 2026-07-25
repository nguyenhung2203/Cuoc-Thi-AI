<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { candidatePortalService } from '../../services/candidate-portal.service'
import Card from '../../components/common/AppCard.vue'
import Badge from '../../components/common/AppBadge.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import { Briefcase, MapPin, Building2, Calendar, FileText, CheckCircle2, XCircle, Clock, Search, ChevronRight, Download, Trash2, Filter, AlertCircle } from 'lucide-vue-next'
import { langStore } from '../../stores/lang.store'

const router = useRouter()
const applications = ref([])
const loading = ref(true)
const searchQuery = ref('')
const currentTab = ref('all') // all, pending, reviewed, interviewing, rejected
const toast = ref(null)

onMounted(async () => {
  try {
    applications.value = await candidatePortalService.getApplications()
  } catch (err) {
    console.error('Lỗi tải danh sách ứng tuyển:', err)
  } finally {
    loading.value = false
  }
})

// Đếm số lượng theo từng trạng thái
const tabCounts = computed(() => {
  const counts = { all: applications.value.length, pending: 0, reviewed: 0, interviewing: 0, rejected: 0 }
  applications.value.forEach(app => {
    if (counts[app.status] !== undefined) {
      counts[app.status]++
    }
  })
  return counts
})

const filteredApps = computed(() => {
  return applications.value.filter(app => {
    const matchesTab = currentTab.value === 'all' || app.status === currentTab.value
    const matchesSearch = !searchQuery.value || 
      app.job_title.toLowerCase().includes(searchQuery.value.toLowerCase()) || 
      app.company_name.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchesTab && matchesSearch
  })
})

const getStatusDetails = (status) => {
  const map = {
    'pending': { text: 'Đang chờ duyệt', class: 'badge-neutral', icon: Clock, desc: 'Hồ sơ đã được gửi đến bộ phận nhân sự và đang trong quá trình tiếp nhận.' },
    'reviewed': { text: 'HR đã xem hồ sơ', class: 'badge-info', icon: CheckCircle2, desc: 'Nhà tuyển dụng đã mở xem CV và hồ sơ năng lực của bạn.' },
    'interviewing': { text: 'Đang phỏng vấn / Test', class: 'badge-primary', icon: Calendar, desc: 'Bạn đã vượt qua vòng hồ sơ và đang tham gia phỏng vấn đánh giá.' },
    'rejected': { text: 'Chưa phù hợp', class: 'badge-danger', icon: XCircle, desc: 'Nhà tuyển dụng đã phản hồi hồ sơ chưa phù hợp với vị trí lúc này.' }
  }
  return map[status] || map['pending']
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return d.toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const handleDownloadCv = (app) => {
  toast.value = { type: 'info', message: `Đang tải xuống CV: ${app.cv_name}...` }
}

const handleCancelApp = async (app) => {
  if (confirm(`Bạn có chắc chắn muốn rút hồ sơ ứng tuyển vị trí "${app.job_title}" tại ${app.company_name}?`)) {
    try {
      await candidatePortalService.cancelApplication(app.id)
      applications.value = applications.value.filter(item => item.id !== app.id)
      toast.value = { type: 'success', message: 'Đã rút/xóa hồ sơ ứng tuyển thành công.' }
    } catch (err) {
      toast.value = { type: 'error', message: err.message || 'Không thể rút hồ sơ ứng tuyển.' }
    }
  }
}
</script>

<template>
  <div class="applications-page space-y-6">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- Khung bọc tiêu đề chuẩn Design System -->
    <div class="header-box animate-rise flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 rounded-2xl bg-[var(--surface)] border border-[var(--border)] shadow-sm gap-4">
      <div>
        <h1 class="text-h1">Việc làm đã ứng tuyển</h1>
        <p class="text-secondary mt-1">Quản lý và theo dõi tiến độ chi tiết từng hồ sơ bạn đã nộp cho doanh nghiệp.</p>
      </div>
      <Button variant="primary" @click="router.push('/job-board')" class="flex items-center gap-2 shrink-0">
        <Search :size="16" /> Khám phá việc làm mới
      </Button>
    </div>

    <!-- Thanh tìm kiếm và bộ lọc nhanh -->
    <div class="bg-[var(--surface)] p-5 rounded-2xl border border-[var(--border)] shadow-sm space-y-4">
      <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
        <!-- Input search -->
        <div class="relative flex-1">
          <Search :size="18" class="absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-muted)] pointer-events-none" />
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Tìm kiếm theo vị trí công việc hoặc tên công ty..." 
            class="w-full pl-10 pr-4 py-2.5 bg-[var(--surface-soft)] border border-[var(--border)] rounded-xl text-sm font-medium text-[var(--text-main)] focus:outline-none focus:border-[var(--primary)] transition-colors"
          />
        </div>
      </div>

      <!-- Tabs chấu mượt mà -->
      <div class="flex items-center gap-2 overflow-x-auto pb-1 border-t border-[var(--border)] pt-4 scrollbar-none">
        <button 
          @click="currentTab = 'all'" 
          :class="['px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 transition-all whitespace-nowrap', currentTab === 'all' ? 'bg-[var(--primary)] text-white shadow-sm' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)] hover:text-[var(--text-main)]']"
        >
          Tất cả <span :class="['px-2 py-0.5 rounded-full text-xs font-bold', currentTab === 'all' ? 'bg-white/20 text-white' : 'bg-slate-200/60 text-[var(--text-secondary)]']">{{ tabCounts.all }}</span>
        </button>

        <button 
          @click="currentTab = 'pending'" 
          :class="['px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 transition-all whitespace-nowrap', currentTab === 'pending' ? 'bg-[var(--primary)] text-white shadow-sm' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)] hover:text-[var(--text-main)]']"
        >
          <Clock :size="15" /> Đang chờ duyệt <span :class="['px-2 py-0.5 rounded-full text-xs font-bold', currentTab === 'pending' ? 'bg-white/20 text-white' : 'bg-slate-200/60 text-[var(--text-secondary)]']">{{ tabCounts.pending }}</span>
        </button>

        <button 
          @click="currentTab = 'reviewed'" 
          :class="['px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 transition-all whitespace-nowrap', currentTab === 'reviewed' ? 'bg-[var(--primary)] text-white shadow-sm' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)] hover:text-[var(--text-main)]']"
        >
          <CheckCircle2 :size="15" /> HR đã xem <span :class="['px-2 py-0.5 rounded-full text-xs font-bold', currentTab === 'reviewed' ? 'bg-white/20 text-white' : 'bg-slate-200/60 text-[var(--text-secondary)]']">{{ tabCounts.reviewed }}</span>
        </button>

        <button 
          @click="currentTab = 'interviewing'" 
          :class="['px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 transition-all whitespace-nowrap', currentTab === 'interviewing' ? 'bg-[var(--primary)] text-white shadow-sm' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)] hover:text-[var(--text-main)]']"
        >
          <Calendar :size="15" /> Đang phỏng vấn <span :class="['px-2 py-0.5 rounded-full text-xs font-bold', currentTab === 'interviewing' ? 'bg-white/20 text-white' : 'bg-slate-200/60 text-[var(--text-secondary)]']">{{ tabCounts.interviewing }}</span>
        </button>

        <button 
          @click="currentTab = 'rejected'" 
          :class="['px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-2 transition-all whitespace-nowrap', currentTab === 'rejected' ? 'bg-[var(--primary)] text-white shadow-sm' : 'bg-[var(--surface-soft)] text-[var(--text-secondary)] hover:text-[var(--text-main)]']"
        >
          <XCircle :size="15" /> Chưa phù hợp <span :class="['px-2 py-0.5 rounded-full text-xs font-bold', currentTab === 'rejected' ? 'bg-white/20 text-white' : 'bg-slate-200/60 text-[var(--text-secondary)]']">{{ tabCounts.rejected }}</span>
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center p-12">
      <div class="spinner"></div>
    </div>

    <!-- Empty state chuẩn Design System -->
    <div v-else-if="filteredApps.length === 0" class="empty-state bg-[var(--surface)] p-12 rounded-2xl border border-dashed border-[var(--border)] text-center max-w-xl mx-auto my-8">
      <div class="w-16 h-16 rounded-full bg-[var(--surface-soft)] text-[var(--text-muted)] flex items-center justify-center mx-auto mb-4">
        <FileText :size="32" />
      </div>
      <h3 class="text-base font-bold text-[var(--text-main)]">Không tìm thấy hồ sơ ứng tuyển nào</h3>
      <p class="text-sm text-[var(--text-secondary)] mt-1.5 leading-relaxed">Bạn chưa nộp hồ sơ vào vị trí nào trong danh mục này hoặc từ khóa tìm kiếm chưa khớp.</p>
      <div class="flex items-center justify-center gap-3 mt-6">
        <Button v-if="searchQuery || currentTab !== 'all'" variant="ghost" @click="searchQuery = ''; currentTab = 'all'">Xóa bộ lọc</Button>
        <Button variant="primary" @click="router.push('/job-board')">Khám phá việc làm ngay</Button>
      </div>
    </div>

    <!-- Danh sách hồ sơ -->
    <div v-else class="space-y-4">
      <Card v-for="app in filteredApps" :key="app.id" class="p-6 rounded-2xl border border-[var(--border)] hover:border-[var(--primary)] transition-all duration-300 shadow-sm animate-rise">
        <div class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-6">
          
          <!-- Thông tin chính -->
          <div class="space-y-3 flex-1">
            <div class="flex items-start justify-between lg:justify-start gap-4">
              <h3 
                class="text-lg font-bold text-[var(--text-main)] hover:text-[var(--primary)] transition-colors cursor-pointer"
                @click="router.push(`/careers/${app.company_id}/jobs/${app.job_id}`)"
              >
                {{ app.job_title }}
              </h3>
            </div>

            <div class="flex flex-wrap items-center gap-4 text-sm text-[var(--text-secondary)]">
              <span class="flex items-center gap-1.5 font-semibold text-[var(--text-main)]">
                <Building2 :size="16" class="text-[var(--primary)]" /> {{ app.company_name }}
              </span>
              <span class="text-[var(--border)]">•</span>
              <span class="flex items-center gap-1.5">
                <Calendar :size="16" class="text-[var(--text-muted)]" /> Ngày nộp: {{ formatDate(app.applied_at) }}
              </span>
            </div>

            <!-- CV và mô tả trạng thái -->
            <div class="flex flex-wrap items-center gap-3 pt-1">
              <div class="inline-flex items-center gap-2 px-3 py-1.5 bg-[var(--surface-soft)] rounded-xl border border-[var(--border)] text-xs font-medium text-[var(--text-main)]">
                <FileText :size="14" class="text-[var(--primary)]" />
                CV đã nộp: <strong class="text-[var(--primary)]">{{ app.cv_name }}</strong>
                <button @click="handleDownloadCv(app)" class="hover:text-[var(--primary)] transition-colors ml-1" title="Tải xuống CV">
                  <Download :size="13" />
                </button>
              </div>

              <span class="text-xs text-[var(--text-secondary)] italic flex items-center gap-1.5">
                <AlertCircle :size="13" /> {{ getStatusDetails(app.status).desc }}
              </span>
            </div>
          </div>

          <!-- Trạng thái & Thao tác -->
          <div class="flex flex-col sm:flex-row lg:flex-col items-start sm:items-center lg:items-end justify-between w-full lg:w-auto gap-4 border-t lg:border-t-0 border-[var(--border)] pt-4 lg:pt-0">
            <Badge :class="[getStatusDetails(app.status).class, 'px-3 py-1.5 text-xs font-bold uppercase tracking-wider flex items-center gap-1.5 shadow-sm']">
              <component :is="getStatusDetails(app.status).icon" :size="14" />
              {{ getStatusDetails(app.status).text }}
            </Badge>

            <div class="flex items-center gap-2">
              <Button variant="ghost" size="sm" @click="router.push(`/careers/${app.company_id}/jobs/${app.job_id}`)" class="text-xs">
                Xem tin <ChevronRight :size="14" />
              </Button>
              <Button variant="ghost" size="sm" @click="handleCancelApp(app)" class="text-xs text-[var(--danger)] hover:bg-red-50 border border-transparent hover:border-red-200" title="Rút hồ sơ">
                <Trash2 :size="14" />
              </Button>
            </div>
          </div>

        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.applications-page {
  animation: fadeIn 0.4s ease-out;
}
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
.scrollbar-none {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
