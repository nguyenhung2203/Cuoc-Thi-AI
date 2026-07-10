<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Modal from '../../components/common/AppModal.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Plus, Search, Filter, Bot, Copy, Edit, Trash2, Tag, Clock, MoreHorizontal } from 'lucide-vue-next'
import { templateService } from '../../services/template.service'
import { authStore } from '../../stores/auth.store'
import { computed } from 'vue'

const templates = ref([])
const searchKeyword = ref('')
const filterType = ref('')
const filteredTemplates = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase()
  return templates.value.filter(t => {
    const matchKw = !kw || (t.title || '').toLowerCase().includes(kw) || (t.desc || '').toLowerCase().includes(kw)
    const matchType = !filterType.value || t.type === filterType.value
    return matchKw && matchType
  })
})
const loading = ref(true)
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const showEditModal = ref(false)
const deletingId = ref(null)
const toast = ref(null)
const saving = ref(false)
const editingTemplate = ref(null)

const newTemplate = ref({
  name: '',
  type: 'Technical',
  description: '',
  duration_minutes: 60,
  config: {}
})

const mapTemplate = (t) => ({
  id: t.id,
  title: t.name,
  type: t.type,
  desc: t.description,
  duration: t.duration_minutes,
  tags: t.type ? [t.type] : [],
  created: t.created_at ? new Date(t.created_at).toLocaleDateString('vi-VN') : '—'
})

const fetchTemplates = async () => {
  loading.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    if (!companyId) {
      throw new Error('Không tìm thấy company ID')
    }
    const response = await templateService.getTemplates(companyId)
    templates.value = (Array.isArray(response) ? response : []).map(mapTemplate)
  } catch (error) {
    console.error('Lỗi tải mẫu:', error)
    toast.value = { type: 'error', message: 'Lỗi tải danh sách mẫu: ' + (error.message || 'Không xác định') }
    templates.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchTemplates)

const confirmDelete = async () => {
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await templateService.deleteTemplate(companyId, deletingId.value)
    templates.value = templates.value.filter(t => t.id !== deletingId.value)
    showDeleteModal.value = false
    toast.value = { type: 'success', message: 'Đã xóa mẫu AI thành công!' }
  } catch (error) {
    showDeleteModal.value = false
    toast.value = { type: 'error', message: 'Lỗi xóa mẫu: ' + (error.message || 'Không xác định') }
  }
}

const handleCreate = async () => {
  if (!newTemplate.value.name) return
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    const created = await templateService.createTemplate(companyId, {
      name: newTemplate.value.name,
      type: newTemplate.value.type,
      description: newTemplate.value.description,
      duration_minutes: newTemplate.value.duration_minutes,
      config: newTemplate.value.config
    })
    newTemplate.value = { name: '', type: 'Technical', description: '', duration_minutes: 60, config: {} }
    showCreateModal.value = false
    toast.value = { type: 'success', message: 'Đã tạo mẫu AI mới!' }
    await fetchTemplates()
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi tạo mẫu: ' + (error.message || 'Không xác định') }
  } finally {
    saving.value = false
  }
}

const handleDuplicate = async (template) => {
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await templateService.createTemplate(companyId, {
      name: (template.title || 'Mẫu') + ' (bản sao)',
      type: template.type,
      description: template.desc,
      duration_minutes: template.duration,
      config: {}
    })
    toast.value = { type: 'success', message: 'Đã nhân bản mẫu AI!' }
    await fetchTemplates()
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi nhân bản: ' + (error.message || 'Không xác định') }
  } finally {
    saving.value = false
  }
}

const openEditModal = (template) => {
  editingTemplate.value = template
  newTemplate.value = {
    name: template.title,
    type: template.type,
    description: template.desc,
    duration_minutes: template.duration,
    config: {}
  }
  showEditModal.value = true
}

const handleUpdate = async () => {
  if (!newTemplate.value.name || !editingTemplate.value) return
  saving.value = true
  try {
    const companyId = authStore.user?.companies?.[0]?.id
    await templateService.updateTemplate(companyId, editingTemplate.value.id, {
      name: newTemplate.value.name,
      type: newTemplate.value.type,
      description: newTemplate.value.description,
      duration_minutes: newTemplate.value.duration_minutes,
      config: newTemplate.value.config
    })
    newTemplate.value = { name: '', type: 'Technical', description: '', duration_minutes: 60, config: {} }
    editingTemplate.value = null
    showEditModal.value = false
    toast.value = { type: 'success', message: 'Đã cập nhật mẫu AI!' }
    await fetchTemplates()
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi cập nhật mẫu: ' + (error.message || 'Không xác định') }
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div>
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Mẫu AI</h1>
        <p class="text-helper" style="margin-top: 4px">Quản lý các mẫu kịch bản phỏng vấn AI cho doanh nghiệp.</p>
      </div>
      <Button @click="showCreateModal = true"><Plus size="16" /> Tạo mẫu mới</Button>
    </div>

    <!-- Toolbar -->
    <div style="display: flex; gap: 16px; margin-bottom: 24px">
      <div style="position: relative; flex: 1; max-width: 400px">
        <Search size="16" style="position: absolute; left: 12px; top: 12px; color: var(--text-muted)" />
        <input
          type="text"
          placeholder="Tìm kiếm mẫu AI..."
          class="input-field"
          style="width: 100%; padding-left: 36px"
          v-model="searchKeyword"
        />
      </div>
      <select class="input-field" style="width: 180px" v-model="filterType">
        <option value="">Tất cả loại kịch bản</option>
        <option value="Technical">Kỹ năng chuyên môn</option>
        <option value="Behavioral">Kỹ năng mềm</option>
        <option value="Management">Quản lý</option>
      </select>
    </div>

    <!-- Loading state -->
    <div v-if="loading" style="text-align: center; padding: 60px">
      <p class="text-helper">Đang tải danh sách mẫu...</p>
    </div>

    <!-- Empty state -->
    <div v-else-if="templates.length === 0">
      <Card style="text-align: center; padding: 60px">
        <Bot size="48" style="margin: 0 auto 16px; opacity: 0.3" />
        <p class="text-body" style="font-weight: 600; margin-bottom: 8px">Chưa có mẫu AI nào</p>
        <p class="text-helper" style="margin-bottom: 24px; max-width: 400px; margin-left: auto; margin-right: auto">
          Tạo mẫu kịch bản phỏng vấn AI để sử dụng cho các việc làm.
        </p>
        <Button @click="showCreateModal = true"><Plus size="16" /> Tạo mẫu đầu tiên</Button>
      </Card>
    </div>

    <!-- Template Grid -->
    <div v-else class="grid" style="grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 24px">
      <Card v-for="template in filteredTemplates" :key="template.id" style="display: flex; flex-direction: column; height: 100%; transition: transform 0.2s, box-shadow 0.2s" class="hover-card">
        <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px">
          <div style="display: flex; gap: 12px; align-items: flex-start">
            <div style="width: 40px; height: 40px; border-radius: 8px; background-color: rgba(37, 99, 235, 0.1); color: var(--primary); display: flex; align-items: center; justify-content: center; flex-shrink: 0">
              <Bot size="20" />
            </div>
            <div>
              <h3 class="text-body" style="font-weight: 600; color: var(--text-main); margin-bottom: 4px; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; text-overflow: ellipsis;" :title="template.title">{{ template.title }}</h3>
              <span style="font-size: 12px; color: var(--text-muted)">{{ template.type }}</span>
            </div>
          </div>
          <div style="display: flex; gap: 4px">
            <Button variant="ghost" style="padding: 4px; color: var(--text-muted)" @click="handleDuplicate(template)"><Copy size="16" /></Button>
            <Button variant="ghost" style="padding: 4px; color: var(--text-muted)" @click="openEditModal(template)"><Edit size="16" /></Button>
            <Button variant="ghost" style="padding: 4px; color: var(--danger)" @click="deletingId = template.id; showDeleteModal = true"><Trash2 size="16" /></Button>
          </div>
        </div>

        <p class="text-helper" style="color: var(--text-secondary); margin-bottom: 24px; flex-grow: 1; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; text-overflow: ellipsis;">
          {{ template.desc }}
        </p>

        <div style="border-top: 1px solid var(--border); padding-top: 16px; margin-top: auto">
          <div style="display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 12px">
            <Badge v-for="tag in template.tags" :key="tag" type="neutral" style="font-size: 11px; padding: 2px 6px"><Tag size="10" style="margin-right: 4px"/>{{ tag }}</Badge>
          </div>
          <div style="display: flex; justify-content: space-between; align-items: center">
            <span class="text-helper" style="display: flex; align-items: center; gap: 4px; color: var(--text-muted)"><Clock size="12" /> Đã tạo: {{ template.created }}</span>
            <Button variant="ghost" style="font-size: 13px; color: var(--primary); padding: 0" @click="openEditModal(template)">Cấu hình &rarr;</Button>
          </div>
        </div>
      </Card>
    </div>

    <!-- Modals -->
    <Modal :isOpen="showCreateModal" @close="showCreateModal = false" title="Tạo mẫu AI mới">
      <Input label="Tên mẫu AI" v-model="newTemplate.name" placeholder="VD: Khối lập trình Frontend..." style="margin-bottom: 16px" required />

      <div class="input-group" style="margin-bottom: 16px">
        <label class="input-label">Loại kịch bản</label>
        <select class="input-field" v-model="newTemplate.type">
          <option value="Technical">Kỹ năng chuyên môn (Technical)</option>
          <option value="Behavioral">Kỹ năng mềm (Behavioral)</option>
          <option value="Management">Quản lý (Management)</option>
          <option value="Custom">Tùy chỉnh</option>
        </select>
      </div>

      <div class="input-group" style="margin-bottom: 16px">
        <label class="input-label">Thời lượng (phút)</label>
        <input type="number" v-model.number="newTemplate.duration_minutes" class="input-field" min="15" step="15" />
      </div>

      <div class="input-group" style="margin-bottom: 24px">
        <label class="input-label">Mô tả ngắn</label>
        <textarea class="input-field" rows="4" v-model="newTemplate.description" placeholder="Mô tả mục đích của bộ câu hỏi và kịch bản này..."></textarea>
      </div>

      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showCreateModal = false">Hủy</Button>
        <Button variant="primary" :disabled="saving" @click="handleCreate">{{ saving ? 'Đang lưu...' : 'Lưu mẫu AI' }}</Button>
      </div>
    </Modal>

    <Modal :isOpen="showEditModal" @close="showEditModal = false" title="Chỉnh sửa mẫu AI">
      <Input label="Tên mẫu AI" v-model="newTemplate.name" placeholder="VD: Khối lập trình Frontend..." style="margin-bottom: 16px" required />

      <div class="input-group" style="margin-bottom: 16px">
        <label class="input-label">Loại kịch bản</label>
        <select class="input-field" v-model="newTemplate.type">
          <option value="Technical">Kỹ năng chuyên môn (Technical)</option>
          <option value="Behavioral">Kỹ năng mềm (Behavioral)</option>
          <option value="Management">Quản lý (Management)</option>
          <option value="Custom">Tùy chỉnh</option>
        </select>
      </div>

      <div class="input-group" style="margin-bottom: 16px">
        <label class="input-label">Thời lượng (phút)</label>
        <input type="number" v-model.number="newTemplate.duration_minutes" class="input-field" min="15" step="15" />
      </div>

      <div class="input-group" style="margin-bottom: 24px">
        <label class="input-label">Mô tả ngắn</label>
        <textarea class="input-field" rows="4" v-model="newTemplate.description" placeholder="Mô tả mục đích của bộ câu hỏi và kịch bản này..."></textarea>
      </div>

      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showEditModal = false">Hủy</Button>
        <Button variant="primary" :disabled="saving" @click="handleUpdate">{{ saving ? 'Đang lưu...' : 'Cập nhật mẫu' }}</Button>
      </div>
    </Modal>

    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn xóa mẫu AI này? Các việc làm đang sử dụng mẫu này có thể phải cấu hình lại.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showDeleteModal = false">Hủy</Button>
        <Button variant="primary" style="background-color: var(--danger); border-color: var(--danger)" @click="confirmDelete">Xóa mẫu</Button>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.hover-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  border-color: rgba(37, 99, 235, 0.3);
}
</style>
