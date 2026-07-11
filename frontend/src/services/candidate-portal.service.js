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

  // Mock functions for missing backend endpoints
  getApplications: () => {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([
          {
            id: 'app-1',
            job_id: 'job-1',
            job_title: 'Senior Backend Engineer (Golang)',
            company_name: 'TechCorp VN',
            applied_at: new Date(Date.now() - 86400000 * 2).toISOString(),
            status: 'interviewing', // pending, reviewed, interviewing, rejected
            cv_name: 'NguyenHung_CV_Backend.pdf'
          },
          {
            id: 'app-2',
            job_id: 'job-2',
            job_title: 'Fullstack Developer (Vue/Node)',
            company_name: 'Innovate AI',
            applied_at: new Date(Date.now() - 86400000 * 5).toISOString(),
            status: 'reviewed',
            cv_name: 'NguyenHung_CV_Fullstack.pdf'
          },
          {
            id: 'app-3',
            job_id: 'job-3',
            job_title: 'Software Engineer',
            company_name: 'Global Corp',
            applied_at: new Date(Date.now() - 86400000 * 10).toISOString(),
            status: 'rejected',
            cv_name: 'NguyenHung_CV.pdf'
          }
        ]);
      }, 500);
    });
  },

  getSavedJobs: () => {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([
          {
            id: 'job-4',
            company_id: 'mock-c-1',
            title: 'AI Prompt Engineer',
            company_name: 'TechCorp VN',
            location: 'Hà Nội',
            salary_min: { Valid: true, Int64: 1500 },
            salary_max: { Valid: true, Int64: 2500 },
            currency: { Valid: true, String: 'USD' },
            saved_at: new Date(Date.now() - 86400000).toISOString()
          },
          {
            id: 'job-5',
            company_id: 'mock-c-2',
            title: 'Product Designer (UI/UX)',
            company_name: 'DesignStudio',
            location: 'TP. HCM',
            salary_min: { Valid: false, Int64: 0 },
            salary_max: { Valid: false, Int64: 0 },
            currency: { Valid: true, String: 'VND' },
            saved_at: new Date(Date.now() - 86400000 * 3).toISOString()
          }
        ]);
      }, 500);
    });
  }
};
