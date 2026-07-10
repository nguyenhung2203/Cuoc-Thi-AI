import { apiService } from './api.service';

/**
 * Auth Service
 * Gọi trực tiếp xuống các endpoint Authentication (Mục 2 của API_SPEC.md)
 */
export const authService = {
  /**
   * Đăng nhập
   * @param {Object} credentials - { email, password }
   * @returns {Promise<Object>} { user, access_token, refresh_token }
   */
  login: (credentials) => {
    return apiService.post('/auth/login', credentials);
  },

  /**
   * Đăng ký
   * @param {Object} payload - { email, password, full_name, role }
   * @returns {Promise<Object>} { user, access_token, refresh_token }
   */
  register: (payload) => {
    return apiService.post('/auth/register', payload);
  },

  /**
   * Lấy thông tin user hiện tại (Me)
   * @returns {Promise<Object>} { id, email, full_name, role, companies }
   */
  getMe: () => {
    return apiService.get('/auth/me');
  },

  /**
   * Đăng xuất (Revoke token trên server)
   * @returns {Promise<void>}
   */
  logout: () => {
    return apiService.post('/auth/logout');
  },

  /**
   * Cập nhật thông tin profile user
   * @param {Object} data { full_name, avatar_url }
   * @returns {Promise<Object>}
   */
  updateProfile: (data) => {
    return apiService.put('/auth/me', data);
  },

  /**
   * Lấy cài đặt user hiện tại
   * @returns {Promise<Object>} { settings: {...} }
   */
  getSettings: () => {
    return apiService.get('/auth/me/settings');
  },

  /**
   * Lưu cài đặt user
   * @param {Object} settings
   * @returns {Promise<Object>}
   */
  saveSettings: (settings) => {
    return apiService.put('/auth/me/settings', { settings });
  }
};
