'use client';
import { useEffect } from 'react';
import { generateBusinessPage } from '@/lib/generateBusinessPage';

export default function PublicProfileClient({ data, systemPrompt, slug }) {
  useEffect(() => {
    const html = generateBusinessPage(data, 'en');
    document.open();
    document.write(html);
    document.close();
  }, []);

  return (
    <div style={{ minHeight: '100vh', background: '#0c0a08', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div style={{ fontSize: 28, color: '#c9a96e', animation: 'bravePulse 1.8s ease-in-out infinite' }}>✦</div>
    </div>
  );
}
