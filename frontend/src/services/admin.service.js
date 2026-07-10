import { apiService } from './api.service';

/**
 * Admin Service
 * Các API dành cho Super Admin (Quản trị hệ thống)
 */
export const adminService = {
  /**
   * Lấy danh sách System Prompts
   */
  getSystemPrompts: () => {
    return apiService.get(`/admin/ai-prompts`);
  },

  /**
   * Cập nhật System Prompt (Tạo phiên bản mới)
   * @param {Object} data { name, content, model, variables_schema, params }
   */
  updateSystemPrompt: (data) => {
    return apiService.post(`/admin/ai-prompts`, data);
  }
};
