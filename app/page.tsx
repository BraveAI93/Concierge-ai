export default function EcosystemHome() {
  return (
    <div style={{ position: 'relative', minHeight: '100vh', background: '#0c0a08', overflow: 'hidden', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', fontFamily: "'Cormorant Garamond', Georgia, serif", color: '#e8dcc8', textAlign: 'center', padding: '80px 20px' }}>
      <div aria-hidden="true" style={{
        position: 'absolute', inset: 0,
        backgroundImage: `
          radial-gradient(1px 1px at 10% 20%, rgba(232,220,200,0.9), transparent),
          radial-gradient(1px 1px at 22% 65%, rgba(232,220,200,0.7), transparent),
          radial-gradient(1.5px 1.5px at 35% 15%, rgba(201,169,110,0.8), transparent),
          radial-gradient(1px 1px at 48% 80%, rgba(232,220,200,0.6), transparent),
          radial-gradient(1.5px 1.5px at 60% 35%, rgba(201,169,110,0.7), transparent),
          radial-gradient(1px 1px at 72% 55%, rgba(232,220,200,0.8), transparent),
          radial-gradient(1px 1px at 85% 20%, rgba(232,220,200,0.6), transparent),
          radial-gradient(1.5px 1.5px at 92% 70%, rgba(201,169,110,0.8), transparent),
          radial-gradient(1px 1px at 15% 88%, rgba(232,220,200,0.5), transparent),
          radial-gradient(1px 1px at 78% 90%, rgba(232,220,200,0.6), transparent),
          radial-gradient(1px 1px at 40% 45%, rgba(232,220,200,0.4), transparent),
          radial-gradient(1px 1px at 55% 8%, rgba(201,169,110,0.6), transparent)
        `,
        backgroundRepeat: 'repeat', backgroundSize: '600px 600px',
        opacity: 0.5, pointerEvents: 'none',
      }} />
      <div style={{ position: 'relative', maxWidth: 640, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <div className="landing-mark" style={{ fontSize: 28, color: '#c9a96e', marginBottom: 18 }}>✦</div>
        <h1 className="landing-title" style={{ fontSize: 'clamp(32px, 6vw, 56px)', fontWeight: 300, marginBottom: 14, lineHeight: 1.1 }}>Brave by Bruno</h1>
        <p className="landing-subtitle" style={{ fontSize: 17, fontFamily: "'Jost', sans-serif", fontWeight: 400, color: '#c9a96e', marginBottom: 22, letterSpacing: '0.02em' }}>
          AI-powered tools for independent professionals, creators, and service businesses.
        </p>
        <p className="landing-lede" style={{ fontSize: 15, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.65)', marginBottom: 44, maxWidth: 560, lineHeight: 1.8 }}>
          We are building The Concierge: a personal AI front desk and owner-side assistant that helps professionals present their work, answer client questions, capture leads, organise bookings, and grow with less friction.
        </p>
        <div className="landing-cta" style={{ display: 'flex', gap: 12, flexWrap: 'wrap', justifyContent: 'center' }}>
          <a href="/theconcierge" style={{ padding: '14px 32px', background: 'linear-gradient(135deg, #c9a96e, #7a4f0e)', borderRadius: 20, color: '#0c0a08', fontFamily: "'Jost', sans-serif", fontSize: 13, fontWeight: 700, letterSpacing: '0.08em', textTransform: 'uppercase', textDecoration: 'none' }}>
            Explore The Concierge →
          </a>
          <a href="/demo/bruno" style={{ padding: '14px 32px', background: 'none', border: '1px solid rgba(201,169,110,0.3)', borderRadius: 20, color: '#c9a96e', fontFamily: "'Jost', sans-serif", fontSize: 13, fontWeight: 500, letterSpacing: '0.06em', textDecoration: 'none' }}>
            Try Demo
          </a>
          <a href="/theconcierge/onboarding" style={{ padding: '14px 32px', background: 'none', border: '1px solid rgba(201,169,110,0.3)', borderRadius: 20, color: '#c9a96e', fontFamily: "'Jost', sans-serif", fontSize: 13, fontWeight: 500, letterSpacing: '0.06em', textDecoration: 'none' }}>
            Build My Concierge
          </a>
        </div>
      </div>
      <p className="landing-footer" style={{ position: 'relative', marginTop: 60, fontSize: 10, fontFamily: "'Jost', sans-serif", color: 'rgba(201,169,110,0.3)', letterSpacing: '0.1em', textTransform: 'uppercase' }}>
        The Concierge · Brave by Bruno · London
      </p>
    </div>
  );
}
