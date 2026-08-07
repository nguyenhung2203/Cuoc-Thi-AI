const STORAGE_KEY = 'candidate_saved_jobs'

/**
 * Read saved jobs from localStorage safely (corrupt JSON → empty list).
 * @returns {Array<object>}
 */
export function loadSavedJobs() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    localStorage.removeItem(STORAGE_KEY)
    return []
  }
}

/**
 * Persist saved jobs list.
 * @param {Array<object>} jobs
 */
export function saveSavedJobs(jobs) {
  const list = Array.isArray(jobs) ? jobs : []
  localStorage.setItem(STORAGE_KEY, JSON.stringify(list))
}

/**
 * @returns {string[]}
 */
export function loadSavedJobIds() {
  return loadSavedJobs().map((j) => j?.id).filter(Boolean)
}

/**
 * Toggle save state for a job snapshot. Returns the updated list.
 * @param {object} job
 * @returns {{ list: Array<object>, saved: boolean }}
 */
export function toggleSavedJob(job) {
  if (!job?.id) return { list: loadSavedJobs(), saved: false }
  let list = loadSavedJobs()
  const exists = list.some((item) => item.id === job.id)
  if (exists) {
    list = list.filter((item) => item.id !== job.id)
    saveSavedJobs(list)
    return { list, saved: false }
  }
  list.push({
    id: job.id,
    company_id: job.company_id,
    title: job.title,
    company_name: job.company_name,
    location: job.location,
    salary_min: job.salary_min != null ? { Valid: true, Int64: job.salary_min || 0 } : { Valid: false, Int64: 0 },
    salary_max: job.salary_max != null ? { Valid: true, Int64: job.salary_max || 0 } : { Valid: false, Int64: 0 },
    currency: { Valid: true, String: job.currency || 'VND' },
    saved_at: new Date().toISOString(),
  })
  saveSavedJobs(list)
  return { list, saved: true }
}
