import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Badge } from '../components/Badge';
import { Modal } from '../components/Modal';
import { Toast } from '../components/Toast';
import { Mic, ArrowRight, Play, Square, RefreshCw, Send, CheckCircle } from 'lucide-react';

export function MockInterviewRoom() {
  const navigate = useNavigate();
  const [recording, setRecording] = useState(false);
  const [questionIndex, setQuestionIndex] = useState(1);
  const [analyzing, setAnalyzing] = useState(false);
  const [feedback, setFeedback] = useState(null);
  const [textAnswer, setTextAnswer] = useState('');
  const [showEndModal, setShowEndModal] = useState(false);

  const location = useLocation();
  const [entryToast, setEntryToast] = useState(location.state?.message ? { type: 'success', message: location.state.message } : null);

  useEffect(() => {
    if (location.state?.message) {
      window.history.replaceState({}, document.title);
    }
  }, [location]);

  const questions = [
    "Hãy giới thiệu ngắn gọn về bản thân và kinh nghiệm làm việc của bạn.",
    "Bạn đã từng sử dụng React trong dự án nào? Hãy mô tả một thử thách lớn nhất bạn gặp phải.",
    "Làm thế nào để bạn quản lý state trong một ứng dụng React lớn?"
  ];

  const handleRecord = () => {
    if (!recording) {
      setRecording(true);
      setFeedback(null);
    } else {
      setRecording(false);
      setAnalyzing(true);
      setTimeout(() => {
        setAnalyzing(false);
        setFeedback({
          score: 'Tốt',
          message: 'Câu trả lời rõ ràng, cấu trúc tốt. Đã đề cập được số năm kinh nghiệm và công nghệ chính.',
          improvement: 'Có thể thêm một ví dụ ngắn về dự án gần nhất để tăng tính thuyết phục.'
        });
      }, 2000);
    }
  };

  const handleSendText = () => {
    if (!textAnswer) return;
    setAnalyzing(true);
    setFeedback(null);
    setTimeout(() => {
      setAnalyzing(false);
      setFeedback({
        score: 'Khá',
        message: 'Bạn đã nêu được các ý chính, câu văn mạch lạc.',
        improvement: 'Cố gắng trả lời bằng giọng nói để rèn luyện sự tự tin tốt hơn nhé!'
      });
    }, 1500);
  };

  const handleNext = () => {
    if (questionIndex < questions.length) {
      setQuestionIndex(prev => prev + 1);
      setFeedback(null);
      setTextAnswer('');
    } else {
      navigate('/mock-results', { state: { message: 'Hoàn thành bài thi thử!' } });
    }
  };

  return (
    <div style={{ height: '100vh', display: 'flex', flexDirection: 'column', backgroundColor: 'var(--background)' }}>
      {/* Header */}
      <div style={{ height: '64px', backgroundColor: 'var(--surface)', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <div>
            <h1 className="text-h2">Luyện tập AI: Frontend Developer</h1>
            <p className="text-helper" style={{ marginTop: '4px' }}>Câu hỏi {questionIndex} / {questions.length}</p>
          </div>
        </div>
        <Button variant="ghost" style={{ color: 'var(--danger)' }} onClick={() => setShowEndModal(true)}>Kết thúc sớm</Button>
      </div>

      <Modal isOpen={showEndModal} onClose={() => setShowEndModal(false)} title="Kết thúc sớm">
        <p className="text-body" style={{ marginBottom: '24px' }}>Bạn có chắc chắn muốn kết thúc bài thi sớm? Kết quả sẽ được tính trên những câu bạn đã trả lời.</p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '12px' }}>
          <Button variant="ghost" onClick={() => setShowEndModal(false)}>Huỷ</Button>
          <Button variant="primary" style={{ backgroundColor: 'var(--danger)', color: 'white', borderColor: 'var(--danger)' }} onClick={() => navigate('/home', { state: { message: 'Đã hủy bài thi thử' } })}>Xác nhận</Button>
        </div>
      </Modal>

      {entryToast && (
        <Toast type={entryToast.type} message={entryToast.message} onClose={() => setEntryToast(null)} />
      )}

      {/* Main Content */}
      <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
        
        {/* Left: AI Interviewer */}
        <div style={{ flex: '4', display: 'flex', flexDirection: 'column', padding: '24px', gap: '24px', borderRight: '1px solid var(--border)', backgroundColor: 'var(--surface)' }}>
          {/* AI Avatar Area */}
          <div style={{ flex: 1, backgroundColor: 'var(--surface-soft)', borderRadius: 'var(--radius-lg)', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', border: '1px solid var(--border)', position: 'relative' }}>
            <div style={{ width: '120px', height: '120px', borderRadius: '50%', backgroundColor: 'var(--accent)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white', marginBottom: '24px', boxShadow: '0 0 0 12px rgba(8, 145, 178, 0.1)' }}>
              <span style={{ fontSize: '32px', fontWeight: 'bold' }}>AI</span>
            </div>
            <h3 className="text-h2" style={{ marginBottom: '8px' }}>AI Interviewer</h3>
            
            <div style={{ padding: '8px 16px', backgroundColor: analyzing ? 'rgba(8, 145, 178, 0.1)' : 'var(--surface)', borderRadius: '20px', border: '1px solid var(--border)', display: 'flex', alignItems: 'center', gap: '8px' }}>
              {analyzing ? (
                <><RefreshCw size={14} className="spin" color="var(--accent)" /> <span className="text-helper" style={{ color: 'var(--accent)' }}>Đang phân tích câu trả lời...</span></>
              ) : (
                <><Play size={14} color="var(--success)" /> <span className="text-helper">Sẵn sàng lắng nghe</span></>
              )}
            </div>
          </div>

          {/* Current Question */}
          <Card style={{ backgroundColor: 'var(--surface)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '12px' }}>
              <Badge type="info">Câu hỏi hiện tại</Badge>
              <Button variant="ghost" style={{ padding: '4px 8px', height: 'auto' }}><Play size={14} style={{ marginRight: '4px' }}/> Nghe lại</Button>
            </div>
            <p className="text-body" style={{ fontSize: '16px', lineHeight: 1.6, fontWeight: 500 }}>
              {questions[questionIndex - 1]}
            </p>
          </Card>
        </div>

        {/* Right: Candidate Answer */}
        <div style={{ flex: '6', display: 'flex', flexDirection: 'column', padding: '24px', gap: '24px' }}>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px', flex: 1 }}>
            {/* Audio Recording UI */}
            <Card style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', padding: '48px 24px', flex: 1 }}>
              <div 
                style={{ 
                  width: '80px', height: '80px', borderRadius: '50%', 
                  backgroundColor: recording ? 'rgba(239, 68, 68, 0.1)' : 'rgba(37, 99, 235, 0.1)', 
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  cursor: 'pointer', transition: 'all 0.3s',
                  border: `2px solid ${recording ? 'var(--danger)' : 'var(--primary)'}`
                }}
                onClick={handleRecord}
              >
                {recording ? <Square size={32} color="var(--danger)" /> : <Mic size={32} color="var(--primary)" />}
              </div>
              <h3 className="text-body" style={{ marginTop: '24px', fontWeight: 600 }}>
                {recording ? 'Đang thu âm...' : 'Nhấn để trả lời bằng giọng nói'}
              </h3>
              <p className="text-helper" style={{ marginTop: '8px' }}>
                {recording ? '00:15' : 'Tối đa 3 phút cho mỗi câu trả lời'}
              </p>

              {/* Fake Waveform when recording */}
              {recording && (
                <div style={{ display: 'flex', gap: '4px', marginTop: '24px', height: '32px', alignItems: 'center' }}>
                  {[1, 2, 3, 4, 5, 4, 3, 2, 1, 2, 4, 5, 3, 2].map((h, i) => (
                    <div key={i} style={{ width: '4px', height: `${h * 20}%`, backgroundColor: 'var(--danger)', borderRadius: '2px', animation: `pulse ${0.5 + (i%3)*0.1}s infinite alternate` }}></div>
                  ))}
                </div>
              )}
            </Card>

            {/* OR Text input */}
            <div style={{ position: 'relative' }}>
              <textarea 
                className="input-field" 
                placeholder="Hoặc nhập câu trả lời bằng văn bản tại đây..."
                style={{ width: '100%', height: '120px' }}
                value={textAnswer}
                onChange={(e) => setTextAnswer(e.target.value)}
                disabled={recording || analyzing}
              ></textarea>
              <Button 
                style={{ position: 'absolute', bottom: '16px', right: '16px' }}
                disabled={!textAnswer || recording || analyzing}
                onClick={handleSendText}
              >
                <Send size={14} style={{ marginRight: '8px' }} /> Gửi
              </Button>
            </div>
          </div>

          {/* Feedback Area */}
          {feedback && (
            <Card style={{ borderTop: '4px solid var(--accent)', backgroundColor: 'var(--surface)', animation: 'fadeIn 0.5s ease-out' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '16px' }}>
                <h3 className="text-h2" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <CheckCircle size={20} color="var(--success)" /> Phản hồi từ AI
                </h3>
                <Badge type="success">Điểm: {feedback.score}</Badge>
              </div>
              <p className="text-body" style={{ marginBottom: '12px' }}>{feedback.message}</p>
              <div style={{ padding: '12px', backgroundColor: 'rgba(217, 119, 6, 0.05)', borderRadius: 'var(--radius)', borderLeft: '3px solid var(--warning)' }}>
                <p className="text-helper" style={{ fontWeight: 600, color: 'var(--warning)', marginBottom: '4px' }}>Gợi ý cải thiện:</p>
                <p className="text-body" style={{ color: 'var(--text-secondary)' }}>{feedback.improvement}</p>
              </div>

              <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '24px' }}>
                <Button onClick={handleNext}>
                  {questionIndex < questions.length ? 'Câu hỏi tiếp theo' : 'Xem báo cáo tổng hợp'} <ArrowRight size={16} />
                </Button>
              </div>
            </Card>
          )}

        </div>
      </div>
      <style>{`
        @keyframes pulse { 0% { height: 20%; } 100% { height: 100%; } }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
        .spin { animation: spin 2s linear infinite; }
        @keyframes spin { 100% { transform: rotate(360deg); } }
      `}</style>
    </div>
  );
}
