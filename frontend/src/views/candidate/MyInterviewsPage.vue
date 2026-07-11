<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { candidatePortalService } from '../../services/candidate-portal.service'
import Button from '../../components/common/AppButton.vue'
import { CalendarDays, Building2, Clock, Video } from 'lucide-vue-next'
import { langStore } from '../../stores/lang.store'

const router = useRouter()
const interviews = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await candidatePortalService.getInterviews()
    interviews.value = data
  } catch (err) {
    console.error('Lỗi tải danh sách phỏng vấn:', err)
  } finally {
    loading.value = false
  }
})

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit' }) + ', ' + d.toLocaleDateString('vi-VN')
}
</script>

<template>
  <div class="space-y-8 pb-12 max-w-6xl mx-auto px-4 md:px-0">
    <div class="header-box animate-rise flex flex-col sm:flex-row items-start sm:items-center justify-between p-6 rounded-2xl bg-[var(--surface)] border border-[var(--border)] shadow-sm gap-4">
      <div>
        <h1 class="text-h1 mb-1.5 text-[var(--text-main)]">{{ langStore.t('interviews', 'title') }}</h1>
        <p class="text-secondary text-sm">{{ langStore.t('interviews', 'subtitle') }}</p>
      </div>
      <div class="shrink-0">
        <Button variant="primary" @click="router.push('/mock-setup')" class="sheen">
          <Video :size="16" class="mr-1.5 shrink-0" />
          <span>Luyện tập với AI ngay</span>
        </Button>
      </div>
    </div>

    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="mi-spinner mb-4"></div>
      <p class="text-helper">{{ langStore.t('interviews', 'loadingText') }}</p>
    </div>

    <div v-else-if="interviews.length === 0" class="mi-empty">
      <div class="mi-empty-icon"><CalendarDays :size="34" /></div>
      <h2 class="mi-empty-title">{{ langStore.t('interviews', 'emptyTitle') }}</h2>
      <p class="mi-empty-desc">
        {{ langStore.t('interviews', 'emptyDesc') }}
      </p>
    </div>

    <div v-else class="space-y-4 stagger">
      <div v-for="iv in interviews" :key="iv.id" class="mi-item card-elevate hover-rail">
        <div class="flex-1 min-w-0">
          <h3 class="mi-title">{{ iv.job_title || iv.title }}</h3>
          <p class="mi-company">
            <Building2 :size="15" />
            {{ iv.company_name || 'Công ty ẩn danh' }}
          </p>

          <div class="flex flex-wrap gap-2 items-center">
            <span class="badge badge-info gap-1.5">
              <Clock :size="14" /> {{ formatDate(iv.scheduled_at) }}
            </span>
            <span class="badge" :class="iv.mode === 'real' ? 'badge-danger' : 'badge-neutral'">
              {{ iv.mode === 'real' ? langStore.t('interviews', 'realMode') : langStore.t('interviews', 'mockMode') }}
            </span>
            <span v-if="iv.status === 'completed'" class="badge badge-success">{{ langStore.t('interviews', 'completed') }}</span>
            <span v-else-if="iv.status === 'cancelled'" class="badge badge-warning">{{ langStore.t('interviews', 'cancelled') }}</span>
          </div>
        </div>

        <div class="flex gap-3 w-full md:w-auto shrink-0">
          <Button
            v-if="iv.mode === 'real' && iv.status !== 'completed' && iv.status !== 'cancelled'"
            variant="primary"
            class="sheen w-full md:w-auto"
            @click="iv.join_link ? router.push(iv.join_link) : null">
            <Video :size="16" /> {{ langStore.t('interviews', 'joinBtn') }}
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mi-spinner { width: 44px; height: 44px; border-radius: 50%; border: 3px solid var(--primary-light); border-top-color: var(--primary); animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.mi-empty { display: flex; flex-direction: column; align-items: center; text-align: center; padding: 56px 24px; background: var(--surface-soft); border: 2px dashed var(--border); border-radius: var(--radius-lg); }
.mi-empty-icon { display: flex; align-items: center; justify-content: center; width: 76px; height: 76px; border-radius: var(--radius-lg); background: var(--primary-light); color: var(--primary); margin-bottom: 22px; }
.mi-empty-title { font-size: 20px; font-weight: 700; color: var(--text-main); margin-bottom: 10px; }
.mi-empty-desc { color: var(--text-secondary); max-width: 30rem; line-height: 1.6; }

.mi-item { display: flex; flex-direction: column; gap: 20px; padding: 22px; padding-left: 26px; background: var(--surface); border: 1px solid var(--border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
@media (min-width: 768px) { .mi-item { flex-direction: row; align-items: center; justify-content: space-between; } }
.mi-title { font-size: 18px; font-weight: 700; color: var(--text-main); margin-bottom: 4px; transition: color 0.2s ease; }
.mi-item:hover .mi-title { color: var(--primary); }
.mi-company { display: flex; align-items: center; gap: 6px; color: var(--text-secondary); font-weight: 500; margin-bottom: 14px; font-size: 14px; }
</style>
