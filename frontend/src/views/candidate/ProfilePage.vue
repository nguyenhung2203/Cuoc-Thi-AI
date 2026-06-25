<script setup>
import { ref } from 'vue'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Input from '../../components/common/AppInput.vue'
import Toast from '../../components/common/AppToast.vue'
import { Upload, FileText, CheckCircle, Save } from 'lucide-vue-next'

const profile = ref({
  name: 'Candidate User',
  email: 'candidate@test.com',
  phone: '0123456789',
  linkedin: 'linkedin.com/in/candidate',
  targetRole: 'Frontend Developer',
  level: 'Middle',
})

const saving = ref(false)
const toast = ref(null)

const handleSave = (e) => {
  e.preventDefault()
  saving.value = true
  setTimeout(() => {
    saving.value = false
    toast.value = { type: 'success', message: 'Hồ sơ cá nhân đã được lưu thành công! Dữ liệu này sẽ được đồng bộ với AI.' }
  }, 800)
}
</script>

<template>
  <div style="max-width: 1000px; margin: 0 auto">
    <div style="margin-bottom: 32px">
      <h1 class="text-h1">Hồ sơ cá nhân & CV</h1>
      <p class="text-helper" style="margin-top: 4px">Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.</p>
    </div>

    <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px">
      
      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card title="Thông tin cơ bản">
          <form @submit="handleSave">
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px">
              <Input label="Họ và Tên" v-model="profile.name" required />
              <Input label="Email" type="email" v-model="profile.email" disabled />
              <Input label="Số điện thoại" v-model="profile.phone" />
              <Input label="LinkedIn Profile" v-model="profile.linkedin" />
            </div>

            <h3 class="text-h2" style="margin-top: 32px; margin-bottom: 16px">Định hướng nghề nghiệp</h3>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px">
              <Input label="Vị trí mục tiêu (Target Role)" v-model="profile.targetRole" />
              <div class="input-group">
                <label class="input-label">Cấp độ hiện tại</label>
                <select class="input-field" v-model="profile.level">
                  <option value="Intern">Intern</option>
                  <option value="Fresher">Fresher</option>
                  <option value="Junior">Junior</option>
                  <option value="Middle">Middle</option>
                  <option value="Senior">Senior</option>
                </select>
              </div>
            </div>

            <div style="display: flex; justify-content: flex-end; margin-top: 24px">
              <Button type="submit" :disabled="saving">
                <Save size="16" /> {{ saving ? 'Đang lưu...' : 'Lưu hồ sơ' }}
              </Button>
            </div>
          </form>
        </Card>

        <Card title="Kỹ năng chuyên môn">
          <textarea class="input-field" rows="4" placeholder="Ví dụ: ReactJS, NodeJS, TypeScript..." style="width: 100%" defaultValue="ReactJS, Redux, JavaScript, HTML, CSS, Git"></textarea>
        </Card>
      </div>

      <div style="display: flex; flex-direction: column; gap: 24px">
        <Card title="CV của bạn" style="background-color: var(--surface)">
          <div style="padding: 16px; background-color: var(--surface-soft); border-radius: var(--radius); display: flex; align-items: center; gap: 12px; border: 1px solid var(--border); margin-bottom: 16px">
            <FileText size="32" color="var(--primary)" />
            <div style="flex: 1">
              <p class="text-body" style="font-weight: 600">NguyenVanA_CV.pdf</p>
              <p class="text-helper" style="color: var(--success); display: flex; align-items: center; gap: 4px">
                <CheckCircle size="14" /> Cập nhật 2 ngày trước
              </p>
            </div>
          </div>

          <div style="border: 2px dashed var(--border); border-radius: var(--radius); padding: 24px; text-align: center; cursor: pointer">
            <Upload size="24" color="var(--text-muted)" style="margin: 0 auto 8px" />
            <p class="text-body" style="font-weight: 500">Tải CV mới lên</p>
            <p class="text-helper">PDF, DOCX (Tối đa 5MB)</p>
          </div>
          
          <p class="text-helper" style="margin-top: 16px; color: var(--text-muted); line-height: 1.5">
            CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.
          </p>
        </Card>
      </div>

    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
  </div>
</template>
