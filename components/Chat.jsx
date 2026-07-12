'use client';
import { useState, useRef, useEffect } from 'react';
import { BACKEND_URL, MAX_TURNS, T, convertCurrency } from '@/lib/constants';
import LegalFormModal from './LegalFormModal';

async function callAPI(messages, systemPrompt, profileId) {
  const r = await fetch(`${BACKEND_URL}/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      messages,
      system_prompt: systemPrompt,
      profile_id: profileId,
      session_id: sessionStorage.getItem('cai_session') || 'anonymous'
    })
  });
  if (!r.ok) { const e = await r.json().catch(() => ({})); throw new Error(e.error || 'Connection error'); }
  const d = await r.json();
  return { reply: d.reply || 'Something went wrong — please try again.', conversationId: d.conversation_id || '' };
}

// Returns the real backend delivery outcome. Never assumes success —
// a network failure or non-2xx response is reported as record failure,
// regardless of what the email sub-status would have been.
async function sendAlert(profileId, topic, excerpt, conversationId) {
  try {
    const r = await fetch(`${BACKEND_URL}/alert`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        profile_id: profileId,
        session_id: sessionStorage.getItem('cai_session') || 'anonymous',
        conversation_id: conversationId || '',
        topic,
        excerpt
      })
    });
    const d = await r.json().catch(() => ({}));
    if (!r.ok || !d.dashboard_record_created) return { recorded: false, emailStatus: null };
    return { recorded: true, emailStatus: d.email_status || null };
  } catch (e) {
    return { recorded: false, emailStatus: null };
  }
}

function Dots({ color }) {
  return (
    <div style={{display:'flex',gap:5,alignItems:'center',padding:'4px 2px'}}>
      {[0,1,2].map(i=>(
        <div key={i} style={{width:6,height:6,borderRadius:'50%',background:color||'#c9a96e',animation:'typeDot 1.2s infinite',animationDelay:`${i*0.18}s`}}/>
      ))}
    </div>
  );
}

function ReviewModal({ reply, name, acc, onApprove, onReject }) {
  const [ed, setEd] = useState(reply);
  const [editing, setEditing] = useState(false);
  return (
    <div style={{position:'fixed',inset:0,background:'rgba(0,0,0,0.75)',display:'flex',alignItems:'center',justifyContent:'center',zIndex:9999,padding:16}}>
      <div style={{width:'100%',maxWidth:500,background:'#181410',border:'1px solid rgba(201,169,110,0.22)',borderRadius:20,overflow:'hidden'}}>
        <div style={{padding:'16px 20px',borderBottom:'1px solid rgba(201,169,110,0.1)',display:'flex',alignItems:'center',gap:9}}>
          <span>👁</span>
          <div>
            <div style={{fontFamily:"'Jost',sans-serif",fontWeight:500,fontSize:11,letterSpacing:'0.09em',textTransform:'uppercase',color:acc}}>Human Review — {name}</div>
            <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(232,220,200,0.35)',marginTop:2}}>Approve, edit or block</div>
          </div>
        </div>
        <div style={{padding:'16px 20px'}}>
          {editing
            ? <textarea value={ed} onChange={e=>setEd(e.target.value)} style={{width:'100%',minHeight:110,background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.25)',borderRadius:20,padding:'11px 13px',color:'#e8dcc8',fontSize:14,fontFamily:"'Cormorant Garamond',serif",lineHeight:1.65,resize:'vertical'}}/>
            : <div style={{background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.09)',borderRadius:20,padding:'12px 14px',fontSize:14,lineHeight:1.7,color:'#e8dcc8',fontFamily:"'Cormorant Garamond',serif",whiteSpace:'pre-wrap'}}>{ed}</div>
          }
        </div>
        <div style={{padding:'0 20px 18px',display:'flex',gap:8}}>
          <button onClick={()=>onApprove(ed)} style={{flex:1,padding:'10px 0',background:`linear-gradient(135deg,${acc},#6b3e10)`,border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.07em',textTransform:'uppercase'}}>Send</button>
          <button onClick={()=>setEditing(!editing)} style={{flex:1,padding:'10px 0',background:'rgba(201,169,110,0.08)',border:`1px solid ${acc}44`,borderRadius:20,cursor:'pointer',color:acc,fontSize:12,fontFamily:"'Jost',sans-serif",letterSpacing:'0.07em',textTransform:'uppercase'}}>{editing?'Preview':'Edit'}</button>
          <button onClick={onReject} style={{padding:'10px 14px',background:'rgba(180,60,60,0.08)',border:'1px solid rgba(180,60,60,0.2)',borderRadius:20,cursor:'pointer',color:'#c06060',fontSize:12,fontFamily:"'Jost',sans-serif",letterSpacing:'0.07em',textTransform:'uppercase'}}>Block</button>
        </div>
      </div>
    </div>
  );
}

export default function Chat({ profile, systemPrompt, profileId, onBack, lang }) {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [alerts, setAlerts] = useState([]);
  // null | 'pending' | 'confirmed' | 'recorded_email_failed' | 'recorded_email_disabled' | 'failed'
  const [alertStatus, setAlertStatus] = useState(null);
  const [reviewData, setReviewData] = useState(null);
  const [blocked, setBlocked] = useState(0);
  const [turns, setTurns] = useState(0);
  const [err, setErr] = useState(null);
  const [leadCaptured, setLeadCaptured] = useState(false);
  const [showLeadForm, setShowLeadForm] = useState(false);
  const [leadName, setLeadName] = useState('');
  const [leadEmail, setLeadEmail] = useState('');
  const [consented, setConsented] = useState(() => {
    if (typeof window === 'undefined') return false;
    return !!sessionStorage.getItem('cai_consent');
  });
  const [showLegalForm, setShowLegalForm] = useState(false);
  const [legalFormDone, setLegalFormDone] = useState(false);
  const [consentSubmitting, setConsentSubmitting] = useState(false);
  const [consentError, setConsentError] = useState(false);
  const bottomRef = useRef(null);
  const inputRef = useRef(null);
  const bookingSubmitted = useRef(false);
  const L = lang || 'en';
  const TL = T[L] || T.en;
  const P = profile;
  const acc = P.accent || '#c9a96e';
  const profileLegalForms = P.legalForms || [];

  const alertMap = {
    'Injuries & pain':['injury','pain','hurt','infortun','dolor'],
    'Pregnancy':['pregnant','pregnancy','incinta','gravidanz'],
    'Price negotiation':['discount','cheaper','sconto','free','gratis','too expensive','troppo caro'],
    'Refunds':['refund','complaint','rimborso','reclamo'],
    'Home visits':['home visit','outcall','hotel','domicilio'],
    'Naturist':['naturist','nude','naked'],
    'Legal advice':['specific advice','my case','processo','specifico'],
    'Allergies':['allerg'],
    'IP & usage rights':['rights','usage','license','copyright'],
  };

  const giveConsent = async () => {
    setConsentSubmitting(true);
    setConsentError(false);
    try {
      const r = await fetch(`${BACKEND_URL}/consent`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          profile_id: profileId,
          session_id: sessionStorage.getItem('cai_session') || 'anonymous',
          client_name: '',
          client_email: '',
          forms_agreed: ['ai_processing'],
          answers: '',
          signature_date: new Date().toISOString()
        })
      });
      if (!r.ok) throw new Error('consent save failed');
      sessionStorage.setItem('cai_consent', new Date().toISOString());
      setConsented(true);
    } catch (e) {
      setConsentError(true);
    } finally {
      setConsentSubmitting(false);
    }
  };

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }); }, [messages, loading]);

  useEffect(() => {
    if (!sessionStorage.getItem('cai_session')) {
      sessionStorage.setItem('cai_session', crypto.randomUUID());
    }
    const q = new URLSearchParams(window.location.search).get('q');
    if (q && messages.length === 0) {
      setTimeout(() => send(decodeURIComponent(q)), 800);
    }
  }, []);

  const send = async (text) => {
    const t = text || input.trim();
    if (!t || loading || turns >= MAX_TURNS) return;
    setInput(''); setErr(null);
    const lower = t.toLowerCase();
    const found = (P.sensitiveTopics||[]).filter(tp => (alertMap[tp]||[]).some(w => lower.includes(w)));
    if (found.length) { setAlerts(found); setAlertStatus('pending'); }
    const newMsgs = [...messages, { role: 'user', content: t }];
    setMessages(newMsgs);
    setLoading(true);
    const newTurn = turns + 1;
    setTurns(newTurn);
    try {
      const { reply, conversationId } = await callAPI(newMsgs, systemPrompt, profileId);
      setLoading(false);
      if (found.length) {
        setReviewData({ reply, snap: newMsgs });
        sendAlert(profileId, found.join(', '), t, conversationId).then(({ recorded, emailStatus }) => {
          if (!recorded) { setAlertStatus('failed'); return; }
          if (emailStatus === 'sent') setAlertStatus('confirmed');
          else if (emailStatus === 'disabled_missing_env') setAlertStatus('recorded_email_disabled');
          else setAlertStatus('recorded_email_failed'); // covers 'failed' and any unexpected value
        });
      }
      else setMessages([...newMsgs, { role: 'assistant', content: convertCurrency(reply, L) }]);
      if (newTurn === 3 && !leadCaptured) setTimeout(() => setShowLeadForm(true), 1200);
      const replyLower = reply.toLowerCase();
      if ((replyLower.includes('sent your request') || replyLower.includes("you'll hear back")) && !bookingSubmitted.current) {
        bookingSubmitted.current = true;
        const allText = newMsgs.map(m => m.content).join(' ');
        const emailMatch = allText.match(/[\w.-]+@[\w.-]+\.\w+/);
        const clientEmail = emailMatch ? emailMatch[0] : '';
        const nameMatch = allText.match(/(?:i'm|i am|my name is|this is)\s+([A-Z][a-z]+)/i);
        const clientName = nameMatch ? nameMatch[1] : '';
        const slotPatterns = [/(?:monday|tuesday|wednesday|thursday|friday|saturday|sunday)[^.!?]*/gi, /\d{1,2}(?:st|nd|rd|th)?\s+(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[^.!?]*/gi];
        const slots = [];
        for (const pattern of slotPatterns) { const m = allText.match(pattern); if (m) slots.push(...m.slice(0, 2)); }
        if (clientEmail || clientName) {
          fetch(`${BACKEND_URL}/booking-request`, {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ profile_id: profileId, session_id: sessionStorage.getItem('cai_session')||'anonymous', client_name: clientName, client_email: clientEmail, primary_slot: slots[0]||'To be confirmed', backup_slot: slots[1]||'', message: 'Booking request via concierge chat' })
          }).catch(() => {});
        }
      }
    } catch(e) { setLoading(false); setErr(e.message || 'Connection error'); }
    setTimeout(() => inputRef.current?.focus(), 100);
  };

  const submitLead = async () => {
    if (!leadName.trim() && !leadEmail.trim()) { setShowLeadForm(false); setLeadCaptured(true); return; }
    setLeadCaptured(true); setShowLeadForm(false);
    try {
      await fetch(`${BACKEND_URL}/lead`, { method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ profile_id: profileId, session_id: sessionStorage.getItem('cai_session')||'anonymous', name: leadName, email: leadEmail }) });
    } catch(e) {}
    setMessages(m => [...m, { role: 'assistant', content: L==='it' ? `Grazie ${leadName||''}! ${P.name} ti contatterà presto. Nel frattempo continua pure a fare domande 😊` : `Thanks ${leadName||''}! ${P.name} will be in touch soon. Feel free to keep asking questions in the meantime 😊` }]);
  };

  const approve = (r) => { setMessages([...reviewData.snap, { role: 'assistant', content: convertCurrency(r, L) }]); setReviewData(null); };
  const reject = () => { setBlocked(b => b+1); setMessages([...reviewData.snap, { role: 'assistant', content: L==='it'?`È qualcosa di cui ${P.name} vorrà parlare direttamente con te — ti contatterà a breve.`:`That's something ${P.name} will want to discuss personally — they'll be in touch shortly.`, _blocked: true }]); setReviewData(null); };
  const hk = (e) => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send(); } };
  const pct = (turns / MAX_TURNS) * 100;
  const tagline = typeof P.tagline === 'object' ? P.tagline[L] || P.tagline.en : P.tagline;
  const qas = P.quickActions?.[L] || P.quickActions?.en || [];

  if (!consented) {
    const isIT = L === 'it';
    return (
      <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',flexDirection:'column',alignItems:'center',justifyContent:'center',padding:'24px 16px',fontFamily:"'Cormorant Garamond',Georgia,serif"}}>
        <div style={{width:'100%',maxWidth:420,background:'rgba(12,10,8,0.82)',border:'1px solid rgba(201,169,110,0.25)',borderRadius:20,padding:'28px 24px',animation:'fadeUp 0.4s ease forwards',backdropFilter:'blur(20px)',boxShadow:'0 0 40px rgba(201,169,110,0.08),0 16px 48px rgba(0,0,0,0.6)'}}>
          <div style={{textAlign:'center',marginBottom:22}}>
            <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",letterSpacing:'0.18em',textTransform:'uppercase',color:'rgba(201,169,110,0.38)',marginBottom:10}}>✦ {P.business||P.name}</div>
            <div style={{fontSize:22,fontWeight:500,color:'#e8dcc8',marginBottom:4}}>{isIT?'Prima di iniziare':'Before we begin'}</div>
            <div style={{fontSize:13,fontStyle:'italic',color:'rgba(232,220,200,0.42)',lineHeight:1.65}}>{isIT?'Questo assistente è alimentato da intelligenza artificiale e risponde per conto di':'This assistant is powered by AI and responds on behalf of'} {P.name}.</div>
          </div>
          <div style={{background:'rgba(201,169,110,0.04)',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,padding:'14px 16px',marginBottom:16}}>
            <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.5)',marginBottom:10}}>{isIT?'Come usiamo i tuoi dati':'How we use your data'}</div>
            {[
              [isIT?'Conversazioni AI':'AI Conversations', isIT?'I messaggi sono elaborati da Anthropic (Claude AI) per generare risposte.':'Messages are processed by Anthropic (Claude AI) to generate responses. Not stored after your session.'],
              [isIT?'Nessun dato personale richiesto':'No personal data required', isIT?'Non ti chiediamo nome, email o dati di pagamento.':'We do not ask for your name, email or payment details.'],
              [isIT?'Sessione temporanea':'Temporary session', isIT?'La conversazione esiste solo durante questa sessione.':'The conversation exists only during this session. Closing the window deletes it.'],
              [isIT?'Nessun cookie di profilazione':'No tracking cookies', isIT?'Non utilizziamo cookie di tracciamento o pubblicità.':'We do not use tracking or advertising cookies.'],
            ].map(([title,desc]) => (
              <div key={title} style={{display:'flex',gap:10,marginBottom:9}}>
                <div style={{color:'rgba(201,169,110,0.6)',fontSize:13,flexShrink:0,marginTop:1}}>✓</div>
                <div>
                  <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:500,color:'rgba(232,220,200,0.75)',marginBottom:2}}>{title}</div>
                  <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(232,220,200,0.38)',lineHeight:1.55}}>{desc}</div>
                </div>
              </div>
            ))}
          </div>
          <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(232,220,200,0.28)',lineHeight:1.6,marginBottom:18,textAlign:'center'}}>
            {isIT?'Continuando accetti i nostri':'By continuing you agree to our'}{' '}
            <a href="/privacy" target="_blank" style={{color:'rgba(201,169,110,0.5)',textDecoration:'underline'}}>{isIT?'Termini di Servizio':'Terms of Service'}</a>
            {' & '}
            <a href="/privacy" target="_blank" style={{color:'rgba(201,169,110,0.5)',textDecoration:'underline'}}>{isIT?'Privacy Policy':'Privacy Policy'}</a>.
            {' '}{isIT?'Elaborazione dati conforme al UK GDPR / GDPR EU.':'Data processing compliant with UK GDPR / EU GDPR.'}
          </div>
          {consentError && (
            <div style={{background:'rgba(180,60,60,0.08)',border:'1px solid rgba(180,60,60,0.22)',borderRadius:12,padding:'9px 12px',marginBottom:10,fontSize:11,fontFamily:"'Jost',sans-serif",color:'#c08080',textAlign:'center',lineHeight:1.5}}>
              {isIT?'Impossibile salvare il consenso. Riprova.':'Could not save your consent. Please try again.'}
            </div>
          )}
          <button onClick={giveConsent} disabled={consentSubmitting} style={{width:'100%',padding:'13px 0',background:consentSubmitting?'rgba(201,169,110,0.25)':`linear-gradient(135deg,${acc},#7a4f0e)`,border:'none',borderRadius:20,cursor:consentSubmitting?'default':'pointer',color:'#0c0a08',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.08em',textTransform:'uppercase',marginBottom:10}}>
            {consentSubmitting ? (isIT?'Salvataggio…':'Saving…') : consentError ? (isIT?'Riprova ✦':'Retry ✦') : (isIT?'Accetto — Inizia la chat ✦':'I agree — Start chatting ✦')}
          </button>
          <div style={{textAlign:'center'}}>
            <button onClick={onBack} style={{background:'none',border:'none',cursor:'pointer',fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.28)',textDecoration:'underline'}}>
              {isIT?'← Torna indietro':'← Go back'}
            </button>
          </div>
          <div style={{textAlign:'center',marginTop:16,fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.18)',letterSpacing:'0.07em',textTransform:'uppercase'}}>
            Powered by The Concierge · Built with Claude (Anthropic)
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{minHeight:'100vh',background:'#0c0a08',display:'flex',flexDirection:'column',alignItems:'center',fontFamily:"'Cormorant Garamond',Georgia,serif",color:'#e8dcc8'}}>
      {/* HEADER */}
      <div style={{width:'100%',maxWidth:680,padding:'15px 18px 9px',borderBottom:'1px solid rgba(201,169,110,0.12)',position:'sticky',top:0,background:'rgba(12,10,8,0.97)',backdropFilter:'blur(20px)',zIndex:10,boxShadow:'0 4px 24px rgba(0,0,0,0.4)'}}>
        <div style={{display:'flex',alignItems:'center',gap:11}}>
          <div style={{width:38,height:38,borderRadius:20,flexShrink:0,background:P.gradient,display:'flex',alignItems:'center',justifyContent:'center',fontSize:17,boxShadow:`0 4px 14px ${acc}33`}}>{P.emoji}</div>
          <div style={{flex:1}}>
            <div style={{fontSize:14,fontWeight:500,letterSpacing:'0.06em',textTransform:'uppercase',color:'#e8dcc8',fontFamily:"'Jost',sans-serif",lineHeight:1.2}}>{P.business||P.name}</div>
            <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",fontWeight:300,color:`${acc}88`,letterSpacing:'0.09em',textTransform:'uppercase',marginTop:2}}>{P.profession} · {P.location}</div>
          </div>
          <div style={{display:'flex',alignItems:'center',gap:7}}>
            {blocked > 0 && <div style={{padding:'2px 7px',borderRadius:20,background:'rgba(180,60,60,0.11)',border:'1px solid rgba(180,60,60,0.2)',fontSize:10,fontFamily:"'Jost',sans-serif",color:'#c06060'}}>{blocked} blocked</div>}
            <button onClick={onBack} style={{padding:'4px 10px',background:'rgba(255,255,255,0.03)',border:'1px solid rgba(255,255,255,0.08)',borderRadius:20,cursor:'pointer',fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.35)'}}>← {L==='it'?'Indietro':'Back'}</button>
            <div style={{display:'flex',alignItems:'center',gap:4,fontSize:10,fontFamily:"'Jost',sans-serif",color:'#7a9e6e'}}>
              <div style={{width:5,height:5,borderRadius:'50%',background:'#7a9e6e',boxShadow:'0 0 5px #7a9e6e',animation:'pulse 2s infinite'}}/>Live
            </div>
          </div>
        </div>
        <div style={{marginTop:7}}>
          <div className="turn-bar"><div className="turn-fill" style={{width:`${pct}%`}}/></div>
          <div style={{display:'flex',justifyContent:'flex-end',marginTop:3}}><span style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.26)',letterSpacing:'0.05em'}}>{MAX_TURNS-turns} {TL.msgRemaining}</span></div>
        </div>
      </div>

      {/* BODY */}
      <div style={{width:'100%',maxWidth:680,flex:1,display:'flex',flexDirection:'column',padding:'0 15px',minHeight:'calc(100vh - 165px)'}}>
        {messages.length === 0 && (
          <div style={{padding:'24px 0 13px',animation:'fadeUp 0.5s ease forwards'}}>
            <div style={{textAlign:'center',marginBottom:22}}>
              <div style={{fontSize:30,color:acc,marginBottom:8,fontWeight:300}}>Benvenuto.</div>
              {tagline && <div style={{fontSize:14,fontStyle:'italic',color:'rgba(232,220,200,0.4)',lineHeight:1.7,maxWidth:310,margin:'0 auto 5px'}}>"{tagline}"</div>}
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(232,220,200,0.17)',letterSpacing:'0.1em',textTransform:'uppercase'}}>— {P.name}</div>
            </div>
            <div style={{display:'grid',gridTemplateColumns:'1fr 1fr',gap:6,marginBottom:8}}>
              {qas.map(q => (
                <button key={q.label} className="qb" onClick={() => send(q.msg)} style={{padding:'9px 11px',background:'rgba(201,169,110,0.045)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:20,cursor:'pointer',color:'#e8dcc8',textAlign:'left',fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:400,lineHeight:1.4,transition:'all 0.2s'}}>{q.label}</button>
              ))}
            </div>
            <button onClick={() => send(L==='it'?"Non trovo quello che cerco — puoi aiutarmi?":"I can't find what I need — can you help?")} style={{width:'100%',padding:'11px 15px',background:'rgba(201,169,110,0.04)',border:'1px dashed rgba(201,169,110,0.2)',borderRadius:10,cursor:'pointer',display:'flex',alignItems:'center',gap:10,transition:'opacity 0.2s'}}>
              <div style={{width:28,height:28,borderRadius:'50%',background:'rgba(201,169,110,0.09)',display:'flex',alignItems:'center',justifyContent:'center',fontSize:13,flexShrink:0}}>✧</div>
              <div style={{textAlign:'left'}}>
                <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",fontWeight:500,color:acc}}>{TL.cantFind}</div>
                <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:300,color:'rgba(232,220,200,0.32)',marginTop:1}}>{TL.tellGoals}</div>
              </div>
              <div style={{marginLeft:'auto',color:`${acc}38`,fontSize:14}}>→</div>
            </button>
          </div>
        )}

        {alerts.length > 0 && (
          <div style={{paddingTop:messages.length>0?11:0}}>
            <div style={{background:'rgba(220,140,60,0.08)',border:'1px solid rgba(220,140,60,0.25)',borderRadius:20,padding:'9px 12px',marginBottom:10}}>
              <div style={{display:'flex',justifyContent:'space-between',marginBottom:5}}>
                <span style={{fontSize:9,fontFamily:"'Jost',sans-serif",fontWeight:500,letterSpacing:'0.1em',textTransform:'uppercase',color:alertStatus==='failed'?'#c06060':'#dc8c3c'}}>
                  {alertStatus==='confirmed' && `🔔 ${P.name} notified`}
                  {alertStatus==='pending' && `Flagging for ${P.name}…`}
                  {alertStatus==='recorded_email_disabled' && `Flagged — saved to dashboard (email not configured)`}
                  {alertStatus==='recorded_email_failed' && `Flagged — saved to dashboard (email couldn't be sent)`}
                  {alertStatus==='failed' && `⚠ Flagged — notification could not be sent`}
                </span>
                <button onClick={() => { setAlerts([]); setAlertStatus(null); }} style={{background:'none',border:'none',cursor:'pointer',color:'rgba(220,140,60,0.42)',fontSize:16,lineHeight:1,padding:0}}>×</button>
              </div>
              <div style={{display:'flex',flexWrap:'wrap',gap:4}}>{alerts.map((a,i) => <span key={i} style={{padding:'2px 8px',background:'rgba(220,140,60,0.09)',borderRadius:20,fontSize:10,fontFamily:"'Jost',sans-serif",color:'#e8a060'}}>{a}</span>)}</div>
            </div>
          </div>
        )}

        {err && <div style={{margin:'9px 0',padding:'8px 12px',background:'rgba(180,60,60,0.06)',border:'1px solid rgba(180,60,60,0.16)',borderRadius:20,fontSize:12,fontFamily:"'Jost',sans-serif",color:'#c08080'}}>{err}</div>}

        <div style={{paddingBottom:20,paddingTop:messages.length>0&&alerts.length===0?15:0}}>
          {messages.map((m,i) => (
            <div key={i} className="ma" style={{marginBottom:11,display:'flex',flexDirection:m.role==='user'?'row-reverse':'row',alignItems:'flex-end',gap:7}}>
              {m.role === 'assistant' && (
                <div style={{width:24,height:24,borderRadius:'50%',flexShrink:0,marginBottom:2,background:m._blocked?'rgba(180,60,60,0.26)':P.gradient,display:'flex',alignItems:'center',justifyContent:'center',fontSize:10}}>{m._blocked?'⚑':P.emoji}</div>
              )}
              <div style={{maxWidth:'75%',padding:'9px 13px',borderRadius:m.role==='user'?'13px 13px 4px 13px':'13px 13px 13px 4px',background:m.role==='user'?P.gradient:m._blocked?'rgba(180,60,60,0.06)':'rgba(255,255,255,0.03)',border:m.role==='user'?'none':m._blocked?'1px solid rgba(180,60,60,0.14)':'1px solid rgba(201,169,110,0.07)',color:m.role==='user'?'#0c0a08':m._blocked?'#b08080':'#e8dcc8',fontSize:15,lineHeight:1.72,fontFamily:m.role==='user'?"'Jost',sans-serif":"'Cormorant Garamond',serif",fontWeight:400,whiteSpace:'pre-wrap',fontStyle:m._blocked?'italic':'normal'}}>
                {m.content}
                {m._blocked && <div style={{marginTop:3,fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(180,120,120,0.4)',letterSpacing:'0.07em',textTransform:'uppercase'}}>— blocked by {P.name}</div>}
              </div>
            </div>
          ))}
          {loading && (
            <div className="ma" style={{display:'flex',alignItems:'flex-end',gap:7,marginBottom:11}}>
              <div style={{width:24,height:24,borderRadius:'50%',background:P.gradient,display:'flex',alignItems:'center',justifyContent:'center',fontSize:10}}>{P.emoji}</div>
              <div style={{padding:'9px 13px',borderRadius:'13px 13px 13px 4px',background:'rgba(255,255,255,0.03)',border:'1px solid rgba(201,169,110,0.07)'}}><Dots color={acc}/></div>
            </div>
          )}
          {turns >= MAX_TURNS && (
            <div style={{textAlign:'center',padding:'16px 0',fontStyle:'italic',color:'rgba(232,220,200,0.25)',fontSize:14}}>
              {TL.demoEnd} <span style={{color:acc,cursor:'pointer'}} onClick={onBack}>{TL.tryAnother}</span>
            </div>
          )}
          <div ref={bottomRef}/>
        </div>
      </div>

      {/* LEAD CAPTURE */}
      {showLeadForm && (
        <div style={{width:'100%',maxWidth:680,padding:'12px 15px',background:'rgba(201,169,110,0.06)',borderTop:'1px solid rgba(201,169,110,0.15)',animation:'fadeUp 0.4s ease forwards'}}>
          <div style={{fontSize:13,fontFamily:"'Cormorant Garamond',serif",fontStyle:'italic',color:'rgba(232,220,200,0.7)',marginBottom:10,lineHeight:1.5}}>
            {L==='it'?`Vuoi che ${P.name} ti ricontatti? Lascia i tuoi dati — opzionale.`:`Want ${P.name} to follow up with you? Leave your details — optional.`}
          </div>
          <div style={{display:'flex',gap:7,flexWrap:'wrap'}}>
            <input value={leadName} onChange={e=>setLeadName(e.target.value)} placeholder={L==='it'?'Il tuo nome':'Your name'} style={{flex:1,minWidth:120,padding:'7px 11px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif"}}/>
            <input value={leadEmail} onChange={e=>setLeadEmail(e.target.value)} placeholder='Email' style={{flex:2,minWidth:160,padding:'7px 11px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.2)',borderRadius:20,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif"}}/>
            <button onClick={submitLead} style={{padding:'7px 14px',background:'linear-gradient(135deg,#c9a96e,#7a4f0e)',border:'none',borderRadius:20,cursor:'pointer',color:'#0c0a08',fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.06em',whiteSpace:'nowrap'}}>
              {L==='it'?'Invia ✦':'Send ✦'}
            </button>
            <button onClick={()=>{setShowLeadForm(false);setLeadCaptured(true);}} style={{padding:'7px 10px',background:'none',border:'1px solid rgba(201,169,110,0.12)',borderRadius:20,cursor:'pointer',color:'rgba(201,169,110,0.4)',fontSize:11,fontFamily:"'Jost',sans-serif"}}>
              {L==='it'?'Salta':'Skip'}
            </button>
          </div>
        </div>
      )}

      {/* LEGAL FORMS TRIGGER */}
      {profileLegalForms.length > 0 && !legalFormDone && turns >= 1 && (
        <div style={{width:'100%',maxWidth:680,padding:'10px 15px',background:'rgba(201,169,110,0.05)',borderTop:'1px solid rgba(201,169,110,0.12)'}}>
          <button onClick={() => setShowLegalForm(true)} style={{width:'100%',padding:'10px 0',background:'rgba(201,169,110,0.08)',border:'1px solid rgba(201,169,110,0.22)',borderRadius:20,cursor:'pointer',color:acc,fontSize:11,fontFamily:"'Jost',sans-serif",fontWeight:600,letterSpacing:'0.06em',textTransform:'uppercase'}}>
            📋 {L==='it'?`Completa i moduli prenotazione (${profileLegalForms.length})`:`Complete booking forms (${profileLegalForms.length})`}
          </button>
        </div>
      )}

      {/* INPUT */}
      <div style={{width:'100%',maxWidth:680,padding:'9px 15px 13px',position:'sticky',bottom:0,background:'rgba(12,10,8,0.98)',backdropFilter:'blur(18px)',borderTop:'1px solid rgba(255,255,255,0.04)'}}>
        <div style={{display:'flex',alignItems:'flex-end',gap:7,background:'rgba(255,255,255,0.022)',border:`1px solid ${acc}20`,borderRadius:20,padding:'7px 7px 7px 13px'}}>
          <textarea ref={inputRef} value={input} onChange={e=>setInput(e.target.value)} onKeyDown={hk}
            placeholder={turns>=MAX_TURNS?TL.demoEnded:TL.askAnything}
            disabled={turns>=MAX_TURNS} rows={1}
            style={{flex:1,background:'transparent',border:'none',color:'#e8dcc8',fontSize:15,fontFamily:"'Jost',sans-serif",fontWeight:300,letterSpacing:'0.02em',lineHeight:1.5,maxHeight:88,overflowY:'auto',outline:'none',resize:'none'}}
            onInput={e=>{e.target.style.height='auto';e.target.style.height=Math.min(e.target.scrollHeight,88)+'px';}}/>
          <button onClick={() => send()} disabled={loading||!input.trim()||turns>=MAX_TURNS} style={{width:31,height:31,borderRadius:'50%',flexShrink:0,background:input.trim()&&!loading&&turns<MAX_TURNS?P.gradient:'rgba(201,169,110,0.08)',border:'none',cursor:input.trim()&&!loading&&turns<MAX_TURNS?'pointer':'default',color:input.trim()&&!loading&&turns<MAX_TURNS?'#0c0a08':'rgba(201,169,110,0.2)',fontSize:13,fontWeight:700,display:'flex',alignItems:'center',justifyContent:'center',transition:'all 0.2s'}}>↑</button>
        </div>
        <div style={{display:'flex',justifyContent:'space-between',marginTop:5,padding:'0 1px'}}>
          <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.12)',letterSpacing:'0.07em',textTransform:'uppercase'}}>{TL.powered}</div>
          <div style={{display:'flex',gap:8}}>
            <span style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(220,140,60,0.32)'}}>🔔 Alerts</span>
            <span style={{fontSize:9,fontFamily:"'Jost',sans-serif",color:'rgba(120,160,100,0.32)'}}>👁 Review</span>
          </div>
        </div>
      </div>

      {reviewData && <ReviewModal reply={reviewData.reply} name={P.name} acc={acc} onApprove={approve} onReject={reject}/>}
      {showLegalForm && <LegalFormModal forms={profileLegalForms} profile={P} profileId={profileId} lang={lang} onClose={() => setShowLegalForm(false)} onComplete={() => setLegalFormDone(true)}/>}
    </div>
  );
}
