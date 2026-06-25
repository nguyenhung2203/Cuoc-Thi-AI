import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Button } from '../components/Button';
import { Modal } from '../components/Modal';
import { Toast } from '../components/Toast';
import { Mic, MicOff, Video, VideoOff, MessageSquare, PhoneOff, CheckCircle, FileText } from 'lucide-react';

export function CandidateInterviewRoom() {
  const navigate = useNavigate();
  const [micOn, setMicOn] = useState(true);
  const [cameraOn, setCameraOn] = useState(true);
  const [activeTab, setActiveTab] = useState('info');

  const [showLeaveModal, setShowLeaveModal] = useState(false);
  const [isLeaving, setIsLeaving] = useState(false);

  const location = useLocation();
  const [entryToast, setEntryToast] = useState(location.state?.message ? { type: 'success', message: location.state.message } : null);

  useEffect(() => {
    if (location.state?.message) {
      window.history.replaceState({}, document.title);
    }
  }, [location]);

  const handleLeave = () => {
    setShowLeaveModal(true);
  };

  const confirmLeave = () => {
    setShowLeaveModal(false);
    setIsLeaving(true);
    setTimeout(() => {
      navigate('/home', { state: { message: 'Rời phòng phỏng vấn thành công' } });
    }, 1500);
  };

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column', backgroundColor: 'var(--background)' }}>
      {entryToast && (
        <Toast type={entryToast.type} message={entryToast.message} onClose={() => setEntryToast(null)} />
      )}
      {isLeaving && (
        <Toast type="info" message="Đang rời phòng phỏng vấn..." duration={1500} />
      )}

      {/* Header */}
      <div style={{ height: '64px', backgroundColor: 'var(--surface)', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <div>
            <h1 className="text-h2">Phỏng vấn: Frontend Developer</h1>
            <p className="text-helper" style={{ marginTop: '4px' }}>Công ty TechCorp Inc.</p>
          </div>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <span className="text-helper" style={{ color: 'var(--success)', display: 'flex', alignItems: 'center', gap: '4px' }}>
            <CheckCircle size={14} /> Đường truyền tốt
          </span>
          <span style={{ width: '1px', height: '24px', backgroundColor: 'var(--border)' }}></span>
          <span className="text-body" style={{ fontWeight: 600 }}>00:15:32</span>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
        
        {/* Left: Video Area (70%) */}
        <div style={{ flex: '7', display: 'flex', flexDirection: 'column', padding: '24px', gap: '16px' }}>
          
          {/* Notification Banner */}
          <div style={{ backgroundColor: 'rgba(37, 99, 235, 0.05)', border: '1px solid rgba(37, 99, 235, 0.2)', padding: '12px 16px', borderRadius: 'var(--radius)', display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--primary)', fontSize: '13px' }}>
            <span style={{ fontWeight: 600 }}>Lời nhắc:</span> Hãy trả lời tự nhiên. Nhà tuyển dụng sẽ dẫn dắt buổi phỏng vấn.
          </div>

          {/* Video Grid */}
          <div style={{ flex: 1, display: 'flex', gap: '16px', position: 'relative' }}>
            {/* Recruiter Video (Large) */}
            <div style={{ flex: 1, backgroundColor: '#1E293B', borderRadius: 'var(--radius-lg)', display: 'flex', alignItems: 'center', justifyContent: 'center', position: 'relative', overflow: 'hidden' }}>
              <div style={{ textAlign: 'center', color: 'white' }}>
                <div style={{ width: '80px', height: '80px', borderRadius: '50%', backgroundColor: 'rgba(255,255,255,0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px', fontSize: '24px', fontWeight: 'bold' }}>R</div>
                <div style={{ fontSize: '16px' }}>Recruiter Camera Feed</div>
              </div>
              <div style={{ position: 'absolute', bottom: '16px', left: '16px', backgroundColor: 'rgba(0,0,0,0.6)', color: 'white', padding: '4px 12px', borderRadius: 'var(--radius)', fontSize: '13px' }}>
                Trần Thị B (HR Manager)
              </div>
            </div>
            
            {/* Candidate Self View (Small/PiP or Side) */}
            <div style={{ width: '240px', backgroundColor: '#334155', borderRadius: 'var(--radius-lg)', display: 'flex', alignItems: 'center', justifyContent: 'center', position: 'absolute', top: '16px', right: '16px', height: '160px', boxShadow: 'var(--shadow-md)', border: '2px solid rgba(255,255,255,0.1)' }}>
               <div style={{ color: 'white', fontSize: '12px' }}>Your Camera</div>
            </div>
          </div>

          {/* Control Bar */}
          <div style={{ height: '72px', backgroundColor: 'var(--surface)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '16px' }}>
            <Button variant={micOn ? 'secondary' : 'primary'} style={{ width: '48px', height: '48px', borderRadius: '50%', padding: 0 }} onClick={() => setMicOn(!micOn)}>
              {micOn ? <Mic size={20} /> : <MicOff size={20} />}
            </Button>
            <Button variant={cameraOn ? 'secondary' : 'primary'} style={{ width: '48px', height: '48px', borderRadius: '50%', padding: 0 }} onClick={() => setCameraOn(!cameraOn)}>
              {cameraOn ? <Video size={20} /> : <VideoOff size={20} />}
            </Button>
            <div style={{ width: '1px', height: '32px', backgroundColor: 'var(--border)', margin: '0 8px' }}></div>
            <Button style={{ backgroundColor: 'var(--danger)', color: 'white', border: 'none', height: '48px', padding: '0 24px', borderRadius: '24px' }} onClick={handleLeave}>
              <PhoneOff size={20} /> Rời phòng
            </Button>
          </div>
        </div>

        {/* Right: Candidate Panel (30%) */}
        <div style={{ flex: '3', backgroundColor: 'var(--surface)', borderLeft: '1px solid var(--border)', display: 'flex', flexDirection: 'column' }}>
          {/* Tabs */}
          <div style={{ display: 'flex', borderBottom: '1px solid var(--border)' }}>
            <div className={`nav-item ${activeTab === 'info' ? 'active' : ''}`} style={{ flex: 1, justifyContent: 'center', cursor: 'pointer', padding: '16px 0' }} onClick={() => setActiveTab('info')}>
              <FileText size={16} /> Thông tin JD
            </div>
            <div className={`nav-item ${activeTab === 'chat' ? 'active' : ''}`} style={{ flex: 1, justifyContent: 'center', cursor: 'pointer', padding: '16px 0' }} onClick={() => setActiveTab('chat')}>
              <MessageSquare size={16} /> Chat
            </div>
          </div>

          {/* Tab Content */}
          <div style={{ flex: 1, overflowY: 'auto', padding: '24px' }}>
            {activeTab === 'info' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <div>
                  <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '8px' }}>Mô tả công việc (Tóm tắt)</h3>
                  <ul className="text-body" style={{ paddingLeft: '20px', color: 'var(--text-secondary)', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    <li>Phát triển các tính năng Frontend sử dụng ReactJS.</li>
                    <li>Tối ưu hóa hiệu suất ứng dụng web.</li>
                    <li>Phối hợp với UI/UX designer và Backend developer.</li>
                  </ul>
                </div>
                
                <div style={{ borderTop: '1px solid var(--border)', paddingTop: '24px' }}>
                  <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '8px' }}>Ghi chú cá nhân của bạn</h3>
                  <textarea 
                    className="input-field" 
                    placeholder="Viết nháp câu trả lời hoặc ghi chú tại đây (Người phỏng vấn không thấy)..."
                    style={{ width: '100%', height: '200px', backgroundColor: 'var(--surface-soft)' }}
                  ></textarea>
                </div>
              </div>
            )}

            {activeTab === 'chat' && (
              <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                <div style={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'flex-end', paddingBottom: '16px' }}>
                  <p className="text-helper" style={{ textAlign: 'center' }}>Chưa có tin nhắn nào.</p>
                </div>
                <div style={{ display: 'flex', gap: '8px' }}>
                  <input type="text" className="input-field" placeholder="Nhập tin nhắn..." style={{ flex: 1 }} />
                  <Button variant="primary">Gửi</Button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <Modal isOpen={showLeaveModal} onClose={() => setShowLeaveModal(false)} title="Xác nhận rời phòng">
        <p className="text-body" style={{ marginBottom: '24px' }}>Bạn có chắc chắn muốn rời khỏi phòng phỏng vấn?</p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
          <Button variant="ghost" onClick={() => setShowLeaveModal(false)}>Huỷ</Button>
          <Button variant="primary" onClick={confirmLeave} style={{ backgroundColor: 'var(--danger)', color: 'white', borderColor: 'var(--danger)' }}>OK</Button>
        </div>
      </Modal>
    </div>
  );
}
