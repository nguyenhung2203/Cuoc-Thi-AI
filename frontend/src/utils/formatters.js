/** Convert absolute VND amounts to TopCV-style "triệu" labels. */
export function formatSalaryTrieu(min, max) {
  const toTrieu = (n) => {
    const num = Number(n)
    if (!Number.isFinite(num) || num <= 0) return null
    const t = num / 1_000_000
    if (Number.isInteger(t)) return String(t)
    return String(Math.round(t * 10) / 10).replace(/\.0$/, '')
  }

  const a = toTrieu(min)
  const b = toTrieu(max)
  if (!a && !b) return 'Thoả thuận'
  if (a && !b) return `Từ ${a} triệu`
  if (!a && b) return `Đến ${b} triệu`
  if (a === b) return `${a} triệu`
  return `${a} - ${b} triệu`
}

const EXPERIENCE_BY_LEVEL = {
  intern: 'Không yêu cầu',
  junior: 'Dưới 1 năm',
  middle: '1 - 3 năm',
  senior: '3 - 5 năm',
  lead: 'Trên 5 năm',
}

/** Map job.level → nhãn kinh nghiệm kiểu bảng tin tuyển dụng. */
export function formatExperience(level) {
  if (!level) return 'Không yêu cầu'
  const key = String(level).toLowerCase().trim()
  return EXPERIENCE_BY_LEVEL[key] || level
}

export function experienceFilterOptions() {
  return Object.entries(EXPERIENCE_BY_LEVEL).map(([value, label]) => ({ value, label }))
}
