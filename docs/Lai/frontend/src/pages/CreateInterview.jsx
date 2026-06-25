import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { ArrowLeft, Calendar, Mail, AlertCircle } from 'lucide-react';
import { mockApi } from '../utils/mockData';

export function CreateInterview() {
  const navigate = useNavigate();
  const [jobs, setJobs] = useState([]);
  const [candidates, setCandidates] = useState([]);
  
  const [interview, setInterview] = useState({
    jobTitle: '',
    candidateName: '',
    datetime: ''
  });
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    mockApi.jobs.getAll().then(setJobs);
    mockApi.candidates.getAll().then(setCandidates);
  }, []);

  const handleSave = async (e) => {
    e.preventDefault();
    if (!interview.jobTitle || !interview.candidateName || !interview.datetime) return;
    
    setSaving(true);
    await mockApi.interviews.create(interview);
    setSaving(false);
    navigate('/interviews', { state: { message: 'Lên lịch phỏng vấn thành công!' } });
  };

  return (
    <div style={{ maxWidth: '700px', margin: '0 auto' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '32px' }}>
        <Button variant="ghost" onClick={() => navigate('/interviews')} style={{ padding: '8px' }}>
          <ArrowLeft size={20} />
        </Button>
        <div>
          <h1 className="text-h1">Lên lịch phỏng vấn</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>Chọn ứng viên và thời gian để hệ thống tạo phòng phỏng vấn ảo.</p>
        </div>
      </div>

      <Card>
        <form onSubmit={handleSave}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            <div className="input-group">
              <label className="input-label">Công việc (Vị trí ứng tuyển)</label>
              <select 
                className="input-field" 
                required
                value={interview.jobTitle}
                onChange={e => setInterview({...interview, jobTitle: e.target.value})}
              >
                <option value="">-- Chọn công việc --</option>
                {jobs.map(j => <option key={j.id} value={j.title}>{j.title}</option>)}
              </select>
            </div>

            <div className="input-group">
              <label className="input-label">Ứng viên</label>
              <select 
                className="input-field" 
                required
                value={interview.candidateName}
                onChange={e => setInterview({...interview, candidateName: e.target.value})}
              >
                <option value="">-- Chọn ứng viên --</option>
                {candidates.map(c => <option key={c.id} value={c.name}>{c.name} ({c.email})</option>)}
              </select>
            </div>

            <div className="input-group">
              <label className="input-label">Ngày và giờ phỏng vấn</label>
              <input 
                type="datetime-local" 
                className="input-field" 
                required
                value={interview.datetime}
                onChange={e => setInterview({...interview, datetime: e.target.value})}
              />
            </div>
            
            <div style={{ padding: '16px', backgroundColor: 'rgba(37, 99, 235, 0.05)', border: '1px solid rgba(37, 99, 235, 0.1)', borderRadius: '8px', display: 'flex', gap: '12px' }}>
              <AlertCircle size={20} color="var(--primary)" style={{ flexShrink: 0 }} />
              <div>
                <p className="text-body" style={{ fontWeight: 500, color: 'var(--primary)', marginBottom: '4px' }}>Thông báo tự động</p>
                <p className="text-helper" style={{ color: 'var(--text-secondary)' }}>
                  Hệ thống sẽ tự động gửi email chứa link phòng phỏng vấn và lịch trình đến ứng viên ngay sau khi bạn bấm Lưu.
                </p>
              </div>
            </div>

            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px', gap: '12px' }}>
              <Button type="button" variant="ghost" onClick={() => navigate('/interviews')}>Hủy</Button>
              <Button type="submit" disabled={saving}>
                <Calendar size={16} /> {saving ? 'Đang tạo lịch...' : 'Lưu và gửi Email mời'}
              </Button>
            </div>
          </div>
        </form>
      </Card>
    </div>
  );
}
