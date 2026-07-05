export default function EcosystemHome() {
  return (
    <div style={{ minHeight: '100vh', background: '#0c0a08', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', fontFamily: "'Cormorant Garamond', Georgia, serif", color: '#e8dcc8', textAlign: 'center', padding: '0 20px' }}>
      <h1 style={{ fontSize: 'clamp(28px, 5vw, 48px)', fontWeight: 300, marginBottom: 16 }}>Brave by Bruno</h1>
      <p style={{ fontSize: 14, fontFamily: "'Jost', sans-serif", fontWeight: 300, color: 'rgba(232,220,200,0.5)', marginBottom: 32 }}>
        More to come.
      </p>
      <a href="/theconcierge" style={{ fontFamily: "'Jost', sans-serif", fontSize: 13, color: '#c9a96e', textDecoration: 'none', letterSpacing: '0.06em' }}>
        The Concierge →
      </a>
    </div>
  );
}
