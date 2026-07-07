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
  <div class="space-y-8 animate-fade-in pb-12 max-w-6xl mx-auto">
    <div class="mb-8">
      <h1 class="text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-indigo-600 to-purple-600">Hồ sơ cá nhân & CV</h1>
      <p class="text-gray-500 mt-2 text-lg">Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      
      <!-- Left Column: Forms -->
      <div class="lg:col-span-2 space-y-8">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden">
          <div class="bg-gradient-to-r from-indigo-50 to-purple-50 p-6 border-b border-indigo-100">
            <h3 class="text-xl font-bold text-indigo-900 flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
              Thông tin cơ bản
            </h3>
          </div>
          
          <div class="p-6">
            <form @submit="handleSave" class="space-y-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Custom styling for inputs can be handled by AppInput if it uses standard HTML attributes -->
                <!-- We'll wrap them in div to provide some spacing if needed -->
                <div class="space-y-1">
                  <Input label="Họ và Tên" v-model="profile.name" required />
                </div>
                <div class="space-y-1">
                  <Input label="Email" type="email" v-model="profile.email" disabled />
                </div>
                <div class="space-y-1">
                  <Input label="Số điện thoại" v-model="profile.phone" />
                </div>
                <div class="space-y-1">
                  <Input label="LinkedIn Profile" v-model="profile.linkedin" />
                </div>
              </div>

              <div class="pt-6 mt-6 border-t border-gray-100">
                <h3 class="text-lg font-bold text-gray-800 mb-6 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                  Định hướng nghề nghiệp
                </h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <Input label="Vị trí mục tiêu (Target Role)" v-model="profile.targetRole" />
                  <div class="flex flex-col gap-2">
                    <label class="text-sm font-semibold text-gray-700">Cấp độ hiện tại</label>
                    <select class="w-full px-4 py-2.5 rounded-xl border border-gray-200 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 outline-none transition-all bg-white text-gray-800" v-model="profile.level">
                      <option value="Intern">Intern</option>
                      <option value="Fresher">Fresher</option>
                      <option value="Junior">Junior</option>
                      <option value="Middle">Middle</option>
                      <option value="Senior">Senior</option>
                    </select>
                  </div>
                </div>
              </div>

              <div class="flex justify-end pt-6">
                <button type="submit" :disabled="saving" class="bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-bold py-3 px-8 rounded-xl shadow-lg hover:shadow-indigo-500/30 transition-all duration-300 transform hover:-translate-y-0.5 disabled:opacity-50 disabled:transform-none flex items-center gap-2">
                  <Save class="w-5 h-5" /> 
                  {{ saving ? 'Đang lưu...' : 'Lưu hồ sơ' }}
                </button>
              </div>
            </form>
          </div>
        </Card>

        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden">
          <div class="bg-gradient-to-r from-emerald-50 to-teal-50 p-6 border-b border-emerald-100">
            <h3 class="text-xl font-bold text-emerald-900 flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" /></svg>
              Kỹ năng chuyên môn
            </h3>
          </div>
          <div class="p-6">
            <textarea class="w-full p-4 border border-gray-200 rounded-xl focus:border-emerald-500 focus:ring-4 focus:ring-emerald-500/20 outline-none transition-all text-gray-800 font-medium resize-none" rows="4" placeholder="Ví dụ: ReactJS, NodeJS, TypeScript..." defaultValue="ReactJS, Redux, JavaScript, HTML, CSS, Git"></textarea>
            <p class="text-sm text-gray-500 mt-2">Phân cách các kỹ năng bằng dấu phẩy (,)</p>
          </div>
        </Card>
      </div>

      <!-- Right Column: CV Upload -->
      <div class="space-y-8">
        <Card class="bg-white/90 backdrop-blur-md shadow-xl border-0 rounded-2xl overflow-hidden sticky top-6">
          <div class="bg-gradient-to-r from-blue-50 to-cyan-50 p-6 border-b border-blue-100">
            <h3 class="text-xl font-bold text-blue-900 flex items-center gap-2">
              <FileText class="w-6 h-6 text-blue-600" />
              CV của bạn
            </h3>
          </div>
          
          <div class="p-6 space-y-6">
            <!-- Current CV -->
            <div class="bg-blue-50/50 p-4 rounded-xl border border-blue-100 flex items-center gap-4 hover:shadow-md transition-shadow group cursor-pointer">
              <div class="w-12 h-12 bg-white rounded-lg shadow-sm flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform">
                <FileText class="w-6 h-6 text-blue-600" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="font-bold text-gray-800 truncate">NguyenVanA_CV.pdf</p>
                <p class="text-sm text-emerald-600 flex items-center gap-1 mt-1 font-medium">
                  <CheckCircle class="w-4 h-4" /> Cập nhật 2 ngày trước
                </p>
              </div>
            </div>

            <!-- Upload Area -->
            <div class="group border-2 border-dashed border-gray-200 hover:border-blue-400 bg-gray-50 hover:bg-blue-50/50 rounded-xl p-8 text-center cursor-pointer transition-all duration-300 relative overflow-hidden" @click="router.push('/cv-upload')">
              <div class="absolute inset-0 bg-gradient-to-br from-blue-400/0 to-indigo-400/0 group-hover:from-blue-400/10 group-hover:to-indigo-400/10 transition-all duration-500"></div>
              <div class="w-16 h-16 bg-white rounded-full shadow-sm flex items-center justify-center mx-auto mb-4 group-hover:-translate-y-2 group-hover:shadow-md transition-all duration-300">
                <Upload class="w-8 h-8 text-blue-500" />
              </div>
              <p class="text-gray-900 font-bold mb-1 group-hover:text-blue-700 transition-colors">Tải CV mới lên</p>
              <p class="text-gray-500 text-sm font-medium">PDF, DOCX (Tối đa 5MB)</p>
            </div>
            
            <div class="bg-amber-50 border border-amber-100 rounded-xl p-4 flex gap-3">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-500 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              <p class="text-amber-800 text-sm leading-relaxed font-medium">
                CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.
              </p>
            </div>
          </div>
        </Card>
      </div>

    </div>

    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
  </div>
</template>

<style scoped>
@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}
</style>
