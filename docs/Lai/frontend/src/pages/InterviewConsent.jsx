import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { AlertTriangle, Video, Mic, ShieldCheck } from 'lucide-react';

export function InterviewConsent() {
  const navigate = useNavigate();
  const [agreed, setAgreed] = useState(false);

  return (
    <div style={{ minHeight: '100vh', backgroundColor: 'var(--background)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '24px' }}>
      <div style={{ maxWidth: '600px', width: '100%' }}>
        <div style={{ textAlign: 'center', marginBottom: '32px' }}>
          <div style={{ width: '64px', height: '64px', backgroundColor: 'var(--primary)', color: 'white', borderRadius: '16px', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px' }}>
            <ShieldCheck size={32} />
          </div>
          <h1 className="text-h1">Chuẩn bị vào phòng phỏng vấn</h1>
          <p className="text-helper" style={{ marginTop: '8px', fontSize: '15px' }}>Vị trí: Frontend Developer - TechCorp Inc.</p>
        </div>

        <Card style={{ padding: '32px' }}>
          <div style={{ backgroundColor: 'rgba(217, 119, 6, 0.1)', border: '1px solid rgba(217, 119, 6, 0.3)', borderRadius: 'var(--radius)', padding: '20px', marginBottom: '24px', display: 'flex', gap: '16px' }}>
            <AlertTriangle size={24} color="var(--warning)" style={{ flexShrink: 0 }} />
            <div>
              <h3 className="text-body" style={{ fontWeight: 600, color: 'var(--warning)', marginBottom: '8px' }}>Thông báo về Quyền riêng tư & AI</h3>
              <p className="text-body" style={{ color: 'var(--text-secondary)', lineHeight: 1.6 }}>
                Buổi phỏng vấn này sẽ có sự hỗ trợ của Trí tuệ nhân tạo (AI). Nhằm mục đích đánh giá công bằng và cung cấp báo cáo chi tiết cho Nhà tuyển dụng, toàn bộ diễn biến bao gồm:
              </p>
              <ul className="text-body" style={{ color: 'var(--text-secondary)', marginTop: '8px', paddingLeft: '24px', lineHeight: 1.6 }}>
                <li>Dữ liệu âm thanh (Giọng nói) sẽ được ghi lại và chuyển ngữ thành văn bản.</li>
                <li>Dữ liệu văn bản sẽ được AI phân tích tự động.</li>
                <li>Hình ảnh từ Camera sẽ được truyền trực tiếp nhưng KHÔNG lưu trữ.</li>
              </ul>
            </div>
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '16px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)', marginBottom: '24px' }}>
            <div style={{ display: 'flex', gap: '8px', color: 'var(--primary)' }}>
              <Video size={20} /> <Mic size={20} />
            </div>
            <p className="text-body">Hệ thống sẽ yêu cầu quyền truy cập Camera và Micro ở bước tiếp theo.</p>
          </div>

          <label style={{ display: 'flex', alignItems: 'flex-start', gap: '12px', cursor: 'pointer', marginBottom: '32px' }}>
            <input 
              type="checkbox" 
              checked={agreed} 
              onChange={(e) => setAgreed(e.target.checked)} 
              style={{ width: '20px', height: '20px', marginTop: '2px', accentColor: 'var(--primary)' }} 
            />
            <span className="text-body" style={{ fontWeight: 500 }}>
              Tôi đã đọc, hiểu rõ và đồng ý với việc sử dụng hệ thống AI phân tích và ghi âm trong buổi phỏng vấn này.
            </span>
          </label>

          <div style={{ display: 'flex', gap: '16px' }}>
            <Button variant="secondary" style={{ flex: 1 }} onClick={() => navigate('/home')}>Từ chối & Quay lại</Button>
            <Button variant="primary" style={{ flex: 1 }} disabled={!agreed} onClick={() => navigate('/candidate-room', { state: { message: 'Vào phòng phỏng vấn thành công!' } })}>
              Tham gia phỏng vấn
            </Button>
          </div>
        </Card>
      </div>
    </div>
  );
}
