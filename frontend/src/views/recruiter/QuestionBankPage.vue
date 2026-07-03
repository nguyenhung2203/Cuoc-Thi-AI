<script setup>
import { ref } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Table from '../../components/common/AppTable.vue'
import Modal from '../../components/common/AppModal.vue'
import Toast from '../../components/common/AppToast.vue'
import Input from '../../components/common/AppInput.vue'
import { Plus, Search, Filter, Layers, Edit, Trash2 } from 'lucide-vue-next'

const questions = ref([
  { id: 1, text: 'Hãy giải thích sự khác biệt giữa var, let và const trong JavaScript.', role: 'Frontend', level: 'Fresher', type: 'Technical' },
  { id: 2, text: 'Bạn đã bao giờ bất đồng quan điểm với Quản lý dự án chưa? Bạn giải quyết thế nào?', role: 'All', level: 'Middle', type: 'Behavioral' },
  { id: 3, text: 'Mô tả nguyên lý hoạt động của Virtual DOM trong React.', role: 'Frontend', level: 'Junior', type: 'Technical' },
  { id: 4, text: 'Làm thế nào để scale một hệ thống chịu tải 1 triệu requests/s?', role: 'Backend', level: 'Senior', type: 'System Design' },
])

const columns = [
  { header: 'Câu hỏi', key: 'text' },
  { header: 'Vị trí (Role)', key: 'role' },
  { header: 'Cấp độ', key: 'level' },
  { header: 'Loại', key: 'type' },
  { header: 'Hành động', key: 'action' }
]

const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const showFilterModal = ref(false)
const deletingId = ref(null)
const toast = ref(null)
const newQuestion = ref({
  text: '',
  role: 'All',
  level: 'Fresher',
  type: 'Technical'
})

const confirmDelete = () => {
  questions.value = questions.value.filter(q => q.id !== deletingId.value)
  showDeleteModal.value = false
  toast.value = { type: 'success', message: 'Đã xóa câu hỏi khỏi kho!' }
}

const handleCreate = () => {
  if (!newQuestion.value.text) return
  questions.value.unshift({
    id: Date.now(),
    text: newQuestion.value.text,
    role: newQuestion.value.role,
    level: newQuestion.value.level,
    type: newQuestion.value.type
  })
  newQuestion.value = {
    text: '',
    role: 'All',
    level: 'Fresher',
    type: 'Technical'
  }
  showCreateModal.value = false
  toast.value = { type: 'success', message: 'Đã thêm câu hỏi mới!' }
}
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px">
      <div>
        <h1 class="text-h1">Kho câu hỏi</h1>
        <p class="text-helper" style="margin-top: 4px">Quản lý ngân hàng câu hỏi dùng chung cho các buổi phỏng vấn.</p>
      </div>
      <Button @click="showCreateModal = true"><Plus size="16" /> Thêm câu hỏi</Button>
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
          />
        </div>
        <select class="input-field" style="width: 180px">
          <option value="">Tất cả vị trí (Role)</option>
          <option value="frontend">Frontend</option>
          <option value="backend">Backend</option>
        </select>
        <select class="input-field" style="width: 180px">
          <option value="">Tất cả cấp độ</option>
          <option value="fresher">Fresher</option>
          <option value="senior">Senior</option>
        </select>
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
            <Button variant="ghost" style="padding: 4px" @click="toast = { type: 'info', message: 'Tính năng chỉnh sửa đang phát triển' }"><Edit size="16" /></Button>
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

    <Modal :isOpen="showCreateModal" @close="showCreateModal = false" title="Thêm câu hỏi mới">
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
        <Button variant="primary" @click="handleCreate">Lưu câu hỏi</Button>
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
          <select class="input-field">
            <option value="">Tất cả</option>
            <option value="Technical">Technical</option>
            <option value="Behavioral">Behavioral</option>
            <option value="System Design">System Design</option>
          </select>
        </div>
        <div class="input-group">
          <label class="input-label">Nguồn câu hỏi</label>
          <select class="input-field">
            <option value="">Tất cả nguồn</option>
            <option value="System">Từ hệ thống AI</option>
            <option value="Custom">Tự tạo</option>
          </select>
        </div>
        <div class="input-group">
          <label class="input-label">Từ khóa/Tags</label>
          <Input placeholder="Nhập tags, cách nhau bởi dấu phẩy..." />
        </div>
      </div>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showFilterModal = false">Xóa bộ lọc</Button>
        <Button variant="primary" @click="showFilterModal = false; toast = { type: 'success', message: 'Đã áp dụng bộ lọc nâng cao!' }">Áp dụng</Button>
      </div>
    </Modal>
  </div>
</template>
