import { VALIDATION_LIMITS, VALIDATION_MESSAGES } from './constants.js'

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const VIETNAMESE_PHONE_PATTERN = /^(?:\+84|0)(?:3|5|7|8|9)\d{8}$/

export const normalizeText = (value) => value == null ? '' : String(value).trim()
export const normalizeEmail = (value) => normalizeText(value).toLowerCase()

export const isEmpty = (value) => {
  if (value == null) return true
  if (typeof value === 'string') return value.trim().length === 0
  if (Array.isArray(value)) return value.length === 0
  return false
}

export const requiredTrim = (value, message = VALIDATION_MESSAGES.REQUIRED) => (
  isEmpty(value) ? message : ''
)

export const isEmail = (value, message = VALIDATION_MESSAGES.INVALID_EMAIL) => {
  if (isEmpty(value)) return ''
  const email = normalizeEmail(value)
  return email.length <= VALIDATION_LIMITS.EMAIL_MAX_LENGTH && EMAIL_PATTERN.test(email) ? '' : message
}

export const minLength = (value, min, message = `Nội dung phải có ít nhất ${min} ký tự.`) => {
  if (isEmpty(value)) return ''
  return normalizeText(value).length >= min ? '' : message
}

export const maxLength = (value, max, message = `Nội dung không được vượt quá ${max} ký tự.`) => {
  if (value == null) return ''
  return String(value).trim().length <= max ? '' : message
}

export const isOneOf = (value, allowed, message = VALIDATION_MESSAGES.INVALID_OPTION) => {
  if (isEmpty(value)) return ''
  return Array.isArray(allowed) && allowed.includes(value) ? '' : message
}

export const validatePassword = (value, options = {}) => {
  if (isEmpty(value)) return ''
  const password = String(value)
  const min = options.min ?? VALIDATION_LIMITS.PASSWORD_MIN_LENGTH
  const max = options.max ?? VALIDATION_LIMITS.PASSWORD_MAX_LENGTH
  const requireStrong = options.requireStrong ?? true

  if (password.length < min) return options.minMessage || `Mật khẩu phải có ít nhất ${min} ký tự.`
  if (password.length > max) return options.maxMessage || `Mật khẩu không được vượt quá ${max} ký tự.`
  if (requireStrong) {
    const hasUpper = /[A-Z]/.test(password)
    const hasLower = /[a-z]/.test(password)
    const hasNumberOrSymbol = /[0-9]|[^A-Za-z0-9]/.test(password)
    if (!hasUpper || !hasLower || !hasNumberOrSymbol) {
      return options.weakMessage || VALIDATION_MESSAGES.PASSWORD_WEAK
    }
  }
  return ''
}

export const passwordsMatch = (password, confirmation, message = VALIDATION_MESSAGES.PASSWORD_MISMATCH) => (
  String(password ?? '') === String(confirmation ?? '') ? '' : message
)

export const confirmPassword = (value, originalValue, message = VALIDATION_MESSAGES.PASSWORD_MISMATCH) => passwordsMatch(originalValue, value, message)

export const isVietnamesePhone = (value, message = VALIDATION_MESSAGES.INVALID_PHONE) => {
  if (isEmpty(value)) return ''
  const phone = String(value).replace(/[\s.-]/g, '')
  return VIETNAMESE_PHONE_PATTERN.test(phone) ? '' : message
}

export const isUrl = (value, message = VALIDATION_MESSAGES.INVALID_URL) => {
  if (isEmpty(value)) return ''
  try {
    const url = new URL(normalizeText(value))
    return url.protocol === 'http:' || url.protocol === 'https:' ? '' : message
  } catch {
    return message
  }
}

export const isLinkedInUrl = (value, message = VALIDATION_MESSAGES.INVALID_LINKEDIN_URL) => {
  const urlError = isUrl(value, message)
  if (urlError || isEmpty(value)) return urlError
  const hostname = new URL(normalizeText(value)).hostname.toLowerCase()
  return hostname === 'linkedin.com' || hostname.endsWith('.linkedin.com') ? '' : message
}

export const isFutureDateTime = (value, bufferMinutes = 0, message = VALIDATION_MESSAGES.FUTURE_DATE_REQUIRED) => {
  if (isEmpty(value)) return ''
  const timestamp = new Date(value).getTime()
  if (!Number.isFinite(timestamp)) return VALIDATION_MESSAGES.INVALID_DATE
  return timestamp > Date.now() + bufferMinutes * 60_000 ? '' : message
}

export const isValidId = (value) => value !== null && value !== undefined && String(value).trim().length > 0

export const isValidDate = (value) => {
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp)
}

export const isNonNegativeNumber = (value, message = 'Giá trị phải là số không âm.') => {
  if (isEmpty(value)) return ''
  const number = Number(value)
  return Number.isFinite(number) && number >= 0 ? '' : message
}

export const maxTrimmedLength = (value, max, message = `Nội dung không được vượt quá ${max} ký tự.`) => maxLength(value, max, message)

export const isValidSalaryRange = (min, max, message = VALIDATION_MESSAGES.INVALID_SALARY_RANGE) => {
  if (isEmpty(min) || isEmpty(max)) return ''
  const minNumber = Number(min)
  const maxNumber = Number(max)
  if (!Number.isFinite(minNumber) || !Number.isFinite(maxNumber) || minNumber < 0 || maxNumber < 0) {
    return 'Mức lương phải là số không âm.'
  }
  return minNumber <= maxNumber ? '' : message
}

const getExtension = (fileName) => {
  const parts = String(fileName || '').toLowerCase().split('.')
  return parts.length > 1 ? parts.pop() : ''
}

export const validateFile = (file, options = {}) => {
  if (!file || typeof file !== 'object') return options.required === false ? '' : VALIDATION_MESSAGES.INVALID_FILE
  if (!Number.isFinite(file.size) || file.size <= 0) return VALIDATION_MESSAGES.EMPTY_FILE

  const maxBytes = options.maxBytes
  if (Number.isFinite(maxBytes) && file.size > maxBytes) {
    return options.sizeMessage || VALIDATION_MESSAGES.FILE_TOO_LARGE
  }

  const mimeTypes = options.mimeTypes || []
  const extensions = (options.extensions || []).map(item => String(item).replace(/^\./, '').toLowerCase())
  const mimeValid = mimeTypes.length === 0 || mimeTypes.includes(file.type)
  const extensionValid = extensions.length === 0 || extensions.includes(getExtension(file.name))

  return mimeValid && extensionValid ? '' : options.typeMessage || VALIDATION_MESSAGES.INVALID_FILE_TYPE
}

export const hasDuplicateNormalized = (values) => {
  if (!Array.isArray(values)) return false
  const normalized = values.map(value => normalizeText(value).toLocaleLowerCase('vi')).filter(Boolean)
  return new Set(normalized).size !== normalized.length
}

export const validateField = (value, rules = [], values = {}) => {
  for (const rule of rules) {
    if (typeof rule !== 'function') continue
    const message = rule(value, values)
    if (message) return message
  }
  return ''
}

export const validateForm = (values, schema) => {
  const errors = {}
  for (const [field, rules] of Object.entries(schema || {})) {
    const message = validateField(values?.[field], Array.isArray(rules) ? rules : [rules], values || {})
    if (message) errors[field] = message
  }
  return { isValid: Object.keys(errors).length === 0, errors }
}

export const hasValidationErrors = (errors) => Boolean(errors && Object.keys(errors).length)

const firstMessage = (value) => {
  if (Array.isArray(value)) return firstMessage(value[0])
  if (typeof value === 'string') return value
  if (value && typeof value === 'object') return value.message || ''
  return ''
}

export const normalizeValidationErrors = (error) => {
  const source = error?.validationErrors || error?.fieldErrors || error?.errors
  const result = {}

  if (source && typeof source === 'object' && !Array.isArray(source)) {
    for (const [field, value] of Object.entries(source)) {
      const message = firstMessage(value)
      if (message) result[field] = message
    }
  }

  const details = error?.details
  if (Array.isArray(details)) {
    for (const detail of details) {
      if (detail && typeof detail === 'object' && detail.field) {
        result[detail.field] = detail.message || VALIDATION_MESSAGES.REQUIRED
      }
    }
  }

  if (error?.field && error?.message) result[error.field] = error.message
  if (!Object.keys(result).length && error?.message) result._form = error.message
  return result
}

export default {
  normalizeText,
  normalizeEmail,
  isEmpty,
  requiredTrim,
  isEmail,
  minLength,
  maxLength,
  isOneOf,
  validatePassword,
  confirmPassword,
  passwordsMatch,
  isVietnamesePhone,
  isUrl,
  isLinkedInUrl,
  isFutureDateTime,
  isValidSalaryRange,
  isNonNegativeNumber,
  isValidId,
  isValidDate,
  maxTrimmedLength,
  validateFile,
  hasDuplicateNormalized,
  validateField,
  validateForm,
  hasValidationErrors,
  normalizeValidationErrors,
}
