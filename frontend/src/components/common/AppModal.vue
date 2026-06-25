<script setup>
import { defineProps, defineEmits } from 'vue'
import { X } from 'lucide-vue-next'

defineProps({
  isOpen: Boolean,
  title: String,
  size: {
    type: String,
    default: 'md' // sm, md, lg, xl
  }
})

const emit = defineEmits(['close'])
</script>

<template>
  <div v-if="isOpen" class="modal-overlay" @click="emit('close')">
    <div :class="['modal-content', `modal-${size}`]" @click.stop>
      <div class="modal-header">
        <h3 class="text-h3" style="margin: 0">{{ title }}</h3>
        <button class="icon-btn" @click="emit('close')">
          <X size="20" />
        </button>
      </div>
      <div class="modal-body">
        <slot></slot>
      </div>
    </div>
  </div>
</template>
