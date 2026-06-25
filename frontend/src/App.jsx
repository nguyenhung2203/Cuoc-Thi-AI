import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout } from './layouts/MainLayout';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { ForgotPassword } from './pages/ForgotPassword';
import { RecruiterDashboard } from './pages/RecruiterDashboard';
import { CandidateHome } from './pages/CandidateHome';
import { Jobs } from './pages/Jobs';
import { JobDetail } from './pages/JobDetail';
import { Candidates } from './pages/Candidates';
import { CandidateDetail } from './pages/CandidateDetail';
import { Interviews } from './pages/Interviews';
import { CreateInterview } from './pages/CreateInterview';
import { InterviewDetail } from './pages/InterviewDetail';
import { RecruiterInterviewRoom } from './pages/RecruiterInterviewRoom';
import { CandidateInterviewRoom } from './pages/CandidateInterviewRoom';
import { MockInterviewRoom } from './pages/MockInterviewRoom';
import { InterviewReport } from './pages/InterviewReport';
import { MockSetup } from './pages/MockSetup';
import { MockResults } from './pages/MockResults';
import { MockResultDetail } from './pages/MockResultDetail';
import { CandidateProfile } from './pages/CandidateProfile';
import { QuestionBank } from './pages/QuestionBank';
import { Settings } from './pages/Settings';
import { InterviewConsent } from './pages/InterviewConsent';
import { InterviewExpired } from './pages/InterviewExpired';

// Protected Route wrapper
function ProtectedRoute({ children, allowedRoles }) {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role');

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && !allowedRoles.includes(role)) {
    return <Navigate to={role === 'recruiter' ? '/dashboard' : '/home'} replace />;
  }

  return children;
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/login" replace />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/forgot-password" element={<ForgotPassword />} />
      
      {/* Rooms without sidebar */}
      <Route path="/recruiter-room" element={<ProtectedRoute allowedRoles={['recruiter']}><RecruiterInterviewRoom /></ProtectedRoute>} />
      <Route path="/candidate-room" element={<ProtectedRoute allowedRoles={['candidate']}><CandidateInterviewRoom /></ProtectedRoute>} />
      <Route path="/mock-room" element={<ProtectedRoute allowedRoles={['candidate']}><MockInterviewRoom /></ProtectedRoute>} />
      <Route path="/interview-consent" element={<ProtectedRoute allowedRoles={['candidate']}><InterviewConsent /></ProtectedRoute>} />
      <Route path="/interview-expired" element={<InterviewExpired />} />

      <Route element={<MainLayout />}>
        {/* Shared Routes */}
        <Route path="/settings" element={<ProtectedRoute><Settings /></ProtectedRoute>} />

        {/* Recruiter Routes */}
        <Route path="/dashboard" element={<ProtectedRoute allowedRoles={['recruiter']}><RecruiterDashboard /></ProtectedRoute>} />
        <Route path="/jobs" element={<ProtectedRoute allowedRoles={['recruiter']}><Jobs /></ProtectedRoute>} />
        <Route path="/jobs/:id" element={<ProtectedRoute allowedRoles={['recruiter']}><JobDetail /></ProtectedRoute>} />
        <Route path="/candidates" element={<ProtectedRoute allowedRoles={['recruiter']}><Candidates /></ProtectedRoute>} />
        <Route path="/candidates/:id" element={<ProtectedRoute allowedRoles={['recruiter']}><CandidateDetail /></ProtectedRoute>} />
        <Route path="/interviews" element={<ProtectedRoute allowedRoles={['recruiter']}><Interviews /></ProtectedRoute>} />
        <Route path="/interviews/new" element={<ProtectedRoute allowedRoles={['recruiter']}><CreateInterview /></ProtectedRoute>} />
        <Route path="/interviews/:id" element={<ProtectedRoute allowedRoles={['recruiter']}><InterviewDetail /></ProtectedRoute>} />
        <Route path="/reports" element={<ProtectedRoute allowedRoles={['recruiter']}><InterviewReport /></ProtectedRoute>} />
        <Route path="/questions" element={<ProtectedRoute allowedRoles={['recruiter']}><QuestionBank /></ProtectedRoute>} />
        
        {/* Candidate Routes */}
        <Route path="/home" element={<ProtectedRoute allowedRoles={['candidate']}><CandidateHome /></ProtectedRoute>} />
        <Route path="/mock-setup" element={<ProtectedRoute allowedRoles={['candidate']}><MockSetup /></ProtectedRoute>} />
        <Route path="/mock-results" element={<ProtectedRoute allowedRoles={['candidate']}><MockResults /></ProtectedRoute>} />
        <Route path="/mock-results/:id" element={<ProtectedRoute allowedRoles={['candidate']}><MockResultDetail /></ProtectedRoute>} />
        <Route path="/profile" element={<ProtectedRoute allowedRoles={['candidate']}><CandidateProfile /></ProtectedRoute>} />
      </Route>
    </Routes>
  );
}

export default App;
