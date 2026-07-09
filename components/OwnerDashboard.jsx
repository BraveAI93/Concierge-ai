'use client';
import { useState, useEffect } from 'react';
import { BACKEND_URL, LEGAL_FORM_FIELDS, FORM_SLUGS, FORM_KEY_TO_SLUG } from '@/lib/constants';
import BravePASettings from './BravePASettings';

export default function OwnerDashboard({ token, slug, onBack, onLogout, onEditProfile, onDataLoaded, onAddProfile, onSwitchProfile }) {
  const [profile, setProfile] = useState(null);
  const [leads, setLeads] = useState([]);
  const [convCount, setConvCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');
  const [calConnected, setCalConnected] = useState(false);
  const [calLoading, setCalLoading] = useState(false);
  const [bookingRequests, setBookingRequests] = useState([]);
  const [respondingId, setRespondingId] = useState(null);
  const [respondStatus, setRespondStatus] = useState('accepted');
  const [respondReply, setRespondReply] = useState('');
  const [respondCounter, setRespondCounter] = useState('');
  const [notes, setNotes] = useState([]);
  const [notifications, setNotifications] = useState([]);
  const [newNote, setNewNote] = useState('');
  const [newNoteType, setNewNoteType] = useState('personal');
  const [newNoteClient, setNewNoteClient] = useState('');
  const [editNoteId, setEditNoteId] = useState(null);
  const [editNoteContent, setEditNoteContent] = useState('');
  const [news, setNews] = useState([]);
  const [newsDate, setNewsDate] = useState('');
  const [newsLoading, setNewsLoading] = useState(false);
  const [digestFreq, setDigestFreq] = useState('none');
  const [digestSaving, setDigestSaving] = useState(false);
  const NOTIF_DEFAULTS = {newLead:true,hotLead:true,newBooking:true,newMessage:false,dailyBriefing:false,soundEnabled:true,soundStyle:'subtle'};
  const [notifPrefs, setNotifPrefs] = useState(() => { try { const s = localStorage.getItem('cai_notif_prefs'); return s ? {...NOTIF_DEFAULTS,...JSON.parse(s)} : NOTIF_DEFAULTS; } catch(e) { return NOTIF_DEFAULTS; } });
  const updateNotif = (key, val) => { const u = {...notifPrefs,[key]:val}; setNotifPrefs(u); localStorage.setItem('cai_notif_prefs', JSON.stringify(u)); };
  const [ownerProfiles, setOwnerProfiles] = useState([]);
  const [profilesLoading, setProfilesLoading] = useState(false);
  const [formSubmissions, setFormSubmissions] = useState([]);
  const [formSubsLoading, setFormSubsLoading] = useState(false);
  const [copiedFormSlug, setCopiedFormSlug] = useState('');
  const [expandedNews, setExpandedNews] = useState(new Set());
  const toggleNews = (i) => setExpandedNews(prev => { const s = new Set(prev); s.has(i) ? s.delete(i) : s.add(i); return s; });

  useEffect(() => {
    setProfilesLoading(true);
    fetch(`${BACKEND_URL}/owner/profiles`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.ok ? r.json() : {profiles:[]})
      .then(d => { setOwnerProfiles(d.profiles||[]); setProfilesLoading(false); })
      .catch(() => setProfilesLoading(false));
  }, [token]);

  useEffect(() => {
    if (activeTab !== 'forms') return;
    setFormSubsLoading(true);
    fetch(`${BACKEND_URL}/owner/form-submissions`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.ok ? r.json() : {submissions:[]})
      .then(d => { setFormSubmissions(d.submissions||[]); setFormSubsLoading(false); })
      .catch(() => setFormSubsLoading(false));
  }, [activeTab, token]);

  useEffect(() => {
    Promise.all([
      fetch(`${BACKEND_URL}/owner/profile`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : null),
      fetch(`${BACKEND_URL}/owner/leads`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : {leads:[]}),
      fetch(`${BACKEND_URL}/owner/conversations`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : {count:0}),
      fetch(`${BACKEND_URL}/owner/bookings`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : {bookings:[]}),
      fetch(`${BACKEND_URL}/owner/notes`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : {notes:[]}),
      fetch(`${BACKEND_URL}/owner/news`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : {items:null,date:''}),
      fetch(`${BACKEND_URL}/owner/notifications`, { headers: { Authorization: `Bearer ${token}` } }).then(r => r.ok ? r.json() : {notifications:[]}),
    ]).then(([p,l,c,br,nt,nw,nfy]) => {
      setProfile(p);
      setLeads(l.leads||[]);
      setConvCount(c.count||0);
      setBookingRequests(br.bookings||[]);
      setNotes(nt.notes||[]);
      setNotifications(nfy.notifications||[]);
      if (p) setDigestFreq(p.digest_frequency||'none');
      if (nw.items) {
        try { setNews(JSON.parse(nw.items)); setNewsDate(nw.date); } catch(e) { setNews([]); }
      } else if (p) {
        generateNewsNow(p);
      }
      if (onDataLoaded) onDataLoaded(p, l.leads||[], c.count||0);
      setLoading(false);
    }).catch(() => setLoading(false));

    const params = new URLSearchParams(typeof window !== 'undefined' ? window.location.search : '');
    if (params.get('cal') === 'connected') {
      setCalConnected(true);
      if (typeof window !== 'undefined') window.history.replaceState({}, '', window.location.pathname);
    }
  }, [token]);

  const hotLeads = leads.filter(l => l.score === 'hot').length;
  const warmLeads = leads.filter(l => l.score === 'warm').length;
  const coldLeads = leads.filter(l => l.score === 'cold' || !l.score).length;
  const pendingBookings = bookingRequests.filter(b => b.status === 'pending').length;

  const avgPrice = (() => {
    try {
      const pd = JSON.parse(profile?.profile_data||'{}');
      const svcs = pd.services||[];
      const prices = svcs.filter(s => s.priceNum).map(s => parseFloat(s.priceNum));
      return prices.length ? prices.reduce((a,b) => a+b,0)/prices.length : 0;
    } catch(e) { return 0; }
  })();
  const projRevenue = Math.round(hotLeads * avgPrice);

  const connectCalendar = () => {
    setCalLoading(true);
    localStorage.setItem('cai_owner_token', token);
    window.location.href = `${BACKEND_URL}/auth/google?Authorization=Bearer+${encodeURIComponent(token)}`;
  };

  const disconnectCalendar = async () => {
    await fetch(`${BACKEND_URL}/auth/google/disconnect`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } });
    setCalConnected(false);
  };

  const respondBooking = async (id) => {
    const body = { status: respondStatus, owner_reply: respondReply, counter_slot: respondCounter };
    await fetch(`${BACKEND_URL}/owner/bookings/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(body) });
    setBookingRequests(brs => brs.map(b => b.id===id ? {...b,...body} : b));
    setRespondingId(null); setRespondReply(''); setRespondCounter(''); setRespondStatus('accepted');
  };

  const addNote = async () => {
    if (!newNote.trim()) return;
    const r = await fetch(`${BACKEND_URL}/owner/notes`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify({ content: newNote, note_type: newNoteType, client_id: newNoteClient }) });
    const d = await r.json();
    setNotes(ns => [{id:d.id,content:newNote,note_type:newNoteType,client_id:newNoteClient,created_at:new Date().toISOString(),updated_at:new Date().toISOString()},...ns]);
    setNewNote(''); setNewNoteClient('');
  };

  const saveEditNote = async (id) => {
    await fetch(`${BACKEND_URL}/owner/notes/${id}`, { method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify({ content: editNoteContent }) });
    setNotes(ns => ns.map(n => n.id===id ? {...n,content:editNoteContent} : n));
    setEditNoteId(null);
  };

  const deleteNote = async (id) => {
    await fetch(`${BACKEND_URL}/owner/notes/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } });
    setNotes(ns => ns.filter(n => n.id !== id));
  };

  const saveDigestFreq = async () => {
    setDigestSaving(true);
    await fetch(`${BACKEND_URL}/owner/digest-prefs`, { method: 'PUT', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify({ frequency: digestFreq }) });
    setDigestSaving(false);
  };

  const generateNewsNow = async (profileOverride) => {
    setNewsLoading(true);
    try {
      const activeProfile = profileOverride || profile;
      const pd = JSON.parse(activeProfile?.profile_data||'{}');
      const services = (pd.services||[]).filter(s => s.name).map(s => s.name).join(', ');
      const r = await fetch(`${BACKEND_URL}/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          profile_id: slug, session_id: `news-gen-${Date.now()}`,
          system_prompt: 'You are a market intelligence assistant. Respond only with valid JSON array, no markdown.',
          messages: [{ role: 'user', content: `Generate 3 short news items relevant to a ${activeProfile?.profession||'professional'} based in ${activeProfile?.location||'London'} who offers: ${services}. Return JSON array: [{"title":"...","summary":"...","relevance":"..."}]` }]
        })
      });
      const d = await r.json();
      const text = d.reply || '[]';
      const parsed = JSON.parse(text.replace(/```json|```/g,'').trim());
      setNews(parsed);
      const today = new Date().toISOString().slice(0,10);
      setNewsDate(today);
      await fetch(`${BACKEND_URL}/owner/news`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify({ items: JSON.stringify(parsed), date: today }) });
    } catch(e) { setNews([]); }
    setNewsLoading(false);
  };

  const origin = typeof window !== 'undefined' ? window.location.origin : 'https://concierge-ai-gamma.vercel.app';

  const sectionLabel = {fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.12em',textTransform:'uppercase',color:'rgba(201,169,110,0.5)',marginBottom:12,marginTop:24};
  const card = {background:'rgba(255,255,255,0.05)',border:'1px solid rgba(255,215,0,0.15)',borderRadius:20,padding:'24px',marginBottom:12,backdropFilter:'blur(10px)',WebkitBackdropFilter:'blur(10px)',boxShadow:'0 8px 32px rgba(0,0,0,0.4),inset 0 1px 0 rgba(255,255,255,0.05)'};
  const tabBtn = (t) => ({padding:'7px 14px',borderRadius:20,border:'1px solid',cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.05em',transition:'all 0.2s',
    background:activeTab===t?'rgba(201,169,110,0.15)':'rgba(255,255,255,0.02)',
    borderColor:activeTab===t?'rgba(201,169,110,0.45)':'rgba(255,255,255,0.07)',
    color:activeTab===t?'#c9a96e':'rgba(232,220,200,0.38)',
    boxShadow:activeTab===t?'0 0 14px rgba(201,169,110,0.2)':'none'
  });

  if (loading) return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',alignItems:'center',justifyContent:'center'}}>
      <div style={{fontSize:28,color:'#c9a96e',animation:'bravePulse 1.8s ease-in-out infinite',textShadow:'0 0 20px rgba(201,169,110,0.6)'}}>✦</div>
    </div>
  );

  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',padding:'20px 16px 60px',fontFamily:"'Cormorant Garamond',Georgia,serif",color:'#e8dcc8'}}>
      <div style={{maxWidth:680,margin:'0 auto'}}>

        {/* Header */}
        <div className="dash-glass" style={{display:'flex',alignItems:'center',justifyContent:'space-between',marginBottom:20,padding:'16px 20px',background:'rgba(255,255,255,0.05)',border:'1px solid rgba(255,215,0,0.15)',borderRadius:20,backdropFilter:'blur(10px)',WebkitBackdropFilter:'blur(10px)',boxShadow:'0 8px 32px rgba(0,0,0,0.4),inset 0 1px 0 rgba(255,255,255,0.05)'}}>
          <div>
            <div style={{fontFamily:"'Jost',sans-serif",fontWeight:500,fontSize:14,letterSpacing:'0.08em',textTransform:'uppercase',color:'#e8dcc8'}}>{profile?.name||'My'} <span style={{letterSpacing:'0.04em'}}>Dashboard</span></div>
            <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.45)',marginTop:2}}>{profile?.profession||''} · <span style={{color:'rgba(201,169,110,0.65)'}}>/{slug}</span></div>
          </div>
          <div style={{display:'flex',gap:8}}>
            <button onClick={onBack} style={{padding:'7px 13px',background:'rgba(255,255,255,0.03)',border:'1px solid rgba(255,255,255,0.09)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.38)'}}>← Back</button>
            <button onClick={onLogout} style={{padding:'7px 13px',background:'rgba(200,80,80,0.08)',border:'1px solid rgba(200,80,80,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(220,140,140,0.7)'}}>Logout</button>
          </div>
        </div>

        <button onClick={onEditProfile} style={{width:'100%',padding:'12px 0',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.07em',textTransform:'uppercase',marginBottom:16}}>
          ✏️ Edit My Profile
        </button>

        {/* Tabs */}
        <div style={{display:'flex',gap:6,flexWrap:'wrap',marginBottom:20}}>
          {[['overview','📊 Overview'],['notifications','🔔 Notifications'+(notifications.length>0?` (${notifications.length})`:'')],['bookings','📅 Bookings'+(pendingBookings>0?` (${pendingBookings})`:'')],['notes','📝 Notes'],['forms','📋 Forms'],['news','📰 Insights'],['profiles','✦ Profiles'],['settings','⚙️ Settings']].map(([t,label])=>(
            <button key={t} onClick={()=>setActiveTab(t)} style={tabBtn(t)}>{label}</button>
          ))}
        </div>

        {/* OVERVIEW */}
        {activeTab==='overview'&&<>
          <div style={{display:'grid',gridTemplateColumns:'1fr 1fr 1fr',gap:10,marginBottom:16}}>
            {[
              [convCount,'💬','Conversations','rgba(201,169,110,0.07)','rgba(201,169,110,0.22)','#c9a96e'],
              [leads.length,'🎯','Total Leads','rgba(201,169,110,0.07)','rgba(201,169,110,0.22)','#c9a96e'],
              [hotLeads,'🔥','Hot Leads','rgba(224,122,110,0.1)','rgba(224,122,110,0.3)','#e07a6e'],
            ].map(([v,icon,label,bg,bdr,numColor])=>(
              <div key={label} className="dash-stat" style={{background:bg,border:`1px solid ${bdr}`,borderRadius:20,padding:'20px 10px 16px',textAlign:'center',backdropFilter:'blur(10px)',WebkitBackdropFilter:'blur(10px)',boxShadow:'0 8px 32px rgba(0,0,0,0.45),inset 0 1px 0 rgba(255,255,255,0.06)'}}>
                <div style={{fontSize:22,marginBottom:8,lineHeight:1,filter:'drop-shadow(0 0 6px rgba(201,169,110,0.3))'}}>{icon}</div>
                <div style={{fontSize:38,fontWeight:300,color:numColor,lineHeight:1,textShadow:`0 0 24px ${numColor}66,0 2px 8px rgba(0,0,0,0.7)`,letterSpacing:'-0.01em',marginBottom:6}}>{v}</div>
                <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.5)',textTransform:'uppercase',letterSpacing:'0.1em'}}>{label}</div>
              </div>
            ))}
          </div>

          <div className="dash-glass" style={{...card,background:'rgba(201,169,110,0.06)',border:'1px solid rgba(201,169,110,0.18)',marginBottom:16}}>
            <div style={sectionLabel}>Lead Temperature</div>
            <div style={{display:'flex',gap:24,marginTop:4,alignItems:'center'}}>
              {[['#e07a6e',hotLeads,'🔥 hot'],['#e0b56e',warmLeads,'🌤 warm'],['#8ea6c9',coldLeads,'❄️ cold']].map(([c,v,label])=>(
                <div key={label} style={{display:'flex',flexDirection:'column',alignItems:'center',gap:3}}>
                  <span style={{color:c,fontSize:28,fontWeight:300,lineHeight:1,textShadow:`0 0 16px ${c}70`}}>{v}</span>
                  <span style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)',letterSpacing:'0.06em'}}>{label}</span>
                </div>
              ))}
            </div>
          </div>

          {avgPrice>0&&<div className="dash-glass" style={{...card,background:'rgba(100,180,80,0.05)',border:'1px solid rgba(100,180,80,0.2)',marginBottom:16}}>
            <div style={sectionLabel}>Revenue Projection</div>
            <div style={{display:'flex',alignItems:'flex-end',gap:16,marginBottom:18}}>
              <div style={{fontSize:40,fontWeight:300,color:'#7aba6a',lineHeight:1,textShadow:'0 0 28px rgba(122,186,106,0.4)',letterSpacing:'-0.02em'}}>{projRevenue>0?`£${projRevenue}`:'-'}</div>
              <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.32)',lineHeight:1.55,paddingBottom:4}}>{hotLeads} hot lead{hotLeads!==1?'s':''}<br/>× avg £{Math.round(avgPrice)}</div>
            </div>
            {(()=>{const mx=Math.max(hotLeads,warmLeads,coldLeads,1);const bars=[{v:hotLeads,c:'#e07a6e',g:'rgba(224,122,110,0.35)',lbl:'Hot'},{v:warmLeads,c:'#e0b56e',g:'rgba(224,181,110,0.3)',lbl:'Warm'},{v:coldLeads,c:'#8ea6c9',g:'rgba(142,166,201,0.25)',lbl:'Cold'}];return <div style={{display:'flex',gap:8,alignItems:'flex-end',height:52}}>
              {bars.map(({v,c,g,lbl})=>{const pct=Math.max((v/mx)*44,v>0?6:0);return <div key={lbl} style={{flex:1,display:'flex',flexDirection:'column',alignItems:'center',gap:3}}>
                <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.55)',fontWeight:500}}>{v}</div>
                <div style={{width:'100%',height:`${pct}px`,background:`linear-gradient(180deg,${c},${g})`,borderRadius:'4px 4px 0 0',transition:'height .6s cubic-bezier(.22,1,.36,1)',boxShadow:`0 0 10px ${c}55`}}/>
                <div style={{fontSize:8,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)',letterSpacing:'0.06em',textTransform:'uppercase'}}>{lbl}</div>
              </div>;})}</div>;})()}
          </div>}

          <div className="dash-glass" style={{...card,background:'rgba(110,120,201,0.06)',border:'1px solid rgba(110,120,201,0.2)',marginBottom:16}}>
            <div style={sectionLabel}>AI Advice</div>
            <div style={{fontSize:13,color:'rgba(232,220,200,0.65)',lineHeight:1.7,marginTop:4}}>
              {hotLeads>0&&<div style={{marginBottom:8}}>🔥 You have {hotLeads} hot lead{hotLeads!==1?'s':''} — follow up within 24 hours while interest is high.</div>}
              {!calConnected&&<div style={{marginBottom:8}}>📅 Connect Google Calendar so your concierge can propose real availability to clients.</div>}
              {leads.length===0&&convCount===0&&<div>✦ Share your public link to start getting enquiries. Put it in your Instagram bio, WhatsApp status, and email signature.</div>}
            </div>
          </div>

          <div className="dash-glass" style={card}>
            <div style={sectionLabel}>Recent Leads</div>
            {leads.length===0&&<div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)'}}>No leads yet — share your link to start receiving enquiries.</div>}
            {leads.slice(0,8).map(l=>(
              <div key={l.id} style={{display:'flex',justifyContent:'space-between',alignItems:'center',padding:'9px 0',borderBottom:'1px solid rgba(255,255,255,0.04)'}}>
                <div>
                  <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'#e8dcc8'}}>{l.name||'Anonymous'}</div>
                  <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)'}}>{l.email||'No email'}</div>
                </div>
                <span style={{padding:'3px 9px',borderRadius:10,fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:600,background:l.score==='hot'?'rgba(224,122,110,0.15)':l.score==='warm'?'rgba(224,181,110,0.15)':'rgba(142,166,201,0.15)',color:l.score==='hot'?'#e07a6e':l.score==='warm'?'#e0b56e':'#8ea6c9'}}>{l.score||'cold'}</span>
              </div>
            ))}
          </div>

          <div className="dash-glass" style={card}>
            <div style={{display:'flex',justifyContent:'space-between',alignItems:'center',marginBottom:8}}>
              <div style={sectionLabel}>Latest Insights</div>
              {newsDate&&<div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)'}}>Updated {newsDate}</div>}
            </div>
            {news.length===0&&<div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)',lineHeight:1.6,marginBottom:10}}>
              No insights yet.
              <button onClick={()=>{setActiveTab('news');generateNewsNow();}} style={{marginLeft:8,padding:'4px 10px',background:'transparent',border:'1px solid rgba(201,169,110,0.25)',borderRadius:20,cursor:'pointer',fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.6)'}}>Generate now</button>
            </div>}
            {news.slice(0,3).map((item,i)=>{
              const open=expandedNews.has(`ov-${i}`);
              return <div key={i} className="news-item" onClick={()=>toggleNews(`ov-${i}`)} style={{padding:'11px 10px',borderRadius:10,marginBottom:i<Math.min(news.length,3)-1?4:0,background:open?'rgba(255,255,255,0.04)':'transparent',border:'1px solid transparent',userSelect:'none',cursor:'pointer'}}>
                <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.4)',marginBottom:4}}>{item.relevance}</div>
                <div style={{fontSize:13,color:'#e8dcc8',lineHeight:1.4,display:'flex',alignItems:'center',justifyContent:'space-between',gap:8}}>
                  <span>{item.title}</span>
                  <span className={`news-chevron${open?' open':''}`}>▾</span>
                </div>
                {open&&<div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.52)',lineHeight:1.6,paddingTop:8}}>{item.summary}</div>}
              </div>;
            })}
            {news.length>0&&<button onClick={(e)=>{e.stopPropagation();setActiveTab('news');}} style={{marginTop:8,width:'100%',padding:'7px 0',background:'transparent',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,cursor:'pointer',fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.45)',letterSpacing:'0.08em',textTransform:'uppercase'}}>View all in Insights →</button>}
          </div>
        </>}

        {/* NOTIFICATIONS */}
        {activeTab==='notifications'&&<>
          <div style={sectionLabel}>Owner Notifications</div>
          {notifications.length===0&&<div style={{...card,fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)'}}>No notifications yet. Sensitive-topic alerts from your Concierge chat will appear here.</div>}
          {notifications.map(n=>{
            let parsed = {};
            try { parsed = JSON.parse(n.content); } catch(e) {}
            const emailBadge = {
              sent: { label: '✉ Email sent', color: 'rgba(120,160,100,0.7)' },
              failed: { label: '✉ Email failed', color: 'rgba(200,80,80,0.7)' },
              disabled_missing_env: { label: '✉ Email not configured', color: 'rgba(232,220,200,0.35)' },
            }[parsed.email_status] || { label: 'Recorded only', color: 'rgba(232,220,200,0.3)' };
            return (
              <div key={n.id} style={{...card,border:'1px solid rgba(220,140,60,0.2)'}}>
                <div style={{display:'flex',justifyContent:'space-between',alignItems:'flex-start',marginBottom:6}}>
                  <span style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(220,140,60,0.7)',textTransform:'uppercase',letterSpacing:'0.08em'}}>🔔 {parsed.topic||'Flagged conversation'}</span>
                  <span style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.25)'}}>{new Date(n.created_at).toLocaleString()}</span>
                </div>
                {parsed.excerpt&&<div style={{fontSize:13,color:'rgba(232,220,200,0.7)',lineHeight:1.6,fontFamily:"'Cormorant Garamond',serif",fontStyle:'italic',marginBottom:6}}>"{parsed.excerpt}"</div>}
                <span style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:emailBadge.color}}>{emailBadge.label}</span>
              </div>
            );
          })}
        </>}

        {/* BOOKINGS */}
        {activeTab==='bookings'&&<>
          <div style={{...card,background:calConnected?'rgba(100,180,80,0.05)':'rgba(201,169,110,0.04)',border:`1px solid ${calConnected?'rgba(100,180,80,0.2)':'rgba(201,169,110,0.15)'}`,marginBottom:16}}>
            <div style={sectionLabel}>Google Calendar</div>
            {calConnected
              ?<div style={{display:'flex',alignItems:'center',justifyContent:'space-between'}}>
                <div style={{fontSize:13,color:'rgba(100,200,80,0.8)',fontFamily:"'Jost',sans-serif"}}>✓ Connected — concierge can read your availability</div>
                <button onClick={disconnectCalendar} style={{padding:'6px 12px',background:'rgba(200,80,80,0.1)',border:'1px solid rgba(200,80,80,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(220,140,140,0.7)'}}>Disconnect</button>
              </div>
              :<div>
                <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.45)',marginBottom:12,lineHeight:1.6}}>Connect your Google Calendar so your concierge can propose real available slots to clients.</div>
                <button onClick={connectCalendar} disabled={calLoading} style={{padding:'10px 18px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600}}>
                  {calLoading?'Connecting...':'📅 Connect Google Calendar'}
                </button>
              </div>}
          </div>
          <div style={sectionLabel}>Booking Requests</div>
          {bookingRequests.length===0&&<div style={{...card,fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)'}}>No booking requests yet.</div>}
          {bookingRequests.map(br=>(
            <div key={br.id} style={{...card,border:`1px solid ${br.status==='pending'?'rgba(201,169,110,0.25)':br.status==='accepted'?'rgba(100,180,80,0.2)':'rgba(200,80,80,0.15)'}`}}>
              <div style={{display:'flex',justifyContent:'space-between',alignItems:'flex-start',marginBottom:10}}>
                <div>
                  <div style={{fontSize:14,fontFamily:"'Jost',sans-serif",fontWeight:500,color:'#e8dcc8'}}>{br.client_name||'Client'}</div>
                  <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)'}}>{br.client_email}</div>
                </div>
                <span style={{padding:'3px 9px',borderRadius:8,fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:600,background:br.status==='pending'?'rgba(201,169,110,0.15)':br.status==='accepted'?'rgba(100,180,80,0.15)':'rgba(200,80,80,0.12)',color:br.status==='pending'?'#c9a96e':br.status==='accepted'?'#7aba6a':'#e07a6e'}}>{br.status}</span>
              </div>
              <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.55)',marginBottom:4}}>📋 {br.service_name}</div>
              <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.55)',marginBottom:2}}>🗓 Primary: <strong style={{color:'#e8dcc8'}}>{br.primary_slot}</strong></div>
              {br.backup_slot&&<div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)',marginBottom:8}}>🔄 Backup: {br.backup_slot}</div>}
              {br.message&&<div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)',fontStyle:'italic',marginBottom:10}}>"{br.message}"</div>}
              {br.status==='pending'&&respondingId!==br.id&&<div style={{display:'flex',gap:8,marginTop:8}}>
                <button onClick={()=>{setRespondingId(br.id);setRespondStatus('accepted');}} style={{flex:1,padding:'8px 0',background:'rgba(100,180,80,0.15)',border:'1px solid rgba(100,180,80,0.3)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'#7aba6a',fontWeight:600}}>✅ Accept</button>
                <button onClick={()=>{setRespondingId(br.id);setRespondStatus('counter');}} style={{flex:1,padding:'8px 0',background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'#c9a96e'}}>🔄 Counter</button>
                <button onClick={()=>{setRespondingId(br.id);setRespondStatus('declined');}} style={{flex:1,padding:'8px 0',background:'rgba(200,80,80,0.08)',border:'1px solid rgba(200,80,80,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(220,140,140,0.7)'}}>❌ Decline</button>
              </div>}
              {respondingId===br.id&&<div style={{marginTop:10,padding:'12px',background:'rgba(255,255,255,0.03)',borderRadius:10}}>
                {respondStatus==='counter'&&<input value={respondCounter} onChange={e=>setRespondCounter(e.target.value)} placeholder="e.g. Thursday 15 Jan at 3pm" style={{width:'100%',padding:'8px 11px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:8,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:8,outline:'none'}}/>}
                <textarea value={respondReply} onChange={e=>setRespondReply(e.target.value)} placeholder="Personal message to client..." rows={2} style={{width:'100%',padding:'8px 11px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.12)',borderRadius:8,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:8,outline:'none'}}/>
                <div style={{display:'flex',gap:8}}>
                  <button onClick={()=>respondBooking(br.id)} style={{flex:1,padding:'9px 0',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600}}>Send Response</button>
                  <button onClick={()=>setRespondingId(null)} style={{padding:'9px 14px',background:'rgba(255,255,255,0.03)',border:'1px solid rgba(255,255,255,0.08)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)'}}>Cancel</button>
                </div>
              </div>}
            </div>
          ))}
        </>}

        {/* NOTES */}
        {activeTab==='notes'&&<>
          <div style={sectionLabel}>Add a Note</div>
          <div style={{...card,marginBottom:16}}>
            <div style={{display:'flex',gap:8,marginBottom:10}}>
              {[['personal','📝 Personal'],['client','👤 Client']].map(([t,label])=>(
                <button key={t} onClick={()=>setNewNoteType(t)} style={{padding:'6px 12px',borderRadius:20,border:'1px solid',cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",background:newNoteType===t?'rgba(201,169,110,0.15)':'transparent',borderColor:newNoteType===t?'rgba(201,169,110,0.4)':'rgba(255,255,255,0.08)',color:newNoteType===t?'#c9a96e':'rgba(232,220,200,0.38)'}}>{label}</button>
              ))}
            </div>
            {newNoteType==='client'&&<input value={newNoteClient} onChange={e=>setNewNoteClient(e.target.value)} placeholder="Client email or identifier" style={{width:'100%',padding:'8px 11px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.15)',borderRadius:20,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:8,outline:'none'}}/>}
            <textarea value={newNote} onChange={e=>setNewNote(e.target.value)} placeholder={newNoteType==='client'?'Notes on this client...':'Personal note, reminder, observation...'} rows={3} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,color:'#e8dcc8',fontSize:13,fontFamily:"'Cormorant Garamond',serif",lineHeight:1.6,marginBottom:10,outline:'none'}}/>
            <button onClick={addNote} disabled={!newNote.trim()} style={{padding:'9px 18px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600}}>Save Note</button>
          </div>
          {notes.length===0&&<div style={{...card,fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)'}}>No notes yet.</div>}
          {notes.map(n=>(
            <div key={n.id} style={{...card,border:`1px solid ${n.note_type==='client'?'rgba(110,143,201,0.2)':'rgba(201,169,110,0.1)'}`}}>
              <div style={{display:'flex',justifyContent:'space-between',alignItems:'flex-start',marginBottom:6}}>
                <span style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:n.note_type==='client'?'rgba(110,143,201,0.7)':'rgba(201,169,110,0.45)',textTransform:'uppercase',letterSpacing:'0.08em'}}>
                  {n.note_type==='client'?`👤 ${n.client_id||'Client'}`:'📝 Personal'}
                </span>
                <div style={{display:'flex',gap:8}}>
                  <span onClick={()=>{setEditNoteId(n.id);setEditNoteContent(n.content);}} style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',cursor:'pointer',textDecoration:'underline'}}>edit</span>
                  <span onClick={()=>deleteNote(n.id)} style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(200,80,80,0.4)',cursor:'pointer',textDecoration:'underline'}}>delete</span>
                </div>
              </div>
              {editNoteId===n.id
                ?<div>
                  <textarea value={editNoteContent} onChange={e=>setEditNoteContent(e.target.value)} rows={3} style={{width:'100%',padding:'9px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,color:'#e8dcc8',fontSize:13,fontFamily:"'Cormorant Garamond',serif",lineHeight:1.6,marginBottom:8,outline:'none'}}/>
                  <div style={{display:'flex',gap:8}}>
                    <button onClick={()=>saveEditNote(n.id)} style={{padding:'7px 14px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:600}}>Save</button>
                    <button onClick={()=>setEditNoteId(null)} style={{padding:'7px 12px',background:'transparent',border:'1px solid rgba(255,255,255,0.08)',borderRadius:20,cursor:'pointer',color:'rgba(232,220,200,0.35)',fontSize:11,fontFamily:"'Jost',sans-serif"}}>Cancel</button>
                  </div>
                </div>
                :<div style={{fontSize:13,color:'rgba(232,220,200,0.7)',lineHeight:1.65,fontFamily:"'Cormorant Garamond',serif"}}>{n.content}</div>}
              <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.25)',marginTop:8}}>{new Date(n.updated_at).toLocaleDateString()}</div>
            </div>
          ))}
        </>}

        {/* NEWS */}
        {activeTab==='news'&&<>
          <div style={{display:'flex',justifyContent:'space-between',alignItems:'center',marginBottom:16}}>
            <div style={sectionLabel}>Market Intelligence</div>
            <button onClick={()=>generateNewsNow()} disabled={newsLoading} style={{padding:'7px 13px',background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'#c9a96e'}}>
              {newsLoading?'Generating...':'↻ Refresh'}
            </button>
          </div>
          {newsDate&&<div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)',marginBottom:12}}>Last updated: {newsDate}</div>}
          {news.length===0&&!newsLoading&&<div style={{...card,fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)',lineHeight:1.6}}>No insights yet. Click "Refresh" to generate today's market intelligence.</div>}
          {newsLoading&&<div style={{...card,textAlign:'center',padding:'24px'}}>
            <div style={{fontSize:20,color:'#c9a96e',animation:'pulse 1.5s infinite',marginBottom:8}}>✦</div>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)'}}>Scanning the market for you...</div>
          </div>}
          {news.map((item,i)=>{
            const open=expandedNews.has(`nt-${i}`);
            return <div key={i} className="news-item dash-glass" onClick={()=>toggleNews(`nt-${i}`)} style={{...card,background:open?'rgba(255,255,255,0.07)':'rgba(255,255,255,0.03)',marginBottom:8,cursor:'pointer',userSelect:'none'}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.12em',textTransform:'uppercase',color:'rgba(201,169,110,0.45)',marginBottom:6}}>{item.relevance}</div>
              <div style={{fontSize:15,color:'#e8dcc8',lineHeight:1.4,display:'flex',alignItems:'center',justifyContent:'space-between',gap:10}}>
                <span>{item.title}</span>
                <span style={{fontSize:12,flexShrink:0}}>▾</span>
              </div>
              {open&&<div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.55)',lineHeight:1.7,paddingTop:10,borderTop:'1px solid rgba(255,215,0,0.08)',marginTop:10}}>{item.summary}</div>}
            </div>;
          })}
          <div className="dash-glass" style={{...card,background:'rgba(201,169,110,0.05)',border:'1px solid rgba(201,169,110,0.18)',marginTop:24}}>
            <div style={sectionLabel}>Email Digest</div>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.45)',marginBottom:12,lineHeight:1.6}}>Receive insights + analytics summary directly by email.</div>
            <div style={{display:'flex',flexWrap:'wrap',gap:6,marginBottom:12}}>
              {[['none','Off'],['biweekly','Twice/week'],['weekly','Weekly'],['bimonthly','Bi-monthly'],['monthly','Monthly']].map(([v,label])=>(
                <button key={v} onClick={()=>setDigestFreq(v)} style={{padding:'6px 12px',borderRadius:20,border:'1px solid',cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",background:digestFreq===v?'rgba(201,169,110,0.15)':'transparent',borderColor:digestFreq===v?'rgba(201,169,110,0.4)':'rgba(255,255,255,0.08)',color:digestFreq===v?'#c9a96e':'rgba(232,220,200,0.38)'}}>{label}</button>
              ))}
            </div>
            <button onClick={saveDigestFreq} disabled={digestSaving} style={{padding:'8px 16px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600}}>
              {digestSaving?'Saving...':'Save preference'}
            </button>
          </div>
        </>}

        {/* SETTINGS */}
        {activeTab==='settings'&&<>
          <BravePASettings slug={slug} onConfigChange={cfg=>{}}/>
          <div style={sectionLabel}>Notifications &amp; Sound</div>
          <div style={{...card,marginBottom:16}}>
            {[['newLead','New lead notification','🔔'],['hotLead','Hot lead alert','🔥'],['newBooking','New booking request','📅'],['newMessage','New message','💬'],['dailyBriefing','Daily briefing','📊']].map(([key,label,icon])=>{
              const on = notifPrefs[key];
              return <div key={key} style={{display:'flex',alignItems:'center',justifyContent:'space-between',padding:'11px 0',borderBottom:'1px solid rgba(201,169,110,0.07)'}}>
                <div style={{display:'flex',alignItems:'center',gap:10}}>
                  <span style={{fontSize:15,lineHeight:1}}>{icon}</span>
                  <span style={{fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:400,color:'rgba(232,220,200,0.75)'}}>{label}</span>
                </div>
                <div onClick={()=>updateNotif(key,!on)} style={{width:40,height:22,borderRadius:20,cursor:'pointer',position:'relative',transition:'background 0.2s',background:on?'rgba(201,169,110,0.45)':'rgba(255,255,255,0.07)',border:'1px solid',borderColor:on?'rgba(201,169,110,0.5)':'rgba(255,255,255,0.1)'}}>
                  <div style={{position:'absolute',top:2,left:on?18:2,width:16,height:16,borderRadius:'50%',transition:'left 0.2s,background 0.2s',background:on?'#c9a96e':'rgba(255,255,255,0.22)',boxShadow:on?'0 0 6px rgba(201,169,110,0.7)':'none'}}/>
                </div>
              </div>;
            })}
          </div>

          <div style={sectionLabel}>Your Public Link</div>
          <div style={{...card,marginBottom:16}}>
            <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'#c9a96e',wordBreak:'break-all',marginBottom:10}}>
              https://concierge-ai-gamma.vercel.app/{slug}
            </div>
            <button onClick={()=>navigator.clipboard?.writeText(`https://concierge-ai-gamma.vercel.app/${slug}`)} style={{padding:'7px 14px',background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'#c9a96e'}}>
              📋 Copy link
            </button>
          </div>

          <div style={sectionLabel}>Account</div>
          <div style={{...card,marginBottom:16}}>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.45)',marginBottom:4}}>Email: {profile?.email||'—'}</div>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.45)'}}>Handle: /{slug}</div>
          </div>

          <button onClick={onLogout} style={{width:'100%',padding:'12px 0',background:'rgba(200,80,80,0.08)',border:'1px solid rgba(200,80,80,0.2)',borderRadius:20,cursor:'pointer',color:'rgba(220,140,140,0.7)',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:500}}>
            Logout
          </button>
        </>}

        {/* FORMS */}
        {activeTab==='forms'&&(()=>{
          const pdata=(()=>{try{return JSON.parse(profile?.profile_data||'{}');}catch(e){return {};}})();
          const myForms=pdata.legalForms||[];
          const copyFormLink=(formKey)=>{
            const fSlug=FORM_KEY_TO_SLUG[formKey];
            if(!fSlug) return;
            const url=`${origin}/forms/${slug}/${fSlug}`;
            navigator.clipboard?.writeText(url).catch(()=>{});
            setCopiedFormSlug(fSlug);
            setTimeout(()=>setCopiedFormSlug(''),2200);
          };
          return <>
            <div style={sectionLabel}>Sendable Form Links</div>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.38)',marginBottom:14,lineHeight:1.6}}>
              Each form has its own shareable link. Send it to clients before their session.
            </div>
            {myForms.length===0?<div style={{...card,textAlign:'center',padding:'22px 18px'}}>
              <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)',lineHeight:1.7}}>No legal forms configured yet.<br/>Edit your profile and select forms under "Client consent & legal".</div>
            </div>:<div style={{marginBottom:20}}>
              {myForms.map(formKey=>{
                const fSlug=FORM_KEY_TO_SLUG[formKey];
                if(!fSlug) return null;
                const def=LEGAL_FORM_FIELDS[formKey];
                const title=def?.en?.title||formKey;
                const formUrl=`${origin}/forms/${slug}/${fSlug}`;
                const copied=copiedFormSlug===fSlug;
                return <div key={formKey} style={{...card,marginBottom:8,padding:'13px 16px'}}>
                  <div style={{display:'flex',alignItems:'center',justifyContent:'space-between',gap:10,flexWrap:'wrap'}}>
                    <div>
                      <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'#e8dcc8',fontWeight:500,marginBottom:3}}>{title}</div>
                      <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',wordBreak:'break-all'}}>{formUrl}</div>
                    </div>
                    <div style={{display:'flex',gap:6,flexShrink:0}}>
                      <button onClick={()=>copyFormLink(formKey)} style={{padding:'6px 13px',background:copied?'rgba(100,190,120,0.15)':'rgba(201,169,110,0.08)',border:'1px solid',borderColor:copied?'rgba(100,190,120,0.35)':'rgba(201,169,110,0.2)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:copied?'rgba(140,210,150,0.9)':'#c9a96e',transition:'all 0.2s',whiteSpace:'nowrap'}}>
                        {copied?'✓ Copied':'📋 Copy link'}
                      </button>
                      <button onClick={()=>window.open(formUrl,'_blank')} style={{padding:'6px 11px',background:'transparent',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.45)'}}>
                        👁
                      </button>
                    </div>
                  </div>
                </div>;
              })}
            </div>}
            <div style={sectionLabel}>Client Responses</div>
            {formSubsLoading?<div style={{textAlign:'center',padding:'24px 0',color:'rgba(201,169,110,0.4)',fontSize:13,fontFamily:"'Jost',sans-serif"}}>Loading…</div>
            :formSubmissions.length===0?<div style={{...card,textAlign:'center',padding:'22px 18px'}}>
              <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)'}}>No responses yet — send a form link to a client to get started.</div>
            </div>:<div>
              {formSubmissions.map(sub=>{
                let resp={};try{resp=JSON.parse(sub.responses||'{}');}catch(e){}
                const def=LEGAL_FORM_FIELDS[FORM_SLUGS[sub.form_type]];
                const title=def?.en?.title||sub.form_type;
                return <div key={sub.id} style={{...card,marginBottom:8,padding:'13px 16px'}}>
                  <div style={{display:'flex',justifyContent:'space-between',alignItems:'flex-start',marginBottom:6}}>
                    <div>
                      <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'#e8dcc8',fontWeight:500}}>{sub.client_name}</div>
                      {sub.client_email&&<div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.45)',marginTop:1}}>{sub.client_email}</div>}
                    </div>
                    <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.4)',textAlign:'right',flexShrink:0}}>
                      <div>{title}</div>
                      <div>{new Date(sub.submitted_at).toLocaleDateString()}</div>
                    </div>
                  </div>
                  {resp.answer&&<div style={{fontSize:12,fontFamily:"'Cormorant Garamond',serif",color:'rgba(232,220,200,0.65)',lineHeight:1.6,padding:'8px 10px',background:'rgba(255,255,255,0.03)',borderRadius:8,borderLeft:'2px solid rgba(201,169,110,0.2)'}}>{resp.answer}</div>}
                  {resp.agreed&&<div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(100,210,140,0.7)',marginTop:6}}>✓ Agreed &amp; signed</div>}
                </div>;
              })}
            </div>}
          </>;
        })()}

        {/* PROFILES */}
        {activeTab==='profiles'&&<>
          <div style={{display:'flex',alignItems:'center',justifyContent:'space-between',marginBottom:16}}>
            <div style={{...sectionLabel,marginTop:0,marginBottom:0}}>My Profiles</div>
            <button onClick={onAddProfile} style={{padding:'7px 18px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.05em'}}>+ New Profile</button>
          </div>
          {profilesLoading&&<div style={{textAlign:'center',padding:'28px 0',color:'rgba(201,169,110,0.4)',fontSize:13,fontFamily:"'Jost',sans-serif"}}>Loading profiles…</div>}
          {ownerProfiles.map(p=>(
            <div key={p.slug} style={{...card,display:'flex',alignItems:'center',gap:12,marginBottom:10}}>
              <div style={{flex:1,minWidth:0}}>
                <div style={{fontSize:14,fontWeight:500,color:'#e8dcc8',marginBottom:2,overflow:'hidden',textOverflow:'ellipsis',whiteSpace:'nowrap'}}>{p.name}</div>
                <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.6)',marginBottom:1}}>/{p.slug}</div>
                <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.32)'}}>{p.profession||''}{p.business&&p.business!==p.name?` · ${p.business}`:''}</div>
              </div>
              <div style={{display:'flex',flexDirection:'column',gap:6,flexShrink:0}}>
                {p.slug===slug
                  ?<span style={{padding:'5px 12px',background:'rgba(201,169,110,0.07)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:20,fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.5)',textAlign:'center',letterSpacing:'0.06em',textTransform:'uppercase'}}>Active</span>
                  :<button onClick={()=>onSwitchProfile&&onSwitchProfile(p.slug)} style={{padding:'5px 12px',background:'rgba(201,169,110,0.1)',border:'1px solid rgba(201,169,110,0.3)',borderRadius:20,cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'#c9a96e'}}>Switch →</button>
                }
                <a href={`https://concierge-ai-gamma.vercel.app/${p.slug}`} target="_blank" rel="noopener noreferrer" style={{padding:'5px 12px',background:'rgba(255,255,255,0.03)',border:'1px solid rgba(255,255,255,0.08)',borderRadius:20,fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)',textDecoration:'none',textAlign:'center'}}>View ↗</a>
              </div>
            </div>
          ))}
          {!profilesLoading&&ownerProfiles.length===0&&(
            <div style={{...card,textAlign:'center',padding:'32px 18px'}}>
              <div style={{fontSize:13,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.3)',marginBottom:14}}>No profiles yet — create your first one.</div>
              <button onClick={onAddProfile} style={{padding:'9px 22px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:700}}>+ New Profile</button>
            </div>
          )}
        </>}

      </div>
    </div>
  );
}
