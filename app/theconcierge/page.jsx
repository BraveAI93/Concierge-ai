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

  const liveToday = [
    'A 24/7 AI concierge answering client questions on your public business page',
    'Lead capture with interest temperature \u2014 hot, warm or cold',
    'Owner dashboard for conversations, leads and profile editing',
    'Guided onboarding with GDPR-compliant consent and legal forms',
    'Contact buttons, media gallery, QR code and calendar links',
  ];

  const inDevelopment = [
    'Calendar-integrated booking confirmation',
    'Installable app with push notifications (PWA)',
    'Deposits and payments via Stripe',
  ];

  return (
    <div style={{ minHeight: '100vh', background: '#0c0a08', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'flex-start', fontFamily: "'Cormorant Garamond', Georgia, serif", color: '#e8dcc8', textAlign: 'center', padding: '80px 20px' }}>
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

      <section style={{ marginTop: 72, width: '100%', maxWidth: 560 }}>
        <h2 style={{ fontSize: 'clamp(22px, 4vw, 32px)', fontWeight: 300, marginBottom: 24, lineHeight: 1.2 }}>Live today</h2>
        <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'flex', flexDirection: 'column', gap: 14 }}>
          {liveToday.map(item => (
            <li key={item} style={{ display: 'flex', gap: 12, alignItems: 'flex-start', textAlign: 'left', fontSize: 15, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.75)', lineHeight: 1.7 }}>
              <span aria-hidden="true" style={{ color: '#c9a96e', fontSize: 13, lineHeight: '1.9' }}>✦</span>
              <span>{item}</span>
            </li>
          ))}
        </ul>
      </section>

      <section style={{ marginTop: 56, width: '100%', maxWidth: 560 }}>
        <h2 style={{ fontSize: 'clamp(22px, 4vw, 32px)', fontWeight: 300, marginBottom: 24, lineHeight: 1.2 }}>In development</h2>
        <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'flex', flexDirection: 'column', gap: 14 }}>
          {inDevelopment.map(item => (
            <li key={item} style={{ display: 'flex', gap: 12, alignItems: 'flex-start', textAlign: 'left', fontSize: 15, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.55)', lineHeight: 1.7 }}>
              <span aria-hidden="true" style={{ color: 'rgba(201,169,110,0.5)', fontSize: 13, lineHeight: '1.9' }}>✧</span>
              <span>{item}</span>
            </li>
          ))}
        </ul>
      </section>

      <section style={{ marginTop: 56, width: '100%', maxWidth: 560 }}>
        <p style={{ fontSize: 15, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.65)', lineHeight: 1.8 }}>
          The Concierge is subscription software by Brave by Bruno, a London-based company founded by Bruno Aversa.
        </p>
        <p style={{ marginTop: 16, fontSize: 15, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.65)', lineHeight: 1.8 }}>
          Contact:{' '}
          <a href="mailto:bruno@bravebybruno.com" style={{ color: '#c9a96e', textDecoration: 'none', borderBottom: '1px solid rgba(201,169,110,0.35)' }}>
            bruno@bravebybruno.com
          </a>
        </p>
      </section>

      <footer style={{ marginTop: 64, paddingTop: 28, borderTop: '1px solid rgba(201,169,110,0.15)', width: '100%', maxWidth: 640, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 14 }}>
        <div style={{ display: 'flex', gap: 22, flexWrap: 'wrap', justifyContent: 'center', fontFamily: "'Jost', sans-serif", fontSize: 12, letterSpacing: '0.04em' }}>
          <a href="mailto:bruno@bravebybruno.com" style={{ color: 'rgba(201,169,110,0.75)', textDecoration: 'none' }}>bruno@bravebybruno.com</a>
          <a href="/privacy" style={{ color: 'rgba(201,169,110,0.75)', textDecoration: 'none' }}>Privacy &amp; Terms</a>
          <a href="/" style={{ color: 'rgba(201,169,110,0.75)', textDecoration: 'none' }}>Brave by Bruno</a>
        </div>
        <p className="landing-footer" style={{ fontSize: 10, fontFamily: "'Jost', sans-serif", color: 'rgba(201,169,110,0.2)', letterSpacing: '0.1em', textTransform: 'uppercase', margin: 0 }}>
          © 2026 Brave by Bruno · London, United Kingdom
        </p>
      </footer>
    </div>
  );
}