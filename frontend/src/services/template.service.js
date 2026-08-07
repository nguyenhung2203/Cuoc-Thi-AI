import { apiService } from './api.service';

/**
 * Template Service
 * CRUD Mẫu kịch bản phỏng vấn AI — API_SPEC
 * Endpoint: /companies/:company_id/templates
 */
export const templateService = {
  /**
   * Lấy danh sách Template của company
   * @param {String} companyId
   * @param {Object} params { page, page_size }
   */
  getTemplates: (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiService.get(`/companies/${companyId}/templates${query ? `?${query}` : ''}`);
  },

  /**
   * Lấy chi tiết Template
   * @param {String} companyId
   * @param {String} templateId
   */
  getTemplate: (companyId, templateId) => {
    return apiService.get(`/companies/${companyId}/templates/${templateId}`);
  },

  /**
   * Tạo Template mới
   * @param {String} companyId
   * @param {Object} data { name, type, duration_minutes, description, config }
   */
  createTemplate: (companyId, data) => {
    return apiService.post(`/companies/${companyId}/templates`, data);
  },

  /**
   * Cập nhật Template
   * @param {String} companyId
   * @param {String} templateId
   * @param {Object} data
   */
  updateTemplate: (companyId, templateId, data) => {
    return apiService.put(`/companies/${companyId}/templates/${templateId}`, data);
  },

  /**
   * Xóa Template
   * @param {String} companyId
   * @param {String} templateId
   */
  deleteTemplate: (companyId, templateId) => {
    return apiService.delete(`/companies/${companyId}/templates/${templateId}`);
  }
};
