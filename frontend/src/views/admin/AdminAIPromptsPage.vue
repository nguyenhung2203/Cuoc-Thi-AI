<script setup>
import { ref, onMounted } from 'vue'
import { adminService } from '../../services/admin.service'
import Card from '../../components/common/AppCard.vue'
import Button from '../../components/common/AppButton.vue'
import Toast from '../../components/common/AppToast.vue'
import { Bot, Save, AlertCircle, Edit2, Play, FileText } from 'lucide-vue-next'

const prompts = ref([])
const loading = ref(true)
const selectedPrompt = ref(null)
const editingContent = ref('')
const toast = ref(null)

onMounted(async () => {
  await fetchPrompts()
})

const fetchPrompts = async () => {
  loading.value = true
  try {
    const data = await adminService.getSystemPrompts()
    // Data is an array of prompt templates
    prompts.value = data || []
    if (prompts.value.length > 0) {
      selectPrompt(prompts.value[0])
    }
  } catch (err) {
    toast.value = { type: 'error', message: 'Lỗi tải danh sách prompts' }
  } finally {
    loading.value = false
  }
}

const selectPrompt = (prompt) => {
  selectedPrompt.value = prompt
  editingContent.value = prompt.content
}

const savePrompt = async () => {
  if (!selectedPrompt.value) return
  
  try {
    const payload = {
      name: selectedPrompt.value.name,
      content: editingContent.value,
      model: selectedPrompt.value.model,
      variables_schema: selectedPrompt.value.variables_schema,
      params: selectedPrompt.value.params
    }
    await adminService.updateSystemPrompt(payload)
    toast.value = { type: 'success', message: 'Cập nhật prompt thành công! (Tạo phiên bản mới)' }
    await fetchPrompts()
  } catch (err) {
    toast.value = { type: 'error', message: 'Có lỗi khi lưu prompt' }
  }
}

</script>

<template>
  <div class="min-h-screen bg-slate-50 p-8">
    <Toast v-if="toast" :type="toast.type" :message="toast.message" @close="toast = null" />
    
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center gap-3 mb-8">
        <div class="p-3 bg-indigo-100 text-indigo-600 rounded-xl">
          <Bot size="28" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-slate-800">Admin AI Prompts Control</h1>
          <p class="text-slate-500">Quản lý các system prompts của nền tảng</p>
        </div>
      </div>

      <div class="grid grid-cols-12 gap-6 h-[75vh]">
        <!-- Sidebar: List of Prompts -->
        <Card class="col-span-4 flex flex-col h-full overflow-hidden">
          <div class="p-4 border-b border-slate-200 bg-slate-50">
            <h3 class="font-semibold text-slate-700">Danh sách Prompt</h3>
          </div>
          <div class="flex-1 overflow-y-auto p-2">
            <div v-if="loading" class="p-4 text-center text-slate-500">Đang tải...</div>
            <div v-else-if="prompts.length === 0" class="p-4 text-center text-slate-500">Không có data</div>
            
            <div v-for="prompt in prompts" :key="prompt.name"
                 @click="selectPrompt(prompt)"
                 :class="[
                   'p-4 rounded-lg cursor-pointer mb-2 transition-colors border',
                   selectedPrompt?.name === prompt.name 
                    ? 'border-indigo-500 bg-indigo-50/50 shadow-sm' 
                    : 'border-transparent hover:bg-slate-100 text-slate-600 hover:text-slate-900'
                 ]">
              <div class="font-medium mb-1">{{ prompt.name }}</div>
              <div class="flex items-center gap-4 text-xs">
                <span class="flex items-center gap-1 text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded">
                  v{{ prompt.version }}
                </span>
                <span class="text-slate-400">Model: {{ prompt.model }}</span>
              </div>
            </div>
          </div>
        </Card>

        <!-- Main Content: Editor -->
        <Card class="col-span-8 flex flex-col h-full overflow-hidden" v-if="selectedPrompt">
          <div class="p-4 border-b border-slate-200 bg-slate-50 flex justify-between items-center">
            <div>
              <h3 class="font-bold text-lg text-slate-800">{{ selectedPrompt.name }}</h3>
              <p class="text-sm text-slate-500">Phiên bản hiện tại: v{{ selectedPrompt.version }}</p>
            </div>
            <Button @click="savePrompt">
              <Save size="18" class="mr-2" /> Lưu phiên bản mới
            </Button>
          </div>
          <div class="flex-1 p-0 flex flex-col relative">
            <div class="bg-amber-50 border-b border-amber-200 px-4 py-3 text-sm text-amber-800 flex items-start gap-2">
              <AlertCircle size="18" class="shrink-0 mt-0.5" />
              <div>
                <strong>Chú ý:</strong> Đây là System Prompt điều khiển hành vi của AI. Mọi thay đổi sẽ lập tức áp dụng cho toàn hệ thống. Hãy kiểm tra kỹ trước khi lưu.
                Khi lưu, hệ thống sẽ tự động tạo một version mới (immutable).
              </div>
            </div>
            <div class="p-4 bg-slate-50 border-b border-slate-200 text-sm font-mono text-slate-600">
              <span class="font-semibold text-slate-700">Variables Schema: </span> 
              {{ selectedPrompt.variables_schema ? Object.keys(selectedPrompt.variables_schema).join(', ') : 'None' }}
            </div>
            <textarea
              v-model="editingContent"
              class="flex-1 p-6 font-mono text-sm resize-none focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500 bg-slate-900 text-green-400 leading-relaxed"
              spellcheck="false"
            ></textarea>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>
