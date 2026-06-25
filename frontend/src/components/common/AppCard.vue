<script setup>
defineProps({
  className: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <div :class="['card', className]" v-bind="$attrs">
    <slot></slot>
  </div>
</template>
