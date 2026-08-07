import { apiService } from './api.service';

export const candidatePortalService = {
  checkApplied: (jobId) => {
    return apiService.get(`/portal/jobs/${jobId}/check-applied`);
  },

  getApplications: () => {
    return apiService.get('/portal/applications');
  },

  getDashboardStats: () => {
    return apiService.get('/portal/dashboard');
  },

  getInterviews: () => {
    return apiService.get('/portal/interviews');
  },

  getProfile: () => {
    return apiService.get('/portal/profile');
  },

  /**
   * Cập nhật thông tin profile candidate
   * @param {Object} data { full_name, avatar_url }
   * @returns {Promise<Object>}
   */
  updateProfile: (data) => {
    return apiService.put('/portal/profile', data);
  },

  uploadCv: (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiService.post('/portal/cv', formData);
  },

  getApplications: () => {
    return apiService.get('/portal/applications');
  },

  /**
   * Lấy độ phù hợp CV↔công việc do AI tính (dùng cho trang ứng tuyển).
   * Trả về { has_cv, fit_score, matched_skills, missing_skills, summary, recommendation }.
   * @param {string} jobID
   * @returns {Promise<Object>}
   */
  getJobMatch: (jobID) => {
    return apiService.get(`/portal/jobs/${jobID}/match`);
  },

  cancelApplication: (id) => {
    return apiService.delete(`/portal/applications/${id}`);
  },

  getSavedJobs: async () => {
    const saved = localStorage.getItem('candidate_saved_jobs');
    return saved ? JSON.parse(saved) : [];
  }
};
