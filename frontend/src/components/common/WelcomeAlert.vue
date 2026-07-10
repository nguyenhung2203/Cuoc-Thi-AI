<script setup>
import { defineProps, defineEmits, ref, onMounted, onBeforeUnmount } from 'vue'
import { PartyPopper, CheckCircle2, Sparkles, X, ArrowRight } from 'lucide-vue-next'

const props = defineProps({
  message: String,
  title: {
    type: String,
    default: 'Thành công!'
  },
  role: {
    type: String,
    default: 'candidate'
  }
})

const emit = defineEmits(['close'])

const progress = ref(100)
let timer = null
let interval = null

onMounted(() => {
  const duration = 4500
  const step = 50
  const decrement = (step / duration) * 100

  interval = setInterval(() => {
    progress.value = Math.max(0, progress.value - decrement)
  }, step)

  timer = setTimeout(() => {
    emit('close')
  }, duration)
})

onBeforeUnmount(() => {
  if (timer) clearTimeout(timer)
  if (interval) clearInterval(interval)
})
</script>

<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-[99999] flex items-center justify-center p-4 bg-slate-950/40 backdrop-blur-md animate-fade-in" @click="emit('close')">
      <!-- Main Welcome Glass Card -->
      <div 
        class="relative w-full max-w-md bg-white dark:bg-slate-900/95 border border-slate-200 dark:border-slate-800/80 rounded-3xl p-7 shadow-2xl shadow-indigo-500/10 overflow-hidden text-center transform transition-all animate-bounce-in"
        @click.stop
      >
        <!-- Top decorative ambient glow -->
        <div class="absolute -top-24 -left-24 w-48 h-48 bg-gradient-to-br from-blue-500/30 to-purple-500/30 rounded-full blur-2xl pointer-events-none"></div>
        <div class="absolute -bottom-24 -right-24 w-48 h-48 bg-gradient-to-br from-emerald-500/20 to-indigo-500/30 rounded-full blur-2xl pointer-events-none"></div>

        <!-- Close Button -->
        <button 
          @click="emit('close')"
          class="absolute top-4 right-4 p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 rounded-full hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors z-20"
          title="Đóng"
        >
          <X size="18" />
        </button>

        <!-- Role Badge -->
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-extrabold uppercase tracking-wider mb-5" :class="role === 'recruiter' ? 'bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 border border-indigo-200 dark:border-indigo-800/60' : 'bg-emerald-50 dark:bg-emerald-950/60 text-emerald-600 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800/60'">
          <Sparkles size="13" class="animate-spin-slow" />
          <span>{{ role === 'recruiter' ? 'Nhà Tuyển Dụng Portal' : 'Ứng Viên AI Portal' }}</span>
        </div>

        <!-- Icon Circle -->
        <div class="relative mx-auto w-20 h-20 flex items-center justify-center mb-5">
          <div class="absolute inset-0 rounded-2xl transform rotate-6 scale-105 opacity-20" :class="role === 'recruiter' ? 'bg-gradient-to-tr from-indigo-600 to-purple-600' : 'bg-gradient-to-tr from-emerald-600 to-teal-600'"></div>
          <div class="w-20 h-20 rounded-2xl flex items-center justify-center shadow-xl relative z-10 text-white transform -rotate-3 transition-transform hover:rotate-0" :class="role === 'recruiter' ? 'bg-gradient-to-tr from-indigo-600 to-purple-600 shadow-indigo-500/30' : 'bg-gradient-to-tr from-emerald-500 to-teal-600 shadow-emerald-500/30'">
            <PartyPopper v-if="role === 'candidate'" size="38" class="animate-bounce" />
            <CheckCircle2 v-else size="38" />
          </div>
        </div>

        <!-- Title & Message -->
        <h2 class="text-2xl font-black tracking-tight text-slate-800 dark:text-white mb-2">
          {{ title }}
        </h2>
        <p class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed max-w-sm mx-auto mb-7 font-medium">
          {{ message }}
        </p>

        <!-- Action Button -->
        <button 
          @click="emit('close')"
          class="w-full py-3.5 px-6 rounded-2xl font-bold text-sm text-white shadow-lg flex items-center justify-center gap-2 transition-all duration-300 transform hover:-translate-y-0.5 active:translate-y-0 relative overflow-hidden group"
          :class="role === 'recruiter' ? 'bg-gradient-to-r from-indigo-600 via-purple-600 to-indigo-600 shadow-indigo-600/30 hover:shadow-indigo-600/40' : 'bg-gradient-to-r from-emerald-600 via-teal-600 to-emerald-600 shadow-emerald-600/30 hover:shadow-emerald-600/40'"
        >
          <span>Bắt đầu trải nghiệm ngay</span>
          <ArrowRight size="16" class="transition-transform group-hover:translate-x-1" />
        </button>

        <!-- Auto Dismiss Progress Bar -->
        <div class="absolute bottom-0 left-0 right-0 h-1 bg-slate-100 dark:bg-slate-800">
          <div 
            class="h-full transition-all duration-100 ease-linear"
            :class="role === 'recruiter' ? 'bg-gradient-to-r from-indigo-500 to-purple-500' : 'bg-gradient-to-r from-emerald-500 to-teal-500'"
            :style="{ width: `${progress}%` }"
          ></div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
.animate-fade-in {
  animation: fadeIn 0.2s ease-out forwards;
}

@keyframes bounceIn {
  0% { opacity: 0; transform: scale(0.8) translateY(20px); }
  70% { transform: scale(1.02) translateY(-4px); }
  100% { opacity: 1; transform: scale(1) translateY(0); }
}
.animate-bounce-in {
  animation: bounceIn 0.35s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
}

@keyframes spinSlow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin-slow {
  animation: spinSlow 8s linear infinite;
}
</style>
