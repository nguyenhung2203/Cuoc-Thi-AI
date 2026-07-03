<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { publicService } from '../../services/public.service'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Briefcase, MapPin, Clock } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const companyId = route.params.company_id
const jobs = ref([])
const loading = ref(true)

const unwrap = (val) => {
  if (!val) return ''
  if (typeof val === 'object' && 'String' in val) {
    return val.Valid ? val.String : ''
  }
  return val
}

onMounted(async () => {
  try {
    const res = await publicService.getCompanyJobs(companyId)
    const rawJobs = Array.isArray(res) ? res : (res?.data || [])
    jobs.value = rawJobs.map(j => ({
      ...j,
      location: unwrap(j.location),
      employment_type: unwrap(j.employment_type),
      department: unwrap(j.department)
    }))
  } catch (error) {
    console.error('Failed to load jobs', error)
  } finally {
    loading.value = false
  }
})

const viewJob = (jobId) => {
  router.push(`/careers/${companyId}/jobs/${jobId}`)
}
</script>

<template>
  <div class="career-page">
    <div style="text-align: center; margin-bottom: 40px;">
      <h1 style="font-size: 32px; color: var(--text-main); margin-bottom: 12px;">Cơ hội nghề nghiệp</h1>
      <p style="color: var(--text-secondary); font-size: 16px;">Khám phá các vị trí đang tuyển dụng và tham gia cùng chúng tôi.</p>
    </div>

    <div v-if="loading" style="text-align: center; padding: 40px;">
      <div class="spinner" style="margin: 0 auto;"></div>
    </div>
    
    <div v-else-if="jobs.length === 0" style="text-align: center; padding: 60px; background: var(--surface); border-radius: var(--radius-lg); border: 1px solid var(--border);">
      <Briefcase size="48" color="var(--border)" style="margin-bottom: 16px;" />
      <h3 style="font-size: 18px; margin-bottom: 8px;">Chưa có vị trí nào đang mở</h3>
      <p style="color: var(--text-secondary);">Vui lòng quay lại sau nhé.</p>
    </div>

    <div v-else class="job-grid">
      <Card v-for="job in jobs" :key="job.id" class="job-card" hover style="display: flex; flex-direction: column;">
        <div style="flex: 1;">
          <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px;">
            <h3 style="font-size: 18px; margin: 0; color: var(--primary);">{{ job.title }}</h3>
            <Badge variant="primary" style="font-size: 12px;">{{ job.department || 'General' }}</Badge>
          </div>
          
          <div style="display: flex; gap: 16px; margin-bottom: 16px; color: var(--text-secondary); font-size: 14px;">
            <span style="display: flex; align-items: center; gap: 4px;">
              <MapPin size="14" /> {{ job.location || 'Bất kỳ' }}
            </span>
            <span style="display: flex; align-items: center; gap: 4px;">
              <Clock size="14" /> {{ job.employment_type || 'Full-time' }}
            </span>
          </div>
          
          <p style="color: var(--text-main); font-size: 14px; line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden;">
            {{ job.description }}
          </p>
        </div>
        
        <div style="margin-top: 24px;">
          <Button variant="primary" style="width: 100%" @click="viewJob(job.id)">Xem chi tiết</Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.job-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}
.job-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.job-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
}
</style>
