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

const role = authStore.user?.role || localStorage.getItem('user_role') || 'candidate'
const activeTab = ref('account')
const toast = ref(null)
const showDeleteModal = ref(false)
const saving = ref(false)
const changingPassword = ref(false)
const passwordForm = ref({ current: '', next: '', confirm: '' })
const settings = ref({
  notify_email_interview: true,
  notify_email_reminder: true,
  notify_push: true,
  public_profile: true,
  share_anon_results: true
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

const handleChangePassword = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  if (!passwordForm.value.current || !passwordForm.value.next) {
    toast.value = { type: 'warning', message: 'Vui lòng nhập đầy đủ mật khẩu hiện tại và mật khẩu mới.' }
    return
  }
  if (passwordForm.value.next.length < 6) {
    toast.value = { type: 'warning', message: 'Mật khẩu mới phải có ít nhất 6 ký tự.' }
    return
  }
  if (passwordForm.value.next !== passwordForm.value.confirm) {
    toast.value = { type: 'warning', message: 'Xác nhận mật khẩu mới không khớp.' }
    return
  }
  changingPassword.value = true
  try {
    await authService.changePassword({
      current_password: passwordForm.value.current,
      new_password: passwordForm.value.next
    })
    passwordForm.value = { current: '', next: '', confirm: '' }
    toast.value = { type: 'success', message: 'Đổi mật khẩu thành công!' }
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi đổi mật khẩu: ' + (error.message || 'Không xác định') }
  } finally {
    changingPassword.value = false
  }
}

const handleDeleteAccount = () => {
  showDeleteModal.value = true
}

const confirmDeleteAccount = async () => {
  showDeleteModal.value = false
  try {
    await authService.deleteAccount()
    toast.value = { type: 'success', message: 'Tài khoản đã được xóa. Đang đăng xuất...' }
    setTimeout(() => { authStore.logout() }, 2000)
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi xóa tài khoản: ' + (error.message || 'Không xác định') }
  }
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
                
              </div>
            </form>

            <div style="border-top: 1px solid var(--border); padding-top: 24px; margin-top: 24px">
              <h3 class="text-body" style="font-weight: 600; margin-bottom: 16px; display: flex; align-items: center; gap: 8px">
                <Key size="16" /> Đổi mật khẩu
              </h3>
              <form @submit="handleChangePassword">
                <div style="display: flex; flex-direction: column; gap: 16px">
                  <Input label="Mật khẩu hiện tại" type="password" v-model="passwordForm.current" />
                  <Input label="Mật khẩu mới" type="password" v-model="passwordForm.next" />
                  <Input label="Xác nhận mật khẩu mới" type="password" v-model="passwordForm.confirm" />
                </div>
                <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                  <Button type="submit" :disabled="changingPassword"><Save size="16" /> {{ changingPassword ? 'Đang đổi...' : 'Đổi mật khẩu' }}</Button>
                </div>
              </form>
            </div>
          </Card>
        </div>

        <div v-if="activeTab === 'notifications'" class="space-y-6">
          <Card title="Thông báo & Cảnh báo">
            <div class="space-y-6">
              <div>
                <h3 class="text-base font-bold text-[var(--text-main)] uppercase tracking-wider mb-3 flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-[var(--primary)]"></span> Thông báo qua Email
                </h3>
                <div class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                  <div class="flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Lịch phỏng vấn mới</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Nhận email thông báo ngay khi có lịch phỏng vấn hoặc phản hồi từ nhà tuyển dụng.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_email_interview" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                    </label>
                  </div>

                  <div v-if="role === 'candidate'" class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Nhắc nhở phỏng vấn (Nhắc trước 1 giờ)</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Gửi email nhắc nhở trước khi buổi phỏng vấn (thật hoặc thử) diễn ra.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_email_reminder" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                    </label>
                  </div>
                </div>
              </div>

              <div>
                <h3 class="text-base font-bold text-[var(--text-main)] uppercase tracking-wider mb-3 flex items-center gap-2">
                  <span class="w-2 h-2 rounded-full bg-[var(--accent)]"></span> Thông báo đẩy (Push Notifications)
                </h3>
                <div class="bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                  <div class="flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Thông báo thời gian thực (Browser push)</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Hiển thị thông báo ngay trên góc trình duyệt khi có cập nhật hoặc phản hồi mới từ AI.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_push" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                    </label>
                  </div>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button @click="handleSave" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu tùy chọn' }}</Button>
              </div>
            </div>
          </Card>
        </div>

        <div v-if="activeTab === 'privacy'" class="space-y-6">
          <Card title="Quyền riêng tư & Bảo mật">
            <div class="space-y-6">
              <div v-if="role === 'candidate'" class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                <div class="flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Hồ sơ công khai (Public Profile)</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Cho phép các nhà tuyển dụng trên hệ thống AI Interview tìm thấy hồ sơ của bạn.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.public_profile" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                  </label>
                </div>

                <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Chia sẻ dữ liệu ẩn danh cải thiện AI</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Cho phép sử dụng các phiên phỏng vấn thử ẩn danh để huấn luyện và nâng cao độ chính xác của AI.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.share_anon_results" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--primary)] shadow-inner"></div>
                  </label>
                </div>
              </div>

              <div class="border-t border-[var(--border)] pt-5">
                <h3 class="text-base font-semibold text-[var(--danger)] mb-3">Quản lý dữ liệu</h3>
                <div class="flex gap-3">
                  <Button variant="ghost" style="color: var(--danger); border-color: var(--danger)" @click="handleDeleteAccount">Xóa tài khoản</Button>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button @click="handleSave" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu tùy chọn' }}</Button>
              </div>
            </div>
          </Card>
        </div>

        <div v-if="activeTab === 'workspace' && role === 'recruiter'" class="space-y-6">
          <Card title="Cấu hình AI Workspace">
            <div class="space-y-6">
              <p class="text-sm text-[var(--text-secondary)]">Tuỳ chỉnh cách AI Assistant hoạt động trong không gian làm việc của công ty bạn.</p>

              <div class="space-y-2">
                <label class="text-base font-semibold text-[var(--text-main)] block">Mô hình AI mặc định</label>
                <select class="input-field w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm font-semibold text-[var(--text-main)]" defaultValue="gpt-4">
                  <option value="gpt-4">GPT-4 (Độ chính xác cao nhất)</option>
                  <option value="gpt-35">GPT-3.5 Turbo (Nhanh nhất)</option>
                  <option value="claude">Claude 3 Haiku (Tối ưu hội thoại)</option>
                </select>
              </div>

              <div class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                <div class="flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Tự động bóc tách JD thành Rubric</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Phân tích yêu cầu tuyển dụng để tạo ra bộ tiêu chí đánh giá chuẩn xác.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" defaultChecked class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--accent)] shadow-inner"></div>
                  </label>
                </div>

                <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Tự động Suggest câu hỏi follow-up trong lúc phỏng vấn</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Gợi ý câu hỏi tiếp theo ngay khi ứng viên trả lời xong để đào sâu kỹ năng.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" defaultChecked class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[var(--accent)] shadow-inner"></div>
                  </label>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button @click="handleSave"><Save size="16" /> Lưu cấu hình</Button>
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
