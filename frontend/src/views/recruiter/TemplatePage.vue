<script setup>
import { ref } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Modal from '../../components/common/AppModal.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import Badge from '../../components/common/AppBadge.vue'
import { Plus, Search, Filter, Bot, Copy, Edit, Trash2, Tag, Clock, MoreHorizontal } from 'lucide-vue-next'

const templates = ref([
  { id: 1, title: 'Frontend ReactJS Fresher', type: 'Technical', desc: 'Bộ câu hỏi tập trung vào nền tảng ReactJS, Javascript và DOM.', tags: ['React', 'Fresher', 'Frontend'], created: '2026-06-25' },
  { id: 2, title: 'Backend Node.js Middle', type: 'Technical', desc: 'Kiểm tra kiến thức System Design cơ bản và API Design.', tags: ['Nodejs', 'Middle', 'Backend'], created: '2026-06-20' },
  { id: 3, title: 'Project Manager (Agile)', type: 'Management', desc: 'Đánh giá khả năng quản lý rủi ro và giải quyết xung đột.', tags: ['Agile', 'Manager', 'Behavioral'], created: '2026-06-15' },
  { id: 4, title: 'Văn hóa Doanh nghiệp', type: 'Behavioral', desc: 'Bộ câu hỏi tiêu chuẩn để đánh giá mức độ phù hợp văn hóa.', tags: ['Culture', 'Soft Skills'], created: '2026-06-10' }
])

const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const deletingId = ref(null)
const toast = ref(null)
const newTemplate = ref({
  title: '',
  type: 'Technical',
  desc: ''
})

const confirmDelete = () => {
  templates.value = templates.value.filter(t => t.id !== deletingId.value)
  showDeleteModal.value = false
  toast.value = { type: 'success', message: 'Đã xóa mẫu AI thành công!' }
}

const handleCreate = () => {
  if (!newTemplate.value.title) return
  templates.value.unshift({
    id: Date.now(),
    title: newTemplate.value.title,
    type: newTemplate.value.type,
    desc: newTemplate.value.desc,
    tags: ['New'],
    created: new Date().toISOString().split('T')[0]
  })
  newTemplate.value = { title: '', type: 'Technical', desc: '' }
  showCreateModal.value = false
  toast.value = { type: 'success', message: 'Đã tạo mẫu AI mới!' }
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
        />
      </div>
      <select class="input-field" style="width: 180px">
        <option value="">Tất cả loại kịch bản</option>
        <option value="Technical">Kỹ năng chuyên môn</option>
        <option value="Behavioral">Kỹ năng mềm</option>
        <option value="Management">Quản lý</option>
      </select>
      <Button variant="secondary"><Filter size="16" /> Lọc nâng cao</Button>
    </div>

    <!-- Template Grid -->
    <div class="grid" style="grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 24px">
      <Card v-for="template in templates" :key="template.id" style="display: flex; flex-direction: column; height: 100%; transition: transform 0.2s, box-shadow 0.2s" class="hover-card">
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
            <Button variant="ghost" style="padding: 4px; color: var(--text-muted)" @click="toast = { type: 'info', message: 'Đã nhân bản mẫu AI!' }"><Copy size="16" /></Button>
            <Button variant="ghost" style="padding: 4px; color: var(--text-muted)" @click="toast = { type: 'info', message: 'Tính năng chỉnh sửa đang phát triển' }"><Edit size="16" /></Button>
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
            <Button variant="ghost" style="font-size: 13px; color: var(--primary); padding: 0">Cấu hình &rarr;</Button>
          </div>
        </div>
      </Card>
    </div>

    <!-- Modals -->
    <Modal :isOpen="showCreateModal" @close="showCreateModal = false" title="Tạo mẫu AI mới">
      <Input label="Tên mẫu AI" v-model="newTemplate.title" placeholder="VD: Khối lập trình Frontend..." style="margin-bottom: 16px" required />
      
      <div class="input-group" style="margin-bottom: 16px">
        <label class="input-label">Loại kịch bản</label>
        <select class="input-field" v-model="newTemplate.type">
          <option value="Technical">Kỹ năng chuyên môn (Technical)</option>
          <option value="Behavioral">Kỹ năng mềm (Behavioral)</option>
          <option value="Management">Quản lý (Management)</option>
          <option value="Custom">Tùy chỉnh</option>
        </select>
      </div>

      <div class="input-group" style="margin-bottom: 24px">
        <label class="input-label">Mô tả ngắn</label>
        <textarea class="input-field" rows="4" v-model="newTemplate.desc" placeholder="Mô tả mục đích của bộ câu hỏi và kịch bản này..."></textarea>
      </div>

      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showCreateModal = false">Hủy</Button>
        <Button variant="primary" @click="handleCreate">Lưu mẫu AI</Button>
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
