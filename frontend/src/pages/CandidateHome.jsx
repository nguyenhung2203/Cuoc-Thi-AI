import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Calendar, Video, Target, Award, ArrowRight } from 'lucide-react';

export function CandidateHome() {
  const navigate = useNavigate();
  return (
    <div>
      <div style={{ marginBottom: '32px' }}>
        <h1 className="text-h1">Chào mừng trở lại, Candidate User</h1>
        <p className="text-helper" style={{ marginTop: '4px' }}>Theo dõi lịch phỏng vấn, luyện tập với AI và cải thiện kỹ năng trả lời.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Lịch phỏng vấn sắp tới">
            <div style={{ display: 'flex', alignItems: 'center', gap: '16px', padding: '16px', border: '1px solid var(--border)', borderRadius: 'var(--radius)', backgroundColor: 'var(--surface)' }}>
              <div style={{ width: '48px', height: '48px', borderRadius: '12px', backgroundColor: 'rgba(37, 99, 235, 0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <Calendar color="var(--primary)" />
              </div>
              <div style={{ flex: 1 }}>
                <h3 className="text-body" style={{ fontWeight: 600 }}>Frontend Developer</h3>
                <p className="text-helper">Công ty: TechCorp Inc.</p>
                <div style={{ display: 'flex', gap: '16px', marginTop: '8px', fontSize: '13px', color: 'var(--text-secondary)' }}>
                  <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}><Calendar size={14} /> Ngày mai, 14:00</span>
                  <span style={{ display: 'flex', alignItems: 'center', gap: '4px' }}><Video size={14} /> Phỏng vấn Online</span>
                </div>
              </div>
              <Button onClick={() => navigate('/interview-consent')}>Vào phòng</Button>
            </div>
          </Card>

          <Card style={{ background: 'linear-gradient(135deg, var(--surface) 0%, var(--surface-soft) 100%)', borderLeft: '4px solid var(--primary)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <h2 className="text-h2" style={{ marginBottom: '8px' }}>Sẵn sàng luyện phỏng vấn với AI?</h2>
                <p className="text-body" style={{ color: 'var(--text-secondary)' }}>AI Interviewer sẽ đóng vai trò HR, đặt câu hỏi và chấm điểm kỹ năng trả lời của bạn.</p>
              </div>
              <Button size="large" style={{ whiteSpace: 'nowrap' }} onClick={() => navigate('/mock-setup')}>Bắt đầu luyện tập <ArrowRight size={16} /></Button>
            </div>
          </Card>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Mức độ hoàn thiện hồ sơ">
            <div style={{ marginBottom: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                <span className="text-body" style={{ fontWeight: 500 }}>Tiến độ</span>
                <span className="text-body" style={{ color: 'var(--success)', fontWeight: 600 }}>80%</span>
              </div>
              <div style={{ height: '8px', width: '100%', backgroundColor: 'var(--surface-soft)', borderRadius: '4px', overflow: 'hidden' }}>
                <div style={{ width: '80%', height: '100%', backgroundColor: 'var(--success)' }}></div>
              </div>
            </div>
            <ul className="text-body" style={{ display: 'flex', flexDirection: 'column', gap: '12px', color: 'var(--text-secondary)' }}>
              <li style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><Target size={16} color="var(--success)" /> CV đã tải lên</li>
              <li style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><Target size={16} color="var(--success)" /> Cập nhật kỹ năng</li>
              <li style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--text-muted)' }}><Target size={16} /> Thêm kinh nghiệm làm việc</li>
            </ul>
          </Card>

          <Card title="AI Career Coach">
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '12px' }}>
              <Award size={20} color="var(--warning)" style={{ flexShrink: 0, marginTop: '2px' }} />
              <div>
                <p className="text-body" style={{ fontWeight: 500, marginBottom: '4px' }}>Gợi ý luyện tập</p>
                <p className="text-helper" style={{ lineHeight: 1.5 }}>Dựa theo JD mới nhất, bạn nên ôn tập thêm các câu hỏi về <span style={{ fontWeight: 600, color: 'var(--text-main)' }}>React Performance Optimization</span> và <span style={{ fontWeight: 600, color: 'var(--text-main)' }}>State Management</span>.</p>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
