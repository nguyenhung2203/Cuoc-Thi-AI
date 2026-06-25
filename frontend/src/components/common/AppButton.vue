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
