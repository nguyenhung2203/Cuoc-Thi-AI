import { apiService } from './api.service';

/**
 * Notification Service
 * Các API liên quan đến quản lý Thông báo
 * Quản lý thông báo người dùng
 */
export const notificationService = {
  /**
   * Lấy danh sách thông báo của user hiện tại
   * @param {Object} params { page, page_size, is_read }
   */
  getNotifications: (params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiService.get(`/notifications${query ? `?${query}` : ''}`);
  },

  /**
   * Đánh dấu một thông báo đã đọc
   * @param {String} notificationId
   */
  markAsRead: (notificationId) => {
    return apiService.put(`/notifications/${notificationId}/read`);
  },

  /**
   * Đánh dấu tất cả thông báo đã đọc
   */
  markAllAsRead: () => {
    return apiService.put('/notifications/read-all');
  }
};

/**
 * Audit Log Service
 * Xem lịch sử hoạt động (chỉ Admin) — API_SPEC §13
 */
export const auditService = {
  /**
   * Lấy danh sách audit log của company
   * API_SPEC §13.1 — GET /companies/:company_id/audit-logs
   * @param {String} companyId
   * @param {Object} params { page, page_size, resource_type }
   */
  getAuditLogs: (companyId, params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiService.get(`/companies/${companyId}/audit-logs${query ? `?${query}` : ''}`);
  },

  /**
   * Lấy danh sách audit log của toàn hệ thống (Admin)
   */
  getSystemAuditLogs: (params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiService.get(`/admin/audit-logs${query ? `?${query}` : ''}`);
  }
};
