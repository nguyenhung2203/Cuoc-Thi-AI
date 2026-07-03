<script setup>
defineProps({
  type: {
    type: String,
    default: 'default' // primary, success, warning, danger, default
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
