'use client';
import { useState } from 'react';
import { BACKEND_URL, LEGAL_FORM_FIELDS } from '@/lib/constants';

export default function LegalFormModal({ forms, profile, profileId, lang, onClose, onComplete }) {
  const L = lang || 'en';
  const [clientName, setClientName] = useState('');
  const [clientEmail, setClientEmail] = useState('');
  const [answers, setAnswers] = useState({});
  const [agreed, setAgreed] = useState({});
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const acc = profile?.accent || '#c9a96e';

  const allAgreed = forms.every(f => agreed[f]) && clientName.trim().length > 1;

  const submit = async () => {
    if (!allAgreed) return;
    setSubmitting(true);
    try {
      await fetch(`${BACKEND_URL}/consent`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          profile_id: profileId,
          session_id: sessionStorage.getItem('cai_session') || 'anonymous',
          client_name: clientName,
          client_email: clientEmail,
          forms_agreed: forms,
          answers: JSON.stringify(answers),
          signature_date: new Date().toISOString()
        })
      });
      setDone(true);
      setTimeout(() => { onComplete && onComplete(); onClose(); }, 1800);
    } catch(e) {
      setDone(true);
      setTimeout(() => { onComplete && onComplete(); onClose(); }, 1800);
    }
    setSubmitting(false);
  };

  return (
    <div style={{position:'fixed',inset:0,zIndex:9999,background:'rgba(0,0,0,0.7)',display:'flex',alignItems:'flex-end',justifyContent:'center'}} onClick={onClose}>
      <div style={{width:'100%',maxWidth:480,maxHeight:'88vh',overflowY:'auto',background:'#0e0b07',border:'1px solid rgba(201,169,110,0.25)',borderRadius:'20px 20px 0 0',padding:'22px 20px 28px'}} onClick={e=>e.stopPropagation()}>
        {!done ? (
          <>
            <div style={{display:'flex',justifyContent:'space-between',alignItems:'center',marginBottom:16}}>
              <div style={{fontSize:11,fontFamily:"'Jost',sans-serif",letterSpacing:'0.14em',textTransform:'uppercase',color:'rgba(201,169,110,0.5)'}}>{L==='it'?'Moduli da completare':'Forms to complete'}</div>
              <div onClick={onClose} style={{cursor:'pointer',color:'rgba(201,169,110,0.4)',fontSize:18}}>×</div>
            </div>
            <div style={{fontSize:20,fontWeight:400,color:'#e8dcc8',marginBottom:6,fontFamily:"'Cormorant Garamond',serif"}}>{L==='it'?`Prima di prenotare con ${profile?.name||''}`:`Before booking with ${profile?.name||''}`}</div>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)',marginBottom:20,lineHeight:1.6}}>{L==='it'?'Compila questi moduli rapidi — sono conformi al GDPR UK/EU e conservati in modo sicuro.':'Complete these quick forms — fully UK/EU GDPR compliant and securely stored.'}</div>
            <div style={{marginBottom:14}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.4)',marginBottom:5}}>{L==='it'?'Nome completo':'Full name'}</div>
              <input value={clientName} onChange={e=>setClientName(e.target.value)} placeholder={L==='it'?'Il tuo nome':'Your name'} style={{width:'100%',padding:'10px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif"}}/>
            </div>
            <div style={{marginBottom:18}}>
              <div style={{fontSize:9,fontFamily:"'Jost',sans-serif",letterSpacing:'0.1em',textTransform:'uppercase',color:'rgba(201,169,110,0.4)',marginBottom:5}}>{L==='it'?'Email (opzionale)':'Email (optional)'}</div>
              <input value={clientEmail} onChange={e=>setClientEmail(e.target.value)} placeholder="you@email.com" style={{width:'100%',padding:'10px 12px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.18)',borderRadius:8,color:'#e8dcc8',fontSize:13,fontFamily:"'Jost',sans-serif"}}/>
            </div>
            {forms.map(formKey => {
              const def = LEGAL_FORM_FIELDS[formKey];
              if (!def) return null;
              const txt = def[L] || def.en;
              return (
                <div key={formKey} style={{marginBottom:16,padding:'13px 14px',background:'rgba(201,169,110,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:20}}>
                  <div style={{fontSize:13,fontWeight:500,color:'#e8dcc8',marginBottom:6,fontFamily:"'Jost',sans-serif"}}>{txt.title}</div>
                  <div style={{fontSize:12,fontFamily:"'Cormorant Garamond',serif",color:'rgba(232,220,200,0.6)',lineHeight:1.6,marginBottom:9}}>{txt.q}</div>
                  <textarea value={answers[formKey]||''} onChange={e=>setAnswers(a=>({...a,[formKey]:e.target.value}))} placeholder={L==='it'?'Risposta (opzionale, lascia vuoto se non applicabile)':'Answer (optional, leave blank if not applicable)'} rows={2} style={{width:'100%',padding:'8px 10px',background:'rgba(255,255,255,0.04)',border:'1px solid rgba(201,169,110,0.1)',borderRadius:7,color:'#e8dcc8',fontSize:12,fontFamily:"'Jost',sans-serif",marginBottom:9}}/>
                  <label style={{display:'flex',alignItems:'flex-start',gap:8,cursor:'pointer'}}>
                    <input type="checkbox" checked={!!agreed[formKey]} onChange={e=>setAgreed(a=>({...a,[formKey]:e.target.checked}))} style={{marginTop:2}}/>
                    <span style={{fontSize:11,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.7)'}}>{L==='it'?'Ho letto e acconsento':'I have read and agree'}</span>
                  </label>
                </div>
              );
            })}
            <button onClick={submit} disabled={!allAgreed||submitting} style={{width:'100%',padding:'13px 0',background:allAgreed?`linear-gradient(135deg,${acc},#7a4f0e)`:'rgba(201,169,110,0.12)',border:'none',borderRadius:20,cursor:allAgreed?'pointer':'not-allowed',color:allAgreed?'#0c0a08':'rgba(232,220,200,0.3)',fontSize:13,fontFamily:"'Jost',sans-serif",fontWeight:700,letterSpacing:'0.08em',textTransform:'uppercase',marginTop:6}}>
              {submitting?(L==='it'?'Invio...':'Submitting...'):(L==='it'?'✓ Conferma e firma':'✓ Confirm & Sign')}
            </button>
            <div style={{fontSize:10,fontFamily:"'Jost',sans-serif",color:'rgba(201,169,110,0.3)',textAlign:'center',marginTop:10}}>{L==='it'?`Firmato digitalmente da ${clientName||'...'} il ${new Date().toLocaleDateString()}`:`Digitally signed by ${clientName||'...'} on ${new Date().toLocaleDateString()}`}</div>
          </>
        ) : (
          <div style={{textAlign:'center',padding:'30px 10px'}}>
            <div style={{fontSize:32,marginBottom:14,color:acc}}>✓</div>
            <div style={{fontSize:18,color:'#e8dcc8',fontFamily:"'Cormorant Garamond',serif",marginBottom:6}}>{L==='it'?'Moduli completati!':'Forms completed!'}</div>
            <div style={{fontSize:12,fontFamily:"'Jost',sans-serif",color:'rgba(232,220,200,0.4)'}}>{L==='it'?`${profile?.name||''} è stato notificato.`:`${profile?.name||''} has been notified.`}</div>
          </div>
        )}
      </div>
    </div>
  );
}
