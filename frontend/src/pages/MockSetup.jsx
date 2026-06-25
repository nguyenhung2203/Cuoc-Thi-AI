import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Play, FileText, Settings, Sparkles } from 'lucide-react';

export function MockSetup() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const handleStart = (e) => {
    e.preventDefault();
    setLoading(true);
    setTimeout(() => {
      navigate('/mock-room', { state: { message: 'Bắt đầu phiên luyện tập thành công!' } });
    }, 800);
  };

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ marginBottom: '32px' }}>
        <h1 className="text-h1">Luyện phỏng vấn cùng AI</h1>
        <p className="text-helper" style={{ marginTop: '4px' }}>AI sẽ đóng vai người phỏng vấn thật, đặt câu hỏi theo vị trí ứng tuyển và đưa feedback chi tiết.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
        <Card title="Cấu hình buổi phỏng vấn">
          <form onSubmit={handleStart}>
            <div className="input-group">
              <label className="input-label">Vị trí ứng tuyển (Target Role)</label>
              <select className="input-field" required>
                <option value="frontend">Frontend Developer</option>
                <option value="backend">Backend Developer</option>
                <option value="fullstack">Fullstack Developer</option>
                <option value="pm">Product Manager</option>
              </select>
            </div>

            <div className="input-group">
              <label className="input-label">Cấp độ (Level)</label>
              <select className="input-field" required>
                <option value="fresher">Fresher</option>
                <option value="junior">Junior</option>
                <option value="middle">Middle</option>
                <option value="senior">Senior</option>
              </select>
            </div>

            <div className="input-group">
              <label className="input-label">Loại phỏng vấn</label>
              <select className="input-field" required>
                <option value="tech">Phỏng vấn Kỹ thuật (Technical)</option>
                <option value="behavior">Phỏng vấn Hành vi (Behavioral)</option>
                <option value="hr">Phỏng vấn Nhân sự (HR)</option>
              </select>
            </div>

            <div className="input-group">
              <label className="input-label">Phong cách AI (Interviewer Style)</label>
              <select className="input-field" required>
                <option value="friendly">Thân thiện, gợi mở</option>
                <option value="professional">Chuyên nghiệp, tiêu chuẩn</option>
                <option value="challenging">Khó tính, hay hỏi xoáy</option>
              </select>
            </div>

            <div style={{ marginTop: '24px', paddingTop: '24px', borderTop: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span className="text-helper" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <FileText size={16} /> Sử dụng CV hiện tại trong Hồ sơ
              </span>
              <Button type="submit" size="large" disabled={loading}>
                <Play size={16} /> {loading ? 'Đang khởi tạo...' : 'Bắt đầu ngay'}
              </Button>
            </div>
          </form>
        </Card>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card style={{ backgroundColor: 'var(--surface-soft)' }}>
            <div style={{ textAlign: 'center', padding: '16px 0' }}>
              <div style={{ width: '80px', height: '80px', borderRadius: '50%', backgroundColor: 'var(--accent)', margin: '0 auto 16px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white' }}>
                <Sparkles size={32} />
              </div>
              <h3 className="text-h2">AI Interviewer</h3>
              <p className="text-helper" style={{ marginTop: '8px' }}>Sẵn sàng hỗ trợ bạn</p>
            </div>
          </Card>
          
          <Card title="Lưu ý">
            <ul className="text-helper" style={{ paddingLeft: '20px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <li>Chuẩn bị micro và không gian yên tĩnh.</li>
              <li>Sẽ có khoảng 3-5 câu hỏi tùy thuộc vào chức danh.</li>
              <li>Bạn có thể trả lời bằng Giọng nói hoặc Văn bản.</li>
            </ul>
          </Card>
        </div>
      </div>
    </div>
  );
}
