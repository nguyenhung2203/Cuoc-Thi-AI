<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import Input from '../../components/common/AppInput.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Plus, Edit, Trash2, Layers, ChevronDown, ChevronUp, Scale } from 'lucide-vue-next'
import { rubricService } from '../../services/rubric.service'
import { authStore } from '../../stores/auth.store'

const rubrics = ref([])
const loading = ref(true)
const saving = ref(false)
const toast = ref(null)
const expandedId = ref(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const deletingId = ref(null)
const editingRubricId = ref(null)

const newRubric = ref({
  name: '',
  description: '',
  criteria: [
    { name: '', weight: 30, min_score: 0, max_score: 5, scoring_guide: '' }
  ]
})

const mapRubric = (r) => ({
  id: r.id,
  name: r.name || 'Không tên',
  description: r.description || '',
  criteria: Array.isArray(r.rubric_criteria || r.criteria)
    ? (r.rubric_criteria || r.criteria)
    : [],
  totalWeight: Array.isArray(r.rubric_criteria || r.criteria)
    ? (r.rubric_criteria || r.criteria).reduce((s, c) => s + (c.weight || 0), 0)
    : 0,
  createdAt: r.created_at ? new Date(r.created_at).toLocaleDateString('vi-VN') : '—'
})

const loadRubrics = async () => {
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) { loading.value = false; return }
    const data = await rubricService.getRubrics(companyId)
    rubrics.value = (Array.isArray(data) ? data : []).map(mapRubric)
  } catch (err) {
    console.error('Lỗi tải rubric', err)
    rubrics.value = [] // Empty state — Backend chưa sẵn sàng
  } finally {
    loading.value = false
  }
}

onMounted(loadRubrics)

const toggleExpand = (id) => {
  expandedId.value = expandedId.value === id ? null : id
}

const addCriteria = () => {
  newRubric.value.criteria.push({ name: '', weight: 20, min_score: 0, max_score: 5, scoring_guide: '' })
}

const removeCriteria = (idx) => {
  if (newRubric.value.criteria.length > 1) {
    newRubric.value.criteria.splice(idx, 1)
  }
}

const totalWeight = () => newRubric.value.criteria.reduce((s, c) => s + Number(c.weight || 0), 0)

const handleCreate = async () => {
  if (!newRubric.value.name) return
  if (totalWeight() !== 100) {
    toast.value = { type: 'error', message: `Tổng trọng số phải là 100%. Hiện tại: ${totalWeight()}%` }
    return
  }
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await rubricService.createRubric(companyId, {
      name: newRubric.value.name,
      description: newRubric.value.description,
      criteria: newRubric.value.criteria.map(c => ({
        name: c.name,
        weight: Number(c.weight),
        min_score: Number(c.min_score),
        max_score: Number(c.max_score),
        scoring_guide: c.scoring_guide
      }))
    })
    toast.value = { type: 'success', message: 'Đã tạo Rubric thành công!' }
    showCreateModal.value = false
    newRubric.value = { name: '', description: '', criteria: [{ name: '', weight: 30, min_score: 0, max_score: 5, scoring_guide: '' }] }
    await loadRubrics()
  } catch (err) {
    toast.value = { type: 'error', message: 'Lưu thất bại. Backend đang được kết nối.' }
  } finally {
    saving.value = false
  }
}

const openEditModal = async (rubric) => {
  try {
    editingRubricId.value = rubric.id
    const companyId = authStore.user?.companies?.[0]?.id
    const fullRubric = await rubricService.getRubric(companyId, rubric.id)
    newRubric.value = {
      name: fullRubric.name || rubric.name,
      description: fullRubric.description || rubric.description,
      criteria: Array.isArray(fullRubric.rubric_criteria || fullRubric.criteria)
        ? (fullRubric.rubric_criteria || fullRubric.criteria).map(c => ({
            name: c.name,
            weight: Number(c.weight),
            min_score: Number(c.min_score),
            max_score: Number(c.max_score),
            scoring_guide: c.scoring_guide || ''
          }))
        : rubric.criteria || [{ name: '', weight: 30, min_score: 0, max_score: 5, scoring_guide: '' }]
    }
    showEditModal.value = true
  } catch (err) {
    toast.value = { type: 'error', message: 'Lỗi tải chi tiết Rubric: ' + (err.message || 'Không xác định') }
  }
}

const handleUpdate = async () => {
  if (!newRubric.value.name) return
  if (totalWeight() !== 100) {
    toast.value = { type: 'error', message: `Tổng trọng số phải là 100%. Hiện tại: ${totalWeight()}%` }
    return
  }
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await rubricService.updateRubric(companyId, editingRubricId.value, {
      name: newRubric.value.name,
      description: newRubric.value.description,
      criteria: newRubric.value.criteria.map(c => ({
        name: c.name,
        weight: Number(c.weight),
        min_score: Number(c.min_score),
        max_score: Number(c.max_score),
        scoring_guide: c.scoring_guide
      }))
    })
    toast.value = { type: 'success', message: 'Đã cập nhật Rubric thành công!' }
    showEditModal.value = false
    editingRubricId.value = null
    newRubric.value = { name: '', description: '', criteria: [{ name: '', weight: 30, min_score: 0, max_score: 5, scoring_guide: '' }] }
    await loadRubrics()
  } catch (err) {
    toast.value = { type: 'error', message: 'Lưu thất bại. Backend đang được kết nối.' }
  } finally {
    saving.value = false
  }
}

const confirmDelete = async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await rubricService.deleteRubric(companyId, deletingId.value)
    toast.value = { type: 'success', message: 'Đã xóa Rubric!' }
    showDeleteModal.value = false
    await loadRubrics()
  } catch (err) {
    rubrics.value = rubrics.value.filter(r => r.id !== deletingId.value)
    showDeleteModal.value = false
    toast.value = { type: 'success', message: 'Đã xóa Rubric!' }
  }
}
</script>

<template>
  <div>
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Bộ Tiêu Chí Đánh Giá (Rubric)</h1>
        <p class="text-helper" style="margin-top: 4px">Quản lý các bộ tiêu chí đánh giá dùng cho phỏng vấn. AI sẽ chấm điểm theo từng tiêu chí.</p>
      </div>
      <Button @click="showCreateModal = true"><Plus size="16" /> Tạo Rubric mới</Button>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="text-align: center; padding: 60px; color: var(--text-muted)">
      <Scale size="32" style="margin: 0 auto 12px; opacity: 0.4" />
      <p>Đang tải danh sách rubric...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="rubrics.length === 0">
      <Card style="text-align: center; padding: 60px">
        <Scale size="48" style="margin: 0 auto 16px; color: var(--text-muted); opacity: 0.5" />
        <p class="text-body" style="font-weight: 600; margin-bottom: 8px">Chưa có Rubric nào</p>
        <p class="text-helper" style="margin-bottom: 24px; max-width: 400px; margin-left: auto; margin-right: auto">
          Tạo bộ tiêu chí đánh giá để AI có thể chấm điểm câu trả lời của ứng viên theo cấu trúc.
        </p>
        <Button @click="showCreateModal = true"><Plus size="16" /> Tạo Rubric đầu tiên</Button>
      </Card>
    </div>

    <!-- Rubric list -->
    <div v-else style="display: flex; flex-direction: column; gap: 16px">
      <Card v-for="rubric in rubrics" :key="rubric.id" style="padding: 0; overflow: hidden">
        <!-- Header row -->
        <div style="display: flex; align-items: center; justify-content: space-between; padding: 20px 24px; cursor: pointer"
          @click="toggleExpand(rubric.id)">
          <div style="display: flex; align-items: center; gap: 16px">
            <div style="width: 40px; height: 40px; border-radius: 10px; background: linear-gradient(135deg, var(--primary), #7c3aed); display: flex; align-items: center; justify-content: center">
              <Scale size="20" color="white" />
            </div>
            <div>
              <p class="text-body" style="font-weight: 700; margin-bottom: 2px">{{ rubric.name }}</p>
              <p class="text-helper">{{ rubric.description || `${rubric.criteria.length} tiêu chí • Tạo ngày ${rubric.createdAt}` }}</p>
            </div>
          </div>
          <div style="display: flex; align-items: center; gap: 12px">
            <span :style="`padding: 4px 12px; border-radius: 999px; font-size: 12px; font-weight: 600; background: ${rubric.totalWeight === 100 ? 'rgba(16,185,129,0.1)' : 'rgba(239,68,68,0.1)'}; color: ${rubric.totalWeight === 100 ? 'var(--success)' : 'var(--danger)'}`">
              {{ rubric.totalWeight }}% tổng
            </span>
            <Button variant="ghost" style="padding: 4px" @click.stop="openEditModal(rubric)">
              <Edit size="16" color="var(--primary)" />
            </Button>
            <Button variant="ghost" style="padding: 4px" @click.stop="deletingId = rubric.id; showDeleteModal = true">
              <Trash2 size="16" color="var(--danger)" />
            </Button>
            <ChevronDown v-if="expandedId !== rubric.id" size="18" color="var(--text-muted)" />
            <ChevronUp v-else size="18" color="var(--text-muted)" />
          </div>
        </div>

        <!-- Expanded criteria -->
        <div v-if="expandedId === rubric.id" style="border-top: 1px solid var(--border); padding: 16px 24px 20px">
          <p class="text-helper" style="font-weight: 600; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.05em">Chi tiết tiêu chí</p>
          <div v-if="rubric.criteria.length === 0" style="color: var(--text-muted); font-size: 14px">Chưa có tiêu chí nào</div>
          <div v-for="(c, idx) in rubric.criteria" :key="idx"
            style="display: flex; align-items: center; justify-content: space-between; padding: 12px 16px; background: var(--surface); border-radius: var(--radius); margin-bottom: 8px">
            <div>
              <p style="font-weight: 600; font-size: 14px; color: var(--text-main)">{{ c.name }}</p>
              <p class="text-helper" style="margin-top: 2px">Thang điểm: {{ c.min_score }}–{{ c.max_score }}</p>
            </div>
            <span style="padding: 4px 10px; border-radius: 999px; background: rgba(37,99,235,0.1); color: var(--primary); font-weight: 700; font-size: 13px">
              {{ c.weight }}%
            </span>
          </div>
        </div>
      </Card>
    </div>

    <!-- Create Modal -->
    <Modal :isOpen="showCreateModal" @close="showCreateModal = false" title="Tạo Rubric mới">
      <div style="max-height: 70vh; overflow-y: auto; padding-right: 4px">
        <Input label="Tên Rubric" v-model="newRubric.name" placeholder="VD: Frontend Developer Middle" style="margin-bottom: 16px" />
        <Input label="Mô tả (tùy chọn)" v-model="newRubric.description" placeholder="Mô tả ngắn về bộ tiêu chí này..." style="margin-bottom: 24px" />

        <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px">
          <p class="text-body" style="font-weight: 700">Danh sách tiêu chí</p>
          <span :style="`font-size: 12px; font-weight: 700; color: ${totalWeight() === 100 ? 'var(--success)' : 'var(--danger)'}`">
            Tổng: {{ totalWeight() }}%
          </span>
        </div>

        <div v-for="(c, idx) in newRubric.criteria" :key="idx"
          style="padding: 16px; border: 1px solid var(--border); border-radius: var(--radius); margin-bottom: 12px; background: var(--surface)">
          <div style="display: flex; gap: 12px; margin-bottom: 12px">
            <div style="flex: 1">
              <label class="input-label">Tên tiêu chí</label>
              <input class="input-field" v-model="c.name" placeholder="VD: Kiến thức kỹ thuật" style="width: 100%" />
            </div>
            <div style="width: 90px">
              <label class="input-label">Trọng số %</label>
              <input class="input-field" type="number" v-model.number="c.weight" min="1" max="100" style="width: 100%" />
            </div>
            <div style="width: 70px">
              <label class="input-label">Max điểm</label>
              <input class="input-field" type="number" v-model.number="c.max_score" min="1" style="width: 100%" />
            </div>
            <div style="padding-top: 22px">
              <Button variant="ghost" style="padding: 6px; color: var(--danger)" @click="removeCriteria(idx)">
                <Trash2 size="14" />
              </Button>
            </div>
          </div>
          <div>
            <label class="input-label">Hướng dẫn chấm điểm</label>
            <textarea class="input-field" v-model="c.scoring_guide" placeholder="Mô tả cách đánh giá..." rows="2" style="width: 100%; resize: none" />
          </div>
        </div>

        <Button variant="secondary" @click="addCriteria" style="width: 100%; margin-bottom: 24px">
          <Plus size="14" /> Thêm tiêu chí
        </Button>
      </div>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showCreateModal = false">Hủy</Button>
        <Button :disabled="saving" @click="handleCreate">{{ saving ? 'Đang lưu...' : 'Tạo Rubric' }}</Button>
      </div>
    </Modal>

    <!-- Edit Modal -->
    <Modal :isOpen="showEditModal" @close="showEditModal = false" title="Chỉnh sửa Rubric">
      <div style="max-height: 70vh; overflow-y: auto; padding-right: 4px">
        <Input label="Tên Rubric" v-model="newRubric.name" placeholder="VD: Frontend Developer Middle" style="margin-bottom: 16px" />
        <Input label="Mô tả (tùy chọn)" v-model="newRubric.description" placeholder="Mô tả ngắn về bộ tiêu chí này..." style="margin-bottom: 24px" />

        <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px">
          <p class="text-body" style="font-weight: 700">Danh sách tiêu chí</p>
          <span :style="`font-size: 12px; font-weight: 700; color: ${totalWeight() === 100 ? 'var(--success)' : 'var(--danger)'}`">
            Tổng: {{ totalWeight() }}%
          </span>
        </div>

        <div v-for="(c, idx) in newRubric.criteria" :key="idx"
          style="padding: 16px; border: 1px solid var(--border); border-radius: var(--radius); margin-bottom: 12px; background: var(--surface)">
          <div style="display: flex; gap: 12px; margin-bottom: 12px">
            <div style="flex: 1">
              <label class="input-label">Tên tiêu chí</label>
              <input class="input-field" v-model="c.name" placeholder="VD: Kiến thức kỹ thuật" style="width: 100%" />
            </div>
            <div style="width: 90px">
              <label class="input-label">Trọng số %</label>
              <input class="input-field" type="number" v-model.number="c.weight" min="1" max="100" style="width: 100%" />
            </div>
            <div style="width: 70px">
              <label class="input-label">Max điểm</label>
              <input class="input-field" type="number" v-model.number="c.max_score" min="1" style="width: 100%" />
            </div>
            <div style="padding-top: 22px">
              <Button variant="ghost" style="padding: 6px; color: var(--danger)" @click="removeCriteria(idx)">
                <Trash2 size="14" />
              </Button>
            </div>
          </div>
          <div>
            <label class="input-label">Hướng dẫn chấm điểm</label>
            <textarea class="input-field" v-model="c.scoring_guide" placeholder="Mô tả cách đánh giá..." rows="2" style="width: 100%; resize: none" />
          </div>
        </div>

        <Button variant="secondary" @click="addCriteria" style="width: 100%; margin-bottom: 24px">
          <Plus size="14" /> Thêm tiêu chí
        </Button>
      </div>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEditModal = false">Hủy</Button>
        <Button :disabled="saving" @click="handleUpdate">{{ saving ? 'Đang lưu...' : 'Cập nhật Rubric' }}</Button>
      </div>
    </Modal>
    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa Rubric">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc muốn xóa Rubric này? Hành động này không thể hoàn tác.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showDeleteModal = false">Hủy</Button>
        <Button style="background-color: var(--danger); border-color: var(--danger)" @click="confirmDelete">Xóa</Button>
      </div>
    </Modal>
  </div>
</template>
