'use client';
import { useState } from 'react';

const BACKEND_URL = 'https://concierge-backend-80rb.onrender.com';

export default function AdminDashboard({ onBack }) {
  const [loading, setLoading] = useState(false);
  const [key, setKey] = useState('');
  const [authed, setAuthed] = useState(false);
  const [err, setErr] = useState('');
  const [realData, setRealData] = useState(null);

  const login = async () => {
    setAuthed(true); setLoading(true); setErr('');
    try {
      const r = await fetch(`${BACKEND_URL}/admin/stats`, { headers: { 'X-Admin-Key': key } });
      if (r.ok) { const d = await r.json(); setRealData(d); }
      else { setAuthed(false); setErr('Invalid admin key'); }
    } catch(e) { setAuthed(false); setErr('Connection error'); }
    setLoading(false);
  };

  if (!authed) return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',alignItems:'center',justifyContent:'center',padding:20}}>
      <div style={{width:'100%',maxWidth:380,background:'rgba(255,255,255,0.03)',border:'1px solid rgba(201,169,110,0.15)',borderRadius:20,padding:'32px 28px'}}>
        <div style={{textAlign:'center',marginBottom:24}}>
          <div style={{fontSize:22,color:'#c9a96e',marginBottom:6}}>⚙</div>
          <div style={{fontFamily:"'Jost',sans-serif",fontWeight:500,fontSize:14,letterSpacing:'0.08em',textTransform:'uppercase',color:'#e8dcc8',marginBottom:4}}>Admin Dashboard</div>
          <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)'}}>Platform-wide stats — Bruno only</div>
        </div>
        <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:7}}>Admin Key</div>
        <input type="password" value={key} onChange={e=>setKey(e.target.value)} onKeyDown={e=>e.key==='Enter'&&login()} placeholder="Enter admin key" style={{width:'100%',padding:'10px 13px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:20,color:'#e8dcc8',fontSize:14,fontFamily:"'Jost',sans-serif",marginBottom:10,outline:'none'}}/>
        {err && <div style={{fontSize:12,color:'#c06060',fontFamily:"'Jost',sans-serif",marginBottom:10}}>{err}</div>}
        <button onClick={login} disabled={loading} style={{width:'100%',padding:'11px 0',background:'linear-gradient(135deg,#c9a96e,#8c5e14)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.07em',textTransform:'uppercase',marginBottom:14}}>{loading?'Checking...':'Access Dashboard →'}</button>
        <div style={{textAlign:'center'}}><button onClick={onBack} style={{background:'none',border:'none',cursor:'pointer',fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)',textDecoration:'underline'}}>← Back</button></div>
      </div>
    </div>
  );

  const s = realData || { total_users:0, total_conversations:0, total_messages:0, total_leads:0, estimated_cost_gbp:'0.00' };
  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',padding:'24px 16px',fontFamily:"'Cormorant Garamond',Georgia,serif",color:'#e8dcc8'}}>
      <div style={{maxWidth:700,margin:'0 auto'}}>
        <div style={{display:'flex',alignItems:'center',justifyContent:'space-between',marginBottom:28}}>
          <div>
            <div style={{fontFamily:"'Jost',sans-serif",fontWeight:500,fontSize:14,letterSpacing:'0.08em',textTransform:'uppercase',color:'#e8dcc8'}}>Admin Dashboard</div>
            <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',marginTop:2}}>Platform-wide · Bruno only</div>
          </div>
          <button onClick={onBack} style={{background:'none',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,padding:'7px 13px',cursor:'pointer',color:'rgba(201,169,110,0.6)',fontSize:11,fontFamily:"'Jost',sans-serif"}}>← Back</button>
        </div>
        <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:12,marginBottom:20}}>
          {[['Total Professionals',s.total_users,'👥'],['Total Conversations',s.total_conversations,'💬'],['Total Messages',s.total_messages,'✉️'],['Total Leads',s.total_leads,'🎯']].map(([label,val,icon])=>(
            <div key={label} style={{background:'rgba(255,255,255,0.025)',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,padding:'16px'}}>
              <div style={{fontSize:18,marginBottom:6}}>{icon}</div>
              <div style={{fontSize:24,fontWeight:300,color:'#e8dcc8'}}>{val}</div>
              <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',textTransform:'uppercase',letterSpacing:'0.06em',marginTop:2}}>{label}</div>
            </div>
          ))}
        </div>
        <div style={{background:'rgba(201,169,110,0.06)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,padding:'18px'}}>
          <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.5)',marginBottom:8}}>Estimated Claude API Cost</div>
          <div style={{fontSize:32,fontWeight:300,color:'#c9a96e'}}>£{s.estimated_cost_gbp}</div>
          <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)',marginTop:6}}>Based on ~£0.003/message average. Check Anthropic console for exact billing.</div>
        </div>
      </div>
    </div>
  );
}
