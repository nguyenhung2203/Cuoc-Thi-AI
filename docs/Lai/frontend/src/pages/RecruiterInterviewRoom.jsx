import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Badge } from '../components/Badge';
import { Modal } from '../components/Modal';
import { Toast } from '../components/Toast';
import { Mic, MicOff, Video, VideoOff, MonitorUp, MessageSquare, PhoneOff, Sparkles, CheckCircle, AlertTriangle } from 'lucide-react';

export function RecruiterInterviewRoom() {
  const navigate = useNavigate();
  const [micOn, setMicOn] = useState(true);
  const [cameraOn, setCameraOn] = useState(true);
  const [activeTab, setActiveTab] = useState('assistant');
  const [transcript, setTranscript] = useState([
    { speaker: 'Candidate', text: 'Vâng, em đã sử dụng React khoảng 3 năm trong các dự án thực tế.', time: '14:02' }
  ]);

  // Simulate incoming transcript
  useEffect(() => {
    const timer = setTimeout(() => {
      setTranscript(prev => [...prev, { speaker: 'AI', text: '[AI Phân tích] Câu trả lời khá tự tin, tuy nhiên chưa nêu rõ dự án cụ thể.', time: '14:03', isAi: true }]);
    }, 4000);
    return () => clearTimeout(timer);
  }, []);

  const [showEndModal, setShowEndModal] = useState(false);
  const [isEnding, setIsEnding] = useState(false);

  const location = useLocation();
  const [entryToast, setEntryToast] = useState(location.state?.message ? { type: 'success', message: location.state.message } : null);

  useEffect(() => {
    if (location.state?.message) {
      window.history.replaceState({}, document.title);
    }
  }, [location]);

  const handleEndCall = () => {
    setShowEndModal(true);
  };

  const confirmEndCall = () => {
    setShowEndModal(false);
    setIsEnding(true);
    setTimeout(() => {
      navigate('/reports', { state: { message: 'Đã lưu kết quả phỏng vấn thành công' } });
    }, 1500);
  };

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column', backgroundColor: 'var(--background)' }}>
      {entryToast && (
        <Toast type={entryToast.type} message={entryToast.message} onClose={() => setEntryToast(null)} />
      )}
      {isEnding && (
        <Toast type="success" message="Đã kết thúc phỏng vấn. Đang lưu kết quả..." duration={1500} />
      )}

      {/* Header */}
      <div style={{ height: '64px', backgroundColor: 'var(--surface)', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <div>
            <h1 className="text-h2">Frontend Developer - Nguyễn Văn A</h1>
            <div style={{ display: 'flex', gap: '12px', marginTop: '4px' }}>
              <span className="text-helper" style={{ color: 'var(--danger)', display: 'flex', alignItems: 'center', gap: '4px' }}>
                <div style={{ width: '8px', height: '8px', borderRadius: '50%', backgroundColor: 'var(--danger)' }}></div> Đang ghi âm & Transcript
              </span>
              <span className="text-helper" style={{ color: 'var(--primary)' }}>00:15:32</span>
            </div>
          </div>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <Badge type="info"><Sparkles size={12} style={{ marginRight: '4px' }} /> AI Active</Badge>
          <Button variant="secondary" onClick={() => navigate('/dashboard')}>Rời tạm thời</Button>
        </div>
      </div>

      {/* Main Content */}
      <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
        
        {/* Left: Video Area (70%) */}
        <div style={{ flex: '7', display: 'flex', flexDirection: 'column', padding: '24px', gap: '16px' }}>
          {/* Video Grid */}
          <div style={{ flex: 1, display: 'flex', gap: '16px', position: 'relative' }}>
            {/* Candidate Video (Large) */}
            <div style={{ flex: 1, backgroundColor: '#1E293B', borderRadius: 'var(--radius-lg)', display: 'flex', alignItems: 'center', justifyContent: 'center', position: 'relative', overflow: 'hidden' }}>
              <div style={{ textAlign: 'center', color: 'white' }}>
                <div style={{ width: '80px', height: '80px', borderRadius: '50%', backgroundColor: 'rgba(255,255,255,0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 16px', fontSize: '24px', fontWeight: 'bold' }}>N</div>
                <div style={{ fontSize: '16px' }}>Candidate Camera Feed</div>
              </div>
              <div style={{ position: 'absolute', bottom: '16px', left: '16px', backgroundColor: 'rgba(0,0,0,0.6)', color: 'white', padding: '4px 12px', borderRadius: 'var(--radius)', fontSize: '13px' }}>
                Nguyễn Văn A (Ứng viên)
              </div>
            </div>
            
            {/* Recruiter Video (Small/PiP or Side) */}
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
            <Button variant="secondary" style={{ width: '48px', height: '48px', borderRadius: '50%', padding: 0 }}>
              <MonitorUp size={20} />
            </Button>
            <Button variant="secondary" style={{ width: '48px', height: '48px', borderRadius: '50%', padding: 0 }}>
              <MessageSquare size={20} />
            </Button>
            <div style={{ width: '1px', height: '32px', backgroundColor: 'var(--border)', margin: '0 8px' }}></div>
            <Button style={{ backgroundColor: 'var(--danger)', color: 'white', border: 'none', height: '48px', padding: '0 24px', borderRadius: '24px' }} onClick={handleEndCall}>
              <PhoneOff size={20} /> Kết thúc
            </Button>
          </div>
        </div>

        {/* Right: AI Panel (30%) */}
        <div style={{ flex: '3', backgroundColor: 'var(--surface)', borderLeft: '1px solid var(--border)', display: 'flex', flexDirection: 'column' }}>
          {/* Tabs */}
          <div style={{ display: 'flex', borderBottom: '1px solid var(--border)' }}>
            <div className={`nav-item ${activeTab === 'assistant' ? 'active' : ''}`} style={{ flex: 1, justifyContent: 'center', cursor: 'pointer', padding: '16px 0' }} onClick={() => setActiveTab('assistant')}>
              <Sparkles size={16} /> AI Assistant
            </div>
            <div className={`nav-item ${activeTab === 'rubric' ? 'active' : ''}`} style={{ flex: 1, justifyContent: 'center', cursor: 'pointer', padding: '16px 0' }} onClick={() => setActiveTab('rubric')}>
              <CheckCircle size={16} /> Tiêu chí
            </div>
            <div className={`nav-item ${activeTab === 'transcript' ? 'active' : ''}`} style={{ flex: 1, justifyContent: 'center', cursor: 'pointer', padding: '16px 0' }} onClick={() => setActiveTab('transcript')}>
              <MessageSquare size={16} /> Transcript
            </div>
          </div>

          {/* Tab Content */}
          <div style={{ flex: 1, overflowY: 'auto', padding: '24px' }}>
            {activeTab === 'assistant' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <Card style={{ backgroundColor: 'rgba(8, 145, 178, 0.05)', borderColor: 'rgba(8, 145, 178, 0.2)' }}>
                  <div style={{ display: 'flex', gap: '8px', alignItems: 'flex-start', color: 'var(--accent)' }}>
                    <Sparkles size={18} style={{ marginTop: '2px' }} />
                    <div>
                      <h4 className="text-body" style={{ fontWeight: 600, marginBottom: '8px' }}>Phân tích Real-time</h4>
                      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                          <span className="text-helper">Mức độ tự tin:</span>
                          <span className="text-helper" style={{ fontWeight: 600, color: 'var(--success)' }}>Tốt (85%)</span>
                        </div>
                        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                          <span className="text-helper">Độ chi tiết:</span>
                          <span className="text-helper" style={{ fontWeight: 600, color: 'var(--warning)' }}>Cần đào sâu hơn</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </Card>

                <div>
                  <h3 className="text-body" style={{ fontWeight: 600, marginBottom: '16px' }}>Gợi ý câu hỏi tiếp theo</h3>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    <div style={{ padding: '12px', border: '1px solid var(--border)', borderRadius: 'var(--radius)', backgroundColor: 'var(--surface)' }}>
                      <p className="text-body" style={{ marginBottom: '12px' }}>"Bạn có thể kể một ví dụ cụ thể về việc tối ưu performance trong dự án React bạn vừa nhắc đến không?"</p>
                      <Button variant="secondary" style={{ width: '100%', fontSize: '13px', height: '32px' }}>Hỏi ngay</Button>
                    </div>
                    <div style={{ padding: '12px', border: '1px solid var(--border)', borderRadius: 'var(--radius)', backgroundColor: 'var(--surface)' }}>
                      <p className="text-body" style={{ marginBottom: '12px' }}>"Trong dự án đó, bạn xử lý state management như thế nào? Dùng Redux hay Context?"</p>
                      <Button variant="secondary" style={{ width: '100%', fontSize: '13px', height: '32px' }}>Hỏi ngay</Button>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'rubric' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                    <span className="text-body" style={{ fontWeight: 500 }}>Technical Knowledge (React)</span>
                    <span className="text-body" style={{ color: 'var(--text-muted)' }}>--/10</span>
                  </div>
                  <div style={{ height: '6px', backgroundColor: 'var(--surface-soft)', borderRadius: '3px' }}></div>
                  <p className="text-helper" style={{ marginTop: '4px' }}>AI đang thu thập thêm bằng chứng...</p>
                </div>
                <div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                    <span className="text-body" style={{ fontWeight: 500 }}>Communication</span>
                    <span className="text-body" style={{ color: 'var(--success)' }}>8/10</span>
                  </div>
                  <div style={{ height: '6px', backgroundColor: 'var(--surface-soft)', borderRadius: '3px' }}>
                    <div style={{ height: '100%', width: '80%', backgroundColor: 'var(--success)', borderRadius: '3px' }}></div>
                  </div>
                  <p className="text-helper" style={{ marginTop: '4px', color: 'var(--text-muted)' }}>Bằng chứng: Trả lời lưu loát, không vấp váp.</p>
                </div>
              </div>
            )}

            {activeTab === 'transcript' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                {transcript.map((item, idx) => (
                  <div key={idx} style={{ padding: '12px', backgroundColor: item.isAi ? 'rgba(8, 145, 178, 0.05)' : 'var(--surface-soft)', borderRadius: 'var(--radius)', border: item.isAi ? '1px dashed rgba(8, 145, 178, 0.3)' : 'none' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                      <span className="text-helper" style={{ fontWeight: 600, color: item.isAi ? 'var(--accent)' : 'var(--primary)' }}>
                        {item.isAi && <Sparkles size={12} style={{ marginRight: '4px', verticalAlign: 'middle' }} />}
                        {item.speaker}
                      </span>
                      <span className="text-helper" style={{ fontSize: '11px' }}>{item.time}</span>
                    </div>
                    <p className="text-body" style={{ color: item.isAi ? 'var(--accent)' : 'var(--text-main)', fontStyle: item.isAi ? 'italic' : 'normal' }}>
                      {item.text}
                    </p>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      <Modal isOpen={showEndModal} onClose={() => setShowEndModal(false)} title="Kết thúc phỏng vấn">
        <p className="text-body" style={{ marginBottom: '24px' }}>Bạn có chắc chắn muốn kết thúc buổi phỏng vấn này?</p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
          <Button variant="ghost" onClick={() => setShowEndModal(false)}>Huỷ</Button>
          <Button variant="primary" onClick={confirmEndCall} style={{ backgroundColor: 'var(--danger)', color: 'white', borderColor: 'var(--danger)' }}>OK</Button>
        </div>
      </Modal>
    </div>
  );
}
