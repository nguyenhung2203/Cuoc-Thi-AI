import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Badge } from '../components/Badge';
import { Toast } from '../components/Toast';
import { mockApi } from '../utils/mockData';
import { ArrowLeft, Video, Copy, Calendar, Clock, User, Briefcase, Mail } from 'lucide-react';

export function InterviewDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  
  const [interview, setInterview] = useState(null);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    mockApi.interviews.getById(id).then(data => {
      setInterview(data);
      setLoading(false);
    });
  }, [id]);

  const copyLink = () => {
    if (interview && interview.link) {
      navigator.clipboard.writeText(interview.link);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  if (loading) return <div style={{ padding: '32px', textAlign: 'center', color: 'var(--text-muted)' }}>Đang tải thông tin phỏng vấn...</div>;
  if (!interview) return <div style={{ padding: '32px', textAlign: 'center', color: 'var(--warning)' }}>Không tìm thấy thông tin phỏng vấn.</div>;

  const dateObj = new Date(interview.datetime);

  return (
    <div>
      {copied && (
        <Toast type="success" message="Link phòng phỏng vấn đã được copy!" />
      )}

      <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '32px' }}>
        <Button variant="ghost" onClick={() => navigate('/interviews')} style={{ padding: '8px' }}>
          <ArrowLeft size={20} />
        </Button>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <h1 className="text-h1">Chi tiết buổi phỏng vấn</h1>
            <Badge type={interview.status === 'Scheduled' ? 'info' : 'neutral'}>
              {interview.status === 'Scheduled' ? 'Đã lên lịch' : interview.status}
            </Badge>
          </div>
          <p className="text-helper" style={{ marginTop: '4px' }}>Lịch phỏng vấn {'>'} {interview.candidateName}</p>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Thông tin chung">
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px' }}>
              <div>
                <p className="text-helper" style={{ marginBottom: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}><User size={14} /> Ứng viên</p>
                <p className="text-body" style={{ fontWeight: 500 }}>{interview.candidateName}</p>
              </div>
              <div>
                <p className="text-helper" style={{ marginBottom: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}><Briefcase size={14} /> Vị trí ứng tuyển</p>
                <p className="text-body" style={{ fontWeight: 500 }}>{interview.jobTitle}</p>
              </div>
              <div>
                <p className="text-helper" style={{ marginBottom: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}><Calendar size={14} /> Ngày phỏng vấn</p>
                <p className="text-body" style={{ fontWeight: 500 }}>{dateObj.toLocaleDateString('vi-VN')}</p>
              </div>
              <div>
                <p className="text-helper" style={{ marginBottom: '4px', display: 'flex', alignItems: 'center', gap: '6px' }}><Clock size={14} /> Thời gian</p>
                <p className="text-body" style={{ fontWeight: 500 }}>{dateObj.toLocaleTimeString('vi-VN', {hour: '2-digit', minute:'2-digit'})}</p>
              </div>
            </div>
            
            <div style={{ marginTop: '24px', paddingTop: '24px', borderTop: '1px solid var(--border)' }}>
              <p className="text-helper" style={{ marginBottom: '12px', fontWeight: 500, color: 'var(--text-main)' }}>Hành động</p>
              <div style={{ display: 'flex', gap: '12px' }}>
                <Button onClick={() => navigate('/recruiter-room', { state: { message: 'Vào phòng phỏng vấn thành công!' } })}>
                  <Video size={16} /> Vào phòng phỏng vấn
                </Button>
                <Button variant="secondary" onClick={copyLink}>
                  <Copy size={16} /> Copy Link Invite
                </Button>
                <Button variant="ghost">
                  <Mail size={16} /> Gửi email nhắc nhở
                </Button>
              </div>
            </div>
          </Card>
          
          <Card title="Cấu hình AI Assistant" style={{ borderTop: '4px solid var(--accent)' }}>
            <ul className="text-body" style={{ paddingLeft: '20px', display: 'flex', flexDirection: 'column', gap: '12px' }}>
              <li><strong>Phân tích realtime:</strong> Bật</li>
              <li><strong>Tạo transcript:</strong> Bật</li>
              <li><strong>Gợi ý câu hỏi (Rubric-based):</strong> Bật</li>
              <li><strong>Tự động chấm điểm:</strong> Bật</li>
            </ul>
          </Card>
        </div>

        <div>
          <Card title="Ghi chú nội bộ">
            <textarea 
              className="input-field" 
              rows="6"
              placeholder="Nhập ghi chú hoặc nhắc nhở trước buổi phỏng vấn (Chỉ recruiter xem được)..."
              style={{ width: '100%', marginBottom: '16px' }}
            ></textarea>
            <Button variant="secondary" style={{ width: '100%' }}>Lưu ghi chú</Button>
          </Card>
        </div>
      </div>
    </div>
  );
}
