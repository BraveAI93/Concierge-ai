'use client';
import { useEffect, useState } from 'react';
import Chat from '@/components/Chat';

function getVideoEmbed(url) {
  if (!url) return null;
  const ytShort = url.match(/youtube\.com\/shorts\/([a-zA-Z0-9_-]+)/);
  if (ytShort) return { src: `https://www.youtube.com/embed/${ytShort[1]}?autoplay=1&mute=1&loop=1&playlist=${ytShort[1]}&playsinline=1`, aspect: '9/16', maxWidth: 280 };
  const yt = url.match(/(?:youtube\.com\/(?:watch\?v=|embed\/)|youtu\.be\/)([a-zA-Z0-9_-]+)/);
  if (yt) return { src: `https://www.youtube.com/embed/${yt[1]}?autoplay=1&mute=1&loop=1&playlist=${yt[1]}&playsinline=1`, aspect: '16/9' };
  const vm = url.match(/vimeo\.com\/(?:video\/)?(\d+)/);
  if (vm) return { src: `https://player.vimeo.com/video/${vm[1]}?autoplay=1&muted=1&loop=1&background=1`, aspect: '16/9' };
  return null;
}

const linkBtnStyle = (primary, accent) => ({
  display: 'flex', alignItems: 'center', gap: 12, padding: '13px 16px',
  background: primary ? `linear-gradient(135deg,${accent},#7a4f0e)` : 'rgba(255,255,255,0.025)',
  border: primary ? 'none' : '1px solid rgba(201,169,110,0.12)',
  borderRadius: 20, textDecoration: 'none', color: primary ? '#0c0a08' : '#e8dcc8',
  marginBottom: 9, fontFamily: "'Jost',sans-serif", cursor: 'pointer', width: '100%',
});
const linkIconStyle = (primary) => ({ width: 36, height: 36, borderRadius: '50%', background: primary ? 'rgba(0,0,0,0.15)' : 'rgba(201,169,110,0.1)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 16, flexShrink: 0 });

export default function PublicProfileClient({ data, systemPrompt, slug }) {
  const [bubbleState, setBubbleState] = useState('closed'); // closed | teaser | chat
  const [hintVisible, setHintVisible] = useState(false);
  const [teaserRevealed, setTeaserRevealed] = useState(false);
  const [copied, setCopied] = useState(false);
  const [publicURL, setPublicURL] = useState(`https://bravebybruno.com/${slug}`);

  useEffect(() => {
    setPublicURL(window.location.href);
    const t1 = setTimeout(() => setHintVisible(true), 2800);
    const t2 = setTimeout(() => setHintVisible(false), 7800);
    return () => { clearTimeout(t1); clearTimeout(t2); };
  }, []);

  useEffect(() => {
    if (bubbleState !== 'teaser') { setTeaserRevealed(false); return; }
    const t = setTimeout(() => setTeaserRevealed(true), 1400);
    return () => clearTimeout(t);
  }, [bubbleState]);

  const accent = data.accent || '#c9a96e';
  const name = data.name || 'Your Name';
  const profession = data.profession || 'Professional';
  const loc = data.loc || '';
  const tag = data.tag || '';
  const services = (data.services || []).filter(s => s.name);
  const videoEmbed = getVideoEmbed(data.video);
  const instagram = data.ig || data.instagram || '';
  const whatsappDigits = (data.wa || '').replace(/[^0-9]/g, '').replace(/^0+/, '');

  const chatProfile = {
    ...data,
    location: data.loc,
    tagline: data.tag || data.tagline,
    gradient: data.gradient || `linear-gradient(135deg,${accent},#5e3a0e)`,
    emoji: data.emoji || '✦',
  };

  const openChat = () => { setBubbleState('chat'); setHintVisible(false); };
  const openTeaser = () => { setBubbleState('teaser'); setHintVisible(false); };

  const copyLink = () => {
    navigator.clipboard?.writeText(publicURL).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div style={{ minHeight: '100vh', background: '#0c0a08', color: '#e8dcc8', fontFamily: "'Cormorant Garamond',serif" }}>
      <div style={{ maxWidth: 480, margin: '0 auto', padding: '32px 18px 100px' }}>

        <div style={{ width: 88, height: 88, borderRadius: '50%', background: `linear-gradient(135deg,${accent},#5e3a0e)`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 32, margin: '0 auto 14px', border: '2px solid rgba(201,169,110,0.3)' }}>✦</div>
        <div style={{ fontSize: 26, fontWeight: 400, textAlign: 'center', color: '#e8dcc8', marginBottom: 4 }}>{name}</div>
        <div style={{ fontSize: 11, fontFamily: "'Jost',sans-serif", fontWeight: 300, letterSpacing: '0.18em', textTransform: 'uppercase', color: 'rgba(201,169,110,0.55)', textAlign: 'center', marginBottom: 4 }}>{profession}</div>
        {loc && <div style={{ fontSize: 11, fontFamily: "'Jost',sans-serif", color: 'rgba(232,220,200,0.35)', textAlign: 'center', marginBottom: 10 }}>📍 {loc}</div>}
        {tag && <div style={{ fontSize: 14, fontStyle: 'italic', color: 'rgba(232,220,200,0.45)', textAlign: 'center', marginBottom: 22, lineHeight: 1.6 }}>"{tag}"</div>}

        {videoEmbed && (
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'center' }}>
            <iframe src={videoEmbed.src} title="video" frameBorder="0" allow="autoplay; fullscreen; picture-in-picture" allowFullScreen
              style={{ width: '100%', maxWidth: videoEmbed.maxWidth || undefined, aspectRatio: videoEmbed.aspect, borderRadius: 12, display: 'block' }} />
          </div>
        )}

        {services.length > 0 && (
          <>
            <div style={{ fontSize: 9, fontFamily: "'Jost',sans-serif", fontWeight: 500, letterSpacing: '0.14em', textTransform: 'uppercase', color: 'rgba(201,169,110,0.38)', marginBottom: 10, marginTop: 24 }}>Services</div>
            {services.map((s, i) => {
              const durTxt = s.durNum ? `${s.durNum} ${s.durUnit === 'h' ? 'hour' + (s.durNum == 1 ? '' : 's') : s.durUnit === 'days' ? 'day' + (s.durNum == 1 ? '' : 's') : 'min'}` : '';
              const priceTxt = s.priceNum ? `${s.currency || '£'}${s.priceNum}` : '';
              return (
                <div key={i} style={{ background: 'rgba(255,255,255,0.025)', border: '1px solid rgba(201,169,110,0.1)', borderRadius: 20, padding: '14px 16px', marginBottom: 10 }}>
                  <div style={{ fontSize: 16, fontWeight: 500, color: '#e8dcc8', marginBottom: 4 }}>{s.name}</div>
                  {durTxt && <div style={{ fontSize: 11, fontFamily: "'Jost',sans-serif", color: 'rgba(232,220,200,0.4)', marginBottom: 2 }}>⏱ {durTxt}</div>}
                  {priceTxt && <div style={{ fontSize: 14, fontFamily: "'Jost',sans-serif", fontWeight: 500, color: accent, marginBottom: 4 }}>{priceTxt}</div>}
                  {s.desc && <div style={{ fontSize: 13, color: 'rgba(232,220,200,0.5)', lineHeight: 1.6, marginTop: 4 }}>{s.desc}</div>}
                </div>
              );
            })}
          </>
        )}

        <div style={{ fontSize: 9, fontFamily: "'Jost',sans-serif", fontWeight: 500, letterSpacing: '0.14em', textTransform: 'uppercase', color: 'rgba(201,169,110,0.38)', marginBottom: 10, marginTop: 24 }}>Connect</div>

        {data.booking && (
          <a href={data.booking} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(true, accent)}>
            <div style={linkIconStyle(true)}>📅</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Book a Session</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>Check availability & book online</div></div>
          </a>
        )}

        <div onClick={openTeaser} style={{ ...linkBtnStyle(false, accent), cursor: 'pointer' }}>
          <div style={linkIconStyle(false)}>💬</div>
          <div><div style={{ fontSize: 13, fontWeight: 500 }}>Ask the Concierge</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>Services, pricing & booking · Available 24/7</div></div>
        </div>

        {data.calendar && (
          <a href={data.calendar} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>🗓️</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Check my availability</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>View calendar</div></div>
          </a>
        )}
        {whatsappDigits.length >= 7 && (
          <a href={`https://wa.me/${whatsappDigits}`} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>💚</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>WhatsApp me</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>{data.wa}</div></div>
          </a>
        )}
        {data.phone && (
          <a href={`tel:${data.phone.replace(/[^0-9+]/g, '')}`} style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>📞</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Call me</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>{data.phone}</div></div>
          </a>
        )}
        {data.email && (
          <a href={`mailto:${data.email}`} style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>✉️</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Email me</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>{data.email}</div></div>
          </a>
        )}
        {instagram && (
          <a href={`https://instagram.com/${instagram.replace('@', '')}`} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>📸</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>@{instagram.replace('@', '')}</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>Follow on Instagram</div></div>
          </a>
        )}
        {data.tg && (
          <a href={`https://t.me/${data.tg.replace('@', '')}`} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>✈️</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Telegram</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>{data.tg}</div></div>
          </a>
        )}
        {data.gallery && (
          <a href={data.gallery} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>🖼️</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>View my gallery</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>Photos & portfolio</div></div>
          </a>
        )}
        {data.video && !videoEmbed && (
          <a href={data.video} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>🎬</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Watch my showreel</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>Video portfolio</div></div>
          </a>
        )}
        {data.lat && data.lng && (
          <a href={`https://maps.google.com/?q=${data.lat},${data.lng}`} target="_blank" rel="noopener noreferrer" style={linkBtnStyle(false, accent)}>
            <div style={linkIconStyle(false)}>📍</div>
            <div><div style={{ fontSize: 13, fontWeight: 500 }}>Find me on the map</div><div style={{ fontSize: 10, fontWeight: 300, opacity: 0.65, marginTop: 1 }}>{loc}</div></div>
          </a>
        )}

        <div style={{ fontSize: 9, fontFamily: "'Jost',sans-serif", fontWeight: 500, letterSpacing: '0.14em', textTransform: 'uppercase', color: 'rgba(201,169,110,0.38)', marginBottom: 10, marginTop: 24 }}>Share</div>
        <div style={{ background: 'rgba(255,255,255,0.025)', border: '1px solid rgba(201,169,110,0.1)', borderRadius: 20, padding: '14px 16px', marginBottom: 10 }}>
          <div style={{ fontSize: 12, fontFamily: "'Jost',sans-serif", color: accent, wordBreak: 'break-all', marginBottom: 10 }}>{publicURL}</div>
          <button onClick={copyLink} style={{ padding: '7px 14px', background: 'rgba(201,169,110,0.08)', border: '1px solid rgba(201,169,110,0.2)', borderRadius: 20, cursor: 'pointer', fontSize: 11, fontFamily: "'Jost',sans-serif", color: accent }}>
            {copied ? '✓ Copied' : '📋 Copy link'}
          </button>
        </div>

        <div style={{ textAlign: 'center', marginTop: 32, fontSize: 9, fontFamily: "'Jost',sans-serif", color: 'rgba(201,169,110,0.2)', letterSpacing: '0.08em' }}>Powered by The Concierge · Brave by Bruno</div>
      </div>

      {bubbleState === 'closed' && (
        <>
          {hintVisible && (
            <div style={{ position: 'fixed', bottom: 90, right: 18, zIndex: 9998, maxWidth: 190, background: 'rgba(14,11,7,0.97)', border: '1px solid rgba(201,169,110,0.28)', borderRadius: '14px 14px 4px 14px', padding: '10px 13px', fontSize: 12, fontFamily: "'Jost',sans-serif", color: '#e8dcc8', lineHeight: 1.55, boxShadow: '0 8px 32px rgba(0,0,0,0.5)' }}>
              Ask the Concierge about services, pricing & booking ✦
            </div>
          )}
          <button onClick={openTeaser} style={{ position: 'fixed', bottom: 22, right: 18, zIndex: 9999, width: 56, height: 56, borderRadius: '50%', border: 'none', cursor: 'pointer', background: `linear-gradient(135deg,${accent},#7a4f0e)`, boxShadow: `0 4px 20px ${accent}66`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 22 }}>✦</button>
        </>
      )}

      {bubbleState === 'teaser' && (
        <>
          <div onClick={() => setBubbleState('closed')} style={{ position: 'fixed', inset: 0, zIndex: 9996, background: 'rgba(0,0,0,0.55)' }} />
          <div style={{ position: 'fixed', bottom: 0, left: 0, right: 0, zIndex: 9997, padding: '0 12px 88px', display: 'flex', flexDirection: 'column', justifyContent: 'flex-end' }}>
            <div style={{ width: '100%', maxWidth: 400, margin: '0 auto', background: 'rgba(14,11,7,0.98)', border: '1px solid rgba(201,169,110,0.22)', borderRadius: 20, padding: 18, boxShadow: '0 -8px 40px rgba(0,0,0,0.6)' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 14 }}>
                <div style={{ width: 38, height: 38, borderRadius: '50%', background: `linear-gradient(135deg,${accent},#7a4f0e)`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 16, flexShrink: 0 }}>✦</div>
                <div>
                  <div style={{ fontSize: 13, fontFamily: "'Jost',sans-serif", fontWeight: 500, color: '#e8dcc8' }}>{name}'s Concierge</div>
                  <div style={{ fontSize: 9, fontFamily: "'Jost',sans-serif", color: 'rgba(120,180,100,0.8)', display: 'flex', alignItems: 'center', gap: 4, marginTop: 1 }}>
                    <div style={{ width: 5, height: 5, borderRadius: '50%', background: '#7aba6a' }} />Available now · 24/7
                  </div>
                </div>
                <div onClick={() => setBubbleState('closed')} style={{ marginLeft: 'auto', background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.08)', borderRadius: '50%', width: 26, height: 26, display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', fontSize: 13, color: 'rgba(232,220,200,0.4)' }}>×</div>
              </div>
              <div style={{ background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(201,169,110,0.1)', borderRadius: '12px 12px 12px 4px', padding: '11px 14px', marginBottom: 12, fontSize: 14, fontFamily: "'Cormorant Garamond',serif", lineHeight: 1.65, color: 'rgba(232,220,200,0.9)' }}>
                {!teaserRevealed
                  ? <div style={{ display: 'flex', gap: 5, alignItems: 'center', padding: '4px 2px' }}>{[0, 1, 2].map(i => <span key={i} style={{ width: 6, height: 6, borderRadius: '50%', background: accent, opacity: 0.5 }} />)}</div>
                  : <div>Hi 👋 I'm {name}'s digital concierge — I can help with services, pricing, booking and availability. What would you like to know?</div>}
              </div>
              {teaserRevealed && (
                <>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6, marginBottom: 12 }}>
                    {[['💆', 'Services & prices'], ['📅', 'How to book'], ['📍', 'Location'], ['✨', `About ${name}`]].map(([icon, label]) => (
                      <div key={label} onClick={openChat} style={{ padding: '6px 12px', borderRadius: 20, cursor: 'pointer', fontSize: 11, fontFamily: "'Jost',sans-serif", background: 'rgba(201,169,110,0.07)', border: '1px solid rgba(201,169,110,0.2)', color: 'rgba(232,220,200,0.75)' }}>{icon} {label}</div>
                    ))}
                  </div>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 7 }}>
                    <div onClick={openChat} style={{ padding: '10px 8px', borderRadius: 11, cursor: 'pointer', textAlign: 'center', fontSize: 11, fontFamily: "'Jost',sans-serif", fontWeight: 600, background: `linear-gradient(135deg,${accent},#7a4f0e)`, color: '#0c0a08', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4 }}>
                      <div style={{ fontSize: 16 }}>💬</div><div style={{ fontSize: 10 }}>Ask the Concierge</div>
                    </div>
                    {data.booking
                      ? <a href={data.booking} target="_blank" rel="noopener noreferrer" style={{ padding: '10px 8px', borderRadius: 11, textAlign: 'center', fontSize: 11, fontFamily: "'Jost',sans-serif", fontWeight: 500, border: '1px solid rgba(201,169,110,0.18)', background: 'rgba(201,169,110,0.06)', color: 'rgba(232,220,200,0.7)', textDecoration: 'none', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4 }}>
                          <div style={{ fontSize: 16 }}>📅</div><div style={{ fontSize: 10 }}>Book a session</div>
                        </a>
                      : <div onClick={openChat} style={{ padding: '10px 8px', borderRadius: 11, cursor: 'pointer', textAlign: 'center', fontSize: 11, fontFamily: "'Jost',sans-serif", fontWeight: 500, border: '1px solid rgba(201,169,110,0.18)', background: 'rgba(201,169,110,0.06)', color: 'rgba(232,220,200,0.7)', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4 }}>
                          <div style={{ fontSize: 16 }}>📅</div><div style={{ fontSize: 10 }}>Book a session</div>
                        </div>}
                  </div>
                </>
              )}
            </div>
          </div>
        </>
      )}

      {bubbleState === 'chat' && (
        <div style={{ position: 'fixed', inset: 0, zIndex: 10000, background: '#0c0a08' }}>
          <Chat profile={chatProfile} systemPrompt={systemPrompt} profileId={slug} onBack={() => setBubbleState('closed')} lang="en" />
        </div>
      )}
    </div>
  );
}
