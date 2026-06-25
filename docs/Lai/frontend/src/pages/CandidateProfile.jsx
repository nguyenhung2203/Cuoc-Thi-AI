import React, { useState } from 'react';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Toast } from '../components/Toast';
import { Upload, FileText, CheckCircle, Save } from 'lucide-react';

export function CandidateProfile() {
  const [profile, setProfile] = useState({
    name: 'Candidate User',
    email: 'candidate@test.com',
    phone: '0123456789',
    linkedin: 'linkedin.com/in/candidate',
    targetRole: 'Frontend Developer',
    level: 'Middle',
  });

  const [saving, setSaving] = useState(false);
  const [toast, setToast] = useState(null);

  const handleSave = (e) => {
    e.preventDefault();
    setSaving(true);
    setTimeout(() => {
      setSaving(false);
      setToast({ type: 'success', message: 'Hồ sơ cá nhân đã được lưu thành công! Dữ liệu này sẽ được đồng bộ với AI.' });
    }, 800);
  };

  return (
    <div style={{ maxWidth: '1000px', margin: '0 auto' }}>
      <div style={{ marginBottom: '32px' }}>
        <h1 className="text-h1">Hồ sơ cá nhân & CV</h1>
        <p className="text-helper" style={{ marginTop: '4px' }}>Cập nhật thông tin để AI có thể đưa ra bài luyện tập chính xác nhất.</p>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '24px' }}>
        
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="Thông tin cơ bản">
            <form onSubmit={handleSave}>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                <Input label="Họ và Tên" value={profile.name} onChange={e => setProfile({...profile, name: e.target.value})} required />
                <Input label="Email" type="email" value={profile.email} disabled />
                <Input label="Số điện thoại" value={profile.phone} onChange={e => setProfile({...profile, phone: e.target.value})} />
                <Input label="LinkedIn Profile" value={profile.linkedin} onChange={e => setProfile({...profile, linkedin: e.target.value})} />
              </div>

              <h3 className="text-h2" style={{ marginTop: '32px', marginBottom: '16px' }}>Định hướng nghề nghiệp</h3>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                <Input label="Vị trí mục tiêu (Target Role)" value={profile.targetRole} onChange={e => setProfile({...profile, targetRole: e.target.value})} />
                <div className="input-group">
                  <label className="input-label">Cấp độ hiện tại</label>
                  <select className="input-field" value={profile.level} onChange={e => setProfile({...profile, level: e.target.value})}>
                    <option value="Intern">Intern</option>
                    <option value="Fresher">Fresher</option>
                    <option value="Junior">Junior</option>
                    <option value="Middle">Middle</option>
                    <option value="Senior">Senior</option>
                  </select>
                </div>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '24px' }}>
                <Button type="submit" disabled={saving}>
                  <Save size={16} /> {saving ? 'Đang lưu...' : 'Lưu hồ sơ'}
                </Button>
              </div>
            </form>
          </Card>

          <Card title="Kỹ năng chuyên môn">
            <textarea className="input-field" rows="4" placeholder="Ví dụ: ReactJS, NodeJS, TypeScript..." style={{ width: '100%' }} defaultValue="ReactJS, Redux, JavaScript, HTML, CSS, Git"></textarea>
          </Card>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card title="CV của bạn" style={{ backgroundColor: 'var(--surface)' }}>
            <div style={{ padding: '16px', backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius)', display: 'flex', alignItems: 'center', gap: '12px', border: '1px solid var(--border)', marginBottom: '16px' }}>
              <FileText size={32} color="var(--primary)" />
              <div style={{ flex: 1 }}>
                <p className="text-body" style={{ fontWeight: 600 }}>NguyenVanA_CV.pdf</p>
                <p className="text-helper" style={{ color: 'var(--success)', display: 'flex', alignItems: 'center', gap: '4px' }}>
                  <CheckCircle size={14} /> Cập nhật 2 ngày trước
                </p>
              </div>
            </div>

            <div style={{ border: '2px dashed var(--border)', borderRadius: 'var(--radius)', padding: '24px', textAlign: 'center', cursor: 'pointer' }}>
              <Upload size={24} color="var(--text-muted)" style={{ margin: '0 auto 8px' }} />
              <p className="text-body" style={{ fontWeight: 500 }}>Tải CV mới lên</p>
              <p className="text-helper">PDF, DOCX (Tối đa 5MB)</p>
            </div>
            
            <p className="text-helper" style={{ marginTop: '16px', color: 'var(--text-muted)', lineHeight: 1.5 }}>
              CV của bạn sẽ được AI dùng làm ngữ cảnh (context) để đặt câu hỏi sát với kinh nghiệm thực tế của bạn trong phần Luyện phỏng vấn Mock Interview.
            </p>
          </Card>
        </div>

      </div>

      {toast && (
        <Toast 
          type={toast.type} 
          message={toast.message} 
          onClose={() => setToast(null)} 
        />
      )}
    </div>
  );
}
