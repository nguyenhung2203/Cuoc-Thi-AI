import { apiService } from './api.service'

export const publicService = {
  getCompanyJobs: async (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString()
    return await apiService.get(`/public/companies/${companyId}/jobs${query ? `?${query}` : ''}`)
  },
  
  getJobDetails: async (companyId, jobId) => {
    return await apiService.get(`/public/companies/${companyId}/jobs/${jobId}`)
  },
  
  applyForJob: async (jobId, formData) => {
    // Note: This goes to the portal protected endpoint since we require login to apply
    return await apiService.post(`/portal/jobs/${jobId}/apply`, formData)
  },

  getCompanyDetails: async (companyId) => {
    return await apiService.get(`/public/companies/${companyId}`)
  }
}
