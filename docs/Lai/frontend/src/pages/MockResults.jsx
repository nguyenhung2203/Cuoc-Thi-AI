import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Badge } from '../components/Badge';
import { Table } from '../components/Table';
import { Target, TrendingUp, CheckCircle, ArrowRight, Eye } from 'lucide-react';

export function MockResults() {
  const navigate = useNavigate();
  const history = [
    { id: 1, role: 'Frontend Developer', level: 'Middle', date: '20/10/2023', score: 8.5 },
    { id: 2, role: 'Frontend Developer', level: 'Junior', date: '15/10/2023', score: 7.0 },
    { id: 3, role: 'React Developer', level: 'Junior', date: '05/10/2023', score: 6.5 },
  ];

  const columns = [
    { header: 'Vị trí luyện tập', accessor: 'role' },
    { header: 'Cấp độ', accessor: 'level' },
    { header: 'Ngày thực hiện', accessor: 'date' },
    { header: 'Điểm số', render: (row) => (
      <span style={{ fontWeight: 'bold', color: row.score >= 8 ? 'var(--success)' : row.score >= 7 ? 'var(--warning)' : 'var(--danger)' }}>
        {row.score}/10
      </span>
    )},
    { header: 'Hành động', render: (row) => (
      <Button variant="ghost" style={{ padding: '4px 8px' }} onClick={() => navigate(`/mock-results/${row.id}`)}>
        <Eye size={16} /> Xem
      </Button>
    )}
  ];

  return (
    <div>
      <div style={{ marginBottom: '32px' }}>
        <h1 className="text-h1">Kết quả & Lịch sử luyện tập</h1>
        <p className="text-helper" style={{ marginTop: '4px' }}>Theo dõi sự tiến bộ của bạn qua các bài luyện tập với AI.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '24px', marginBottom: '32px' }}>
        <Card style={{ padding: '24px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            <div style={{ width: '48px', height: '48px', borderRadius: '12px', backgroundColor: 'rgba(22, 163, 74, 0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <Target color="var(--success)" />
            </div>
            <div>
              <p className="text-helper">Điểm trung bình</p>
              <h2 className="text-h1">7.3<span style={{ fontSize: '16px', color: 'var(--text-muted)', fontWeight: 400 }}>/10</span></h2>
            </div>
          </div>
        </Card>
        
        <Card style={{ padding: '24px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            <div style={{ width: '48px', height: '48px', borderRadius: '12px', backgroundColor: 'rgba(37, 99, 235, 0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <TrendingUp color="var(--primary)" />
            </div>
            <div>
              <p className="text-helper">Số lần luyện tập</p>
              <h2 className="text-h1">3<span style={{ fontSize: '16px', color: 'var(--text-muted)', fontWeight: 400 }}> lần</span></h2>
            </div>
          </div>
        </Card>

        <Card style={{ padding: '24px', backgroundColor: 'var(--surface-soft)' }}>
          <p className="text-helper" style={{ marginBottom: '8px', fontWeight: 600 }}>Điểm cần khắc phục</p>
          <ul className="text-body" style={{ paddingLeft: '20px', color: 'var(--danger)' }}>
            <li>Trình bày cấu trúc câu trả lời (STAR)</li>
            <li>Đưa ra ví dụ thực tế</li>
          </ul>
        </Card>
      </div>

      <Card title="Lịch sử bài luyện tập">
        <Table columns={columns} data={history} />
      </Card>
    </div>
  );
}
