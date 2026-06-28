'use client';
import { useState, useEffect } from 'react';
import { PA_PERSONALITIES } from '@/lib/constants';

function Toggle({ value, onChange }) {
  return (
    <div onClick={onChange} style={{ width:40, height:22, borderRadius:11, background:value?'#C9A96E':'rgba(255,255,255,0.1)', position:'relative', cursor:'pointer', transition:'background 0.2s', flexShrink:0 }}>
      <div style={{ position:'absolute', top:3, left:value?20:3, width:16, height:16, borderRadius:'50%', background:'white', transition:'left 0.2s', boxShadow:'0 1px 4px rgba(0,0,0,0.3)' }}/>
    </div>
  );
}

export default function BravePASettings({ slug, onConfigChange }) {
  const [config, setConfig] = useState({ paName:'Brave PA', personality:'professional', proactiveAlerts:true, morningBriefing:true, soundEnabled:true, soundStyle:'subtle' });
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    const s = localStorage.getItem(`brave_pa_config_${slug}`);
    if (s) try { setConfig(prev => ({...prev, ...JSON.parse(s)})); } catch(e) {}
  }, [slug]);

  const save = () => {
    localStorage.setItem(`brave_pa_config_${slug}`, JSON.stringify(config));
    onConfigChange && onConfigChange(config);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  const gold = '#C9A96E', cream = '#E8DCC8', dark = '#0C0A08';
  const sL = { fontSize:10, fontFamily:"'Jost',sans-serif", fontWeight:500, letterSpacing:'0.12em', textTransform:'uppercase', color:'rgba(201,169,110,0.45)', marginBottom:10, marginTop:20 };
  const card = { background:'rgba(255,255,255,0.025)', border:'1px solid rgba(201,169,110,0.1)', borderRadius:13, padding:'16px 18px', marginBottom:12 };

  return (
    <div>
      <div style={sL}>✦ Your PA — Brave PA</div>
      <div style={card}>
        <div style={{ fontSize:11, color:'rgba(201,169,110,0.5)', marginBottom:6, fontFamily:"'Jost',sans-serif" }}>PA Name</div>
        <input value={config.paName} onChange={e=>setConfig({...config,paName:e.target.value})} placeholder="Brave PA" style={{ width:'100%', padding:'8px 12px', background:'rgba(255,255,255,0.04)', border:'1px solid rgba(201,169,110,0.15)', borderRadius:8, color:cream, fontSize:13, fontFamily:"'Jost',sans-serif", outline:'none' }}/>
      </div>
      <div style={sL}>Personality</div>
      <div style={{ display:'flex', gap:8, flexWrap:'wrap', marginBottom:16 }}>
        {Object.entries(PA_PERSONALITIES).map(([key,p]) => (
          <button key={key} onClick={()=>setConfig({...config,personality:key})} title={p.description} style={{ padding:'7px 14px', borderRadius:8, border:'1px solid', cursor:'pointer', fontSize:11, fontFamily:"'Jost',sans-serif", background:config.personality===key?'rgba(201,169,110,0.15)':'transparent', borderColor:config.personality===key?'rgba(201,169,110,0.4)':'rgba(255,255,255,0.08)', color:config.personality===key?gold:'rgba(232,220,200,0.38)' }}>
            {p.label}
          </button>
        ))}
      </div>
      <div style={sL}>Notifications &amp; Sound</div>
      <div style={card}>
        {[['proactiveAlerts','Proactive PA messages'],['morningBriefing','Morning briefing at 7am'],['soundEnabled','Sound on']].map(([key,label]) => (
          <div key={key} style={{ display:'flex', justifyContent:'space-between', alignItems:'center', padding:'8px 0', borderBottom:'1px solid rgba(255,255,255,0.04)' }}>
            <span style={{ fontSize:13, color:cream, fontFamily:"'Jost',sans-serif" }}>{label}</span>
            <Toggle value={config[key]} onChange={()=>setConfig({...config,[key]:!config[key]})}/>
          </div>
        ))}
        {config.soundEnabled && (
          <div style={{ marginTop:12 }}>
            <div style={{ fontSize:11, color:'rgba(201,169,110,0.5)', marginBottom:8, fontFamily:"'Jost',sans-serif" }}>Sound style</div>
            <div style={{ display:'flex', gap:8 }}>
              {['subtle','standard','bold'].map(style => (
                <button key={style} onClick={()=>setConfig({...config,soundStyle:style})} style={{ padding:'5px 12px', borderRadius:7, border:'1px solid', cursor:'pointer', fontSize:11, fontFamily:"'Jost',sans-serif", textTransform:'capitalize', background:config.soundStyle===style?'rgba(201,169,110,0.15)':'transparent', borderColor:config.soundStyle===style?'rgba(201,169,110,0.4)':'rgba(255,255,255,0.08)', color:config.soundStyle===style?gold:'rgba(232,220,200,0.38)' }}>
                  {style}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
      <button onClick={save} style={{ padding:'10px 20px', background:`linear-gradient(135deg,${gold},#7a4f0e)`, border:'none', borderRadius:10, cursor:'pointer', color:dark, fontSize:13, fontFamily:"'Jost',sans-serif", fontWeight:600 }}>
        {saved ? '✓ Saved' : 'Save PA Settings'}
      </button>
    </div>
  );
}
