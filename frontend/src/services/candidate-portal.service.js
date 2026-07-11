import { apiService } from './api.service';

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

  getApplications: () => {
    return apiService.get('/portal/applications');
  },

  cancelApplication: (id) => {
    return apiService.delete(`/portal/applications/${id}`);
  },

  getSavedJobs: async () => {
    const saved = localStorage.getItem('candidate_saved_jobs');
    return saved ? JSON.parse(saved) : [];
  }
};
