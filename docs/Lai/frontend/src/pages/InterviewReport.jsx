import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Badge } from '../components/Badge';
import { ArrowLeft, Download, Share2, CheckCircle, AlertTriangle, FileText } from 'lucide-react';

export function InterviewReport() {
  const navigate = useNavigate();

  return (
    <div style={{ maxWidth: '900px', margin: '0 auto', paddingBottom: '64px' }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '32px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <Button variant="ghost" onClick={() => navigate('/interviews')} style={{ padding: '8px' }}>
            <ArrowLeft size={20} />
          </Button>
          <div>
            <h1 className="text-h1">Báo cáo Phỏng vấn: Nguyễn Văn A</h1>
            <p className="text-helper" style={{ marginTop: '4px' }}>Frontend Developer • 14/10/2023</p>
          </div>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          <Button variant="secondary"><Share2 size={16} /> Chia sẻ</Button>
          <Button><Download size={16} /> Xuất PDF</Button>
        </div>
      </div>

      {/* Top Overview */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '16px', marginBottom: '24px' }}>
        <Card style={{ padding: '20px', textAlign: 'center', backgroundColor: 'var(--surface)' }}>
          <p className="text-helper" style={{ textTransform: 'uppercase', marginBottom: '8px', fontWeight: 600 }}>Điểm tổng quan</p>
          <p style={{ fontSize: '36px', fontWeight: 700, color: 'var(--primary)' }}>8.5<span style={{ fontSize: '18px', color: 'var(--text-muted)' }}>/10</span></p>
        </Card>
        <Card style={{ padding: '20px', textAlign: 'center', backgroundColor: 'var(--surface)' }}>
          <p className="text-helper" style={{ textTransform: 'uppercase', marginBottom: '8px', fontWeight: 600 }}>Đề xuất từ AI</p>
          <Badge type="success" style={{ fontSize: '14px', padding: '6px 16px', marginTop: '4px' }}>Nên tuyển (Hire)</Badge>
        </Card>
        <Card style={{ padding: '20px', gridColumn: 'span 2', backgroundColor: 'rgba(37, 99, 235, 0.05)', border: '1px solid rgba(37, 99, 235, 0.1)' }}>
          <p className="text-helper" style={{ textTransform: 'uppercase', marginBottom: '8px', fontWeight: 600, color: 'var(--primary)' }}>Nhận xét cốt lõi</p>
          <p className="text-body" style={{ fontWeight: 500 }}>Ứng viên có kiến thức chuyên môn sâu về React, thái độ giao tiếp chuyên nghiệp. Phù hợp văn hóa công ty.</p>
        </Card>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '24px' }}>
        
        {/* Left: Rubric Scores */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Điểm chi tiết (Rubric)">
            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
              {[
                { name: 'Technical Knowledge', score: 9, color: 'var(--success)' },
                { name: 'Problem Solving', score: 8, color: 'var(--primary)' },
                { name: 'Communication', score: 9, color: 'var(--success)' },
                { name: 'Experience Relevance', score: 7, color: 'var(--warning)' },
                { name: 'Culture Fit', score: 8.5, color: 'var(--primary)' }
              ].map(item => (
                <div key={item.name}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                    <span className="text-body" style={{ fontWeight: 500 }}>{item.name}</span>
                    <span className="text-body" style={{ fontWeight: 600 }}>{item.score}/10</span>
                  </div>
                  <div style={{ height: '8px', backgroundColor: 'var(--surface-soft)', borderRadius: '4px', overflow: 'hidden' }}>
                    <div style={{ height: '100%', width: `${item.score * 10}%`, backgroundColor: item.color, borderRadius: '4px' }}></div>
                  </div>
                </div>
              ))}
            </div>
          </Card>
          
          <Card title="Quyết định của Bạn">
            <select className="input-field" style={{ width: '100%', marginBottom: '16px' }}>
              <option value="">-- Chọn quyết định --</option>
              <option value="offer">Gửi Offer</option>
              <option value="reject">Từ chối</option>
              <option value="next_round">Phỏng vấn vòng sau</option>
            </select>
            <textarea className="input-field" placeholder="Nhập ghi chú HR..." style={{ width: '100%', height: '100px', marginBottom: '16px' }}></textarea>
            <Button style={{ width: '100%' }}>Lưu quyết định</Button>
          </Card>
        </div>

        {/* Right: Insights & Evidence */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          
          <Card title="Phân tích Điểm mạnh & Rủi ro" style={{ borderTop: '4px solid var(--accent)' }}>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px' }}>
              <div>
                <h4 className="text-body" style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--success)', fontWeight: 600, marginBottom: '12px' }}>
                  <CheckCircle size={18} /> Điểm mạnh
                </h4>
                <ul className="text-body" style={{ paddingLeft: '24px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  <li>Nắm rất vững vòng đời (Lifecycle) và Hooks trong React.</li>
                  <li>Cách trình bày vấn đề có cấu trúc rõ ràng (STAR method).</li>
                </ul>
              </div>
              <div>
                <h4 className="text-body" style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--warning)', fontWeight: 600, marginBottom: '12px' }}>
                  <AlertTriangle size={18} /> Điểm rủi ro (Cần lưu ý)
                </h4>
                <ul className="text-body" style={{ paddingLeft: '24px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  <li>Chưa có nhiều kinh nghiệm lead team trực tiếp.</li>
                  <li>Mức lương đề xuất hơi cao so với budget ban đầu.</li>
                </ul>
              </div>
            </div>
          </Card>

          <Card title="Trích xuất Transcript (Bằng chứng)">
            <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
              
              <div style={{ padding: '16px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)', borderLeft: '4px solid var(--success)' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
                  <Badge type="success">Bằng chứng: Technical</Badge>
                  <span className="text-helper">14:15:20</span>
                </div>
                <p className="text-body" style={{ fontStyle: 'italic', color: 'var(--text-secondary)' }}>
                  "...Trong dự án đó, khi list data lên tới 10,000 items, em đã sử dụng kỹ thuật Virtualization (react-window) kết hợp với useMemo để tránh re-render những components không cần thiết, giúp FPS tăng từ 20 lên 60."
                </p>
              </div>

              <div style={{ padding: '16px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)', borderLeft: '4px solid var(--warning)' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
                  <Badge type="warning">Bằng chứng: Leadership</Badge>
                  <span className="text-helper">14:30:10</span>
                </div>
                <p className="text-body" style={{ fontStyle: 'italic', color: 'var(--text-secondary)' }}>
                  "...Thường thì em sẽ nhận task từ Tech Lead và tự implement độc lập. Thỉnh thoảng em có review code cho các bạn Junior nhưng chưa chính thức lead dự án nào."
                </p>
              </div>

            </div>
            <Button variant="ghost" style={{ width: '100%', marginTop: '16px' }}><FileText size={16} /> Xem toàn bộ Transcript</Button>
          </Card>

        </div>
      </div>
    </div>
  );
}
