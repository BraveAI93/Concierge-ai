'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function TheConciergeHome() {
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const hasCookieToken = document.cookie.includes('cai_token=');
    const hasLocalToken = localStorage.getItem('ownerToken') || localStorage.getItem('cai_owner_token');
    if (hasCookieToken || hasLocalToken) {
      fetch('/api/auth/session').then(r => r.json()).then(d => {
        if (d.token) setIsLoggedIn(true);
      });
    }
  }, []);

  return (
    <div style={{ minHeight: '100vh', background: '#0c0a08', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', fontFamily: "'Cormorant Garamond', Georgia, serif", color: '#e8dcc8', textAlign: 'center', padding: '0 20px' }}>
      <div className="landing-mark" style={{ fontSize: 32, color: '#c9a96e', marginBottom: 20 }}>✦</div>
      <h1 className="landing-title" style={{ fontSize: 'clamp(32px, 6vw, 64px)', fontWeight: 300, marginBottom: 16, lineHeight: 1.1 }}>The Concierge</h1>
      <p className="landing-subtitle" style={{ fontSize: 16, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.5)', marginBottom: 48, maxWidth: 560, lineHeight: 1.7 }}>
        An AI-powered client concierge and owner assistant for independent professionals — built to answer questions, capture leads, support bookings, and help your business grow.
      </p>
      <div className="landing-demos" style={{ display: 'flex', gap: 12, flexWrap: 'wrap', justifyContent: 'center', marginBottom: 60 }}>
        {[
          { id: 'bruno', emoji: '💪', name: 'Bruno', role: 'Personal Trainer' },
          { id: 'marco', emoji: '🍽', name: 'Marco', role: 'Private Chef' },
          { id: 'nour', emoji: '📸', name: 'Nour', role: 'Photographer' },
          { id: 'sofia', emoji: '🩰', name: 'Sofia', role: 'Dance Teacher' },
          { id: 'alex', emoji: '🧠', name: 'Alex', role: 'Consultant' },
        ].map(d => (
          <button key={d.id} onClick={() => router.push(`/demo/${d.id}`)}
            style={{ padding: '14px 20px', background: 'rgba(201,169,110,0.08)', border: '1px solid rgba(201,169,110,0.2)', borderRadius: 20, cursor: 'pointer', color: '#e8dcc8', fontFamily: "'Jost', sans-serif", fontSize: 13, display: 'flex', alignItems: 'center', gap: 8 }}>
            <span>{d.emoji}</span>
            <div style={{ textAlign: 'left' }}>
              <div style={{ fontWeight: 500 }}>{d.name}</div>
              <div style={{ fontSize: 10, color: 'rgba(201,169,110,0.5)', letterSpacing: '0.06em', textTransform: 'uppercase' }}>{d.role}</div>
            </div>
          </button>
        ))}
      </div>
      <div className="landing-cta" style={{ display: 'flex', gap: 12, flexWrap: 'wrap', justifyContent: 'center' }}>
        <button onClick={() => router.push('/theconcierge/onboarding')}
          style={{ padding: '14px 32px', background: 'linear-gradient(135deg, #c9a96e, #7a4f0e)', border: 'none', borderRadius: 20, cursor: 'pointer', color: '#0c0a08', fontFamily: "'Jost', sans-serif", fontSize: 13, fontWeight: 700, letterSpacing: '0.08em', textTransform: 'uppercase' }}>
          Build My Concierge ✦
        </button>
        {isLoggedIn ? (
          <button onClick={() => router.push('/theconcierge/dashboard')}
            style={{ padding: '14px 32px', background: 'none', border: '1px solid rgba(201,169,110,0.3)', borderRadius: 20, cursor: 'pointer', color: '#c9a96e', fontFamily: "'Jost', sans-serif", fontSize: 13, fontWeight: 500, letterSpacing: '0.06em' }}>
            Go to Dashboard
          </button>
        ) : (
          <button onClick={() => router.push('/theconcierge/owner-auth')}
            style={{ padding: '14px 32px', background: 'none', border: '1px solid rgba(201,169,110,0.3)', borderRadius: 20, cursor: 'pointer', color: '#c9a96e', fontFamily: "'Jost', sans-serif", fontSize: 13, fontWeight: 500, letterSpacing: '0.06em' }}>
            Owner Login
          </button>
        )}
      </div>
      <p className="landing-footer" style={{ marginTop: 60, fontSize: 10, fontFamily: "'Jost', sans-serif", color: 'rgba(201,169,110,0.2)', letterSpacing: '0.1em', textTransform: 'uppercase' }}>
        Brave by Bruno · London
      </p>
    </div>
  );
}
