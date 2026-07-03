import { apiService } from './api.service';

/**
 * Job Service
 * Các API liên quan đến quản lý Job
 */
export const jobService = {
  /**
   * Lấy danh sách Job
   * @param {String} companyId 
   * @param {Object} params { status, keyword, page, page_size }
   */
  getJobs: (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString();
    const url = `/companies/${companyId}/jobs${query ? `?${query}` : ''}`;
    return apiService.get(url);
  },

  /**
   * Tạo Job mới
   * @param {String} companyId 
   * @param {Object} data 
   */
  createJob: (companyId, data) => {
    return apiService.post(`/companies/${companyId}/jobs`, data);
  },

  /**
   * Lấy chi tiết Job
   * @param {String} companyId 
   * @param {String} jobId 
   */
  getJob: (companyId, jobId) => {
    return apiService.get(`/companies/${companyId}/jobs/${jobId}`);
  },

  /**
   * Cập nhật Job
   * @param {String} companyId 
   * @param {String} jobId 
   * @param {Object} data 
   */
  updateJob: (companyId, jobId, data) => {
    return apiService.put(`/companies/${companyId}/jobs/${jobId}`, data);
  },

  /**
   * Xóa Job
   * @param {String} companyId 
   * @param {String} jobId 
   */
  deleteJob: (companyId, jobId) => {
    return apiService.delete(`/companies/${companyId}/jobs/${jobId}`);
  },

  /**
   * Phân tích JD bằng AI
   * @param {String} companyId 
   * @param {String} jobId 
   * @param {Boolean} forceRefresh 
   */
  analyzeJD: (companyId, jobId, forceRefresh = false) => {
    return apiService.post(`/companies/${companyId}/jobs/${jobId}/analyze`, { force_refresh: forceRefresh });
  }
};
