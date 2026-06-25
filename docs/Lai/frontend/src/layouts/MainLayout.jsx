import React from 'react';
import { Outlet, Link, useNavigate, useLocation } from 'react-router-dom';
import { Home, Briefcase, Users, Calendar, LogOut, FileText, Settings, Target, Database } from 'lucide-react';
import { Toast } from '../components/Toast';
import { useState, useEffect } from 'react';

export function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const role = localStorage.getItem('role') || 'recruiter';
  const [globalToast, setGlobalToast] = useState(null);

  useEffect(() => {
    if (location.state && location.state.message) {
      setGlobalToast({ type: location.state.type || 'success', message: location.state.message });
      // Clean up the state so refresh doesn't trigger it again
      window.history.replaceState({}, document.title);
    }
  }, [location]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    navigate('/login');
  };

  const recruiterNav = [
    { name: 'Tổng quan', path: '/dashboard', icon: Home },
    { name: 'Việc làm', path: '/jobs', icon: Briefcase },
    { name: 'Ứng viên', path: '/candidates', icon: Users },
    { name: 'Lịch phỏng vấn', path: '/interviews', icon: Calendar },
    { name: 'Kho câu hỏi', path: '/questions', icon: Database },
    { name: 'Cài đặt', path: '/settings', icon: Settings }
  ];

  const candidateNav = [
    { name: 'Tổng quan', path: '/home', icon: Home },
    { name: 'Luyện phỏng vấn AI', path: '/mock-setup', icon: Target },
    { name: 'Kết quả luyện tập', path: '/mock-results', icon: FileText },
    { name: 'Hồ sơ của tôi', path: '/profile', icon: Users },
    { name: 'Cài đặt', path: '/settings', icon: Settings }
  ];

  const navItems = role === 'recruiter' ? recruiterNav : candidateNav;

  return (
    <div className="main-layout">
      <div className="sidebar">
        <div className="sidebar-header">
          <div className="sidebar-logo">Interview AI</div>
          <div className="sidebar-subtitle">{role === 'recruiter' ? 'Enterprise Tier' : 'Candidate Portal'}</div>
        </div>
        <div className="sidebar-nav">
          {navItems.map(item => (
            <Link
              key={item.path}
              to={item.path}
              className={`nav-item ${location.pathname.startsWith(item.path) ? 'active' : ''}`}
            >
              <item.icon size={18} />
              {item.name}
            </Link>
          ))}
        </div>
        <div style={{ padding: '16px', borderTop: '1px solid var(--border)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '10px', marginBottom: '12px' }}>
            <div style={{ width: '32px', height: '32px', borderRadius: '50%', backgroundColor: 'var(--primary)', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold' }}>
              {role === 'recruiter' ? 'R' : 'C'}
            </div>
            <div style={{ overflow: 'hidden' }}>
              <div style={{ fontSize: '14px', fontWeight: 500, whiteSpace: 'nowrap', textOverflow: 'ellipsis' }}>
                {role === 'recruiter' ? 'Recruiter User' : 'Candidate User'}
              </div>
              <div style={{ fontSize: '12px', color: 'var(--text-muted)' }}>
                {role === 'recruiter' ? 'HR Department' : 'Applicant'}
              </div>
            </div>
          </div>
          <button className="btn btn-ghost" onClick={handleLogout} style={{ width: '100%', justifyContent: 'flex-start', color: 'var(--danger)' }}>
            <LogOut size={18} /> Đăng xuất
          </button>
        </div>
      </div>
      <div className="content-area">
        {globalToast && (
          <Toast type={globalToast.type} message={globalToast.message} onClose={() => setGlobalToast(null)} />
        )}
        <div className="navbar">
          <span className="text-helper">{new Date().toLocaleDateString('vi-VN', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</span>
        </div>
        <div className="page-content">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
