import { computed, reactive } from 'vue'
import { validateField as runFieldValidation, validateForm as runFormValidation } from '../utils/validators.js'

export function useFormValidation(initialValues = {}, schema = {}) {
  const errors = reactive({})
  const touched = reactive({})

  const setError = (field, message) => {
    if (message) errors[field] = message
    else delete errors[field]
  }
  const clearError = (field) => { delete errors[field] }
  const touchField = (field) => { touched[field] = true }
  const validateField = (field, values = initialValues) => {
    const rules = schema[field]
    const message = runFieldValidation(values?.[field], Array.isArray(rules) ? rules : [rules], values)
    setError(field, message)
    return !message
  }
  const validateForm = (values = initialValues) => {
    const result = runFormValidation(values, schema)
    Object.keys(errors).forEach(key => delete errors[key])
    Object.assign(errors, result.errors)
    Object.keys(schema).forEach(key => { touched[key] = true })
    return result.isValid
  }
  const resetErrors = () => {
    Object.keys(errors).forEach(key => delete errors[key])
    Object.keys(touched).forEach(key => delete touched[key])
  }
  return { errors, touched, setError, clearError, validateField, validateForm, touchField, resetErrors, hasErrors: computed(() => Object.keys(errors).length > 0) }
}

export default useFormValidation
