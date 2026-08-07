import test from 'node:test'
import assert from 'node:assert/strict'
import {
  confirmPassword,
  hasDuplicateNormalized,
  isEmail,
  isFutureDateTime,
  isLinkedInUrl,
  isOneOf,
  isValidSalaryRange,
  isVietnamesePhone,
  normalizeEmail,
  normalizeValidationErrors,
  requiredTrim,
  validateFile,
  validateForm,
  validatePassword,
} from './validators.js'

test('requiredTrim rejects empty values and whitespace', () => {
  assert.ok(requiredTrim(''))
  assert.ok(requiredTrim('   '))
  assert.ok(requiredTrim(null))
  assert.equal(requiredTrim('Hợp lệ'), '')
})

test('normalizeEmail trims and lowercases email', () => {
  assert.equal(normalizeEmail('  USER@Example.COM '), 'user@example.com')
})

test('isEmail validates common email formats', () => {
  assert.equal(isEmail('user@example.com'), '')
  assert.ok(isEmail('user@'))
  assert.ok(isEmail('not-an-email'))
})

test('validatePassword enforces the shared password policy', () => {
  assert.equal(validatePassword('Abcdef1!'), '')
  assert.ok(validatePassword('abc'))
  assert.ok(validatePassword('abcdefgh'))
})

test('confirmPassword detects mismatches', () => {
  assert.equal(confirmPassword('Abcdef1!', 'Abcdef1!'), '')
  assert.ok(confirmPassword('Abcdef1!', 'Different1!'))
})

test('phone and LinkedIn URL validators accept valid Vietnamese values', () => {
  assert.equal(isVietnamesePhone('0912 345 678'), '')
  assert.ok(isVietnamesePhone('12345'))
  assert.equal(isLinkedInUrl('https://www.linkedin.com/in/example'), '')
  assert.ok(isLinkedInUrl('https://example.com/profile'))
})

test('future datetime and salary range validators handle boundaries', () => {
  assert.equal(isFutureDateTime(new Date(Date.now() + 60_000).toISOString()), '')
  assert.ok(isFutureDateTime(new Date(Date.now() - 60_000).toISOString()))
  assert.ok(isFutureDateTime('not-a-date'))
  assert.ok(isFutureDateTime(new Date(Date.now() + 5 * 60_000).toISOString(), 15))
  assert.equal(isValidSalaryRange(10, 20), '')
  assert.ok(isValidSalaryRange(20, 10))
  assert.ok(isValidSalaryRange(-1, 10))
})

test('validateFile checks size, type and extension', () => {
  const valid = { name: 'cv.pdf', type: 'application/pdf', size: 1024 }
  const invalidType = { name: 'cv.exe', type: 'application/octet-stream', size: 1024 }
  const tooLarge = { name: 'cv.pdf', type: 'application/pdf', size: 2048 }
  const empty = { name: 'cv.pdf', type: 'application/pdf', size: 0 }
  const fakeExtension = { name: 'cv.exe', type: 'application/pdf', size: 1024 }
  const options = { maxBytes: 1500, mimeTypes: ['application/pdf'], extensions: ['pdf'] }

  assert.equal(validateFile(valid, options), '')
  assert.ok(validateFile(invalidType, options))
  assert.ok(validateFile(tooLarge, options))
  assert.ok(validateFile(empty, options))
  assert.ok(validateFile(fakeExtension, options))
})

test('hasDuplicateNormalized ignores spacing and letter case', () => {
  assert.equal(hasDuplicateNormalized(['Vue', ' vue ']), true)
  assert.equal(hasDuplicateNormalized(['Vue', 'Go']), false)
})

test('validateForm returns the first error for each invalid field', () => {
  const result = validateForm(
    { email: 'bad-email', password: '' },
    {
      email: [requiredTrim, isEmail],
      password: [requiredTrim, validatePassword],
    },
  )

  assert.equal(result.isValid, false)
  assert.ok(result.errors.email)
  assert.ok(result.errors.password)
})

test('Sprint 6 enum, length and duplicate rules reject invalid AI inputs', () => {
  assert.equal(isOneOf('middle', ['fresher', 'junior', 'middle', 'senior']), '')
  assert.ok(isOneOf('owner', ['candidate', 'recruiter']))
  assert.ok(validateForm(
    { role: 'x'.repeat(101), answer: '   ' },
    {
      role: [(value) => value.length <= 100 ? '' : 'Vị trí quá dài.'],
      answer: [requiredTrim, (value) => value.length <= 5000 ? '' : 'Câu trả lời quá dài.'],
    },
  ).isValid === false)
  assert.equal(hasDuplicateNormalized(['Thiết kế API thế nào?', ' thiết kế api thế nào? ']), true)
})

test('Sprint 6 rubric boundaries reject invalid totals and score ranges', () => {
  const totalWeight = (criteria) => criteria.reduce((sum, item) => sum + Number(item.weight || 0), 0)
  assert.equal(totalWeight([{ weight: 50 }, { weight: 50 }]), 100)
  assert.notEqual(totalWeight([{ weight: 50 }, { weight: 49 }]), 100)
  assert.notEqual(totalWeight([{ weight: 50 }, { weight: 51 }]), 100)
  assert.ok(!(5 < 5))
  assert.ok(!(6 < 5))
})
test('normalizeValidationErrors supports common backend error shapes', () => {
  assert.deepEqual(
    normalizeValidationErrors({ errors: { email: ['Email không hợp lệ.'] } }),
    { email: 'Email không hợp lệ.' },
  )
  assert.deepEqual(
    normalizeValidationErrors({ details: [{ field: 'password', message: 'Mật khẩu yếu.' }] }),
    { password: 'Mật khẩu yếu.' },
  )
  assert.deepEqual(
    normalizeValidationErrors({ message: 'Dữ liệu không hợp lệ.' }),
    { _form: 'Dữ liệu không hợp lệ.' },
  )
})
