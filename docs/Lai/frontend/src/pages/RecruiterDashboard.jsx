import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Users, Briefcase, Calendar, Plus, FileText, Sparkles, TrendingUp } from 'lucide-react';

export function RecruiterDashboard() {
  const navigate = useNavigate();
  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <div>
          <h1 className="text-h1">Tổng quan</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>Theo dõi lịch phỏng vấn, ứng viên và báo cáo AI hôm nay.</p>
        </div>
        <div style={{ display: 'flex', gap: '16px' }}>
          <Button onClick={() => navigate('/interviews/new')}><Plus size={16} /> Tạo lịch phỏng vấn</Button>
        </div>
      </div>

      <div className="stats-grid">
        <Card>
          <div className="stat-card">
            <span className="text-helper" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Briefcase size={16} color="var(--primary)" /> Jobs đang mở
            </span>
            <span className="stat-value">12</span>
          </div>
        </Card>
        <Card>
          <div className="stat-card">
            <span className="text-helper" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Users size={16} color="var(--accent)" /> Ứng viên mới
            </span>
            <span className="stat-value">48</span>
          </div>
        </Card>
        <Card>
          <div className="stat-card">
            <span className="text-helper" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Calendar size={16} color="var(--warning)" /> Phỏng vấn hôm nay
            </span>
            <span className="stat-value">5</span>
          </div>
        </Card>
        <Card>
          <div className="stat-card">
            <span className="text-helper" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <FileText size={16} color="var(--success)" /> Báo cáo chờ xem
            </span>
            <span className="stat-value">3</span>
          </div>
        </Card>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Lịch phỏng vấn sắp tới">
            <div style={{ padding: '32px', textAlign: 'center' }}>
              <Calendar size={32} color="var(--text-muted)" style={{ margin: '0 auto 16px' }} />
              <p className="text-body" style={{ color: 'var(--text-muted)' }}>Chưa có lịch phỏng vấn nào sắp tới.</p>
            </div>
          </Card>
          <Card title="Ứng viên mới nhất">
            <div style={{ padding: '32px', textAlign: 'center' }}>
              <Users size={32} color="var(--text-muted)" style={{ margin: '0 auto 16px' }} />
              <p className="text-body" style={{ color: 'var(--text-muted)' }}>Chưa có ứng viên mới ứng tuyển.</p>
            </div>
          </Card>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="AI Insights" style={{ borderTop: '4px solid var(--accent)' }}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              <div style={{ display: 'flex', gap: '12px', alignItems: 'flex-start' }}>
                <Sparkles size={16} color="var(--accent)" style={{ marginTop: '2px', flexShrink: 0 }} />
                <div>
                  <p className="text-body" style={{ fontWeight: 500 }}>3 ứng viên có mức phù hợp cao</p>
                  <p className="text-helper">Vừa nộp đơn vào vị trí Frontend Developer.</p>
                </div>
              </div>
              <div style={{ display: 'flex', gap: '12px', alignItems: 'flex-start' }}>
                <TrendingUp size={16} color="var(--warning)" style={{ marginTop: '2px', flexShrink: 0 }} />
                <div>
                  <p className="text-body" style={{ fontWeight: 500 }}>2 buổi phỏng vấn cần review</p>
                  <p className="text-helper">Báo cáo AI đã sẵn sàng để bạn đưa ra quyết định.</p>
                </div>
              </div>
            </div>
          </Card>
          
          <Card title="Thao tác nhanh">
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <Button variant="secondary" style={{ justifyContent: 'flex-start' }} onClick={() => navigate('/jobs/new')}><Plus size={16} /> Tạo job mới</Button>
              <Button variant="secondary" style={{ justifyContent: 'flex-start' }} onClick={() => navigate('/candidates/new')}><Plus size={16} /> Thêm ứng viên</Button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
}
