<script setup>
defineProps({
  variant: {
    type: String,
    default: 'primary'
  },
  className: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <button :class="['btn', `btn-${variant}`, className]" v-bind="$attrs">
    <slot></slot>
  </button>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 20px;
  font-family: var(--sans);
  font-size: 14px;
  font-weight: 600;
  border-radius: var(--radius-full) !important;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  outline: none;
  white-space: nowrap;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

.btn:active:not(:disabled) {
  transform: scale(0.96);
}

.btn-primary {
  background: var(--primary);
  color: white;
  box-shadow: var(--shadow-sm);
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-hover);
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

.btn-secondary {
  background: var(--surface);
  color: var(--text-main);
  border-color: var(--border);
  box-shadow: var(--shadow-sm);
}

.btn-secondary:hover:not(:disabled) {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--surface-soft);
}

.btn-danger {
  background: var(--danger);
  color: white;
  box-shadow: var(--shadow-sm);
}

.btn-danger:hover:not(:disabled) {
  background: #b91c1c;
  box-shadow: var(--shadow-md);
}

.btn-ghost {
  background: transparent;
  color: var(--text-secondary);
}

.btn-ghost:hover:not(:disabled) {
  background: var(--surface-soft);
  color: var(--primary);
}

.btn-outline {
  background: transparent;
  color: var(--primary);
  border-color: var(--primary-light);
}

.btn-outline:hover:not(:disabled) {
  background: var(--primary-light);
  border-color: var(--primary);
}
</style>
