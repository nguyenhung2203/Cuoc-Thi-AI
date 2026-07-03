import { apiService } from './api.service';

/**
 * Notification Service
 * Các API liên quan đến quản lý Thông báo
 */
export const notificationService = {
  /**
   * Lấy danh sách thông báo
   */
  getNotifications: () => {
    return apiService.get(`/notifications`);
  },

  /**
   * Đánh dấu một thông báo là đã đọc
   * @param {String} notificationId 
   */
  markAsRead: (notificationId) => {
    return apiService.put(`/notifications/${notificationId}/read`);
  },

  /**
   * Đánh dấu tất cả thông báo là đã đọc
   */
  markAllAsRead: () => {
    return apiService.put(`/notifications/read-all`);
  }
};
