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

  uploadCv: (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiService.post('/portal/cv', formData);
  }
};
