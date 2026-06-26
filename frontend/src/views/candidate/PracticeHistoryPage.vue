<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import { Target, TrendingUp, Eye } from 'lucide-vue-next'

const router = useRouter()
const history = ref([
  { id: 1, role: 'Frontend Developer', level: 'Middle', date: '20/10/2023', score: 8.5 },
  { id: 2, role: 'Frontend Developer', level: 'Junior', date: '15/10/2023', score: 7.0 },
  { id: 3, role: 'React Developer', level: 'Junior', date: '05/10/2023', score: 6.5 },
])

const columns = [
  { header: 'Vị trí luyện tập', key: 'role' },
  { header: 'Cấp độ', key: 'level' },
  { header: 'Ngày thực hiện', key: 'date' },
  { header: 'Điểm số', key: 'score' },
  { header: 'Hành động', key: 'action' }
]
</script>

<template>
  <div>
    <div style="margin-bottom: 32px">
      <h1 class="text-h1">Kết quả & Lịch sử luyện tập</h1>
      <p class="text-helper" style="margin-top: 4px">Theo dõi sự tiến bộ của bạn qua các bài luyện tập với AI.</p>
    </div>

    <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; margin-bottom: 32px">
      <Card style="padding: 24px">
        <div style="display: flex; align-items: center; gap: 16px">
          <div style="width: 48px; height: 48px; border-radius: 12px; background-color: rgba(22, 163, 74, 0.1); display: flex; align-items: center; justify-content: center">
            <Target color="var(--success)" />
          </div>
          <div>
            <p class="text-helper">Điểm trung bình</p>
            <h2 class="text-h1">7.3<span style="font-size: 16px; color: var(--text-muted); font-weight: 400">/10</span></h2>
          </div>
        </div>
      </Card>
      
      <Card style="padding: 24px">
        <div style="display: flex; align-items: center; gap: 16px">
          <div style="width: 48px; height: 48px; border-radius: 12px; background-color: rgba(37, 99, 235, 0.1); display: flex; align-items: center; justify-content: center">
            <TrendingUp color="var(--primary)" />
          </div>
          <div>
            <p class="text-helper">Số lần luyện tập</p>
            <h2 class="text-h1">3<span style="font-size: 16px; color: var(--text-muted); font-weight: 400"> lần</span></h2>
          </div>
        </div>
      </Card>

      <Card style="padding: 24px; background-color: var(--surface-soft)">
        <p class="text-helper" style="margin-bottom: 8px; font-weight: 600">Điểm cần khắc phục</p>
        <ul class="text-body" style="padding-left: 20px; color: var(--danger)">
          <li>Trình bày cấu trúc câu trả lời (STAR)</li>
          <li>Đưa ra ví dụ thực tế</li>
        </ul>
      </Card>
    </div>

    <Card title="Lịch sử bài luyện tập">
      <Table :columns="columns" :data="history">
        <template #score="{ row }">
          <span :style="{ fontWeight: 'bold', color: row.score >= 8 ? 'var(--success)' : row.score >= 7 ? 'var(--warning)' : 'var(--danger)' }">
            {{ row.score }}/10
          </span>
        </template>
        <template #action="{ row }">
          <Button variant="ghost" style="padding: 4px 8px" @click="router.push(`/mock-results/${row.id}`)">
            <Eye size="16" /> Xem
          </Button>
        </template>
      </Table>
    </Card>
  </div>
</template>
