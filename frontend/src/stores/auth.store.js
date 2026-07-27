import { reactive } from 'vue';
import { authService } from '../services/auth.service';

/**
 * Auth Store (State Management)
 * Lưu trữ trạng thái đăng nhập toàn cục và token.
 */
export const authStore = reactive({
  user: null, // { id, email, full_name, role, companies }
  isAuthenticated: !!localStorage.getItem('access_token'),
  isLoading: false,
  error: null,

  /**
   * Khởi tạo store: Lấy lại thông tin user nếu có token
   */
  async init() {
    if (this.isAuthenticated) {
      try {
        this.isLoading = true;
        let userData = await authService.getMe();
        
        // Auto create company for recruiter if missing
        if (userData.role === 'recruiter' && (!userData.companies || userData.companies.length === 0)) {
          const { apiService } = await import('../services/api.service');
          await apiService.post('/companies', { 
            name: `Công ty của ${userData.full_name}`,
            website: '',
            industry: 'IT',
            size: '1-50'
          });
          userData = await authService.getMe();
        }
        
        this.user = userData;
        if (userData?.role) {
          localStorage.setItem('user_role', userData.role);
        }
        if (userData?.email) {
          localStorage.setItem('user_email', userData.email);
        }
        if (userData?.id) {
          localStorage.setItem('user_id', userData.id);
        }
      } catch (err) {
        console.error('Failed to init auth store:', err);
        // api.service.js đã tự động đá về /login nếu 401
      } finally {
        this.isLoading = false;
      }
    }
  },

  /**
   * Gọi API đăng nhập và lưu token
   */
  async login(email, password) {
    this.isLoading = true;
    this.error = null;
    try {
      const data = await authService.login({ email, password });
      
      // Lưu token
      localStorage.setItem('access_token', data.access_token);
      if (data.refresh_token) {
        localStorage.setItem('refresh_token', data.refresh_token);
      }
      
      this.user = data.user;
      
      this.isAuthenticated = true;
      
      // Lấy danh sách companies đầy đủ
      let fullUserData = await authService.getMe();
      
      // Auto create company for recruiter if missing
      if (fullUserData.role === 'recruiter' && (!fullUserData.companies || fullUserData.companies.length === 0)) {
        const { apiService } = await import('../services/api.service');
        await apiService.post('/companies', { 
          name: `Công ty của ${fullUserData.full_name}`,
          website: '',
          industry: 'IT',
          size: '1-50'
        });
        fullUserData = await authService.getMe();
      }
      
      this.user = fullUserData;
      localStorage.setItem('user_role', fullUserData.role);
      
      return fullUserData;
    } catch (err) {
      this.error = err.message || 'Đăng nhập thất bại';
      throw err;
    } finally {
      this.isLoading = false;
    }
  },
  
  /**
   * Đăng ký
   */
  async register(email, password, full_name, role, otp) {
    this.isLoading = true;
    this.error = null;
    try {
      const data = await authService.register({ email, password, full_name, role, otp });
      return data;
    } catch (err) {
      this.error = err.message || 'Đăng ký thất bại';
      throw err;
    } finally {
      this.isLoading = false;
    }
  },

  /**
   * Đăng xuất
   */
  async logout() {
    try {
      await authService.logout();
    } catch (err) {
      console.error('Logout error:', err);
    } finally {
      // Dù API thành công hay lỗi, vẫn xóa ở client
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user_role');
      this.user = null;
      this.isAuthenticated = false;
      window.location.href = '/login';
    }
  },

  /**
   * Đăng nhập thật bằng Google OAuth.
   * Yêu cầu id_token thật do Google Identity Services cấp — backend sẽ verify token
   * với Google và lấy email/tên từ token đã xác thực (không tin payload từ client).
   * @param {string} idToken - JWT credential từ Google Identity Services
   * @param {string} role - 'candidate' | 'recruiter'
   */
  async loginWithGoogle(idToken, role = 'candidate') {
    this.isLoading = true;
    this.error = null;
    try {
      if (!idToken) {
        throw new Error('Thiếu Google credential. Vui lòng thử đăng nhập lại.');
      }

      // Gọi xuống API thực tế trên Go Backend (/api/v1/auth/google-login)
      const data = await authService.googleLogin({
        id_token: idToken,
        role: role
      });

      localStorage.setItem('access_token', data.access_token);
      if (data.refresh_token) {
        localStorage.setItem('refresh_token', data.refresh_token);
      }

      this.user = data.user;
      this.isAuthenticated = true;

      let fullUserData = await authService.getMe();
      if (fullUserData.role === 'recruiter' && (!fullUserData.companies || fullUserData.companies.length === 0)) {
        const { apiService } = await import('../services/api.service');
        await apiService.post('/companies', {
          name: `Doanh nghiệp AI (${fullUserData.full_name})`,
          website: '',
          industry: 'Technology & AI',
          size: '50-200'
        });
        fullUserData = await authService.getMe();
      }

      this.user = fullUserData;
      localStorage.setItem('user_role', fullUserData.role);
      return fullUserData;
    } catch (err) {
      this.error = err.message || 'Lỗi xác thực Google OAuth';
      throw err;
    } finally {
      this.isLoading = false;
    }
  },
});
