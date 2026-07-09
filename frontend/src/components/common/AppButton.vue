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
  background: linear-gradient(135deg, var(--primary), var(--accent));
  color: white;
  box-shadow: var(--shadow-md);
}

.btn-primary:hover:not(:disabled) {
  box-shadow: var(--shadow-glow);
  transform: translateY(-1px);
}

.btn-secondary {
  background: var(--surface);
  color: var(--text-h);
  border-color: var(--border);
  box-shadow: var(--shadow-sm);
}

.btn-secondary:hover:not(:disabled) {
  border-color: var(--primary-light);
  color: var(--primary);
  background: var(--bg);
}

.btn-danger {
  background: linear-gradient(135deg, var(--danger), #f43f5e);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  box-shadow: 0 8px 24px -4px rgba(239, 68, 68, 0.4);
}

.btn-ghost {
  background: transparent;
  color: var(--text);
}

.btn-ghost:hover:not(:disabled) {
  background: var(--accent-bg);
  color: var(--primary);
}
</style>
