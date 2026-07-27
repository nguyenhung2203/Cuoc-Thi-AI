export const VALIDATION_LIMITS = Object.freeze({
  NAME_MIN_LENGTH: 2,
  NAME_MAX_LENGTH: 255,
  EMAIL_MAX_LENGTH: 254,
  PASSWORD_MIN_LENGTH: 8,
  PASSWORD_MAX_LENGTH: 128,
  PHONE_MIN_LENGTH: 9,
  PHONE_MAX_LENGTH: 15,
  SHORT_TEXT_MAX_LENGTH: 255,
  LONG_TEXT_MAX_LENGTH: 5000,
  CV_MAX_BYTES: 5 * 1024 * 1024,
  AVATAR_MAX_BYTES: 2 * 1024 * 1024,
})

export const VALIDATION_MESSAGES = Object.freeze({
  REQUIRED: 'Vui lòng nhập thông tin này.',
  INVALID_EMAIL: 'Địa chỉ email không đúng định dạng.',
  INVALID_PHONE: 'Số điện thoại không đúng định dạng.',
  INVALID_URL: 'Đường dẫn không đúng định dạng.',
  INVALID_LINKEDIN_URL: 'Đường dẫn LinkedIn không đúng định dạng.',
  PASSWORD_TOO_SHORT: `Mật khẩu phải có ít nhất ${VALIDATION_LIMITS.PASSWORD_MIN_LENGTH} ký tự.`,
  PASSWORD_TOO_LONG: `Mật khẩu không được vượt quá ${VALIDATION_LIMITS.PASSWORD_MAX_LENGTH} ký tự.`,
  PASSWORD_WEAK: 'Mật khẩu phải có chữ hoa, chữ thường và số hoặc ký tự đặc biệt.',
  PASSWORD_MISMATCH: 'Xác nhận mật khẩu không khớp.',
  INVALID_OPTION: 'Giá trị đã chọn không hợp lệ.',
  INVALID_DATE: 'Ngày giờ không hợp lệ.',
  FUTURE_DATE_REQUIRED: 'Ngày giờ phải nằm trong tương lai.',
  INVALID_SALARY_RANGE: 'Mức lương tối thiểu không được lớn hơn mức lương tối đa.',
  INVALID_FILE: 'Vui lòng chọn một tệp hợp lệ.',
  EMPTY_FILE: 'Tệp được chọn không có nội dung.',
  INVALID_FILE_TYPE: 'Định dạng tệp không được hỗ trợ.',
  FILE_TOO_LARGE: 'Tệp vượt quá dung lượng cho phép.',
})

export const USER_ROLES = Object.freeze(['candidate', 'recruiter'])
export const JOB_STATUSES = Object.freeze(['draft', 'open', 'closed'])
export const CANDIDATE_STATUSES = Object.freeze(['new', 'screening', 'interviewing', 'offered', 'hired', 'rejected'])
export const JOB_LEVELS = Object.freeze(['intern', 'fresher', 'junior', 'middle', 'senior', 'lead'])
export const EMPLOYMENT_TYPES = Object.freeze(['full_time', 'part_time', 'contract', 'internship', 'remote'])
export const INTERVIEW_STATUSES = Object.freeze(['scheduled', 'active', 'completed', 'cancelled'])

export const CV_FILE_RULES = Object.freeze({
  maxBytes: VALIDATION_LIMITS.CV_MAX_BYTES,
  mimeTypes: Object.freeze([
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  ]),
  extensions: Object.freeze(['pdf', 'doc', 'docx']),
})

export const AVATAR_FILE_RULES = Object.freeze({
  maxBytes: VALIDATION_LIMITS.AVATAR_MAX_BYTES,
  mimeTypes: Object.freeze(['image/jpeg', 'image/png', 'image/webp']),
  extensions: Object.freeze(['jpg', 'jpeg', 'png', 'webp']),
})
