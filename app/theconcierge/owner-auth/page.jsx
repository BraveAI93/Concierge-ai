'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function OwnerAuthPage() {
  const router = useRouter();
  const [mode, setMode] = useState('login');
  const [slug, setSlug] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState('');

  useEffect(() => {
    fetch('/api/auth/session').then(r => r.json()).then(d => {
      if (d.token) router.replace('/theconcierge/dashboard');
    });
  }, []);

  const submit = async () => {
    setErr(''); setLoading(true);
    try {
      const endpoint = mode === 'signup' ? '/api/auth/signup' : '/api/auth/login';
      const body = mode === 'signup' ? { slug, email, password } : { email, password };
      const r = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      const d = await r.json();
      if (!r.ok) { setErr(d.error || 'Something went wrong'); setLoading(false); return; }
      localStorage.setItem('cai_owner_token', d.token);
      localStorage.setItem('cai_owner_slug', d.slug || '');
      localStorage.setItem('ownerToken', d.token);
      localStorage.setItem('ownerProfile', JSON.stringify({ slug: d.slug || '' }));
      router.push('/theconcierge/dashboard');
    } catch(e) {
      setErr('Connection error');
    }
    setLoading(false);
  };

  const acc = '#c9a96e';
  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',alignItems:'center',justifyContent:'center',padding:20,fontFamily:"'Cormorant Garamond',Georgia,serif"}}>
      <div style={{width:'100%',maxWidth:380,background:'rgba(255,255,255,0.03)',border:'1px solid rgba(201,169,110,0.15)',borderRadius:20,padding:'32px 28px'}}>
        <div style={{textAlign:'center',marginBottom:24}}>
          <div style={{fontSize:22,color:acc,marginBottom:6}}>🔐</div>
          <div style={{fontFamily:"'Jost',sans-serif",fontWeight:500,fontSize:14,letterSpacing:'0.08em',textTransform:'uppercase',color:'#e8dcc8',marginBottom:4}}>
            {mode === 'signup' ? 'Create Your Account' : 'Welcome back'}
          </div>
          <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)'}}>
            {mode === 'signup' ? 'Link an email & password to your profile' : 'Login to manage your concierge'}
          </div>
        </div>
        {mode === 'signup' && (
          <div style={{marginBottom:12}}>
            <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:6}}>Your profile link (slug)</div>
            <input value={slug} onChange={e=>setSlug(e.target.value)} placeholder="e.g. bruno-massage-london" style={{width:'100%',padding:'10px 13px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:20,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
          </div>
        )}
        <div style={{marginBottom:12}}>
          <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:6}}>Email</div>
          <input type="email" value={email} onChange={e=>setEmail(e.target.value)} placeholder="you@email.com" style={{width:'100%',padding:'10px 13px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:20,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
        </div>
        <div style={{marginBottom:14}}>
          <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:6}}>Password</div>
          <input type="password" value={password} onChange={e=>setPassword(e.target.value)} onKeyDown={e=>e.key==='Enter'&&submit()} placeholder="6+ characters" style={{width:'100%',padding:'10px 13px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:20,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif",outline:'none'}}/>
        </div>
        {err && <div style={{fontSize:12,color:'#c06060',fontFamily:"'Jost',sans-serif",marginBottom:10}}>{err}</div>}
        <button onClick={submit} disabled={loading} style={{width:'100%',padding:'11px 0',background:`linear-gradient(135deg,${acc},#8c5e14)`,border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.07em',textTransform:'uppercase',marginBottom:14}}>
          {loading ? '...' : (mode === 'signup' ? 'Create Account →' : 'Login →')}
        </button>
        <div style={{textAlign:'center',fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)'}}>
          {mode === 'signup'
            ? <span>Already have an account? <span onClick={()=>{setMode('login');setErr('');}} style={{color:acc,cursor:'pointer',textDecoration:'underline'}}>Login</span></span>
            : <span>New here? <span onClick={()=>{setMode('signup');setErr('');}} style={{color:acc,cursor:'pointer',textDecoration:'underline'}}>Create an account</span></span>
          }
        </div>
        <div style={{textAlign:'center',marginTop:14}}>
          <button onClick={()=>router.push('/')} style={{background:'none',border:'none',cursor:'pointer',fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)',textDecoration:'underline'}}>← Back to home</button>
        </div>
      </div>
    </div>
  );
}
