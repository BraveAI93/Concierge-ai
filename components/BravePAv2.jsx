'use client';
import { useState, useRef, useEffect, useMemo } from 'react';
import { BACKEND_URL } from '@/lib/constants';
import { buildBravePAPrompt, generateProactiveMessage, getQuickActions } from '@/lib/buildPrompt';
import { useFeatureFlags } from '@/lib/useFeatureFlags';

export default function BravePAv2({ token, slug, profile, leads, convCount }) {
  const { flags } = useFeatureFlags();
  const [mode, setMode] = useState('bubble');
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [paConfig, setPAConfig] = useState({ paName: 'Brave PA', personality: 'professional', proactiveAlerts: true, morningBriefing: true, soundEnabled: true });
  const [unread, setUnread] = useState(0);
  const [proactive, setProactive] = useState(null);
  const [activeTab, setActiveTab] = useState('chat');
  const [tasks, setTasks] = useState([]);
  const bottomRef = useRef(null);
  const inputRef = useRef(null);

  const pdata = useMemo(() => { try { return JSON.parse(profile?.profile_data || '{}'); } catch(e) { return {}; } }, [profile]);
  const businessData = { leads: leads || [], convCount: convCount || 0, services: pdata.services || [] };

  useEffect(() => {
    const saved = localStorage.getItem(`brave_pa_config_${slug}`);
    if (saved) try { setPAConfig(prev => ({...prev, ...JSON.parse(saved)})); } catch(e) {}
    const savedTasks = localStorage.getItem(`brave_pa_tasks_${slug}`);
    if (savedTasks) try { setTasks(JSON.parse(savedTasks)); } catch(e) {}
  }, [slug]);

  useEffect(() => {
    if (!profile || messages.length > 0) return;
    const hour = new Date().getHours();
    const greeting = hour < 12 ? 'Good morning' : hour < 17 ? 'Good afternoon' : 'Good evening';
    const firstName = profile.name?.split(' ')[0] || '';
    const hotLeads = (leads || []).filter(l => l.score === 'hot');
    let intro = `${greeting}${firstName ? `, ${firstName}` : ''}! I'm ${paConfig.paName} — your personal assistant. `;
    if (hotLeads.length > 0) intro += `Quick heads up: you have ${hotLeads.length} hot lead${hotLeads.length > 1 ? 's' : ''} waiting. `;
    intro += `What do you need?`;
    setMessages([{ role: 'assistant', content: intro, time: new Date(), type: 'text' }]);
  }, [profile]);

  useEffect(() => {
    if (!profile || !paConfig.proactiveAlerts) return;
    const timer = setTimeout(() => {
      const msg = generateProactiveMessage(businessData, paConfig);
      if (msg) {
        setProactive(msg);
        if (mode === 'bubble') setUnread(prev => prev + 1);
        setTimeout(() => setProactive(null), 12000);
      }
    }, 4000);
    return () => clearTimeout(timer);
  }, [profile, leads]);

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }); }, [messages]);

  const sendMessage = async (text) => {
    const t = (text || input).trim();
    if (!t || loading) return;
    setInput('');
    const userMsg = { role: 'user', content: t, time: new Date() };
    const newMessages = [...messages, userMsg];
    setMessages(newMessages);
    setLoading(true);
    try {
      const systemPrompt = buildBravePAPrompt(profile, paConfig, businessData, flags);
      const apiMessages = newMessages
        .filter(m => !m.type || m.type === 'text')
        .map(m => ({ role: m.role === 'assistant' ? 'assistant' : 'user', content: m.content }));
      const response = await fetch(`${BACKEND_URL}/chat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ messages: apiMessages, system_prompt: systemPrompt, profile_id: slug || 'brave-pa', session_id: `brave-pa-${slug}` })
      });
      const data = await response.json();
      const reply = data.reply || "I'm on it — give me a moment.";
      const calendarKeywords = ['add to calendar', 'schedule', 'calendar', 'book', 'appointment'];
      const suggestsCalendar = calendarKeywords.some(kw => reply.toLowerCase().includes(kw));
      const links = reply.match(/https?:\/\/[^\s)]+/g) || [];
      setMessages(prev => [...prev, { role: 'assistant', content: reply, time: new Date(), type: 'text', links, suggestsCalendar }]);
    } catch(e) {
      setMessages(prev => [...prev, { role: 'assistant', content: 'Connection blip — try again in a sec.', time: new Date(), type: 'text' }]);
    }
    setLoading(false);
    setTimeout(() => inputRef.current?.focus(), 100);
  };

  const addTask = (text) => {
    const newTask = { id: Date.now(), text, done: false, created: new Date().toISOString() };
    const updated = [...tasks, newTask];
    setTasks(updated);
    localStorage.setItem(`brave_pa_tasks_${slug}`, JSON.stringify(updated));
  };

  const toggleTask = (id) => {
    const updated = tasks.map(t => t.id === id ? {...t, done: !t.done} : t);
    setTasks(updated);
    localStorage.setItem(`brave_pa_tasks_${slug}`, JSON.stringify(updated));
  };

  const openPA = () => { setMode('floating'); setUnread(0); setProactive(null); setTimeout(() => inputRef.current?.focus(), 300); };
  const closePA = () => setMode('bubble');
  const toggleFullscreen = () => setMode(mode === 'fullscreen' ? 'floating' : 'fullscreen');

  const gold = '#C9A96E', dark = '#0C0A08', cream = '#E8DCC8';
  const quickActions = getQuickActions(businessData, new Date().getHours());

  if (!token || !profile) return null;

  const isFullscreen = mode === 'fullscreen';

  const winWidth = typeof window !== 'undefined' ? window.innerWidth : 400;
  const winHeight = typeof window !== 'undefined' ? window.innerHeight : 700;

  const BubbleEl = (
    <div style={{ position:'fixed', bottom:24, right:20, zIndex:9999, display:'flex', flexDirection:'column', alignItems:'flex-end', gap:8 }}>
      {proactive && mode === 'bubble' && (
        <div style={{ background:'rgba(12,10,8,0.97)', border:`1px solid ${gold}`, borderRadius:14, padding:'14px 16px', maxWidth:260, fontFamily:"'Jost',sans-serif", fontSize:13, color:cream, lineHeight:1.6, boxShadow:'0 12px 40px rgba(0,0,0,0.6)', animation:'paSlideUp 0.4s ease' }}>
          <div style={{ fontSize:9, color:'rgba(201,169,110,0.5)', letterSpacing:'0.12em', textTransform:'uppercase', marginBottom:8 }}>✦ {paConfig.paName}</div>
          {proactive.msg}
          <button onClick={openPA} style={{ display:'block', marginTop:10, fontSize:11, color:gold, background:'none', border:'none', cursor:'pointer', padding:0, fontFamily:"'Jost',sans-serif" }}>Reply →</button>
        </div>
      )}
      <button onClick={openPA} style={{ width:56, height:56, borderRadius:'50%', background:`linear-gradient(135deg,${gold},#7a4f0e)`, border:'none', cursor:'pointer', boxShadow:`0 6px 24px rgba(201,169,110,0.45),0 0 0 ${unread>0?'3px':'0px'} rgba(201,169,110,0.3)`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:22, color:dark, position:'relative', animation:unread>0?'paPulse 2s infinite':'none', transition:'all 0.2s ease' }}>
        ✦
        {unread > 0 && <div style={{ position:'absolute', top:-3, right:-3, width:20, height:20, borderRadius:'50%', background:'#e07a6e', color:'white', fontSize:10, fontWeight:700, fontFamily:"'Jost',sans-serif", display:'flex', alignItems:'center', justifyContent:'center', border:`2px solid ${dark}` }}>{unread}</div>}
      </button>
    </div>
  );

  const ChatWindow = mode !== 'bubble' && (
    <div style={{ position:'fixed', ...(isFullscreen?{inset:0,borderRadius:0,maxWidth:'100vw'}:{bottom:88,right:16,width:Math.min(360,winWidth-32),height:Math.min(560,winHeight-120),borderRadius:20}), background:dark, border:'1px solid rgba(201,169,110,0.2)', boxShadow:'0 24px 80px rgba(0,0,0,0.9)', display:'flex', flexDirection:'column', zIndex:9998, overflow:'hidden', fontFamily:"'Jost',sans-serif", transition:'all 0.3s cubic-bezier(0.34,1.56,0.64,1)' }}>
      {/* Header */}
      <div style={{ padding:'14px 16px', borderBottom:'1px solid rgba(201,169,110,0.1)', background:'rgba(201,169,110,0.04)', flexShrink:0, display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <div style={{ display:'flex', alignItems:'center', gap:10 }}>
          <div style={{ width:36, height:36, borderRadius:'50%', background:`linear-gradient(135deg,${gold},#7a4f0e)`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:16, color:dark, flexShrink:0 }}>✦</div>
          <div>
            <div style={{ fontSize:14, fontWeight:500, color:cream, letterSpacing:'0.04em' }}>{paConfig.paName}</div>
            <div style={{ fontSize:10, color:loading?gold:'rgba(100,200,100,0.7)', letterSpacing:'0.1em' }}>{loading?'✦ thinking...':'● online'}</div>
          </div>
        </div>
        <div style={{ display:'flex', gap:6 }}>
          <button onClick={toggleFullscreen} style={{ background:'rgba(255,255,255,0.04)', border:'1px solid rgba(255,255,255,0.08)', borderRadius:6, padding:'4px 8px', cursor:'pointer', color:'rgba(201,169,110,0.6)', fontSize:12, fontFamily:"'Jost',sans-serif" }}>{isFullscreen?'⊡':'⊞'}</button>
          <button onClick={closePA} style={{ background:'rgba(255,255,255,0.04)', border:'1px solid rgba(255,255,255,0.08)', borderRadius:6, padding:'4px 8px', cursor:'pointer', color:'rgba(232,220,200,0.4)', fontSize:14 }}>×</button>
        </div>
      </div>
      {/* Tabs */}
      <div style={{ display:'flex', borderBottom:'1px solid rgba(255,255,255,0.06)', flexShrink:0 }}>
        {[['chat','💬 Chat'],['tasks','✓ Tasks'],['calendar','📅 Calendar']].map(([tab,label]) => (
          <button key={tab} onClick={()=>setActiveTab(tab)} style={{ flex:1, padding:'8px 4px', background:'transparent', border:'none', borderBottom:`2px solid ${activeTab===tab?gold:'transparent'}`, cursor:'pointer', fontSize:11, fontFamily:"'Jost',sans-serif", color:activeTab===tab?gold:'rgba(232,220,200,0.35)', transition:'all 0.2s' }}>{label}</button>
        ))}
      </div>
      {/* Chat */}
      {activeTab === 'chat' && <>
        <div style={{ flex:1, overflowY:'auto', padding:'14px 12px', display:'flex', flexDirection:'column', gap:14 }}>
          {messages.map((msg, i) => (
            <div key={i} style={{ display:'flex', flexDirection:msg.role==='user'?'row-reverse':'row', gap:8, alignItems:'flex-end' }}>
              {msg.role === 'assistant' && <div style={{ width:26, height:26, borderRadius:'50%', background:`linear-gradient(135deg,${gold},#7a4f0e)`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:11, flexShrink:0, color:dark }}>✦</div>}
              <div style={{ maxWidth:'82%' }}>
                <div style={{ background:msg.role==='user'?'rgba(201,169,110,0.12)':'rgba(255,255,255,0.04)', border:`1px solid ${msg.role==='user'?'rgba(201,169,110,0.2)':'rgba(255,255,255,0.06)'}`, borderRadius:msg.role==='user'?'16px 16px 4px 16px':'16px 16px 16px 4px', padding:'10px 13px', fontSize:13, color:cream, lineHeight:1.65, whiteSpace:'pre-wrap' }}>{msg.content}</div>
                {msg.suggestsCalendar && msg.role==='assistant' && <button onClick={()=>sendMessage('Yes, please add it to my calendar')} style={{ marginTop:6, padding:'5px 10px', background:'rgba(201,169,110,0.08)', border:'1px solid rgba(201,169,110,0.2)', borderRadius:8, cursor:'pointer', fontSize:11, color:gold, fontFamily:"'Jost',sans-serif" }}>📅 Add to calendar</button>}
                <div style={{ fontSize:9, color:'rgba(201,169,110,0.25)', marginTop:4, textAlign:msg.role==='user'?'right':'left', fontFamily:"'Jost',sans-serif" }}>{new Date(msg.time).toLocaleTimeString('en-GB',{hour:'2-digit',minute:'2-digit'})}</div>
              </div>
            </div>
          ))}
          {loading && (
            <div style={{ display:'flex', gap:8, alignItems:'flex-end' }}>
              <div style={{ width:26, height:26, borderRadius:'50%', background:`linear-gradient(135deg,${gold},#7a4f0e)`, display:'flex', alignItems:'center', justifyContent:'center', fontSize:11, color:dark }}>✦</div>
              <div style={{ background:'rgba(255,255,255,0.04)', border:'1px solid rgba(255,255,255,0.06)', borderRadius:'16px 16px 16px 4px', padding:'12px 16px' }}>
                <div style={{ display:'flex', gap:5, alignItems:'center' }}>{[0,1,2].map(i=><div key={i} style={{ width:7, height:7, borderRadius:'50%', background:gold, animation:`paBounce 1.2s ${i*0.2}s infinite ease-in-out` }}/>)}</div>
              </div>
            </div>
          )}
          {messages.length <= 2 && !loading && (
            <div style={{ marginTop:4 }}>
              <div style={{ fontSize:9, color:'rgba(201,169,110,0.35)', letterSpacing:'0.12em', textTransform:'uppercase', marginBottom:10 }}>Suggested</div>
              <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:7 }}>
                {quickActions.map((action,i) => <button key={i} onClick={()=>sendMessage(action.msg)} style={{ background:'rgba(201,169,110,0.05)', border:'1px solid rgba(201,169,110,0.12)', borderRadius:10, padding:'9px 10px', cursor:'pointer', fontSize:11, color:'rgba(201,169,110,0.7)', textAlign:'left', fontFamily:"'Jost',sans-serif", lineHeight:1.4 }}>{action.label}</button>)}
              </div>
            </div>
          )}
          <div ref={bottomRef}/>
        </div>
        <div style={{ padding:'10px 12px', borderTop:'1px solid rgba(255,255,255,0.06)', display:'flex', gap:8, alignItems:'flex-end', background:'rgba(0,0,0,0.3)', flexShrink:0 }}>
          <textarea ref={inputRef} value={input} onChange={e=>setInput(e.target.value)} onKeyDown={e=>{if(e.key==='Enter'&&!e.shiftKey){e.preventDefault();sendMessage();}}} placeholder={`Ask ${paConfig.paName} anything...`} rows={1} style={{ flex:1, background:'rgba(255,255,255,0.05)', border:'1px solid rgba(201,169,110,0.12)', borderRadius:12, padding:'9px 13px', color:cream, fontSize:13, fontFamily:"'Jost',sans-serif", resize:'none', outline:'none', maxHeight:100, overflowY:'auto', lineHeight:1.5 }}/>
          <button onClick={()=>sendMessage()} disabled={!input.trim()||loading} style={{ width:38, height:38, borderRadius:'50%', flexShrink:0, background:input.trim()&&!loading?`linear-gradient(135deg,${gold},#7a4f0e)`:'rgba(255,255,255,0.05)', border:'none', cursor:input.trim()&&!loading?'pointer':'default', color:input.trim()&&!loading?dark:'rgba(255,255,255,0.15)', fontSize:16, transition:'all 0.2s', display:'flex', alignItems:'center', justifyContent:'center' }}>→</button>
        </div>
      </>}
      {/* Tasks */}
      {activeTab === 'tasks' && (
        <div style={{ flex:1, overflowY:'auto', padding:14, display:'flex', flexDirection:'column', gap:10 }}>
          <input placeholder="Add a task..." onKeyDown={e=>{if(e.key==='Enter'&&e.target.value.trim()){addTask(e.target.value.trim());e.target.value='';}}} style={{ padding:'8px 12px', background:'rgba(255,255,255,0.05)', border:'1px solid rgba(201,169,110,0.12)', borderRadius:10, color:cream, fontSize:13, fontFamily:"'Jost',sans-serif", outline:'none' }}/>
          {tasks.length===0 && <div style={{ fontSize:12, color:'rgba(232,220,200,0.3)', textAlign:'center', marginTop:20 }}>No tasks yet. Ask your PA to create some.</div>}
          {tasks.map(task => (
            <div key={task.id} onClick={()=>toggleTask(task.id)} style={{ display:'flex', alignItems:'center', gap:10, padding:'10px 12px', background:'rgba(255,255,255,0.03)', border:'1px solid rgba(255,255,255,0.06)', borderRadius:10, cursor:'pointer', opacity:task.done?0.5:1 }}>
              <div style={{ width:18, height:18, borderRadius:'50%', border:`1.5px solid ${task.done?gold:'rgba(201,169,110,0.3)'}`, background:task.done?gold:'transparent', display:'flex', alignItems:'center', justifyContent:'center', flexShrink:0, fontSize:10, color:dark }}>{task.done?'✓':''}</div>
              <span style={{ fontSize:13, color:cream, textDecoration:task.done?'line-through':'none', fontFamily:"'Jost',sans-serif" }}>{task.text}</span>
            </div>
          ))}
          <button onClick={()=>{setActiveTab('chat');sendMessage('Help me prioritise my tasks and add any I might be missing for today');}} style={{ marginTop:8, padding:'9px 14px', background:'rgba(201,169,110,0.08)', border:'1px solid rgba(201,169,110,0.2)', borderRadius:10, cursor:'pointer', fontSize:12, color:gold, fontFamily:"'Jost',sans-serif" }}>✦ Ask PA to help prioritise</button>
        </div>
      )}
      {/* Calendar */}
      {activeTab === 'calendar' && (
        <div style={{ flex:1, overflowY:'auto', padding:14, display:'flex', flexDirection:'column', gap:10 }}>
          <div style={{ fontSize:12, color:'rgba(232,220,200,0.5)', lineHeight:1.6 }}>{profile?.google_refresh_token?'Your Google Calendar is connected.':'Connect Google Calendar to see your events here.'}</div>
          <button onClick={()=>{setActiveTab('chat');sendMessage("What do I have on my calendar today and this week?");}} style={{ padding:'9px 14px', background:'rgba(201,169,110,0.08)', border:'1px solid rgba(201,169,110,0.2)', borderRadius:10, cursor:'pointer', fontSize:12, color:gold, fontFamily:"'Jost',sans-serif" }}>📅 Ask PA about my schedule</button>
          <button onClick={()=>{setActiveTab('chat');sendMessage('Add a new event to my calendar');}} style={{ padding:'9px 14px', background:'rgba(255,255,255,0.03)', border:'1px solid rgba(255,255,255,0.08)', borderRadius:10, cursor:'pointer', fontSize:12, color:'rgba(232,220,200,0.5)', fontFamily:"'Jost',sans-serif" }}>+ Add event via PA</button>
        </div>
      )}
    </div>
  );

  return (
    <>
      <style>{`
        @keyframes paSlideUp { from{opacity:0;transform:translateY(12px);}to{opacity:1;transform:translateY(0);} }
        @keyframes paPulse { 0%,100%{box-shadow:0 6px 24px rgba(201,169,110,0.45);}50%{box-shadow:0 6px 32px rgba(201,169,110,0.7),0 0 0 6px rgba(201,169,110,0.1);} }
        @keyframes paBounce { 0%,60%,100%{transform:translateY(0);}30%{transform:translateY(-5px);} }
      `}</style>
      {BubbleEl}
      {ChatWindow}
    </>
  );
}
