import React from 'react';

export function Table({ columns, data }) {
  return (
    <div style={{ overflowX: 'auto', margin: '0 -24px' }}>
      <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
        <thead>
          <tr style={{ borderBottom: '1px solid var(--border)' }}>
            {columns.map((col, index) => (
              <th key={index} style={{ 
                padding: '16px 24px', 
                color: 'var(--text-muted)', 
                fontWeight: 500,
                fontSize: '12px',
                textTransform: 'uppercase',
                letterSpacing: '0.05em'
              }}>
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td colSpan={columns.length} style={{ textAlign: 'center', padding: '32px 24px', color: 'var(--text-muted)' }}>
                Không có dữ liệu
              </td>
            </tr>
          ) : (
            data.map((row, rowIndex) => (
              <tr key={rowIndex} style={{ borderBottom: '1px solid var(--border)' }}>
                {columns.map((col, colIndex) => (
                  <td key={colIndex} style={{ padding: '16px 24px', height: '64px', fontSize: '14px' }}>
                    {col.accessor ? row[col.accessor] : col.render(row)}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
