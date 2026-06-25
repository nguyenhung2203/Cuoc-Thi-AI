import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card } from '../components/Card';
import { Button } from '../components/Button';
import { Badge } from '../components/Badge';
import { ArrowLeft, MessageCircle, CheckCircle, AlertCircle } from 'lucide-react';

export function MockResultDetail() {
  const { id } = useParams();
  const navigate = useNavigate();

  // Mock data cho một kết quả chi tiết
  const sessionData = {
    role: 'Frontend Developer',
    level: 'Middle',
    date: '20/10/2023',
    overallScore: 8.5,
    summaryFeedback: 'Bạn đã làm rất tốt trong việc giải thích các khái niệm cốt lõi của React. Tuy nhiên, phần System Design cần cụ thể hơn về cách scale ứng dụng.',
    questions: [
      {
        id: 1,
        question: 'Bạn hãy giải thích cơ chế Virtual DOM trong React và tại sao nó lại giúp tăng hiệu suất?',
        candidateAnswer: 'Virtual DOM là một bản copy của Real DOM. Khi state thay đổi, React tạo ra một Virtual DOM mới, so sánh với cái cũ (diffing), và chỉ cập nhật những node bị thay đổi lên Real DOM.',
        score: 9,
        feedback: 'Câu trả lời rất chính xác, ngắn gọn và đi thẳng vào trọng tâm. Bạn có thể bổ sung thêm về quá trình Reconciliation để đạt điểm tuyệt đối.',
        goodPoints: ['Hiểu rõ khái niệm bản copy', 'Nắm được quá trình diffing'],
        improvePoints: ['Thiếu key term Reconciliation']
      },
      {
        id: 2,
        question: 'Làm thế nào để tối ưu hóa hiệu suất (performance) của một ứng dụng React lớn?',
        candidateAnswer: 'Tôi thường dùng useMemo và useCallback để tránh re-render. Ngoài ra cũng dùng React.lazy để code splitting.',
        score: 7.5,
        feedback: 'Các ý chính đều đúng, tuy nhiên bạn cần giải thích rõ HƯỚNG áp dụng thực tế thay vì chỉ liệt kê hooks. Khi nào KHÔNG NÊN dùng useMemo cũng là một ý quan trọng.',
        goodPoints: ['Đề cập đúng các công cụ tối ưu (useMemo, React.lazy)'],
        improvePoints: ['Cần ví dụ thực tế', 'Thiếu cân nhắc trade-off khi lạm dụng useMemo']
      }
    ]
  };

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '32px' }}>
        <Button variant="ghost" onClick={() => navigate('/mock-results')} style={{ padding: '8px' }}>
          <ArrowLeft size={20} />
        </Button>
        <div>
          <h1 className="text-h1">Chi tiết kết quả luyện tập</h1>
          <p className="text-helper" style={{ marginTop: '4px' }}>{sessionData.role} - Cấp độ {sessionData.level} - Ngày {sessionData.date}</p>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '24px', marginBottom: '24px' }}>
        <Card>
          <div style={{ textAlign: 'center', marginBottom: '24px' }}>
            <div style={{ display: 'inline-flex', alignItems: 'center', justifyContent: 'center', width: '80px', height: '80px', borderRadius: '50%', backgroundColor: 'rgba(22, 163, 74, 0.1)', color: 'var(--success)', fontSize: '28px', fontWeight: 'bold', marginBottom: '16px' }}>
              {sessionData.overallScore}
            </div>
            <h2 className="text-body" style={{ fontWeight: 600 }}>Điểm tổng kết</h2>
          </div>
          
          <div style={{ paddingTop: '16px', borderTop: '1px solid var(--border)' }}>
            <h3 className="text-helper" style={{ fontWeight: 600, marginBottom: '8px' }}>Nhận xét chung:</h3>
            <p className="text-body" style={{ color: 'var(--text-secondary)' }}>{sessionData.summaryFeedback}</p>
          </div>
        </Card>

        <Card title="Đánh giá chi tiết từng câu hỏi">
          <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
            {sessionData.questions.map((q, index) => (
              <div key={q.id} style={{ paddingBottom: '24px', borderBottom: index < sessionData.questions.length - 1 ? '1px solid var(--border)' : 'none' }}>
                <div style={{ display: 'flex', gap: '12px', marginBottom: '12px' }}>
                  <MessageCircle size={20} color="var(--primary)" style={{ flexShrink: 0, marginTop: '2px' }} />
                  <div>
                    <h4 className="text-body" style={{ fontWeight: 600 }}>Câu hỏi {index + 1}: {q.question}</h4>
                  </div>
                </div>

                <div style={{ marginLeft: '32px', marginBottom: '16px', padding: '12px', backgroundColor: 'var(--surface-soft)', borderRadius: '8px' }}>
                  <p className="text-helper" style={{ fontWeight: 600, marginBottom: '4px', color: 'var(--text-main)' }}>Câu trả lời của bạn:</p>
                  <p className="text-body" style={{ color: 'var(--text-secondary)' }}>"{q.candidateAnswer}"</p>
                </div>

                <div style={{ marginLeft: '32px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                    <p className="text-helper" style={{ fontWeight: 600, color: 'var(--primary)' }}>AI Feedback</p>
                    <Badge type={q.score >= 8 ? 'success' : 'warning'}>Điểm: {q.score}/10</Badge>
                  </div>
                  <p className="text-body" style={{ marginBottom: '12px' }}>{q.feedback}</p>
                  
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                    <div style={{ padding: '12px', border: '1px solid #bbf7d0', backgroundColor: '#f0fdf4', borderRadius: '8px' }}>
                      <p className="text-helper" style={{ color: '#166534', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '4px', marginBottom: '8px' }}><CheckCircle size={14} /> Điểm tốt</p>
                      <ul className="text-body" style={{ paddingLeft: '16px', color: '#166534', margin: 0 }}>
                        {q.goodPoints.map((p, i) => <li key={i}>{p}</li>)}
                      </ul>
                    </div>
                    <div style={{ padding: '12px', border: '1px solid #fef08a', backgroundColor: '#fefce8', borderRadius: '8px' }}>
                      <p className="text-helper" style={{ color: '#854d0e', fontWeight: 600, display: 'flex', alignItems: 'center', gap: '4px', marginBottom: '8px' }}><AlertCircle size={14} /> Cần cải thiện</p>
                      <ul className="text-body" style={{ paddingLeft: '16px', color: '#854d0e', margin: 0 }}>
                        {q.improvePoints.map((p, i) => <li key={i}>{p}</li>)}
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>
    </div>
  );
}
