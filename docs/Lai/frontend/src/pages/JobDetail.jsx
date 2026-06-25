import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Badge } from '../components/Badge';
import { mockApi } from '../utils/mockData';
import { Sparkles, ArrowLeft, Save, AlertCircle, CheckCircle } from 'lucide-react';

export function JobDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = id === 'new';
  
  const [job, setJob] = useState({ title: '', status: 'Open', description: '' });
  const [loading, setLoading] = useState(!isNew);
  const [saving, setSaving] = useState(false);
  const [aiAnalyzing, setAiAnalyzing] = useState(false);
  const [rubric, setRubric] = useState(null);

  useEffect(() => {
    if (!isNew) {
      mockApi.jobs.getById(id).then(data => {
        if (data) setJob(data);
        setLoading(false);
      });
    }
  }, [id, isNew]);

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    if (isNew) {
      await mockApi.jobs.create({ ...job, created: new Date().toISOString().split('T')[0] });
    } else {
      await new Promise(r => setTimeout(r, 600));
    }
    setSaving(false);
    navigate('/jobs', { state: { message: 'Lưu thông tin công việc thành công!' } });
  };

  const handleAiAnalyze = () => {
    setAiAnalyzing(true);
    setTimeout(() => {
      setAiAnalyzing(false);
      setRubric([
        { criterion: 'Technical Skills', weight: '40%' },
        { criterion: 'Communication', weight: '30%' },
        { criterion: 'Problem Solving', weight: '30%' }
      ]);
    }, 1500);
  };

  if (loading) return <div>Đang tải...</div>;

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '32px' }}>
        <Button variant="ghost" onClick={() => navigate('/jobs')} style={{ padding: '8px' }}>
          <ArrowLeft size={20} />
        </Button>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <h1 className="text-h1">{isNew ? 'Tạo Job Mới' : job.title}</h1>
            {!isNew && <Badge type={job.status === 'Open' ? 'success' : 'neutral'}>{job.status}</Badge>}
          </div>
          <p className="text-helper" style={{ marginTop: '4px' }}>Jobs {'>'} {isNew ? 'New' : job.title}</p>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Thông tin chung">
            <form onSubmit={handleSave}>
              <Input 
                label="Tiêu đề công việc" 
                value={job.title} 
                onChange={e => setJob({...job, title: e.target.value})} 
                required 
              />
              
              <div className="input-group">
                <label className="input-label">Trạng thái</label>
                <select 
                  className="input-field" 
                  value={job.status}
                  onChange={e => setJob({...job, status: e.target.value})}
                >
                  <option value="Open">Đang mở (Open)</option>
                  <option value="Closed">Đã đóng (Closed)</option>
                </select>
              </div>

              <div className="input-group">
                <label className="input-label">Mô tả công việc (JD)</label>
                <textarea 
                  className="input-field" 
                  rows="10"
                  value={job.description}
                  onChange={e => setJob({...job, description: e.target.value})}
                  placeholder="Nhập yêu cầu công việc..."
                ></textarea>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '24px', gap: '12px' }}>
                <Button type="button" variant="ghost" onClick={() => navigate('/jobs')}>Hủy</Button>
                <Button type="submit" disabled={saving}>
                  <Save size={16} /> {saving ? 'Đang lưu...' : 'Lưu thông tin'}
                </Button>
              </div>
            </form>
          </Card>
          
          {!isNew && (
            <Card title="Danh sách ứng viên">
              <div style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '32px 0' }}>
                <p className="text-body">Chưa có ứng viên nào nộp đơn.</p>
              </div>
            </Card>
          )}
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="AI Gợi ý tối ưu JD" style={{ borderTop: '4px solid var(--accent)' }}>
            <p className="text-body" style={{ color: 'var(--text-muted)', marginBottom: '16px' }}>
              Hệ thống AI sẽ tự động đọc JD và phân tích ra bộ tiêu chí chấm điểm (Rubric) phù hợp nhất.
            </p>
            <Button 
              variant="secondary" 
              style={{ width: '100%', borderColor: 'var(--accent)', color: 'var(--accent)' }} 
              onClick={handleAiAnalyze}
              disabled={aiAnalyzing || !job.description}
            >
              <Sparkles size={16} /> 
              {aiAnalyzing ? 'Đang phân tích...' : 'Bóc tách Rubric bằng AI'}
            </Button>

            {rubric && (
              <div style={{ marginTop: '24px', borderTop: '1px solid var(--border)', paddingTop: '16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '12px' }}>
                  <CheckCircle size={16} color="var(--success)" />
                  <h4 className="text-h2" style={{ fontSize: '14px' }}>Bộ tiêu chí đề xuất</h4>
                </div>
                <ul className="text-body" style={{ paddingLeft: '24px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                  {rubric.map((r, i) => (
                    <li key={i}>
                      {r.criterion} - <span style={{ color: 'var(--primary)', fontWeight: 600 }}>{r.weight}</span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
            {!job.description && !rubric && (
               <div style={{ marginTop: '16px', display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--warning)', fontSize: '13px' }}>
                 <AlertCircle size={14} /> Cần nhập JD để AI có thể phân tích.
               </div>
            )}
          </Card>
        </div>
      </div>
    </div>
  );
}
