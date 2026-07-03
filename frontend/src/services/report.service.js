import { apiService } from './api.service';

/**
 * Report Service
 * Quản lý Báo cáo AI sau phỏng vấn
 */
export const reportService = {
  /**
   * Lấy báo cáo phân tích AI của buổi phỏng vấn
   * @param {String} companyId 
   * @param {String} interviewId 
   */
  getReport: async (companyId, interviewId) => {
    return apiService.get(`/companies/${companyId}/interviews/${interviewId}/report`);
  },

  /**
   * Lưu quyết định của Recruiter
   */
  saveDecision: async (companyId, interviewId, data) => {
    return apiService.put(`/companies/${companyId}/interviews/${interviewId}/report/decision`, data);
  }
};
