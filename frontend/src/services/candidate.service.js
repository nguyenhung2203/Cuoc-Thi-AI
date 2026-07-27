import { apiService } from './api.service';

/**
 * Candidate Service
 * Các API liên quan đến quản lý Ứng viên (Candidate)
 */
export const candidateService = {
  /**
   * Lấy danh sách Candidate
   * @param {String} companyId 
   * @param {Object} params { job_id, status, keyword, page, page_size }
   */
  getCandidates: (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString();
    const url = `/companies/${companyId}/candidates${query ? `?${query}` : ''}`;
    return apiService.get(url);
  },

  /**
   * Thêm Candidate mới
   * @param {String} companyId 
   * @param {Object} data { full_name, email, phone, source, job_id }
   */
  createCandidate: (companyId, data) => {
    return apiService.post(`/companies/${companyId}/candidates`, data);
  },

  /**
   * Lấy chi tiết Candidate
   * @param {String} companyId 
   * @param {String} candidateId 
   */
  getCandidate: (companyId, candidateId) => {
    return apiService.get(`/companies/${companyId}/candidates/${candidateId}`);
  },

  /**
   * Cập nhật Candidate
   * @param {String} companyId 
   * @param {String} candidateId 
   * @param {Object} data 
   */
  updateCandidate: (companyId, candidateId, data) => {
    return apiService.put(`/companies/${companyId}/candidates/${candidateId}`, data);
  },

  /**
   * Upload CV cho Candidate
   * Dùng API File Upload trước, sau đó lưu lại file_id
   * (Hoặc gọi API upload riêng của candidate tùy backend)
   * Ở đây theo API_SPEC.md 5.5
   * @param {String} companyId 
   * @param {String} candidateId 
   * @param {File} file 
   */
  uploadCV: (companyId, candidateId, file) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiService.post(`/companies/${companyId}/candidates/${candidateId}/cv`, formData);
  },

  /**
   * Kích hoạt AI phân tích CV
   * @param {String} companyId 
   * @param {String} candidateId 
   */
  parseCV: (companyId, candidateId) => {
    return apiService.post(`/companies/${companyId}/candidates/${candidateId}/parse-cv`);
  },

  /**
   * Lấy link file tải về hoặc xem chi tiết
   * @param {String} fileId 
   * @param {String} companyId 
   */
  getCVUrl: (fileId, companyId) => {
    return apiService.get(`/files/${fileId}/signed-url?company_id=${companyId}`);
  }
};
