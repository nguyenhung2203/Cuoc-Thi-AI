<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  label: String,
  type: {
    type: String,
    default: 'text'
  },
  modelValue: [String, Number],
  placeholder: String,
  required: Boolean,
  error: String,
  className: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div :class="['input-group', className]">
    <label v-if="label" class="input-label">
      {{ label }} <span v-if="required" style="color: var(--danger)">*</span>
    </label>
    <input 
      :type="type" 
      class="input-field" 
      :placeholder="placeholder"
      :value="modelValue"
      @input="emit('update:modelValue', $event.target.value)"
      :required="required"
      v-bind="$attrs"
    />
    <span v-if="error" class="error-text" style="font-size: 12px; margin-top: 4px; display: block;">{{ error }}</span>
  </div>
</template>
