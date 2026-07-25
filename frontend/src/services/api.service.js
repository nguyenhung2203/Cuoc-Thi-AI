/**
 * API Service 
 * Wrapper xung quanh native fetch API để tự động xử lý Token và Lỗi theo API_SPEC.
 * Không dùng axios để tránh thay đổi package.json gây conflict cho team.
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

// Lấy Token từ LocalStorage
const getToken = () => localStorage.getItem('access_token');

// Hàm wrapper chính
const request = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // Mặc định Headers
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // Gắn Token nếu có
  const token = getToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // Xóa Content-Type nếu gửi FormData (trình duyệt tự tính toán boundary)
  if (options.body instanceof FormData) {
    delete headers['Content-Type'];
  }

  const config = {
    ...options,
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
      // Bắt lỗi 401 Unauthorized -> Đẩy về login (trừ khi đang ở API đăng nhập/đăng ký)
      if (response.status === 401 && !endpoint.includes('/auth/login') && !endpoint.includes('/auth/register')) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login'; 
        return;
      }

      // Bắt lỗi 403 Forbidden hoặc lỗi cấm quyền -> Đẩy về trang /403 cảnh báo
      const errorCode = (data.error?.code || '').toUpperCase();
      if ((response.status === 403 || errorCode === 'FORBIDDEN' || errorCode === 'UNAUTHORIZED_ROLE') && !endpoint.includes('/interviews/join')) {
        window.location.href = '/403?reason=unauthorized&attempted=' + encodeURIComponent(window.location.pathname);
        return;
      }
      if (errorCode === 'ACCOUNT_LOCKED' || errorCode === 'ACCOUNT_BLOCKED') {
        window.location.href = '/403?reason=account_blocked';
        return;
      }
      if (errorCode === 'ACCOUNT_PENDING') {
        window.location.href = '/403?reason=pending_approval';
        return;
      }
      
      // Quăng lỗi ra ngoài để component tự xử lý (hiển thị Toast)
      const errorPayload = data.error || { message: 'Đã xảy ra lỗi không xác định.' };
      
      // Tự động dịch lỗi Validation từ Backend sang Tiếng Việt
      if (errorPayload.message === 'validation failed' && Array.isArray(errorPayload.details)) {
        const translatedDetails = errorPayload.details.map(detail => {
          const d = detail.toLowerCase();
          if (d.includes('description: failed min')) return 'Mô tả công việc (JD) phải dài ít nhất 10 ký tự.';
          if (d.includes('title: failed min')) return 'Tiêu đề công việc phải dài ít nhất 2 ký tự.';
          if (d.includes('fullname: failed min')) return 'Họ tên phải dài ít nhất 2 ký tự.';
          if (d.includes('email: failed email')) return 'Địa chỉ email không đúng định dạng.';
          if (d.includes('jobid: failed required')) return 'Vui lòng chọn Vị trí ứng tuyển.';
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
      throw { message: 'Không thể kết nối đến máy chủ. Vui lòng kiểm tra lại hệ thống backend.' };
    }
    
    throw error;
  }
};

export const apiService = {
  get: (endpoint, options = {}) => request(endpoint, { method: 'GET', ...options }),
  
  getWithMeta: async (endpoint, options = {}) => {
    const url = `${API_BASE_URL}${endpoint}`;
    const token = localStorage.getItem('access_token');
    const headers = { 'Content-Type': 'application/json', ...options.headers };
    if (token) headers['Authorization'] = `Bearer ${token}`;
    
    const response = await fetch(url, { method: 'GET', headers, ...options });
    const data = await response.json();
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
