import { apiService } from './api.service';

/**
 * Company Service
 * Các API liên quan đến quản lý Company
 */
export const companyService = {
  /**
   * Lấy danh sách companies của user hiện tại
   */
  getCompanies: () => {
    return apiService.get('/companies');
  },

  /**
   * Tạo company mới
   * @param {Object} data { name, website, industry, size }
   */
  createCompany: (data) => {
    return apiService.post('/companies', data);
  },

  /**
   * Lấy chi tiết company
   * @param {String} companyId 
   */
  getCompany: (companyId) => {
    return apiService.get(`/companies/${companyId}`);
  },

  /**
   * Cập nhật thông tin company
   * @param {String} companyId 
   * @param {Object} data 
   */
  updateCompany: (companyId, data) => {
    return apiService.put(`/companies/${companyId}`, data);
  }
};
