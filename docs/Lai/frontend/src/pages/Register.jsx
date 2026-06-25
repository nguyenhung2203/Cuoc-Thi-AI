import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Input } from '../components/Input';
import { Button } from '../components/Button';

export function Register() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('candidate');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleRegister = (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    setTimeout(() => {
      setLoading(false);
      // Giả lập đăng ký thành công
      navigate('/login', { state: { message: 'Đăng ký thành công! Vui lòng đăng nhập.' } });
    }, 1000);
  };

  return (
    <div className="auth-container">
      <div style={{ display: 'flex', flexDirection: 'column', gap: '32px', width: '100%', maxWidth: '400px' }}>
        <div style={{ textAlign: 'center' }}>
          <h1 className="text-h1" style={{ marginBottom: '8px', color: 'var(--primary)' }}>Interview AI</h1>
          <p className="text-body" style={{ color: 'var(--text-secondary)' }}>Tạo tài khoản mới</p>
        </div>
        
        <Card>
          <form onSubmit={handleRegister}>
            <Input 
              label="Họ và tên" 
              type="text" 
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              placeholder="Nguyễn Văn A"
            />
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

            <div style={{ marginBottom: '16px' }}>
              <label className="text-helper" style={{ display: 'block', marginBottom: '8px', fontWeight: 500, color: 'var(--text-main)' }}>
                Bạn là:
              </label>
              <div style={{ display: 'flex', gap: '16px' }}>
                <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                  <input 
                    type="radio" 
                    name="role" 
                    value="candidate" 
                    checked={role === 'candidate'} 
                    onChange={(e) => setRole(e.target.value)} 
                  />
                  <span className="text-body">Ứng viên</span>
                </label>
                <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                  <input 
                    type="radio" 
                    name="role" 
                    value="recruiter" 
                    checked={role === 'recruiter'} 
                    onChange={(e) => setRole(e.target.value)} 
                  />
                  <span className="text-body">Nhà tuyển dụng</span>
                </label>
              </div>
            </div>
            
            {error && <div className="error-text" style={{ marginBottom: '16px' }}>{error}</div>}
            
            <Button type="submit" style={{ width: '100%', marginTop: '16px' }} disabled={loading}>
              {loading ? 'Đang xử lý...' : 'Đăng ký'}
            </Button>
          </form>
          
          <div style={{ marginTop: '24px', textAlign: 'center' }}>
            <p className="text-body" style={{ color: 'var(--text-secondary)' }}>
              Đã có tài khoản? <span style={{ color: 'var(--primary)', cursor: 'pointer', fontWeight: 500 }} onClick={() => navigate('/login')}>Đăng nhập</span>
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
}
