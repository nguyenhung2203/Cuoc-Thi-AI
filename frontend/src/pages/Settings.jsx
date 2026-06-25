import React, { useState } from 'react';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import { Toast } from '../components/Toast';
import { Modal } from '../components/Modal';
import { Save, Key, Bell, Shield, User, Monitor } from 'lucide-react';

export function Settings() {
  const role = localStorage.getItem('role') || 'recruiter';
  const [activeTab, setActiveTab] = useState('account');
  const [toast, setToast] = useState(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const handleSave = (e) => {
    if (e) e.preventDefault();
    setToast({ type: 'success', message: 'Các thay đổi đã được lưu thành công!' });
  };

  const handleDeleteAccount = () => {
    setShowDeleteModal(true);
  };

  const confirmDeleteAccount = () => {
    setShowDeleteModal(false);
    setToast({ type: 'success', message: 'Yêu cầu xóa tài khoản đã được gửi.' });
  };

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <div style={{ marginBottom: '32px' }}>
        <h1 className="text-h1">Cài đặt Hệ thống</h1>
        <p className="text-helper" style={{ marginTop: '4px' }}>Tùy chỉnh thông tin tài khoản và cấu hình hệ thống.</p>
      </div>

      {toast && (
        <Toast type={toast.type} message={toast.message} onClose={() => setToast(null)} />
      )}

      <div style={{ display: 'grid', gridTemplateColumns: '240px 1fr', gap: '32px' }}>
        
        {/* Settings Sidebar Menu */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
          <div 
            className={`nav-item ${activeTab === 'account' ? 'active' : ''}`} 
            style={{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'account' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'account' ? 'var(--primary)' : 'var(--text-secondary)' }}
            onClick={() => setActiveTab('account')}
          >
            <User size={16} style={{ marginRight: '8px' }} /> Tài khoản
          </div>
          <div 
            className={`nav-item ${activeTab === 'notifications' ? 'active' : ''}`} 
            style={{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'notifications' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'notifications' ? 'var(--primary)' : 'var(--text-secondary)' }}
            onClick={() => setActiveTab('notifications')}
          >
            <Bell size={16} style={{ marginRight: '8px' }} /> Thông báo
          </div>
          <div 
            className={`nav-item ${activeTab === 'privacy' ? 'active' : ''}`} 
            style={{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'privacy' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'privacy' ? 'var(--primary)' : 'var(--text-secondary)' }}
            onClick={() => setActiveTab('privacy')}
          >
            <Shield size={16} style={{ marginRight: '8px' }} /> Quyền riêng tư & Bảo mật
          </div>
          {role === 'recruiter' && (
            <div 
              className={`nav-item ${activeTab === 'workspace' ? 'active' : ''}`} 
              style={{ borderRadius: '8px', cursor: 'pointer', padding: '12px 16px', borderRight: 'none', backgroundColor: activeTab === 'workspace' ? 'rgba(37, 99, 235, 0.08)' : 'transparent', color: activeTab === 'workspace' ? 'var(--primary)' : 'var(--text-secondary)' }}
              onClick={() => setActiveTab('workspace')}
            >
              <Monitor size={16} style={{ marginRight: '8px' }} /> Cấu hình AI Workspace
            </div>
          )}
        </div>

        {/* Settings Content */}
        <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
          
          {activeTab === 'account' && (
            <Card title="Thông tin tài khoản">
              <form onSubmit={handleSave}>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <Input label="Tên người dùng" defaultValue={role === 'recruiter' ? 'Recruiter User' : 'Candidate User'} />
                  <Input label="Email đăng nhập" type="email" defaultValue={role === 'recruiter' ? 'recruiter@test.com' : 'candidate@test.com'} disabled />
                  
                  <div style={{ borderTop: '1px solid var(--border)', paddingTop: '24px', marginTop: '8px' }}>
                    <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '16px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <Key size={16} /> Đổi mật khẩu
                    </h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                      <Input label="Mật khẩu hiện tại" type="password" />
                      <Input label="Mật khẩu mới" type="password" />
                      <Input label="Xác nhận mật khẩu mới" type="password" />
                    </div>
                  </div>

                  <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px' }}>
                    <Button><Save size={16} /> Lưu thay đổi</Button>
                  </div>
                </div>
              </form>
            </Card>
          )}

          {activeTab === 'notifications' && (
            <Card title="Thông báo & Cảnh báo">
              <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                <div>
                  <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '12px' }}>Qua Email</h3>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                      <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                      <span className="text-body">Nhận email thông báo khi có lịch phỏng vấn mới</span>
                    </label>
                    {role === 'recruiter' && (
                      <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                        <span className="text-body">Nhận email khi AI Report đã xử lý xong</span>
                      </label>
                    )}
                    {role === 'candidate' && (
                      <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                        <span className="text-body">Nhận email nhắc nhở trước 1 tiếng khi diễn ra phỏng vấn</span>
                      </label>
                    )}
                  </div>
                </div>

                <div style={{ borderTop: '1px solid var(--border)', paddingTop: '20px' }}>
                  <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '12px' }}>Thông báo đẩy (Push Notifications)</h3>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                      <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                      <span className="text-body">Hiển thị thông báo trên trình duyệt (Browser push)</span>
                    </label>
                  </div>
                </div>

                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px' }}>
                  <Button onClick={handleSave}><Save size={16} /> Lưu tùy chọn</Button>
                </div>
              </div>
            </Card>
          )}

          {activeTab === 'privacy' && (
            <Card title="Quyền riêng tư & Bảo mật">
              <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                
                {role === 'candidate' && (
                  <div>
                    <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '12px' }}>Hiển thị hồ sơ</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                        <span className="text-body">Cho phép các nhà tuyển dụng khác xem hồ sơ của tôi (Public Profile)</span>
                      </label>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                        <span className="text-body">Chia sẻ ẩn danh kết quả Mock Interview để cải thiện AI</span>
                      </label>
                    </div>
                  </div>
                )}

                {role === 'recruiter' && (
                  <div>
                    <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '12px' }}>Bảo mật dữ liệu công ty</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                        <span className="text-body">Mã hóa ghi âm/video các cuộc phỏng vấn (E2E Encryption)</span>
                      </label>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                        <input type="checkbox" style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                        <span className="text-body">Yêu cầu xác thực 2 bước (2FA) khi đăng nhập nội bộ</span>
                      </label>
                    </div>
                  </div>
                )}

                <div style={{ borderTop: '1px solid var(--border)', paddingTop: '20px' }}>
                  <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '12px', color: 'var(--danger)' }}>Quản lý dữ liệu</h3>
                  <div style={{ display: 'flex', gap: '12px' }}>
                    <Button variant="secondary" onClick={() => setToast({ type: 'info', message: 'Đang chuẩn bị dữ liệu xuất...' })}>Xuất toàn bộ dữ liệu (Export Data)</Button>
                    <Button variant="ghost" style={{ color: 'var(--danger)', borderColor: 'var(--danger)' }} onClick={handleDeleteAccount}>Xóa tài khoản</Button>
                  </div>
                </div>

                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px' }}>
                  <Button onClick={handleSave}><Save size={16} /> Lưu tùy chọn</Button>
                </div>
              </div>
            </Card>
          )}

          {activeTab === 'workspace' && role === 'recruiter' && (
            <Card title="Cấu hình AI Workspace">
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <p className="text-body" style={{ color: 'var(--text-secondary)', marginBottom: '8px' }}>Tuỳ chỉnh cách AI Assistant hoạt động trong không gian làm việc của công ty bạn.</p>
                
                <div className="input-group">
                  <label className="input-label">Mô hình AI mặc định</label>
                  <select className="input-field" defaultValue="gpt-4">
                    <option value="gpt-4">GPT-4 (Độ chính xác cao nhất)</option>
                    <option value="gpt-35">GPT-3.5 Turbo (Nhanh nhất)</option>
                    <option value="claude">Claude 3 Haiku (Tối ưu hội thoại)</option>
                  </select>
                </div>

                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', marginTop: '12px' }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                    <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                    <span className="text-body">Tự động bóc tách JD thành Rubric</span>
                  </label>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }}>
                    <input type="checkbox" defaultChecked style={{ width: '16px', height: '16px', accentColor: 'var(--primary)' }} />
                    <span className="text-body">Tự động Suggest câu hỏi follow-up trong lúc phỏng vấn</span>
                  </label>
                </div>

                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '16px' }}>
                  <Button onClick={handleSave}><Save size={16} /> Lưu cấu hình</Button>
                </div>
              </div>
            </Card>
          )}

        </div>

      </div>

      <Modal isOpen={showDeleteModal} onClose={() => setShowDeleteModal(false)} title="Xác nhận xóa tài khoản">
        <p className="text-body" style={{ marginBottom: '24px' }}>Bạn có chắc chắn muốn xóa tài khoản này? Hành động này không thể hoàn tác và toàn bộ dữ liệu của bạn sẽ bị xóa vĩnh viễn khỏi hệ thống.</p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
          <Button variant="ghost" onClick={() => setShowDeleteModal(false)}>Huỷ</Button>
          <Button variant="primary" style={{ backgroundColor: 'var(--danger)', color: 'white', borderColor: 'var(--danger)' }} onClick={confirmDeleteAccount}>Xác nhận xóa</Button>
        </div>
      </Modal>
    </div>
  );
}
