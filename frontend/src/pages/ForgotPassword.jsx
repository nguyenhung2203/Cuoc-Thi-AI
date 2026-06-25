import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Input } from '../components/Input';
import { Button } from '../components/Button';
import { Toast } from '../components/Toast';
import { Mail, KeyRound, Lock, ArrowLeft } from 'lucide-react';

export function ForgotPassword() {
  const navigate = useNavigate();
  const [step, setStep] = useState(1);
  const [email, setEmail] = useState('');
  const [otp, setOtp] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [toast, setToast] = useState(null);
  
  // Timer state
  const [timeLeft, setTimeLeft] = useState(60);
  const [canResend, setCanResend] = useState(false);

  useEffect(() => {
    let timer;
    if (step === 2 && timeLeft > 0) {
      timer = setInterval(() => {
        setTimeLeft(prev => prev - 1);
      }, 1000);
    } else if (timeLeft === 0) {
      setCanResend(true);
    }
    return () => clearInterval(timer);
  }, [step, timeLeft]);

  const handleSendEmail = (e) => {
    e.preventDefault();
    if (!email) return;
    setLoading(true);
    // Giả lập API call
    setTimeout(() => {
      setLoading(false);
      setStep(2);
      setTimeLeft(60);
      setCanResend(false);
      setToast({ type: 'success', message: 'Mã xác nhận đã được gửi đến email của bạn!' });
    }, 1000);
  };

  const handleResendOTP = () => {
    if (!canResend) return;
    setToast({ type: 'success', message: 'Đã gửi lại mã xác nhận mới!' });
    setTimeLeft(60);
    setCanResend(false);
  };

  const handleVerifyOTP = (e) => {
    e.preventDefault();
    if (!otp) return;
    if (otp !== '123456') {
      setToast({ type: 'error', message: 'Mã xác nhận không hợp lệ. Vui lòng thử lại!' });
      return;
    }
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      setStep(3);
    }, 800);
  };

  const handleResetPassword = (e) => {
    e.preventDefault();
    if (!password || !confirmPassword) return;
    if (password !== confirmPassword) {
      setToast({ type: 'error', message: 'Mật khẩu xác nhận không khớp!' });
      return;
    }
    if (password.length < 6) {
      setToast({ type: 'error', message: 'Mật khẩu phải có ít nhất 6 ký tự!' });
      return;
    }
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      navigate('/login', { state: { message: 'Đổi mật khẩu thành công! Vui lòng đăng nhập lại.' } });
    }, 1000);
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: 'var(--background)' }}>
      {toast && (
        <Toast type={toast.type} message={toast.message} onClose={() => setToast(null)} />
      )}
      
      <div style={{ position: 'absolute', top: '32px', left: '32px' }}>
        <Button variant="ghost" onClick={() => navigate('/login')}>
          <ArrowLeft size={18} style={{ marginRight: '8px' }} /> Quay lại Đăng nhập
        </Button>
      </div>

      <Card style={{ width: '100%', maxWidth: '440px', padding: '40px', position: 'relative', overflow: 'hidden' }}>
        
        {/* STEP 1: Nhập Email */}
        {step === 1 && (
          <div style={{ animation: 'fadeIn 0.4s ease-out' }}>
            <div style={{ textAlign: 'center', marginBottom: '32px' }}>
              <div style={{ width: '56px', height: '56px', backgroundColor: 'rgba(37, 99, 235, 0.1)', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px', color: 'var(--primary)' }}>
                <KeyRound size={28} />
              </div>
              <h1 className="text-h1" style={{ marginBottom: '8px' }}>Quên mật khẩu?</h1>
              <p className="text-helper" style={{ lineHeight: 1.5 }}>Nhập địa chỉ email được liên kết với tài khoản của bạn để nhận mã xác nhận.</p>
            </div>

            <form onSubmit={handleSendEmail} style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
              <Input 
                label="Địa chỉ Email" 
                type="email" 
                placeholder="Ví dụ: candidate@test.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              <Button type="submit" disabled={!email || loading} style={{ height: '44px' }}>
                {loading ? 'Đang gửi...' : 'Gửi mã xác nhận'}
              </Button>
            </form>
          </div>
        )}

        {/* STEP 2: Xác thực OTP */}
        {step === 2 && (
          <div style={{ animation: 'fadeIn 0.4s ease-out' }}>
            <div style={{ textAlign: 'center', marginBottom: '32px' }}>
              <div style={{ width: '56px', height: '56px', backgroundColor: 'rgba(8, 145, 178, 0.1)', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px', color: 'var(--accent)' }}>
                <Mail size={28} />
              </div>
              <h1 className="text-h1" style={{ marginBottom: '8px' }}>Nhập mã xác nhận</h1>
              <p className="text-helper" style={{ lineHeight: 1.5 }}>
                Chúng tôi đã gửi mã xác nhận 6 chữ số tới email<br/>
                <strong style={{ color: 'var(--text-main)' }}>{email}</strong>
              </p>
            </div>

            <form onSubmit={handleVerifyOTP} style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
              <div>
                <Input 
                  placeholder="Nhập mã 6 chữ số (VD: 123456)" 
                  value={otp}
                  onChange={(e) => setOtp(e.target.value.replace(/[^0-9]/g, '').slice(0, 6))}
                  style={{ textAlign: 'center', fontSize: '20px', letterSpacing: '4px', fontWeight: 600 }}
                  required
                />
              </div>
              
              <Button type="submit" disabled={otp.length < 6 || loading} style={{ height: '44px' }}>
                {loading ? 'Đang xác thực...' : 'Xác nhận mã'}
              </Button>

              <div style={{ textAlign: 'center' }}>
                <p className="text-helper">
                  Chưa nhận được mã?{' '}
                  <span 
                    onClick={handleResendOTP} 
                    style={{ 
                      color: canResend ? 'var(--primary)' : 'var(--text-muted)', 
                      cursor: canResend ? 'pointer' : 'not-allowed',
                      fontWeight: 500
                    }}
                  >
                    Gửi lại mã {canResend ? '' : `(${timeLeft}s)`}
                  </span>
                </p>
              </div>
            </form>
          </div>
        )}

        {/* STEP 3: Đặt lại mật khẩu */}
        {step === 3 && (
          <div style={{ animation: 'fadeIn 0.4s ease-out' }}>
            <div style={{ textAlign: 'center', marginBottom: '32px' }}>
              <div style={{ width: '56px', height: '56px', backgroundColor: 'rgba(16, 185, 129, 0.1)', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px', color: 'var(--success)' }}>
                <Lock size={28} />
              </div>
              <h1 className="text-h1" style={{ marginBottom: '8px' }}>Tạo mật khẩu mới</h1>
              <p className="text-helper" style={{ lineHeight: 1.5 }}>Mật khẩu mới của bạn phải khác với mật khẩu sử dụng trước đó.</p>
            </div>

            <form onSubmit={handleResetPassword} style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
              <Input 
                label="Mật khẩu mới" 
                type="password" 
                placeholder="Tối thiểu 6 ký tự"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              <Input 
                label="Xác nhận mật khẩu mới" 
                type="password" 
                placeholder="Nhập lại mật khẩu mới"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
              
              <Button type="submit" disabled={!password || !confirmPassword || loading} style={{ height: '44px', marginTop: '8px' }}>
                {loading ? 'Đang đổi mật khẩu...' : 'Xác nhận đổi mật khẩu'}
              </Button>
            </form>
          </div>
        )}

      </Card>
      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: translateX(20px); }
          to { opacity: 1; transform: translateX(0); }
        }
      `}</style>
    </div>
  );
}
