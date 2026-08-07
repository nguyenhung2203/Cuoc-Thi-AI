import { apiService } from './api.service';
import { loadSavedJobs } from '../utils/savedJobs';

export const candidatePortalService = {
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

  deleteCv: (fileId) => {
    return apiService.delete(`/portal/cv/${fileId}`);
  },

  /** Re-run AI CV parse and persist on the file (for portal CVs without applications). */
  reparseCv: () => {
    return apiService.post('/portal/cv/reparse', {});
  },

  /**
   * AI góp ý sửa lỗi / cải thiện CV (không phải chấm phỏng vấn).
   * @returns {Promise<{ summary, issues, suggestions, missing_sections, strengths }>}
   */
  reviewCv: () => {
    return apiService.post('/portal/cv/review', {});
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
    return loadSavedJobs();
  }
};
