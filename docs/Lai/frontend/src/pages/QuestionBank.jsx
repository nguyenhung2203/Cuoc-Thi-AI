import React, { useState } from 'react';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Table } from '../components/Table';
import { Plus, Search, Filter, Layers, Edit, Trash2 } from 'lucide-react';

export function QuestionBank() {
  const [questions] = useState([
    { id: 1, text: 'Hãy giải thích sự khác biệt giữa var, let và const trong JavaScript.', role: 'Frontend', level: 'Fresher', type: 'Technical' },
    { id: 2, text: 'Bạn đã bao giờ bất đồng quan điểm với Quản lý dự án chưa? Bạn giải quyết thế nào?', role: 'All', level: 'Middle', type: 'Behavioral' },
    { id: 3, text: 'Mô tả nguyên lý hoạt động của Virtual DOM trong React.', role: 'Frontend', level: 'Junior', type: 'Technical' },
    { id: 4, text: 'Làm thế nào để scale một hệ thống chịu tải 1 triệu requests/s?', role: 'Backend', level: 'Senior', type: 'System Design' },
  ]);

  const columns = [
    { header: 'Câu hỏi', render: (row) => <span style={{ fontWeight: 500, color: 'var(--text-main)' }}>{row.text}</span> },
    { header: 'Vị trí (Role)', accessor: 'role' },
    { header: 'Cấp độ', render: (row) => <span className="badge badge-neutral">{row.level}</span> },
    { header: 'Loại', accessor: 'type' },
    { header: 'Hành động', render: (row) => (
      <div style={{ display: 'flex', gap: '8px' }}>
        <Button variant="ghost" style={{ padding: '4px' }}><Edit size={16} /></Button>
        <Button variant="ghost" style={{ padding: '4px', color: 'var(--danger)' }}><Trash2 size={16} /></Button>
      </div>
    )}
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <div>
          <h1 className="text-h1">Kho câu hỏi</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>Quản lý ngân hàng câu hỏi dùng chung cho các buổi phỏng vấn.</p>
        </div>
        <Button><Plus size={16} /> Thêm câu hỏi</Button>
      </div>

      <Card>
        <div style={{ display: 'flex', gap: '16px', marginBottom: '24px' }}>
          <div style={{ position: 'relative', flex: 1, maxWidth: '400px' }}>
            <Search size={16} style={{ position: 'absolute', left: '12px', top: '12px', color: 'var(--text-muted)' }} />
            <input 
              type="text" 
              placeholder="Tìm kiếm nội dung câu hỏi..." 
              className="input-field"
              style={{ width: '100%', paddingLeft: '36px' }}
            />
          </div>
          <select className="input-field" style={{ width: '180px' }}>
            <option value="">Tất cả vị trí (Role)</option>
            <option value="frontend">Frontend</option>
            <option value="backend">Backend</option>
          </select>
          <select className="input-field" style={{ width: '180px' }}>
            <option value="">Tất cả cấp độ</option>
            <option value="fresher">Fresher</option>
            <option value="senior">Senior</option>
          </select>
          <Button variant="secondary"><Filter size={16} /> Lọc nâng cao</Button>
        </div>

        <Table columns={columns} data={questions} />
      </Card>
      
      <div style={{ marginTop: '24px', padding: '16px', backgroundColor: 'rgba(37, 99, 235, 0.05)', borderRadius: 'var(--radius)', border: '1px solid rgba(37, 99, 235, 0.2)', display: 'flex', alignItems: 'flex-start', gap: '12px' }}>
        <Layers size={20} color="var(--primary)" style={{ flexShrink: 0, marginTop: '2px' }} />
        <div>
          <p className="text-body" style={{ fontWeight: 600, color: 'var(--primary)' }}>AI Suggestion Engine đang hoạt động</p>
          <p className="text-helper" style={{ color: 'var(--text-secondary)' }}>Trong lúc phỏng vấn, AI sẽ tự động tìm kiếm các câu hỏi liên quan trong Kho câu hỏi này dựa trên ngữ cảnh để gợi ý cho bạn.</p>
        </div>
      </div>
    </div>
  );
}
