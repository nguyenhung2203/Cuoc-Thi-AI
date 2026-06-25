import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Table } from '../components/Table';
import { Badge } from '../components/Badge';
import { Toast } from '../components/Toast';
import { Plus, Video, Copy, ExternalLink } from 'lucide-react';
import { mockApi } from '../utils/mockData';

export function Interviews() {
  const navigate = useNavigate();
  const [interviews, setInterviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    mockApi.interviews.getAll().then(data => {
      setInterviews(data);
      setLoading(false);
    });
  }, []);

  const copyLink = (link) => {
    navigator.clipboard.writeText(link);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const columns = [
    { header: 'Vị trí & Ứng viên', render: (row) => (
      <div>
        <div style={{ fontWeight: 500, color: 'var(--text-main)' }}>{row.candidateName}</div>
        <div className="text-helper">Ứng tuyển: {row.jobTitle}</div>
      </div>
    )},
    { header: 'Thời gian', render: (row) => {
      const date = new Date(row.datetime);
      return (
        <div>
          <div style={{ fontWeight: 500, color: 'var(--text-main)' }}>{date.toLocaleTimeString('vi-VN', {hour: '2-digit', minute:'2-digit'})}</div>
          <div className="text-helper">{date.toLocaleDateString('vi-VN')}</div>
        </div>
      );
    }},
    { header: 'Trạng thái', render: (row) => (
      <Badge type={row.status === 'Scheduled' ? 'info' : 'neutral'}>
        {row.status === 'Scheduled' ? 'Đã lên lịch' : row.status}
      </Badge>
    )},
    { header: 'Phòng phỏng vấn', render: (row) => (
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <a href="#" onClick={(e) => { e.preventDefault(); navigate('/recruiter-room', { state: { message: 'Vào phòng phỏng vấn thành công!' } }); }} style={{ display: 'inline-flex', alignItems: 'center', gap: '4px', color: 'var(--primary)', fontWeight: 500, fontSize: '13px' }}>
          <Video size={14} /> Vào phòng
        </a>
        <Button variant="ghost" onClick={() => copyLink(row.link)} style={{ padding: '4px' }} title="Copy Link">
          <Copy size={14} />
        </Button>
      </div>
    )},
    { header: 'Hành động', render: (row) => (
      <Button variant="ghost" style={{ padding: '4px 8px', fontSize: '13px' }} onClick={() => navigate(`/interviews/${row.id}`)}>
        <ExternalLink size={14} style={{ marginRight: '4px' }} /> Chi tiết
      </Button>
    )}
  ];

  return (
    <div>
      {copied && (
        <Toast type="success" message="Link phòng phỏng vấn đã được copy!" />
      )}

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <div>
          <h1 className="text-h1">Lịch phỏng vấn</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>Quản lý các buổi phỏng vấn trực tiếp với ứng viên.</p>
        </div>
        <Button onClick={() => navigate('/interviews/new')}><Plus size={16} /> Tạo lịch phỏng vấn</Button>
      </div>

      <Card>
        {loading ? (
          <div style={{ padding: '32px', textAlign: 'center', color: 'var(--text-muted)' }}>Đang tải lịch phỏng vấn...</div>
        ) : (
          <Table columns={columns} data={interviews} />
        )}
      </Card>
    </div>
  );
}
