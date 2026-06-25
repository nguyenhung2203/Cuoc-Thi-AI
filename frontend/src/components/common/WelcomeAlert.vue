<script setup>
import { defineProps, defineEmits, onMounted } from 'vue'
import { PartyPopper, CheckCircle } from 'lucide-vue-next'
import Button from './AppButton.vue'

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

onMounted(() => {
  // Auto close after 4 seconds
  setTimeout(() => {
    emit('close')
  }, 4000)
})
</script>

<template>
  <div class="alert-overlay" @click="emit('close')">
    <div class="alert-box" @click.stop>
      <div class="alert-icon-wrapper" :class="role === 'recruiter' ? 'bg-primary' : 'bg-success'">
        <PartyPopper v-if="role === 'candidate'" size="32" color="white" />
        <CheckCircle v-else size="32" color="white" />
      </div>
      
      <h2 class="alert-title">{{ title }}</h2>
      <p class="alert-message">{{ message }}</p>
      
      <Button 
        :variant="role === 'recruiter' ? 'primary' : 'primary'" 
        style="width: 100%; margin-top: 24px"
        @click="emit('close')"
      >
        Bắt đầu làm việc
      </Button>
    </div>
  </div>
</template>

<style scoped>
.alert-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10000;
  animation: fadeIn 0.3s ease-out;
}

.alert-box {
  background-color: var(--surface);
  border-radius: 20px;
  padding: 32px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  animation: scaleUp 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.alert-icon-wrapper {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

.bg-primary {
  background: linear-gradient(135deg, var(--primary), #60A5FA);
}

.bg-success {
  background: linear-gradient(135deg, var(--success), #34D399);
}

.alert-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-main);
  margin-bottom: 8px;
}

.alert-message {
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.5;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes scaleUp {
  from { transform: scale(0.8); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
</style>
