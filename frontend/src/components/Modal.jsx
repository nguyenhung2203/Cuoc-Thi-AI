import React from 'react';
import { X } from 'lucide-react';
import { Button } from './Button';

export function Modal({ isOpen, onClose, title, children }) {
  if (!isOpen) return null;

  return (
    <div style={{
      position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
      backgroundColor: 'rgba(0,0,0,0.5)', display: 'flex',
      alignItems: 'center', justifyContent: 'center', zIndex: 50
    }}>
      <div style={{
        backgroundColor: 'var(--surface)', padding: '1.5rem',
        borderRadius: 'var(--radius)', minWidth: '400px',
        boxShadow: 'var(--shadow-lg)'
      }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h2 style={{ fontSize: '1.25rem', fontWeight: 600 }}>{title}</h2>
          <Button variant="ghost" onClick={onClose} style={{ padding: '0.25rem' }}>
            <X size={20} />
          </Button>
        </div>
        <div>
          {children}
        </div>
      </div>
    </div>
  );
}
