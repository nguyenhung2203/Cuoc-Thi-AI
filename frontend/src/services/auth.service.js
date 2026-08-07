import { apiService } from './api.service';

/**
 * Auth Service
 * Gọi trực tiếp xuống các endpoint Authentication (Mục 2 của API_SPEC.md)
 */
export const authService = {
  /**
   * Đăng nhập
   * @param {Object} credentials - { email, password }
   * @returns {Promise<Object>} { user, access_token }; refresh token is an HttpOnly cookie
   */
  login: (credentials) => {
    return apiService.post('/auth/login', credentials);
  },

  /**
   * Đăng ký (payload cần kèm otp đã xác thực qua email)
   * @param {Object} payload - { email, password, full_name, role, otp }
   * @returns {Promise<Object>} { user, access_token }; refresh token is an HttpOnly cookie
   */
  register: (payload) => {
    return apiService.post('/auth/register', payload);
  },

  /**
   * Gửi mã OTP xác nhận đăng ký về email
   * @param {String} email
   */
  sendRegisterOtp: (email) => {
    return apiService.post('/auth/register/send-otp', { email });
  },

  /**
   * Quên mật khẩu — gửi OTP đặt lại về email
   * @param {String} email
   */
  forgotPassword: (email) => {
    return apiService.post('/auth/forgot-password', { email });
  },

  /**
   * Xác thực OTP đặt lại mật khẩu (chưa tiêu thụ mã)
   * @param {Object} data { email, otp }
   */
  verifyResetOtp: (data) => {
    return apiService.post('/auth/verify-reset-otp', data);
  },

  /**
   * Đặt lại mật khẩu với OTP đã xác thực
   * @param {Object} data { email, otp, new_password }
   */
  resetPassword: (data) => {
    return apiService.post('/auth/reset-password', data);
  },

  /**
   * Đăng nhập thật bằng Google OAuth token
   * @param {Object} payload - { id_token, email, full_name, avatar, role }
   */
  googleLogin: (payload) => {
    return apiService.post('/auth/google-login', payload);
  },

  /**
   * Lấy thông tin user hiện tại (Me)
   * @returns {Promise<Object>} { id, email, full_name, role, companies }
   */
  getMe: () => {
    return apiService.get('/auth/me');
  },

  /**
   * Làm mới access token bằng HttpOnly refresh-token cookie.
   */
  refresh: () => {
    return apiService.post('/auth/refresh', null, { skipAuthRefresh: true });
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
  },

  /**
   * Xác thực tài liệu
   * @param {string} fileId
   * @returns {Promise<Object>}
   */
  verifyDocument: (fileId) => {
    return apiService.post('/auth/me/verify-document', { file_id: fileId });
  },

  /**
   * Đổi mật khẩu
   * @param {Object} data { current_password, new_password }
   * @returns {Promise<Object>}
   */
  changePassword: (data) => {
    return apiService.put('/auth/me/password', data);
  },

  /**
   * Xóa tài khoản (soft-delete) của chính mình
   * @returns {Promise<Object>}
   */
  deleteAccount: () => {
    return apiService.delete('/auth/me');
  }
};
