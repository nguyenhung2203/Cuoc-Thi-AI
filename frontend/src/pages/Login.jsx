import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Card } from '../components/Card';
import { Input } from '../components/Input';
import { Button } from '../components/Button';
import { Toast } from '../components/Toast';

export function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [entryToast, setEntryToast] = useState(location.state?.message ? { type: 'success', message: location.state.message } : null);

  useEffect(() => {
    if (location.state?.message) {
      window.history.replaceState({}, document.title);
    }
  }, [location]);

  const handleLogin = (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    setTimeout(() => {
      setLoading(false);
      localStorage.setItem('token', 'mock-token');
      
      const emailLower = email.toLowerCase();
      if (emailLower.includes('hr') || emailLower.includes('recruiter') || emailLower.includes('admin')) {
        localStorage.setItem('role', 'recruiter');
        navigate('/dashboard', { state: { message: 'Đăng nhập thành công với quyền Nhà tuyển dụng!' } });
      } else {
        localStorage.setItem('role', 'candidate');
        navigate('/home', { state: { message: 'Đăng nhập thành công với quyền Ứng viên!' } });
      }
    }, 1000);
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: 'var(--background)' }}>
      {entryToast && (
        <Toast type={entryToast.type} message={entryToast.message} onClose={() => setEntryToast(null)} />
      )}
      <div style={{ display: 'flex', flexDirection: 'column', gap: '32px', width: '100%', maxWidth: '400px' }}>
        <div style={{ textAlign: 'center' }}>
          <h1 className="text-h1" style={{ marginBottom: '8px', color: 'var(--primary)' }}>Interview AI</h1>
          <p className="text-body" style={{ color: 'var(--text-secondary)' }}>Đăng nhập vào hệ thống tuyển dụng</p>
        </div>
        
        <Card>
          <form onSubmit={handleLogin}>
            <Input 
              label="Email" 
              type="email" 
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              placeholder="nhapemail@congty.com"
            />
            <Input 
              label="Mật khẩu" 
              type="password" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="••••••••"
            />
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '-8px', marginBottom: '16px' }}>
              <span onClick={() => navigate('/forgot-password')} style={{ fontSize: '13px', color: 'var(--primary)', fontWeight: 500, cursor: 'pointer' }}>Quên mật khẩu?</span>
            </div>
            
            {error && <div className="error-text" style={{ marginBottom: '16px' }}>{error}</div>}
            
            <Button type="submit" style={{ width: '100%', marginTop: '16px' }} disabled={loading}>
              {loading ? 'Đang xử lý...' : 'Đăng nhập'}
            </Button>
          </form>
          
          <div style={{ marginTop: '24px', textAlign: 'center' }}>
            <p className="text-body" style={{ color: 'var(--text-secondary)' }}>
              Chưa có tài khoản? <span style={{ color: 'var(--primary)', cursor: 'pointer', fontWeight: 500 }} onClick={() => navigate('/register')}>Đăng ký ngay</span>
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
}
