import React from 'react';

export function Badge({ children, type = 'neutral', className = '' }) {
  return (
    <span className={`badge badge-${type} ${className}`}>
      {children}
    </span>
  );
}
