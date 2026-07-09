<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import { Target, TrendingUp, Eye } from 'lucide-vue-next'
import { mockService } from '../../services/mock.service'

const router = useRouter()
const history = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    // API_SPEC §11.5 — GET /mock-interviews/my
    const data = await mockService.listMyMockInterviews()
    // Map API response fields sang format hiển thị
    history.value = (Array.isArray(data) ? data : []).map(item => ({
      id: item.id,
      role: item.target_role || 'Không rõ',
      level: item.target_level || '—',
      date: item.created_at ? new Date(item.created_at).toLocaleDateString('vi-VN') : '—',
      score: item.final_score != null ? Number(item.final_score).toFixed(1) : '—',
      status: item.status || 'completed'
    }))
  } catch (err) {
    console.error('Lỗi tải lịch sử mock', err)
    // Giữ nguyên empty array — Backend chưa sẵn sàng
    history.value = []
  } finally {
    loading.value = false
  }
})

const averageScore = computed(() => {
  const scored = history.value.filter(h => h.score !== '—')
  if (scored.length === 0) return 0
  const sum = scored.reduce((acc, curr) => acc + Number(curr.score), 0)
  return (sum / scored.length).toFixed(1)
})

const columns = [
  { header: 'Vị trí luyện tập', key: 'role' },
  { header: 'Cấp độ', key: 'level' },
  { header: 'Ngày thực hiện', key: 'date' },
  { header: 'Điểm số', key: 'score' },
  { header: 'Hành động', key: 'action' }
]
</script>

<template>
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-indigo-600">Kết quả & Lịch sử luyện tập</h1>
      <p class="text-gray-500 mt-2 text-lg">Theo dõi sự tiến bộ của bạn qua các bài luyện tập với AI.</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <!-- Average Score Card -->
      <Card class="bg-white/90 backdrop-blur-md shadow-lg border-0 rounded-2xl overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 relative">
        <div class="absolute -right-6 -top-6 w-24 h-24 bg-emerald-100 rounded-full opacity-50 blur-2xl pointer-events-none"></div>
        <div class="p-6">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-emerald-400 to-teal-500 flex items-center justify-center text-white shadow-lg shadow-emerald-500/30 shrink-0">
              <Target class="w-7 h-7" />
            </div>
            <div>
              <p class="text-sm font-bold text-gray-500 uppercase tracking-wider mb-1">Điểm trung bình</p>
              <h2 class="text-3xl font-black text-gray-800 flex items-baseline gap-1">
                {{ averageScore }}<span class="text-lg font-bold text-gray-400">/10</span>
              </h2>
            </div>
          </div>
        </div>
      </Card>
      
      <!-- Practice Count Card -->
      <Card class="bg-white/90 backdrop-blur-md shadow-lg border-0 rounded-2xl overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 relative">
        <div class="absolute -right-6 -top-6 w-24 h-24 bg-blue-100 rounded-full opacity-50 blur-2xl pointer-events-none"></div>
        <div class="p-6">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white shadow-lg shadow-blue-500/30 shrink-0">
              <TrendingUp class="w-7 h-7" />
            </div>
            <div>
              <p class="text-sm font-bold text-gray-500 uppercase tracking-wider mb-1">Số lần luyện tập</p>
              <h2 class="text-3xl font-black text-gray-800 flex items-baseline gap-1">
                {{ history.length }}<span class="text-lg font-bold text-gray-400">lần</span>
              </h2>
            </div>
          </div>
        </div>
      </Card>

      <!-- Improvement Areas Card -->
      <Card class="bg-gradient-to-br from-amber-50 to-orange-50 backdrop-blur-md shadow-lg border border-amber-100 rounded-2xl overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 relative">
        <div class="absolute -right-6 -top-6 w-24 h-24 bg-amber-200 rounded-full opacity-30 blur-2xl pointer-events-none"></div>
        <div class="p-6">
          <p class="text-sm font-bold text-amber-800 uppercase tracking-wider mb-3 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z" clip-rule="evenodd" /></svg>
            Điểm cần khắc phục
          </p>
          <ul class="space-y-2 text-amber-700 text-sm font-medium">
            <li class="flex items-start gap-2">
              <span class="w-1.5 h-1.5 bg-amber-400 rounded-full mt-1.5 shrink-0"></span>
              Trình bày cấu trúc câu trả lời (STAR)
            </li>
            <li class="flex items-start gap-2">
              <span class="w-1.5 h-1.5 bg-amber-400 rounded-full mt-1.5 shrink-0"></span>
              Đưa ra thêm nhiều ví dụ thực tế
            </li>
          </ul>
        </div>
      </Card>
    </div>

    <!-- History Table Area -->
    <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden p-0">
      <div class="p-6 border-b border-gray-100 bg-gray-50/50">
        <h2 class="text-xl font-bold text-gray-800">Lịch sử bài luyện tập</h2>
      </div>
      
      <div v-if="loading" class="p-12 text-center flex flex-col items-center">
        <div class="w-10 h-10 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mb-4"></div>
        <p class="text-gray-500 font-medium">Đang tải lịch sử...</p>
      </div>
      
      <div v-else class="p-0 overflow-x-auto">
        <!-- Replace AppTable with a beautifully styled custom table wrapper if AppTable doesn't support Tailwind injection easily -->
        <!-- But assuming AppTable just renders a standard table that we can style from outside or it passes through classes. -->
        <!-- For maximum visual impact matching the requirement, we will wrap it or rely on scoped styles if needed. -->
        <!-- In Vue, if AppTable renders standard HTML table tags, the scoped CSS below will hit them. -->
        <div class="modern-table-container">
          <Table :columns="columns" :data="history" class="w-full text-left border-collapse">
            <template #date="{ row }">
              <span class="text-gray-600 font-medium">{{ new Date(row.date).toLocaleDateString('vi-VN') }}</span>
            </template>
            <template #score="{ row }">
              <span class="px-3 py-1 rounded-full text-sm font-bold"
                    :class="row.score >= 8 ? 'bg-emerald-100 text-emerald-700' : row.score >= 7 ? 'bg-amber-100 text-amber-700' : 'bg-rose-100 text-rose-700'">
                {{ row.score }}/10
              </span>
            </template>
            <template #action="{ row }">
              <button @click="router.push(`/mock-results/${row.id}`)" class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-bold text-blue-600 hover:text-white border border-blue-200 hover:bg-blue-600 rounded-lg transition-colors group">
                <Eye class="w-4 h-4 group-hover:scale-110 transition-transform" /> Xem chi tiết
              </button>
            </template>
          </Table>
        </div>
      </div>
    </Card>
  </div>
</template>

<style>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}

/* Custom Table Styles targeting AppTable internal elements if possible */
.modern-table-container table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
}
.modern-table-container th {
  background-color: #f8fafc;
  color: #475569;
  font-weight: 700;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
}
.modern-table-container td {
  padding: 1rem 1.5rem;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
  font-weight: 500;
}
.modern-table-container tr:last-child td {
  border-bottom: none;
}
.modern-table-container tbody tr {
  transition: all 0.2s ease;
}
.modern-table-container tbody tr:hover {
  background-color: #f8fafc;
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
  position: relative;
  z-index: 10;
}
</style>
