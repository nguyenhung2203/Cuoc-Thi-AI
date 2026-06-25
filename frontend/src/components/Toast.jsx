import React, { useState, useEffect } from 'react';
import { CheckCircle, AlertTriangle, XCircle, Info, X } from 'lucide-react';

export function Toast({ type = 'info', message, onClose, duration = 3000 }) {
  const [visible, setVisible] = useState(true);

  useEffect(() => {
    if (duration) {
      const timer = setTimeout(() => {
        handleClose();
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [duration]);

  const handleClose = () => {
    setVisible(false);
    setTimeout(() => {
      if (onClose) onClose();
    }, 300); // Wait for fade out animation
  };

  const getIcon = () => {
    switch (type) {
      case 'success': return <CheckCircle size={20} color="var(--success)" />;
      case 'error': return <XCircle size={20} color="var(--danger)" />;
      case 'warning': return <AlertTriangle size={20} color="var(--warning)" />;
      default: return <Info size={20} color="var(--primary)" />;
    }
  };

  const getBorderColor = () => {
    switch (type) {
      case 'success': return 'var(--success)';
      case 'error': return 'var(--danger)';
      case 'warning': return 'var(--warning)';
      default: return 'var(--primary)';
    }
  };

  if (!visible && !onClose) return null;

  return (
    <div 
      style={{ 
        position: 'fixed', 
        bottom: '24px', 
        right: '24px', 
        backgroundColor: 'var(--surface)', 
        boxShadow: 'var(--shadow-md)', 
        border: '1px solid var(--border)',
        borderLeft: `4px solid ${getBorderColor()}`,
        borderRadius: 'var(--radius)', 
        padding: '16px', 
        display: 'flex', 
        alignItems: 'flex-start', 
        gap: '12px',
        zIndex: 9999,
        animation: visible ? 'slideInRight 0.3s ease-out forwards' : 'fadeOutRight 0.3s ease-in forwards',
        minWidth: '300px',
        maxWidth: '400px'
      }}
    >
      <div style={{ marginTop: '2px' }}>{getIcon()}</div>
      <div style={{ flex: 1 }}>
        <p className="text-body" style={{ fontWeight: 500, lineHeight: 1.4 }}>{message}</p>
      </div>
      <button 
        onClick={handleClose} 
        style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-muted)' }}
      >
        <X size={16} />
      </button>

      <style>{`
        @keyframes slideInRight {
          from { transform: translateX(100%); opacity: 0; }
          to { transform: translateX(0); opacity: 1; }
        }
        @keyframes fadeOutRight {
          from { transform: translateX(0); opacity: 1; }
          to { transform: translateX(100%); opacity: 0; }
        }
      `}</style>
    </div>
  );
}
