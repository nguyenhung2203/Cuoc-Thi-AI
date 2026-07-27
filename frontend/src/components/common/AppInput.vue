<script setup>
import { computed, defineProps, defineEmits } from 'vue'

const props = defineProps({
  label: String, type: { type: String, default: 'text' }, modelValue: [String, Number], placeholder: String,
  required: Boolean, error: String, minlength: [String, Number], maxlength: [String, Number], pattern: String,
  autocomplete: String, id: String, name: String, disabled: Boolean, className: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue', 'blur'])
const inputId = computed(() => props.id || props.name || `input-${Math.random().toString(36).slice(2)}`)
const errorId = computed(() => `${inputId.value}-error`)
</script>

<template>
  <div :class="['input-group', className]">
    <label v-if="label" :for="inputId" class="input-label">{{ label }} <span v-if="required" style="color: var(--danger)">*</span></label>
    <input :id="inputId" :name="name" :type="type" class="input-field" :placeholder="placeholder" :value="modelValue"
      @input="emit('update:modelValue', $event.target.value)" @blur="emit('blur', $event)" :required="required"
      :minlength="minlength" :maxlength="maxlength" :pattern="pattern" :autocomplete="autocomplete" :disabled="disabled"
      :aria-invalid="error ? 'true' : 'false'" :aria-describedby="error ? errorId : undefined" v-bind="$attrs" :class="{ 'has-error': error }" />
    <span v-if="error" :id="errorId" class="error-text" style="font-size: 12px; margin-top: 4px; display: block; color: var(--danger);">{{ error }}</span>
  </div>
</template>
