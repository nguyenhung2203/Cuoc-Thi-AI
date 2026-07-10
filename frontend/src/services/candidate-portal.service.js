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

  uploadCv: (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiService.post('/portal/cv', formData);
  }
};
