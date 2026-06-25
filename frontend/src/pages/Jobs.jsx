import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Table } from '../components/Table';
import { Badge } from '../components/Badge';
import { Plus, Search, Eye, MoreHorizontal } from 'lucide-react';
import { mockApi } from '../utils/mockData';

export function Jobs() {
  const navigate = useNavigate();
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    mockApi.jobs.getAll().then(data => {
      setJobs(data);
      setLoading(false);
    });
  }, []);

  const columns = [
    { header: 'Job Title', render: (row) => <span style={{ fontWeight: 500, color: 'var(--text-main)' }}>{row.title}</span> },
    { header: 'Status', render: (row) => (
      <Badge type={row.status === 'Open' ? 'success' : 'neutral'}>
        {row.status}
      </Badge>
    )},
    { header: 'Created', accessor: 'created' },
    { header: 'Applicants', accessor: 'applicants' },
    { header: 'Action', render: (row) => (
      <div style={{ display: 'flex', gap: '8px' }}>
        <Button variant="ghost" onClick={() => navigate(`/jobs/${row.id}`)} style={{ padding: '0 8px' }}>
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
          <h1 className="text-h1">Việc làm</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>Quản lý các vị trí tuyển dụng và pipeline ứng viên.</p>
        </div>
        <Button onClick={() => navigate('/jobs/new')}><Plus size={16} /> Tạo job mới</Button>
      </div>

      <Card>
        <div style={{ display: 'flex', gap: '16px', marginBottom: '24px' }}>
          <div style={{ position: 'relative', flex: 1, maxWidth: '320px' }}>
            <Search size={16} style={{ position: 'absolute', left: '12px', top: '12px', color: 'var(--text-muted)' }} />
            <input 
              type="text" 
              placeholder="Tìm kiếm công việc..." 
              className="input-field"
              style={{ width: '100%', paddingLeft: '36px' }}
            />
          </div>
          <Button variant="secondary">Lọc theo trạng thái</Button>
        </div>

        {loading ? (
          <div style={{ padding: '32px', textAlign: 'center', color: 'var(--text-muted)' }}>Đang tải dữ liệu...</div>
        ) : (
          <Table columns={columns} data={jobs} />
        )}
      </Card>
    </div>
  );
}
