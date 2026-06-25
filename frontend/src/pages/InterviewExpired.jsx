import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Clock, Mail, ArrowLeft } from 'lucide-react';

export function InterviewExpired() {
  const navigate = useNavigate();

  return (
    <div style={{ minHeight: '100vh', backgroundColor: 'var(--background)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '24px' }}>
      <div style={{ maxWidth: '480px', width: '100%' }}>
        <Card style={{ padding: '40px 32px', textAlign: 'center' }}>
          <div style={{ width: '80px', height: '80px', backgroundColor: 'rgba(220, 38, 38, 0.1)', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px' }}>
            <Clock size={40} color="var(--danger)" />
          </div>
          
          <h1 className="text-h1" style={{ marginBottom: '12px' }}>Link phỏng vấn đã hết hạn</h1>
          <p className="text-body" style={{ color: 'var(--text-secondary)', marginBottom: '32px', lineHeight: 1.6 }}>
            Buổi phỏng vấn này đã kết thúc hoặc đường dẫn bạn truy cập không còn hiệu lực. Vui lòng kiểm tra lại thời gian trong email mời.
          </p>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <Button variant="primary" style={{ width: '100%' }} onClick={() => window.location.href = 'mailto:hr@techcorp.com'}>
              <Mail size={18} /> Liên hệ với Bộ phận Nhân sự
            </Button>
            <Button variant="ghost" style={{ width: '100%' }} onClick={() => navigate('/home')}>
              <ArrowLeft size={18} /> Quay về Trang chủ
            </Button>
          </div>
        </Card>
      </div>
    </div>
  );
}
