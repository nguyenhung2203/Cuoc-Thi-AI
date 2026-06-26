<script setup>
import { defineProps, defineEmits, onMounted, onUnmounted } from 'vue'
import { X, CheckCircle, AlertCircle, Info } from 'lucide-vue-next'

const props = defineProps({
  message: String,
  type: {
    type: String,
    default: 'info' // success, error, info, warning
  },
  duration: {
    type: Number,
    default: 3000
  },
  position: {
    type: String,
    default: 'bottom-right' // bottom-right, top-center
  }
})

const emit = defineEmits(['close'])

let timer = null

onMounted(() => {
  if (props.duration > 0) {
    timer = setTimeout(() => {
      emit('close')
    }, props.duration)
  }
})

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div :class="['toast-container', `pos-${position}`]">
    <div :class="['toast', `toast-${type}`]">
      <div class="toast-icon">
        <CheckCircle v-if="type === 'success'" size="20" />
        <AlertCircle v-else-if="type === 'error' || type === 'warning'" size="20" />
        <Info v-else size="20" />
      </div>
      <div class="toast-content">{{ message }}</div>
      <button class="toast-close" @click="emit('close')">
        <X size="16" />
      </button>
    </div>
  </div>
</template>
