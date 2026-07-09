<script setup>
// Using inheritAttrs: false so we can manually merge the outer class
// onto the inner .card div via v-bind="$attrs"
defineOptions({ inheritAttrs: false })
</script>

<template>
  <div class="card" v-bind="$attrs">
    <slot></slot>
  </div>
</template>

<style scoped>
.card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border);
  transition: box-shadow 0.3s ease, transform 0.3s ease;
  position: relative;
  overflow: hidden;
}

.card:hover {
  box-shadow: var(--shadow-md);
}

.card.glass-panel {
  background: var(--surface-soft);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}
</style>
