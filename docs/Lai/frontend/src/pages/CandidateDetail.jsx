import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Badge } from '../components/Badge';
import { mockApi } from '../utils/mockData';
import { ArrowLeft, Save, Upload, FileText, Sparkles, CheckCircle } from 'lucide-react';

export function CandidateDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const fileInputRef = useRef(null);
  const isNew = id === 'new';
  
  const [candidate, setCandidate] = useState({ name: '', email: '', appliedJob: '', status: 'New' });
  const [loading, setLoading] = useState(!isNew);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [parsedData, setParsedData] = useState(null);

  useEffect(() => {
    if (!isNew) {
      mockApi.candidates.getById(id).then(data => {
        if (data) setCandidate(data);
        setLoading(false);
      });
    }
  }, [id, isNew]);

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    if (isNew) {
      await mockApi.candidates.create(candidate);
    } else {
      await new Promise(r => setTimeout(r, 600));
    }
    setSaving(false);
    navigate('/candidates', { state: { message: 'Lưu thông tin ứng viên thành công!' } });
  };

  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setUploading(true);
    setUploadProgress(0);

    const interval = setInterval(() => {
      setUploadProgress(prev => {
        if (prev >= 100) {
          clearInterval(interval);
          setUploading(false);
          setCandidate({...candidate, cv: file.name});
          
          setTimeout(() => {
            setParsedData({
              skills: ['React', 'JavaScript', 'Node.js', 'TypeScript', 'CSS'],
              experience: '3 years Frontend',
              education: 'BS Computer Science'
            });
          }, 1000);
          return 100;
        }
        return prev + 25;
      });
    }, 400);
  };

  if (loading) return <div>Đang tải...</div>;

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '32px' }}>
        <Button variant="ghost" onClick={() => navigate('/candidates')} style={{ padding: '8px' }}>
          <ArrowLeft size={20} />
        </Button>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <h1 className="text-h1">{isNew ? 'Thêm Ứng Viên Mới' : candidate.name}</h1>
            {!isNew && <Badge type="info">{candidate.status}</Badge>}
          </div>
          <p className="text-helper" style={{ marginTop: '4px' }}>Candidates {'>'} {isNew ? 'New' : candidate.name}</p>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px' }}>
        <Card title="Thông tin cá nhân">
          <form onSubmit={handleSave}>
            <Input 
              label="Họ và Tên" 
              value={candidate.name} 
              onChange={e => setCandidate({...candidate, name: e.target.value})} 
              required 
            />
            <Input 
              label="Email" 
              type="email"
              value={candidate.email} 
              onChange={e => setCandidate({...candidate, email: e.target.value})} 
              required 
            />
            <Input 
              label="Vị trí ứng tuyển" 
              value={candidate.appliedJob} 
              onChange={e => setCandidate({...candidate, appliedJob: e.target.value})} 
              required 
            />
            
            <div className="input-group">
              <label className="input-label">Trạng thái</label>
              <select 
                className="input-field" 
                value={candidate.status}
                onChange={e => setCandidate({...candidate, status: e.target.value})}
              >
                <option value="New">Mới (New)</option>
                <option value="Interviewing">Đang phỏng vấn (Interviewing)</option>
                <option value="Offered">Đã gửi Offer (Offered)</option>
                <option value="Rejected">Từ chối (Rejected)</option>
              </select>
            </div>

            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '24px', gap: '12px' }}>
              <Button type="button" variant="ghost" onClick={() => navigate('/candidates')}>Hủy</Button>
              <Button type="submit" disabled={saving}>
                <Save size={16} /> {saving ? 'Đang lưu...' : 'Lưu ứng viên'}
              </Button>
            </div>
          </form>
        </Card>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Hồ sơ (CV & Resume)">
            {candidate.cv ? (
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px', padding: '16px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)' }}>
                <FileText size={32} color="var(--primary)" />
                <div style={{ flex: 1 }}>
                  <p className="text-body" style={{ fontWeight: 500 }}>{candidate.cv}</p>
                  <p className="text-helper" style={{ color: 'var(--success)', display: 'flex', alignItems: 'center', gap: '4px', marginTop: '4px' }}>
                    <CheckCircle size={14} /> Tải lên thành công
                  </p>
                </div>
                <Button variant="secondary" onClick={() => fileInputRef.current?.click()}>Thay đổi</Button>
              </div>
            ) : (
              <div 
                style={{ 
                  border: '2px dashed var(--border)', 
                  borderRadius: 'var(--radius-lg)', 
                  padding: '40px 24px', 
                  textAlign: 'center',
                  cursor: 'pointer',
                  backgroundColor: 'var(--surface-soft)'
                }}
                onClick={() => fileInputRef.current?.click()}
              >
                {uploading ? (
                  <div>
                    <p className="text-body" style={{ marginBottom: '12px', fontWeight: 500 }}>Đang tải lên... {uploadProgress}%</p>
                    <div style={{ height: '6px', background: 'var(--border)', borderRadius: '3px', overflow: 'hidden' }}>
                      <div style={{ height: '100%', background: 'var(--primary)', width: `${uploadProgress}%`, transition: 'width 0.3s' }}></div>
                    </div>
                  </div>
                ) : (
                  <>
                    <Upload size={32} color="var(--text-muted)" style={{ margin: '0 auto 12px' }} />
                    <p className="text-body" style={{ fontWeight: 500, marginBottom: '4px' }}>Click để tải CV lên</p>
                    <p className="text-helper">Hỗ trợ PDF, DOCX tối đa 5MB</p>
                  </>
                )}
              </div>
            )}
            <input 
              type="file" 
              ref={fileInputRef} 
              style={{ display: 'none' }} 
              accept=".pdf,.doc,.docx" 
              onChange={handleFileUpload}
            />
          </Card>

          {parsedData && (
            <Card title="AI Bóc tách dữ liệu CV" style={{ borderTop: '4px solid var(--accent)' }}>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--accent)' }}>
                  <Sparkles size={16} /> <span style={{ fontSize: '14px', fontWeight: 500 }}>Hoàn tất phân tích tự động</span>
                </div>
                <div>
                  <h4 className="text-helper" style={{ textTransform: 'uppercase', marginBottom: '8px', letterSpacing: '0.05em' }}>Kỹ năng nổi bật</h4>
                  <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
                    {parsedData.skills.map(s => (
                      <span key={s} style={{ backgroundColor: 'var(--surface-soft)', padding: '4px 10px', borderRadius: '6px', fontSize: '13px', border: '1px solid var(--border)' }}>
                        {s}
                      </span>
                    ))}
                  </div>
                </div>
                <div>
                  <h4 className="text-helper" style={{ textTransform: 'uppercase', marginBottom: '4px', letterSpacing: '0.05em' }}>Kinh nghiệm</h4>
                  <p className="text-body">{parsedData.experience}</p>
                </div>
                <div>
                  <h4 className="text-helper" style={{ textTransform: 'uppercase', marginBottom: '4px', letterSpacing: '0.05em' }}>Học vấn</h4>
                  <p className="text-body">{parsedData.education}</p>
                </div>
              </div>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
