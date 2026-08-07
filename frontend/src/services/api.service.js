/**
 * API Service 
 * Wrapper xung quanh native fetch API để tự động xử lý Token và Lỗi theo API_SPEC.
 * Không dùng axios để tránh thay đổi package.json gây conflict cho team.
 */

import { normalizeValidationErrors } from '../utils/validators.js';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

// Lấy Token từ LocalStorage
const getToken = () => localStorage.getItem('access_token');
let refreshPromise = null;

const isAuthEndpoint = (endpoint) => endpoint.includes('/auth/login')
  || endpoint.includes('/auth/register')
  || endpoint.includes('/auth/refresh');

const clearSessionAndRedirect = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('user_role');
  window.location.href = '/login';
};

const refreshAccessToken = async () => {
  if (!refreshPromise) {
    refreshPromise = fetch(`${API_BASE_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: 'null',
    })
      .then(async response => {
        const data = await response.json();
        if (!response.ok || data.success === false || !data.data?.access_token) {
          throw data.error || { message: 'Phiên đăng nhập đã hết hạn.' };
        }
        localStorage.setItem('access_token', data.data.access_token);
        return data.data.access_token;
      })
      .finally(() => { refreshPromise = null; });
  }
  return refreshPromise;
};

// Hàm wrapper chính
const request = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const { skipAuthRefresh = false, ...requestOptions } = options;

  // Mặc định Headers
  const headers = {
    'Content-Type': 'application/json',
    ...requestOptions.headers,
  };

  // Gắn Token nếu có
  const token = getToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // Xóa Content-Type nếu gửi FormData (trình duyệt tự tính toán boundary)
  if (requestOptions.body instanceof FormData) {
    delete headers['Content-Type'];
  }

  const config = {
    ...requestOptions,
    credentials: 'include',
    headers,
  };

  try {
    const response = await fetch(url, config);
    
    // Nếu API trả về 204 No Content
    if (response.status === 204) {
      return { success: true, data: null };
    }

    const data = await response.json();

    // Xử lý lỗi HTTP hoặc lỗi từ cấu trúc trả về (success: false)
    if (!response.ok || data.success === false) {
      const errorCode = (data.error?.code || '').toUpperCase();

      // Tài khoản bị khoá — phải bắt TRƯỚC nhánh 403 chung để hiện đúng
      // trang /403?reason=account_blocked (backend trả 403 ACCOUNT_BLOCKED).
      if (errorCode === 'ACCOUNT_LOCKED' || errorCode === 'ACCOUNT_BLOCKED') {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_role');
        window.location.href = '/403?reason=account_blocked';
        // Vẫn throw để caller (vd authStore.login) không đọc data.access_token
        // của response lỗi trong lúc trang đang điều hướng.
        throw data.error || { message: 'Tài khoản đã bị khoá.' };
      }
      if (errorCode === 'ACCOUNT_PENDING') {
        window.location.href = '/403?reason=pending_approval';
        throw data.error || { message: 'Tài khoản đang chờ phê duyệt.' };
      }

      // Access token hết hạn: refresh một lần rồi retry request gốc.
      if (response.status === 401 && !skipAuthRefresh && !isAuthEndpoint(endpoint)) {
        try {
          const newToken = await refreshAccessToken();
          const retryHeaders = { ...headers, Authorization: `Bearer ${newToken}` };
          const retryResponse = await fetch(url, { ...config, headers: retryHeaders });
          const retryData = retryResponse.status === 204 ? null : await retryResponse.json();
          if (!retryResponse.ok || retryData?.success === false) {
            throw retryData?.error || { message: 'Phiên đăng nhập đã hết hạn.' };
          }
          return retryData?.data;
        } catch (refreshError) {
          clearSessionAndRedirect();
          throw refreshError;
        }
      }

      // Bắt lỗi 401 Unauthorized -> Đẩy về login (trừ khi đang ở API đăng nhập/đăng ký)
      if (response.status === 401 && !endpoint.includes('/auth/login') && !endpoint.includes('/auth/register')) {
        clearSessionAndRedirect();
        return;
      }

      // Bắt lỗi 403 Forbidden hoặc lỗi cấm quyền -> Đẩy về trang /403 cảnh báo
      if ((response.status === 403 || errorCode === 'FORBIDDEN' || errorCode === 'UNAUTHORIZED_ROLE') && !endpoint.includes('/interviews/join')) {
        window.location.href = '/403?reason=unauthorized&attempted=' + encodeURIComponent(window.location.pathname);
        return;
      }

      // Quăng lỗi ra ngoài để component tự xử lý (hiển thị Toast)
      const errorPayload = data.error || { message: 'Đã xảy ra lỗi không xác định.' };
      const errorMessages = {
        INVALID_CURRENT_PASSWORD: 'Mật khẩu hiện tại không chính xác. Vui lòng kiểm tra và thử lại.',
        INVALID_CREDENTIALS: 'Email hoặc mật khẩu không chính xác. Vui lòng kiểm tra lại.',
        EMAIL_ALREADY_EXISTS: 'Email này đã được sử dụng.',
        INVALID_OTP: 'Mã OTP không chính xác.',
        OTP_EXPIRED: 'Mã OTP đã hết hạn. Vui lòng yêu cầu mã mới.',
        FILE_TOO_LARGE: 'Tệp vượt quá dung lượng cho phép.',
        UNSUPPORTED_FILE_TYPE: 'Định dạng tệp không được hỗ trợ.',
        ALREADY_APPLIED: 'Bạn đã ứng tuyển công việc này.',
        JOB_CLOSED: 'Công việc này đã đóng và không còn nhận hồ sơ.',
        INTERVIEW_CONFLICT: 'Thời gian phỏng vấn bị trùng với một lịch đã có.',
      };
      if (errorMessages[errorCode]) errorPayload.message = errorMessages[errorCode];
      if (/current password is incorrect/i.test(errorPayload.message || '')) {
        errorPayload.message = errorMessages.INVALID_CURRENT_PASSWORD;
      }
      errorPayload.validationErrors = normalizeValidationErrors(errorPayload);

      // Tự động dịch lỗi Validation từ Backend sang Tiếng Việt
      if (errorPayload.message === 'validation failed' && Array.isArray(errorPayload.details)) {
        const translatedDetails = errorPayload.details.map(detail => {
          const d = detail.toLowerCase();
          if (d.includes('description: failed min')) return 'Mô tả công việc (JD) phải dài ít nhất 10 ký tự.';
          if (d.includes('title: failed min')) return 'Tiêu đề công việc phải dài ít nhất 2 ký tự.';
          if (d.includes('fullname: failed min')) return 'Họ tên phải dài ít nhất 2 ký tự.';
          if (d.includes('fullname: failed required')) return 'Vui lòng nhập họ tên.';
          if (d.includes('email: failed email')) return 'Địa chỉ email không đúng định dạng.';
          if (d.includes('email: failed required')) return 'Vui lòng nhập địa chỉ email.';
          if (d.includes('password: failed min')) return 'Mật khẩu phải dài ít nhất 6 ký tự.';
          if (d.includes('password: failed required')) return 'Vui lòng nhập mật khẩu.';
          if (d.includes('name: failed required')) return 'Vui lòng nhập tên.';
          if (d.includes('jobid: failed required')) return 'Vui lòng chọn Vị trí ứng tuyển.';
          if (d.includes('current password is incorrect')) return 'Mật khẩu hiện tại không chính xác.';
          return 'Dữ liệu nhập vào chưa hợp lệ: ' + detail;
        });
        errorPayload.message = translatedDetails.join(' ');
      }
      
      throw errorPayload;
    }

    // Trả về data (bóc vỏ success)
    return data.data;

  } catch (error) {
    console.error(`[API Error] ${options.method || 'GET'} ${endpoint}:`, error);
    
    // Xử lý lỗi Network (không kết nối được tới server)
    if (error instanceof TypeError && error.message.includes('Failed to fetch')) {
      throw { message: 'Không thể kết nối đến máy chủ.' };
    }
    
    throw error;
  }
};

export const apiService = {
  get: (endpoint, options = {}) => request(endpoint, { method: 'GET', ...options }),
  
  getWithMeta: async (endpoint, options = {}) => {
    const url = `${API_BASE_URL}${endpoint}`;
    const { skipAuthRefresh = false, ...requestOptions } = options;
    const token = localStorage.getItem('access_token');
    const headers = { 'Content-Type': 'application/json', ...requestOptions.headers };
    if (token) headers['Authorization'] = `Bearer ${token}`;
    const fetchOptions = {
      method: 'GET',
      ...requestOptions,
      credentials: 'include',
      headers,
    };

    const response = await fetch(url, fetchOptions);
    const data = await response.json();
    if (response.status === 401 && !skipAuthRefresh && !isAuthEndpoint(endpoint)) {
      try {
        const newToken = await refreshAccessToken();
        const retryHeaders = { ...headers, Authorization: `Bearer ${newToken}` };
        const retryResponse = await fetch(url, { ...fetchOptions, headers: retryHeaders, credentials: 'include' });
        const retryData = await retryResponse.json();
        if (!retryResponse.ok || retryData.success === false) {
          throw retryData.error || { message: 'Phiên đăng nhập đã hết hạn.' };
        }
        return retryData;
      } catch (refreshError) {
        clearSessionAndRedirect();
        throw refreshError;
      }
    }
    if (!response.ok || data.success === false) throw data.error || { message: 'Lỗi' };
    return data; // Returns { success, data, meta }
  },
  
  post: (endpoint, body, options = {}) => request(endpoint, { 
    method: 'POST', 
    body: body instanceof FormData ? body : JSON.stringify(body), 
    ...options 
  }),
  
  put: (endpoint, body, options = {}) => request(endpoint, { 
    method: 'PUT', 
    body: body instanceof FormData ? body : JSON.stringify(body), 
    ...options 
  }),
  
  delete: (endpoint, options = {}) => request(endpoint, { method: 'DELETE', ...options }),
  
  uploadFile: (endpoint, file, fieldName = 'file') => {
    const formData = new FormData();
    formData.append(fieldName, file);
    return request(endpoint, {
      method: 'POST',
      body: formData
    });
  }
};
