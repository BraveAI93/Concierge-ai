'use client';
import { useState, useEffect } from 'react';
import { BACKEND_URL, LEGAL_FORM_FIELDS, FORM_SLUGS } from '@/lib/constants';

export default function FormPage({ profileSlug, formTypeSlug }) {
  const formKey = FORM_SLUGS[formTypeSlug];
  const [profile, setProfile] = useState(null);
  const [notFound, setNotFound] = useState(false);
  const [clientName, setClientName] = useState('');
  const [clientEmail, setClientEmail] = useState('');
  const [answer, setAnswer] = useState('');
  const [agreed, setAgreed] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const [err, setErr] = useState('');

  useEffect(() => {
    if (!profileSlug || !formKey) { setNotFound(true); return; }
    fetch(`${BACKEND_URL}/forms/${profileSlug}/${formTypeSlug}`)
      .then(r => r.ok ? r.json() : null)
      .then(d => { if (d && d.name) setProfile(d); else setNotFound(true); })
      .catch(() => setNotFound(true));
  }, []);

  const submit = async () => {
    if (!agreed || !clientName.trim()) return;
    setSubmitting(true); setErr('');
    try {
      const responses = JSON.stringify({ answer, agreed: true });
      const r = await fetch(`${BACKEND_URL}/forms/${profileSlug}/${formTypeSlug}`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ client_name: clientName, client_email: clientEmail, responses })
      });
      if (r.ok) { setDone(true); }
      else { const d = await r.json().catch(() => ({})); setErr(d.error || 'Submission failed — please try again.'); }
    } catch(e) { setErr('Network error — please try again.'); }
    setSubmitting(false);
  };

  const acc = (profile?.accent) || '#c9a96e';
  const def = formKey ? LEGAL_FORM_FIELDS[formKey] : null;
  const txt = def?.en;

  if (notFound) return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',alignItems:'center',justifyContent:'center',flexDirection:'column',gap:12}}>
      <div style={{fontSize:32,color:'rgba(201,169,110,0.3)'}}>✦</div>
      <div style={{fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)',fontSize:13}}>Form not found</div>
    </div>
  );

  if (!profile) return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',alignItems:'center',justifyContent:'center'}}>
      <div style={{fontSize:28,color:acc,animation:'bravePulse 1.8s ease-in-out infinite',textShadow:'0 0 20px rgba(201,169,110,0.6)'}}>✦</div>
    </div>
  );

  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',padding:'0 0 60px',fontFamily:"'Cormorant Garamond',Georgia,serif",color:'#e8dcc8'}}>
      <div style={{background:'rgba(12,10,8,0.95)',borderBottom:'1px solid rgba(201,169,110,0.15)',padding:'22px 20px 18px',textAlign:'center',backdropFilter:'blur(20px)'}}>
        <div style={{fontSize:22,color:acc,marginBottom:4,letterSpacing:'0.06em'}}>✦</div>
        <div style={{fontSize:22,fontWeight:400,color:'#e8dcc8',marginBottom:3}}>{profile.name}</div>
        <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.5)',letterSpacing:'0.1em',textTransform:'uppercase'}}>{profile.profession}</div>
      </div>
      <div style={{maxWidth:520,margin:'0 auto',padding:'28px 18px'}}>
        {done ? (
          <div style={{textAlign:'center',padding:'40px 0'}}>
            <div style={{fontSize:42,marginBottom:16,color:acc}}>✓</div>
            <div style={{fontSize:22,fontWeight:400,color:'#e8dcc8',marginBottom:8,fontFamily:"'Cormorant Garamond',serif"}}>Form submitted</div>
            <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.45)',lineHeight:1.7}}>Thank you, {clientName}.<br/>Your response has been securely saved.</div>
          </div>
        ) : (
          <>
            <div style={{marginBottom:24}}>
              <div style={{fontSize:24,fontWeight:400,color:'#e8dcc8',marginBottom:6,lineHeight:1.3}}>{txt?.title}</div>
              <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.38)',lineHeight:1.6}}>Complete this form before your session with {profile.name}. Fully GDPR compliant and securely stored.</div>
            </div>
            <div style={{marginBottom:14}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:5}}>Full name <span style={{color:'rgba(200,80,80,0.7)'}}>*</span></div>
              <input value={clientName} onChange={e=>setClientName(e.target.value)} placeholder="Your full name" style={{width:'100%',padding:'11px 14px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:10,color:'#e8dcc8',fontSize:14,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
            </div>
            <div style={{marginBottom:20}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:5}}>Email (optional)</div>
              <input value={clientEmail} onChange={e=>setClientEmail(e.target.value)} placeholder="you@email.com" type="email" style={{width:'100%',padding:'11px 14px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:10,color:'#e8dcc8',fontSize:14,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
            </div>
            <div style={{marginBottom:20,padding:'16px 18px',background:'rgba(201,169,110,0.04)',border:'1px solid rgba(201,169,110,0.12)',borderRadius:14}}>
              <div style={{fontSize:14,fontWeight:500,color:'#e8dcc8',marginBottom:6,fontFamily:"'Jost',sans-serif"}}>{txt?.title}</div>
              <div style={{fontSize:13,fontFamily:"'Cormorant Garamond',serif",color:'rgba(232,220,200,0.65)',lineHeight:1.7,marginBottom:12}}>{txt?.q}</div>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.08em',textTransform:'uppercase',color:'rgba(201,169,110,0.4)',marginBottom:5}}>Your response (optional)</div>
              <textarea value={answer} onChange={e=>setAnswer(e.target.value)} placeholder="Leave blank if not applicable" rows={3} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",marginBottom:12,outline:'none'}}/>
              <label style={{display:'flex',alignItems:'flex-start',gap:10,cursor:'pointer'}}>
                <input type="checkbox" checked={agreed} onChange={e=>setAgreed(e.target.checked)} style={{marginTop:3,accentColor:acc}}/>
                <span style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.7)',lineHeight:1.5}}>I have read and understood the above, and I agree</span>
              </label>
            </div>
            {err && <div style={{marginBottom:14,padding:'10px 14px',background:'rgba(200,80,80,0.08)',border:'1px solid rgba(200,80,80,0.2)',borderRadius:8,fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(220,140,140,0.9)'}}>{err}</div>}
            <button onClick={submit} disabled={!agreed||!clientName.trim()||submitting} style={{width:'100%',padding:'14px 0',background:(agreed&&clientName.trim())?`linear-gradient(135deg,${acc},#7a4f0e)`:'rgba(201,169,110,0.1)',border:'none',borderRadius:20,cursor:(agreed&&clientName.trim())?'pointer':'not-allowed',color:(agreed&&clientName.trim())?'#0c0a08':'rgba(232,220,200,0.25)',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.08em',textTransform:'uppercase'}}>
              {submitting ? 'Submitting…' : '✓ Submit & Sign'}
            </button>
            <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)',textAlign:'center',marginTop:10}}>
              Digitally signed by {clientName||'…'} · {new Date().toLocaleDateString()} · Secured by Concierge AI
            </div>
          </>
        )}
      </div>
    </div>
  );
}
