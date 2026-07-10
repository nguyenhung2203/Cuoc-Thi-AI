<script setup>
import { ref, onMounted } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import Modal from '../../components/common/AppModal.vue'
import { authStore } from '../../stores/auth.store'
import { authService } from '../../services/auth.service'
import { jobService } from '../../services/job.service'
import { candidateService } from '../../services/candidate.service'
import { interviewService } from '../../services/interview.service'
import { Save, Key, Bell, Shield, User, Monitor } from 'lucide-vue-next'
import * as XLSX from 'xlsx'

const role = localStorage.getItem('role') || 'recruiter'
const activeTab = ref('account')
const toast = ref(null)
const showDeleteModal = ref(false)

const notificationSettings = ref({
  browserPush: false
})
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)

onMounted(() => {
  const saved = localStorage.getItem('recruiterNotificationSettings')
  if (saved) {
    notificationSettings.value = JSON.parse(saved)
  }
})

const handleSave = (e) => {
  if (e && e.preventDefault) e.preventDefault()
  if (activeTab.value === 'notifications') {
    localStorage.setItem('recruiterNotificationSettings', JSON.stringify(notificationSettings.value))
    if (notificationSettings.value.browserPush && Notification.permission !== 'granted') {
      Notification.requestPermission()
    }
  }
  toast.value = { type: 'success', message: 'Các thay đổi đã được lưu thành công!' }
}

const handleChangePassword = async (e) => {
  if (e && e.preventDefault) e.preventDefault()
  
  if (!oldPassword.value || !newPassword.value || !confirmPassword.value) {
    toast.value = { type: 'error', message: 'Vui lòng điền đầy đủ các trường mật khẩu!' }
    return
  }
  
  if (newPassword.value !== confirmPassword.value) {
    toast.value = { type: 'error', message: 'Mật khẩu mới và xác nhận không khớp!' }
    return
  }

  if (newPassword.value.length < 6) {
    toast.value = { type: 'error', message: 'Mật khẩu mới phải dài ít nhất 6 ký tự!' }
    return
  }

  loading.value = true
  try {
    await authService.changePassword(oldPassword.value, newPassword.value)
    toast.value = { type: 'success', message: 'Đổi mật khẩu thành công!' }
    oldPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
  } catch (error) {
    console.error(error)
    toast.value = { type: 'error', message: error.response?.data?.message || 'Không thể đổi mật khẩu, vui lòng kiểm tra lại mật khẩu cũ.' }
  } finally {
    loading.value = false
  }
}

const handleDeleteAccount = () => {
  showDeleteModal.value = true
}

const confirmDeleteAccount = () => {
  showDeleteModal.value = false
  toast.value = { type: 'success', message: 'Yêu cầu xóa tài khoản đã được gửi.' }
}

const exportToExcel = async () => {
  toast.value = { type: 'info', message: 'Đang trích xuất dữ liệu, vui lòng đợi...' }
  try {
    const wb = XLSX.utils.book_new()
    
    // 1. Account info
    const userInfo = [
      { 
        "Họ Tên": authStore.user?.full_name || '', 
        "Email": authStore.user?.email || '',
        "Quyền": authStore.user?.role || '',
        "Công ty": authStore.user?.companies?.[0]?.name || '' 
      }
    ]
    const wsAccount = XLSX.utils.json_to_sheet(userInfo)
    XLSX.utils.book_append_sheet(wb, wsAccount, "Tài khoản")

    // 2. Fetch jobs
    const companyId = authStore.user?.companies?.[0]?.id
    if (companyId) {
      const jobsRes = await jobService.getJobs(companyId)
      const jobs = jobsRes.data || jobsRes || []
      if (jobs.length > 0) {
        const wsJobs = XLSX.utils.json_to_sheet(jobs.map(j => ({
          "ID": j.id,
          "Tiêu đề": j.title,
          "Phòng ban": j.department || '',
          "Địa điểm": j.location || '',
          "Mức lương": j.salary_range || '',
          "Trạng thái": j.status,
          "Ngày tạo": new Date(j.created_at).toLocaleDateString('vi-VN')
        })))
        XLSX.utils.book_append_sheet(wb, wsJobs, "Việc làm")
      }
    }

    // 3. Fetch Candidates
    if (companyId) {
      const candidatesRes = await candidateService.getCandidates(companyId)
      const candidates = candidatesRes.data || candidatesRes || []
      if (candidates.length > 0) {
        const wsCandidates = XLSX.utils.json_to_sheet(candidates.map(c => ({
          "Họ tên": c.full_name,
          "Email": c.email,
          "SĐT": c.phone || '',
          "Vị trí ứng tuyển": c.job_title || c.job_id || '',
          "Trạng thái": c.status || 'Mới',
          "Nguồn": c.source || '',
          "Ngày ứng tuyển": new Date(c.created_at).toLocaleDateString('vi-VN')
        })))
        XLSX.utils.book_append_sheet(wb, wsCandidates, "Ứng viên")
      }
    }

    // 4. Fetch Interviews
    if (companyId) {
      const interviewsRes = await interviewService.getInterviews(companyId)
      const interviews = interviewsRes.data || interviewsRes || []
      if (interviews.length > 0) {
        const wsInterviews = XLSX.utils.json_to_sheet(interviews.map(i => ({
          "Ứng viên": i.candidate_name || i.candidate_id || '',
          "Vị trí": i.job_title || i.job_id || '',
          "Lịch hẹn": new Date(i.scheduled_at).toLocaleString('vi-VN'),
          "Thời lượng (phút)": i.duration_minutes || 60,
          "Hình thức": i.mode || 'Online',
          "Trạng thái": i.status || 'Chờ phỏng vấn'
        })))
        XLSX.utils.book_append_sheet(wb, wsInterviews, "Lịch phỏng vấn")
      }
    }

    // Save file
    XLSX.writeFile(wb, `Data_Export_${new Date().toISOString().slice(0,10)}.xlsx`)
    toast.value = { type: 'success', message: 'Xuất dữ liệu Excel thành công!' }
  } catch (error) {
    console.error('Export error:', error)
    toast.value = { type: 'error', message: 'Có lỗi xảy ra khi xuất dữ liệu' }
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
            <form @submit="handleChangePassword">
              <div style="display: flex; flex-direction: column; gap: 16px">
                <Input label="Tên người dùng" :modelValue="authStore.user?.full_name || ''" disabled />
                <Input label="Email đăng nhập" type="email" :modelValue="authStore.user?.email || ''" disabled />
                
                <div style="border-top: 1px solid var(--border); padding-top: 24px; margin-top: 8px">
                  <h3 class="text-body" style="font-weight: 600; margin-bottom: 16px; display: flex; align-items: center; gap: 8px">
                    <Key size="16" /> Đổi mật khẩu
                  </h3>
                  <div style="display: flex; flex-direction: column; gap: 16px">
                    <Input label="Mật khẩu hiện tại" type="password" v-model="oldPassword" />
                    <Input label="Mật khẩu mới" type="password" v-model="newPassword" />
                    <Input label="Xác nhận mật khẩu mới" type="password" v-model="confirmPassword" />
                  </div>
                </div>

                <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                  <Button type="submit" :disabled="loading"><Save size="16" /> {{ loading ? 'Đang lưu...' : 'Lưu thay đổi' }}</Button>
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
                  <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Nhận email thông báo khi có lịch phỏng vấn mới <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                  <label v-if="role === 'recruiter'" style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Nhận email khi AI Report đã xử lý xong <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                  <label v-if="role === 'candidate'" style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Nhận email nhắc nhở trước 1 tiếng khi diễn ra phỏng vấn <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                </div>
              </div>

              <div style="border-top: 1px solid var(--border); padding-top: 20px">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px">Thông báo đẩy (Push Notifications)</h3>
                <div style="display: flex; flex-direction: column; gap: 12px">
                  <label style="display: flex; align-items: center; gap: 12px; cursor: pointer">
                    <input type="checkbox" v-model="notificationSettings.browserPush" style="width: 16px; height: 16px; accent-color: var(--primary)" />
                    <span class="text-body">Hiển thị thông báo trên trình duyệt (Browser push)</span>
                  </label>
                </div>
              </div>

              <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                <Button @click="handleSave"><Save size="16" /> Lưu tùy chọn</Button>
              </div>
            </div>
          </Card>
        </div>

        <div v-if="activeTab === 'privacy'">
          <Card title="Quyền riêng tư & Bảo mật">
            <div style="display: flex; flex-direction: column; gap: 20px">
              
              <div v-if="role === 'candidate'">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px">Hiển thị hồ sơ</h3>
                <div style="display: flex; flex-direction: column; gap: 12px">
                  <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Cho phép các nhà tuyển dụng khác xem hồ sơ của tôi (Public Profile) <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                  <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Chia sẻ ẩn danh kết quả Mock Interview để cải thiện AI <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                </div>
              </div>

              <div v-if="role === 'recruiter'">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px">Bảo mật dữ liệu công ty</h3>
                <div style="display: flex; flex-direction: column; gap: 12px">
                  <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Mã hóa ghi âm/video các cuộc phỏng vấn (E2E Encryption) <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                  <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                    <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                    <span class="text-body flex items-center gap-2">Yêu cầu xác thực 2 bước (2FA) khi đăng nhập nội bộ <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                  </label>
                </div>
              </div>

              <div style="border-top: 1px solid var(--border); padding-top: 20px">
                <h3 class="text-body" style="font-weight: 600; margin-bottom: 12px; color: var(--danger)">Quản lý dữ liệu</h3>
                <div style="display: flex; gap: 12px">
                  <Button variant="secondary" @click="exportToExcel">Xuất toàn bộ dữ liệu (Export Data)</Button>
                  <Button variant="ghost" style="color: var(--danger); border-color: var(--danger)" @click="handleDeleteAccount">Xóa tài khoản</Button>
                </div>
              </div>

              <div style="display: flex; justify-content: flex-end; margin-top: 16px">
                <Button @click="handleSave"><Save size="16" /> Lưu tùy chọn</Button>
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
                <select class="input-field" defaultValue="gpt-4">
                  <option value="gpt-4">GPT-4 (Độ chính xác cao nhất)</option>
                  <option value="gpt-35">GPT-3.5 Turbo (Nhanh nhất)</option>
                  <option value="claude">Claude 3 Haiku (Tối ưu hội thoại)</option>
                </select>
              </div>

              <div style="display: flex; flex-direction: column; gap: 12px; margin-top: 12px">
                <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                  <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                  <span class="text-body flex items-center gap-2">Tự động bóc tách JD thành Rubric <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                </label>
                <label style="display: flex; align-items: center; gap: 12px; cursor: not-allowed; opacity: 0.6">
                  <input type="checkbox" disabled style="width: 16px; height: 16px;" />
                  <span class="text-body flex items-center gap-2">Tự động Suggest câu hỏi follow-up trong lúc phỏng vấn <span class="badge badge-neutral" style="font-size: 10px; padding: 2px 6px">Đang phát triển</span></span>
                </label>
              </div>

              <div style="display: flex; justify-content: flex-end; margin-top: 16px">
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
