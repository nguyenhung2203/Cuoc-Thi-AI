<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import Modal from '../../components/common/AppModal.vue'
import { authStore } from '../../stores/auth.store'
import { authService } from '../../services/auth.service'
import { Save, Key, Bell, Shield, User, Monitor } from 'lucide-vue-next'

const role = localStorage.getItem('role') || 'recruiter'
const activeTab = ref('account')
const toast = ref(null)
const showDeleteModal = ref(false)
const saving = ref(false)
const settings = ref({
  notify_email_interview: true,
  notify_email_report: true,
  notify_push: true,
  encrypt_recordings: true,
  require_2fa: false,
  ai_model: 'gpt-4',
  auto_rubric: true,
  auto_suggest_followup: true
})

onMounted(async () => {
  try {
    const res = await authService.getSettings()
    if (res?.settings) {
      settings.value = { ...settings.value, ...res.settings }
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
  }
})

const handleSave = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  saving.value = true
  try {
    await authService.saveSettings(settings.value)
    toast.value = { type: 'success', message: 'Cài đặt đã được lưu thành công!' }
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi lưu cài đặt: ' + (error.message || 'Không xác định') }
  } finally {
    saving.value = false
  }
}

const handleDeleteAccount = () => {
  showDeleteModal.value = true
}

const confirmDeleteAccount = () => {
  showDeleteModal.value = false
  toast.value = { type: 'success', message: 'Yêu cầu xóa tài khoản đã được gửi.' }
}
</script>

<template>
  <div style="max-width: 800px; margin: 0 auto">
    <div style="margin-bottom: 32px">
      <h1 class="text-h1">Cài đặt Hệ thống</h1>
      <p class="text-helper" style="margin-top: 4px">Tùy chỉnh thông tin tài khoản và cấu hình hệ thống.</p>
    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <div style="display: grid; grid-template-columns: 240px 1fr; gap: 32px">
      
      <!-- Settings Sidebar Menu -->
      <div style="display: flex; flex-direction: column; gap: 8px">
        <div 
          :class="['nav-item', activeTab === 'account' ? 'active' : '']" 
          :style="{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'account' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'account' ? 'var(--primary)' : 'var(--text-secondary)' }"
          @click="activeTab = 'account'"
        >
          <User size="16" style="margin-right: 8px" /> Tài khoản
        </div>
        <div 
          :class="['nav-item', activeTab === 'notifications' ? 'active' : '']" 
          :style="{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'notifications' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'notifications' ? 'var(--primary)' : 'var(--text-secondary)' }"
          @click="activeTab = 'notifications'"
        >
          <Bell size="16" style="margin-right: 8px" /> Thông báo
        </div>
        <div 
          :class="['nav-item', activeTab === 'privacy' ? 'active' : '']" 
          :style="{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'privacy' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'privacy' ? 'var(--primary)' : 'var(--text-secondary)' }"
          @click="activeTab = 'privacy'"
        >
          <Shield size="16" style="margin-right: 8px" /> Quyền riêng tư & Bảo mật
        </div>
        <div v-if="role === 'recruiter'"
          :class="['nav-item', activeTab === 'workspace' ? 'active' : '']" 
          :style="{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'workspace' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'workspace' ? 'var(--primary)' : 'var(--text-secondary)' }"
          @click="activeTab = 'workspace'"
        >
          <Monitor size="16" style="margin-right: 8px" /> Cấu hình AI Workspace
        </div>
      </div>

      <!-- Settings Content -->
      <div style="display: flex; flex-direction: column; gap: 24px">
        
        <div v-if="activeTab === 'account'">
          <Card title="Thông tin tài khoản">
            <form @submit="handleSave">
              <div style="display: flex; flex-direction: column; gap: 16px">
                <Input label="Tên người dùng" :modelValue="authStore.user?.full_name || ''" disabled />
                <Input label="Email đăng nhập" type="email" :modelValue="authStore.user?.email || ''" disabled />
                
                <div style="border-top: 1px solid var(--border); padding-top: 24px; margin-top: 8px">
                  <h3 class="text-body" style="font-weight: 600; margin-bottom: 16px; display: flex; align-items: center; gap: 8px">
                    <Key size="16" /> Đổi mật khẩu
                  </h3>
                  <div style="display: flex; flex-direction: column; gap: 16px">
                    <Input label="Mật khẩu hiện tại" type="password" />
                    <Input label="Mật khẩu mới" type="password" />
                    <Input label="Xác nhận mật khẩu mới" type="password" />
                  </div>
                </div>

                <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                  <Button><Save size="16" /> Lưu thay đổi</Button>
                </div>
              </div>
            </form>
          </Card>
        </div>

        <div v-if="activeTab === 'notifications'">
          <Card title="Thông báo & Cảnh báo">
            <div style="display: flex; flex-direction: column; gap: 20px">
              <div>
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px">Qua Email</h3>
                <div style="display: flex; flex-direction: column; gap: 12px">
                  <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                    <input type="checkbox" v-model="settings.notify_email_interview" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                    <span class="text-body">Nhận email thông báo khi có lịch phỏng vấn mới</span>
                  </label>
                  <label v-if="role === 'recruiter'" style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                    <input type="checkbox" v-model="settings.notify_email_report" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                    <span class="text-body">Nhận email khi AI Report đã xử lý xong</span>
                  </label>
                </div>
              </div>

              <div style="border-top: 1px solid var(--border); padding-top: 20px">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px">Thông báo đẩy (Push Notifications)</h3>
                <div style="display: flex; flex-direction: column; gap: 12px">
                  <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                    <input type="checkbox" v-model="settings.notify_push" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                    <span class="text-body">Hiển thị thông báo trên trình duyệt (Browser push)</span>
                  </label>
                </div>
              </div>

              <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                <Button @click="handleSave" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu tùy chọn' }}</Button>
              </div>
            </div>
          </Card>
        </div>

        <div v-if="activeTab === 'privacy'">
          <Card title="Quyền riêng tư & Bảo mật">
            <div style="display: flex; flex-direction: column; gap: 20px">

              <div v-if="role === 'recruiter'">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px">Bảo mật dữ liệu công ty</h3>
                <div style="display: flex; flex-direction: column; gap: 12px">
                  <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                    <input type="checkbox" v-model="settings.encrypt_recordings" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                    <span class="text-body">Mã hóa ghi âm/video các cuộc phỏng vấn (E2E Encryption)</span>
                  </label>
                  <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                    <input type="checkbox" v-model="settings.require_2fa" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                    <span class="text-body">Yêu cầu xác thực 2 bước (2FA) khi đăng nhập nội bộ</span>
                  </label>
                </div>
              </div>

              <div style="border-top: 1px solid var(--border); padding-top: 20px">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px; color: var(--danger)">Quản lý dữ liệu</h3>
                <div style="display: flex; gap: 12px">
                  <Button variant="ghost" style="color: var(--danger); border-color: var(--danger)" @click="handleDeleteAccount">Xóa tài khoản</Button>
                </div>
              </div>

              <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                <Button @click="handleSave" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu tùy chọn' }}</Button>
              </div>
            </div>
          </Card>
        </div>

        <div v-if="activeTab === 'workspace' && role === 'recruiter'">
          <Card title="Cấu hình AI Workspace">
            <div style="display: flex; flex-direction: column; gap: 16px">
              <p class="text-body" style="color: var(--text-secondary); margin-bottom: 8px">Tuỳ chỉnh cách AI Assistant hoạt động trong không gian làm việc của công ty bạn.</p>

              <div class="input-group">
                <label class="input-label">Mô hình AI mặc định</label>
                <select class="input-field" v-model="settings.ai_model">
                  <option value="gpt-4">GPT-4 (Độ chính xác cao nhất)</option>
                  <option value="gpt-35">GPT-3.5 Turbo (Nhanh nhất)</option>
                  <option value="claude">Claude 3 Haiku (Tối ưu hội thoại)</option>
                </select>
              </div>

              <div style="display: flex; flex-direction: column; gap: 12px; margin-top: 12px">
                <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                  <input type="checkbox" v-model="settings.auto_rubric" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                  <span class="text-body">Tự động bóc tách JD thành Rubric</span>
                </label>
                <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                  <input type="checkbox" v-model="settings.auto_suggest_followup" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                  <span class="text-body">Tự động Suggest câu hỏi follow-up trong lúc phỏng vấn</span>
                </label>
              </div>

              <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                <Button @click="handleSave" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu cấu hình' }}</Button>
              </div>
            </div>
          </Card>
        </div>

      </div>

    </div>

    <Modal :isOpen="showDeleteModal" @close="showDeleteModal = false" title="Xác nhận xóa tài khoản">
      <p class="text-body" style="margin-bottom: 24px">Bạn có chắc chắn muốn xóa tài khoản này? Hành động này không thể hoàn tác và toàn bộ dữ liệu của bạn sẽ bị xóa vĩnh viễn khỏi hệ thống.</p>
      <div style="display: flex; justify-content: flex-end; gap: 12px">
        <Button variant="ghost" @click="showDeleteModal = false">Huỷ</Button>
        <Button variant="primary" style="background-color: var(--danger); color: white; border-color: var(--danger)" @click="confirmDeleteAccount">Xác nhận xóa</Button>
      </div>
    </Modal>
  </div>
</template>
