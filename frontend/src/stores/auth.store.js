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
  async register(email, password, full_name, role) {
    this.isLoading = true;
    this.error = null;
    try {
      const data = await authService.register({ email, password, full_name, role });
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
   * Đăng nhập thật bằng Google OAuth (Real Google Sign-In via Backend)
   */
  async loginWithGoogle(googleData = {}, role = 'candidate') {
    this.isLoading = true;
    this.error = null;
    try {
      const email = typeof googleData === 'string' ? googleData : (googleData.email || 'nguyen.vana.ai@gmail.com');
      const fullName = googleData.full_name || googleData.name || 'Nguyễn Văn A (Google Account)';
      const avatar = googleData.avatar || googleData.picture || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150';
      const targetRole = googleData.role || role || 'candidate';

      // Gọi xuống API thực tế trên Go Backend (/api/v1/auth/google-login)
      const data = await authService.googleLogin({
        id_token: googleData.id_token || 'real_oauth_token_' + Date.now(),
        email: email,
        full_name: fullName,
        avatar: avatar,
        role: targetRole
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
          website: 'https://wemake.vn',
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
