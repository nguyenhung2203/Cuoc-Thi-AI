import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Table } from '../components/Table';
import { Badge } from '../components/Badge';
import { Plus, Search, Eye, MoreHorizontal } from 'lucide-react';
import { mockApi } from '../utils/mockData';

export function Candidates() {
  const navigate = useNavigate();
  const [candidates, setCandidates] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    mockApi.candidates.getAll().then(data => {
      setCandidates(data);
      setLoading(false);
    });
  }, []);

  const getStatusType = (status) => {
    if (status === 'New') return 'info';
    if (status === 'Interviewing') return 'warning';
    if (status === 'Offered') return 'success';
    if (status === 'Rejected') return 'danger';
    return 'neutral';
  };

  const columns = [
    { header: 'Ứng viên', render: (row) => (
      <div>
        <div style={{ fontWeight: 500, color: 'var(--text-main)' }}>{row.name}</div>
        <div className="text-helper">{row.email}</div>
      </div>
    )},
    { header: 'Vị trí ứng tuyển', accessor: 'appliedJob' },
    { header: 'Trạng thái', render: (row) => (
      <Badge type={getStatusType(row.status)}>
        {row.status}
      </Badge>
    )},
    { header: 'Hành động', render: (row) => (
      <div style={{ display: 'flex', gap: '8px' }}>
        <Button variant="ghost" onClick={() => navigate(`/candidates/${row.id}`)} style={{ padding: '0 8px' }}>
          <Eye size={16} />
        </Button>
        <Button variant="ghost" style={{ padding: '0 8px' }}>
          <MoreHorizontal size={16} />
        </Button>
      </div>
    )}
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <div>
          <h1 className="text-h1">Ứng viên</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>Theo dõi ứng viên, điểm phù hợp và lịch sử phỏng vấn.</p>
        </div>
        <Button onClick={() => navigate('/candidates/new')}><Plus size={16} /> Thêm ứng viên</Button>
      </div>

      <Card>
        <div style={{ display: 'flex', gap: '16px', marginBottom: '24px' }}>
          <div style={{ position: 'relative', flex: 1, maxWidth: '320px' }}>
            <Search size={16} style={{ position: 'absolute', left: '12px', top: '12px', color: 'var(--text-muted)' }} />
            <input 
              type="text" 
              placeholder="Tìm tên, email..." 
              className="input-field"
              style={{ width: '100%', paddingLeft: '36px' }}
            />
          </div>
          <Button variant="secondary">Lọc theo Job</Button>
          <Button variant="secondary">Lọc theo trạng thái</Button>
        </div>

        {loading ? (
          <div style={{ padding: '32px', textAlign: 'center', color: 'var(--text-muted)' }}>Đang tải dữ liệu...</div>
        ) : (
          <Table columns={columns} data={candidates} />
        )}
      </Card>
    </div>
  );
}
