<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import Modal from '../../components/common/AppModal.vue'
import { authStore } from '../../stores/auth.store'
import { authService } from '../../services/auth.service'
import { companyService } from '../../services/company.service'
import { confirmPassword, requiredTrim, validateForm, validatePassword } from '../../utils/validators.js'
import { Settings, Save, Key, Bell, Shield, User, Monitor, Building2 } from 'lucide-vue-next'

const role = localStorage.getItem('role') || 'recruiter'
const route = useRoute()
const activeTab = ref(route.query.tab || 'account')

watch(() => route.query.tab, (newTab) => {
  if (newTab) activeTab.value = newTab
})
const toast = ref(null)
const showDeleteModal = ref(false)
const saving = ref(false)
const changingPassword = ref(false)
const passwordForm = ref({ current: '', next: '', confirm: '' })
const passwordErrors = ref({})
const settings = ref({
  notify_email_interview: true,
  notify_email_report: true,
  notify_push: true,
  encrypt_recordings: true,
  require_2fa: false,
  ai_model: 'gpt-4',
  auto_rubric: true,
  auto_rubric: true,
  auto_suggest_followup: true
})
const companyForm = ref({ id: '', name: '', industry: 'IT', size: '1-50', website: '' })

onMounted(async () => {
  try {
    const res = await authService.getSettings()
    if (res?.settings) {
      settings.value = { ...settings.value, ...res.settings }
    }
    // Load company info if recruiter
    if (role === 'recruiter' && authStore.user?.companies?.length > 0) {
      const c = authStore.user.companies[0]
      companyForm.value = { id: c.id, name: c.name || '', industry: c.industry || 'IT', size: c.size || '1-50', website: c.website || '' }
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

const handleSaveCompany = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  if (!companyForm.value.id || !companyForm.value.name) return
  saving.value = true
  try {
    await companyService.updateCompany(companyForm.value.id, {
      name: companyForm.value.name,
      industry: companyForm.value.industry,
      size: companyForm.value.size,
      website: companyForm.value.website
    })
    toast.value = { type: 'success', message: 'Hồ sơ công ty đã được cập nhật thành công!' }
    // Update local state
    if (authStore.user?.companies?.[0]) {
      authStore.user.companies[0].name = companyForm.value.name
      authStore.user.companies[0].industry = companyForm.value.industry
      authStore.user.companies[0].size = companyForm.value.size
      authStore.user.companies[0].website = companyForm.value.website
    }
  } catch (error) {
    toast.value = { type: 'error', message: 'Lỗi lưu hồ sơ công ty: ' + (error.message || 'Không xác định') }
  } finally {
    saving.value = false
  }
}

const handleChangePassword = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  if (changingPassword.value) return

  const validation = validateForm(passwordForm.value, {
    current: [(value) => requiredTrim(value, 'Vui lòng nhập mật khẩu hiện tại.')],
    next: [
      (value) => requiredTrim(value, 'Vui lòng nhập mật khẩu mới.'),
      validatePassword,
      (value, values) => value === values.current ? 'Mật khẩu mới phải khác mật khẩu hiện tại.' : '',
    ],
    confirm: [
      (value) => requiredTrim(value, 'Vui lòng xác nhận mật khẩu mới.'),
      (value, values) => confirmPassword(value, values.next),
    ],
  })
  passwordErrors.value = validation.errors
  if (!validation.isValid) return

  changingPassword.value = true
  try {
    await authService.changePassword({
      current_password: passwordForm.value.current,
      new_password: passwordForm.value.next
    })
    passwordForm.value = { current: '', next: '', confirm: '' }
    passwordErrors.value = {}
    toast.value = { type: 'success', message: 'Đổi mật khẩu thành công!' }
  } catch (error) {
    const message = error?.message || 'Không thể đổi mật khẩu. Vui lòng thử lại.'
    toast.value = { type: 'error', message }
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
  <div style="max-width: 1000px; margin: 0 auto; padding-bottom: 48px">
    <!-- Page Header -->
    <Card class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-6 rounded-2xl shadow-sm mb-6">
      <div>
        <h1 class="text-2xl font-bold text-slate-800 dark:text-white flex items-center gap-2.5">
          <Settings size="26" class="text-blue-600 dark:text-blue-400" />
          Cài đặt & Cấu hình Hệ thống
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm mt-1">
          Tùy chỉnh thông tin tài khoản cá nhân, hồ sơ công ty và cấu hình nền tảng AI Workspace.
        </p>
      </div>
    </Card>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />

    <!-- Tabs Navigation -->
    <div class="flex items-center gap-2 border-b border-[var(--border)] pb-2 mb-6" style="overflow-x: auto; white-space: nowrap;">
      <button 
        @click="activeTab = 'account'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :style="activeTab === 'account' ? 'background: var(--primary); color: #fff; box-shadow: 0 4px 12px rgba(37,99,235,0.2)' : 'color: var(--text-secondary); background: transparent'"
        :class="activeTab !== 'account' ? 'hover:bg-[var(--surface-soft)]' : ''"
      >
        <User size="18" />
        <span>Tài khoản</span>
      </button>

      <button v-if="role === 'recruiter'"
        @click="activeTab = 'company'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :style="activeTab === 'company' ? 'background: var(--primary); color: #fff; box-shadow: 0 4px 12px rgba(37,99,235,0.2)' : 'color: var(--text-secondary); background: transparent'"
        :class="activeTab !== 'company' ? 'hover:bg-[var(--surface-soft)]' : ''"
      >
        <Building2 size="18" />
        <span>Hồ sơ công ty</span>
      </button>

      <button 
        @click="activeTab = 'notifications'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :style="activeTab === 'notifications' ? 'background: var(--primary); color: #fff; box-shadow: 0 4px 12px rgba(37,99,235,0.2)' : 'color: var(--text-secondary); background: transparent'"
        :class="activeTab !== 'notifications' ? 'hover:bg-[var(--surface-soft)]' : ''"
      >
        <Bell size="18" />
        <span>Thông báo</span>
      </button>

      <button 
        @click="activeTab = 'privacy'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :style="activeTab === 'privacy' ? 'background: var(--primary); color: #fff; box-shadow: 0 4px 12px rgba(37,99,235,0.2)' : 'color: var(--text-secondary); background: transparent'"
        :class="activeTab !== 'privacy' ? 'hover:bg-[var(--surface-soft)]' : ''"
      >
        <Shield size="18" />
        <span>Quyền riêng tư & Bảo mật</span>
      </button>

      <button v-if="role === 'recruiter'"
        @click="activeTab = 'workspace'"
        class="flex items-center gap-2 px-5 py-2.5 rounded-xl font-bold text-sm transition-all"
        :style="activeTab === 'workspace' ? 'background: var(--primary); color: #fff; box-shadow: 0 4px 12px rgba(37,99,235,0.2)' : 'color: var(--text-secondary); background: transparent'"
        :class="activeTab !== 'workspace' ? 'hover:bg-[var(--surface-soft)]' : ''"
      >
        <Monitor size="18" />
        <span>Cấu hình AI Workspace</span>
      </button>
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
                  <Input label="Mật khẩu hiện tại" type="password" v-model="passwordForm.current" :error="passwordErrors.current" required />
                  <Input label="Mật khẩu mới" type="password" v-model="passwordForm.next" :error="passwordErrors.next" placeholder="Ít nhất 8 ký tự, có chữ hoa, chữ thường và số" required />
                  <Input label="Xác nhận mật khẩu mới" type="password" v-model="passwordForm.confirm" :error="passwordErrors.confirm" required />
                </div>
                <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                  <Button type="submit" :disabled="changingPassword"><Save size="16" /> {{ changingPassword ? 'Đang đổi...' : 'Đổi mật khẩu' }}</Button>
                </div>
              </form>
            </div>
          </Card>
        </div>


        <div v-if="activeTab === 'company' && role === 'recruiter'">
          <Card title="Hồ sơ Công ty">
            <form @submit="handleSaveCompany">
              <div style="display: flex; flex-direction: column; gap: 16px">
                <Input label="Tên công ty / Doanh nghiệp" v-model="companyForm.name" required placeholder="VD: FPT Software, VNG..." />
                
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px">
                  <div class="input-group">
                    <label class="input-label" style="display: block; font-size: 13px; font-weight: 600; color: var(--text-main); margin-bottom: 6px">Quy mô nhân sự</label>
                    <select class="input-field" v-model="companyForm.size" style="width: 100%; height: 40px; padding: 0 12px; border-radius: var(--radius); border: 1px solid var(--border); background: #fff; font-size: 14px; color: var(--text-main); outline: none; transition: border-color 0.2s">
                      <option value="1-50">1 - 50 nhân viên</option>
                      <option value="51-200">51 - 200 nhân viên</option>
                      <option value="201-1000">201 - 1000 nhân viên</option>
                      <option value="1000+">Hơn 1000 nhân viên</option>
                    </select>
                  </div>
                  
                  <div class="input-group">
                    <label class="input-label" style="display: block; font-size: 13px; font-weight: 600; color: var(--text-main); margin-bottom: 6px">Lĩnh vực hoạt động</label>
                    <select class="input-field" v-model="companyForm.industry" style="width: 100%; height: 40px; padding: 0 12px; border-radius: var(--radius); border: 1px solid var(--border); background: #fff; font-size: 14px; color: var(--text-main); outline: none; transition: border-color 0.2s">
                      <option value="IT">Công nghệ thông tin (IT)</option>
                      <option value="Finance">Tài chính / Ngân hàng</option>
                      <option value="Education">Giáo dục / Đào tạo</option>
                      <option value="Healthcare">Y tế / Chăm sóc sức khỏe</option>
                      <option value="Manufacturing">Sản xuất</option>
                      <option value="Retail">Bán lẻ / Thương mại điện tử</option>
                      <option value="Other">Khác</option>
                    </select>
                  </div>
                </div>

                <Input label="Website công ty" type="url" v-model="companyForm.website" placeholder="https://..." />
              </div>
              <div style="display: flex; justify-content: flex-end; margin-top: 24px">
                <Button type="submit" variant="primary" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu hồ sơ công ty' }}</Button>
              </div>
            </form>
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
                      <span class="text-base font-semibold text-[var(--text-main)] block">Nhận email thông báo khi có lịch phỏng vấn mới</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Hệ thống gửi email ngay khi có ứng viên xác nhận hoặc đặt lịch phỏng vấn với doanh nghiệp.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_email_interview" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
                    </label>
                  </div>

                  <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                    <div class="space-y-1 pr-4">
                      <span class="text-base font-semibold text-[var(--text-main)] block">Nhận email khi AI Report đã xử lý xong</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Gửi báo cáo phân tích và điểm số đánh giá chi tiết của ứng viên từ mô hình AI vào email HR.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_email_report" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
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
                      <span class="text-base font-semibold text-[var(--text-main)] block">Hiển thị thông báo trên trình duyệt (Browser push)</span>
                      <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Nhận thông báo realtime tức thì trên màn hình làm việc khi ứng viên tham gia phòng phỏng vấn.</span>
                    </div>
                    <label class="relative inline-flex items-center cursor-pointer shrink-0">
                      <input type="checkbox" v-model="settings.notify_push" class="sr-only peer" />
                      <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
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
              <div class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                <div class="flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Mã hóa ghi âm/video các cuộc phỏng vấn (E2E Encryption)</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Bảo vệ toàn vẹn dữ liệu cuộc họp bằng cơ chế mã hóa đầu cuối tiêu chuẩn Enterprise.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.encrypt_recordings" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
                  </label>
                </div>

                <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Yêu cầu xác thực 2 bước (2FA) khi đăng nhập nội bộ</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Đảm bảo an toàn tài khoản tuyển dụng bằng mã xác thực qua ứng dụng authenticator.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.require_2fa" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
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
                <select class="input-field w-full px-4 py-3 bg-[var(--surface)] border border-[var(--border)] rounded-xl text-sm font-semibold text-[var(--text-main)]" v-model="settings.ai_model">
                  <option value="gpt-4">GPT-4 (Độ chính xác cao nhất)</option>
                  <option value="gpt-35">GPT-3.5 Turbo (Nhanh nhất)</option>
                  <option value="claude">Claude 3 Haiku (Tối ưu hội thoại)</option>
                </select>
              </div>

              <div class="space-y-4 bg-[var(--surface-soft)] p-6 rounded-2xl border border-[var(--border)] shadow-sm">
                <div class="flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Tự động bóc tách JD thành Rubric</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Phân tích yêu cầu tuyển dụng để tạo ra bộ tiêu chí đánh giá chuẩn xác cho từng vị trí.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.auto_rubric" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
                  </label>
                </div>

                <div class="border-t border-[var(--border)] pt-4 flex items-center justify-between gap-6">
                  <div class="space-y-1 pr-4">
                    <span class="text-base font-semibold text-[var(--text-main)] block">Tự động Suggest câu hỏi follow-up trong lúc phỏng vấn</span>
                    <span class="text-sm text-[var(--text-secondary)] block leading-relaxed">Đưa ra gợi ý câu hỏi tiếp theo ngay khi ứng viên vừa trả lời xong để đào sâu kỹ năng.</span>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer shrink-0">
                    <input type="checkbox" v-model="settings.auto_suggest_followup" class="sr-only peer" />
                    <div class="w-12 h-6 bg-slate-300 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600 shadow-inner"></div>
                  </label>
                </div>
              </div>

              <div class="flex justify-end pt-2">
                <Button @click="handleSave" :disabled="saving"><Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu cấu hình' }}</Button>
              </div>
            </div>
          </Card>
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
