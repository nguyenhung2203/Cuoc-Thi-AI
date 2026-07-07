<script setup>
defineProps({
  type: {
    type: String,
    default: 'default' // primary, success, warning, danger, default, info
  },
  className: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <span :class="['badge', `badge-${type}`, className]" v-bind="$attrs">
    <slot></slot>
  </span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-radius: var(--radius-full);
  white-space: nowrap;
}

.badge-default {
  background: var(--bg-secondary);
  color: var(--text);
  border: 1px solid var(--border);
}

.badge-primary {
  background: var(--accent-bg);
  color: var(--primary);
  border: 1px solid var(--accent-border);
}

.badge-info {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.badge-success {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.badge-warning {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.badge-danger {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}
</style>
